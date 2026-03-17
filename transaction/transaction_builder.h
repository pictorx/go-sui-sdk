/**
 * transaction_builder.h
 *
 * C/C++ interface for the transaction-builder Rust library.
 *
 * Linking
 * -------
 * Native static lib (recommended — no WASM runtime needed):
 *   cc -I. your_file.c libtransaction_builder.a -lpthread -ldl -lm -o your_app
 *   (add -lresolv on Linux; nothing extra on macOS)
 *
 * WASM (original path):
 *   Load transaction_builder.wasm through your runtime as before.
 *
 * Both targets expose exactly this ABI; only the link step changes.
 *
 * Memory contract
 * ---------------
 *  - Buffers you pass IN  are owned by you; free them whenever you like.
 *  - `build_transaction` returns a buffer you must free with `free_bytes`.
 *  - Use `alloc` / `dealloc` only when you need the Rust allocator to own
 *    a staging buffer (e.g. when bridging from an environment that can only
 *    write into Rust-side memory).
 *
 * Error codes
 * -----------
 * Functions that return int32_t / int64_t use:
 *   >= 0   success  (value is an Argument ID or a boolean 1)
 *     -1   JSON/UTF-8 parse error  (or "count == 0" sentinel)
 *     -2   bad digest / bad type-tag
 *     -3   unknown object kind / bad function identifier
 *
 * Copyright (c) Mysten Labs, Inc.  SPDX-License-Identifier: Apache-2.0
 */

#ifndef TRANSACTION_BUILDER_H
#define TRANSACTION_BUILDER_H

#include <stddef.h>   /* size_t          */
#include <stdint.h>   /* uint8_t … int64_t */

