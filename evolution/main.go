package peer

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"sort"
// 	"strconv"
// )

// // SavePeersToJSON saves the Peers.Map to a JSON file as a slice of Peer objects and logs the count.
// // SavePeersToJSON saves the Peers.Map to a JSON file with specific fields and uses the int64 version of LoopringID.
// // SavePeersToJSON saves the Peers.Map to a JSON file as a slice of Peer objects, sorted by LoopringIDINT in ascending order.
// func (p *Peers) SavePeersToJSON(filename string) error {
// 	p.Mu.RLock()
// 	defer p.Mu.RUnlock()

// 	// Define a custom struct for serialization
// 	type PeerOutput struct {
// 		Address     string `json:"address"`
// 		ENS         string `json:"ens"`
// 		LoopringENS string `json:"loopringEns"`
// 		LoopringID  int64  `json:"loopringId"`
// 	}

// 	// Create a slice to hold the serialized Peer objects
// 	var peersSlice []PeerOutput
// 	for _, peer := range p.Map {
// 		peersSlice = append(peersSlice, PeerOutput{
// 			Address:     peer.Address,
// 			ENS:         peer.ENS,
// 			LoopringENS: peer.LoopringENS,
// 			LoopringID:  peer.LoopringIDINT, // Use the int64 version of LoopringID
// 		})
// 	}

// 	// Sort the slice by LoopringIDINT in ascending order
// 	sort.Slice(peersSlice, func(i, j int) bool {
// 		return peersSlice[i].LoopringID < peersSlice[j].LoopringID
// 	})

// 	// Log the number of peers
// 	fmt.Printf("Number of peers to save: %d\n", len(peersSlice))

// 	// Marshal the slice to JSON
// 	data, err := json.MarshalIndent(peersSlice, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal peers to JSON: %w", err)
// 	}

// 	// Write the JSON data to a file
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to create JSON file: %w", err)
// 	}
// 	defer file.Close()

// 	if _, err := file.Write(data); err != nil {
// 		return fmt.Errorf("failed to write JSON data to file: %w", err)
// 	}

// 	fmt.Printf("Peers saved to JSON file: %s\n", filename)
// 	return nil
// }

// // UpdateLoopringIDInt updates the LoopringIDINT field in memory by converting LoopringID strings to integers.
// // Assigns -1 to LoopringIDINT if the LoopringID contains a ".".
// func (p *Peers) UpdateLoopringIDInt() error {
// 	p.Mu.RLock()
// 	defer p.Mu.RUnlock()

// 	for _, peer := range p.Map {
// 		// Check if LoopringID contains a "."
// 		if peer.LoopringID == "." {
// 			peer.LoopringIDINT = -1
// 			continue
// 		}

// 		// Convert LoopringID string to int64
// 		loopringIDInt, err := strconv.ParseInt(peer.LoopringID, 10, 64)
// 		if err != nil {
// 			peer.LoopringIDINT = -1
// 			fmt.Printf("Failed to convert LoopringID '%s' to int for address '%s': %v\n", peer.LoopringID, peer.Address, err)
// 			continue
// 		}

// 		// Update the in-memory map
// 		peer.LoopringIDINT = loopringIDInt
// 	}

// 	fmt.Println("LoopringIDINT field updated successfully in memory.")
// 	return nil
// }

// // InsertBlock inserts a block into the database.
// func (l *Loopring) InsertBlock(in *BlockIn) error {
// 	query := `
//         INSERT INTO loopring (block_id, block_size, created, tx_hash, transactions)
//         VALUES ($1, $2, $3, $4, $5)
//         ON CONFLICT (created) DO NOTHING
//     `
// 	transactions, err := json.Marshal(in.Transactions)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal transactions: %w", err)
// 	}

// 	if _, err := l.Factory.Db.Exec(query, in.Number, in.Size, in.Created, in.TxHash, transactions); err != nil {
// 		return fmt.Errorf("failed to insert in into database: %w", err)
// 	}
// 	l.Factory.Json.Print(in.Number)
// 	return nil
// }

// // QualityControl checks if each block in the database has transactions.
// // If a block does not have transactions, it fetches the block and updates the database.
// func (l *Loopring) QualityControl() error {
// 	query := `SELECT block_id, transactions FROM loopring`

