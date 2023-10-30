package tests

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// TestManager wraps the EVM client binding and everything needed to interact with it for the Rocket Pool unit tests
type TestManager struct {
	StorageAddress   common.Address
	OwnerAccount     *Account
	NonOwnerAccounts []*Account
	RocketPool       *rocketpool.RocketPool
	Client           core.ExecutionClient

	baselineSnapshotID    string
	initializedSnapshotID string
	rpcClient             *rpc.Client
}

var singleton *TestManager

// Creates a new TestManager or returns the instance that's already been created
func NewTestManager() (*TestManager, error) {
	// Return the instance if it's already been made
	if singleton != nil {
		return singleton, nil
	}

	// Create the rpcClient bindings
	rpcClient, err := rpc.Dial(evmRpcAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating RPC client binding: %w", err)
	}
	client, err := ethclient.Dial(evmRpcAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating ETH client binding: %w", err)
	}

	// Create the account bindings
	chainID := big.NewInt(int64(evmChainID))
	owner, err := CreateAccountFromPrivateKey(AccountPrivateKeys[0], chainID)
	if err != nil {
		return nil, fmt.Errorf("error creating owner account binding: %w", err)
	}
	accounts := make([]*Account, len(AccountPrivateKeys)-1)
	for i := 0; i < len(accounts); i++ {
		account, err := CreateAccountFromPrivateKey(AccountPrivateKeys[i+1], chainID)
		if err != nil {
			return nil, fmt.Errorf("error creating account binding %d: %w", i+1, err)
		}
		accounts[i] = account
	}

	// Contract addresses
	storageAddress := common.HexToAddress(rocketStorageAddress)
	multicallAddress := common.HexToAddress(multicallAddress)
	balanceBatcherAddress := common.HexToAddress(balanceBatcherAddress)

	// Create the RP binding
	rp, err := rocketpool.NewRocketPool(client, storageAddress, multicallAddress, balanceBatcherAddress)
	if err != nil {
		log.Fatal(err)
	}
	err = rp.LoadAllContracts(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating RP binding: %w", err)
	}

	// Create the test manager
	m := &TestManager{
		StorageAddress:   storageAddress,
		OwnerAccount:     owner,
		NonOwnerAccounts: accounts,
		RocketPool:       rp,
		Client:           client,

		rpcClient: rpcClient,
	}

	// Create the baseline snapshot
	baselineSnapshotID, err := m.takeSnapshot()
	if err != nil {
		return nil, fmt.Errorf("error creating baseline snapshot: %w", err)
	}
	m.baselineSnapshotID = baselineSnapshotID

	return m, nil
}

// Bootstraps the Oracle DAO and common network settings
func (m *TestManager) InitializeDeployment() error {
	rp := m.RocketPool

	// Get the contract bindings
	odaoMgr, err := oracle.NewOracleDaoManager(rp)
	if err != nil {
		return fmt.Errorf("error creating oDAO manager binding: %w", err)
	}

	// Get some contract state
	err = rp.Query(nil, nil, odaoMgr.MemberCount)
	if err != nil {
		return fmt.Errorf("error getting network info: %w", err)
	}

	// If there aren't members, bootstrap the protocol
	if odaoMgr.MemberCount.Formatted() == 0 {
		err = m.initializeImpl(odaoMgr)
		if err != nil {
			return fmt.Errorf("error initializing network: %w", err)
		}
	}

	// Generate the initialized snapshot
	initializedSnapshotID, err := m.takeSnapshot()
	if err != nil {
		return fmt.Errorf("error creating baseline snapshot: %w", err)
	}
	m.initializedSnapshotID = initializedSnapshotID
	return nil
}

// Reverts the EVM to the baseline snapshot
func (m *TestManager) RevertToBaseline() error {
	err := m.revertToSnapshot(m.baselineSnapshotID)
	if err != nil {
		return fmt.Errorf("error reverting to baseline snapshot: %w", err)
	}

	// Regenerate the baseline snapshot since Hardhat can't revert to it multiple times
	baselineSnapshotID, err := m.takeSnapshot()
	if err != nil {
		return fmt.Errorf("error creating baseline snapshot: %w", err)
	}
	m.baselineSnapshotID = baselineSnapshotID
	return nil
}

// Reverts the EVM to the initialized snapshot
func (m *TestManager) RevertToInitialized() error {
	err := m.revertToSnapshot(m.initializedSnapshotID)
	if err != nil {
		return fmt.Errorf("error reverting to initialized snapshot: %w", err)
	}

	// Regenerate the initialized snapshot since Hardhat can't revert to it multiple times
	initializedSnapshotID, err := m.takeSnapshot()
	if err != nil {
		return fmt.Errorf("error creating initialized snapshot: %w", err)
	}
	m.initializedSnapshotID = initializedSnapshotID
	return nil
}

