package value

// const (
// 	byAddress = "https://api3.loopring.io/api/v3/account?owner=%s"
// 	byId      = "https://api3.loopring.io/api/v3/account?accountId=%d"
// 	dotLoop   = "https://api3.loopring.io/api/wallet/v3/resolveName?owner=%s"
// )

// func (p *Peers) Format(address string) string {
// 	address = strings.ToLower(address)
// 	if strings.HasPrefix(address, "0x") || strings.HasSuffix(address, ".eth") {
// 		return address
// 	}
// 	return address
// }

// // hex -> .eth
// func (p *Peers) GetENS(peer *One) *One {
// 	if peer.ENS == "." || (peer.ENS != "" && peer.ENS != "!") {
// 		// Return immediately if ENS is already set or marked as checked
// 		return peer
// 	}

// 	ensName, err := ens.ReverseResolve(p.Factory.Eth, common.HexToAddress(peer.Address.(string)))
// 	if err != nil || ensName == "" {
// 		peer.ENS = "." // Mark as checked with no ENS
// 		return peer
// 	}
// 	peer.ENS = p.Format(ensName)
// 	return peer
// }

// // ENS -> hex
// func (p *Peers) GetAddress(peer *One) *One {
// 	address, err := ens.Resolve(p.Factory.Eth, peer.ENS.(string))
// 	if err != nil {
// 		peer.Address = peer.ENS
// 		return peer
// 	}
// 	peer.Address = p.Format(address.Hex())
// 	return peer
// }

// // hex -> LoopringENS [.loopring.eth] or "."
// func (p *Peers) GetLoopringENS(peer *One) *One {
// 	if peer.LoopringENS == "." || (peer.LoopringENS != "" && peer.LoopringENS != "!") {
// 		// Return immediately if LoopringENS is already set or marked as checked
// 		return peer
// 	}

// 	url := fmt.Sprintf(dotLoop, peer.Address)
// 	var response struct {
// 		Loopring string `json:"data"`
// 	}

// 	data, err := p.Factory.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
// 	if err != nil || json.Unmarshal(data, &response) != nil || response.Loopring == "" {
// 		peer.LoopringENS = "." // Mark as checked with no valid LoopringENS
// 		return peer
// 	}

// 	peer.LoopringENS = p.Format(response.Loopring)
// 	return peer
// }

// // hex -> LoopringId or -1
// func (p *Peers) GetLoopringID(peer *One) *One {
// 	if peer.LoopringID == "." || (peer.LoopringID != "" && peer.LoopringID != "!") {
// 		return peer
// 	}
// 	url := fmt.Sprintf(byAddress, peer.Address)
// 	var response struct {
// 		ID int64 `json:"accountId"`
// 	}

// 	data, err := p.Factory.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
// 	if err != nil || json.Unmarshal(data, &response) != nil || response.ID == 0 {
// 		peer.LoopringID = "."
// 		return peer
// 	}
// 	peer.LoopringID = strconv.FormatInt(response.ID, 10)
// 	return peer
// }

// // LoopringId -> hex
// func (p *Peers) GetLoopringAddress(peer *One) *One {
// 	if peer.Address == "." || (peer.Address != "" && peer.Address != "!") {
// 		return peer
// 	}
// 	accountID, err := strconv.Atoi(peer.LoopringID.(string))
// 	if err != nil {
// 		peer.Address = "." // Mark as checked with no valid Address
// 		return peer
// 	}

// 	url := fmt.Sprintf(byId, accountID)
// 	var response struct {
// 		Address string `json:"owner"`
// 	}

// 	data, err := p.Factory.Json.In(url, os.Getenv("LOOPRING_API_KEY"))
// 	if err != nil || json.Unmarshal(data, &response) != nil || response.Address == "" {
// 		peer.Address = "." // Mark as checked with no valid Address
// 		return peer
// 	}
// 	peer.Address = p.Format(response.Address)
// 	return peer
// }
