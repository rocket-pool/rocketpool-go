package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNodeDeposit
type NodeDeposit struct {
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new NodeDeposit contract binding
func NewNodeDeposit(rp *rocketpool.RocketPool) (*NodeDeposit, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketNodeDeposit)
	if err != nil {
		return nil, fmt.Errorf("error getting node deposit contract: %w", err)
	}

	return &NodeDeposit{
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the amount of ETH in the node's deposit credit bank
func (c *NodeDeposit) GetNodeDepositCredit(mc *multicall.MultiCaller, nodeAddress common.Address, credit_Out **big.Int) {
	multicall.AddCall(mc, c.contract, credit_Out, "getNodeDepositCredit", nodeAddress)
}

// ====================
// === Transactions ===
// ====================

// Get info for making a node deposit
func (c *NodeDeposit) Deposit(bondAmount *big.Int, minimumNodeFee float64, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, salt *big.Int, expectedMinipoolAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "deposit", opts, bondAmount, eth.EthToWei(minimumNodeFee), validatorPubkey[:], validatorSignature[:], depositDataRoot, salt, expectedMinipoolAddress)
}

// Get info for making a node deposit by using the credit balance
func (c *NodeDeposit) DepositWithCredit(bondAmount *big.Int, minimumNodeFee float64, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, salt *big.Int, expectedMinipoolAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "depositWithCredit", opts, bondAmount, eth.EthToWei(minimumNodeFee), validatorPubkey[:], validatorSignature[:], depositDataRoot, salt, expectedMinipoolAddress)
}

// Get info for making a vacant minipool for solo staker migration
func (c *NodeDeposit) CreateVacantMinipool(bondAmount *big.Int, minimumNodeFee float64, validatorPubkey rptypes.ValidatorPubkey, salt *big.Int, expectedMinipoolAddress common.Address, currentBalance *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "createVacantMinipool", opts, bondAmount, eth.EthToWei(minimumNodeFee), validatorPubkey[:], salt, expectedMinipoolAddress, currentBalance)
}
