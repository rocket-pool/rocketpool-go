package v10_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"runtime/debug"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/v2/node"
	"github.com/rocket-pool/rocketpool-go/v2/rewards"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
	"github.com/rocket-pool/rocketpool-go/v2/tests"
	"github.com/rocket-pool/rocketpool-go/v2/tokens"
	merkletree "github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
)

type rewardsInfo struct {
	CollateralRpl    *big.Int
	OracleDaoRpl     *big.Int
	SmoothingPoolEth *big.Int
	MerkleData       []byte
	MerkleProof      []common.Hash
}

func TestRewards(t *testing.T) {
	// Revert to the initialized state at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToInitialized()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to initialized snapshot: %w", err))
		}
	})
	defer func() {
		r := recover()
		if r != nil {
			t.Logf("Recovered from panic: %s\nReverting to baseline...", r)
			err := mgr.RevertToBaseline()
			if err != nil {
				t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
			}
			debug.PrintStack()
			t.FailNow()
		}
	}()

	// Create some bindings
	txMgr := rp.GetTransactionManager()
	rplBinding, err := tokens.NewTokenRpl(rp)
	if err != nil {
		t.Fatal("error creating RPL: %w", err)
	}
	vault, err := rp.GetContract(rocketpool.ContractName_RocketVault)
	if err != nil {
		t.Fatal("error creating vault: %w", err)
	}
	smoothingPool, err := rp.GetContract(rocketpool.ContractName_RocketSmoothingPool)
	if err != nil {
		t.Fatal("error creating smoothing pool: %w", err)
	}
	rewardsPool, err := rewards.NewRewardsPool(rp)
	if err != nil {
		t.Fatal("error creating rewards pool: %w", err)
	}
	mdr, err := rewards.NewMerkleDistributorMainnet(rp)
	if err != nil {
		t.Fatal("error creating merkle distributor: %w", err)
	}

	// Query some initial settings
	var initialVaultRpl *big.Int
	err = rp.Query(func(mc *batch.MultiCaller) error {
		rplBinding.BalanceOf(mc, &initialVaultRpl, vault.Address)
		eth.AddQueryablesToMulticall(mc,
			rplBinding.InflationInterval,
			rplBinding.InflationIntervalStartTime,
			rewardsPool.RewardIndex,
		)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying initial settings: %w", err))
	}

	// Some initial info
	t.Logf("Initial Vault RPL: %s", initialVaultRpl.String())
	t.Logf("Current rewards interval: %d", rewardsPool.RewardIndex.Formatted())

	// Register the nodes
	errs := []error{
		deployNode(node1, t.Logf),
		deployNode(node2, t.Logf),
		deployNode(node3, t.Logf),
		deployNode(node4, t.Logf),
		deployNode(node5, t.Logf),
		deployNode(node6, t.Logf),
		deployNode(node7, t.Logf),
		deployNode(node8, t.Logf),
	}
	if err := errors.Join(errs...); err != nil {
		t.Fatal(err)
	}

	// Set the withdrawal addresses
	errs = []error{
		setPrimaryWithdrawalAddress(node2, node2Primary.Address, t.Logf),

		setRplWithdrawalAddress(node3, node3Rpl.Address, node3, t.Logf),

		setPrimaryWithdrawalAddress(node4, node4Primary.Address, t.Logf),
		setRplWithdrawalAddress(node4, node4Rpl.Address, node4Primary, t.Logf),

		setPrimaryWithdrawalAddress(node5, multiReceiver.Address, t.Logf),
		setRplWithdrawalAddress(node5, multiReceiver.Address, multiReceiver, t.Logf),
		setPrimaryWithdrawalAddress(node6, multiReceiver.Address, t.Logf),
		setRplWithdrawalAddress(node6, multiReceiver.Address, multiReceiver, t.Logf),
		setPrimaryWithdrawalAddress(node7, multiReceiver.Address, t.Logf),
		setRplWithdrawalAddress(node7, multiReceiver.Address, multiReceiver, t.Logf),
		setPrimaryWithdrawalAddress(node8, multiReceiver.Address, t.Logf),
		setRplWithdrawalAddress(node8, multiReceiver.Address, multiReceiver, t.Logf),
	}
	if err := errors.Join(errs...); err != nil {
		t.Fatal(err)
	}

	// Start by minting inflation so there's some RPL to distribute
	latestHeader, err := rp.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting latest header: %w", err))
	}
	currentTime := time.Unix(int64(latestHeader.Time), 0)
	timeUntilStart := rplBinding.InflationIntervalStartTime.Formatted().Sub(currentTime)
	timeToWait := timeUntilStart + rplBinding.InflationInterval.Formatted()
	timeToWaitSeconds := int(timeToWait.Seconds())
	err = mgr.IncreaseTime(timeToWaitSeconds)
	if err != nil {
		t.Fatal(fmt.Errorf("error increasing time: %w", err))
	}
	t.Logf("Increased time by %d seconds", timeToWaitSeconds)
	txInfo, err := rplBinding.MintInflationRPL(odao1.Transactor)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating mint Tx: %w", err))
	}
	tx, err := txMgr.ExecuteTransaction(txInfo, odao1.Transactor)
	if err != nil {
		t.Fatal(fmt.Errorf("error executing mint Tx: %w", err))
	}
	err = txMgr.WaitForTransaction(tx)
	if err != nil {
		t.Fatal(fmt.Errorf("error waiting for mint Tx: %w", err))
	}
	t.Log("Inflation successfully minted")

	// Make sure the vault has the new inflation
	var vaultRpl *big.Int
	err = rp.Query(func(mc *batch.MultiCaller) error {
		rplBinding.BalanceOf(mc, &vaultRpl, vault.Address)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying vault RPL: %w", err))
	}
	if vaultRpl.Cmp(initialVaultRpl) == 0 {
		t.Fatal("Vault RPL didn't increase")
	}
	t.Logf("Vault RPL increased to %s", vaultRpl.String())

	// Send some ETH to the Smoothing Pool
	smoothingPoolEth := 10.0
	smoothingPoolEthWei := eth.EthToWei(smoothingPoolEth)
	sender := odao1.Transactor
	newOpts := &bind.TransactOpts{
		From:  sender.From,
		Value: smoothingPoolEthWei,
	}
	txInfo = txMgr.CreateTransactionInfoRaw(smoothingPool.Address, nil, newOpts)
	tx, err = txMgr.ExecuteTransaction(txInfo, sender)
	if err != nil {
		t.Fatal(fmt.Errorf("error sending ETH to SP: %w", err))
	}
	err = txMgr.WaitForTransaction(tx)
	if err != nil {
		t.Fatal(fmt.Errorf("error waiting for sending ETH to SP tx: %w", err))
	}
	t.Logf("Sent %.0f ETH to the Smoothing Pool", smoothingPoolEth)

	// Get some stats of the current state
	latestHeader, err = rp.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting latest header: %w", err))
	}

	// Create a new rewards snapshot
	rewardsMap, root, err := createRewardsTree(t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating rewards tree: %w", err))
	}
	odaoRpl := big.NewInt(0)
	collateralRpl := big.NewInt(0)
	spEth := big.NewInt(0)
	t.Log("Rewards submission created!")
	for address, rewards := range rewardsMap {
		odaoRpl.Add(odaoRpl, rewards.OracleDaoRpl)
		collateralRpl.Add(collateralRpl, rewards.CollateralRpl)
		spEth.Add(spEth, rewards.SmoothingPoolEth)
		t.Logf("%s:", accountNames[address])
		t.Logf("\tCollateral RPL: %s", rewards.CollateralRpl.String())
		t.Logf("\tOracle DAO RPL: %s", rewards.OracleDaoRpl.String())
		t.Logf("\tSmoothing Pool ETH: %s", rewards.SmoothingPoolEth.String())
	}
	t.Log()
	rewardSnapshot := rewards.RewardSubmission{
		RewardIndex:     big.NewInt(0),
		ExecutionBlock:  latestHeader.Number,
		ConsensusBlock:  latestHeader.Number,
		MerkleRoot:      root,
		MerkleTreeCID:   "",
		IntervalsPassed: big.NewInt(1),
		TreasuryRPL:     big.NewInt(300),
		TrustedNodeRPL: []*big.Int{
			odaoRpl,
		},
		NodeRPL: []*big.Int{
			collateralRpl,
		},
		NodeETH: []*big.Int{
			spEth,
		},
		UserETH: big.NewInt(400),
	}

	// Submit it with 2 Oracles
	err = submitRewardSnapshot(rewardsPool, rewardSnapshot, odao1, t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error submitting rewards snapshot from ODAO 1: %w", err))
	}
	err = submitRewardSnapshot(rewardsPool, rewardSnapshot, odao2, t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error submitting rewards snapshot from ODAO 2: %w", err))
	}

	// Ensure the interval was incremented and the snapshot is canon
	oldInterval := rewardsPool.RewardIndex.Formatted()
	err = rp.Query(nil, nil, rewardsPool.RewardIndex)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting new interval: %w", err))
	}
	interval := rewardsPool.RewardIndex.Formatted()
	if oldInterval == interval {
		t.Fatal("Interval wasn't incremented")
	}
	t.Logf("Interval incremented to %d successfully", interval)

	// ========
	// Confirm the rewards are expected and claims work properly
	// ========

	balanceAddresses := []common.Address{
		node1.Address,
		node2.Address,
		node3.Address,
		node4.Address,

		node2Primary.Address,
		node3Rpl.Address,
		node4Primary.Address,
		node4Rpl.Address,

		node5.Address,
		node6.Address,
		node7.Address,
		node8.Address,

		multiReceiver.Address,
	}

	// Get the initial balances
	type balance struct {
		Eth *big.Int
		Rpl *big.Int
	}
	initBalances := make([]*balance, len(balanceAddresses))
	err = rp.Query(func(mc *batch.MultiCaller) error {
		for i, address := range balanceAddresses {
			initBalances[i] = &balance{}
			rplBinding.BalanceOf(mc, &initBalances[i].Rpl, address)
		}
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting initial RPL balances: %w", err))
	}
	ethBalances, err := rp.BalanceBatcher.GetEthBalances(balanceAddresses, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting initial ETH balances: %w", err))
	}
	for i, balance := range initBalances {
		balance.Eth = ethBalances[i]
	}

	// Claim the rewards
	gas, err := claimRewards(mdr, rewardsMap[node1.Address], node1.Address, node1, t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error claiming rewards for node 1: %w", err))
	}
	initBalances[0].Eth.Sub(initBalances[0].Eth, gas)
	gas, err = claimRewards(mdr, rewardsMap[node2.Address], node2.Address, node2, t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error claiming rewards for node 2: %w", err))
	}
	initBalances[1].Eth.Sub(initBalances[1].Eth, gas)
	gas, err = claimRewards(mdr, rewardsMap[node3.Address], node3.Address, node3Rpl, t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error claiming rewards for node 3 from RPL address: %w", err))
	}
	initBalances[5].Eth.Sub(initBalances[5].Eth, gas)
	gas, err = claimRewards(mdr, rewardsMap[node4Primary.Address], node4Primary.Address, node4Primary, t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error claiming rewards for node 4 primary: %w", err))
	}
	initBalances[6].Eth.Sub(initBalances[6].Eth, gas)
	gas, err = claimRewards(mdr, rewardsMap[node4Rpl.Address], node4Rpl.Address, node4Rpl, t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error claiming rewards for node 4 RPL: %w", err))
	}
	initBalances[7].Eth.Sub(initBalances[7].Eth, gas)
	gas, err = claimRewards(mdr, rewardsMap[multiReceiver.Address], multiReceiver.Address, multiReceiver, t.Logf)
	if err != nil {
		t.Fatal(fmt.Errorf("error claiming rewards for %s: %w", accountNames[multiReceiver.Address], err))
	}
	initBalances[12].Eth.Sub(initBalances[12].Eth, gas)

	// Get the new balances
	latestBalances := make([]*balance, len(balanceAddresses))
	err = rp.Query(func(mc *batch.MultiCaller) error {
		for i, address := range balanceAddresses {
			latestBalances[i] = &balance{}
			rplBinding.BalanceOf(mc, &latestBalances[i].Rpl, address)
		}
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting updated RPL balances: %w", err))
	}
	ethBalances, err = rp.BalanceBatcher.GetEthBalances(balanceAddresses, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting updated ETH balances: %w", err))
	}
	for i, balance := range latestBalances {
		balance.Eth = ethBalances[i]
	}
	t.Log("Got latest balances:")
	for i := 0; i < len(balanceAddresses); i++ {
		t.Logf("\t%s - ETH delta: %s, RPL delta: %s", accountNames[balanceAddresses[i]],
			big.NewInt(0).Sub(latestBalances[i].Eth, initBalances[i].Eth).String(),
			big.NewInt(0).Sub(latestBalances[i].Rpl, initBalances[i].Rpl).String(),
		)
	}
	t.Log()

	// Check the balances to ensure the proper amount of rewards were claimed
	expectedBalances := []*balance{
		{
			// Node 1
			Eth: big.NewInt(0).Add(rewardsMap[node1.Address].SmoothingPoolEth, initBalances[0].Eth),
			Rpl: big.NewInt(0).Add(rewardsMap[node1.Address].CollateralRpl, initBalances[0].Rpl),
		},
		{
			// Node 2
			Eth: big.NewInt(0).Set(initBalances[1].Eth),
			Rpl: big.NewInt(0).Set(initBalances[1].Rpl),
		},
		{
			// Node 3
			Eth: big.NewInt(0).Add(rewardsMap[node3.Address].SmoothingPoolEth, initBalances[2].Eth),
			Rpl: big.NewInt(0).Set(initBalances[2].Rpl),
		},
		{
			// Node 4
			Eth: big.NewInt(0).Set(initBalances[3].Eth),
			Rpl: big.NewInt(0).Set(initBalances[3].Rpl),
		},
		{
			// Node 2 Primary
			Eth: big.NewInt(0).Add(rewardsMap[node2.Address].SmoothingPoolEth, initBalances[4].Eth),
			Rpl: big.NewInt(0).Add(rewardsMap[node2.Address].CollateralRpl, initBalances[4].Rpl),
		},
		{
			// Node 3 RPL
			Eth: big.NewInt(0).Set(initBalances[5].Eth),
			Rpl: big.NewInt(0).Add(rewardsMap[node3.Address].CollateralRpl, initBalances[5].Rpl),
		},
		{
			// Node 4 Primary
			Eth: big.NewInt(0).Add(rewardsMap[node4Primary.Address].SmoothingPoolEth, initBalances[6].Eth),
			Rpl: big.NewInt(0).Set(initBalances[6].Rpl),
		},
		{
			// Node 4 RPL
			Eth: big.NewInt(0).Set(initBalances[7].Eth),
			Rpl: big.NewInt(0).Add(rewardsMap[node4Rpl.Address].CollateralRpl, initBalances[7].Rpl),
		},
		{
			// Node 5
			Eth: big.NewInt(0).Set(initBalances[8].Eth),
			Rpl: big.NewInt(0).Set(initBalances[8].Rpl),
		},
		{
			// Node 6
			Eth: big.NewInt(0).Set(initBalances[9].Eth),
			Rpl: big.NewInt(0).Set(initBalances[9].Rpl),
		},
		{
			// Node 7
			Eth: big.NewInt(0).Set(initBalances[10].Eth),
			Rpl: big.NewInt(0).Set(initBalances[10].Rpl),
		},
		{
			// Node 8
			Eth: big.NewInt(0).Set(initBalances[11].Eth),
			Rpl: big.NewInt(0).Set(initBalances[11].Rpl),
		},
		{
			// Multi Primary
			Eth: big.NewInt(0).Add(rewardsMap[multiReceiver.Address].SmoothingPoolEth, initBalances[12].Eth),
			Rpl: big.NewInt(0).Add(rewardsMap[multiReceiver.Address].CollateralRpl, initBalances[12].Rpl),
		},
	}
	for i, expectedBalance := range expectedBalances {
		if latestBalances[i].Eth.Cmp(expectedBalance.Eth) != 0 {
			t.Fatalf("ETH balance for %s is incorrect: expected %s but got %s", accountNames[balanceAddresses[i]], expectedBalance.Eth.String(), latestBalances[i].Eth.String())
		}
		if latestBalances[i].Rpl.Cmp(expectedBalance.Rpl) != 0 {
			t.Fatalf("RPL balance for %s is incorrect: expected %s but got %s", accountNames[balanceAddresses[i]], expectedBalance.Rpl.String(), latestBalances[i].Rpl.String())
		}
	}

	t.Log("Rewards test passed")
}

func createRewardsTree(logger func(format string, args ...any)) (map[common.Address]*rewardsInfo, common.Hash, error) {
	rewards := map[common.Address]*rewardsInfo{
		odao1.Address: {
			CollateralRpl:    big.NewInt(0),
			OracleDaoRpl:     big.NewInt(101),
			SmoothingPoolEth: big.NewInt(0),
		},
		odao2.Address: {
			CollateralRpl:    big.NewInt(0),
			OracleDaoRpl:     big.NewInt(102),
			SmoothingPoolEth: big.NewInt(0),
		},
		odao3.Address: {
			CollateralRpl:    big.NewInt(0),
			OracleDaoRpl:     big.NewInt(103),
			SmoothingPoolEth: big.NewInt(0),
		},
		node1.Address: {
			CollateralRpl:    big.NewInt(200),
			OracleDaoRpl:     big.NewInt(0),
			SmoothingPoolEth: big.NewInt(225),
		},
		node2.Address: {
			CollateralRpl:    big.NewInt(201),
			OracleDaoRpl:     big.NewInt(0),
			SmoothingPoolEth: big.NewInt(226),
		},
		node3.Address: {
			CollateralRpl:    big.NewInt(202),
			OracleDaoRpl:     big.NewInt(0),
			SmoothingPoolEth: big.NewInt(227),
		},
		node4.Address: {
			CollateralRpl:    big.NewInt(203),
			OracleDaoRpl:     big.NewInt(0),
			SmoothingPoolEth: big.NewInt(228),
		},
		node5.Address: {
			CollateralRpl:    big.NewInt(204),
			OracleDaoRpl:     big.NewInt(0),
			SmoothingPoolEth: big.NewInt(229),
		},
		node6.Address: {
			CollateralRpl:    big.NewInt(205),
			OracleDaoRpl:     big.NewInt(0),
			SmoothingPoolEth: big.NewInt(230),
		},
		node7.Address: {
			CollateralRpl:    big.NewInt(206),
			OracleDaoRpl:     big.NewInt(0),
			SmoothingPoolEth: big.NewInt(231),
		},
		node8.Address: {
			CollateralRpl:    big.NewInt(207),
			OracleDaoRpl:     big.NewInt(0),
			SmoothingPoolEth: big.NewInt(233),
		},
	}

	// v10 calculation
	v10Rewards := map[common.Address]*rewardsInfo{}

	for address, ogRewards := range rewards {
		// Get the info
		node, err := node.NewNode(rp, address)
		if err != nil {
			return nil, common.Hash{}, fmt.Errorf("error creating node binding: %w", err)
		}
		err = rp.Query(nil, nil,
			node.IsRplWithdrawalAddressSet,
			node.RplWithdrawalAddress,
			node.PrimaryWithdrawalAddress,
		)
		if err != nil {
			return nil, common.Hash{}, fmt.Errorf("error getting node info: %w", err)
		}

		if !node.IsRplWithdrawalAddressSet.Get() {
			logger("%s has no RPL withdrawal address set, using address as entry", accountNames[address])
			logger("\tCollateral RPL: %s", ogRewards.CollateralRpl)
			logger("\tOracle DAO RPL: %s", ogRewards.OracleDaoRpl)
			logger("\tETH: %s", ogRewards.SmoothingPoolEth)
			addClaimer(v10Rewards, address, ogRewards)
		} else if node.PrimaryWithdrawalAddress.Get() == address {
			logger("%s has an RPL withdrawal address set but no primary, using address as entry", accountNames[address])
			logger("\tCollateral RPL: %s", ogRewards.CollateralRpl)
			logger("\tOracle DAO RPL: %s", ogRewards.OracleDaoRpl)
			logger("\tETH: %s", ogRewards.SmoothingPoolEth)
			addClaimer(v10Rewards, address, ogRewards)
		} else {
			logger("%s has an RPL withdrawal address set and a different primary, invoking ruleset v10", accountNames[address])
			logger("\tRPL (%s) Collateral RPL: %s", accountNames[node.RplWithdrawalAddress.Get()], ogRewards.CollateralRpl)
			logger("\tRPL (%s) Oracle DAO RPL: %s", accountNames[node.RplWithdrawalAddress.Get()], ogRewards.OracleDaoRpl)
			logger("\tPrimary (%s) ETH: %s", accountNames[node.PrimaryWithdrawalAddress.Get()], ogRewards.SmoothingPoolEth)
			addClaimer(v10Rewards, node.PrimaryWithdrawalAddress.Get(), &rewardsInfo{
				// Primary gets the ETH
				CollateralRpl:    big.NewInt(0),
				OracleDaoRpl:     big.NewInt(0),
				SmoothingPoolEth: big.NewInt(0).Set(ogRewards.SmoothingPoolEth),
			})
			addClaimer(v10Rewards, node.RplWithdrawalAddress.Get(), &rewardsInfo{
				// RPL gets the RPL
				CollateralRpl:    big.NewInt(0).Set(ogRewards.CollateralRpl),
				OracleDaoRpl:     big.NewInt(0).Set(ogRewards.OracleDaoRpl),
				SmoothingPoolEth: big.NewInt(0),
			})
		}
	}

	root, err := generateMerkleTree(v10Rewards)
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("error generating Merkle Tree: %w", err)
	}

	logger("Rewards Merkle Tree generated with root %s", root.Hex())
	return v10Rewards, root, nil
}

