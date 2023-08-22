package settings_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/settings"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"golang.org/x/sync/errgroup"
)

func Test_BootstrapCreateAuctionLotEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Auction.IsCreateLotEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Auction.IsCreateLotEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapCreateAuctionLotEnabled(newVal, opts)
	})
}

func Test_BootstrapBidOnAuctionLotEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Auction.IsBidOnLotEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Auction.IsBidOnLotEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapBidOnAuctionLotEnabled(newVal, opts)
	})
}

func Test_BootstrapAuctionLotMinimumEthValue(t *testing.T) {
	newVal := big.NewInt(0).Add(pdaoDefaults.Auction.LotMinimumEthValue, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotMinimumEthValue = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapAuctionLotMinimumEthValue(newVal, opts)
	})
}

func Test_BootstrapAuctionLotMaximumEthValue(t *testing.T) {
	newVal := big.NewInt(0).Add(pdaoDefaults.Auction.LotMaximumEthValue, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotMaximumEthValue = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapAuctionLotMaximumEthValue(newVal, opts)
	})
}

func Test_BootstrapAuctionLotDuration(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(pdaoDefaults.Auction.LotDuration.Formatted() + 1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotDuration.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapAuctionLotDuration(newVal, opts)
	})
}

func Test_BootstrapAuctionLotStartingPriceRatio(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Auction.LotStartingPriceRatio.Formatted() - 0.2)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotStartingPriceRatio.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapAuctionLotStartingPriceRatio(newVal, opts)
	})
}

func Test_BootstrapAuctionLotReservePriceRatio(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Auction.LotReservePriceRatio.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapAuctionLotReservePriceRatio(newVal, opts)
	})
}

func Test_BootstrapPoolDepositEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Deposit.IsDepositingEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.IsDepositingEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapPoolDepositEnabled(newVal, opts)
	})
}

func Test_BootstrapAssignPoolDepositsEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Deposit.AreDepositAssignmentsEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.AreDepositAssignmentsEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapAssignPoolDepositsEnabled(newVal, opts)
	})
}

func Test_BootstrapMinimumPoolDeposit(t *testing.T) {
	newVal := big.NewInt(0).Add(pdaoDefaults.Deposit.MinimumDeposit, eth.EthToWei(0.01))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MinimumDeposit = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMinimumPoolDeposit(newVal, opts)
	})
}

func Test_BootstrapMaximumDepositPoolSize(t *testing.T) {
	newVal := big.NewInt(0).Add(pdaoDefaults.Deposit.MaximumDepositPoolSize, eth.EthToWei(100))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumDepositPoolSize = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMaximumDepositPoolSize(newVal, opts)
	})
}

func Test_BootstrapMaximumPoolDepositAssignments(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(pdaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Formatted() + 10)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumAssignmentsPerDeposit.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMaximumPoolDepositAssignments(newVal, opts)
	})
}

func Test_BootstrapMaximumSocialisedPoolDepositAssignments(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(pdaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Formatted() + 5)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMaximumSocialisedPoolDepositAssignments(newVal, opts)
	})
}

func Test_BootstrapDepositFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.RawValue = big.NewInt(0).Add(pdaoDefaults.Deposit.DepositFee.RawValue, eth.EthToWei(0.1))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Deposit.DepositFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapDepositFee(newVal, opts)
	})
}

func Test_BootstrapInflationIntervalRate(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.RawValue = big.NewInt(0).Add(pdaoDefaults.Inflation.IntervalRate.RawValue, eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Inflation.IntervalRate.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapInflationIntervalRate(newVal, opts)
	})
}

func Test_BootstrapInflationIntervalStartTime(t *testing.T) {
	newVal := core.Parameter[time.Time]{}
	newVal.Set(pdaoDefaults.Inflation.StartTime.Formatted().Add(24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Inflation.StartTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapInflationIntervalStartTime(newVal, opts)
	})
}

func Test_BootstrapSubmitWithdrawableEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Minipool.IsSubmitWithdrawableEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.IsSubmitWithdrawableEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapSubmitWithdrawableEnabled(newVal, opts)
	})
}

func Test_BootstrapBondReductionEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Minipool.IsBondReductionEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.IsBondReductionEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapBondReductionEnabled(newVal, opts)
	})
}

func Test_BootstrapMinipoolLaunchTimeout(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(pdaoDefaults.Minipool.LaunchTimeout.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.LaunchTimeout.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMinipoolLaunchTimeout(newVal, opts)
	})
}

