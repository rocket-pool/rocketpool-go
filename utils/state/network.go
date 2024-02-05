package state

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/eth"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"

	"golang.org/x/sync/errgroup"
)

const (
	networkEffectiveStakeBatchSize int = 250
)

type NetworkDetails struct {
	// Redstone
	RplPrice                      *big.Int
	MinCollateralFraction         *big.Int
	MaxCollateralFraction         *big.Int
	IntervalDuration              time.Duration
	IntervalStart                 time.Time
	NodeOperatorRewardsPercent    *big.Int
	OracleDaoRewardsPercent       *big.Int
	ProtocolDaoRewardsPercent     *big.Int
	PendingRPLRewards             *big.Int
	RewardIndex                   uint64
	ScrubPeriod                   time.Duration
	SmoothingPoolAddress          common.Address
	DepositPoolBalance            *big.Int
	DepositPoolExcess             *big.Int
	TotalQueueCapacity            *big.Int
	EffectiveQueueCapacity        *big.Int
	QueueLength                   *big.Int
	RPLInflationIntervalRate      *big.Int
	RPLTotalSupply                *big.Int
	PricesBlock                   uint64
	LatestReportablePricesBlock   uint64
	ETHUtilizationRate            float64
	StakingETHBalance             *big.Int
	RETHExchangeRate              float64
	TotalETHBalance               *big.Int
	RETHBalance                   *big.Int
	TotalRETHSupply               *big.Int
	TotalRPLStake                 *big.Int
	SmoothingPoolBalance          *big.Int
	NodeFee                       float64
	BalancesBlock                 *big.Int
	LatestReportableBalancesBlock uint64
	SubmitBalancesEnabled         bool
	SubmitPricesEnabled           bool
	MinipoolLaunchTimeout         *big.Int

	// Atlas
	PromotionScrubPeriod      time.Duration
	BondReductionWindowStart  time.Duration
	BondReductionWindowLength time.Duration
	DepositPoolUserBalance    *big.Int

	// Houston
	PricesSubmissionFrequency   uint64
	BalancesSubmissionFrequency uint64
}

