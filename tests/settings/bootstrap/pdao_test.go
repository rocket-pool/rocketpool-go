package bootstrap_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/tests"
	settings_test "github.com/rocket-pool/rocketpool-go/tests/settings"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"golang.org/x/sync/errgroup"
)

func Test_BootstrapCreateAuctionLotEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Auction.IsCreateLotEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.IsCreateLotEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapCreateAuctionLotEnabled(newVal, opts)
	})
}

func Test_BootstrapBidOnAuctionLotEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Auction.IsBidOnLotEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.IsBidOnLotEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapBidOnAuctionLotEnabled(newVal, opts)
	})
}

func Test_BootstrapAuctionLotMinimumEthValue(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMinimumEthValue, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotMinimumEthValue = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapAuctionLotMinimumEthValue(newVal, opts)
	})
}

func Test_BootstrapAuctionLotMaximumEthValue(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMaximumEthValue, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotMaximumEthValue = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapAuctionLotMaximumEthValue(newVal, opts)
	})
}

func Test_BootstrapAuctionLotDuration(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Auction.LotDuration.Formatted() + 1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotDuration.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapAuctionLotDuration(newVal, opts)
	})
}

func Test_BootstrapAuctionLotStartingPriceRatio(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Auction.LotStartingPriceRatio.Formatted() - 0.2)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotStartingPriceRatio.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapAuctionLotStartingPriceRatio(newVal, opts)
	})
}

func Test_BootstrapAuctionLotReservePriceRatio(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotReservePriceRatio.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapAuctionLotReservePriceRatio(newVal, opts)
	})
}

func Test_BootstrapPoolDepositEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Deposit.IsDepositingEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.IsDepositingEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapPoolDepositEnabled(newVal, opts)
	})
}

func Test_BootstrapAssignPoolDepositsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Deposit.AreDepositAssignmentsEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.AreDepositAssignmentsEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapAssignPoolDepositsEnabled(newVal, opts)
	})
}

func Test_BootstrapMinimumPoolDeposit(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MinimumDeposit, eth.EthToWei(0.01))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MinimumDeposit = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMinimumPoolDeposit(newVal, opts)
	})
}

func Test_BootstrapMaximumDepositPoolSize(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MaximumDepositPoolSize, eth.EthToWei(100))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumDepositPoolSize = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMaximumDepositPoolSize(newVal, opts)
	})
}

func Test_BootstrapMaximumPoolDepositAssignments(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Formatted() + 10)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumAssignmentsPerDeposit.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMaximumPoolDepositAssignments(newVal, opts)
	})
}

func Test_BootstrapMaximumSocialisedPoolDepositAssignments(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Formatted() + 5)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMaximumSocialisedPoolDepositAssignments(newVal, opts)
	})
}

func Test_BootstrapDepositFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.RawValue = big.NewInt(0).Add(tests.PDaoDefaults.Deposit.DepositFee.RawValue, eth.EthToWei(0.1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.DepositFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapDepositFee(newVal, opts)
	})
}

func Test_BootstrapInflationIntervalRate(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.RawValue = big.NewInt(0).Add(tests.PDaoDefaults.Inflation.IntervalRate.RawValue, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Inflation.IntervalRate.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapInflationIntervalRate(newVal, opts)
	})
}

func Test_BootstrapInflationIntervalStartTime(t *testing.T) {
	newVal := core.Parameter[time.Time]{}
	newVal.Set(tests.PDaoDefaults.Inflation.StartTime.Formatted().Add(24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Inflation.StartTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapInflationIntervalStartTime(newVal, opts)
	})
}

func Test_BootstrapSubmitWithdrawableEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.IsSubmitWithdrawableEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapSubmitWithdrawableEnabled(newVal, opts)
	})
}

func Test_BootstrapBondReductionEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Minipool.IsBondReductionEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.IsBondReductionEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapBondReductionEnabled(newVal, opts)
	})
}

func Test_BootstrapMinipoolLaunchTimeout(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Minipool.LaunchTimeout.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.LaunchTimeout.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMinipoolLaunchTimeout(newVal, opts)
	})
}

func Test_BootstrapMaximumMinipoolCount(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Minipool.MaximumCount.Formatted() + 1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.MaximumCount.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMaximumMinipoolCount(newVal, opts)
	})
}

func Test_BootstrapUserDistributeWindowStart(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowStart.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.UserDistributeWindowStart.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapUserDistributeWindowStart(newVal, opts)
	})
}

