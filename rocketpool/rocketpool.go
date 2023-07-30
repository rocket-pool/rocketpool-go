package rocketpool

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

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

// Load Rocket Pool contract addresses
func (rp *RocketPool) GetAddress(contractName string, opts *bind.CallOpts) (common.Address, error) {
	var address common.Address
	err := rp.Query(func(mc *multicall.MultiCaller) {
		rp.Storage.GetAddress(mc, &address, contractName)
	}, opts)
	return address, err
}

// Load Rocket Pool contract ABIs
func (rp *RocketPool) GetABI(contractName string, opts *bind.CallOpts) (*abi.ABI, error) {
	// Get the encoded ABI
	var abiEncoded string
	err := rp.Query(func(mc *multicall.MultiCaller) {
		rp.Storage.GetAbi(mc, &abiEncoded, contractName)
	}, opts)
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

	err := rp.Query(func(mc *multicall.MultiCaller) {
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

// =========================
// === Multicall Helpers ===
// =========================

// Run a multicall query that doesn't perform any return type allocation
func (rp *RocketPool) Query(query func(*multicall.MultiCaller), opts *bind.CallOpts) error {
	// Create the multicaller
	mc, err := multicall.NewMultiCaller(rp.Client, *rp.MulticallAddress)
	if err != nil {
		return fmt.Errorf("error creating multicaller: %w", err)
	}

	// Run the query
	query(mc)

	// Execute the multicall
	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return fmt.Errorf("error executing multicall: %w", err)
	}

	return nil
}

// Run a multicall query that doesn't perform any return type allocation - used when the query itself can return an error
func (rp *RocketPool) QueryWithError(query func(*multicall.MultiCaller) error, opts *bind.CallOpts) error {
	// Create the multicaller
	mc, err := multicall.NewMultiCaller(rp.Client, *rp.MulticallAddress)
	if err != nil {
		return fmt.Errorf("error creating multicaller: %w", err)
	}

	// Run the query
	err = query(mc)
	if err != nil {
		return fmt.Errorf("error running multicall query: %w", err)
	}

	// Execute the multicall
	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return fmt.Errorf("error executing multicall: %w", err)
	}

	return nil
}

// Create and execute a multicall query that is too big for one call and must be run in batches
func BatchQuery[ObjType any](rp *RocketPool, count uint64, batchSize uint64, createAndQuery func(*multicall.MultiCaller, uint64) (*ObjType, error), opts *bind.CallOpts) ([]*ObjType, error) {
	// Create the array of query objects
	objs := make([]*ObjType, count)

	// Sync
	var wg errgroup.Group
	wg.SetLimit(int(batchSize))

	// Run getters in batches
	for i := uint64(0); i < count; i += batchSize {
		i := i
		max := i + batchSize
		if max > count {
			max = count
		}

		// Load details
		wg.Go(func() error {
			mc, err := multicall.NewMultiCaller(rp.Client, *rp.MulticallAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				obj, err := createAndQuery(mc, j)
				if err != nil {
					return fmt.Errorf("error running query adder: %w", err)
				}
				objs[j] = obj
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	// Wait for them all to complete
	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error during multicall query: %w", err)
	}

	// Return
	return objs, nil
}
