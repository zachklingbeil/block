package zero

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/timefactoryio/block/zero/proto/bytecodedb"
	"github.com/timefactoryio/block/zero/proto/sigprovider"
	"github.com/timefactoryio/block/zero/proto/userops"
	"github.com/timefactoryio/block/zero/proto/verifier"
)

type Zero struct {
	Rpc       *rpc.Client
	Eth       *ethclient.Client
	DB        *sql.DB
	Http      *http.Client
	Sig       sigprovider.AbiServiceClient
	ByteDB    bytecodedb.DatabaseClient
	Ops       userops.UserOpsServiceClient
	Solidity  verifier.SolidityVerifierClient
	Vyper     verifier.VyperVerifierClient
	Sourcify  verifier.SourcifyVerifierClient
	grpcConns []*grpc.ClientConn
	context.Context
	*sync.RWMutex
	*sync.Cond
}

func Init(url string) *Zero {
	ctx := context.Background()

	var rpcClient *rpc.Client
	var err error

	if url == "" {
		rpcClient, err = rpc.DialIPC(ctx, "/.ethereum/geth.ipc")
	} else {
		rpcClient, err = rpc.DialContext(ctx, url)
	}

	if err != nil {
		log.Fatalf("failed to connect to Ethereum client: %v", err)
	}

	rw := &sync.RWMutex{}
	return &Zero{
		RWMutex: rw,
		Cond:    sync.NewCond(rw),
		Context: ctx,
		Rpc:     rpcClient,
		Eth:     ethclient.NewClient(rpcClient),
	}
}

func (z *Zero) dial(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	z.grpcConns = append(z.grpcConns, conn)
	return conn, nil
}

func (z *Zero) Blockscout() error {
	sig, err := z.dial("sig-provider:8051")
	if err != nil {
		return fmt.Errorf("sig-provider: %w", err)
	}
	z.Sig = sigprovider.NewAbiServiceClient(sig)

	db, err := z.dial("eth-bytecode-db:8051")
	if err != nil {
		return fmt.Errorf("eth-bytecode-db: %w", err)
	}
	z.ByteDB = bytecodedb.NewDatabaseClient(db)

	ops, err := z.dial("user-ops-indexer:8051")
	if err != nil {
		return fmt.Errorf("user-ops-indexer: %w", err)
	}
	z.Ops = userops.NewUserOpsServiceClient(ops)

	sc, err := z.dial("smart-contract-verifier:8051")
	if err != nil {
		return fmt.Errorf("smart-contract-verifier: %w", err)
	}
	z.Solidity = verifier.NewSolidityVerifierClient(sc)
	z.Vyper = verifier.NewVyperVerifierClient(sc)
	z.Sourcify = verifier.NewSourcifyVerifierClient(sc)

	return nil
}

func (z *Zero) Close() {
	if z.Rpc != nil {
		z.Rpc.Close()
	}
	if z.DB != nil {
		z.DB.Close()
	}
	for _, c := range z.grpcConns {
		c.Close()
	}
}

func (z *Zero) ConnectPostgres(dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database '%s': %w", dbName, err)
	}
	z.DB = db
	return db, nil
}
