package fx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type ABIEntry struct {
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Inputs    []ABIParam `json:"inputs"`
	Anonymous bool       `json:"anonymous"`
}

type ABIParam struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed"`
}

type EventABI struct {
	Name    string
	Sig     string
	Indexed []ABIParam
	Data    []ABIParam
}

type sourcifyResponse struct {
	Match   *string         `json:"match"`
	Address string          `json:"address"`
	ABI     json.RawMessage `json:"abi"`
}

// ContractABI fetches the ABI from the local Sourcify v2 API.
func (fx *Fx) ContractABI(addr common.Address) ([]ABIEntry, error) {
	url := fmt.Sprintf("http://sourcify:5555/server/v2/contract/1/%s?fields=abi", addr.Hex())
	resp, err := fx.Http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("sourcify request [%s]: %w", addr.Hex(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("contract not verified [%s]", addr.Hex())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sourcify status %d [%s]", resp.StatusCode, addr.Hex())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response [%s]: %w", addr.Hex(), err)
	}

	var result sourcifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response [%s]: %w", addr.Hex(), err)
	}

	if result.Match == nil || len(result.ABI) == 0 {
		return nil, fmt.Errorf("no abi [%s]", addr.Hex())
	}

	var abi []ABIEntry
	if err := json.Unmarshal(result.ABI, &abi); err != nil {
		return nil, fmt.Errorf("parse abi [%s]: %w", addr.Hex(), err)
	}

	return abi, nil
}

func EventSignature(e ABIEntry) string {
	types := make([]string, len(e.Inputs))
	for i, p := range e.Inputs {
		types[i] = p.Type
	}
	return fmt.Sprintf("%s(%s)", e.Name, strings.Join(types, ","))
}

func ParseEvents(abi []ABIEntry) map[common.Hash]*EventABI {
	events := make(map[common.Hash]*EventABI)
	for _, entry := range abi {
		if entry.Type != "event" {
			continue
		}

		sig := EventSignature(entry)
		topic0 := crypto.Keccak256Hash([]byte(sig))

		var indexed, data []ABIParam
		for _, p := range entry.Inputs {
			if p.Indexed {
				indexed = append(indexed, p)
			} else {
				data = append(data, p)
			}
		}

		events[topic0] = &EventABI{
			Name:    entry.Name,
			Sig:     sig,
			Indexed: indexed,
			Data:    data,
		}
	}
	return events
}

func decodeSlot(paramType string, slot common.Hash) string {
	switch {
	case paramType == "address":
		return common.BytesToAddress(slot.Bytes()).Hex()
	case paramType == "bool":
		if slot[31] == 1 {
			return "true"
		}
		return "false"
	case strings.HasPrefix(paramType, "uint"):
		return new(big.Int).SetBytes(slot.Bytes()).String()
	case strings.HasPrefix(paramType, "int"):
		v := new(big.Int).SetBytes(slot.Bytes())
		if slot[0]&0x80 != 0 {
			v.Sub(v, new(big.Int).Lsh(big.NewInt(1), 256))
		}
		return v.String()
	case strings.HasPrefix(paramType, "bytes"):
		return "0x" + hex.EncodeToString(slot.Bytes())
	default:
		return slot.Hex()
	}
}

// DecodeIndexed maps indexed parameter names to their decoded values from log topics.
// topics[0] is the event signature (topic0), topics[1:] are the indexed values.
func DecodeIndexed(event *EventABI, topics []common.Hash) map[string]string {
	result := make(map[string]string, len(event.Indexed))
	for i, param := range event.Indexed {
		if i+1 >= len(topics) {
			break
		}
		result[param.Name] = decodeSlot(param.Type, topics[i+1])
	}
	return result
}