func addClaimer(v10Rewards map[common.Address]*rewardsInfo, claimer common.Address, rewards *rewardsInfo) {
	v10entry, exists := v10Rewards[claimer]
	if !exists {
		v10entry = rewards
	} else {
		v10entry.CollateralRpl.Add(v10entry.CollateralRpl, rewards.CollateralRpl)
		v10entry.OracleDaoRpl.Add(v10entry.OracleDaoRpl, rewards.OracleDaoRpl)
		v10entry.SmoothingPoolEth.Add(v10entry.SmoothingPoolEth, rewards.SmoothingPoolEth)
	}
	v10Rewards[claimer] = v10entry
}

func deployNode(account *tests.Account, logger func(format string, args ...any)) error {
	txMgr := rp.GetTransactionManager()
	node1Binding, err := node.NewNode(rp, account.Address)
	if err != nil {
		return fmt.Errorf("error creating %s binding: %w", accountNames[account.Address], err)
	}
	txInfo, err := node1Binding.Register("Etc/UTC", account.Transactor)
	if err != nil {
		return fmt.Errorf("error creating %s registration: %w", accountNames[account.Address], err)
	}
	tx, err := txMgr.ExecuteTransaction(txInfo, account.Transactor)
	if err != nil {
		return fmt.Errorf("error executing %s registration: %w", accountNames[account.Address], err)
	}
	err = txMgr.WaitForTransaction(tx)
	if err != nil {
		return fmt.Errorf("error waiting for %s registration: %w", accountNames[account.Address], err)
	}
	logger("%s (%s) successfully registered", accountNames[account.Address], account.Address.Hex())
	return nil
}

