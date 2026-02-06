package ethereum

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// )

// type Signature struct {
// 	Hex  string `json:"hex"`
// 	Text string `json:"text"`
// }

// func (e *Ethereum) LoadSignatures() error {
// 	keys := []string{"signature", "event_signature"}
// 	m := make(map[string]string)
// 	for _, redisKey := range keys {
// 		source, err := e.Factory.Data.RB.SMembers(e.Factory.Ctx, redisKey).Result()
// 		if err != nil {
// 			return fmt.Errorf("failed to fetch %s from Redis set: %v", redisKey, err)
// 		}
// 		for _, sigJSON := range source {
// 			var sig Signature
// 			if err := json.Unmarshal([]byte(sigJSON), &sig); err != nil {
// 				log.Printf("Skipping invalid %s: %v (data: %s)", redisKey, err, sigJSON)
// 				continue
// 			}
// 			m[sig.Hex] = sig.Text
// 		}
// 	}
// 	e.Factory.Rw.Lock()
// 	e.Signature = m
// 	e.Factory.Rw.Unlock()
// 	return nil
// }

// // Concurrent read method for both function and event signatures
// func (e *Ethereum) GetSignature(hex string) (string, bool) {
// 	e.Factory.Rw.RLock()
// 	defer e.Factory.Rw.RUnlock()
// 	val, ok := e.Signature[hex]
// 	return val, ok
// }
