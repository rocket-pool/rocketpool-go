package rewards

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketMerkleDistributorMainnet
type MerkleDistributorMainnet struct {
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new MerkleDistributorMainnet contract binding
func NewMerkleDistributorMainnet(rp *rocketpool.RocketPool) (*MerkleDistributorMainnet, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketMerkleDistributorMainnet)
	if err != nil {
		return nil, fmt.Errorf("error getting merkle distributor mainnet contract: %w", err)
	}

	return &MerkleDistributorMainnet{
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Check if the given node has already claimed rewards for the given interval
func (c *MerkleDistributorMainnet) GetTotalRPLBalance(index *big.Int, claimerAddress common.Address, claimed_Out *bool, mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, claimed_Out, "isClaimed", index, claimerAddress)
}

// Get the Merkle root for an interval
func (c *MerkleDistributorMainnet) GetMerkleRoot(interval *big.Int, root_Out *common.Hash, mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, root_Out, "merkleRoots", interval)
}

// ====================
// === Transactions ===
// ====================

// Get info for claiming rewards
func (c *MerkleDistributorMainnet) Claim(address common.Address, indices []*big.Int, amountRPL []*big.Int, amountETH []*big.Int, merkleProofs [][]common.Hash, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "claim", opts, address, indices, amountRPL, amountETH, merkleProofs)
}

// Get info for claiming and restaking rewards
func (c *MerkleDistributorMainnet) ClaimAndStake(address common.Address, indices []*big.Int, amountRPL []*big.Int, amountETH []*big.Int, merkleProofs [][]common.Hash, stakeAmount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "claimAndStake", opts, address, indices, amountRPL, amountETH, merkleProofs, stakeAmount)
}