#ifdef __cplusplus
extern "C" {
#endif

/* ── Opaque handle ───────────────────────────────────────────────────────── */

/**
 * Opaque pointer to a TransactionBuilder.
 * Create with new_builder(), destroy with free_builder() or (implicitly) with
 * a successful build_transaction() call.
 */
typedef struct TransactionBuilder TransactionBuilder;


/* ── Memory management ───────────────────────────────────────────────────── */

/**
 * alloc(len)
 * Allocate `len` bytes inside the Rust heap and return a pointer to them.
 * Use this when your runtime needs to write data into Rust-managed memory.
 * Must be freed with dealloc(ptr, len).
 */
uint8_t *alloc(size_t len);

/**
 * dealloc(ptr, len)
 * Free a buffer previously allocated by alloc().
 * `len` must equal the value passed to the matching alloc() call.
 */
void dealloc(uint8_t *ptr, size_t len);


/* ── Builder lifecycle ───────────────────────────────────────────────────── */

/**
 * new_builder()
 * Create a new, empty TransactionBuilder.
 * Returns an opaque pointer; never NULL.
 */
TransactionBuilder *new_builder(void);

/**
 * free_builder(builder)
 * Destroy a builder that was NOT consumed by build_transaction().
 * Safe to call with NULL.
 * Do NOT call this after a successful build_transaction() call — the builder
 * is already gone.
 */
void free_builder(TransactionBuilder *builder);


/* ── Configuration ───────────────────────────────────────────────────────── */

/**
 * set_config(builder, json_ptr, json_len)
 * Set sender address, optional gas_budget, and optional gas_price from JSON.
 *
 * JSON shape:
 *   {"sender":"0x<hex>","gas_budget":10000000,"gas_price":1000}
 *
 * Returns  1 on success, -1 on JSON parse error.
 */
int32_t set_config(TransactionBuilder *builder,
                   const uint8_t     *json_ptr,
                   size_t             json_len);


/* ── Gas objects ─────────────────────────────────────────────────────────── */

/**
 * add_gas_object(builder, json_ptr, json_len)
 * Add an owned gas object from JSON.
 *
 * JSON shape:
 *   {"id":"0x<hex>","version":2,"digest":"<base58>"}
 *
 * Returns  1 on success, -1 on JSON parse error, -2 on invalid digest.
 */
int32_t add_gas_object(TransactionBuilder *builder,
                       const uint8_t      *json_ptr,
                       size_t              json_len);

/**
 * gas_argument(builder)
 * Register (or retrieve) the gas-coin pseudo-input.
 * Idempotent — always returns the same Argument ID within one builder.
 * Pass the returned ID wherever a gas-coin argument is expected.
 */
uint64_t gas_argument(TransactionBuilder *builder);


/* ── Object inputs ───────────────────────────────────────────────────────── */

/**
 * input_object(builder, json_ptr, json_len)
 * Push an object input (owned / immutable / receiving / shared).
 *
 * JSON shapes:
 *   Owned / immutable / receiving:
 *     {"id":"0x<hex>","version":N,"digest":"<base58>","kind":"owned"}
 *   Shared:
 *     {"id":"0x<hex>","version":N,"mutable":true,"kind":"shared"}
 *
 * Returns Argument ID (>= 0) on success.
 *   -1  JSON parse error
 *   -2  bad or missing digest
 *   -3  unknown kind string
 */
int64_t input_object(TransactionBuilder *builder,
                     const uint8_t      *json_ptr,
                     size_t              json_len);


/* ── Pure-value helpers ──────────────────────────────────────────────────── */

/** Push a BCS-encoded bool. Returns Argument ID. value=0 → false, else true. */
int64_t pure_bool(TransactionBuilder *builder, uint8_t value);

/** Push a BCS-encoded u8. Returns Argument ID. */
int64_t pure_u8(TransactionBuilder *builder, uint8_t value);

/** Push a BCS-encoded u16. Returns Argument ID. */
int64_t pure_u16(TransactionBuilder *builder, uint16_t value);

/** Push a BCS-encoded u32. Returns Argument ID. */
int64_t pure_u32(TransactionBuilder *builder, uint32_t value);

/** Push a BCS-encoded u64. Returns Argument ID. */
int64_t pure_u64(TransactionBuilder *builder, uint64_t value);

/**
 * pure_u128(builder, lo, hi)
 * Push a BCS-encoded u128 supplied as two u64 halves.
 *   value = (hi << 64) | lo
 * Returns Argument ID.
 */
int64_t pure_u128(TransactionBuilder *builder, uint64_t lo, uint64_t hi);

/**
 * pure_address(builder, ptr, len)
 * Push a BCS-encoded Sui address.
 * `ptr` points to a UTF-8 hex string (e.g. "0xabc…") — no JSON quotes.
 * Returns Argument ID on success, -1 on parse error.
 */
int64_t pure_address(TransactionBuilder *builder,
                     const uint8_t      *ptr,
                     size_t              len);

/**
 * pure_raw_bcs(builder, ptr, len)
 * Push raw pre-BCS-encoded bytes as a pure argument.
 * Use this when you have already serialised the value yourself.
 * Returns Argument ID.
 */
int64_t pure_raw_bcs(TransactionBuilder *builder,
                     const uint8_t      *ptr,
                     size_t              len);


/* ── Nested result helper ─────────────────────────────────────────────────── */

/**
 * nested_result(builder, base_id, sub_index)
 * Mint a new Argument ID addressing the `sub_index`-th sub-result of a
 * multi-output command (e.g. the k-th coin from command_split_coins).
 *
 *   base_id   – Argument ID returned by command_split_coins (or similar).
 *   sub_index – 0-based index into the result tuple.
 *
 * Returns the new Argument ID.
 * Example: to access the 2nd split coin — nested_result(builder, base, 1).
 */
int64_t nested_result(TransactionBuilder *builder,
                      uint64_t            base_id,
                      uint64_t            sub_index);


/* ── Commands ────────────────────────────────────────────────────────────── */

/**
 * command_move_call(builder, json_ptr, json_len)
 * Generic Move function call.
 *
 * JSON shape:
 * {
 *   "package":   "0x2",
 *   "module":    "coin",
 *   "function":  "split",
 *   "type_args": ["0x2::sui::SUI"],          // optional
 *   "arguments": [                            // optional
 *     {"id": 3},                              // reference an existing Argument
 *     {"pure_bcs": [1,0,0,0,0,0,0,0]}        // raw BCS bytes inline
 *   ]
 * }
 *
 * Returns result Argument ID (>= 0) on success.
 *   -1  JSON parse error
 *   -2  invalid module identifier
 *   -3  invalid function identifier
 */
int64_t command_move_call(TransactionBuilder *builder,
                          const uint8_t      *json_ptr,
                          size_t              json_len);

/**
 * command_split_coins(builder, coin_arg_id, amount_arg_ids_ptr, count)
 * Split one coin into `count` coins of specified amounts.
 *
 *   coin_arg_id        – Argument ID of the source coin (use gas_argument()
 *                        for the gas coin).
 *   amount_arg_ids_ptr – C array of `count` Argument IDs, each from pure_u64().
 *   count              – number of output coins (must be >= 1).
 *
 * Returns the BASE Argument ID shared by all result coins.
 * Address individual results with nested_result(base, 0..count-1).
 * Returns -1 if count == 0.
 */
int64_t command_split_coins(TransactionBuilder *builder,
                            uint64_t            coin_arg_id,
                            const uint64_t     *amount_arg_ids_ptr,
                            size_t              count);

/**
 * command_merge_coins(builder, target_coin_arg_id, source_arg_ids_ptr, count)
 * Merge `count` source coins into `target_coin_arg_id` (in-place, no result).
 *
 *   target_coin_arg_id – Argument ID of the coin to merge into.
 *   source_arg_ids_ptr – C array of `count` Argument IDs to merge.
 *   count              – number of source coins (must be >= 1).
 *
 * Returns 1 on success, -1 if count == 0.
 */
int32_t command_merge_coins(TransactionBuilder *builder,
                            uint64_t            target_coin_arg_id,
                            const uint64_t     *source_arg_ids_ptr,
                            size_t              count);

/**
 * command_transfer_objects(builder, object_arg_ids_ptr, count, recipient_arg_id)
 * Transfer a list of objects to a recipient.
 *
 *   object_arg_ids_ptr – C array of `count` Argument IDs.
 *   count              – number of objects (must be >= 1).
 *   recipient_arg_id   – Argument ID from pure_address().
 *
 * Returns 1 on success, -1 if count == 0.
 */
int32_t command_transfer_objects(TransactionBuilder *builder,
                                 const uint64_t     *object_arg_ids_ptr,
                                 size_t              count,
                                 uint64_t            recipient_arg_id);

/**
 * command_make_move_vec(builder,
 *                       type_tag_ptr, type_tag_len,
 *                       elem_arg_ids_ptr, count)
 * Construct a Move vector<T> from a list of existing arguments.
 *
 *   type_tag_ptr / type_tag_len – UTF-8 type-tag string e.g. "0x2::sui::SUI".
 *                                 Pass ptr=NULL / len=0 to infer from elements.
 *   elem_arg_ids_ptr            – C array of `count` Argument IDs (may be
 *                                 NULL when count == 0).
 *   count                       – number of elements.
 *
 * Returns result Argument ID (>= 0) on success.
 *   -1  bad type-tag UTF-8
 *   -2  type-tag parse error
 */
int64_t command_make_move_vec(TransactionBuilder *builder,
                              const uint8_t      *type_tag_ptr,
                              size_t              type_tag_len,
                              const uint64_t     *elem_arg_ids_ptr,
                              size_t              count);

/**
 * command_publish(builder, json_ptr, json_len)
 * Publish new Move modules.
 *
 * JSON shape:
 * {
 *   "modules":      [[...bytecode bytes...], ...],
 *   "dependencies": ["0x1", "0x2", ...]
 * }
 *
 * Returns the UpgradeCap Argument ID (>= 0), or -1 on JSON parse error.
 */
int64_t command_publish(TransactionBuilder *builder,
                        const uint8_t      *json_ptr,
                        size_t              json_len);

/**
 * command_upgrade(builder, json_ptr, json_len)
 * Upgrade an existing Move package.
 *
 * JSON shape:
 * {
 *   "modules":       [[...bytecode bytes...], ...],
 *   "dependencies":  ["0x1", "0x2"],
 *   "package":       "0xCAFE...",
 *   "ticket_arg_id": 7
 * }
 * `ticket_arg_id` must point to the UpgradeTicket from
 * 0x2::package::authorize_upgrade.
 *
 * Returns the UpgradeReceipt Argument ID (>= 0), or -1 on JSON parse error.
 */
int64_t command_upgrade(TransactionBuilder *builder,
                        const uint8_t      *json_ptr,
                        size_t              json_len);


/* ── Finalisation ────────────────────────────────────────────────────────── */

/**
 * build_transaction(builder)
 * Serialise the fully-built transaction to BCS and return a heap buffer.
 *
 * Buffer layout:
 *   [4 bytes: uint32_t payload_len, little-endian]
 *   [payload_len bytes: BCS-encoded transaction data]
 *
 * To consume:
 *   1. Read the first 4 bytes as a little-endian uint32_t → payload_len.
 *   2. Read the next payload_len bytes → your BCS payload.
 *   3. Call free_bytes(ptr, payload_len) to release the buffer.
 *
 * Returns NULL on any build or serialisation error.
 *
 * IMPORTANT: This call CONSUMES (drops) the builder.
 *            Do NOT call free_builder() after a successful call.
 *            On NULL return the builder is still valid and must be freed with
 *            free_builder().
 */
uint8_t *build_transaction(TransactionBuilder *builder);

/**
 * free_bytes(ptr, payload_len)
 * Free the buffer returned by build_transaction().
 * `payload_len` is the uint32_t read from the first 4 bytes of the buffer
 * (NOT the total buffer size — the library adds 4 internally).
 */
void free_bytes(uint8_t *ptr, size_t payload_len);


#ifdef __cplusplus
} /* extern "C" */
#endif

#endif /* TRANSACTION_BUILDER_H */