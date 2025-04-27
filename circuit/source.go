package circuit

import "fmt"

func (c *Circuit) Load() error {
	if err := c.Factory.Data.Source("peers", &c.Peers); err != nil {
		return fmt.Errorf("failed to load peers")
	}
	if err := c.Factory.Data.Source("tokens", &c.Tokens); err != nil {
		return fmt.Errorf("failed to load tokens")
	}
	fmt.Printf("%d tokensloaded\n", len(c.Tokens))
	fmt.Printf("%d peersloaded\n", len(c.Peers))
	return nil
}
