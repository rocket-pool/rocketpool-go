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
	newVal := !tests.PDaoDefaults.Auction.IsCreateLotEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.IsCreateLotEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.IsCreateLotEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBidOnAuctionLotEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Auction.IsBidOnLotEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.IsBidOnLotEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.IsBidOnLotEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAuctionLotMinimumEthValue(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMinimumEthValue.Value, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotMinimumEthValue.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotMinimumEthValue.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAuctionLotMaximumEthValue(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMaximumEthValue.Value, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotMaximumEthValue.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotMaximumEthValue.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAuctionLotDuration(t *testing.T) {
	newVal := core.Uint256Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Auction.LotDuration.Value.Formatted() + 1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotDuration.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotDuration.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAuctionLotStartingPriceRatio(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Auction.LotStartingPriceRatio.Value.Formatted() - 0.2)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotStartingPriceRatio.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotStartingPriceRatio.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAuctionLotReservePriceRatio(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Auction.LotReservePriceRatio.Value.Formatted() - 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotReservePriceRatio.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotReservePriceRatio.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapPoolDepositEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Deposit.IsDepositingEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.IsDepositingEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.IsDepositingEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAssignPoolDepositsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Deposit.AreDepositAssignmentsEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.AreDepositAssignmentsEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.AreDepositAssignmentsEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMinimumPoolDeposit(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MinimumDeposit.Value, eth.EthToWei(0.01))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MinimumDeposit.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.MinimumDeposit.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMaximumDepositPoolSize(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MaximumDepositPoolSize.Value, eth.EthToWei(100))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumDepositPoolSize.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.MaximumDepositPoolSize.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMaximumPoolDepositAssignments(t *testing.T) {
	newVal := core.Uint256Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Value.Formatted() + 10)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumAssignmentsPerDeposit.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.MaximumAssignmentsPerDeposit.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMaximumSocialisedPoolDepositAssignments(t *testing.T) {
	newVal := core.Uint256Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Value.Formatted() + 5)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapDepositFee(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.RawValue = big.NewInt(0).Add(tests.PDaoDefaults.Deposit.DepositFee.Value.RawValue, eth.EthToWei(0.1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.DepositFee.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.DepositFee.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapInflationIntervalRate(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.RawValue = big.NewInt(0).Add(tests.PDaoDefaults.Inflation.IntervalRate.Value.RawValue, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Inflation.IntervalRate.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Inflation.IntervalRate.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapInflationIntervalStartTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Time]{}
	newVal.Set(tests.PDaoDefaults.Inflation.StartTime.Value.Formatted().Add(24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Inflation.StartTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Inflation.StartTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSubmitWithdrawableEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.IsSubmitWithdrawableEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.IsSubmitWithdrawableEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBondReductionEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Minipool.IsBondReductionEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.IsBondReductionEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.IsBondReductionEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMinipoolLaunchTimeout(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Minipool.LaunchTimeout.Value.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.LaunchTimeout.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.LaunchTimeout.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMaximumMinipoolCount(t *testing.T) {
	newVal := core.Uint256Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Minipool.MaximumCount.Value.Formatted() + 1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.MaximumCount.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.MaximumCount.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapUserDistributeWindowStart(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowStart.Value.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.UserDistributeWindowStart.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.UserDistributeWindowStart.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapUserDistributeWindowLength(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowLength.Value.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.UserDistributeWindowLength.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.UserDistributeWindowLength.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapOracleDaoConsensusThreshold(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.OracleDaoConsensusThreshold.Value.Formatted() + 0.15)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.OracleDaoConsensusThreshold.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.OracleDaoConsensusThreshold.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSubmitBalancesEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitBalancesEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitBalancesEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.IsSubmitBalancesEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSubmitBalancesFrequency(t *testing.T) {
	newVal := core.Uint256Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Network.SubmitBalancesFrequency.Value.Formatted() + 100)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.SubmitBalancesFrequency.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.SubmitBalancesFrequency.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSubmitPricesEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitPricesEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitPricesEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.IsSubmitPricesEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSubmitPricesFrequency(t *testing.T) {
	newVal := core.Uint256Parameter[uint64]{}
	newVal.Set(tests.PDaoDefaults.Network.SubmitPricesFrequency.Value.Formatted() + 100)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.SubmitPricesFrequency.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.SubmitPricesFrequency.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMinimumNodeFee(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.MinimumNodeFee.Value.Formatted() - 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.MinimumNodeFee.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.MinimumNodeFee.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapTargetNodeFee(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.TargetNodeFee.Value.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.TargetNodeFee.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.TargetNodeFee.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMaximumNodeFee(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.MaximumNodeFee.Value.Formatted() + 0.3)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.MaximumNodeFee.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.MaximumNodeFee.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapNodeFeeDemandRange(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Network.NodeFeeDemandRange.Value, eth.EthToWei(100))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.NodeFeeDemandRange.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.NodeFeeDemandRange.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapTargetRethCollateralRate(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.TargetRethCollateralRate.Value.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.TargetRethCollateralRate.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.TargetRethCollateralRate.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapNodePenaltyThreshold(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.NodePenaltyThreshold.Value.Formatted() + 0.15)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.NodePenaltyThreshold.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.NodePenaltyThreshold.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapPerPenaltyRate(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Network.PerPenaltyRate.Value.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.PerPenaltyRate.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.PerPenaltyRate.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapRethDepositDelay(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Network.RethDepositDelay.Value.Formatted() + time.Hour)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.RethDepositDelay.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.RethDepositDelay.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSubmitRewardsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitRewardsEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitRewardsEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.IsSubmitRewardsEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapNodeRegistrationEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsRegistrationEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsRegistrationEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.IsRegistrationEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSmoothingPoolRegistrationEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsSmoothingPoolRegistrationEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.IsSmoothingPoolRegistrationEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapNodeDepositEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsDepositingEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsDepositingEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.IsDepositingEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapVacantMinipoolsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.AreVacantMinipoolsEnabled.Value
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.AreVacantMinipoolsEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.AreVacantMinipoolsEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMinimumPerMinipoolStake(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Node.MinimumPerMinipoolStake.Value.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.MinimumPerMinipoolStake.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.MinimumPerMinipoolStake.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMaximumPerMinipoolStake(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.PDaoDefaults.Node.MaximumPerMinipoolStake.Value.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Node.MaximumPerMinipoolStake.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.MaximumPerMinipoolStake.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapRewardsIntervalTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.PDaoDefaults.Rewards.IntervalTime.Value.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettingsDetails) {
		newSettings.Rewards.IntervalTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Rewards.IntervalTime.Bootstrap(newVal, opts)
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
	pdaoMgr, err := protocol.NewProtocolDaoManager(mgr.RocketPool)
	if err != nil {
		err = fmt.Errorf("error creating protocol DAO manager: %w", err)
		return
	}
	newPdaoSettings := *pdaoMgr.Settings.ProtocolDaoSettingsDetails
	newPdaoSettings.Auction.IsCreateLotEnabled.Value = !tests.PDaoDefaults.Auction.IsCreateLotEnabled.Value
	newPdaoSettings.Auction.IsBidOnLotEnabled.Value = !tests.PDaoDefaults.Auction.IsBidOnLotEnabled.Value
	newPdaoSettings.Auction.LotMinimumEthValue.Value = big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMinimumEthValue.Value, eth.EthToWei(1))
	newPdaoSettings.Auction.LotMaximumEthValue.Value = big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMaximumEthValue.Value, eth.EthToWei(1))
	newPdaoSettings.Auction.LotDuration.Value.Set(tests.PDaoDefaults.Auction.LotDuration.Value.Formatted() + 1)
	newPdaoSettings.Auction.LotStartingPriceRatio.Value.Set(tests.PDaoDefaults.Auction.LotReservePriceRatio.Value.Formatted() - 0.2)
	newPdaoSettings.Auction.LotReservePriceRatio.Value.Set(tests.PDaoDefaults.Auction.LotReservePriceRatio.Value.Formatted() - 0.1)
	newPdaoSettings.Deposit.IsDepositingEnabled.Value = !tests.PDaoDefaults.Deposit.IsDepositingEnabled.Value
	newPdaoSettings.Deposit.AreDepositAssignmentsEnabled.Value = !tests.PDaoDefaults.Deposit.AreDepositAssignmentsEnabled.Value
	newPdaoSettings.Deposit.MinimumDeposit.Value = big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MinimumDeposit.Value, eth.EthToWei(0.01))
	newPdaoSettings.Deposit.MaximumDepositPoolSize.Value = big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MaximumDepositPoolSize.Value, eth.EthToWei(100))
	newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit.Value.Set(tests.PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Value.Formatted() + 10)
	newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Value.Set(tests.PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Value.Formatted() + 5)
	newPdaoSettings.Deposit.DepositFee.Value.RawValue = big.NewInt(0).Add(tests.PDaoDefaults.Deposit.DepositFee.Value.RawValue, eth.EthToWei(0.1))
	newPdaoSettings.Inflation.IntervalRate.Value.RawValue = big.NewInt(0).Add(tests.PDaoDefaults.Inflation.IntervalRate.Value.RawValue, eth.EthToWei(1))
	newPdaoSettings.Inflation.StartTime.Value.Set(tests.PDaoDefaults.Inflation.StartTime.Value.Formatted().Add(24 * time.Hour))
	newPdaoSettings.Minipool.IsSubmitWithdrawableEnabled.Value = !tests.PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled.Value
	newPdaoSettings.Minipool.IsBondReductionEnabled.Value = !tests.PDaoDefaults.Minipool.IsBondReductionEnabled.Value
	newPdaoSettings.Minipool.LaunchTimeout.Value.Set(tests.PDaoDefaults.Minipool.LaunchTimeout.Value.Formatted() + (24 * time.Hour))
	newPdaoSettings.Minipool.MaximumCount.Value.Set(tests.PDaoDefaults.Minipool.MaximumCount.Value.Formatted() + 1)
	newPdaoSettings.Minipool.UserDistributeWindowStart.Value.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowStart.Value.Formatted() + (24 * time.Hour))
	newPdaoSettings.Minipool.UserDistributeWindowLength.Value.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowLength.Value.Formatted() + (24 * time.Hour))
	newPdaoSettings.Network.OracleDaoConsensusThreshold.Value.Set(tests.PDaoDefaults.Network.OracleDaoConsensusThreshold.Value.Formatted() + 0.15)
	newPdaoSettings.Network.IsSubmitBalancesEnabled.Value = !tests.PDaoDefaults.Network.IsSubmitBalancesEnabled.Value
	newPdaoSettings.Network.SubmitBalancesFrequency.Value.Set(tests.PDaoDefaults.Network.SubmitBalancesFrequency.Value.Formatted() + 100)
	newPdaoSettings.Network.IsSubmitPricesEnabled.Value = !tests.PDaoDefaults.Network.IsSubmitPricesEnabled.Value
	newPdaoSettings.Network.SubmitPricesFrequency.Value.Set(tests.PDaoDefaults.Network.SubmitPricesFrequency.Value.Formatted() + 100)
	newPdaoSettings.Network.MinimumNodeFee.Value.Set(tests.PDaoDefaults.Network.MinimumNodeFee.Value.Formatted() - 0.1)
	newPdaoSettings.Network.TargetNodeFee.Value.Set(tests.PDaoDefaults.Network.TargetNodeFee.Value.Formatted() + 0.1)
	newPdaoSettings.Network.MaximumNodeFee.Value.Set(tests.PDaoDefaults.Network.MaximumNodeFee.Value.Formatted() + 0.3)
	newPdaoSettings.Network.NodeFeeDemandRange.Value = big.NewInt(0).Add(tests.PDaoDefaults.Network.NodeFeeDemandRange.Value, eth.EthToWei(100))
	newPdaoSettings.Network.TargetRethCollateralRate.Value.Set(tests.PDaoDefaults.Network.TargetRethCollateralRate.Value.Formatted() + 0.1)
	newPdaoSettings.Network.NodePenaltyThreshold.Value.Set(tests.PDaoDefaults.Network.NodePenaltyThreshold.Value.Formatted() + 0.15)
	newPdaoSettings.Network.PerPenaltyRate.Value.Set(tests.PDaoDefaults.Network.PerPenaltyRate.Value.Formatted() + 0.1)
	newPdaoSettings.Network.RethDepositDelay.Value.Set(tests.PDaoDefaults.Network.RethDepositDelay.Value.Formatted() + time.Hour)
	newPdaoSettings.Network.IsSubmitRewardsEnabled.Value = !tests.PDaoDefaults.Network.IsSubmitRewardsEnabled.Value
	newPdaoSettings.Node.IsRegistrationEnabled.Value = !tests.PDaoDefaults.Node.IsRegistrationEnabled.Value
	newPdaoSettings.Node.IsSmoothingPoolRegistrationEnabled.Value = !tests.PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled.Value
	newPdaoSettings.Node.IsDepositingEnabled.Value = !tests.PDaoDefaults.Node.IsDepositingEnabled.Value
	newPdaoSettings.Node.AreVacantMinipoolsEnabled.Value = !tests.PDaoDefaults.Node.AreVacantMinipoolsEnabled.Value
	newPdaoSettings.Node.MinimumPerMinipoolStake.Value.Set(tests.PDaoDefaults.Node.MinimumPerMinipoolStake.Value.Formatted() + 0.1)
	newPdaoSettings.Node.MaximumPerMinipoolStake.Value.Set(tests.PDaoDefaults.Node.MaximumPerMinipoolStake.Value.Formatted() + 0.1)
	newPdaoSettings.Rewards.IntervalTime.Value.Set(tests.PDaoDefaults.Rewards.IntervalTime.Value.Formatted() + (24 * time.Hour))

	// Ensure they're all different from the default
	settings_test.EnsureDifferentDetails(t.Fatalf, &tests.PDaoDefaults, &newPdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	txInfos := []*core.TransactionInfo{}
	bootstrappers := []func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.IsCreateLotEnabled.Bootstrap(newPdaoSettings.Auction.IsCreateLotEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.IsBidOnLotEnabled.Bootstrap(newPdaoSettings.Auction.IsBidOnLotEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotMinimumEthValue.Bootstrap(newPdaoSettings.Auction.LotMinimumEthValue.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotMaximumEthValue.Bootstrap(newPdaoSettings.Auction.LotMaximumEthValue.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotDuration.Bootstrap(newPdaoSettings.Auction.LotDuration.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotStartingPriceRatio.Bootstrap(newPdaoSettings.Auction.LotStartingPriceRatio.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotReservePriceRatio.Bootstrap(newPdaoSettings.Auction.LotReservePriceRatio.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.IsDepositingEnabled.Bootstrap(newPdaoSettings.Deposit.IsDepositingEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.AreDepositAssignmentsEnabled.Bootstrap(newPdaoSettings.Deposit.AreDepositAssignmentsEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.MinimumDeposit.Bootstrap(newPdaoSettings.Deposit.MinimumDeposit.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.MaximumDepositPoolSize.Bootstrap(newPdaoSettings.Deposit.MaximumDepositPoolSize.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.MaximumAssignmentsPerDeposit.Bootstrap(newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Bootstrap(newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.DepositFee.Bootstrap(newPdaoSettings.Deposit.DepositFee.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Inflation.IntervalRate.Bootstrap(newPdaoSettings.Inflation.IntervalRate.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Inflation.StartTime.Bootstrap(newPdaoSettings.Inflation.StartTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.IsSubmitWithdrawableEnabled.Bootstrap(newPdaoSettings.Minipool.IsSubmitWithdrawableEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.IsBondReductionEnabled.Bootstrap(newPdaoSettings.Minipool.IsBondReductionEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.LaunchTimeout.Bootstrap(newPdaoSettings.Minipool.LaunchTimeout.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.MaximumCount.Bootstrap(newPdaoSettings.Minipool.MaximumCount.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.UserDistributeWindowStart.Bootstrap(newPdaoSettings.Minipool.UserDistributeWindowStart.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.UserDistributeWindowLength.Bootstrap(newPdaoSettings.Minipool.UserDistributeWindowLength.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.OracleDaoConsensusThreshold.Bootstrap(newPdaoSettings.Network.OracleDaoConsensusThreshold.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.IsSubmitBalancesEnabled.Bootstrap(newPdaoSettings.Network.IsSubmitBalancesEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.SubmitPricesFrequency.Bootstrap(newPdaoSettings.Network.SubmitPricesFrequency.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.IsSubmitPricesEnabled.Bootstrap(newPdaoSettings.Network.IsSubmitPricesEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.SubmitPricesFrequency.Bootstrap(newPdaoSettings.Network.SubmitPricesFrequency.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.MinimumNodeFee.Bootstrap(newPdaoSettings.Network.MinimumNodeFee.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.TargetNodeFee.Bootstrap(newPdaoSettings.Network.TargetNodeFee.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.MaximumNodeFee.Bootstrap(newPdaoSettings.Network.MaximumNodeFee.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.NodeFeeDemandRange.Bootstrap(newPdaoSettings.Network.NodeFeeDemandRange.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.TargetRethCollateralRate.Bootstrap(newPdaoSettings.Network.TargetRethCollateralRate.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.NodePenaltyThreshold.Bootstrap(newPdaoSettings.Network.NodePenaltyThreshold.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.PerPenaltyRate.Bootstrap(newPdaoSettings.Network.PerPenaltyRate.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.RethDepositDelay.Bootstrap(newPdaoSettings.Network.RethDepositDelay.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.IsSubmitRewardsEnabled.Bootstrap(newPdaoSettings.Network.IsSubmitRewardsEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.IsRegistrationEnabled.Bootstrap(newPdaoSettings.Node.IsRegistrationEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.IsSmoothingPoolRegistrationEnabled.Bootstrap(newPdaoSettings.Node.IsSmoothingPoolRegistrationEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.IsDepositingEnabled.Bootstrap(newPdaoSettings.Node.IsDepositingEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.AreVacantMinipoolsEnabled.Bootstrap(newPdaoSettings.Node.AreVacantMinipoolsEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.MinimumPerMinipoolStake.Bootstrap(newPdaoSettings.Node.MinimumPerMinipoolStake.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.MaximumPerMinipoolStake.Bootstrap(newPdaoSettings.Node.MaximumPerMinipoolStake.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Rewards.IntervalTime.Bootstrap(newPdaoSettings.Rewards.IntervalTime.Value, opts)
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
	settings_test.EnsureSameDetails(t.Fatalf, &newPdaoSettings, pdaoMgr.Settings.ProtocolDaoSettingsDetails)
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
	pdaoMgr, err := protocol.NewProtocolDaoManager(mgr.RocketPool)
	if err != nil {
		err = fmt.Errorf("error creating protocol DAO manager: %w", err)
		return
	}
	settings := *pdaoMgr.Settings.ProtocolDaoSettingsDetails
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
	settings_test.EnsureSameDetails(t.Fatalf, &settings, pdaoMgr.Settings.ProtocolDaoSettingsDetails)
	t.Log("New settings match expected settings")
}
