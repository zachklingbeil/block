package fx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	abiCache   = make(map[common.Address]*abi.ABI)
	abiCacheMu sync.RWMutex
)

func (fx *Fx) fetchABIs(contracts map[common.Address]struct{}) map[common.Address]*abi.ABI {
	out := make(map[common.Address]*abi.ABI, len(contracts))

	var missing []common.Address
	abiCacheMu.RLock()
	for addr := range contracts {
		if a, ok := abiCache[addr]; ok {
			out[addr] = a
		} else {
			missing = append(missing, addr)
		}
	}
	abiCacheMu.RUnlock()

	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, addr := range missing {
		wg.Add(1)
		go func(addr common.Address) {
			defer wg.Done()
			a := fx.sourcifyABI(addr)
			if a == nil {
				return
			}
			mu.Lock()
			out[addr] = a
			mu.Unlock()

			abiCacheMu.Lock()
			abiCache[addr] = a
			abiCacheMu.Unlock()
		}(addr)
	}
	wg.Wait()

	return out
}

func (fx *Fx) sourcifyABI(addr common.Address) *abi.ABI {
	url := fmt.Sprintf("http://sourcify:5555/v2/contract/1/%s?fields=abi", addr.Hex())
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil
	}

	var result struct {
		ABI json.RawMessage `json:"abi"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result.ABI) == 0 {
		return nil
	}

	parsed, err := abi.JSON(strings.NewReader(string(result.ABI)))
	if err != nil {
		return nil
	}
	return &parsed
}

func (fx *Fx) decodeInput(abis map[common.Address]*abi.ABI, to common.Address, data []byte) (string, map[string]any) {
	a, ok := abis[to]
	if !ok || len(data) < 4 {
		return "", nil
	}

	method, err := a.MethodById(data[:4])
	if err != nil {
		return "", nil
	}

	args := make(map[string]any)
	if len(data) > 4 {
		vals, err := method.Inputs.Unpack(data[4:])
		if err == nil {
			for i, input := range method.Inputs {
				args[input.Name] = vals[i]
			}
		}
	}
	return method.Name, args
}

func (fx *Fx) decodeLog(abis map[common.Address]*abi.ABI, log *types.Log) (string, map[string]any) {
	a, ok := abis[log.Address]
	if !ok || len(log.Topics) == 0 {
		return "", nil
	}

	event, err := a.EventByID(log.Topics[0])
	if err != nil {
		return "", nil
	}

	args := make(map[string]any)

	indexed := make([]abi.Argument, 0)
	nonIndexed := make([]abi.Argument, 0)
	for _, input := range event.Inputs {
		if input.Indexed {
			indexed = append(indexed, input)
		} else {
			nonIndexed = append(nonIndexed, input)
		}
	}

	for i, input := range indexed {
		if i+1 < len(log.Topics) {
			val, err := abi.Arguments{input}.Unpack(log.Topics[i+1].Bytes())
			if err == nil && len(val) > 0 {
				args[input.Name] = val[0]
			}
		}
	}

	if len(log.Data) > 0 {
		vals, err := abi.Arguments(nonIndexed).Unpack(log.Data)
		if err == nil {
			for i, input := range nonIndexed {
				args[input.Name] = vals[i]
			}
		}
	}

	return event.Name, args
}
