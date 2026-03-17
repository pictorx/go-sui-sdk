package v2

// builder_cgo.go
//
// CGo bindings for the transaction_builder static library.
//
// This file is a drop-in replacement for builder.go (the WASM version).
// The public API is identical — only the build tag and import differ.
//
// Build tags
// ----------
// WASM build (original):
//   go build                          (no tag — builder.go is active)
//
// Static-lib build (this file):
//   go build -tags txbuilder_cgo      (builder_cgo.go is active instead)
//
// Both files live in the same package; exactly one is compiled at a time.
//
// Prerequisites for the static-lib build
// ---------------------------------------
//  1. Build the Rust library for your native target:
//       cargo build --release
//     This produces  target/release/libtransaction_builder.a  (Linux/macOS)
//     or             target/release/transaction_builder.lib   (Windows).
//
//  2. Set the CGo search paths, either via environment variables:
//       export CGO_LDFLAGS="-L/path/to/target/release"
//       export CGO_CFLAGS="-I/path/to/transaction_builder_header"
//     or by editing the #cgo directives below (see "Adjust paths" comment).
//
//  3. Build:
//       go build -tags txbuilder_cgo ./...
//
// Linking notes
// -------------
//  • Linux  : the linker flags below pull in -lpthread -ldl -lm (all
//             standard; no extra packages needed).
//  • macOS  : -framework Security is added for the TLS stack used by
//             sui-rpc inside the static lib.
//  • Windows: replace -ltransaction_builder with transaction_builder.lib
//             and ensure the .lib is on the LIB search path.

/*
// ── Adjust paths if you prefer not to use env vars ───────────────────────────
#cgo CFLAGS:  -I${SRCDIR}/../transaction
#cgo LDFLAGS: -L${SRCDIR}/../transaction/target/release
// ─────────────────────────────────────────────────────────────────────────────

#cgo linux   LDFLAGS: -L${SRCDIR}/../transaction/target/release -Wl,-Bstatic -ltransaction_builder -Wl,-Bdynamic -lpthread -ldl -lm -lresolv
#cgo darwin  LDFLAGS: -L${SRCDIR}/../transaction/target/release -ltransaction_builder -lpthread -ldl -lm -framework Security
#cgo windows LDFLAGS: -L${SRCDIR}/../transaction/target/release -ltransaction_builder -lws2_32 -luserenv -lntdll -lbcrypt

#include "../transaction/transaction_builder.h"
#include <stdlib.h>   // free(), malloc()
#include <string.h>   // memcpy()
*/
import "C"

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"unsafe"
)

// ── Builder ───────────────────────────────────────────────────────────────────

// Builder wraps the native (CGo) TransactionBuilder pointer.
// It is NOT safe for concurrent use.
type Builder struct {
	ptr *C.TransactionBuilder // nil after Build() or Free()
}

// NewBuilder instantiates a fresh TransactionBuilder inside the static library.
// The ctx parameter is accepted for API compatibility with the WASM version
// but is not used — native calls are synchronous.
func NewBuilder(_ interface{}, _ interface{}) *Builder {
	return &Builder{ptr: C.new_builder()}
}

// Free releases a builder that was NOT consumed by Build().
// After a successful Build() the builder is already freed — do not call Free()
// in that case.  Safe to call on a nil/already-freed builder.
func (b *Builder) Free() {
	if b.ptr != nil {
		C.free_builder(b.ptr)
		b.ptr = nil
	}
}

// ── Configuration ─────────────────────────────────────────────────────────────

// SetConfig sets the sender address, gas budget, and gas price.
// sender must be a 0x-prefixed 32-byte hex string.
func (b *Builder) SetConfig(sender string, gasBudget, gasPrice uint64) error {
	payload, _ := json.Marshal(map[string]any{
		"sender":     sender,
		"gas_budget": gasBudget,
		"gas_price":  gasPrice,
	})
	cptr, clen := goBytesCopy(payload)
	defer C.free(unsafe.Pointer(cptr))
	code := C.set_config(b.ptr, (*C.uint8_t)(cptr), C.size_t(clen))
	if code != 1 {
		return fmt.Errorf("set_config failed (code %d) — check sender address format", code)
	}
	return nil
}

// ── Gas objects ───────────────────────────────────────────────────────────────

// AddGasObject adds an owned gas coin identified by its object ID, version,
// and base-58 digest string.
func (b *Builder) AddGasObject(id string, version uint64, digest string) error {
	payload, _ := json.Marshal(map[string]any{
		"id":      id,
		"version": version,
		"digest":  digest,
	})
	cptr, clen := goBytesCopy(payload)
	defer C.free(unsafe.Pointer(cptr))
	code := C.add_gas_object(b.ptr, (*C.uint8_t)(cptr), C.size_t(clen))
	switch code {
	case 1:
		return nil
	case -2:
		return fmt.Errorf("add_gas_object: invalid digest %q", digest)
	default:
		return fmt.Errorf("add_gas_object failed (code %d)", code)
	}
}

