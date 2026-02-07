package fx

// func TestBlock(t *testing.T) {
// 	f := Init()
// 	defer f.Rpc.Close()

// 	block, err := f.Block(nil)
// 	if err != nil {
// 		t.Fatalf("Block: %v", err)
// 	}

// 	t.Logf("Block #%s (%s) — %d txs", block.Header.Number, block.Header.Hash(), len(block.Transactions))

// 	output, err := json.MarshalIndent(block, "", "  ")
// 	if err != nil {
// 		t.Fatalf("Marshal: %v", err)
// 	}

// 	if err := os.WriteFile("../output/block.json", output, 0644); err != nil {
// 		t.Fatalf("WriteFile: %v", err)
// 	}

// 	t.Log("Block written to output/block.json")
// }
