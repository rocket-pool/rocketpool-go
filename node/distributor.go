package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
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
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for distributing the contract's balance to the rETH contract and the user
func (c *NodeDistributor) Distribute(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "distribute", opts)
}