// ── Gas pseudo-input ──────────────────────────────────────────────────────────

// GasArgument returns the Argument ID for the transaction's gas coin.
// Idempotent — always returns the same ID within one builder.
func (b *Builder) GasArgument() uint64 {
	return uint64(C.gas_argument(b.ptr))
}

// ── Object inputs ─────────────────────────────────────────────────────────────

// ObjectKind describes how an object is used as an input.
type ObjectKind string

const (
	ObjectKindOwned     ObjectKind = "owned"
	ObjectKindImmutable ObjectKind = "immutable"
	ObjectKindReceiving ObjectKind = "receiving"
	ObjectKindShared    ObjectKind = "shared"
)

// InputObject pushes an object input and returns its Argument ID.
//
// For owned / immutable / receiving: supply id, version, digest, kind.
// For shared: supply id, version, mutable, kind="shared" (digest is ignored).
func (b *Builder) InputObject(id string, version uint64, digest string, kind ObjectKind, mutable bool) (uint64, error) {
	m := map[string]any{
		"id":      id,
		"version": version,
		"kind":    string(kind),
	}
	if kind == ObjectKindShared {
		m["mutable"] = mutable
	} else {
		m["digest"] = digest
	}
	payload, _ := json.Marshal(m)
	cptr, clen := goBytesCopy(payload)
	defer C.free(unsafe.Pointer(cptr))
	res := int64(C.input_object(b.ptr, (*C.uint8_t)(cptr), C.size_t(clen)))
	if res < 0 {
		return 0, fmt.Errorf("input_object failed (code %d)", res)
	}
	return uint64(res), nil
}

// ── Pure-value helpers ────────────────────────────────────────────────────────

// PureBool pushes a BCS-encoded bool and returns its Argument ID.
func (b *Builder) PureBool(v bool) uint64 {
	var u C.uint8_t
	if v {
		u = 1
	}
	return uint64(C.pure_bool(b.ptr, u))
}

// PureU8 pushes a BCS-encoded u8 and returns its Argument ID.
func (b *Builder) PureU8(v uint8) uint64 {
	return uint64(C.pure_u8(b.ptr, C.uint8_t(v)))
}

// PureU16 pushes a BCS-encoded u16 and returns its Argument ID.
func (b *Builder) PureU16(v uint16) uint64 {
	return uint64(C.pure_u16(b.ptr, C.uint16_t(v)))
}

// PureU32 pushes a BCS-encoded u32 and returns its Argument ID.
func (b *Builder) PureU32(v uint32) uint64 {
	return uint64(C.pure_u32(b.ptr, C.uint32_t(v)))
}

// PureU64 pushes a BCS-encoded u64 and returns its Argument ID.
func (b *Builder) PureU64(v uint64) uint64 {
	return uint64(C.pure_u64(b.ptr, C.uint64_t(v)))
}

// PureU128 pushes a BCS-encoded u128 (supplied as high/low uint64 halves)
// and returns its Argument ID.
func (b *Builder) PureU128(hi, lo uint64) uint64 {
	// CGo signature: pure_u128(builder, lo, hi) — lo first, matching Rust.
	return uint64(C.pure_u128(b.ptr, C.uint64_t(lo), C.uint64_t(hi)))
}

// PureAddress pushes a BCS-encoded Sui address (bare 0x-prefixed hex string)
// and returns its Argument ID.
func (b *Builder) PureAddress(addr string) (uint64, error) {
	cptr, clen := goBytesCopy([]byte(addr))
	defer C.free(unsafe.Pointer(cptr))
	res := int64(C.pure_address(b.ptr, (*C.uint8_t)(cptr), C.size_t(clen)))
	if res < 0 {
		return 0, fmt.Errorf("pure_address: invalid address %q", addr)
	}
	return uint64(res), nil
}

// PureRawBCS pushes already-BCS-encoded bytes as a pure argument and returns
// its Argument ID.
func (b *Builder) PureRawBCS(bcsBytes []byte) uint64 {
	cptr, clen := goBytesCopy(bcsBytes)
	defer C.free(unsafe.Pointer(cptr))
	return uint64(C.pure_raw_bcs(b.ptr, (*C.uint8_t)(cptr), C.size_t(clen)))
}

// ── Nested result ─────────────────────────────────────────────────────────────

// NestedResult returns the Argument ID for the Nth sub-result of a
// multi-output command (e.g. the Kth coin from SplitCoins).
func (b *Builder) NestedResult(baseID, subIndex uint64) uint64 {
	return uint64(C.nested_result(b.ptr, C.uint64_t(baseID), C.uint64_t(subIndex)))
}

