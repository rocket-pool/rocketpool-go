package trustednode

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

const (
	// Contract names
	DaoNodeTrustedActions_ContractName string = "rocketDAONodeTrustedActions"

	// Transactions
	daoNodeTrustedActions_actionJoin            string = "actionJoin"
	daoNodeTrustedActions_actionLeave           string = "actionLeave"
	daoNodeTrustedActions_actionChallengeMake   string = "actionChallengeMake"
	daoNodeTrustedActions_actionChallengeDecide string = "actionChallengeDecide"
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
	contract, err := rp.GetContract(DaoNodeTrustedActions_ContractName, opts)
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
	return rocketpool.NewTransactionInfo(c.contract, daoNodeTrustedActions_actionJoin, opts)
}

// Get info for leaving the Oracle DAO
func (c *DaoNodeTrustedActions) Leave(rplBondRefundAddress common.Address, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, daoNodeTrustedActions_actionLeave, opts, rplBondRefundAddress)
}

// Get info for making a challenge to an Oracle DAO member
func (c *DaoNodeTrustedActions) MakeChallenge(memberAddress common.Address, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, daoNodeTrustedActions_actionChallengeMake, opts, memberAddress)
}

// Get info for deciding a challenge to an Oracle DAO member
func (c *DaoNodeTrustedActions) DecideChallenge(memberAddress common.Address, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, daoNodeTrustedActions_actionChallengeDecide, opts, memberAddress)
}
