package circuit

import (
	"fmt"
)

func (c *Circuit) Load() error {
	peersAny := c.Factory.Data.Source("peers")
	peers := make([]Peer, len(peersAny))
	for i, p := range peersAny {
		peer, ok := p.(Peer)
		if !ok {
			return fmt.Errorf("invalid type assertion for peer at index %d", i)
		}
		peers[i] = peer
	}
	c.Peers = peers

	tokensAny := c.Factory.Data.Source("tokens")
	tokens := make([]Token, len(tokensAny))
	for i, t := range tokensAny {
		token, ok := t.(Token)
		if !ok {
			return fmt.Errorf("invalid type assertion for token at index %d", i)
		}
		tokens[i] = token
	}
	c.Tokens = tokens

	fmt.Printf("%d tokens loaded\n", len(c.Tokens))
	fmt.Printf("%d peers loaded\n", len(c.Peers))
	return nil
}
