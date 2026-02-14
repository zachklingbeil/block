package fx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	ABI       *abi.ABI
	Templates map[[4]byte]*Template
}

type Template struct {
	Method   [4]byte
	Outcomes []Outcome
}

type Outcome struct {
	Status uint64         // 1 = success, 0 = fail
	Events []OutcomeEvent // events emitted on success
	Error  *OutcomeEvent  // revert on failure
	Count  uint64
}

type OutcomeEvent struct {
	Contract common.Address
	Selector [4]byte
	ParamMap map[string]string // event param → method input (nil for errors)
}

// GetContract returns the Contract for an address, fetching the ABI if needed.
func (fx *Fx) GetContract(addr common.Address) (*Contract, bool) {
	fx.RLock()
	c, ok := fx.Contracts[addr]
	fx.RUnlock()
	if ok {
		return c, true
	}
	if err := fx.Fetch(addr, true); err != nil {
		return nil, false
	}
	fx.RLock()
	c, ok = fx.Contracts[addr]
	fx.RUnlock()
	if !ok {
		return nil, false
	}
	return c, true
}

func (fx *Fx) Load(addr common.Address, raw string) error {
	a, err := abi.JSON(strings.NewReader(raw))
	if err != nil {
		return fmt.Errorf("abi %s: %w", addr.Hex(), err)
	}
	fx.Lock()
	fx.Contracts[addr] = &Contract{ABI: &a, Templates: make(map[[4]byte]*Template)}
	fx.Unlock()
	return nil
}

func (fx *Fx) Fetch(addr common.Address, limited bool) error {
	if limited {
		if err := fx.Limiter.Wait(fx.Context); err != nil {
			return fmt.Errorf("rate limit %s: %w", addr.Hex(), err)
		}
	}

	url := fmt.Sprintf("https://sourcify.dev/server/v2/contract/1/%s?fields=abi", addr.Hex())

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", addr.Hex(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch %s: status %d", addr.Hex(), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read %s: %w", addr.Hex(), err)
	}

	var result struct {
		ABI json.RawMessage `json:"abi"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("unmarshal %s: %w", addr.Hex(), err)
	}

	a, err := abi.JSON(strings.NewReader(string(result.ABI)))
	if err != nil {
		return fmt.Errorf("parse abi %s: %w", addr.Hex(), err)
	}

	fx.Lock()
	fx.Contracts[addr] = &Contract{ABI: &a, Templates: make(map[[4]byte]*Template)}
	fx.Unlock()
	return nil
}