// Create a snapshot of all of the network's details
func NewNetworkDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts, isHoustonDeployed bool) (*NetworkDetails, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	details := &NetworkDetails{}

	// Local vars for things that need to be converted
	var rewardIndex *big.Int
	var intervalStart *big.Int
	var intervalDuration *big.Int
	var scrubPeriodSeconds *big.Int
	var totalQueueCapacity *big.Int
	var effectiveQueueCapacity *big.Int
	var totalQueueLength *big.Int
	var pricesBlock *big.Int
	var latestReportablePricesBlock *big.Int
	var pricesSubmissionFrequency *big.Int
	var ethUtilizationRate *big.Int
	var rETHExchangeRate *big.Int
	var nodeFee *big.Int
	var balancesBlock *big.Int
	var latestReportableBalancesBlock *big.Int
	var balancesSubmissionFrequency *big.Int
	var minipoolLaunchTimeout *big.Int
	var promotionScrubPeriodSeconds *big.Int
	var windowStartRaw *big.Int
	var windowLengthRaw *big.Int

	// Multicall getters
	mc, err := batch.NewMultiCaller(rp.Client, contracts.MulticallerAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating multicaller: %w", err)
	}
	core.AddCall(mc, contracts.RocketNetworkPrices, &details.RplPrice, "getRPLPrice")
	core.AddCall(mc, contracts.RocketDAOProtocolSettingsNode, &details.MinCollateralFraction, "getMinimumPerMinipoolStake")
	core.AddCall(mc, contracts.RocketDAOProtocolSettingsNode, &details.MaxCollateralFraction, "getMaximumPerMinipoolStake")
	core.AddCall(mc, contracts.RocketRewardsPool, &rewardIndex, "getRewardIndex")
	core.AddCall(mc, contracts.RocketRewardsPool, &intervalStart, "getClaimIntervalTimeStart")
	core.AddCall(mc, contracts.RocketRewardsPool, &intervalDuration, "getClaimIntervalTime")
	core.AddCall(mc, contracts.RocketRewardsPool, &details.NodeOperatorRewardsPercent, "getClaimingContractPerc", "rocketClaimNode")
	core.AddCall(mc, contracts.RocketRewardsPool, &details.OracleDaoRewardsPercent, "getClaimingContractPerc", "rocketClaimTrustedNode")
	core.AddCall(mc, contracts.RocketRewardsPool, &details.ProtocolDaoRewardsPercent, "getClaimingContractPerc", "rocketClaimDAO")
	core.AddCall(mc, contracts.RocketRewardsPool, &details.PendingRPLRewards, "getPendingRPLRewards")
	core.AddCall(mc, contracts.RocketDAONodeTrustedSettingsMinipool, &scrubPeriodSeconds, "getScrubPeriod")
	core.AddCall(mc, contracts.RocketDepositPool, &details.DepositPoolBalance, "getBalance")
	core.AddCall(mc, contracts.RocketDepositPool, &details.DepositPoolExcess, "getExcessBalance")
	core.AddCall(mc, contracts.RocketMinipoolQueue, &totalQueueCapacity, "getTotalCapacity")
	core.AddCall(mc, contracts.RocketMinipoolQueue, &effectiveQueueCapacity, "getEffectiveCapacity")
	core.AddCall(mc, contracts.RocketMinipoolQueue, &totalQueueLength, "getTotalLength")
	core.AddCall(mc, contracts.RocketTokenRPL, &details.RPLInflationIntervalRate, "getInflationIntervalRate")
	core.AddCall(mc, contracts.RocketTokenRPL, &details.RPLTotalSupply, "totalSupply")
	core.AddCall(mc, contracts.RocketNetworkPrices, &pricesBlock, "getPricesBlock")
	core.AddCall(mc, contracts.RocketNetworkBalances, &ethUtilizationRate, "getETHUtilizationRate")
	core.AddCall(mc, contracts.RocketNetworkBalances, &details.StakingETHBalance, "getStakingETHBalance")
	core.AddCall(mc, contracts.RocketTokenRETH, &rETHExchangeRate, "getExchangeRate")
	core.AddCall(mc, contracts.RocketNetworkBalances, &details.TotalETHBalance, "getTotalETHBalance")
	core.AddCall(mc, contracts.RocketTokenRETH, &details.TotalRETHSupply, "totalSupply")
	core.AddCall(mc, contracts.RocketNodeStaking, &details.TotalRPLStake, "getTotalRPLStake")
	core.AddCall(mc, contracts.RocketNetworkFees, &nodeFee, "getNodeFee")
	core.AddCall(mc, contracts.RocketNetworkBalances, &balancesBlock, "getBalancesBlock")
	core.AddCall(mc, contracts.RocketDAOProtocolSettingsNetwork, &details.SubmitBalancesEnabled, "getSubmitBalancesEnabled")
	core.AddCall(mc, contracts.RocketDAOProtocolSettingsNetwork, &details.SubmitPricesEnabled, "getSubmitPricesEnabled")
	core.AddCall(mc, contracts.RocketDAOProtocolSettingsMinipool, &minipoolLaunchTimeout, "getLaunchTimeout")

	// Atlas things
	core.AddCall(mc, contracts.RocketDAONodeTrustedSettingsMinipool, &promotionScrubPeriodSeconds, "getPromotionScrubPeriod")
	core.AddCall(mc, contracts.RocketDAONodeTrustedSettingsMinipool, &windowStartRaw, "getBondReductionWindowStart")
	core.AddCall(mc, contracts.RocketDAONodeTrustedSettingsMinipool, &windowLengthRaw, "getBondReductionWindowLength")
	core.AddCall(mc, contracts.RocketDepositPool, &details.DepositPoolUserBalance, "getUserBalance")

	// Houston
	if isHoustonDeployed {
		core.AddCall(mc, contracts.RocketDAOProtocolSettingsNetwork, &pricesSubmissionFrequency, "getSubmitPricesFrequency")
		core.AddCall(mc, contracts.RocketDAOProtocolSettingsNetwork, &balancesSubmissionFrequency, "getSubmitBalancesFrequency")
	} else {
		// getLatestReportableBlock was deprecated on Houston
		core.AddCall(mc, contracts.RocketNetworkPrices, &latestReportablePricesBlock, "getLatestReportableBlock")
		core.AddCall(mc, contracts.RocketNetworkBalances, &latestReportableBalancesBlock, "getLatestReportableBlock")
	}

	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return nil, fmt.Errorf("error executing multicall: %w", err)
	}

	// Conversion for raw parameters
	details.RewardIndex = rewardIndex.Uint64()
	details.IntervalStart = convertToTime(intervalStart)
	details.IntervalDuration = convertToDuration(intervalDuration)
	details.ScrubPeriod = convertToDuration(scrubPeriodSeconds)
	details.SmoothingPoolAddress = contracts.RocketSmoothingPool.Address
	details.TotalQueueCapacity = totalQueueCapacity
	details.EffectiveQueueCapacity = effectiveQueueCapacity
	details.QueueLength = totalQueueLength
	details.PricesBlock = pricesBlock.Uint64()
	if !isHoustonDeployed {
		details.LatestReportablePricesBlock = latestReportablePricesBlock.Uint64()
		details.LatestReportableBalancesBlock = latestReportableBalancesBlock.Uint64()
	} else {
		details.PricesSubmissionFrequency = pricesSubmissionFrequency.Uint64()
		details.BalancesSubmissionFrequency = balancesSubmissionFrequency.Uint64()
	}
	details.ETHUtilizationRate = eth.WeiToEth(ethUtilizationRate)
	details.RETHExchangeRate = eth.WeiToEth(rETHExchangeRate)
	details.NodeFee = eth.WeiToEth(nodeFee)
	details.BalancesBlock = balancesBlock
	details.MinipoolLaunchTimeout = minipoolLaunchTimeout
	details.PromotionScrubPeriod = convertToDuration(promotionScrubPeriodSeconds)
	details.BondReductionWindowStart = convertToDuration(windowStartRaw)
	details.BondReductionWindowLength = convertToDuration(windowLengthRaw)

	// Get various balances
	addresses := []common.Address{
		contracts.RocketSmoothingPool.Address,
		contracts.RocketTokenRETH.Address,
	}
	balances, err := contracts.BalanceBatcher.GetEthBalances(addresses, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting contract balances: %w", err)
	}
	details.SmoothingPoolBalance = balances[0]
	details.RETHBalance = balances[1]

	return details, nil
}

