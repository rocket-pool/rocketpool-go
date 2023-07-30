package storage

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketStorage
type Storage struct {
	contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new Storage contract binding
func NewStorage(client core.ExecutionClient, rocketStorageAddress common.Address) (*Storage, error) {
	// Create a Contract for the underlying raw RocketStorage binding
	rsAbi, err := abi.JSON(strings.NewReader(RocketStorageABI))
	if err != nil {
		return nil, err
	}
	contract := &core.Contract{
		Contract: bind.NewBoundContract(rocketStorageAddress, rsAbi, client, client, client),
		Address:  &rocketStorageAddress,
		ABI:      &rsAbi,
		Client:   client,
	}

	return &Storage{
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get a boolean value
func (c *Storage) GetBool(mc *multicall.MultiCaller, result_Out *bool, key common.Hash) {
	multicall.AddCall(mc, c.contract, result_Out, "getBool", key)
}

// Get a uint value
func (c *Storage) GetUint(mc *multicall.MultiCaller, result_Out **big.Int, key common.Hash) {
	multicall.AddCall(mc, c.contract, result_Out, "getUint", key)
}

// Get an address
func (c *Storage) GetAddress(mc *multicall.MultiCaller, address_Out *common.Address, contractName string) {
	key := crypto.Keccak256Hash([]byte("contract.address"), []byte(contractName))
	multicall.AddCall(mc, c.contract, address_Out, "getAddress", key)
}

// Get an ABI
func (c *Storage) GetAbi(mc *multicall.MultiCaller, abiEncoded_Out *string, contractName string) {
	key := crypto.Keccak256Hash([]byte("contract.abi"), []byte(contractName))
	multicall.AddCall(mc, c.contract, abiEncoded_Out, "getString", key)
}

// Get a node's withdrawal address
func (c *Storage) GetNodeWithdrawalAddress(mc *multicall.MultiCaller, result_Out *common.Address, nodeAddress common.Address) {
	multicall.AddCall(mc, c.contract, result_Out, "getNodeWithdrawalAddress", nodeAddress)
}

// Get a node's pending withdrawal address
func (c *Storage) GetNodePendingWithdrawalAddress(mc *multicall.MultiCaller, result_Out *common.Address, nodeAddress common.Address) {
	multicall.AddCall(mc, c.contract, result_Out, "getNodePendingWithdrawalAddress", nodeAddress)
}

// ====================
// === Transactions ===
// ====================

// Get info for setting a node's withdrawal address
func (c *Storage) SetWithdrawalAddress(nodeAddress common.Address, withdrawalAddress common.Address, confirm bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "setWithdrawalAddress", opts, nodeAddress, withdrawalAddress, confirm)
}

// Get info for confirming a node's withdrawal address
func (c *Storage) ConfirmWithdrawalAddress(nodeAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "confirmWithdrawalAddress", opts, nodeAddress)
}