func Test_BootstrapUserDistributeWindowLength(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowLength.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.UserDistributeWindowLength.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapUserDistributeWindowLength(newVal, opts)
	})
}

func Test_BootstrapOracleDaoConsensusThreshold(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.OracleDaoConsensusThreshold.Formatted() + 0.15)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.OracleDaoConsensusThreshold.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapOracleDaoConsensusThreshold(newVal, opts)
	})
}

func Test_BootstrapSubmitBalancesEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitBalancesEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitBalancesEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapSubmitBalancesEnabled(newVal, opts)
	})
}

func Test_BootstrapSubmitBalancesFrequency(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Network.SubmitBalancesFrequency.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.SubmitBalancesFrequency.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapSubmitBalancesFrequency(newVal, opts)
	})
}

func Test_BootstrapSubmitPricesEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitPricesEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitPricesEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapSubmitPricesEnabled(newVal, opts)
	})
}

func Test_BootstrapSubmitPricesFrequency(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Network.SubmitPricesFrequency.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.SubmitPricesFrequency.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapSubmitPricesFrequency(newVal, opts)
	})
}

func Test_BootstrapMinimumNodeFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.MinimumNodeFee.Formatted() - 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.MinimumNodeFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMinimumNodeFee(newVal, opts)
	})
}

func Test_BootstrapTargetNodeFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.TargetNodeFee.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.TargetNodeFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapTargetNodeFee(newVal, opts)
	})
}

func Test_BootstrapMaximumNodeFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.MaximumNodeFee.Formatted() + 0.3)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.MaximumNodeFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMaximumNodeFee(newVal, opts)
	})
}

func Test_BootstrapNodeFeeDemandRange(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Network.NodeFeeDemandRange, eth.EthToWei(100))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.NodeFeeDemandRange = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapNodeFeeDemandRange(newVal, opts)
	})
}

func Test_BootstrapTargetRethCollateralRate(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.TargetRethCollateralRate.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.TargetRethCollateralRate.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapTargetRethCollateralRate(newVal, opts)
	})
}

func Test_BootstrapNodePenaltyThreshold(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.NodePenaltyThreshold.Formatted() + 0.15)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.NodePenaltyThreshold.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapNodePenaltyThreshold(newVal, opts)
	})
}

func Test_BootstrapPerPenaltyRate(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.PerPenaltyRate.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.PerPenaltyRate.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapPerPenaltyRate(newVal, opts)
	})
}

func Test_BootstrapRethDepositDelay(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Network.RethDepositDelay.Formatted() + time.Hour)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.RethDepositDelay.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapRethDepositDelay(newVal, opts)
	})
}

func Test_BootstrapSubmitRewardsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitRewardsEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitRewardsEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapSubmitRewardsEnabled(newVal, opts)
	})
}

func Test_BootstrapNodeRegistrationEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsRegistrationEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsRegistrationEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapNodeRegistrationEnabled(newVal, opts)
	})
}

func Test_BootstrapSmoothingPoolRegistrationEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsSmoothingPoolRegistrationEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapSmoothingPoolRegistrationEnabled(newVal, opts)
	})
}

func Test_BootstrapNodeDepositEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsDepositingEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsDepositingEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapNodeDepositEnabled(newVal, opts)
	})
}

func Test_BootstrapVacantMinipoolsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.AreVacantMinipoolsEnabled
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.AreVacantMinipoolsEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapVacantMinipoolsEnabled(newVal, opts)
	})
}

func Test_BootstrapMinimumPerMinipoolStake(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Node.MinimumPerMinipoolStake.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.MinimumPerMinipoolStake.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMinimumPerMinipoolStake(newVal, opts)
	})
}

func Test_BootstrapMaximumPerMinipoolStake(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Node.MaximumPerMinipoolStake.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.MaximumPerMinipoolStake.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapMaximumPerMinipoolStake(newVal, opts)
	})
}

func Test_BootstrapRewardsIntervalTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Rewards.IntervalTime.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Rewards.IntervalTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.BootstrapRewardsIntervalTime(newVal, opts)
	})
}

