package circuit

import (
	"encoding/json"
	"fmt"
)

func (c *Circuit) Continue() error {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()

	if err := c.SourcePeers(); err != nil {
		return fmt.Errorf("failed to source peers: %w", err)
	}

	if err := c.SourceTokens(); err != nil {
		return fmt.Errorf("failed to source tokens: %w", err)
	}

	fmt.Printf("%d\n", len(c.Map))
	return nil
}

func (c *Circuit) SourceTokens() error {
	source, err := c.Factory.Redis.SMembers(c.Factory.Ctx, "tokens").Result()
	if err != nil {
		return fmt.Errorf("failed to load tokens from Redis: %w", err)
	}

	c.Tokens = make([]Token, 0, len(source))
	for _, i := range source {
		var token Token
		if err := json.Unmarshal([]byte(i), &token); err != nil {
			return fmt.Errorf("failed to unmarshal token: %w", err)
		}
		c.Tokens = append(c.Tokens, token)
	}

	for _, token := range c.Tokens {
		c.TokenMap[token.TokenId] = &token
		c.Map[token.Symbol] = &token
		c.Map[token.Address] = &token
		c.Map[token.LoopringID] = &token
	}
	return nil
}

func (c *Circuit) SourcePeers() error {
	source, err := c.Factory.Redis.SMembers(c.Factory.Ctx, "peers").Result()
	if err != nil {
		return fmt.Errorf("failed to load peers from Redis: %w", err)
	}

	c.Peers = make([]Peer, 0, len(source))
	for _, i := range source {
		var peer Peer
		if err := json.Unmarshal([]byte(i), &peer); err != nil {
			return fmt.Errorf("failed to unmarshal peer: %w", err)
		}
		c.Peers = append(c.Peers, peer)
	}

	for _, peer := range c.Peers {
		c.Map[peer.Address] = &peer
		c.Map[peer.ENS] = &peer
		c.Map[peer.LoopringENS] = &peer
		c.Map[peer.LoopringID] = &peer
	}
	return nil
}
