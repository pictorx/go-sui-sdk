package main

import (
	"context"
	"fmt"
	"log"
	"os"

	gosuisdk "github.com/pictorx/go-sui-sdk"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	ctx := context.Background()

	// ── WASM runtime ─────────────────────────────────────────────────────
	rt := wazero.NewRuntime(ctx)
	defer rt.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, rt)

	wasmBytes, err := os.ReadFile("../transaction/target/wasm32-wasip1/release/transaction_builder.wasm")
	if err != nil {
		panic(err)
	}
	mod, err := rt.Instantiate(ctx, wasmBytes)
	if err != nil {
		panic(err)
	}

	// ── gRPC connection ───────────────────────────────────────────────────
	conn, err := grpc.Dial(
		"fullnode.testnet.sui.io:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		log.Fatalf("grpc.Dial: %v", err)
	}
	defer conn.Close()

	// ── Constants ────────────────────────────────────────────────────────
	const (
		sender    = "0x8aeec8403b86f22e58d87bdc85ff78c87b69dce58b8f651900b9eb5644f45180"
		recipient = sender
		splitMIST = uint64(100_000_000) // 0.1 SUI in MIST
	)

	// ── Fetch gas coin from chain ─────────────────────────────────────────
	ownedObjs, err := gosuisdk.ListOwnedObjects(conn, sender, nil, nil, ctx)
	if err != nil {
		panic(err)
	}
	suiCoins := gosuisdk.OwnedCoins(ownedObjs, gosuisdk.SuiCoin.String(), sender)
	if len(suiCoins) == 0 {
		log.Fatal("no SUI coins found for sender")
	}
	gasCoin, err := gosuisdk.GetObject(conn, *suiCoins[0].ObjectId, suiCoins[0].Version, ctx)
	if err != nil {
		panic(err)
	}

	// ── Build transaction ─────────────────────────────────────────────────
	b := gosuisdk.NewBuilder(ctx, mod)

	if err := b.SetConfig(sender, 10_000_000, 1_000); err != nil {
		log.Fatalf("SetConfig: %v", err)
	}
	if err := b.AddGasObject(
		fmt.Sprintf("%v", gasCoin.Object.ObjectId),
		uint64(*gasCoin.Object.Version),
		fmt.Sprintf("%s", *gasCoin.Object.Digest),
	); err != nil {
		log.Fatalf("AddGasObject: %v", err)
	}

	// SplitCoins(Gas, [splitMIST])
	gasArgID := b.GasArgument()
	amountID := b.PureU64(splitMIST)
	baseID, err := b.SplitCoins(gasArgID, []uint64{amountID})
	if err != nil {
		log.Fatalf("SplitCoins: %v", err)
	}

	// The single result coin is at sub-index 0.
	coinID := b.NestedResult(baseID, 0)

	// TransferObjects([coin], @recipient)
	recipientID, err := b.PureAddress(recipient)
	if err != nil {
		log.Fatalf("PureAddress: %v", err)
	}
	if err := b.TransferObjects([]uint64{coinID}, recipientID); err != nil {
		log.Fatalf("TransferObjects: %v", err)
	}

	// Serialise — consumes the builder.
	bcsBytes, err := b.Build()
	if err != nil {
		log.Fatalf("Build: %v", err)
	}

	fmt.Printf("Transaction BCS (%d bytes): %x\n", len(bcsBytes), bcsBytes)

	// TODO: sign bcsBytes with your keypair and submit via Sui gRPC.
}