func Test_AllPDaoBoostrapFunctions(t *testing.T) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Create new settings
	newPdaoSettings := protocol.ProtocolDaoSettingsDetails{}
	newPdaoSettings.Auction.IsCreateLotEnabled = !tests.PDaoDefaults.Auction.IsCreateLotEnabled
	newPdaoSettings.Auction.IsBidOnLotEnabled = !tests.PDaoDefaults.Auction.IsBidOnLotEnabled
	newPdaoSettings.Auction.LotMinimumEthValue = big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMinimumEthValue, eth.EthToWei(1))
	newPdaoSettings.Auction.LotMaximumEthValue = big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMaximumEthValue, eth.EthToWei(1))
	newPdaoSettings.Auction.LotDuration.Set(tests.PDaoDefaults.Auction.LotDuration.Formatted() + 1)
	newPdaoSettings.Auction.LotStartingPriceRatio.Set(tests.PDaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.2)
	newPdaoSettings.Auction.LotReservePriceRatio.Set(tests.PDaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.1)
	newPdaoSettings.Deposit.IsDepositingEnabled = !tests.PDaoDefaults.Deposit.IsDepositingEnabled
	newPdaoSettings.Deposit.AreDepositAssignmentsEnabled = !tests.PDaoDefaults.Deposit.AreDepositAssignmentsEnabled
	newPdaoSettings.Deposit.MinimumDeposit = big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MinimumDeposit, eth.EthToWei(0.01))
	newPdaoSettings.Deposit.MaximumDepositPoolSize = big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MaximumDepositPoolSize, eth.EthToWei(100))
	newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit.Set(tests.PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Formatted() + 10)
	newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Set(tests.PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Formatted() + 5)
	newPdaoSettings.Deposit.DepositFee.RawValue = big.NewInt(0).Add(tests.PDaoDefaults.Deposit.DepositFee.RawValue, eth.EthToWei(0.1))
	newPdaoSettings.Inflation.IntervalRate.RawValue = big.NewInt(0).Add(tests.PDaoDefaults.Inflation.IntervalRate.RawValue, eth.EthToWei(1))
	newPdaoSettings.Inflation.StartTime.Set(tests.PDaoDefaults.Inflation.StartTime.Formatted().Add(24 * time.Hour))
	newPdaoSettings.Minipool.IsSubmitWithdrawableEnabled = !tests.PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled
	newPdaoSettings.Minipool.IsBondReductionEnabled = !tests.PDaoDefaults.Minipool.IsBondReductionEnabled
	newPdaoSettings.Minipool.LaunchTimeout.Set(tests.PDaoDefaults.Minipool.LaunchTimeout.Formatted() + (24 * time.Hour))
	newPdaoSettings.Minipool.MaximumCount.Set(tests.PDaoDefaults.Minipool.MaximumCount.Formatted() + 1)
	newPdaoSettings.Minipool.UserDistributeWindowStart.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowStart.Formatted() + (24 * time.Hour))
	newPdaoSettings.Minipool.UserDistributeWindowLength.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowLength.Formatted() + (24 * time.Hour))
	newPdaoSettings.Network.OracleDaoConsensusThreshold.Set(tests.PDaoDefaults.Network.OracleDaoConsensusThreshold.Formatted() + 0.15)
	newPdaoSettings.Network.IsSubmitBalancesEnabled = !tests.PDaoDefaults.Network.IsSubmitBalancesEnabled
	newPdaoSettings.Network.SubmitBalancesFrequency.Set(tests.PDaoDefaults.Network.SubmitBalancesFrequency.Formatted() + (24 * time.Hour))
	newPdaoSettings.Network.IsSubmitPricesEnabled = !tests.PDaoDefaults.Network.IsSubmitPricesEnabled
	newPdaoSettings.Network.SubmitPricesFrequency.Set(tests.PDaoDefaults.Network.SubmitPricesFrequency.Formatted() + (24 * time.Hour))
	newPdaoSettings.Network.MinimumNodeFee.Set(tests.PDaoDefaults.Network.MinimumNodeFee.Formatted() - 0.1)
	newPdaoSettings.Network.TargetNodeFee.Set(tests.PDaoDefaults.Network.TargetNodeFee.Formatted() + 0.1)
	newPdaoSettings.Network.MaximumNodeFee.Set(tests.PDaoDefaults.Network.MaximumNodeFee.Formatted() + 0.3)
	newPdaoSettings.Network.NodeFeeDemandRange = big.NewInt(0).Add(tests.PDaoDefaults.Network.NodeFeeDemandRange, eth.EthToWei(100))
	newPdaoSettings.Network.TargetRethCollateralRate.Set(tests.PDaoDefaults.Network.TargetRethCollateralRate.Formatted() + 0.1)
	newPdaoSettings.Network.NodePenaltyThreshold.Set(tests.PDaoDefaults.Network.NodePenaltyThreshold.Formatted() + 0.15)
	newPdaoSettings.Network.PerPenaltyRate.Set(tests.PDaoDefaults.Network.PerPenaltyRate.Formatted() + 0.1)
	newPdaoSettings.Network.RethDepositDelay.Set(tests.PDaoDefaults.Network.RethDepositDelay.Formatted() + time.Hour)
	newPdaoSettings.Network.IsSubmitRewardsEnabled = !tests.PDaoDefaults.Network.IsSubmitRewardsEnabled
	newPdaoSettings.Node.IsRegistrationEnabled = !tests.PDaoDefaults.Node.IsRegistrationEnabled
	newPdaoSettings.Node.IsSmoothingPoolRegistrationEnabled = !tests.PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled
	newPdaoSettings.Node.IsDepositingEnabled = !tests.PDaoDefaults.Node.IsDepositingEnabled
	newPdaoSettings.Node.AreVacantMinipoolsEnabled = !tests.PDaoDefaults.Node.AreVacantMinipoolsEnabled
	newPdaoSettings.Node.MinimumPerMinipoolStake.Set(tests.PDaoDefaults.Node.MinimumPerMinipoolStake.Formatted() + 0.1)
	newPdaoSettings.Node.MaximumPerMinipoolStake.Set(tests.PDaoDefaults.Node.MaximumPerMinipoolStake.Formatted() + 0.1)
	newPdaoSettings.Rewards.IntervalTime.Set(tests.PDaoDefaults.Rewards.IntervalTime.Formatted() + (24 * time.Hour))

	// Ensure they're all different from the default
	settings_test.EnsureDifferentDetails(t.Fatalf, &tests.PDaoDefaults, &newPdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	txInfos := []*core.TransactionInfo{}
	bootstrappers := []func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapCreateAuctionLotEnabled(newPdaoSettings.Auction.IsCreateLotEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapBidOnAuctionLotEnabled(newPdaoSettings.Auction.IsBidOnLotEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapAuctionLotMinimumEthValue(newPdaoSettings.Auction.LotMinimumEthValue, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapAuctionLotMaximumEthValue(newPdaoSettings.Auction.LotMaximumEthValue, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapAuctionLotDuration(newPdaoSettings.Auction.LotDuration, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapAuctionLotStartingPriceRatio(newPdaoSettings.Auction.LotStartingPriceRatio, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapAuctionLotReservePriceRatio(newPdaoSettings.Auction.LotReservePriceRatio, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapPoolDepositEnabled(newPdaoSettings.Deposit.IsDepositingEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapAssignPoolDepositsEnabled(newPdaoSettings.Deposit.AreDepositAssignmentsEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMinimumPoolDeposit(newPdaoSettings.Deposit.MinimumDeposit, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMaximumDepositPoolSize(newPdaoSettings.Deposit.MaximumDepositPoolSize, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMaximumPoolDepositAssignments(newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMaximumSocialisedPoolDepositAssignments(newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapDepositFee(newPdaoSettings.Deposit.DepositFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapInflationIntervalRate(newPdaoSettings.Inflation.IntervalRate, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapInflationIntervalStartTime(newPdaoSettings.Inflation.StartTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapSubmitWithdrawableEnabled(newPdaoSettings.Minipool.IsSubmitWithdrawableEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapBondReductionEnabled(newPdaoSettings.Minipool.IsBondReductionEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMinipoolLaunchTimeout(newPdaoSettings.Minipool.LaunchTimeout, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMaximumMinipoolCount(newPdaoSettings.Minipool.MaximumCount, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapUserDistributeWindowStart(newPdaoSettings.Minipool.UserDistributeWindowStart, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapUserDistributeWindowLength(newPdaoSettings.Minipool.UserDistributeWindowLength, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapOracleDaoConsensusThreshold(newPdaoSettings.Network.OracleDaoConsensusThreshold, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapSubmitBalancesEnabled(newPdaoSettings.Network.IsSubmitBalancesEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapSubmitBalancesFrequency(newPdaoSettings.Network.SubmitPricesFrequency, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapSubmitPricesEnabled(newPdaoSettings.Network.IsSubmitPricesEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapSubmitPricesFrequency(newPdaoSettings.Network.SubmitPricesFrequency, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMinimumNodeFee(newPdaoSettings.Network.MinimumNodeFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapTargetNodeFee(newPdaoSettings.Network.TargetNodeFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMaximumNodeFee(newPdaoSettings.Network.MaximumNodeFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapNodeFeeDemandRange(newPdaoSettings.Network.NodeFeeDemandRange, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapTargetRethCollateralRate(newPdaoSettings.Network.TargetRethCollateralRate, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapNodePenaltyThreshold(newPdaoSettings.Network.NodePenaltyThreshold, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapPerPenaltyRate(newPdaoSettings.Network.PerPenaltyRate, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapRethDepositDelay(newPdaoSettings.Network.RethDepositDelay, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapSubmitRewardsEnabled(newPdaoSettings.Network.IsSubmitRewardsEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapNodeRegistrationEnabled(newPdaoSettings.Node.IsRegistrationEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapSmoothingPoolRegistrationEnabled(newPdaoSettings.Node.IsSmoothingPoolRegistrationEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapNodeDepositEnabled(newPdaoSettings.Node.IsDepositingEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapVacantMinipoolsEnabled(newPdaoSettings.Node.AreVacantMinipoolsEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMinimumPerMinipoolStake(newPdaoSettings.Node.MinimumPerMinipoolStake, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapMaximumPerMinipoolStake(newPdaoSettings.Node.MaximumPerMinipoolStake, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.BootstrapRewardsIntervalTime(newPdaoSettings.Rewards.IntervalTime, opts)
		},
	}
	for i, bootstrapper := range bootstrappers {
		txInfo, err := bootstrapper()
		if err != nil {
			t.Fatalf("error running boostrapper %d: %s", i, err.Error())
			continue
		}
		if txInfo.SimError != "" {
			t.Fatalf("error simming boostrapper %d: %s", i, txInfo.SimError)
		}
		txInfos = append(txInfos, txInfo)
	}
	t.Log("Bootstrappers constructed")

	// Run the transactions
	txs, err := rp.BatchExecuteTransactions(txInfos, opts)
	if err != nil {
		t.Fatalf("error submitting transactions: %s", err.Error())
	}
	t.Log("Bootstrappers submitted")

	var wg errgroup.Group
	for _, tx := range txs {
		tx := tx
		wg.Go(func() error {
			return rp.WaitForTransaction(tx)
		})
	}

	err = wg.Wait()
	if err != nil {
		t.Fatalf("error waiting for transactions: %s", err.Error())
	}
	t.Log("Bootstrappers executed")

	// Get new values and make sure they match
	err = rp.Query(func(mc *batch.MultiCaller) error {
		pdaoMgr.Settings.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	settings_test.EnsureSameDetails(t.Fatalf, &newPdaoSettings, &pdaoMgr.Settings.Details)
	t.Log("New settings match expected settings")
}

func testPdaoParameterBootstrap(t *testing.T, setter func(*protocol.ProtocolDaoSettingsDetails), bootstrapper func() (*core.TransactionInfo, error)) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Get the original settings
	var settings protocol.ProtocolDaoSettingsDetails
	settings_test.Clone(t, &tests.PDaoDefaults, &settings)
	pass := settings_test.EnsureSameDetails(t.Errorf, &tests.PDaoDefaults, &settings)
	if !pass {
		t.Fatalf("Details differed unexpectedly, can't continue")
	}
	t.Log("Cloned default settings")

	// Set the new setting
	setter(&settings)
	pass = settings_test.EnsureSameDetails(t.Logf, &tests.PDaoDefaults, &settings)
	if pass {
		t.Fatalf("Details were same, ineffective setter")
	}
	t.Log("Applied new setting")

	// Run the bootstrapper
	txInfo, err := bootstrapper()
	if err != nil {
		t.Fatalf("error running boostrapper: %s", err.Error())
	}
	if txInfo.SimError != "" {
		t.Fatalf("error simming boostrapper: %s", txInfo.SimError)
	}
	t.Log("Bootstrapper constructed")

	opts := mgr.OwnerAccount.Transactor
	tx, err := rp.ExecuteTransaction(txInfo, opts)
	if err != nil {
		t.Fatalf("error submitting transaction: %s", err.Error())
	}
	t.Log("Bootstrapper submitted")

	err = rp.WaitForTransaction(tx)
	if err != nil {
		t.Fatalf("error waiting for transactions: %s", err.Error())
	}
	t.Log("Bootstrapper executed")

	// Get new values and make sure they match
	err = rp.Query(func(mc *batch.MultiCaller) error {
		pdaoMgr.Settings.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	settings_test.EnsureSameDetails(t.Fatalf, &settings, &pdaoMgr.Settings.Details)
	t.Log("New settings match expected settings")
}
