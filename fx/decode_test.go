package fx

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	f := Init()
	defer f.Rpc.Close()

	block, err := f.Block(nil)
	if err != nil {
		t.Fatalf("Block: %v", err)
	}

	decoded, err := f.Decode(f.Context, block)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	t.Logf("Decoded %d transactions", len(decoded.Decoded))

	for i, dt := range decoded.Decoded {
		switch {
		case dt.Deploy:
			t.Logf("  tx[%d] %s DEPLOY", i, dt.Hash)
		case dt.Method != nil && dt.Method.Name != "":
			t.Logf("  tx[%d] %s %s() — %d events", i, dt.Hash, dt.Method.Name, len(dt.Events))
		case dt.Method != nil:
			t.Logf("  tx[%d] %s %s — %d events", i, dt.Hash, dt.Method.Selector, len(dt.Events))
		default:
			t.Logf("  tx[%d] %s transfer — %d events", i, dt.Hash, len(dt.Events))
		}

		if len(dt.UserOps) > 0 {
			t.Logf("    %d user operations", len(dt.UserOps))
		}
	}

	output, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if err := os.WriteFile("../decoded.json", output, 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	t.Log("Decoded block written to decoded.json")
}
