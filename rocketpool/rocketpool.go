package rocketpool

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/go-version"
	"github.com/rocket-pool/node-manager-core/eth"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/storage"
)

const (
	defaultAddressBatchSize         int = 1000
	defaultContractVersionBatchSize int = 500
	defaultBalanceBatchSize         int = 1000
)

// Rocket Pool contract manager
type RocketPool struct {
	Client                   eth.IExecutionClient
	Storage                  *storage.Storage
	MulticallAddress         *common.Address
	BalanceBatcher           *batch.BalanceBatcher
	VersionManager           *VersionManager
	ConcurrentCallLimit      int
	AddressBatchSize         int
	ContractVersionBatchSize int

	// Internal fields
	contracts    map[ContractName]*core.Contract
	instanceAbis map[ContractName]*abi.ABI // Used for instanced contracts like minipools or fee distributors
	txMgr        *eth.TransactionManager
	queryMgr     *eth.QueryManager
}

// Create new contract manager
func NewRocketPool(client eth.IExecutionClient, rocketStorageAddress common.Address, multicallAddress common.Address, balanceBatcherAddress common.Address) (*RocketPool, error) {
	// Create the RocketStorage binding
	storage, err := storage.NewStorage(client, rocketStorageAddress)
	if err != nil {
		return nil, fmt.Errorf("error initializing Rocket Pool storage contract: %w", err)
	}

	// Create the balance batcher
	concurrentCallLimit := runtime.NumCPU() / 2
	balanceBatcher, err := batch.NewBalanceBatcher(client, balanceBatcherAddress, defaultBalanceBatchSize, concurrentCallLimit)
	if err != nil {
		return nil, fmt.Errorf("error creating balance batcher: %w", err)
	}

	// Create the binding
	rp := &RocketPool{
		Client:                   client,
		Storage:                  storage,
		MulticallAddress:         &multicallAddress,
		BalanceBatcher:           balanceBatcher,
		ConcurrentCallLimit:      concurrentCallLimit,
		AddressBatchSize:         defaultAddressBatchSize,
		ContractVersionBatchSize: defaultContractVersionBatchSize,
		contracts:                map[ContractName]*core.Contract{},
		instanceAbis:             map[ContractName]*abi.ABI{},
	}
	rp.VersionManager = NewVersionManager(rp)
	rp.txMgr, err = eth.NewTransactionManager(client, eth.DefaultSafeGasBuffer, eth.DefaultSafeGasMultiplier)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction manager: %w", err)
	}
	rp.queryMgr = eth.NewQueryManager(client, multicallAddress, rp.ConcurrentCallLimit)
	return rp, nil
}

// Load all network contracts
func (rp *RocketPool) LoadAllContracts(opts *bind.CallOpts) error {
	err := rp.LoadContracts(opts, ContractNames...)
	if err != nil {
		return fmt.Errorf("error loading contracts: %w", err)
	}

	err = rp.LoadInstanceABIs(opts, InstanceContractNames...)
	if err != nil {
		return fmt.Errorf("error loading instance contract ABIs: %w", err)
	}
	return nil
}

// Load only the provided specific contracts by their name
func (rp *RocketPool) LoadContracts(opts *bind.CallOpts, contractNames ...ContractName) error {
	addresses := make([]common.Address, len(contractNames))
	abiStrings := make([]string, len(contractNames))

	// Load the details via multicall
	results, err := rp.FlexQuery(func(mc *batch.MultiCaller) error {
		for i, contractName := range contractNames {
			rp.Storage.GetAddress(mc, &addresses[i], string(contractName))
			rp.Storage.GetAbi(mc, &abiStrings[i], string(contractName))
		}
		return nil
	}, opts)
	if err != nil {
		return fmt.Errorf("error getting addresses and ABIs: %w", err)
	}
	for i, result := range results {
		if !result {
			contractName := contractNames[i]
			return fmt.Errorf("failed getting address and ABI for contract %s", contractName)
		}
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
			Contract: &eth.Contract{
				Name:         string(contractName),
				ContractImpl: bind.NewBoundContract(addresses[i], *abi, rp.Client, rp.Client, rp.Client),
				Address:      addresses[i],
				ABI:          abi,
			},
		}
		rp.contracts[contractName] = contract
	}

	// Get the versions of each contract
	results, err = rp.FlexQuery(func(mc *batch.MultiCaller) error {
		for _, contractName := range contractNames {
			contract := rp.contracts[contractName]
			err := GetContractVersion(mc, &contract.Version, contract.Address)
			if err != nil {
				return fmt.Errorf("error getting version for contract %s: %w", string(contractName), err)
			}
		}
		return nil
	}, opts)
	if err != nil {
		return fmt.Errorf("error getting contract versions: %w", err)
	}
	for i, result := range results {
		if !result {
			contract := rp.contracts[contractNames[i]]
			contract.Version = 1 // If the contract doesn't have a version() in its ABI then it's v1
		}
	}

	return nil
}