// 	rows, err := l.Factory.Db.Query(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to query blocks from the database: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var blockID int64
// 		var transactionsJSON []byte

// 		if err := rows.Scan(&blockID, &transactionsJSON); err != nil {
// 			return fmt.Errorf("failed to scan block row: %w", err)
// 		}

// 		// Check if the transactions slice is empty
// 		if len(transactionsJSON) == 0 || string(transactionsJSON) == "[]" {
// 			fmt.Printf("BlockIn %d has no transactions. Fetching block data...\n", blockID)
// 			if err := l.GetBlock(int(blockID)); err != nil {
// 				fmt.Printf("Failed to fetch block %d: %v\n", blockID, err)
// 				continue
// 			}
// 			fmt.Printf("Successfully updated block %d with transactions.\n", blockID)
// 		}
// 	}
// 	return nil
// }

// // LoadBlocks queries the loopring table, processes the data, and inserts Blocks into the coords table
// func (l *Loopring) LoadBlocks() error {
// 	// Query the loopring table to fetch data
// 	query := `
//         SELECT created, block_id, block_size
//         FROM loopring
//     `
// 	rows, err := l.Factory.Db.Query(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to query loopring table: %w", err)
// 	}
// 	defer rows.Close()

// 	// Create a slice of BlockIn
// 	var b []BlockIn
// 	for rows.Next() {
// 		var in BlockIn
// 		if err := rows.Scan(&in.Created, &in.Number, &in.Size); err != nil {
// 			return fmt.Errorf("failed to scan row: %w", err)
// 		}
// 		b = append(b, in)
// 	}

// 	// Process the b into Blocks
// 	blocks := l.ProcessInputs(b)

// 	// Insert the Blocks into the coords table
// 	for _, block := range blocks {
// 		if err := l.InsertBlockToCoords(&block); err != nil {
// 			return fmt.Errorf("failed to insert block into coords table: %w", err)
// 		}
// 	}

// 	return nil
// }

// // InsertBlockToCoords inserts a block into the coords table
// func (l *Loopring) InsertBlockToCoords(o *Output) error {
// 	query := `
//         INSERT INTO coords (block_id, block_size, created, coords)
//         VALUES ($1, $2, $3, $4)
//         ON CONFLICT (coords) DO NOTHING
//     `
// 	if _, err := l.Factory.Db.Exec(query, o.Number, o.Size, o.Timestamp, o.Coords); err != nil {
// 		return fmt.Errorf("failed to insert block into coords table: %w", err)
// 	}
// 	return nil
// }

// func (l *Loopring) CreateCoordsTable() error {
// 	query := `
//         CREATE TABLE IF NOT EXISTS coords (
//             block_id BIGINT NOT NULL,
//             block_size BIGINT NOT NULL,
//             created BIGINT NOT NULL,
//             coords TEXT NOT NULL,
//             PRIMARY KEY (coords) -- Use coords as the primary key
//         )
//     `
// 	if _, err := l.Factory.Db.Exec(query); err != nil {
// 		return fmt.Errorf("failed to create coords table: %w", err)
// 	}
// 	return nil
// }

// func (l *Loopring) OutputCoordsAsJSON() error {
// 	query := `
//         SELECT block_id, block_size, created, coords
//         FROM coords
//     `
// 	rows, err := l.Factory.Db.Query(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to query coords table: %w", err)
// 	}
// 	defer rows.Close()

// 	// Create a slice to hold the results
// 	var results []Output
// 	for rows.Next() {
// 		var output Output
// 		if err := rows.Scan(&output.Number, &output.Size, &output.Timestamp, &output.Coords); err != nil {
// 			return fmt.Errorf("failed to scan row: %w", err)
// 		}
// 		results = append(results, output)
// 	}

// 	// Convert the results to JSON
// 	jsonData, err := json.MarshalIndent(results, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal results to JSON: %w", err)
// 	}

// 	// Write JSON to a file or print to console
// 	file, err := os.Create("coords.json")
// 	if err != nil {
// 		return fmt.Errorf("failed to create JSON file: %w", err)
// 	}
// 	defer file.Close()

