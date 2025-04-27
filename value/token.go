package value

// func NewToken(factory *factory.Factory) []One {

// 	source, err := factory.Data.RB.SMembers(factory.Ctx, "tokens").Result()
// 	if err != nil {
// 		log.Fatalf("Failed to fetch tokens from Redis: %v", err)
// 	}

// 	for _, tokenJSON := range source {
// 		var token One
// 		if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
// 			log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
// 			continue
// 		}
// 		t.Slice = append(t.Slice, token)
// 	}
// 	factory.State.Add("tokens", len(t.Slice))
// 	factory.State.Add("tokenId", len(t.ID))
// 	return t
// }
