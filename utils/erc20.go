package utils

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	Erc20AbiString string = `[
		{
			"constant": true,
			"inputs": [],
			"name": "name",
			"outputs": [
			{
				"name": "",
				"type": "string"
			}
			],
			"payable": false,
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "decimals",
			"outputs": [
			{
				"name": "",
				"type": "uint8"
			}
			],
			"payable": false,
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [
			{
				"name": "_owner",
				"type": "address"
			}
			],
			"name": "balanceOf",
			"outputs": [
			{
				"name": "balance",
				"type": "uint256"
			}
			],
			"payable": false,
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "symbol",
			"outputs": [
			{
				"name": "",
				"type": "string"
			}
			],
			"payable": false,
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
			{
				"name": "_to",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
			],
			"name": "transfer",
			"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
			],
			"payable": false,
			"type": "function"
		}
	]`
)

// Global container for the parsed ABI above
var erc20Abi *abi.ABI

// ===============
// === Structs ===
// ===============

// Binding for ERC20 contracts
type Erc20Contract struct {
	Details  Erc20ContractDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for ERC20 contracts
type Erc20ContractDetails struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}

// ====================
// === Constructors ===
// ====================

// Creates a contract wrapper for the ERC20 at the given address
func NewErc20Contract(rp *rocketpool.RocketPool, address common.Address, client core.ExecutionClient, opts *bind.CallOpts) (*Erc20Contract, error) {
	// Parse the ABI
	if erc20Abi == nil {
		abiParsed, err := abi.JSON(strings.NewReader(Erc20AbiString))
		if err != nil {
			return nil, fmt.Errorf("error parsing ERC20 ABI: %w", err)
		}
		erc20Abi = &abiParsed
	}

	// Create contract
	contract := &core.Contract{
		Contract: bind.NewBoundContract(address, *erc20Abi, client, client, client),
		Address:  &address,
		ABI:      erc20Abi,
		Client:   client,
	}

	// Create the wrapper
	wrapper := &Erc20Contract{
		Details:  Erc20ContractDetails{},
		contract: contract,
	}

	// Get the details
	err := rp.Query(func(mc *multicall.MultiCaller) error {
		multicall.AddCall(mc, contract, &wrapper.Details.Name, "name")
		multicall.AddCall(mc, contract, &wrapper.Details.Symbol, "symbol")
		multicall.AddCall(mc, contract, &wrapper.Details.Decimals, "decimals")
		return nil
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting ERC-20 details of token %s: %w", address.Hex(), err)
	}

	return wrapper, nil
}

// =============
// === Calls ===
// =============

// Get the token balance for an address
func (c *Erc20Contract) BalanceOf(mc *multicall.MultiCaller, balance_Out **big.Int, address common.Address) {
	multicall.AddCall(mc, c.contract, balance_Out, "balanceOf", address)
}

// ====================
// === Transactions ===
// ====================

// Get info for transferring the ERC20 to another address
func (c *Erc20Contract) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "transfer", opts, to, amount)
}
