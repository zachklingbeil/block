package token

// type Tokens struct {
//     Factory  *factory.Factory
//     TokenMap map[int64]*Token // Use int64 as the key for fast lookups
// }

// func NewTokens(factory *factory.Factory) *Tokens {
//     return &Tokens{
//         Factory:  factory,
//         TokenMap: make(map[int64]*Token),
//     }
// }

// func (t *Tokens) LoadTokens(ctx context.Context, redis *redis.Client) error {
//     hashKey := "token"
//     source, err := redis.HGetAll(ctx, hashKey).Result()
//     if err != nil {
//         return fmt.Errorf("failed to fetch tokens from Redis hash: %v", err)
//     }

//     t.TokenMap = make(map[int64]*Token)

//     for _, tokenJSON := range source {
//         var token Token
//         if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
//             log.Printf("Skipping invalid token: %v (data: %s)", err, tokenJSON)
//             continue
//         }
//         t.TokenMap[token.TokenInt] = &token
//     }
//     return nil
// }
