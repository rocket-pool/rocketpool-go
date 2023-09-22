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
		PDaoDefaults = protocol.ProtocolDaoSettingsDetails{}

		// Auction
		PDaoDefaults.Auction.IsCreateLotEnabled = true
		PDaoDefaults.Auction.IsBidOnLotEnabled = true
		PDaoDefaults.Auction.LotMinimumEthValue = eth.EthToWei(1)
		PDaoDefaults.Auction.LotMaximumEthValue = eth.EthToWei(10)
		PDaoDefaults.Auction.LotDuration.Set(40320)
		PDaoDefaults.Auction.LotStartingPriceRatio.Set(1)  // 100%
		PDaoDefaults.Auction.LotReservePriceRatio.Set(0.5) // 50%

		// Deposit
		PDaoDefaults.Deposit.IsDepositingEnabled = false
		PDaoDefaults.Deposit.AreDepositAssignmentsEnabled = true
		PDaoDefaults.Deposit.MinimumDeposit = eth.EthToWei(0.01)
		PDaoDefaults.Deposit.MaximumDepositPoolSize = eth.EthToWei(160)
		PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Set(90)
		PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Set(2)
		PDaoDefaults.Deposit.DepositFee.RawValue = big.NewInt(int64(1e18 * 5 / 10000)) // Set to approx. 1 day of rewards at 18.25% APR; supposed to be 0.0005 but have to do it via integer math cause of floating point errors

		// Inflation
		PDaoDefaults.Inflation.IntervalRate.RawValue = big.NewInt(1000133680617113500) // 5% annual calculated on a daily interval - Calculate in js example: let dailyInflation = web3.utils.toBN((1 + 0.05) ** (1 / (365)) * 1e18);
		PDaoDefaults.Inflation.StartTime.Set(startTime.Add(24 * time.Hour))            // Set the default start date for inflation to begin as 1 day after deployment

		// Minipool
		PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled = false
		PDaoDefaults.Minipool.IsBondReductionEnabled = true
		PDaoDefaults.Minipool.LaunchTimeout.Set(72 * time.Hour)
		PDaoDefaults.Minipool.MaximumCount.Set(14)
		PDaoDefaults.Minipool.UserDistributeWindowStart.Set(90 * 24 * time.Hour) // 90 days
		PDaoDefaults.Minipool.UserDistributeWindowLength.Set(2 * 24 * time.Hour) // 2 days

		// Network
		PDaoDefaults.Network.OracleDaoConsensusThreshold.Set(0.51) // 51%
		PDaoDefaults.Network.IsSubmitBalancesEnabled = true
		PDaoDefaults.Network.SubmitBalancesFrequency.Set(5760 * time.Second) // ~24 hours
		PDaoDefaults.Network.IsSubmitPricesEnabled = true
		PDaoDefaults.Network.SubmitPricesFrequency.Set(5760 * time.Second) // ~24 hours
		PDaoDefaults.Network.MinimumNodeFee.Set(0.14)
		PDaoDefaults.Network.TargetNodeFee.Set(0.14)
		PDaoDefaults.Network.MaximumNodeFee.Set(0.14)
		PDaoDefaults.Network.NodeFeeDemandRange = eth.EthToWei(160)
		PDaoDefaults.Network.TargetRethCollateralRate.Set(0.1)
		PDaoDefaults.Network.NodePenaltyThreshold.Set(0.51) // Consensus for penalties requires 51% vote
		PDaoDefaults.Network.PerPenaltyRate.Set(0.1)        // 10% per penalty
		PDaoDefaults.Network.RethDepositDelay.Set(0)
		PDaoDefaults.Network.IsSubmitRewardsEnabled = true

		// Node
		PDaoDefaults.Node.IsRegistrationEnabled = false
		PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled = true
		PDaoDefaults.Node.IsDepositingEnabled = false
		PDaoDefaults.Node.AreVacantMinipoolsEnabled = true
		PDaoDefaults.Node.MinimumPerMinipoolStake.Set(0.1) // 10% of user ETH value (matched ETH)
		PDaoDefaults.Node.MaximumPerMinipoolStake.Set(1.5) // 150% of node ETH value (provided ETH)

		// Rewards
		PDaoDefaults.Rewards.IntervalTime.Set(28 * 24 * time.Hour) // 28 days

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
		ODaoDefaults.Members.ChallengeCooldown.Value.Set(7 * 24 * time.Hour) // 7 days
		ODaoDefaults.Members.ChallengeCost.Value = eth.EthToWei(1)
		ODaoDefaults.Members.ChallengeWindow.Value.Set(7 * 24 * time.Hour) // 7 days
		ODaoDefaults.Members.Quorum.Value.Set(0.51)
		ODaoDefaults.Members.RplBond.Value = eth.EthToWei(1750)
		ODaoDefaults.Members.UnbondedMinipoolMax.Value.Set(30)
		ODaoDefaults.Members.UnbondedMinipoolMinFee.Value.Set(0.8)

		// Minipools
		ODaoDefaults.Minipools.BondReductionWindowStart.Value.Set(12 * time.Hour)
		ODaoDefaults.Minipools.BondReductionWindowLength.Value.Set(2 * 24 * time.Hour)
		ODaoDefaults.Minipools.IsScrubPenaltyEnabled.Value = true
		ODaoDefaults.Minipools.ScrubPeriod.Value.Set(12 * time.Hour)
		ODaoDefaults.Minipools.ScrubQuorum.Value.Set(0.51)
		ODaoDefaults.Minipools.PromotionScrubPeriod.Value.Set(3 * 24 * time.Hour) // 3 days
		ODaoDefaults.Minipools.BondReductionCancellationQuorum.Value.Set(0.51)

		// Proposals
		ODaoDefaults.Proposals.ActionTime.Value.Set(4 * 7 * 24 * time.Hour)  // 4 weeks
		ODaoDefaults.Proposals.CooldownTime.Value.Set(2 * 24 * time.Hour)    // 2 days
		ODaoDefaults.Proposals.ExecuteTime.Value.Set(4 * 7 * 24 * time.Hour) // 4 weeks
		ODaoDefaults.Proposals.VoteTime.Value.Set(2 * 7 * 24 * time.Hour)    // 2 weeks
		ODaoDefaults.Proposals.VoteDelayTime.Value.Set(7 * 24 * time.Hour)   // 1 week
	})

	return err
}