// Load the ABIs for instances contracts (like minipools or fee distributors)
func (rp *RocketPool) LoadInstanceABIs(opts *bind.CallOpts, contractNames ...ContractName) error {
	abiStrings := make([]string, len(contractNames))

	// Load the details via multicall
	err := rp.Query(func(mc *batch.MultiCaller) error {
		for i, contractName := range contractNames {
			rp.Storage.GetAbi(mc, &abiStrings[i], string(contractName))
		}
		return nil
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

// Get several network contracts
func (rp *RocketPool) GetContracts(contractNames ...ContractName) ([]*core.Contract, error) {
	contracts := make([]*core.Contract, len(contractNames))
	for i, contractName := range contractNames {
		contract, exists := rp.contracts[contractName]
		if !exists {
			return nil, fmt.Errorf("contract %s has not been loaded yet", string(contractName))
		}
		contracts[i] = contract
	}
	return contracts, nil
}

// Create a binding for a network contract instance
func (rp *RocketPool) MakeContract(contractName ContractName, address common.Address) (*core.Contract, error) {
	abi, err := rp.GetAbi(contractName)
	if err != nil {
		return nil, err
	}

	// Create and return
	return &core.Contract{
		Contract: &eth.Contract{
			Name:         string(contractName),
			ContractImpl: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
			Address:      address,
			ABI:          abi,
		},
	}, nil
}

// Get the ABI for a network contract (typically used for instances like minipools or fee distributors)
func (rp *RocketPool) GetAbi(contractName ContractName) (*abi.ABI, error) {
	abi, exists := rp.instanceAbis[contractName]
	if !exists {
		return nil, fmt.Errorf("ABI for contract %s has not been loaded yet", string(contractName))
	}
	return abi, nil
}

func (rp *RocketPool) GetTransactionManager() *eth.TransactionManager {
	return rp.txMgr
}

func (rp *RocketPool) GetQueryManager() *eth.QueryManager {
	return rp.queryMgr
}

// =============
// === Utils ===
// =============

func (rp *RocketPool) GetProtocolVersion(opts *bind.CallOpts) (*version.Version, error) {
	// Try getting the version from storage directly if present
	var protocolVersion string
	err := rp.Query(func(mc *batch.MultiCaller) error {
		rp.Storage.GetProtocolVersion(mc, &protocolVersion)
		return nil
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving protocol version: %w", err)
	}

	// Convert it to a semver
	if protocolVersion != "" {
		semver, err := version.NewSemver(protocolVersion)
		if err != nil {
			return nil, fmt.Errorf("protocol version stored in contracts [%s] could not be parsed: %w", protocolVersion, err)
		}
		return semver, nil
	}

	// Fall back to the legacy checking behavior
	nodeStaking, err := rp.GetContract(ContractName_RocketNodeStaking)
	if err != nil {
		return nil, fmt.Errorf("error getting node staking contract: %w", err)
	}
	nodeMgr, err := rp.GetContract(ContractName_RocketNodeManager)
	if err != nil {
		return nil, fmt.Errorf("error getting node manager contract: %w", err)
	}

	nodeStakingVersion := nodeStaking.Version
	nodeMgrVersion := nodeMgr.Version
	// Check for v1.2 (Atlas)
	if nodeStakingVersion > 3 {
		return version.NewSemver("1.2.0")
	}

	// Check for v1.1 (Redstone)
	if nodeMgrVersion > 1 {
		return version.NewSemver("1.1.0")
	}

	// v1.0 (Classic)
	return version.NewSemver("1.0.0")
}

// Create a contract directly from its ABI, encoded in string form
func (rp *RocketPool) CreateMinipoolContractFromEncodedAbi(address common.Address, encodedAbi string) (*core.Contract, error) {
	// Decode ABI
	abi, err := core.DecodeAbi(encodedAbi)
	if err != nil {
		return nil, fmt.Errorf("error decoding minipool %s ABI: %w", address, err)
	}

	// Create and return
	return &core.Contract{
		Contract: &eth.Contract{
			ContractImpl: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
			Address:      address,
			ABI:          abi,
		},
	}, nil
}

// Create a contract directly from its ABI
func (rp *RocketPool) CreateMinipoolContractFromAbi(address common.Address, abi *abi.ABI) (*core.Contract, error) {
	// Create and return
	return &core.Contract{
		Contract: &eth.Contract{
			ContractImpl: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
			Address:      address,
			ABI:          abi,
		},
	}, nil
}

// =========================
// === Multicall Helpers ===
// =========================

// Run a multicall query that doesn't perform any return type allocation.
// The 'query' function is an optional general-purpose function you can use to add whatever you want to the multicall
// before running it. The 'queryables' can be used to simply list a collection of IQueryable objects, each of which will
// run 'AddToQuery()' on the multicall for convenience.
func (rp *RocketPool) Query(query func(*batch.MultiCaller) error, opts *bind.CallOpts, queryables ...eth.IQueryable) error {
	return rp.queryMgr.Query(query, opts, queryables...)
}

// Run a multicall query that doesn't perform any return type allocation
// Use this if one of the calls is allowed to fail without interrupting the others; the returned result array provides information about the success of each call.
// The 'query' function is an optional general-purpose function you can use to add whatever you want to the multicall
// before running it. The 'queryables' can be used to simply list a collection of IQueryable objects, each of which will
// run 'AddToQuery()' on the multicall for convenience.
func (rp *RocketPool) FlexQuery(query func(*batch.MultiCaller) error, opts *bind.CallOpts, queryables ...eth.IQueryable) ([]bool, error) {
	return rp.queryMgr.FlexQuery(query, opts, queryables...)
}

// Create and execute a multicall query that is too big for one call and must be run in batches
func (rp *RocketPool) BatchQuery(count int, batchSize int, query func(*batch.MultiCaller, int) error, opts *bind.CallOpts) error {
	return rp.queryMgr.BatchQuery(count, batchSize, query, opts)
}

// Create and execute a multicall query that is too big for one call and must be run in batches.
// Use this if one of the calls is allowed to fail without interrupting the others; the returned result array provides information about the success of each call.
func (rp *RocketPool) FlexBatchQuery(count int, batchSize int, query func(*batch.MultiCaller, int) error, handleResult func(bool, int) error, opts *bind.CallOpts) error {
	return rp.queryMgr.FlexBatchQuery(count, batchSize, query, handleResult, opts)
}

// ===========================
// === Transaction Helpers ===
// ===========================

// Signs a transaction but does not submit it to the network. Use this if you want to sign something offline and submit it later,
// or submit it as part of a bundle.
func (rp *RocketPool) SignTransaction(txInfo *eth.TransactionInfo, opts *bind.TransactOpts) (*types.Transaction, error) {
	return rp.txMgr.SignTransaction(txInfo, opts)
}

// Signs and submits a transaction to the network.
// The nonce and gas fee info in the provided opts will be used.
// The value will come from the provided txInfo. It will *not* use the value in the provided opts.
func (rp *RocketPool) ExecuteTransaction(txInfo *eth.TransactionInfo, opts *bind.TransactOpts) (*types.Transaction, error) {
	return rp.txMgr.ExecuteTransaction(txInfo, opts)
}

// Creates, signs, and submits a transaction to the network using the nonce and value from the original TX info.
// Use this if you don't care about the estimated gas cost and just want to run it as quickly as possible.
// If failOnSimErrors is true, it will treat a simualtion / gas estimation error as a failure and stop before the transaction is submitted to the network.
func (rp *RocketPool) CreateAndExecuteTransaction(creator func() (*eth.TransactionInfo, error), failOnSimError bool, opts *bind.TransactOpts) (*types.Transaction, error) {

	txInfo, err := creator()
	if err != nil {
		return nil, fmt.Errorf("error creating TX info: %w", err)
	}
	if failOnSimError && txInfo.SimulationResult.SimulationError != "" {
		return nil, fmt.Errorf("error simulating TX: %s", txInfo.SimulationResult.SimulationError)
	}

	return rp.ExecuteTransaction(txInfo, opts)
}

// Creates, signs, submits, and waits for the transaction to be included in a block.
// The nonce and gas fee info in the provided opts will be used.
// The value will come from the provided txInfo. It will *not* use the value in the provided opts.
// Use this if you don't care about the estimated gas cost and just want to run it as quickly as possible.
// If failOnSimErrors is true, it will treat a simualtion / gas estimation error as a failure and stop before the transaction is submitted to the network.
func (rp *RocketPool) CreateAndWaitForTransaction(creator func() (*eth.TransactionInfo, error), failOnSimError bool, opts *bind.TransactOpts) error {
	// Create the TX
	txInfo, err := creator()
	if err != nil {
		return fmt.Errorf("error creating TX info: %w", err)
	}
	if failOnSimError && txInfo.SimulationResult.SimulationError != "" {
		return fmt.Errorf("error simulating TX: %s", txInfo.SimulationResult.SimulationError)
	}

	// Execute the TX
	tx, err := rp.ExecuteTransaction(txInfo, opts)
	if err != nil {
		return fmt.Errorf("error executing TX: %w", err)
	}

	// Wait for the TX
	err = rp.WaitForTransaction(tx)
	if err != nil {
		return fmt.Errorf("error waiting for TX: %w", err)
	}

	return nil
}

// Signs and submits a bundle of transactions to the network that are all sent from the same address.
// The values for each TX will be in each TX info; the value specified in the opts argument is not used.
// The GasFeeCap and GasTipCap from opts will be used for all transactions.
// NOTE: this assumes the bundle is meant to be submitted sequentially, so the nonce of each one will be incremented.
// Assign the Nonce in the opts tto the nonce you want to use for the first transaction.
func (rp *RocketPool) BatchExecuteTransactions(txSubmissions []*eth.TransactionSubmission, opts *bind.TransactOpts) ([]*types.Transaction, error) {
	if opts.Nonce == nil {
		// Get the latest nonce and use that as the nonce for the first TX
		nonce, err := rp.Client.NonceAt(context.Background(), opts.From, nil)
		if err != nil {
			return nil, fmt.Errorf("error getting latest nonce for node: %w", err)
		}
		opts.Nonce = big.NewInt(0).SetUint64(nonce)
	}

	txs := make([]*types.Transaction, len(txSubmissions))
	for i, txSubmission := range txSubmissions {
		txInfo := txSubmission.TxInfo
		opts.GasLimit = txSubmission.GasLimit
		tx, err := rp.txMgr.ExecuteTransaction(txInfo, opts)
		if err != nil {
			return nil, fmt.Errorf("error creating transaction %d in bundle: %w", i, err)
		}
		txs[i] = tx

		// Increment the nonce for the next TX
		opts.Nonce.Add(opts.Nonce, common.Big1)
	}
	return txs, nil
}

// Creates, signs, and submits a collection of transactions to the network that are all sent from the same address.
// The values for each TX will be in each TX info; the value specified in the opts argument is not used.
// The GasFeeCap and GasTipCap from opts will be used for all transactions.
// Use this if you don't care about the estimated gas costs and just want to run them as quickly as possible.
// If failOnSimErrors is true, it will treat simualtion / gas estimation errors as failures and stop before any of transactions are submitted to the network.
// NOTE: this assumes the bundle is meant to be submitted sequentially, so the nonce of each one will be incremented.
// Assign the Nonce in the opts tto the nonce you want to use for the first transaction.
func (rp *RocketPool) BatchCreateAndExecuteTransactions(creators []func() (*eth.TransactionSubmission, error), failOnSimErrors bool, opts *bind.TransactOpts) ([]*types.Transaction, error) {
	// Create the TXs
	txSubmissions := make([]*eth.TransactionSubmission, len(creators))
	for i, creator := range creators {
		txSubmission, err := creator()
		if err != nil {
			return nil, fmt.Errorf("error creating TX submission for TX %d: %w", i, err)
		}
		if failOnSimErrors && txSubmission.TxInfo.SimulationResult.SimulationError != "" {
			return nil, fmt.Errorf("error simulating TX %d: %s", i, txSubmission.TxInfo.SimulationResult.SimulationError)
		}
		txSubmissions[i] = txSubmission
	}

	// Run the TXs
	return rp.BatchExecuteTransactions(txSubmissions, opts)
}

// Creates, signs, and submits a collection of transactions to the network that are all sent from the same address, then waits for them all to complete.
// The values for each TX will be in each TX info; the value specified in the opts argument is not used.
// The GasFeeCap and GasTipCap from opts will be used for all transactions.
// Use this if you don't care about the estimated gas costs and just want to run them as quickly as possible.
// If failOnSimErrors is true, it will treat simualtion / gas estimation errors as failures and stop before any of transactions are submitted to the network.
// NOTE: this assumes the bundle is meant to be submitted sequentially, so the nonce of each one will be incremented.
// Assign the Nonce in the opts tto the nonce you want to use for the first transaction.
func (rp *RocketPool) BatchCreateAndWaitForTransactions(creators []func() (*eth.TransactionSubmission, error), failOnSimErrors bool, opts *bind.TransactOpts) error {
	// Create the TXs
	txSubmissions := make([]*eth.TransactionSubmission, len(creators))
	for i, creator := range creators {
		txSubmission, err := creator()
		if err != nil {
			return fmt.Errorf("error creating TX submission for TX %d: %w", i, err)
		}
		if failOnSimErrors && txSubmission.TxInfo.SimulationResult.SimulationError != "" {
			return fmt.Errorf("error simulating TX %d: %s", i, txSubmission.TxInfo.SimulationResult.SimulationError)
		}
		txSubmissions[i] = txSubmission
	}

	// Run the TXs
	txs, err := rp.BatchExecuteTransactions(txSubmissions, opts)
	if err != nil {
		return fmt.Errorf("error running TXs: %w", err)
	}

	// Wait for the TXs
	err = rp.WaitForTransactions(txs)
	if err != nil {
		return fmt.Errorf("error waiting for TXs: %w", err)
	}

	return nil
}

// Wait for a transaction to get included in blocks
func (rp *RocketPool) WaitForTransaction(tx *types.Transaction) error {
	return rp.txMgr.WaitForTransaction(tx)
}

// Wait for a set of transactions to get included in blocks
func (rp *RocketPool) WaitForTransactions(txs []*types.Transaction) error {
	return rp.txMgr.WaitForTransactions(txs)
}

// Wait for a transaction to get included in blocks
func (rp *RocketPool) WaitForTransactionByHash(hash common.Hash) error {
	return rp.txMgr.WaitForTransactionByHash(hash)
}

// Wait for a set of transactions to get included in blocks
func (rp *RocketPool) WaitForTransactionsByHash(hashes []common.Hash) error {
	return rp.txMgr.WaitForTransactionsByHash(hashes)
}

// Get a TX from its hash
func (rp *RocketPool) getTransactionFromHash(hash common.Hash) (*types.Transaction, error) {
	// Retry for 30 sec if the TX wasn't found
	for i := 0; i < 30; i++ {
		tx, _, err := rp.Client.TransactionByHash(context.Background(), hash)
		if err != nil {
			if err.Error() == "not found" {
				time.Sleep(1 * time.Second)
				continue
			}
			return nil, err
		}

		return tx, nil
	}

	return nil, fmt.Errorf("transaction not found after 30 seconds")
}
