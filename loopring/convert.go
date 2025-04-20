package loopring

import (
	"encoding/json"
	"fmt"
)

func (l *Loopring) SwapToTx(swap Swap) Tx {
	return Tx{
		Zero:        swap.ZeroId,
		One:         swap.OneId,
		Value:       swap.ZeroValue,
		Token:       swap.ZeroToken,
		OneValue:    swap.OneValue,
		OneToken:    swap.OneToken,
		Fee:         swap.ZeroFee,
		OneFeeToken: swap.OneFee,
		Type:        "swap",
		Coordinates: swap.Coordinates,
		Index:       swap.Index,
	}
}

func (l *Loopring) TransferToTx(transfer Transfer) Tx {
	return Tx{
		Zero:        transfer.ZeroId,
		One:         transfer.One,
		Value:       transfer.Value,
		Token:       transfer.Token,
		Fee:         transfer.Fee,
		FeeToken:    transfer.FeeToken,
		Type:        "transfer",
		Coordinates: transfer.Coordinates,
		Index:       transfer.Index,
	}
}

func (l *Loopring) DepositToTx(dw DW) Tx {
	return Tx{
		Zero:        dw.Zero,
		One:         dw.One,
		Value:       dw.Value,
		Token:       dw.Token,
		Type:        "deposit",
		Coordinates: dw.Coordinates,
		Index:       dw.Index,
	}
}

func (l *Loopring) WithdrawToTx(dw DW) Tx {
	return Tx{
		Zero:        dw.Zero,
		One:         dw.One,
		Value:       dw.Value,
		Token:       dw.Token,
		Fee:         dw.Fee,
		FeeToken:    dw.FeeToken,
		Type:        "withdraw",
		Coordinates: dw.Coordinates,
		Index:       dw.Index,
	}
}

func (l *Loopring) AccountUpdateToTx(accountUpdate AccountUpdate) Tx {
	return Tx{
		Zero:        accountUpdate.ZeroId,
		Type:        "accountUpdate",
		Coordinates: accountUpdate.Coordinates,
		Index:       accountUpdate.Index,
	}
}

func (l *Loopring) AmmUpdateToTx(ammUpdate AmmUpdate) Tx {
	return Tx{
		Zero:        ammUpdate.ZeroId,
		Type:        "ammUpdate",
		Coordinates: ammUpdate.Coordinates,
		Index:       ammUpdate.Index,
	}
}

func (l *Loopring) MintToTx(mint Mint) Tx {
	return Tx{
		Zero:        mint.Zero,
		Value:       mint.Quantity,
		Token:       mint.NftAddress,
		Fee:         mint.Fee,
		FeeToken:    mint.FeeToken,
		Type:        "mint",
		Coordinates: mint.Coordinates,
		Index:       mint.Index,
	}
}

func (l *Loopring) NftDataToTx(nftData NftData) Tx {
	raw, err := json.Marshal(nftData)
	if err != nil {
		fmt.Printf("Error marshaling NftData to raw JSON: %v\n", err)
	}

	var rawMap map[string]any
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		fmt.Printf("Error unmarshaling NftData to map: %v\n", err)
	}

	delete(rawMap, "accountId")
	delete(rawMap, "txType")
	delete(rawMap, "coordinates")

	filteredRaw, err := json.Marshal(rawMap)
	if err != nil {
		fmt.Printf("Error marshaling filtered raw map to JSON: %v\n", err)
	}
	return Tx{
		Zero:        nftData.ZeroId,
		Type:        "nftData",
		Coordinates: nftData.Coordinates,
		Raw:         filteredRaw,
		Index:       nftData.Index,
	}
}
