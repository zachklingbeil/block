package zero

import (
	"encoding/json"
	"fmt"
	"log"
)

func (z *Zero) CreateContractTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS contracts (
            contract_address TEXT PRIMARY KEY,
            abi JSONB NOT NULL
        );`
	_, err := z.ContractDb.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create contracts table: %w", err)
	}
	return nil
}

type compilationArtifacts struct {
	ABI json.RawMessage `json:"abi"`
}

// PopulateContracts queries the sourcify database for all verified contracts
// on chain_id=1 and inserts their address + ABI into the contractdb contracts table.
func (z *Zero) PopulateContracts() (int, error) {
	query := `
        SELECT cd.address, cc.compilation_artifacts
        FROM contract_deployments cd
        JOIN contracts c ON c.id = cd.contract_id
        JOIN verified_contracts vc ON vc.deployment_id = cd.id
        JOIN compiled_contracts cc ON cc.id = vc.compilation_id
        WHERE cd.chain_id = 1
    `

	rows, err := z.Sourcify.Query(query)
	if err != nil {
		return 0, fmt.Errorf("query sourcify: %w", err)
	}
	defer rows.Close()

	insert := `
        INSERT INTO contracts (contract_address, abi)
        VALUES ($1, $2)
        ON CONFLICT (contract_address) DO NOTHING
    `

	count := 0
	for rows.Next() {
		var addrBytes []byte
		var artifactsRaw []byte

		if err := rows.Scan(&addrBytes, &artifactsRaw); err != nil {
			log.Printf("scan row: %v", err)
			continue
		}

		var artifacts compilationArtifacts
		if err := json.Unmarshal(artifactsRaw, &artifacts); err != nil {
			log.Printf("unmarshal artifacts for %x: %v", addrBytes, err)
			continue
		}

		if len(artifacts.ABI) == 0 {
			continue
		}

		// Convert raw address bytes to 0x-prefixed hex string
		contractAddr := fmt.Sprintf("0x%x", addrBytes)

		if _, err := z.ContractDb.Exec(insert, contractAddr, artifacts.ABI); err != nil {
			log.Printf("insert %s: %v", contractAddr, err)
			continue
		}
		count++
	}

	if err := rows.Err(); err != nil {
		return count, fmt.Errorf("rows iteration: %w", err)
	}

	return count, nil
}
