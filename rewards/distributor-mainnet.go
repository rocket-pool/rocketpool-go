package rewards

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/eth"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketMerkleDistributorMainnet
type MerkleDistributorMainnet struct {
	// === Internal fields ===
	rp    *rocketpool.RocketPool
	rmdm  *core.Contract
	txMgr *eth.TransactionManager
}

// ====================
// === Constructors ===
// ====================

// Creates a new MerkleDistributorMainnet contract binding
func NewMerkleDistributorMainnet(rp *rocketpool.RocketPool) (*MerkleDistributorMainnet, error) {
	// Create the contract
	rmdm, err := rp.GetContract(rocketpool.ContractName_RocketMerkleDistributorMainnet)
	if err != nil {
		return nil, fmt.Errorf("error getting merkle distributor mainnet contract: %w", err)
	}

	return &MerkleDistributorMainnet{
		rp:    rp,
		rmdm:  rmdm,
		txMgr: rp.GetTransactionManager(),
	}, nil
}

// =============
// === Calls ===
// =============

// Check if the given node has already claimed rewards for the given interval
func (c *MerkleDistributorMainnet) HasNodeClaimedRewards(index *big.Int, claimerAddress common.Address, claimed_Out *bool, mc *batch.MultiCaller) {
	core.AddCall(mc, c.rmdm, claimed_Out, "isClaimed", index, claimerAddress)
}

// Get the Merkle root for an interval
func (c *MerkleDistributorMainnet) GetMerkleRoot(interval *big.Int, root_Out *common.Hash, mc *batch.MultiCaller) {
	core.AddCall(mc, c.rmdm, root_Out, "merkleRoots", interval)
}

// ====================
// === Transactions ===
// ====================

// Get info for claiming rewards
func (c *MerkleDistributorMainnet) Claim(address common.Address, indices []*big.Int, amountRPL []*big.Int, amountETH []*big.Int, merkleProofs [][]common.Hash, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.rmdm.Contract, "claim", opts, address, indices, amountRPL, amountETH, merkleProofs)
}

// Get info for claiming and restaking rewards
func (c *MerkleDistributorMainnet) ClaimAndStake(address common.Address, indices []*big.Int, amountRPL []*big.Int, amountETH []*big.Int, merkleProofs [][]common.Hash, stakeAmount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.rmdm.Contract, "claimAndStake", opts, address, indices, amountRPL, amountETH, merkleProofs, stakeAmount)
}