// 	if _, err := file.Write(jsonData); err != nil {
// 		return fmt.Errorf("failed to write JSON to file: %w", err)
// 	}

// 	fmt.Println("Coords table exported to coords.json")
// 	return nil
// }

// func (l *Loopring) OutputLoopringAsJSON() error {
// 	query := `
//         SELECT block_id, block_size, created, tx_hash, transactions
//         FROM loopring
//     `
// 	rows, err := l.Factory.Db.Query(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to query loopring table: %w", err)
// 	}
// 	defer rows.Close()

// 	// Create a slice to hold the results
// 	var results []map[string]interface{}
// 	for rows.Next() {
// 		var blockID int64
// 		var blockSize int64
// 		var created int64
// 		var txHash string
// 		var transactionsJSON []byte

// 		if err := rows.Scan(&blockID, &blockSize, &created, &txHash, &transactionsJSON); err != nil {
// 			return fmt.Errorf("failed to scan row: %w", err)
// 		}

// 		// Store the data in a map to preserve all fields
// 		blockData := map[string]interface{}{
// 			"block_id":     blockID,
// 			"block_size":   blockSize,
// 			"created":      created,
// 			"tx_hash":      txHash,
// 			"transactions": json.RawMessage(transactionsJSON), // Preserve raw JSON
// 		}

// 		results = append(results, blockData)
// 	}

// 	// Convert the results to JSON
// 	jsonData, err := json.MarshalIndent(results, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal results to JSON: %w", err)
// 	}

// 	// Write JSON to a file or print to console
// 	file, err := os.Create("loopring.json")
// 	if err != nil {
// 		return fmt.Errorf("failed to create JSON file: %w", err)
// 	}
// 	defer file.Close()

// 	if _, err := file.Write(jsonData); err != nil {
// 		return fmt.Errorf("failed to write JSON to file: %w", err)
// 	}

// 	fmt.Println("Loopring table exported to loopring.json")
// 	return nil
// }

// // getPeers extracts unique addresses from a block's transactions.
// func (l *Loopring) getPeers(block *BlockIn) {
// 	one := make(map[int64]string)

// 	for _, tx := range block.Transactions {
// 		if tx.To != 0 {
// 			one[tx.To] = tx.ToAddress
// 		}
// 		if tx.From != 0 && one[tx.From] == "" {
// 			one[tx.From] = ""
// 		}
// 	}
// 	addresses := make([]string, 0, len(one))
// 	for _, address := range one {
// 		if address != "" {
// 			addresses = append(addresses, address)
// 		}
// 	}
// 	go func() {
// 		l.Factory.Peer.NewBlock(addresses)
// 	}()
// }

// // ProcessInputs converts a slice of Block into a slice of Output
// func (l *Loopring) ProcessInputs(in []BlockIn) []Output {
// 	blocks := make([]Output, len(in))

// 	for i, block := range in {
// 		blocks[i] = fx(block)
// 	}
// 	return blocks
// }

// func (l *Loopring) StoreTransactions(blockNumber int64, transactions []any) error {
// 	// Ensure the table exists
// 	createTableQuery := `
//         CREATE TABLE IF NOT EXISTS loopring (
//             block BIGINT PRIMARY KEY,
//             transactions JSONB NOT NULL
//         );
//     `
// 	_, err := l.Factory.Db.Exec(createTableQuery)
// 	if err != nil {
// 		return fmt.Errorf("failed to create table: %w", err)
// 	}

// 	// Convert transactions to JSON
// 	txJSON, err := json.Marshal(transactions)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal transactions: %w", err)
// 	}

// 	// Insert into the database
// 	query := `
//         INSERT INTO loopring (block, transactions)
//         VALUES ($1, $2)
//         ON CONFLICT (block) DO UPDATE
//         SET transactions = EXCLUDED.transactions;
//     `
// 	_, err = l.Factory.Db.Exec(query, blockNumber, txJSON)
// 	if err != nil {
// 		return fmt.Errorf("failed to store transactions: %w", err)
// 	}

// 	return nil
// }

