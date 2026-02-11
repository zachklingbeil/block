package zero

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Zero struct {
	Rpc        *rpc.Client
	Eth        *ethclient.Client
	Sourcify   *sql.DB
	ContractDb *sql.DB
	Http       *http.Client
	context.Context
	*sync.RWMutex
	*sync.Cond
}

func Init(password string) *Zero {
	ctx := context.Background()

	rpcClient, err := rpc.DialIPC(ctx, "/.ethereum/geth.ipc")
	if err != nil {
		log.Fatalf("ethereum: %v", err)
	}

	sourcifyConn := fmt.Sprintf("user=postgres password=%s dbname=sourcify host=postgres port=5432 sslmode=disable", password)
	sourcifyDb, err := sql.Open("postgres", sourcifyConn)
	if err != nil {
		log.Fatalf("postgres open sourcify: %v", err)
	}
	if err := sourcifyDb.Ping(); err != nil {
		log.Fatalf("postgres ping sourcify: %v", err)
	}

	contractConn := fmt.Sprintf("user=postgres password=%s dbname=contractdb host=postgres port=5432 sslmode=disable", password)
	contractDb, err := sql.Open("postgres", contractConn)
	if err != nil {
		log.Fatalf("postgres open contractdb: %v", err)
	}
	if err := contractDb.Ping(); err != nil {
		log.Fatalf("postgres ping contractdb: %v", err)
	}

	rw := &sync.RWMutex{}
	return &Zero{
		RWMutex:    rw,
		Cond:       sync.NewCond(rw),
		Context:    ctx,
		Http:       &http.Client{},
		Rpc:        rpcClient,
		Eth:        ethclient.NewClient(rpcClient),
		Sourcify:   sourcifyDb,
		ContractDb: contractDb,
	}
}

func (z *Zero) Close() {
	if z.Rpc != nil {
		z.Rpc.Close()
	}
	if z.Sourcify != nil {
		z.Sourcify.Close()
	}
	if z.ContractDb != nil {
		z.ContractDb.Close()
	}
}

func (z *Zero) ConnectPostgres(dbName, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=%s dbname=%s host=postgres port=5432 sslmode=disable", password, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database '%s': %w", dbName, err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database '%s': %w", dbName, err)
	}
	z.Sourcify = db
	return db, nil
}
