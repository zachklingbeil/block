package universe

// // AddStructToMap now works with the Key interface
// func (u *Universe) AddStructToMap(key Key) {
// 	address := key.GetAddress()
// 	if !common.IsHexAddress(address) {
// 		log.Printf("Invalid Ethereum address: %s. Entry not added to the map.", address)
// 		return
// 	}

// 	u.Map[common.HexToAddress(address)] = &struct{}{}
// 	log.Printf("Added struct to map with key: %s", address)
// }
