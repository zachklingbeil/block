package fx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/timefactoryio/block/zero"
	"golang.org/x/time/rate"
)

type Fx struct {
	*zero.Zero
	Chain     *params.ChainConfig
	Contracts map[common.Address]*abi.ABI
	Limiter   *rate.Limiter
}

func Init(password string) *Fx {
	fx := &Fx{
		Zero:      zero.Init(password),
		Chain:     params.MainnetChainConfig,
		Contracts: make(map[common.Address]*abi.ABI),
		Limiter:   rate.NewLimiter(rate.Limit(3), 1),
	}
	return fx
}

func (fx *Fx) GetABI(addr common.Address) (*abi.ABI, bool) {
	fx.RLock()
	a, ok := fx.Contracts[addr]
	fx.RUnlock()
	if ok {
		return a, true
	}
	if err := fx.FetchLimited(addr); err != nil {
		return nil, false
	}
	fx.RLock()
	a, ok = fx.Contracts[addr]
	fx.RUnlock()
	return a, ok
}

func (fx *Fx) Load(addr common.Address, raw string) error {
	a, err := abi.JSON(strings.NewReader(raw))
	if err != nil {
		return fmt.Errorf("abi %s: %w", addr.Hex(), err)
	}
	fx.Lock()
	fx.Contracts[addr] = &a
	fx.Unlock()
	return nil
}

func (fx *Fx) Fetch(addr common.Address) error {
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
	fx.Contracts[addr] = &a
	fx.Unlock()
	return nil
}

func (fx *Fx) FetchLimited(addr common.Address) error {
	if err := fx.Limiter.Wait(fx.Context); err != nil {
		return fmt.Errorf("rate limit %s: %w", addr.Hex(), err)
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
	fx.Contracts[addr] = &a
	fx.Unlock()
	return nil
}

func (fx *Fx) Test() error {
	block, err := fx.Block(nil)
	if err != nil {
		return fmt.Errorf("Block: %w", err)
	}

	output, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		return fmt.Errorf("MarshalIndent: %w", err)
	}

	if err := os.WriteFile("../output/block.json", output, 0644); err != nil {
		return fmt.Errorf("WriteFile: %w", err)
	}

	fmt.Println("Block written to output/block.json")
	return nil
}
