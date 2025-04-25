package circuit

import (
	"fmt"
	"strings"

	"github.com/zachklingbeil/factory"
)

type Circuit struct {
	Factory        *factory.Factory
	Map            map[string]any
	TokenMap       map[int64]*Token
	Tokens         []Token
	Peers          []Peer
	LoopringApiKey string
}

func NewCircuit(factory *factory.Factory) *Circuit {
	circuit := &Circuit{
		Factory:  factory,
		Map:      make(map[string]any),
		TokenMap: make(map[int64]*Token),
	}

	return circuit
}
func (c *Circuit) Get(key any) any {
	c.Factory.Rw.Lock()
	defer c.Factory.Rw.Unlock()

	switch k := key.(type) {
	case string:
		if value, ok := c.Map[strings.ToLower(k)]; ok {
			return value
		}
		fmt.Printf("Key not found: %v (string)\n", k)
	case int:
		int64Key := int64(k)
		if value, ok := c.TokenMap[int64Key]; ok {
			return value
		}
		fmt.Printf("Key not found: %v (int)\n", k)
	default:
		fmt.Printf("Unsupported key type: %T\n", key)
	}

	return nil
}

func (c *Circuit) Add(key any, value any) {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()

	switch k := key.(type) {
	case string:
		c.Map[strings.ToLower(k)] = value
	case int:
		int64Key := int64(k)
		c.TokenMap[int64Key] = value.(*Token)
	default:
		fmt.Printf("Unsupported key type: %T\n", key)
		return
	}
}
