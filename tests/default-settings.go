package tests

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/v2/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/v2/dao/protocol"
)

var once sync.Once
var PDaoDefaults protocol.ProtocolDaoSettings
var ODaoDefaults oracle.OracleDaoSettings

func CreateDefaults(mgr *TestManager) error {
	var err error
	once.Do(func() {
		// Get the timestamp of the deploy block from hardhat - needed for inflation's default
		var fromBlock *big.Int
		err := mgr.RocketPool.Query(func(mc *batch.MultiCaller) error {
			mgr.RocketPool.Storage.GetDeployBlock(mc, &fromBlock)
			return nil
		}, nil)
		if err != nil {
			err = fmt.Errorf("error getting deployment block: %w", err)
			return
		}
		targetBlock := big.NewInt(0).Add(fromBlock, big.NewInt(34)) // Inflation timing started 34 blocks after the deploy block
		var header *types.Header
		header, err = mgr.Client.HeaderByNumber(context.Background(), targetBlock)
		if err != nil {
			err = fmt.Errorf("error getting header: %w", err)
			return
		}
		startTime := time.Unix(int64(header.Time), 0)

		// ====================
		// === Protocol DAO ===
		// ====================
		pdaoMgr, err := protocol.NewProtocolDaoManager(mgr.RocketPool)
		if err != nil {
			err = fmt.Errorf("error creating protocol DAO manager: %w", err)
			return
		}
		PDaoDefaults = *pdaoMgr.Settings

		// Auction
		PDaoDefaults.Auction.IsCreateLotEnabled.Set(true)
		PDaoDefaults.Auction.IsBidOnLotEnabled.Set(true)
		PDaoDefaults.Auction.LotMinimumEthValue.Set(eth.EthToWei(1))
		PDaoDefaults.Auction.LotMaximumEthValue.Set(eth.EthToWei(10))
		PDaoDefaults.Auction.LotDuration.Set(40320)
		PDaoDefaults.Auction.LotStartingPriceRatio.Set(1)  // 100%
		PDaoDefaults.Auction.LotReservePriceRatio.Set(0.5) // 50%

		// Deposit
		PDaoDefaults.Deposit.IsDepositingEnabled.Set(false)
		PDaoDefaults.Deposit.AreDepositAssignmentsEnabled.Set(true)
		PDaoDefaults.Deposit.MinimumDeposit.Set(eth.EthToWei(0.01))
		PDaoDefaults.Deposit.MaximumDepositPoolSize.Set(eth.EthToWei(160))
		PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Set(90)
		PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Set(2)
		PDaoDefaults.Deposit.DepositFee.SetRawValue(big.NewInt(int64(1e18 * 5 / 10000))) // Set to approx. 1 day of rewards at 18.25% APR; supposed to be 0.0005 but have to do it via integer math cause of floating point errors

		// Inflation
		PDaoDefaults.Inflation.IntervalRate.SetRawValue(big.NewInt(1000133680617113500)) // 5% annual calculated on a daily interval - Calculate in js example: let dailyInflation = web3.utils.toBN((1 + 0.05) ** (1 / (365)) * 1e18);
		PDaoDefaults.Inflation.StartTime.Set(startTime.Add(24 * time.Hour))              // Set the default start date for inflation to begin as 1 day after deployment

		// Minipool
		PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled.Set(false)
		PDaoDefaults.Minipool.IsBondReductionEnabled.Set(true)
		PDaoDefaults.Minipool.LaunchTimeout.Set(72 * time.Hour)
		PDaoDefaults.Minipool.MaximumCount.Set(14)
		PDaoDefaults.Minipool.UserDistributeWindowStart.Set(90 * 24 * time.Hour) // 90 days
		PDaoDefaults.Minipool.UserDistributeWindowLength.Set(2 * 24 * time.Hour) // 2 days

		// Network
		PDaoDefaults.Network.OracleDaoConsensusThreshold.Set(0.51) // 51%
		PDaoDefaults.Network.IsSubmitBalancesEnabled.Set(true)
		PDaoDefaults.Network.SubmitBalancesFrequency.Set(5760) // ~24 hours
		PDaoDefaults.Network.IsSubmitPricesEnabled.Set(true)
		PDaoDefaults.Network.SubmitPricesFrequency.Set(5760) // ~24 hours
		PDaoDefaults.Network.MinimumNodeFee.Set(0.14)
		PDaoDefaults.Network.TargetNodeFee.Set(0.14)
		PDaoDefaults.Network.MaximumNodeFee.Set(0.14)
		PDaoDefaults.Network.NodeFeeDemandRange.Set(eth.EthToWei(160))
		PDaoDefaults.Network.TargetRethCollateralRate.Set(0.1)
		PDaoDefaults.Network.NodePenaltyThreshold.Set(0.51) // Consensus for penalties requires 51% vote
		PDaoDefaults.Network.PerPenaltyRate.Set(0.1)        // 10% per penalty
		PDaoDefaults.Network.IsSubmitRewardsEnabled.Set(true)

		// Node
		PDaoDefaults.Node.IsRegistrationEnabled.Set(false)
		PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled.Set(true)
		PDaoDefaults.Node.IsDepositingEnabled.Set(false)
		PDaoDefaults.Node.AreVacantMinipoolsEnabled.Set(true)
		PDaoDefaults.Node.MinimumPerMinipoolStake.Set(0.1) // 10% of user ETH value (matched ETH)
		PDaoDefaults.Node.MaximumPerMinipoolStake.Set(1.5) // 150% of node ETH value (provided ETH)

		// Rewards
		PDaoDefaults.Rewards.IntervalPeriods.Set(28) // 28 periods

		// ==================
		// === Oracle DAO ===
		// ==================
		odaoMgr, err := oracle.NewOracleDaoManager(mgr.RocketPool)
		if err != nil {
			err = fmt.Errorf("error creating oracle DAO manager: %w", err)
			return
		}
		ODaoDefaults = *odaoMgr.Settings

		// Members
		ODaoDefaults.Member.ChallengeCooldown.Set(7 * 24 * time.Hour) // 7 days
		ODaoDefaults.Member.ChallengeCost.Set(eth.EthToWei(1))
		ODaoDefaults.Member.ChallengeWindow.Set(7 * 24 * time.Hour) // 7 days
		ODaoDefaults.Member.Quorum.Set(0.51)
		ODaoDefaults.Member.RplBond.Set(eth.EthToWei(1750))

		// Minipools
		ODaoDefaults.Minipool.BondReductionWindowStart.Set(2 * 24 * time.Hour)
		ODaoDefaults.Minipool.BondReductionWindowLength.Set(2 * 24 * time.Hour)
		ODaoDefaults.Minipool.IsScrubPenaltyEnabled.Set(false)
		ODaoDefaults.Minipool.ScrubPeriod.Set(12 * time.Hour)
		ODaoDefaults.Minipool.ScrubQuorum.Set(0.51)
		ODaoDefaults.Minipool.PromotionScrubPeriod.Set(3 * 24 * time.Hour) // 3 days
		ODaoDefaults.Minipool.BondReductionCancellationQuorum.Set(0.51)

		// Proposals
		ODaoDefaults.Proposal.ActionTime.Set(4 * 7 * 24 * time.Hour)  // 4 weeks
		ODaoDefaults.Proposal.CooldownTime.Set(2 * 24 * time.Hour)    // 2 days
		ODaoDefaults.Proposal.ExecuteTime.Set(4 * 7 * 24 * time.Hour) // 4 weeks
		ODaoDefaults.Proposal.VoteTime.Set(2 * 7 * 24 * time.Hour)    // 2 weeks
		ODaoDefaults.Proposal.VoteDelayTime.Set(7 * 24 * time.Hour)   // 1 week
	})

	return err
}
