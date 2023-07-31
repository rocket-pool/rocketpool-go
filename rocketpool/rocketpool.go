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

	// Internal fields
	contracts    map[ContractName]*core.Contract
	instanceAbis map[ContractName]*abi.ABI // Used for instanced contracts like minipools or fee distributors
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

// Load the provided contracts by their name
func (rp *RocketPool) LoadContracts(opts *bind.CallOpts, contractNames ...ContractName) error {
	addresses := make([]common.Address, len(contractNames))
	abiStrings := make([]string, len(contractNames))

	// Load the details via multicall
	err := rp.Query(func(mc *multicall.MultiCaller) {
		for i, contractName := range contractNames {
			rp.Storage.GetAddress(mc, &addresses[i], string(contractName))
			rp.Storage.GetAbi(mc, &abiStrings[i], string(contractName))
		}
	}, opts)
	if err != nil {
		return fmt.Errorf("error getting addresses and ABIs: %w", err)
	}

	// Create the contract objects
	for i, contractName := range contractNames {
		// Decode the ABI
		abi, err := core.DecodeAbi(abiStrings[i])
		if err != nil {
			return fmt.Errorf("error decoding contract %s ABI: %w", string(contractNames[i]), err)
		}

		// Make the contract binding
		contract := &core.Contract{
			Contract: bind.NewBoundContract(addresses[i], *abi, rp.Client, rp.Client, rp.Client),
		}
		rp.contracts[contractName] = contract
	}

	// Get the versions of each contract
	emptyAddress := common.Address{}
	versions := make([]uint8, len(contractNames))
	err = rp.Query(func(mc *multicall.MultiCaller) {
		for i, contractName := range contractNames {
			address := addresses[i]
			if address != emptyAddress {
				multicall.AddCall[uint8](mc, rp.contracts[contractName], &versions[i], "version") // TODO: use the contract version getter once it's ready
			}
		}
	}, opts)
	if err != nil {
		return fmt.Errorf("error getting contract versions: %w", err)
	}

	// Assign the contract versions
	for i, contractName := range contractNames {
		rp.contracts[contractName].Version = versions[i]
	}

	return nil
}

// Load the ABIs for instances contracts (like minipools or fee distributors)
func (rp *RocketPool) LoadInstanceABIs(opts *bind.CallOpts, contractNames ...ContractName) error {
	abiStrings := make([]string, len(contractNames))

	// Load the details via multicall
	err := rp.Query(func(mc *multicall.MultiCaller) {
		for i, contractName := range contractNames {
			rp.Storage.GetAbi(mc, &abiStrings[i], string(contractName))
		}
	}, opts)
	if err != nil {
		return fmt.Errorf("error getting instanced ABIs: %w", err)
	}

	// Create the contract objects
	for i, contractName := range contractNames {
		// Decode the ABI
		abi, err := core.DecodeAbi(abiStrings[i])
		if err != nil {
			return fmt.Errorf("error decoding contract %s ABI: %w", string(contractNames[i]), err)
		}
		rp.instanceAbis[contractName] = abi
	}

	return nil
}

// Get a network contract
func (rp *RocketPool) GetContract(contractName ContractName) (*core.Contract, error) {
	contract, exists := rp.contracts[contractName]
	if !exists {
		return nil, fmt.Errorf("contract %s has not been loaded yet", string(contractName))
	}
	return contract, nil
}

// Create a binding for a network contract instance
func (rp *RocketPool) MakeContract(contractName ContractName, address common.Address) (*core.Contract, error) {
	abi, exists := rp.instanceAbis[contractName]
	if !exists {
		return nil, fmt.Errorf("ABI for contract %s has not been loaded yet", string(contractName))
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