func Test_BootstrapMaximumMinipoolCount(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(pdaoDefaults.Minipool.MaximumCount.Formatted() + 1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.MaximumCount.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMaximumMinipoolCount(newVal, opts)
	})
}

func Test_BootstrapUserDistributeWindowStart(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(pdaoDefaults.Minipool.UserDistributeWindowStart.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.UserDistributeWindowStart.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapUserDistributeWindowStart(newVal, opts)
	})
}

func Test_BootstrapUserDistributeWindowLength(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(pdaoDefaults.Minipool.UserDistributeWindowLength.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Minipool.UserDistributeWindowLength.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapUserDistributeWindowLength(newVal, opts)
	})
}

func Test_BootstrapOracleDaoConsensusThreshold(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Network.OracleDaoConsensusThreshold.Formatted() + 0.15)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.OracleDaoConsensusThreshold.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapOracleDaoConsensusThreshold(newVal, opts)
	})
}

func Test_BootstrapSubmitBalancesEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Network.IsSubmitBalancesEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitBalancesEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapSubmitBalancesEnabled(newVal, opts)
	})
}

func Test_BootstrapSubmitBalancesFrequency(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(pdaoDefaults.Network.SubmitBalancesFrequency.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.SubmitBalancesFrequency.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapSubmitBalancesFrequency(newVal, opts)
	})
}

func Test_BootstrapSubmitPricesEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Network.IsSubmitPricesEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitPricesEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapSubmitPricesEnabled(newVal, opts)
	})
}

func Test_BootstrapSubmitPricesFrequency(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(pdaoDefaults.Network.SubmitPricesFrequency.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.SubmitPricesFrequency.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapSubmitPricesFrequency(newVal, opts)
	})
}

func Test_BootstrapMinimumNodeFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Network.MinimumNodeFee.Formatted() - 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.MinimumNodeFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMinimumNodeFee(newVal, opts)
	})
}

func Test_BootstrapTargetNodeFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Network.TargetNodeFee.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.TargetNodeFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapTargetNodeFee(newVal, opts)
	})
}

func Test_BootstrapMaximumNodeFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Network.MaximumNodeFee.Formatted() + 0.3)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.MaximumNodeFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMaximumNodeFee(newVal, opts)
	})
}

func Test_BootstrapNodeFeeDemandRange(t *testing.T) {
	newVal := big.NewInt(0).Add(pdaoDefaults.Network.NodeFeeDemandRange, eth.EthToWei(100))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.NodeFeeDemandRange = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapNodeFeeDemandRange(newVal, opts)
	})
}

func Test_BootstrapTargetRethCollateralRate(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Network.TargetRethCollateralRate.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.TargetRethCollateralRate.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapTargetRethCollateralRate(newVal, opts)
	})
}

func Test_BootstrapNodePenaltyThreshold(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Network.NodePenaltyThreshold.Formatted() + 0.15)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.NodePenaltyThreshold.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapNodePenaltyThreshold(newVal, opts)
	})
}

func Test_BootstrapPerPenaltyRate(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Network.PerPenaltyRate.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.PerPenaltyRate.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapPerPenaltyRate(newVal, opts)
	})
}

func Test_BootstrapRethDepositDelay(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(pdaoDefaults.Network.RethDepositDelay.Formatted() + time.Hour)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.RethDepositDelay.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapRethDepositDelay(newVal, opts)
	})
}

func Test_BootstrapSubmitRewardsEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Network.IsSubmitRewardsEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Network.IsSubmitRewardsEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapSubmitRewardsEnabled(newVal, opts)
	})
}

func Test_BootstrapNodeRegistrationEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Node.IsRegistrationEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsRegistrationEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapNodeRegistrationEnabled(newVal, opts)
	})
}

func Test_BootstrapSmoothingPoolRegistrationEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Node.IsSmoothingPoolRegistrationEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsSmoothingPoolRegistrationEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapSmoothingPoolRegistrationEnabled(newVal, opts)
	})
}

func Test_BootstrapNodeDepositEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Node.IsDepositingEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Node.IsDepositingEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapNodeDepositEnabled(newVal, opts)
	})
}

func Test_BootstrapVacantMinipoolsEnabled(t *testing.T) {
	newVal := !pdaoDefaults.Node.AreVacantMinipoolsEnabled
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Node.AreVacantMinipoolsEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapVacantMinipoolsEnabled(newVal, opts)
	})
}

func Test_BootstrapMinimumPerMinipoolStake(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Node.MinimumPerMinipoolStake.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Node.MinimumPerMinipoolStake.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMinimumPerMinipoolStake(newVal, opts)
	})
}

