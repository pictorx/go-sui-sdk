package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/pictorx/go-sui-sdk/sui_rpc_proto/generated" // Replace with your actual import path for the generated package (e.g., github.com/yourrepo/sui/rpc/v2)
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// Endpoint for Sui mainnet full node gRPC (adjust for testnet: fullnode.testnet.sui.io:443)
	endpoint := "fullnode.testnet.sui.io:443"

	// Dial the gRPC server with TLS (using system cert pool for verification)
	conn, err := grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Create clients for specific services (from your generated code)
	ledgerClient := pb.NewLedgerServiceClient(conn) // For reading chain data like checkpoints, objects, transactions
	// txExecClient := pb.NewTransactionExecutionServiceClient(conn) // For submitting/executing transactions
	// Add other clients as needed, e.g., pb.NewMovePackageServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := ledgerClient.GetEpoch(ctx, &pb.GetEpochRequest{})
	if err != nil {
		log.Fatalf("GetEpoch failed: %v", err)
	}

	// Print the response (adjust based on your needs)
	fmt.Printf("%+v\n", resp)
}
