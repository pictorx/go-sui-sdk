package gosuisdk

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/block-vision/sui-go-sdk/signer"
	pb "github.com/pictorx/go-sui-sdk/sui_rpc_proto/generated"
	"github.com/tetratelabs/wazero/api"
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

type Coin struct {
	Type string
}

func (c *Coin) String() string {
	prefix := "0x0000000000000000000000000000000000000000000000000000000000000002::coin::Coin"
	return prefix + "<" + c.Type + ">"
}

var SuiCoin Coin = Coin{
	Type: "0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI",
}

func VerifySignature(conn *grpc.ClientConn, txBytes, signature []byte, ctx context.Context) {
	client := pb.NewSignatureVerificationServiceClient(conn)
	resp, err := client.VerifySignature(ctx, &pb.VerifySignatureRequest{
		Message: &pb.Bcs{Value: txBytes},
		Signature: &pb.UserSignature{
			Bcs:    &pb.Bcs{Value: signature},
			Scheme: pb.SignatureScheme_ED25519.Enum(),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

var schemeMap = map[byte]pb.SignatureScheme{
	0x00: pb.SignatureScheme_ED25519,
	0x01: pb.SignatureScheme_SECP256K1,
	0x02: pb.SignatureScheme_SECP256R1,
}

func SignExecuteTransaction(conn *grpc.ClientConn, txBytes, signature []byte, ctx context.Context) (*pb.ExecuteTransactionResponse, error) {
	// The serialized signature format is: [flag: 1 byte][sig: 64 bytes][pubkey: 32 bytes]
	if len(signature) != 97 {
		return nil, fmt.Errorf("invalid signature length: expected 97, got %d", len(signature))
	}

	// Extract components
	flagByte := signature[0]        // Should be 0x00 for Ed25519
	sigBytes := signature[1:65]     // 64-byte signature
	pubKeyBytes := signature[65:97] // 32-byte public key

	scheme, exists := schemeMap[flagByte]
	if !exists {
		return nil, fmt.Errorf("Unsupported signature scheme flag: 0x%02x", flagByte)
	}

	client := pb.NewTransactionExecutionServiceClient(conn)
	resp, err := client.ExecuteTransaction(ctx, &pb.ExecuteTransactionRequest{
		Transaction: &pb.Transaction{
			Bcs: &pb.Bcs{Value: txBytes},
		},
		Signatures: []*pb.UserSignature{
			{
				Scheme: scheme.Enum(),
				Signature: &pb.UserSignature_Simple{
					Simple: &pb.SimpleSignature{
						Scheme:    scheme.Enum(),
						Signature: sigBytes,
						PublicKey: pubKeyBytes,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetGas(conn *grpc.ClientConn, ctx context.Context) (*pb.GetEpochResponse, error) {
	client := pb.NewLedgerServiceClient(conn)
	resp, err := client.GetEpoch(ctx, &pb.GetEpochRequest{})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SimulateTransaction(conn *grpc.ClientConn, txBytes []byte, ctx context.Context) (*pb.SimulateTransactionResponse, error) {
	client := pb.NewTransactionExecutionServiceClient(conn)
	resp, err := client.SimulateTransaction(ctx, &pb.SimulateTransactionRequest{
		Transaction: &pb.Transaction{
			Bcs: &pb.Bcs{Value: txBytes},
		},
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Helper function to extract estimated budget from simulation response
func EstimateGasBudget(resp *pb.SimulateTransactionResponse) (uint64, error) {
	effects := resp.Transaction.GetEffects()
	if !effects.GetStatus().GetSuccess() {
		return 0, fmt.Errorf("simulation failed: %s", effects.GetStatus().GetError())
	}

	gasUsed := effects.GetGasUsed()

	// Budget must cover Computation + Storage
	// We do NOT subtract the rebate here; the rebate is a refund applied *after* execution.
	estimatedCost := gasUsed.GetComputationCost() + gasUsed.GetStorageCost()

	// Add a small safety buffer (e.g., 5-10%) just to be safe against slight network fluctuations
	// 2.97M becomes ~3.1M
	buffer := estimatedCost / 10
	finalBudget := estimatedCost + buffer

	return finalBudget, nil
}

type SplitCoin struct {
	Sender    string
	Recipient string
	Gasbudget uint64
	Gasprice  uint64
	Amount    uint64
	GasCoin   *pb.GetObjectResponse
}

func (split *SplitCoin) buildTx(mod api.Module, ctx context.Context) ([]byte, error) {
	b := NewBuilder(ctx, mod)
	// Set config with the specific budget passed in
	if err := b.SetConfig(split.Sender, split.Gasbudget, split.Gasprice); err != nil {
		return nil, err
	}

	// Add Gas Object
	if err := b.AddGasObject(*split.GasCoin.Object.ObjectId, uint64(*split.GasCoin.Object.Version), *split.GasCoin.Object.Digest); err != nil {
		return nil, err
	}

	// ... Add your transaction commands (SplitCoins, Transfer, etc) ...
	// (Copy your existing logic here)
	gasArg := b.GasArgument()
	amt := b.PureU64(split.Amount)
	res, _ := b.SplitCoins(gasArg, []uint64{amt})
	coin := b.NestedResult(res, 0)
	rec, _ := b.PureAddress(split.Recipient)
	b.TransferObjects([]uint64{coin}, rec)

	return b.Build()
}

func (split *SplitCoin) SignExecuteTx(conn *grpc.ClientConn, mod api.Module, account *signer.Signer, ctx context.Context) (*pb.ExecuteTransactionResponse, error) {
	simBytes, err := split.buildTx(mod, ctx)
	if err != nil {
		return nil, err
	}
	simResp, err := SimulateTransaction(conn, simBytes, ctx)
	if err != nil {
		return nil, err
	}

	optimalBudget, err := EstimateGasBudget(simResp)
	if err != nil {
		return nil, err
	}

	split.Gasbudget = optimalBudget
	execBytes, err := split.buildTx(mod, ctx)
	if err != nil {
		return nil, err
	}

	signed, err := SignTransaction(execBytes, account)
	if err != nil {
		return nil, err
	}

	txBytesRaw, err := base64.StdEncoding.DecodeString(signed.TxBytes)
	if err != nil {
		return nil, err
	}

	signatureRaw, err := base64.StdEncoding.DecodeString(signed.Signature)
	if err != nil {
		return nil, err
	}

	resp, err := SignExecuteTransaction(
		conn,
		txBytesRaw,
		signatureRaw,
		ctx,
	)

	return resp, err

}