// Creates a snapshot of the EVM's current state, returning the snapshot ID - this can be used in RevertToCustomSnapshot()
func (m *TestManager) CreateCustomSnapshot() (string, error) {
	return m.takeSnapshot()
}

// Reverts the EVM's current state to a previously taken snapshot
func (m *TestManager) RevertToCustomSnapshot(snapshotID string) error {
	return m.revertToSnapshot(snapshotID)
}

// Mine a number of blocks
func (m *TestManager) MineBlocks(numBlocks int) error {
	for bi := 0; bi < numBlocks; bi++ {
		err := m.rpcClient.Call(nil, "evm_mine")
		if err != nil {
			return fmt.Errorf("error mining blocks: %w", err)
		}
	}
	return nil
}

// Fast forward to some number of seconds
func (m *TestManager) IncreaseTime(time int) error {
	// Increase the time
	err := m.rpcClient.Call(nil, "evm_increaseTime", time)
	if err != nil {
		return fmt.Errorf("error increasing time: %w", err)
	}

	// Mine a new block so the time increase is captured on-chain
	err = m.MineBlocks(1)
	if err != nil {
		return fmt.Errorf("error mining a block after time increase: %w", err)
	}

	return nil
}

// Take a snapshot of the EVM's state
func (m *TestManager) takeSnapshot() (string, error) {
	var response string
	err := m.rpcClient.Call(&response, "evm_snapshot")
	if err != nil {
		return "", fmt.Errorf("error creating snapshot: %w", err)
	}
	return response, nil
}

// Revert the EVM to a snapshot state
func (m *TestManager) revertToSnapshot(snapshotID string) error {
	err := m.rpcClient.Call(nil, "evm_revert", snapshotID)
	if err != nil {
		return fmt.Errorf("error reverting to snapshot %s: %w", snapshotID, err)
	}
	return nil
}

// Implementation for initialization
func (m *TestManager) initializeImpl(odaoMgr *oracle.OracleDaoManager) error {
	rp := m.RocketPool

	// Prep the accounts
	owner := m.OwnerAccount
	odao1 := m.NonOwnerAccounts[0]
	odao2 := m.NonOwnerAccounts[1]
	odao3 := m.NonOwnerAccounts[2]

	// Get the settings managers
	pdaoMgr, err := protocol.NewProtocolDaoManager(rp)
	if err != nil {
		return fmt.Errorf("error creating protocol DAO manager binding: %w", err)
	}

	// Bootstrap all of the relevant parameters
	err = rp.BatchCreateAndWaitForTransactions([]func() (*core.TransactionSubmission, error){
		func() (*core.TransactionSubmission, error) {
			return core.CreateTxSubmissionFromInfo(pdaoMgr.Settings.Node.IsRegistrationEnabled.Bootstrap(true, owner.Transactor))
		},
	}, true, owner.Transactor)
	if err != nil {
		return fmt.Errorf("error bootstrapping settings: %w", err)
	}

	// Bootstrap the nodes onto the oDAO
	_, err = BootstrapNodeToOdao(m.RocketPool, m.OwnerAccount, odao1, "Etc/UTC", "Node 1", "node1.com")
	if err != nil {
		return fmt.Errorf("error bootstrapping node1: %w", err)
	}
	_, err = BootstrapNodeToOdao(m.RocketPool, m.OwnerAccount, odao2, "Australia/Brisbane", "Node 2", "node2.com")
	if err != nil {
		return fmt.Errorf("error bootstrapping node2: %w", err)
	}
	_, err = BootstrapNodeToOdao(m.RocketPool, m.OwnerAccount, odao3, "America/New_York", "Node 3", "node3.com")
	if err != nil {
		return fmt.Errorf("error bootstrapping node3: %w", err)
	}

	// Make sure they're on it
	err = rp.Query(nil, nil, odaoMgr.MemberCount)
	if err != nil {
		return fmt.Errorf("error getting oDAO member count: %w", err)
	}
	memberCount := odaoMgr.MemberCount.Formatted()
	if memberCount != 3 {
		return fmt.Errorf("expected 3 oDAO members but there are %d", memberCount)
	}
	addresses, err := odaoMgr.GetMemberAddresses(memberCount, nil)
	if err != nil {
		return fmt.Errorf("error getting oDAO member addresses: %w", err)
	}
	if addresses[0] != odao1.Address {
		return fmt.Errorf("address 1 is not on the oDAO (expected %s but got %s)", odao1.Address.Hex(), addresses[0].Hex())
	}
	if addresses[1] != odao2.Address {
		return fmt.Errorf("address 2 is not on the oDAO (expected %s but got %s)", odao2.Address.Hex(), addresses[1].Hex())
	}
	if addresses[2] != odao3.Address {
		return fmt.Errorf("address 3 is not on the oDAO (expected %s but got %s)", odao3.Address.Hex(), addresses[2].Hex())
	}

	return nil
}
