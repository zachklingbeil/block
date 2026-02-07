package fx

// func TestDecode(t *testing.T) {
// 	f := Init()
// 	defer f.Rpc.Close()

// 	block, err := f.Block(nil)
// 	if err != nil {
// 		t.Fatalf("Block: %v", err)
// 	}

// 	decoded, err := f.Decode(f.Context, block)
// 	if err != nil {
// 		t.Fatalf("Decode: %v", err)
// 	}

// 	t.Logf("Block %s: %d transactions", decoded.Hash.Hex(), len(decoded.Txs))

// 	for i, dt := range decoded.Txs {
// 		switch {
// 		case dt.Deploy:
// 			t.Logf("  tx[%d] %s DEPLOY to %s", i, dt.Hash.Hex(), dt.ContractAddress.Hex())
// 		case dt.Method != nil && dt.Method.Name != "":
// 			t.Logf("  tx[%d] %s %s() — %d events", i, dt.Hash.Hex(), dt.Method.Name, len(dt.Events))
// 		case dt.Method != nil:
// 			t.Logf("  tx[%d] %s %s — %d events", i, dt.Hash.Hex(), dt.Method.Selector, len(dt.Events))
// 		default:
// 			t.Logf("  tx[%d] %s transfer — %d events", i, dt.Hash.Hex(), len(dt.Events))
// 		}

// 		if len(dt.UserOps) > 0 {
// 			t.Logf("    %d user operations", len(dt.UserOps))
// 		}
// 	}

// 	output, err := json.MarshalIndent(decoded, "", "  ")
// 	if err != nil {
// 		t.Fatalf("Marshal: %v", err)
// 	}

// 	if err := os.WriteFile("../output/decoded.json", output, 0644); err != nil {
// 		t.Fatalf("WriteFile: %v", err)
// 	}
// 	t.Log("Decoded block written to output/decoded.json")
// }
