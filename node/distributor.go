package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNodeDistributorDelegate
type NodeDistributor struct {
	Details  NodeDistributorDetails
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Details for RocketNodeDistributorDelegate
type NodeDistributorDetails struct {
	NodeAddress        common.Address `json:"nodeAddress"`
	DistributorAddress common.Address `json:"distributorAddress"`
	NodeShare          *big.Int       `json:"nodeShare"`
	UserShare          *big.Int       `json:"userShare"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NodeDistributor contract binding
func NewNodeDistributor(rp *rocketpool.RocketPool, nodeAddress common.Address, distributorAddress common.Address, opts *bind.CallOpts) (*NodeDistributor, error) {
	// Create the contract
	contract, err := rp.MakeContract("rocketNodeDistributorDelegate", distributorAddress, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node distributor delegate contract for node %s at %s: %w", nodeAddress.Hex(), distributorAddress.Hex(), err)
	}

	return &NodeDistributor{
		Details: NodeDistributorDetails{
			NodeAddress:        nodeAddress,
			DistributorAddress: distributorAddress,
		},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Gets the node share of the distributor's current balance
func (c *NodeDistributor) GetNodeShare(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.NodeShare, "getNodeShare")
}

// Gets the user share of the distributor's current balance
func (c *NodeDistributor) GetUserShare(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.UserShare, "getUserShare")
}

// Get all basic details
func (c *NodeDistributor) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetNodeShare(mc)
	c.GetUserShare(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for distributing the contract's balance to the rETH contract and the user
func (c *NodeDistributor) PlaceBid(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "distribute", opts)
}
