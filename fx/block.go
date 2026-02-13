package fx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const signatureDB = "http://fourbyte:80/signature-database/v1/lookup"

type Block struct {
	Number       *big.Int       `json:"number"`
	Hash         common.Hash    `json:"hash"`
	Timestamp    uint64         `json:"timestamp"`
	GasLimit     uint64         `json:"gasLimit"`
	GasUsed      uint64         `json:"gasUsed"`
	BaseFee      *big.Int       `json:"baseFeePerGas,omitempty"`
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
	TxHash            common.Hash     `json:"hash"`
	TxIndex           uint            `json:"index"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to,omitempty"`
	Value             *big.Int        `json:"value,omitempty"`
	Input             string          `json:"input,omitempty"`
	Status            uint64          `json:"status"`
	Gas               uint64          `json:"gas"`
	EffectiveGasPrice *big.Int        `json:"gasPrice"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	Logs              []*Event        `json:"logs,omitempty"`
	Decoded           []*Decoded      `json:"decoded,omitempty"`
}

type Event struct {
	Address common.Address `json:"contract"`
	Topics  []string       `json:"topics"`
	Data    string         `json:"data,omitempty"`
}

type Decoded struct {
	Contract  common.Address `json:"contract"`
	Signature string         `json:"signature"`
	Values    string         `json:"values"`
}

func (fx *Fx) Block(number *big.Int) (*Block, error) {
	block, err := fx.Eth.BlockByNumber(fx.Context, number)
	if err != nil {
		return nil, fmt.Errorf("block: %w", err)
	}

	receipts, err := fx.blockReceipts(block.Number())
	if err != nil {
		return nil, err
	}

	eventHashes := make(map[string]struct{})
	for _, r := range receipts {
		for _, l := range r.Logs {
			if len(l.Topics) > 0 {
				eventHashes[l.Topics[0].Hex()] = struct{}{}
			}
		}
	}
	signatures := lookupSignatures(eventHashes)

	signer := types.MakeSigner(fx.Chain, block.Number(), block.Time())
	txs := make([]*Transaction, len(block.Transactions()))

	for i, tx := range block.Transactions() {
		r := receipts[i]
		from, _ := types.Sender(signer, tx)

		var contract *common.Address
		if r.ContractAddress != (common.Address{}) {
			contract = &r.ContractAddress
		}

		txs[i] = &Transaction{
			TxHash:            tx.Hash(),
			TxIndex:           uint(i),
			From:              from,
			To:                tx.To(),
			Value:             tx.Value(),
			Input:             hexEncode(tx.Data()),
			Status:            r.Status,
			Gas:               r.GasUsed,
			EffectiveGasPrice: r.EffectiveGasPrice,
			ContractAddress:   contract,
			Logs:              events(r.Logs),
			Decoded:           decode(r.Logs, signatures),
		}
	}

	return &Block{
		Number:       block.Number(),
		Hash:         block.Hash(),
		Timestamp:    block.Time(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		BaseFee:      block.BaseFee(),
		Transactions: txs,
	}, nil
}

func (fx *Fx) blockReceipts(number *big.Int) ([]*types.Receipt, error) {
	var receipts []*types.Receipt
	arg := "latest"
	if number != nil {
		arg = fmt.Sprintf("0x%x", number)
	}
	if err := fx.Rpc.CallContext(fx.Context, &receipts, "eth_getBlockReceipts", arg); err != nil {
		return nil, fmt.Errorf("block receipts: %w", err)
	}
	return receipts, nil
}

func events(raw []*types.Log) []*Event {
	out := make([]*Event, len(raw))
	for i, l := range raw {
		topics := make([]string, len(l.Topics))
		for j, t := range l.Topics {
			topics[j] = t.Hex()
		}
		out[i] = &Event{
			Address: l.Address,
			Topics:  topics,
			Data:    hexEncode(l.Data),
		}
	}
	return out
}

func decode(raw []*types.Log, signatures map[string]string) []*Decoded {
	if len(raw) == 0 {
		return nil
	}
	var out []*Decoded
	for _, l := range raw {
		if len(l.Topics) == 0 {
			continue
		}
		sig, ok := signatures[l.Topics[0].Hex()]
		if !ok {
			continue
		}

		paramTypes := parseParams(sig)
		var vals []string
		topicIdx := 1
		dataOffset := 0
		data := l.Data

		for _, typ := range paramTypes {
			if topicIdx < len(l.Topics) && isIndexable(typ) {
				vals = append(vals, decodeWord(l.Topics[topicIdx].Bytes(), typ))
				topicIdx++
			} else if dataOffset+32 <= len(data) {
				vals = append(vals, decodeWord(data[dataOffset:dataOffset+32], typ))
				dataOffset += 32
			}
		}

		out = append(out, &Decoded{
			Contract:  l.Address,
			Signature: sig,
			Values:    strings.Join(vals, ","),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseParams(sig string) []string {
	start := strings.Index(sig, "(")
	end := strings.LastIndex(sig, ")")
	if start < 0 || end <= start+1 {
		return nil
	}
	inner := sig[start+1 : end]
	var parts []string
	depth := 0
	current := 0
	for i, c := range inner {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(inner[current:i]))
				current = i + 1
			}
		}
	}
	parts = append(parts, strings.TrimSpace(inner[current:]))
	return parts
}

func isIndexable(typ string) bool {
	switch {
	case typ == "address", typ == "bool":
		return true
	case strings.HasPrefix(typ, "uint"), strings.HasPrefix(typ, "int"):
		return true
	case strings.HasPrefix(typ, "bytes") && !strings.HasSuffix(typ, "[]"):
		return typ != "bytes" // dynamic bytes is hashed
	default:
		return false
	}
}

func decodeWord(b []byte, typ string) string {
	switch {
	case typ == "address":
		return common.BytesToAddress(b).Hex()
	case typ == "bool":
		if len(b) > 0 && b[len(b)-1] != 0 {
			return "true"
		}
		return "false"
	case strings.HasPrefix(typ, "uint"):
		return new(big.Int).SetBytes(b).String()
	case strings.HasPrefix(typ, "int"):
		v := new(big.Int).SetBytes(b)
		if len(b) > 0 && b[0]&0x80 != 0 {
			v.Sub(v, new(big.Int).Lsh(big.NewInt(1), uint(len(b)*8)))
		}
		return v.String()
	default:
		return "0x" + hex.EncodeToString(b)
	}
}

type signatureResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Event map[string][]struct {
			Name string `json:"name"`
		} `json:"event"`
	} `json:"result"`
}

func lookupSignatures(hashes map[string]struct{}) map[string]string {
	resolved := make(map[string]string, len(hashes))
	if len(hashes) == 0 {
		return resolved
	}
	keys := make([]string, 0, len(hashes))
	for h := range hashes {
		keys = append(keys, h)
	}
	resp, err := http.Get(fmt.Sprintf("%s?event=%s", signatureDB, strings.Join(keys, ",")))
	if err != nil {
		return resolved
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resolved
	}
	var body signatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil || !body.OK {
		return resolved
	}
	for hash, sigs := range body.Result.Event {
		if len(sigs) > 0 && sigs[0].Name != "" {
			resolved[hash] = sigs[0].Name
		}
	}
	return resolved
}

func hexEncode(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}