// Gets the details for a node using the efficient multicall contract
func GetTotalEffectiveRplStake(rp *rocketpool.RocketPool, contracts *NetworkContracts) (*big.Int, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	// Get the list of node addresses
	addresses, err := getNodeAddressesFast(rp, contracts, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node addresses: %w", err)
	}
	count := len(addresses)
	minimumStakes := make([]*big.Int, count)
	effectiveStakes := make([]*big.Int, count)

	// Sync
	var wg errgroup.Group
	wg.SetLimit(threadLimit)

	// Run the getters in batches
	for i := 0; i < count; i += networkEffectiveStakeBatchSize {
		i := i
		max := i + networkEffectiveStakeBatchSize
		if max > count {
			max = count
		}

		wg.Go(func() error {
			var err error
			mc, err := batch.NewMultiCaller(rp.Client, contracts.MulticallerAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				address := addresses[j]
				core.AddCall(mc, contracts.RocketNodeStaking, &minimumStakes[j], "getNodeMinimumRPLStake", address)
				core.AddCall(mc, contracts.RocketNodeStaking, &effectiveStakes[j], "getNodeEffectiveRPLStake", address)
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting effective stakes for all nodes: %w", err)
	}

	totalEffectiveStake := big.NewInt(0)
	for i, effectiveStake := range effectiveStakes {
		minimumStake := minimumStakes[i]
		// Fix the effective stake
		if effectiveStake.Cmp(minimumStake) >= 0 {
			totalEffectiveStake.Add(totalEffectiveStake, effectiveStake)
		}
	}

	return totalEffectiveStake, nil
}
