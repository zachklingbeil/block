package circuit

import (
	"encoding/json"
	"fmt"
)

func (c *Circuit) Continue() error {
	c.Factory.Mu.Lock()
	defer c.Factory.Mu.Unlock()

	if err := c.SourceValues(); err != nil {
		return fmt.Errorf("failed to source values: %w", err)
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

func (c *Circuit) SourceValues() error {
	source, err := c.Factory.Redis.SMembers(c.Factory.Ctx, "value").Result()
	if err != nil {
		return fmt.Errorf("failed to load values from Redis: %w", err)
	}

	c.Values = make([]Value, 0, len(source))
	for _, i := range source {
		var value Value
		if err := json.Unmarshal([]byte(i), &value); err != nil {
			return fmt.Errorf("failed to unmarshal value: %w", err)
		}
		c.Map[value.Address] = &value
		c.Map[value.ENS] = &value
		c.Map[value.LoopringENS] = &value
		c.Map[value.LoopringID] = &value
		c.Map[value.Symbol] = &value
		c.Map[value.Address] = &value
		c.Map[value.LoopringID] = &value
		c.Values = append(c.Values, value)
	}

	return nil
}
