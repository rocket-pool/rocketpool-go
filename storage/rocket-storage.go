package storage

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nodeset-org/eth-utils/eth"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketStorage
type Storage struct {
	Contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new Storage contract binding
func NewStorage(client eth.IExecutionClient, rocketStorageAddress common.Address) (*Storage, error) {
	// Create a Contract for the underlying raw RocketStorage binding
	rsAbi, err := abi.JSON(strings.NewReader(RocketStorageABI))
	if err != nil {
		return nil, err
	}
	contract := &core.Contract{
		Contract: &eth.Contract{
			ContractImpl: bind.NewBoundContract(rocketStorageAddress, rsAbi, client, client, client),
			Address:      rocketStorageAddress,
			ABI:          &rsAbi,
		},
	}

	return &Storage{
		Contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get a boolean value
func (c *Storage) GetBool(mc *batch.MultiCaller, result_Out *bool, key common.Hash) {
	core.AddCall(mc, c.Contract, result_Out, "getBool", key)
}

// Get a uint value
func (c *Storage) GetUint(mc *batch.MultiCaller, result_Out **big.Int, key common.Hash) {
	core.AddCall(mc, c.Contract, result_Out, "getUint", key)
}

// Get an address
func (c *Storage) GetAddress(mc *batch.MultiCaller, address_Out *common.Address, contractName string) {
	key := crypto.Keccak256Hash([]byte("contract.address"), []byte(contractName))
	core.AddCall(mc, c.Contract, address_Out, "getAddress", key)
}

// Get an ABI
func (c *Storage) GetAbi(mc *batch.MultiCaller, abiEncoded_Out *string, contractName string) {
	key := crypto.Keccak256Hash([]byte("contract.abi"), []byte(contractName))
	core.AddCall(mc, c.Contract, abiEncoded_Out, "getString", key)
}

// Get the number of the block that Rocket Pool was deployed on
func (c *Storage) GetDeployBlock(mc *batch.MultiCaller, result_Out **big.Int) {
	deployBlockHash := crypto.Keccak256Hash([]byte("deploy.block"))
	c.GetUint(mc, result_Out, deployBlockHash)
}