func setPrimaryWithdrawalAddress(nodeAccount *tests.Account, address common.Address, logger func(format string, args ...any)) error {
	txMgr := rp.GetTransactionManager()
	nodeBinding, err := node.NewNode(rp, nodeAccount.Address)
	if err != nil {
		return fmt.Errorf("error creating node binding: %w", err)
	}
	txInfo, err := nodeBinding.SetPrimaryWithdrawalAddress(address, true, nodeAccount.Transactor)
	if err != nil {
		return fmt.Errorf("error creating set primary withdrawal address Tx: %w", err)
	}
	tx, err := txMgr.ExecuteTransaction(txInfo, nodeAccount.Transactor)
	if err != nil {
		return fmt.Errorf("error executing set primary withdrawal address Tx: %w", err)
	}
	err = txMgr.WaitForTransaction(tx)
	if err != nil {
		return fmt.Errorf("error waiting for set primary withdrawal address Tx: %w", err)
	}
	logger("Primary withdrawal address for %s set to %s", accountNames[nodeAccount.Address], accountNames[address])
	return nil
}

func setRplWithdrawalAddress(nodeAccount *tests.Account, address common.Address, sender *tests.Account, logger func(format string, args ...any)) error {
	txMgr := rp.GetTransactionManager()
	nodeBinding, err := node.NewNode(rp, nodeAccount.Address)
	if err != nil {
		return fmt.Errorf("error creating node binding: %w", err)
	}
	txInfo, err := nodeBinding.SetRplWithdrawalAddress(address, true, sender.Transactor)
	if err != nil {
		return fmt.Errorf("error creating set RPL withdrawal address Tx: %w", err)
	}
	tx, err := txMgr.ExecuteTransaction(txInfo, sender.Transactor)
	if err != nil {
		return fmt.Errorf("error executing set RPL withdrawal address Tx: %w", err)
	}
	err = txMgr.WaitForTransaction(tx)
	if err != nil {
		return fmt.Errorf("error waiting for set RPL withdrawal address Tx: %w", err)
	}
	logger("RPL withdrawal address for %s set to %s", accountNames[nodeAccount.Address], accountNames[address])
	return nil
}

