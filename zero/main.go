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
	Rpc      *rpc.Client
	Eth      *ethclient.Client
	postgres *sql.DB
	Http     *http.Client

	Sig      sigprovider.AbiServiceClient
	DB       bytecodedb.DatabaseClient
	Ops      userops.UserOpsServiceClient
	Solidity verifier.SolidityVerifierClient
	Vyper    verifier.VyperVerifierClient
	Sourcify verifier.SourcifyVerifierClient

	grpcConns []*grpc.ClientConn

	context.Context
	*sync.RWMutex
	*sync.Cond
}

func Init() *Zero {
	rw := &sync.RWMutex{}
	return &Zero{
		RWMutex: rw,
		Cond:    sync.NewCond(rw),
		Context: context.Background(),
	}
}

func (f *Zero) Node() error {
	rpc, err := rpc.DialIPC(f.Context, "/.ethereum/geth.ipc")
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return nil
	}
	f.Rpc = rpc
	f.Eth = ethclient.NewClient(rpc)
	return nil
}

func (f *Zero) NodeDial(url string) error {
	client, err := rpc.DialContext(f.Context, url)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", url, err)
		return err
	}
	f.Rpc = client
	f.Eth = ethclient.NewClient(client)
	return nil
}

func (f *Zero) dial(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	f.grpcConns = append(f.grpcConns, conn)
	return conn, nil
}

// Blockscout connects all gRPC microservices on the bridge network.
func (f *Zero) Blockscout() error {
	sig, err := f.dial("sig-provider:8051")
	if err != nil {
		return fmt.Errorf("sig-provider: %w", err)
	}
	f.Sig = sigprovider.NewAbiServiceClient(sig)

	db, err := f.dial("eth-bytecode-db:8051")
	if err != nil {
		return fmt.Errorf("eth-bytecode-db: %w", err)
	}
	f.DB = bytecodedb.NewDatabaseClient(db)

	ops, err := f.dial("user-ops-indexer:8051")
	if err != nil {
		return fmt.Errorf("user-ops-indexer: %w", err)
	}
	f.Ops = userops.NewUserOpsServiceClient(ops)

	sc, err := f.dial("smart-contract-verifier:8051")
	if err != nil {
		return fmt.Errorf("smart-contract-verifier: %w", err)
	}
	f.Solidity = verifier.NewSolidityVerifierClient(sc)
	f.Vyper = verifier.NewVyperVerifierClient(sc)
	f.Sourcify = verifier.NewSourcifyVerifierClient(sc)

	return nil
}

func (f *Zero) Close() {
	if f.Rpc != nil {
		f.Rpc.Close()
	}
	if f.postgres != nil {
		f.postgres.Close()
	}
	for _, c := range f.grpcConns {
		c.Close()
	}
}

func (f *Zero) ConnectPostgres(dbName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=postgres dbname=%s host=postgres port=5432 sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database '%s': %w", dbName, err)
	}
	f.postgres = db
	return db, nil
}