// // DiskToMem converts tables into slices of structs.
// func (d *Database) DiskToMem(table string, result any) error {
// 	query := fmt.Sprintf("SELECT * FROM %s", table)
// 	rows, err := d.Query(query)
// 	if err != nil {
// 		return fmt.Errorf("query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	cols, err := rows.Columns()
// 	if err != nil {
// 		return fmt.Errorf("failed to get columns: %w", err)
// 	}

// 	colCount := len(cols)
// 	jsonRows := make([]map[string]any, 0, 10000)
// 	values := make([]any, colCount)
// 	valuePtrs := make([]any, colCount)

// 	for i := range values {
// 		valuePtrs[i] = &values[i]
// 	}

// 	for rows.Next() {
// 		rowMap := make(map[string]any, colCount)

// 		if err := rows.Scan(valuePtrs...); err != nil {
// 			return fmt.Errorf("scan failed: %w", err)
// 		}

// 		for i, col := range cols {
// 			rowMap[col] = values[i]
// 		}
// 		jsonRows = append(jsonRows, rowMap)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return fmt.Errorf("iteration error: %w", err)
// 	}

// 	jsonData, err := json.Marshal(jsonRows)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal results to JSON: %w", err)
// 	}
// 	if err := json.Unmarshal(jsonData, result); err != nil {
// 		return fmt.Errorf("failed to unmarshal JSON into result: %w", err)
// 	}
// 	return nil
// }

// func (d *Database) ColumnToSlice(table string, column string, result any) error {
// 	query := fmt.Sprintf("SELECT %s FROM %s", column, table)
// 	rows, err := d.Query(query)
// 	if err != nil {
// 		return fmt.Errorf("query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	slice := make([]any, 0, 250000) // Preallocate a slice with an initial capacity
// 	for rows.Next() {
// 		var value any
// 		if err := rows.Scan(&value); err != nil {
// 			return fmt.Errorf("scan failed: %w", err)
// 		}
// 		slice = append(slice, value)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return fmt.Errorf("iteration error: %w", err)
// 	}

// 	// Marshal the slice into JSON and unmarshal it into the provided result
// 	jsonData, err := json.Marshal(slice)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal results to JSON: %w", err)
// 	}
// 	if err := json.Unmarshal(jsonData, result); err != nil {
// 		return fmt.Errorf("failed to unmarshal JSON into result: %w", err)
// 	}
// 	return nil
// }

// func (p *Peers) LoadPeersFromJSON() error {
// 	// Use the embedded JSON data
// 	data := []byte(embeddedPeersJSON)

// 	// Unmarshal the JSON data into a slice of Peer objects
// 	var peers []Peer
// 	if err := json.Unmarshal(data, &peers); err != nil {
// 		return fmt.Errorf("failed to unmarshal embedded JSON data: %w", err)
// 	}

// 	// Populate the Map and Addresses fields
// 	p.Mu.Lock()
// 	defer p.Mu.Unlock()

// 	for _, peer := range peers {
// 		// Add the peer to the Map
// 		p.Map[peer.Address] = &peer

// 		// Add the peer's address to the Addresses slice if fields are invalid
// 		if peer.ENS == "." || peer.LoopringENS == "." || peer.LoopringID == -1 {
// 			p.Addresses = append(p.Addresses, peer.Address)
// 		}
// 	}

// 	fmt.Printf("%d peers loaded from embedded JSON\n", len(peers))
// 	return nil
// }

// func (p *Peers) LoadAndSavePeersFromJSON() error {
// 	// Load peers from the embedded JSON
// 	if err := p.LoadPeersFromJSON(); err != nil {
// 		return fmt.Errorf("failed to load peers from embedded JSON: %w", err)
// 	}

// 	// Save all peers to the database
// 	if err := p.SavePeers(); err != nil {
// 		return fmt.Errorf("failed to save peers to the database: %w", err)
// 	}

// 	return nil
// }

// func (p *Peers) CreatePeersTable() error {
// 	query := `
//     CREATE TABLE IF NOT EXISTS peers (
//         address TEXT PRIMARY KEY,
//         ens TEXT NOT NULL,
//         loopringEns TEXT NOT NULL,
//         loopringId BIGINT NOT NULL
//     );
//     `
// 	_, err := p.Db.Exec(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to create peers table: %w", err)
// 	}
// 	fmt.Println("Peers table created or already exists.")
// 	return nil
// }
