package tests

import (
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// Mint old RPL for unit testing
func MintLegacyRpl(rp *rocketpool.RocketPool, ownerAccount *Account, toAccount *Account, amount *big.Int) (*core.TransactionInfo, error) {
	fsrpl, err := rp.GetContract(rocketpool.ContractName_RocketTokenRPLFixedSupply)
	if err != nil {
		return nil, fmt.Errorf("error creating legacy RPL contract: %w", err)
	}

	return core.NewTransactionInfo(fsrpl, "mint", ownerAccount.Transactor, toAccount.Address, amount)
}