func generateMerkleTree(rewards map[common.Address]*rewardsInfo) (common.Hash, error) {
	// Generate the leaf data for each claimer
	totalData := make([][]byte, 0, len(rewards))
	for address, rewardsForClaimer := range rewards {
		// Ignore claimers that didn't receive any rewards
		if rewardsForClaimer.CollateralRpl.Cmp(common.Big0) == 0 && rewardsForClaimer.OracleDaoRpl.Cmp(common.Big0) == 0 && rewardsForClaimer.SmoothingPoolEth.Cmp(common.Big0) == 0 {
			continue
		}

		// Claimer data is address[20] :: network[32] :: RPL[32] :: ETH[32]
		claimerData := make([]byte, 0, 20+32*3)

		// Claimer address
		addressBytes := address.Bytes()
		claimerData = append(claimerData, addressBytes...)

		// Claimer network
		network := big.NewInt(0)
		networkBytes := make([]byte, 32)
		network.FillBytes(networkBytes)
		claimerData = append(claimerData, networkBytes...)

		// RPL rewards
		rplRewards := big.NewInt(0)
		rplRewards.Add(rewardsForClaimer.CollateralRpl, rewardsForClaimer.OracleDaoRpl)
		rplRewardsBytes := make([]byte, 32)
		rplRewards.FillBytes(rplRewardsBytes)
		claimerData = append(claimerData, rplRewardsBytes...)

		// ETH rewards
		ethRewardsBytes := make([]byte, 32)
		rewardsForClaimer.SmoothingPoolEth.FillBytes(ethRewardsBytes)
		claimerData = append(claimerData, ethRewardsBytes...)

		// Assign it to the claimer rewards tracker and add it to the leaf data slice
		rewardsForClaimer.MerkleData = claimerData
		totalData = append(totalData, claimerData)
	}

	// Generate the tree
	tree, err := merkletree.NewUsing(totalData, keccak256.New(), false, true)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error generating Merkle Tree: %w", err)
	}

	// Generate the proofs for each claimer
	for address, rewardsForClaimer := range rewards {
		// Get the proof
		proof, err := tree.GenerateProof(rewardsForClaimer.MerkleData, 0)
		if err != nil {
			return common.Hash{}, fmt.Errorf("error generating proof for claimer %s: %w", address.Hex(), err)
		}

		// Convert the proof into hex strings
		proofHashes := make([]common.Hash, len(proof.Hashes))
		for i, hash := range proof.Hashes {
			proofHashes[i] = common.BytesToHash(hash)
		}

		// Assign the proof hashes to the claimer rewards struct
		rewardsForClaimer.MerkleProof = proofHashes
	}

	merkleRoot := common.BytesToHash(tree.Root())
	return merkleRoot, nil
}