func Test_BootstrapMaximumPerMinipoolStake(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(pdaoDefaults.Node.MaximumPerMinipoolStake.Formatted() + 0.1)
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Node.MaximumPerMinipoolStake.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapMaximumPerMinipoolStake(newVal, opts)
	})
}

func Test_BootstrapRewardsIntervalTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(pdaoDefaults.Rewards.IntervalTime.Formatted() + (24 * time.Hour))
	testPdaoParameterBootstrap(t, func(newSettings *settings.ProtocolDaoSettingsDetails) {
		newSettings.Rewards.IntervalTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return pdao.BootstrapRewardsIntervalTime(newVal, opts)
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
	newPdaoSettings := settings.ProtocolDaoSettingsDetails{}
	newPdaoSettings.Auction.IsCreateLotEnabled = !pdaoDefaults.Auction.IsCreateLotEnabled
	newPdaoSettings.Auction.IsBidOnLotEnabled = !pdaoDefaults.Auction.IsBidOnLotEnabled
	newPdaoSettings.Auction.LotMinimumEthValue = big.NewInt(0).Add(pdaoDefaults.Auction.LotMinimumEthValue, eth.EthToWei(1))
	newPdaoSettings.Auction.LotMaximumEthValue = big.NewInt(0).Add(pdaoDefaults.Auction.LotMaximumEthValue, eth.EthToWei(1))
	newPdaoSettings.Auction.LotDuration.Set(pdaoDefaults.Auction.LotDuration.Formatted() + 1)
	newPdaoSettings.Auction.LotStartingPriceRatio.Set(pdaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.2)
	newPdaoSettings.Auction.LotReservePriceRatio.Set(pdaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.1)
	newPdaoSettings.Deposit.IsDepositingEnabled = !pdaoDefaults.Deposit.IsDepositingEnabled
	newPdaoSettings.Deposit.AreDepositAssignmentsEnabled = !pdaoDefaults.Deposit.AreDepositAssignmentsEnabled
	newPdaoSettings.Deposit.MinimumDeposit = big.NewInt(0).Add(pdaoDefaults.Deposit.MinimumDeposit, eth.EthToWei(0.01))
	newPdaoSettings.Deposit.MaximumDepositPoolSize = big.NewInt(0).Add(pdaoDefaults.Deposit.MaximumDepositPoolSize, eth.EthToWei(100))
	newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit.Set(pdaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Formatted() + 10)
	newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Set(pdaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Formatted() + 5)
	newPdaoSettings.Deposit.DepositFee.RawValue = big.NewInt(0).Add(pdaoDefaults.Deposit.DepositFee.RawValue, eth.EthToWei(0.1))
	newPdaoSettings.Inflation.IntervalRate.RawValue = big.NewInt(0).Add(pdaoDefaults.Inflation.IntervalRate.RawValue, eth.EthToWei(1))
	newPdaoSettings.Inflation.StartTime.Set(pdaoDefaults.Inflation.StartTime.Formatted().Add(24 * time.Hour))
	newPdaoSettings.Minipool.IsSubmitWithdrawableEnabled = !pdaoDefaults.Minipool.IsSubmitWithdrawableEnabled
	newPdaoSettings.Minipool.IsBondReductionEnabled = !pdaoDefaults.Minipool.IsBondReductionEnabled
	newPdaoSettings.Minipool.LaunchTimeout.Set(pdaoDefaults.Minipool.LaunchTimeout.Formatted() + (24 * time.Hour))
	newPdaoSettings.Minipool.MaximumCount.Set(pdaoDefaults.Minipool.MaximumCount.Formatted() + 1)
	newPdaoSettings.Minipool.UserDistributeWindowStart.Set(pdaoDefaults.Minipool.UserDistributeWindowStart.Formatted() + (24 * time.Hour))
	newPdaoSettings.Minipool.UserDistributeWindowLength.Set(pdaoDefaults.Minipool.UserDistributeWindowLength.Formatted() + (24 * time.Hour))
	newPdaoSettings.Network.OracleDaoConsensusThreshold.Set(pdaoDefaults.Network.OracleDaoConsensusThreshold.Formatted() + 0.15)
	newPdaoSettings.Network.IsSubmitBalancesEnabled = !pdaoDefaults.Network.IsSubmitBalancesEnabled
	newPdaoSettings.Network.SubmitBalancesFrequency.Set(pdaoDefaults.Network.SubmitBalancesFrequency.Formatted() + (24 * time.Hour))
	newPdaoSettings.Network.IsSubmitPricesEnabled = !pdaoDefaults.Network.IsSubmitPricesEnabled
	newPdaoSettings.Network.SubmitPricesFrequency.Set(pdaoDefaults.Network.SubmitPricesFrequency.Formatted() + (24 * time.Hour))
	newPdaoSettings.Network.MinimumNodeFee.Set(pdaoDefaults.Network.MinimumNodeFee.Formatted() - 0.1)
	newPdaoSettings.Network.TargetNodeFee.Set(pdaoDefaults.Network.TargetNodeFee.Formatted() + 0.1)
	newPdaoSettings.Network.MaximumNodeFee.Set(pdaoDefaults.Network.MaximumNodeFee.Formatted() + 0.3)
	newPdaoSettings.Network.NodeFeeDemandRange = big.NewInt(0).Add(pdaoDefaults.Network.NodeFeeDemandRange, eth.EthToWei(100))
	newPdaoSettings.Network.TargetRethCollateralRate.Set(pdaoDefaults.Network.TargetRethCollateralRate.Formatted() + 0.1)
	newPdaoSettings.Network.NodePenaltyThreshold.Set(pdaoDefaults.Network.NodePenaltyThreshold.Formatted() + 0.15)
	newPdaoSettings.Network.PerPenaltyRate.Set(pdaoDefaults.Network.PerPenaltyRate.Formatted() + 0.1)
	newPdaoSettings.Network.RethDepositDelay.Set(pdaoDefaults.Network.RethDepositDelay.Formatted() + time.Hour)
	newPdaoSettings.Network.IsSubmitRewardsEnabled = !pdaoDefaults.Network.IsSubmitRewardsEnabled
	newPdaoSettings.Node.IsRegistrationEnabled = !pdaoDefaults.Node.IsRegistrationEnabled
	newPdaoSettings.Node.IsSmoothingPoolRegistrationEnabled = !pdaoDefaults.Node.IsSmoothingPoolRegistrationEnabled
	newPdaoSettings.Node.IsDepositingEnabled = !pdaoDefaults.Node.IsDepositingEnabled
	newPdaoSettings.Node.AreVacantMinipoolsEnabled = !pdaoDefaults.Node.AreVacantMinipoolsEnabled
	newPdaoSettings.Node.MinimumPerMinipoolStake.Set(pdaoDefaults.Node.MinimumPerMinipoolStake.Formatted() + 0.1)
	newPdaoSettings.Node.MaximumPerMinipoolStake.Set(pdaoDefaults.Node.MaximumPerMinipoolStake.Formatted() + 0.1)
	newPdaoSettings.Rewards.IntervalTime.Set(pdaoDefaults.Rewards.IntervalTime.Formatted() + (24 * time.Hour))

	// Ensure they're all different from the default
	EnsureDifferentDetails(t.Fatalf, &pdaoDefaults, &newPdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	txInfos := []*core.TransactionInfo{}
	bootstrappers := []func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapCreateAuctionLotEnabled(newPdaoSettings.Auction.IsCreateLotEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapBidOnAuctionLotEnabled(newPdaoSettings.Auction.IsBidOnLotEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapAuctionLotMinimumEthValue(newPdaoSettings.Auction.LotMinimumEthValue, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapAuctionLotMaximumEthValue(newPdaoSettings.Auction.LotMaximumEthValue, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapAuctionLotDuration(newPdaoSettings.Auction.LotDuration, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapAuctionLotStartingPriceRatio(newPdaoSettings.Auction.LotStartingPriceRatio, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapAuctionLotReservePriceRatio(newPdaoSettings.Auction.LotReservePriceRatio, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapPoolDepositEnabled(newPdaoSettings.Deposit.IsDepositingEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapAssignPoolDepositsEnabled(newPdaoSettings.Deposit.AreDepositAssignmentsEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMinimumPoolDeposit(newPdaoSettings.Deposit.MinimumDeposit, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMaximumDepositPoolSize(newPdaoSettings.Deposit.MaximumDepositPoolSize, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMaximumPoolDepositAssignments(newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMaximumSocialisedPoolDepositAssignments(newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapDepositFee(newPdaoSettings.Deposit.DepositFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapInflationIntervalRate(newPdaoSettings.Inflation.IntervalRate, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapInflationIntervalStartTime(newPdaoSettings.Inflation.StartTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapSubmitWithdrawableEnabled(newPdaoSettings.Minipool.IsSubmitWithdrawableEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapBondReductionEnabled(newPdaoSettings.Minipool.IsBondReductionEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMinipoolLaunchTimeout(newPdaoSettings.Minipool.LaunchTimeout, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMaximumMinipoolCount(newPdaoSettings.Minipool.MaximumCount, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapUserDistributeWindowStart(newPdaoSettings.Minipool.UserDistributeWindowStart, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapUserDistributeWindowLength(newPdaoSettings.Minipool.UserDistributeWindowLength, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapOracleDaoConsensusThreshold(newPdaoSettings.Network.OracleDaoConsensusThreshold, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapSubmitBalancesEnabled(newPdaoSettings.Network.IsSubmitBalancesEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapSubmitBalancesFrequency(newPdaoSettings.Network.SubmitPricesFrequency, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapSubmitPricesEnabled(newPdaoSettings.Network.IsSubmitPricesEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapSubmitPricesFrequency(newPdaoSettings.Network.SubmitPricesFrequency, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMinimumNodeFee(newPdaoSettings.Network.MinimumNodeFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapTargetNodeFee(newPdaoSettings.Network.TargetNodeFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMaximumNodeFee(newPdaoSettings.Network.MaximumNodeFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapNodeFeeDemandRange(newPdaoSettings.Network.NodeFeeDemandRange, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapTargetRethCollateralRate(newPdaoSettings.Network.TargetRethCollateralRate, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapNodePenaltyThreshold(newPdaoSettings.Network.NodePenaltyThreshold, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapPerPenaltyRate(newPdaoSettings.Network.PerPenaltyRate, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapRethDepositDelay(newPdaoSettings.Network.RethDepositDelay, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapSubmitRewardsEnabled(newPdaoSettings.Network.IsSubmitRewardsEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapNodeRegistrationEnabled(newPdaoSettings.Node.IsRegistrationEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapSmoothingPoolRegistrationEnabled(newPdaoSettings.Node.IsSmoothingPoolRegistrationEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapNodeDepositEnabled(newPdaoSettings.Node.IsDepositingEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapVacantMinipoolsEnabled(newPdaoSettings.Node.AreVacantMinipoolsEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMinimumPerMinipoolStake(newPdaoSettings.Node.MinimumPerMinipoolStake, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapMaximumPerMinipoolStake(newPdaoSettings.Node.MaximumPerMinipoolStake, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapRewardsIntervalTime(newPdaoSettings.Rewards.IntervalTime, opts)
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
	txs, err := rp.SubmitTransactions(txInfos, opts)
	if err != nil {
		t.Fatalf("error submitting transactions: %s", err.Error())
	}
	t.Log("Bootstrappers submitted")

	var wg errgroup.Group
	for _, tx := range txs {
		tx := tx
		wg.Go(func() error {
			_, err := utils.WaitForTransaction(rp.Client, tx.Hash())
			return err
		})
	}

	err = wg.Wait()
	if err != nil {
		t.Fatalf("error waiting for transactions: %s", err.Error())
	}
	t.Log("Bootstrappers executed")

	// Get new values and make sure they match
	err = rp.Query(func(mc *batch.MultiCaller) error {
		pdao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	EnsureSameDetails(t.Fatalf, &newPdaoSettings, &pdao.Details)
	t.Log("New settings match expected settings")
}

func testPdaoParameterBootstrap(t *testing.T, setter func(*settings.ProtocolDaoSettingsDetails), bootstrapper func() (*core.TransactionInfo, error)) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Get the original settings
	var settings settings.ProtocolDaoSettingsDetails
	Clone(t, &pdaoDefaults, &settings)
	pass := EnsureSameDetails(t.Errorf, &pdaoDefaults, &settings)
	if !pass {
		t.Fatalf("Details differed unexpectedly, can't continue")
	}
	t.Log("Cloned default settings")

	// Set the new setting
	setter(&settings)
	pass = EnsureSameDetails(t.Logf, &pdaoDefaults, &settings)
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
	tx, err := rp.SubmitTransaction(txInfo, opts)
	if err != nil {
		t.Fatalf("error submitting transaction: %s", err.Error())
	}
	t.Log("Bootstrapper submitted")

	_, err = utils.WaitForTransaction(rp.Client, tx.Hash())
	if err != nil {
		t.Fatalf("error waiting for transactions: %s", err.Error())
	}
	t.Log("Bootstrapper executed")

	// Get new values and make sure they match
	err = rp.Query(func(mc *batch.MultiCaller) error {
		pdao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	EnsureSameDetails(t.Fatalf, &settings, &pdao.Details)
	t.Log("New settings match expected settings")
}
