package trustednode

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAONodeTrustedActions
type DaoNodeTrustedActions struct {
	Contract *core.Contract
	rp       *rocketpool.RocketPool
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrustedActions contract binding
func NewDaoNodeTrustedActions(rp *rocketpool.RocketPool) (*DaoNodeTrustedActions, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrustedActions)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted actions contract: %w", err)
	}

	return &DaoNodeTrustedActions{
		Contract: contract,
		rp:       rp,
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for joining the Oracle DAO
func (c *DaoNodeTrustedActions) Join(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "actionJoin", opts)
}

// Get info for leaving the Oracle DAO
func (c *DaoNodeTrustedActions) Leave(rplBondRefundAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "actionLeave", opts, rplBondRefundAddress)
}

// Get info for making a challenge to an Oracle DAO member
func (c *DaoNodeTrustedActions) MakeChallenge(memberAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "actionChallengeMake", opts, memberAddress)
}

// Get info for deciding a challenge to an Oracle DAO member
func (c *DaoNodeTrustedActions) DecideChallenge(memberAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "actionChallengeDecide", opts, memberAddress)
}

// =============
// === Utils ===
// =============

// Returns the most recent block number that the number of trusted nodes changed since fromBlock
func (c *DaoNodeTrustedActions) GetLatestMemberCountChangedBlock(fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) (uint64, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.Contract.Address}
	topicFilter := [][]common.Hash{{
		c.Contract.ABI.Events["ActionJoined"].ID,
		c.Contract.ABI.Events["ActionLeave"].ID,
		c.Contract.ABI.Events["ActionKick"].ID,
		c.Contract.ABI.Events["ActionChallengeDecided"].ID,
	}}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(fromBlock)), nil, nil)
	if err != nil {
		return 0, err
	}

	for i := range logs {
		log := logs[len(logs)-i-1]
		if log.Topics[0] == c.Contract.ABI.Events["ActionChallengeDecided"].ID {
			values := make(map[string]interface{})
			// Decode the event
			if c.Contract.ABI.Events["ActionChallengeDecided"].Inputs.UnpackIntoMap(values, log.Data) != nil {
				return 0, err
			}
			if values["success"].(bool) {
				return log.BlockNumber, nil
			}
		} else {
			return log.BlockNumber, nil
		}
	}
	return fromBlock, nil
}
