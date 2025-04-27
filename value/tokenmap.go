package value

// type Token struct {
// 	Token      string `json:"token"`
// 	TokenId    int64  `json:"tokenId"`
// 	LoopringID int64  `json:"loopringId,omitempty"`
// 	Decimals   int64  `json:"decimals"`
// 	Address    string `json:"address"`
// }

// func NewTokens(factory *factory.Factory) []One {
// 	t := &Tokens{
// 		Slice: make([]Token, 270),
// 		Map:   make(map[int64]*Token),
// 	}

// 	source, err := factory.Data.RB.SMembers(factory.Ctx, "tokens").Result()
// 	if err != nil {
// 		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
// 	}

// 	for _, tokenJSON := range source {
// 		var token Token
// 		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
// 			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
// 			continue
// 		}
// 		t.Slice = append(t.Slice, token)
// 		t.ID[token.TokenId] = &t.Slice[len(t.Slice)-1]
// 	}
// 	factory.State.Add("tokens", len(t.Slice))
// 	factory.State.Add("tokenId", len(t.ID))
// 	return t
// }
