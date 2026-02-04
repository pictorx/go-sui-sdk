package gosuisdk

import (
	"context"
	"fmt"

	pb "github.com/pictorx/go-sui-sdk/sui_rpc_proto/generated"
	"google.golang.org/grpc"
)

func GetEpoch(conn *grpc.ClientConn, ctx context.Context) (*pb.GetEpochResponse, error) {
	client := pb.NewLedgerServiceClient(conn)
	resp, err := client.GetEpoch(ctx, &pb.GetEpochRequest{})

	if err != nil {
		return nil, err
	}

	return resp, err
}

func GetServiceInfo(conn *grpc.ClientConn, ctx context.Context) (*pb.GetServiceInfoResponse, error) {
	client := pb.NewLedgerServiceClient(conn)
	resp, err := client.GetServiceInfo(ctx, &pb.GetServiceInfoRequest{})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetObject(conn *grpc.ClientConn, objectId string, version *uint64, ctx context.Context) (*pb.GetObjectResponse, error) {
	client := pb.NewLedgerServiceClient(conn)
	resp, err := client.GetObject(ctx, &pb.GetObjectRequest{
		ObjectId: &objectId,
		Version:  version,
	})
	if err != nil {
		return nil, err
	}

	return resp, err
}

func GetTransaction(conn *grpc.ClientConn, digest string, ctx context.Context) (*pb.GetTransactionResponse, error) {
	client := pb.NewLedgerServiceClient(conn)
	resp, err := client.GetTransaction(ctx, &pb.GetTransactionRequest{
		Digest: &digest,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func BatchGetObjects(conn *grpc.ClientConn, objects map[string]*uint64, ctx context.Context) (*pb.BatchGetObjectsResponse, error) {
	if len(objects) == 0 {
		return nil, fmt.Errorf("objects cannot be zero")
	}

	client := pb.NewLedgerServiceClient(conn)
	ObjectRequests := []*pb.GetObjectRequest{}

	for objectId, version := range objects {
		ObjectRequests = append(ObjectRequests, &pb.GetObjectRequest{
			ObjectId: &objectId,
			Version:  version,
		})
	}

	resp, err := client.BatchGetObjects(ctx, &pb.BatchGetObjectsRequest{
		Requests: ObjectRequests,
	})
	if err != nil {
		return nil, err
	}

	return resp, err
}

func BatchGetTransactions(conn *grpc.ClientConn, digests []string, ctx context.Context) (*pb.BatchGetTransactionsResponse, error) {
	if len(digests) == 0 {
		return nil, fmt.Errorf("digests parameter cannot be empty")
	}

	client := pb.NewLedgerServiceClient(conn)
	resp, err := client.BatchGetTransactions(ctx, &pb.BatchGetTransactionsRequest{
		Digests: digests,
	})

	if err != nil {
		return nil, err
	}

	return resp, err
}

func GetPackage(conn *grpc.ClientConn, packageId string, ctx context.Context) (*pb.GetPackageResponse, error) {
	client := pb.NewMovePackageServiceClient(conn)
	resp, err := client.GetPackage(ctx, &pb.GetPackageRequest{
		PackageId: &packageId,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetFunction(conn *grpc.ClientConn, packageId, module, funcName string, ctx context.Context) (*pb.GetFunctionResponse, error) {
	client := pb.NewMovePackageServiceClient(conn)
	resp, err := client.GetFunction(ctx, &pb.GetFunctionRequest{
		PackageId:  &packageId,
		ModuleName: &module,
		Name:       &funcName,
	})

	if err != nil {
		return nil, err
	}

	return resp, err
}

func GetDatatype(conn *grpc.ClientConn, packageId, module, dataTypeName string, ctx context.Context) (*pb.GetDatatypeResponse, error) {
	client := pb.NewMovePackageServiceClient(conn)

	resp, err := client.GetDatatype(ctx, &pb.GetDatatypeRequest{
		PackageId:  &packageId,
		ModuleName: &module,
		Name:       &dataTypeName,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ListPackageVersions(conn *grpc.ClientConn, packageId string, pagesize *uint32, pagetoken []byte, ctx context.Context) (*pb.ListPackageVersionsResponse, error) {
	client := pb.NewMovePackageServiceClient(conn)

	resp, err := client.ListPackageVersions(ctx, &pb.ListPackageVersionsRequest{
		PackageId: &packageId,
		PageSize:  pagesize,
		PageToken: pagetoken,
	})

	if err != nil {
		return nil, err
	}

	return resp, err
}

func GetBalance(conn *grpc.ClientConn, owner, cointype string, ctx context.Context) (*pb.GetBalanceResponse, error) {
	client := pb.NewStateServiceClient(conn)
	resp, err := client.GetBalance(ctx, &pb.GetBalanceRequest{
		Owner:    &owner,
		CoinType: &cointype,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetCoinInfo(conn *grpc.ClientConn, cointype string, ctx context.Context) (*pb.GetCoinInfoResponse, error) {
	client := pb.NewStateServiceClient(conn)
	resp, err := client.GetCoinInfo(ctx, &pb.GetCoinInfoRequest{
		CoinType: &cointype,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ListBalances(conn *grpc.ClientConn, owner string, pagesize *uint32, pagetoken []byte, ctx context.Context) (*pb.ListBalancesResponse, error) {
	client := pb.NewStateServiceClient(conn)
	resp, err := client.ListBalances(ctx, &pb.ListBalancesRequest{
		Owner:     &owner,
		PageSize:  pagesize,
		PageToken: pagetoken,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ListOwnedObjects(conn *grpc.ClientConn, owner string, pagesize *uint32, pagetoken []byte, ctx context.Context) (*pb.ListOwnedObjectsResponse, error) {
	client := pb.NewStateServiceClient(conn)
	resp, err := client.ListOwnedObjects(ctx, &pb.ListOwnedObjectsRequest{
		Owner:     &owner,
		PageSize:  pagesize,
		PageToken: pagetoken,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ListDynamicFields(conn *grpc.ClientConn, objectId string, pagesize *uint32, pagetoken []byte, ctx context.Context) (*pb.ListDynamicFieldsResponse, error) {
	client := pb.NewStateServiceClient(conn)
	resp, err := client.ListDynamicFields(ctx, &pb.ListDynamicFieldsRequest{
		Parent:    &objectId,
		PageSize:  pagesize,
		PageToken: pagetoken,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func OwnedCoins(listownedobjects *pb.ListOwnedObjectsResponse, cointype, owner string) []*pb.Object {
	list := listownedobjects

	var coins []*pb.Object
	for _, v := range list.GetObjects() {

		if *v.ObjectType == cointype {

			coins = append(coins, v)
		}
	}

	return coins
}

func GasPayment(coins []*pb.Object, owner string, price uint64, budget uint64) *pb.GasPayment {
	coinObjects := []*pb.ObjectReference{}
	for _, coin := range coins {
		coinObjects = append(coinObjects, &pb.ObjectReference{
			ObjectId: coin.ObjectId,
			Version:  coin.Version,
			Digest:   coin.Digest,
		})
	}
	return &pb.GasPayment{
		Objects: coinObjects,
		Owner:   &owner,
		Price:   &price,
		Budget:  &budget,
	}
}
