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
)

type Zero struct {
	Rpc  *rpc.Client
	Eth  *ethclient.Client
	DB   *sql.DB
	Http *http.Client
	context.Context
	*sync.RWMutex
	*sync.Cond
}

func Init() *Zero {
	ctx := context.Background()
	var rpcClient *rpc.Client
	var err error

	rpcClient, err = rpc.DialIPC(ctx, "/.ethereum/geth.ipc")
	if err != nil {
		log.Fatalf("ethereum: %v", err)
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

func (z *Zero) Close() {
	if z.Rpc != nil {
		z.Rpc.Close()
	}
	if z.DB != nil {
		z.DB.Close()
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
	z.DB = db
	return db, nil
}
