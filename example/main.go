package main

import (
	"context"
	"fmt"
	"log"
	"time"

	gosuisdk "github.com/pictorx/go-sui-sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	endpoint := "fullnode.testnet.sui.io:443"

	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	getEpoch(conn, ctx)
	getObject(conn, ctx)
	getTransaction(conn, ctx)
	batchGetObjects(conn, ctx)
	batchGetTransactions(conn, ctx)
	getServiceInfo(conn, ctx)
	getPackage(conn, ctx)
	getFunction(conn, ctx)
	getDatatype(conn, ctx)
	listPackageVersions(conn, ctx)
	getBalance(conn, ctx)
	getCoinInfo(conn, ctx)
	listBalances(conn, ctx)
	listOwnedObjects(conn, ctx)
	listDynamicFields(conn, ctx)

}

func response(r any, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", r)
}

func getEpoch(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetEpoch(conn, ctx)
	response(resp, err)
}

func getObject(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetObject(conn, "0xe567f65413d10d585bbec909f46f11c5fc666061e1ede471b840f3d1cf25eaaa",
		nil, ctx)
	response(resp, err)
}

func getTransaction(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetTransaction(conn, "6SCQPcEvHHE2c34Hxse7PyxCdWHWYgZQhQpUNpw9314n", ctx)
	response(resp, err)
}

func batchGetObjects(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.BatchGetObjects(conn, map[string]*uint64{
		"0xe567f65413d10d585bbec909f46f11c5fc666061e1ede471b840f3d1cf25eaaa": nil,
		"0x491d113c71512f45612e34eb2a25f7dd1014aa8bdb0718a22c4eea54f2a06340": nil,
	}, ctx)
	response(resp, err)
}

func batchGetTransactions(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.BatchGetTransactions(conn, []string{
		"6SCQPcEvHHE2c34Hxse7PyxCdWHWYgZQhQpUNpw9314n",
		"3WZieTbDv4CnBjbeBGs4AuakJAHySTEqXu6gu1bxY8rV",
	}, ctx)
	response(resp, err)
}

func getServiceInfo(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetServiceInfo(conn, ctx)
	response(resp, err)
}

func getPackage(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetPackage(conn, "0xd84704c17fc870b8764832c535aa6b11f21a95cd6f5bb38a9b07d2cf42220c66", ctx)
	response(resp, err)
}

func getFunction(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetFunction(conn, "0xd84704c17fc870b8764832c535aa6b11f21a95cd6f5bb38a9b07d2cf42220c66", "system", "reserve_space", ctx)
	response(resp, err)
}

func getDatatype(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetDatatype(conn, "0xd84704c17fc870b8764832c535aa6b11f21a95cd6f5bb38a9b07d2cf42220c66", "system", "System", ctx)
	response(resp, err)
}

func listPackageVersions(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.ListPackageVersions(conn, "0xd84704c17fc870b8764832c535aa6b11f21a95cd6f5bb38a9b07d2cf42220c66", nil, nil, ctx)
	response(resp, err)
}

func getBalance(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetBalance(conn, "0x8aeec8403b86f22e58d87bdc85ff78c87b69dce58b8f651900b9eb5644f45180", "0x8270feb7375eee355e64fdb69c50abb6b5f9393a722883c1cf45f8e26048810a::wal::WAL", ctx)
	response(resp, err)
}

func getCoinInfo(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.GetCoinInfo(conn, "0x8270feb7375eee355e64fdb69c50abb6b5f9393a722883c1cf45f8e26048810a::wal::WAL", ctx)
	response(resp, err)
}

func listBalances(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.ListBalances(conn, "0x8aeec8403b86f22e58d87bdc85ff78c87b69dce58b8f651900b9eb5644f45180", nil, nil, ctx)
	response(resp, err)
}

func listOwnedObjects(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.ListOwnedObjects(conn, "0x8aeec8403b86f22e58d87bdc85ff78c87b69dce58b8f651900b9eb5644f45180", nil, nil, ctx)
	response(resp, err)
}

func listDynamicFields(conn *grpc.ClientConn, ctx context.Context) {
	resp, err := gosuisdk.ListDynamicFields(conn, "0x6c2547cbbc38025cf3adac45f63cb0a8d12ecf777cdc75a4971612bf97fdf6af", nil, nil, ctx)
	response(resp, err)
}
