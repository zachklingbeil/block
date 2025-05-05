package value

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
	byId      = "https://api3.loopring.io/api/v3/account?accountId=%s"
	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
)

func (v *Value) Format(address string) string {
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
		return address
	}
	return address
}

func (v *Value) input(url string, response any) error {
	data, err := v.Factory.Json.In(url, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}

func (v *Value) Refresh() {
	for i := range v.Peers {
		fmt.Printf("%d\n", i)
		peer := v.Peers[i]
		v.Format(peer.Address)
		v.Save(peer)
	}
}
