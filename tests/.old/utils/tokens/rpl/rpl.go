package rpl

import (
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
	"github.com/rocket-pool/rocketpool-go/v2/tests"
	"github.com/rocket-pool/rocketpool-go/v2/tokens"
)

// Mint an amount of RPL to an account
func MintRPL(rp *rocketpool.RocketPool, ownerAccount *tests.Account, toAccount *tests.Account, amount *big.Int) error {

	// Get RPL token contracts
	rpl, err := tokens.NewTokenRpl(rp)
	if err != nil {
		return fmt.Errorf("error getting RPL binding: %w", err)
	}
	legacyRpl, err := tokens.NewTokenRplFixedSupply(rp)
	if err != nil {
		return fmt.Errorf("error getting legacy RPL binding: %w", err)
	}

	// Mint, approve & swap fixed-supply RPL
	err = MintFixedSupplyRPL(rp, ownerAccount, toAccount, amount)
	if err != nil {
		return fmt.Errorf("error minting legacy RPL: %w", err)
	}

	// Approve legacy RPL usage
	txInfo, err := legacyRpl.Approve(*rpl.Contract.Address, amount, toAccount.Transactor)
	if err != nil {
		return fmt.Errorf("error getting approval info: %w", err)
	}
	if txInfo.SimError != "" {
		return fmt.Errorf("simulating approval failed: %s", txInfo.SimError)
	}

	if _, err := tokens.SwapFixedSupplyRPLForRPL(rp, amount, toAccount.GetTransactor()); err != nil {
		return err
	}

	// Return
	return nil

}

// Mint an amount of fixed-supply RPL to an account
func MintFixedSupplyRPL(rp *rocketpool.RocketPool, ownerAccount *tests.Account, toAccount *tests.Account, amount *big.Int) error {
	rocketTokenFixedSupplyRPL, err := rp.GetContract("rocketTokenRPLFixedSupply")
	if err != nil {
		return err
	}
	if _, err := rocketTokenFixedSupplyRPL.Transact(ownerAccount.GetTransactor(), "mint", toAccount.Address, amount); err != nil {
		return fmt.Errorf("Could not mint fixed-supply RPL tokens to %s: %w", toAccount.Address.Hex(), err)
	}
	return nil
}
