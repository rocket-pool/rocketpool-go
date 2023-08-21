package settings_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rocket-pool/rocketpool-go/settings"
	"github.com/rocket-pool/rocketpool-go/tests"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

var once sync.Once
var pdaoDefaults settings.ProtocolDaoSettingsDetails
var odaoDefaults settings.OracleDaoSettingsDetails

func createDefaults(mgr *tests.TestManager) error {
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
		pdaoDefaults = settings.ProtocolDaoSettingsDetails{}

		// Auction
		pdaoDefaults.Auction.IsCreateLotEnabled = true
		pdaoDefaults.Auction.IsBidOnLotEnabled = true
		pdaoDefaults.Auction.LotMinimumEthValue = eth.EthToWei(1)
		pdaoDefaults.Auction.LotMaximumEthValue = eth.EthToWei(10)
		pdaoDefaults.Auction.LotDuration.Set(40320)
		pdaoDefaults.Auction.LotStartingPriceRatio.Set(1)  // 100%
		pdaoDefaults.Auction.LotReservePriceRatio.Set(0.5) // 50%

		// Deposit
		pdaoDefaults.Deposit.IsDepositingEnabled = false
		pdaoDefaults.Deposit.AreDepositAssignmentsEnabled = true
		pdaoDefaults.Deposit.MinimumDeposit = eth.EthToWei(0.01)
		pdaoDefaults.Deposit.MaximumDepositPoolSize = eth.EthToWei(160)
		pdaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Set(90)
		pdaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Set(2)
		pdaoDefaults.Deposit.DepositFee.RawValue = big.NewInt(int64(1e18 * 5 / 10000)) // Set to approx. 1 day of rewards at 18.25% APR; supposed to be 0.0005 but have to do it via integer math cause of floating point errors

		// Inflation
		pdaoDefaults.Inflation.IntervalRate.RawValue = big.NewInt(1000133680617113500) // 5% annual calculated on a daily interval - Calculate in js example: let dailyInflation = web3.utils.toBN((1 + 0.05) ** (1 / (365)) * 1e18);
		pdaoDefaults.Inflation.StartTime.Set(startTime.Add(24 * time.Hour))            // Set the default start date for inflation to begin as 1 day after deployment

		// Minipool
		pdaoDefaults.Minipool.LaunchBalance = eth.EthToWei(32)
		pdaoDefaults.Minipool.PrelaunchValue = eth.EthToWei(1)
		pdaoDefaults.Minipool.FullDepositUserAmount = eth.EthToWei(16)
		pdaoDefaults.Minipool.HalfDepositUserAmount = eth.EthToWei(16)
		pdaoDefaults.Minipool.VariableDepositAmount = eth.EthToWei(31)
		pdaoDefaults.Minipool.IsSubmitWithdrawableEnabled = false
		pdaoDefaults.Minipool.IsBondReductionEnabled = true
		pdaoDefaults.Minipool.LaunchTimeout.Set(72 * time.Hour)
		pdaoDefaults.Minipool.MaximumCount.Set(14)
		pdaoDefaults.Minipool.UserDistributeWindowStart.Set(90 * 24 * time.Hour) // 90 days
		pdaoDefaults.Minipool.UserDistributeWindowLength.Set(2 * 24 * time.Hour) // 2 days

		// Network
		pdaoDefaults.Network.OracleDaoConsensusThreshold.Set(0.51) // 51%
		pdaoDefaults.Network.IsSubmitBalancesEnabled = true
		pdaoDefaults.Network.SubmitBalancesFrequency.Set(5760 * time.Second) // ~24 hours
		pdaoDefaults.Network.IsSubmitPricesEnabled = true
		pdaoDefaults.Network.SubmitPricesFrequency.Set(5760 * time.Second) // ~24 hours
		pdaoDefaults.Network.MinimumNodeFee.Set(0.14)
		pdaoDefaults.Network.TargetNodeFee.Set(0.14)
		pdaoDefaults.Network.MaximumNodeFee.Set(0.14)
		pdaoDefaults.Network.NodeFeeDemandRange = eth.EthToWei(160)
		pdaoDefaults.Network.TargetRethCollateralRate.Set(0.1)
		pdaoDefaults.Network.NodePenaltyThreshold.Set(0.51) // Consensus for penalties requires 51% vote
		pdaoDefaults.Network.PerPenaltyRate.Set(0.1)        // 10% per penalty
		pdaoDefaults.Network.RethDepositDelay.Set(0)
		pdaoDefaults.Network.IsSubmitRewardsEnabled = true

		// Node
		pdaoDefaults.Node.IsRegistrationEnabled = false
		pdaoDefaults.Node.IsSmoothingPoolRegistrationEnabled = true
		pdaoDefaults.Node.IsDepositingEnabled = false
		pdaoDefaults.Node.AreVacantMinipoolsEnabled = true
		pdaoDefaults.Node.MinimumPerMinipoolStake.Set(0.1) // 10% of user ETH value (matched ETH)
		pdaoDefaults.Node.MaximumPerMinipoolStake.Set(1.5) // 150% of node ETH value (provided ETH)

		// Rewards
		pdaoDefaults.Rewards.PercentageTotal.Set(1)
		pdaoDefaults.Rewards.IntervalTime.Set(28 * 24 * time.Hour) // 28 days

		// ==================
		// === Oracle DAO ===
		// ==================
		odaoDefaults = settings.OracleDaoSettingsDetails{}

		// Members
		odaoDefaults.Members.ChallengeCooldown.Set(7 * 24 * time.Hour) // 7 days
		odaoDefaults.Members.ChallengeCost = eth.EthToWei(1)
		odaoDefaults.Members.ChallengeWindow.Set(7 * 24 * time.Hour) // 7 days
		odaoDefaults.Members.Quorum.Set(0.51)
		odaoDefaults.Members.RplBond = eth.EthToWei(1750)
		odaoDefaults.Members.UnbondedMinipoolMax.Set(30)
		odaoDefaults.Members.UnbondedMinipoolMinFee.Set(0.8)

		// Minipools
		odaoDefaults.Minipools.BondReductionWindowStart.Set(12 * time.Hour)
		odaoDefaults.Minipools.BondReductionWindowLength.Set(2 * 24 * time.Hour)
		odaoDefaults.Minipools.IsScrubPenaltyEnabled = true
		odaoDefaults.Minipools.ScrubPeriod.Set(12 * time.Hour)
		odaoDefaults.Minipools.ScrubQuorum.Set(0.51)
		odaoDefaults.Minipools.PromotionScrubPeriod.Set(3 * 24 * time.Hour) // 3 days
		odaoDefaults.Minipools.BondReductionCancellationQuorum.Set(0.51)

		// Proposals
		odaoDefaults.Proposals.ActionTime.Set(4 * 7 * 24 * time.Hour)  // 4 weeks
		odaoDefaults.Proposals.CooldownTime.Set(2 * 24 * time.Hour)    // 2 days
		odaoDefaults.Proposals.ExecuteTime.Set(4 * 7 * 24 * time.Hour) // 4 weeks
		odaoDefaults.Proposals.VoteTime.Set(2 * 7 * 24 * time.Hour)    // 2 weeks
		odaoDefaults.Proposals.VoteDelayTime.Set(7 * 24 * time.Hour)   // 1 week
	})

	return err
}
