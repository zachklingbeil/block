package fx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Log struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    []byte         `json:"data,omitempty"`
	Index   uint           `json:"logIndex"`
	TxIndex uint           `json:"transactionIndex"`
	Removed bool           `json:"removed,omitempty"`
	Event   string         `json:"event,omitempty"`
	Args    []*Arg         `json:"args,omitempty"`
}

type Arg struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed,omitempty"`
	Value   string `json:"value"`
}

// Logs converts raw receipt logs into decoded Logs, resolving event
// signatures from the local Sourcify DB and decoding arguments via geth ABI.
func (fx *Fx) Logs(raw []*types.Log) []*Log {
	logs := make([]*Log, len(raw))
	for i, l := range raw {
		logs[i] = &Log{
			Address: l.Address,
			Topics:  l.Topics,
			Data:    l.Data,
			Index:   l.Index,
			TxIndex: l.TxIndex,
			Removed: l.Removed,
		}
	}

	if len(logs) == 0 {
		return logs
	}

	// Collect unique topic0 hashes
	seen := make(map[common.Hash]struct{})
	var topic0s []common.Hash
	for _, l := range logs {
		if len(l.Topics) == 0 {
			continue
		}
		if _, ok := seen[l.Topics[0]]; !ok {
			seen[l.Topics[0]] = struct{}{}
			topic0s = append(topic0s, l.Topics[0])
		}
	}

	sigs := fx.lookupEventSignatures(topic0s)

	// Decode each log
	for _, l := range logs {
		if len(l.Topics) == 0 {
			continue
		}
		sig, ok := sigs[l.Topics[0]]
		if !ok {
			continue
		}
		decodeLog(l, sig)
	}

	return logs
}

// lookupEventSignatures batch-queries the local Sourcify DB for event signatures.
// Prefers verified event signatures via compiled_contracts_signatures.
func (fx *Fx) lookupEventSignatures(hashes []common.Hash) map[common.Hash]string {
	if len(hashes) == 0 {
		return nil
	}

	placeholders := make([]string, len(hashes))
	args := make([]any, len(hashes))
	for i, h := range hashes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = h.Bytes()
	}

	query := fmt.Sprintf(`
        SELECT DISTINCT ON (s.signature_hash_32)
            s.signature_hash_32,
            s.signature
        FROM signatures s
        LEFT JOIN compiled_contracts_signatures ccs
            ON ccs.signature_hash_32 = s.signature_hash_32
            AND ccs.signature_type = 'event'
        WHERE s.signature_hash_32 IN (%s)
        ORDER BY s.signature_hash_32,
            CASE WHEN ccs.signature_type IS NOT NULL THEN 0 ELSE 1 END,
            length(s.signature)
    `, strings.Join(placeholders, ","))

	rows, err := fx.DB.QueryContext(fx.Context, query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make(map[common.Hash]string, len(hashes))
	for rows.Next() {
		var hashBytes []byte
		var sig string
		if err := rows.Scan(&hashBytes, &sig); err != nil {
			continue
		}
		result[common.BytesToHash(hashBytes)] = sig
	}
	return result
}

// decodeLog decodes a single log's indexed topics and non-indexed data.
func decodeLog(l *Log, sig string) {
	name, allArgs, err := parseSignature(sig, len(l.Topics))
	if err != nil {
		return
	}
	l.Event = name

	var indexedArgs, nonIndexedArgs abi.Arguments
	for _, a := range allArgs {
		if a.Indexed {
			indexedArgs = append(indexedArgs, a)
		} else {
			nonIndexedArgs = append(nonIndexedArgs, a)
		}
	}

	topicValues := make(map[string]any)
	if len(indexedArgs) > 0 && len(l.Topics) > 1 {
		if err := abi.ParseTopicsIntoMap(topicValues, indexedArgs, l.Topics[1:]); err != nil {
			return
		}
	}

	dataValues := make(map[string]any)
	if len(nonIndexedArgs) > 0 && len(l.Data) > 0 {
		if unpacked, err := nonIndexedArgs.UnpackValues(l.Data); err == nil {
			for i, a := range nonIndexedArgs {
				dataValues[a.Name] = unpacked[i]
			}
		}
	}

	args := make([]*Arg, len(allArgs))
	for i, a := range allArgs {
		var val any
		if a.Indexed {
			val = topicValues[a.Name]
		} else {
			val = dataValues[a.Name]
		}
		args[i] = &Arg{
			Name:    a.Name,
			Type:    a.Type.String(),
			Indexed: a.Indexed,
			Value:   formatValue(val),
		}
	}
	l.Args = args
}

func parseSignature(sig string, topicCount int) (string, abi.Arguments, error) {
	open := strings.Index(sig, "(")
	if open < 0 {
		return "", nil, fmt.Errorf("invalid signature: %s", sig)
	}
	name := sig[:open]
	inner := sig[open+1 : len(sig)-1]
	if inner == "" {
		return name, nil, nil
	}

	typeStrs := splitTopLevelParams(inner)
	args := make(abi.Arguments, len(typeStrs))
	indexedCount := topicCount - 1

	for i, ts := range typeStrs {
		typ, err := abi.NewType(strings.TrimSpace(ts), "", nil)
		if err != nil {
			return "", nil, fmt.Errorf("type %q: %w", ts, err)
		}
		args[i] = abi.Argument{
			Name:    fmt.Sprintf("arg%d", i),
			Type:    typ,
			Indexed: i < indexedCount,
		}
	}
	return name, args, nil
}

func splitTopLevelParams(s string) []string {
	var parts []string
	depth, start := 0, 0
	for i, c := range s {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	return append(parts, s[start:])
}

func formatValue(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case common.Address:
		return val.Hex()
	case common.Hash:
		return val.Hex()
	case *big.Int:
		return val.String()
	case []byte:
		return "0x" + hex.EncodeToString(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", val)
	}
}
