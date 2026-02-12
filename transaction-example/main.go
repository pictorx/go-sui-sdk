package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/block-vision/sui-go-sdk/signer"
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

	gas, err := gosuisdk.GetGas(conn, ctx)
	if err != nil {
		panic(err)
	}
	gasPrice := gas.Epoch.ReferenceGasPrice

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

	account, err := signer.NewSignerWithSecretKey("example_priv_key")
	if err != nil {
		panic(err)
	}

	// ── Build transaction ─────────────────────────────────────────────────
	split := gosuisdk.SplitCoin{
		Gasbudget: 100_000_000,
		Gasprice:  *gasPrice,
		Amount:    splitMIST / 10,
		GasCoin:   gasCoin,
		Sender:    sender,
		Recipient: sender,
	}

	resp, err := split.SignExecuteTx(conn, mod, account, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
