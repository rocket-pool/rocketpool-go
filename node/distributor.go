package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/eth"

	"github.com/rocket-pool/rocketpool-go/v2/core"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNodeDistributorDelegate
type NodeDistributor struct {
	// The address of the distributor
	DistributorAddress common.Address

	// The address of the node that owns this distributor
	NodeAddress common.Address

	// The node share of the distributor's current balance
	NodeShare *core.SimpleField[*big.Int]

	// The user share of the distributor's current balance
	UserShare *core.SimpleField[*big.Int]

	// === Internal fields ===
	rp       *rocketpool.RocketPool
	contract *core.Contract
	txMgr    *eth.TransactionManager
}

// ====================
// === Constructors ===
// ====================

// Creates a new NodeDistributor contract binding
func NewNodeDistributor(rp *rocketpool.RocketPool, nodeAddress common.Address, distributorAddress common.Address) (*NodeDistributor, error) {
	// Create the contract
	contract, err := rp.MakeContract(rocketpool.ContractName_RocketNodeDistributorDelegate, distributorAddress)
	if err != nil {
		return nil, fmt.Errorf("error getting node distributor delegate contract for node %s at %s: %w", nodeAddress.Hex(), distributorAddress.Hex(), err)
	}

	return &NodeDistributor{
		NodeAddress:        nodeAddress,
		DistributorAddress: distributorAddress,
		NodeShare:          core.NewSimpleField[*big.Int](contract, "getNodeShare"),
		UserShare:          core.NewSimpleField[*big.Int](contract, "getUserShare"),

		rp:       rp,
		contract: contract,
		txMgr:    rp.GetTransactionManager(),
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for distributing the contract's balance to the rETH contract and the user
func (c *NodeDistributor) Distribute(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.contract.Contract, "distribute", opts)
}