func submitRewardSnapshot(rewardsPool *rewards.RewardsPool, submission rewards.RewardSubmission, sender *tests.Account, logger func(format string, args ...any)) error {
	txMgr := rp.GetTransactionManager()
	txInfo, err := rewardsPool.SubmitRewardSnapshot(submission, sender.Transactor)
	if err != nil {
		return fmt.Errorf("error creating rewards snapshot Tx: %w", err)
	}
	tx, err := txMgr.ExecuteTransaction(txInfo, sender.Transactor)
	if err != nil {
		return fmt.Errorf("error executing rewards snapshot Tx: %w", err)
	}
	err = txMgr.WaitForTransaction(tx)
	if err != nil {
		return fmt.Errorf("error waiting for rewards snapshot Tx: %w", err)
	}
	logger("Rewards snapshot submitted by %s", accountNames[sender.Address])
	return nil
}

func claimRewards(mdr *rewards.MerkleDistributorMainnet, rewards *rewardsInfo, rewardsAddress common.Address, sender *tests.Account, logger func(format string, args ...any)) (*big.Int, error) {
	// Check if the rewards address is a node, just for posterity
	node, err := node.NewNode(rp, rewardsAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating node binding: %w", err)
	}
	err = rp.Query(nil, nil,
		node.Exists,
	)
	if err != nil {
		return nil, fmt.Errorf("error checking if rewards address is a node: %w", err)
	}
	if !node.Exists.Get() {
		logger("Rewards address %s is NOT a registered node", accountNames[rewardsAddress])
	} else {
		logger("Rewards address %s is a registered node", accountNames[rewardsAddress])
	}

	txMgr := rp.GetTransactionManager()
	txInfo, err := mdr.Claim(
		rewardsAddress,
		[]*big.Int{
			big.NewInt(0),
		},
		[]*big.Int{
			big.NewInt(0).Add(rewards.CollateralRpl, rewards.OracleDaoRpl),
		},
		[]*big.Int{
			rewards.SmoothingPoolEth,
		},
		[][]common.Hash{
			rewards.MerkleProof,
		},
		sender.Transactor,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating rewards claim Tx: %w", err)
	}
	tx, err := txMgr.ExecuteTransaction(txInfo, sender.Transactor)
	if err != nil {
		return nil, fmt.Errorf("error executing rewards claim Tx: %w", err)
	}
	err = txMgr.WaitForTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for rewards claim Tx: %w", err)
	}

	// Get the block the TX was in
	receipt, err := rp.Client.TransactionReceipt(context.Background(), tx.Hash())
	gasUsed := big.NewInt(int64(receipt.GasUsed))
	cost := big.NewInt(0).Mul(receipt.EffectiveGasPrice, gasUsed)

	logger("Rewards claimed for %s by %s, cost: %s", accountNames[rewardsAddress], accountNames[sender.Address], cost)
	return cost, nil
}