// ── Commands ──────────────────────────────────────────────────────────────────

// MoveCallArg describes a single argument to a Move call.
// Supply exactly one of ArgID (existing Argument) or PureBCS (raw bytes).
type MoveCallArg struct {
	ArgID   *uint64
	PureBCS []byte
}

// ArgID is a convenience constructor for a MoveCallArg that references an
// existing Argument by ID.
func ArgID(id uint64) MoveCallArg { return MoveCallArg{ArgID: &id} }

// ArgBCS is a convenience constructor for a MoveCallArg that passes raw
// pre-encoded BCS bytes.
func ArgBCS(bcs []byte) MoveCallArg { return MoveCallArg{PureBCS: bcs} }

// MoveCall executes an entry or public Move function and returns the result
// Argument ID.
func (b *Builder) MoveCall(pkg, module, function string, typeArgs []string, args []MoveCallArg) (uint64, error) {
	type callArgJSON struct {
		ID      *uint64 `json:"id,omitempty"`
		PureBCS []byte  `json:"pure_bcs,omitempty"`
	}
	jsonArgs := make([]callArgJSON, len(args))
	for i, a := range args {
		if a.ArgID != nil {
			jsonArgs[i] = callArgJSON{ID: a.ArgID}
		} else {
			jsonArgs[i] = callArgJSON{PureBCS: a.PureBCS}
		}
	}
	payload, _ := json.Marshal(map[string]any{
		"package":   pkg,
		"module":    module,
		"function":  function,
		"type_args": typeArgs,
		"arguments": jsonArgs,
	})
	cptr, clen := goBytesCopy(payload)
	defer C.free(unsafe.Pointer(cptr))
	res := int64(C.command_move_call(b.ptr, (*C.uint8_t)(cptr), C.size_t(clen)))
	if res < 0 {
		return 0, fmt.Errorf("command_move_call failed (code %d)", res)
	}
	return uint64(res), nil
}

// SplitCoins splits coinArgID into len(amountArgIDs) new coins.
// amountArgIDs must be Argument IDs returned by PureU64.
// Returns the base Argument ID; use NestedResult(base, i) to address coin i.
func (b *Builder) SplitCoins(coinArgID uint64, amountArgIDs []uint64) (uint64, error) {
	if len(amountArgIDs) == 0 {
		return 0, fmt.Errorf("SplitCoins: at least one amount required")
	}
	cAmounts, cCount := goU64SliceCopy(amountArgIDs)
	defer C.free(unsafe.Pointer(cAmounts))
	res := int64(C.command_split_coins(b.ptr, C.uint64_t(coinArgID), cAmounts, cCount))
	if res < 0 {
		return 0, fmt.Errorf("command_split_coins failed (code %d)", res)
	}
	return uint64(res), nil
}

// MergeCoins merges sourceArgIDs into targetCoinArgID.
// Produces no result; the target coin absorbs all sources.
func (b *Builder) MergeCoins(targetCoinArgID uint64, sourceArgIDs []uint64) error {
	if len(sourceArgIDs) == 0 {
		return fmt.Errorf("MergeCoins: at least one source required")
	}
	cSrcs, cCount := goU64SliceCopy(sourceArgIDs)
	defer C.free(unsafe.Pointer(cSrcs))
	code := int32(C.command_merge_coins(b.ptr, C.uint64_t(targetCoinArgID), cSrcs, cCount))
	if code != 1 {
		return fmt.Errorf("command_merge_coins failed (code %d)", code)
	}
	return nil
}

// TransferObjects sends objectArgIDs to the address identified by recipientArgID.
// recipientArgID must be an Argument ID returned by PureAddress.
func (b *Builder) TransferObjects(objectArgIDs []uint64, recipientArgID uint64) error {
	if len(objectArgIDs) == 0 {
		return fmt.Errorf("TransferObjects: at least one object required")
	}
	cObjs, cCount := goU64SliceCopy(objectArgIDs)
	defer C.free(unsafe.Pointer(cObjs))
	code := int32(C.command_transfer_objects(b.ptr, cObjs, cCount, C.uint64_t(recipientArgID)))
	if code != 1 {
		return fmt.Errorf("command_transfer_objects failed (code %d)", code)
	}
	return nil
}

