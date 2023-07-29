package rocketpool

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/storage"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// Rocket Pool contract manager
type RocketPool struct {
	Client                core.ExecutionClient
	Storage               *storage.Storage
	MulticallAddress      *common.Address
	BalanceBatcherAddress *common.Address
	VersionManager        *VersionManager
}

// Create new contract manager
func NewRocketPool(client core.ExecutionClient, rocketStorageAddress common.Address, multicallAddress common.Address, balanceBatcherAddress common.Address) (*RocketPool, error) {

	// Create the RocketStorage binding
	storage, err := storage.NewStorage(client, rocketStorageAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating rocket storage binding: %w", err)
	}

	// Create and return
	rp := &RocketPool{
		Client:                client,
		Storage:               storage,
		MulticallAddress:      &multicallAddress,
		BalanceBatcherAddress: &balanceBatcherAddress,
	}
	rp.VersionManager = NewVersionManager(rp)

	return rp, nil

}

func SingleCall[outType core.CallReturnType](rp *RocketPool, query func(*multicall.MultiCaller, *outType), opts *bind.CallOpts) (outType, error) {
	// Run the query
	outPtr, err := multicall.MulticallQuery[outType](
		rp.Client,
		*rp.MulticallAddress,
		func(mc *multicall.MultiCaller) (*outType, error) {
			result := new(outType)
			query(mc, result)
			return result, nil
		},
		nil,
		opts,
	)
	// Handle errors
	if err != nil {
		var result outType
		return result, err
	}

	// Return the result
	return *outPtr, nil
}

func SingleCallWithError[outType core.CallReturnType](rp *RocketPool, query func(*multicall.MultiCaller, *outType, ...any) error, opts *bind.CallOpts, args ...any) (outType, error) {
	// Run the query
	outPtr, err := multicall.MulticallQuery[outType](
		rp.Client,
		*rp.MulticallAddress,
		func(mc *multicall.MultiCaller) (*outType, error) {
			result := new(outType)
			err := query(mc, result, args)
			return result, err
		},
		nil,
		opts,
	)
	// Handle errors
	if err != nil {
		var result outType
		return result, err
	}

	// Return the result
	return *outPtr, nil
}

// Load Rocket Pool contract addresses
func (rp *RocketPool) GetAddress(contractName string, opts *bind.CallOpts) (common.Address, error) {
	return SingleCall[common.Address](
		rp,
		func(mc *multicall.MultiCaller, out *common.Address) {
			rp.Storage.GetAddress(mc, out, contractName)
		},
		opts,
	)
}

// Load Rocket Pool contract ABIs
func (rp *RocketPool) GetABI(contractName string, opts *bind.CallOpts) (*abi.ABI, error) {
	// Get the encoded ABI
	abiEncoded, err := SingleCall[string](
		rp,
		func(mc *multicall.MultiCaller, out *string) {
			rp.Storage.GetAbi(mc, out, contractName)
		},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting encoded ABI: %w", err)
	}

	// Decode ABI
	abi, err := core.DecodeAbi(abiEncoded)
	if err != nil {
		return nil, fmt.Errorf("Could not decode contract %s ABI: %w", contractName, err)
	}

	// Return
	return abi, nil
}

// Load Rocket Pool contracts
func (rp *RocketPool) GetContract(contractName string, opts *bind.CallOpts) (*core.Contract, error) {

	var address common.Address
	var abiEncoded string

	err := multicall.MulticallQuery2(rp.Client, *rp.MulticallAddress, func(mc *multicall.MultiCaller) {
		rp.Storage.GetAddress(mc, &address, contractName)
		rp.Storage.GetAbi(mc, &abiEncoded, contractName)
	}, opts)
	if err != nil {
		return nil, err
	}

	// Decode ABI
	abi, err := core.DecodeAbi(abiEncoded)
	if err != nil {
		return nil, fmt.Errorf("Could not decode contract %s ABI: %w", contractName, err)
	}

	// Create contract
	contract := &core.Contract{
		Contract: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
		Address:  &address,
		ABI:      abi,
		Client:   rp.Client,
	}

	// Return
	return contract, nil

}

// Create a Rocket Pool contract instance
func (rp *RocketPool) MakeContract(contractName string, address common.Address, opts *bind.CallOpts) (*core.Contract, error) {

	// Load ABI
	abi, err := rp.GetABI(contractName, opts)
	if err != nil {
		return nil, err
	}

	// Create and return
	return &core.Contract{
		Contract: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
		Address:  &address,
		ABI:      abi,
		Client:   rp.Client,
	}, nil

}
