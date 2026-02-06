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
	Rpc      *rpc.Client
	Eth      *ethclient.Client
	postgres *sql.DB
	Http     *http.Client
	context.Context
	*sync.RWMutex
	*sync.Cond
}

func Init() *Zero {
	rw := &sync.RWMutex{}
	zero := &Zero{
		RWMutex: rw,
		Cond:    sync.NewCond(rw),
		Context: context.Background(),
	}
	return zero
}

// Establish geth.ipc connection
func (f *Zero) Node() error {
	rpc, err := rpc.DialIPC(f.Context, "/.ethereum/geth.ipc")
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return nil
	}
	eth := ethclient.NewClient(rpc)
	f.Rpc = rpc
	f.Eth = eth
	return nil
}

// Establish JSON-RPC connection (http, https, ws, wss)
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
