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

	"github.com/timefactoryio/block/zero/proto/sigprovider"
)

type Zero struct {
	Rpc  *rpc.Client
	Eth  *ethclient.Client
	DB   *sql.DB
	Http *http.Client
	Sig  sigprovider.AbiServiceClient
	sig  *grpc.ClientConn
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
		log.Fatalf("ethereum: %v", err)
	}

	sig, err := grpc.NewClient("sig-provider:8051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("sig-provider: %v", err)
	}

	rw := &sync.RWMutex{}
	return &Zero{
		RWMutex: rw,
		Cond:    sync.NewCond(rw),
		Context: ctx,
		Rpc:     rpcClient,
		Eth:     ethclient.NewClient(rpcClient),
		sig:     sig,
		Sig:     sigprovider.NewAbiServiceClient(sig),
	}
}

func (z *Zero) Close() {
	if z.Rpc != nil {
		z.Rpc.Close()
	}
	if z.DB != nil {
		z.DB.Close()
	}
	if z.sig != nil {
		z.sig.Close()
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
