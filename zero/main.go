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
	Rpc  *rpc.Client
	Eth  *ethclient.Client
	Http *http.Client
	Db   *sql.DB
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

	db, err := ConnectPostgres(password)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}

	rw := &sync.RWMutex{}
	return &Zero{
		RWMutex: rw,
		Cond:    sync.NewCond(rw),
		Context: ctx,
		Http:    &http.Client{},
		Db:      db,
		Rpc:     rpcClient,
		Eth:     ethclient.NewClient(rpcClient),
	}
}

func (z *Zero) Close() {
	if z.Rpc != nil {
		z.Rpc.Close()
	}

}

func ConnectPostgres(password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=%s dbname=ethereum host=postgres port=5432 sslmode=disable", password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func (z *Zero) CreateContractTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS contracts (
            contract TEXT PRIMARY KEY,
            abi JSONB NOT NULL
        );`
	_, err := z.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create contracts table: %w", err)
	}
	return nil
}
