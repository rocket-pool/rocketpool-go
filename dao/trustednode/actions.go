package trustednode

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAONodeTrustedActions
type DaoNodeTrustedActions struct {
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrustedActions contract binding
func NewDaoNodeTrustedActions(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*DaoNodeTrustedActions, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketDAONodeTrustedActions", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted actions contract: %w", err)
	}

	return &DaoNodeTrustedActions{
		rp:       rp,
		contract: contract,
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for joining the Oracle DAO
func (c *DaoNodeTrustedActions) Join(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "actionJoin", opts)
}

// Get info for leaving the Oracle DAO
func (c *DaoNodeTrustedActions) Leave(rplBondRefundAddress common.Address, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "actionLeave", opts, rplBondRefundAddress)
}

// Get info for making a challenge to an Oracle DAO member
func (c *DaoNodeTrustedActions) MakeChallenge(memberAddress common.Address, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "actionChallengeMake", opts, memberAddress)
}

// Get info for deciding a challenge to an Oracle DAO member
func (c *DaoNodeTrustedActions) DecideChallenge(memberAddress common.Address, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "actionChallengeDecide", opts, memberAddress)
}