// MakeMoveVec constructs a Move vector<T> from elemArgIDs.
// typeTag is the element type string (e.g. "0x2::sui::SUI"); pass "" to infer.
// Returns the result Argument ID.
func (b *Builder) MakeMoveVec(typeTag string, elemArgIDs []uint64) (uint64, error) {
	var ttPtr *C.uint8_t
	var ttLen C.size_t
	if typeTag != "" {
		p, l := goBytesCopy([]byte(typeTag))
		defer C.free(unsafe.Pointer(p))
		ttPtr = (*C.uint8_t)(p)
		ttLen = C.size_t(l)
	}

	var elemsPtr *C.uint64_t
	var elemsCount C.size_t
	if len(elemArgIDs) > 0 {
		p, c := goU64SliceCopy(elemArgIDs)
		defer C.free(unsafe.Pointer(p))
		elemsPtr = p
		elemsCount = c
	}

	res := int64(C.command_make_move_vec(b.ptr, ttPtr, ttLen, elemsPtr, elemsCount))
	if res < 0 {
		return 0, fmt.Errorf("command_make_move_vec failed (code %d)", res)
	}
	return uint64(res), nil
}

// Publish publishes a new Move package.
// modules is a slice of compiled module bytecodes.
// dependencies is a slice of 0x-prefixed package IDs this package depends on.
// Returns the UpgradeCap Argument ID.
func (b *Builder) Publish(modules [][]byte, dependencies []string) (uint64, error) {
	payload, _ := json.Marshal(map[string]any{
		"modules":      modules,
		"dependencies": dependencies,
	})
	cptr, clen := goBytesCopy(payload)
	defer C.free(unsafe.Pointer(cptr))
	res := int64(C.command_publish(b.ptr, (*C.uint8_t)(cptr), C.size_t(clen)))
	if res < 0 {
		return 0, fmt.Errorf("command_publish failed (code %d)", res)
	}
	return uint64(res), nil
}

// Upgrade upgrades an existing Move package.
// packageID is the on-chain ID of the package being upgraded.
// ticketArgID is the Argument ID of the UpgradeTicket from authorize_upgrade.
// Returns the UpgradeReceipt Argument ID.
func (b *Builder) Upgrade(modules [][]byte, dependencies []string, packageID string, ticketArgID uint64) (uint64, error) {
	payload, _ := json.Marshal(map[string]any{
		"modules":       modules,
		"dependencies":  dependencies,
		"package":       packageID,
		"ticket_arg_id": ticketArgID,
	})
	cptr, clen := goBytesCopy(payload)
	defer C.free(unsafe.Pointer(cptr))
	res := int64(C.command_upgrade(b.ptr, (*C.uint8_t)(cptr), C.size_t(clen)))
	if res < 0 {
		return 0, fmt.Errorf("command_upgrade failed (code %d)", res)
	}
	return uint64(res), nil
}

// ── Finalisation ─────────────────────────────────────────────────────────────

// Build serialises the transaction to BCS bytes and returns them.
// The builder is consumed — do NOT call Free() after a successful Build().
// Returns an error if any required field is missing or if build fails.
func (b *Builder) Build() ([]byte, error) {
	// build_transaction consumes the builder regardless of outcome.
	raw := C.build_transaction(b.ptr)
	b.ptr = nil

	if raw == nil {
		return nil, fmt.Errorf("build_transaction failed — ensure sender, gas object, gas_budget, gas_price and at least one command are set")
	}

	// Buffer layout: [4-byte LE uint32 payload_len][payload_len bytes BCS]
	rawSlice := unsafe.Slice((*byte)(unsafe.Pointer(raw)), 4)
	payloadLen := binary.LittleEndian.Uint32(rawSlice)

	bcsData := C.GoBytes(unsafe.Pointer(uintptr(unsafe.Pointer(raw))+4), C.int(payloadLen))

	// Free the Rust-owned buffer before returning.
	C.free_bytes(raw, C.size_t(payloadLen))

	return bcsData, nil
}

// ── internal CGo helpers ──────────────────────────────────────────────────────

// goBytesCopy copies a Go []byte into a C.malloc buffer.
// The caller must C.free the returned pointer.
func goBytesCopy(data []byte) (unsafe.Pointer, int) {
	if len(data) == 0 {
		return nil, 0
	}
	ptr := C.malloc(C.size_t(len(data)))
	C.memcpy(ptr, unsafe.Pointer(&data[0]), C.size_t(len(data)))
	return ptr, len(data)
}

// goU64SliceCopy copies a Go []uint64 into a C.malloc buffer of C.uint64_t.
// The caller must C.free the returned pointer.
func goU64SliceCopy(ids []uint64) (*C.uint64_t, C.size_t) {
	if len(ids) == 0 {
		return nil, 0
	}
	byteLen := len(ids) * 8
	ptr := C.malloc(C.size_t(byteLen))
	dst := unsafe.Slice((*byte)(ptr), byteLen)
	for i, v := range ids {
		binary.LittleEndian.PutUint64(dst[i*8:], v)
	}
	return (*C.uint64_t)(ptr), C.size_t(len(ids))
}
