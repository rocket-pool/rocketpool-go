package tests

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rocket-pool/rocketpool-go/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

var once sync.Once
var PDaoDefaults protocol.ProtocolDaoSettingsDetails
var ODaoDefaults oracle.OracleDaoSettingsDetails

func CreateDefaults(mgr *TestManager) error {
	var err error
	once.Do(func() {
		// Get the timestamp of the deploy block from hardhat - needed for inflation's default
		targetBlock := big.NewInt(0).Add(mgr.RocketPool.DeployBlock, big.NewInt(34)) // Inflation timing started 34 blocks after the deploy block
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
		PDaoDefaults = *pdaoMgr.Settings.ProtocolDaoSettingsDetails

		// Auction
		PDaoDefaults.Auction.IsCreateLotEnabled.Value = true
		PDaoDefaults.Auction.IsBidOnLotEnabled.Value = true
		PDaoDefaults.Auction.LotMinimumEthValue.Value = eth.EthToWei(1)
		PDaoDefaults.Auction.LotMaximumEthValue.Value = eth.EthToWei(10)
		PDaoDefaults.Auction.LotDuration.Value.Set(40320)
		PDaoDefaults.Auction.LotStartingPriceRatio.Value.Set(1)  // 100%
		PDaoDefaults.Auction.LotReservePriceRatio.Value.Set(0.5) // 50%

		// Deposit
		PDaoDefaults.Deposit.IsDepositingEnabled.Value = false
		PDaoDefaults.Deposit.AreDepositAssignmentsEnabled.Value = true
		PDaoDefaults.Deposit.MinimumDeposit.Value = eth.EthToWei(0.01)
		PDaoDefaults.Deposit.MaximumDepositPoolSize.Value = eth.EthToWei(160)
		PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Value.Set(90)
		PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Value.Set(2)
		PDaoDefaults.Deposit.DepositFee.Value.RawValue = big.NewInt(int64(1e18 * 5 / 10000)) // Set to approx. 1 day of rewards at 18.25% APR; supposed to be 0.0005 but have to do it via integer math cause of floating point errors

		// Inflation
		PDaoDefaults.Inflation.IntervalRate.Value.RawValue = big.NewInt(1000133680617113500) // 5% annual calculated on a daily interval - Calculate in js example: let dailyInflation = web3.utils.toBN((1 + 0.05) ** (1 / (365)) * 1e18);
		PDaoDefaults.Inflation.StartTime.Value.Set(startTime.Add(24 * time.Hour))            // Set the default start date for inflation to begin as 1 day after deployment

		// Minipool
		PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled.Value = false
		PDaoDefaults.Minipool.IsBondReductionEnabled.Value = true
		PDaoDefaults.Minipool.LaunchTimeout.Value.Set(72 * time.Hour)
		PDaoDefaults.Minipool.MaximumCount.Value.Set(14)
		PDaoDefaults.Minipool.UserDistributeWindowStart.Value.Set(90 * 24 * time.Hour) // 90 days
		PDaoDefaults.Minipool.UserDistributeWindowLength.Value.Set(2 * 24 * time.Hour) // 2 days

		// Network
		PDaoDefaults.Network.OracleDaoConsensusThreshold.Value.Set(0.51) // 51%
		PDaoDefaults.Network.IsSubmitBalancesEnabled.Value = true
		PDaoDefaults.Network.SubmitBalancesFrequency.Value.Set(5760) // ~24 hours
		PDaoDefaults.Network.IsSubmitPricesEnabled.Value = true
		PDaoDefaults.Network.SubmitPricesFrequency.Value.Set(5760) // ~24 hours
		PDaoDefaults.Network.MinimumNodeFee.Value.Set(0.14)
		PDaoDefaults.Network.TargetNodeFee.Value.Set(0.14)
		PDaoDefaults.Network.MaximumNodeFee.Value.Set(0.14)
		PDaoDefaults.Network.NodeFeeDemandRange.Value = eth.EthToWei(160)
		PDaoDefaults.Network.TargetRethCollateralRate.Value.Set(0.1)
		PDaoDefaults.Network.NodePenaltyThreshold.Value.Set(0.51) // Consensus for penalties requires 51% vote
		PDaoDefaults.Network.PerPenaltyRate.Value.Set(0.1)        // 10% per penalty
		PDaoDefaults.Network.RethDepositDelay.Value.Set(0)
		PDaoDefaults.Network.IsSubmitRewardsEnabled.Value = true

		// Node
		PDaoDefaults.Node.IsRegistrationEnabled.Value = false
		PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled.Value = true
		PDaoDefaults.Node.IsDepositingEnabled.Value = false
		PDaoDefaults.Node.AreVacantMinipoolsEnabled.Value = true
		PDaoDefaults.Node.MinimumPerMinipoolStake.Value.Set(0.1) // 10% of user ETH value (matched ETH)
		PDaoDefaults.Node.MaximumPerMinipoolStake.Value.Set(1.5) // 150% of node ETH value (provided ETH)

		// Rewards
		PDaoDefaults.Rewards.IntervalTime.Value.Set(28 * 24 * time.Hour) // 28 days

		// ==================
		// === Oracle DAO ===
		// ==================
		odaoMgr, err := oracle.NewOracleDaoManager(mgr.RocketPool)
		if err != nil {
			err = fmt.Errorf("error creating oracle DAO manager: %w", err)
			return
		}
		ODaoDefaults = *odaoMgr.Settings.OracleDaoSettingsDetails

		// Members
		ODaoDefaults.Member.ChallengeCooldown.Value.Set(7 * 24 * time.Hour) // 7 days
		ODaoDefaults.Member.ChallengeCost.Value = eth.EthToWei(1)
		ODaoDefaults.Member.ChallengeWindow.Value.Set(7 * 24 * time.Hour) // 7 days
		ODaoDefaults.Member.Quorum.Value.Set(0.51)
		ODaoDefaults.Member.RplBond.Value = eth.EthToWei(1750)
		ODaoDefaults.Member.UnbondedMinipoolMax.Value.Set(30)
		ODaoDefaults.Member.UnbondedMinipoolMinFee.Value.Set(0.8)

		// Minipools
		ODaoDefaults.Minipool.BondReductionWindowStart.Value.Set(12 * time.Hour)
		ODaoDefaults.Minipool.BondReductionWindowLength.Value.Set(2 * 24 * time.Hour)
		ODaoDefaults.Minipool.IsScrubPenaltyEnabled.Value = true
		ODaoDefaults.Minipool.ScrubPeriod.Value.Set(12 * time.Hour)
		ODaoDefaults.Minipool.ScrubQuorum.Value.Set(0.51)
		ODaoDefaults.Minipool.PromotionScrubPeriod.Value.Set(3 * 24 * time.Hour) // 3 days
		ODaoDefaults.Minipool.BondReductionCancellationQuorum.Value.Set(0.51)

		// Proposals
		ODaoDefaults.Proposal.ActionTime.Value.Set(4 * 7 * 24 * time.Hour)  // 4 weeks
		ODaoDefaults.Proposal.CooldownTime.Value.Set(2 * 24 * time.Hour)    // 2 days
		ODaoDefaults.Proposal.ExecuteTime.Value.Set(4 * 7 * 24 * time.Hour) // 4 weeks
		ODaoDefaults.Proposal.VoteTime.Value.Set(2 * 7 * 24 * time.Hour)    // 2 weeks
		ODaoDefaults.Proposal.VoteDelayTime.Value.Set(7 * 24 * time.Hour)   // 1 week
	})

	return err
}
