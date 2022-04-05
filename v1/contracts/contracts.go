package contracts

import (
	"encoding/json"

	"github.com/algorand/go-algorand-sdk/crypto"

	tTypes "github.com/synycboom/tinyman-go-sdk/types"
	tUtils "github.com/synycboom/tinyman-go-sdk/utils"
)

//go:generate ./bundle_asc_json.sh

var asc tTypes.ASC

func init() {
	if err := json.Unmarshal(ascJson, &asc); err != nil {
		panic(err)
	}
}

// PoolLogicSigAccount creates a logic signature account of the pool
func PoolLogicSigAccount(validatorAppID, asset1ID, asset2ID uint64) (*crypto.LogicSigAccount, error) {
	if asset2ID > asset1ID {
		asset1ID, asset2ID = asset2ID, asset1ID
	}

	program, err := tUtils.Program(asc.Contracts.PoolLogicSig.Logic, map[string]uint64{
		"validator_app_id": validatorAppID,
		"asset_id_1":       asset1ID,
		"asset_id_2":       asset2ID,
	})
	if err != nil {
		return nil, err
	}

	poolAccount := crypto.MakeLogicSigAccountEscrow(program, nil)

	return &poolAccount, nil
}
