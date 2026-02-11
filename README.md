# go-sui-sdk
This is a golang sdk implementing gRPC

## cmd
protoc --go_out=paths=import:./sui_rpc_proto --go-grpc_out=paths=import:./sui_rpc_pr
oto sui_rpc_proto/*.proto

cargo build --release --target wasm32-wasip1 --no-default-features

## dependencies
https://github.com/MystenLabs/sui-rust-sdk/tree/9b29d6040c3409de996d8b50d95961d9a660f14b/crates/sui-rpc/vendored/proto
https://github.com/MystenLabs/sui-rust-sdk/tree/master/crates/sui-transaction-builder
