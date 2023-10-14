package bootstrap_test

import (
	"fmt"
	"math/big"
	"runtime/debug"
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
	newVal := !tests.PDaoDefaults.Auction.IsCreateLotEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Auction.IsCreateLotEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.IsCreateLotEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBidOnAuctionLotEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Auction.IsBidOnLotEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Auction.IsBidOnLotEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.IsBidOnLotEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAuctionLotMinimumEthValue(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMinimumEthValue.Get(), eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Auction.LotMinimumEthValue.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotMinimumEthValue.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAuctionLotMaximumEthValue(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMaximumEthValue.Get(), eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Auction.LotMaximumEthValue.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotMaximumEthValue.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAuctionLotDuration(t *testing.T) {
	newVal := tests.PDaoDefaults.Auction.LotDuration.Formatted() + 1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Auction.LotDuration.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotDuration.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapAuctionLotStartingPriceRatio(t *testing.T) {
	newVal := tests.PDaoDefaults.Auction.LotStartingPriceRatio.Formatted() - 0.2
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Auction.LotStartingPriceRatio.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotStartingPriceRatio.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapAuctionLotReservePriceRatio(t *testing.T) {
	newVal := tests.PDaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Auction.LotReservePriceRatio.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Auction.LotReservePriceRatio.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapPoolDepositEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Deposit.IsDepositingEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Deposit.IsDepositingEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.IsDepositingEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapAssignPoolDepositsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Deposit.AreDepositAssignmentsEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Deposit.AreDepositAssignmentsEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.AreDepositAssignmentsEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMinimumPoolDeposit(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MinimumDeposit.Get(), eth.EthToWei(0.01))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Deposit.MinimumDeposit.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.MinimumDeposit.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMaximumDepositPoolSize(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MaximumDepositPoolSize.Get(), eth.EthToWei(100))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Deposit.MaximumDepositPoolSize.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.MaximumDepositPoolSize.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMaximumPoolDepositAssignments(t *testing.T) {
	newVal := tests.PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Formatted() + 10
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Deposit.MaximumAssignmentsPerDeposit.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.MaximumAssignmentsPerDeposit.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapMaximumSocialisedPoolDepositAssignments(t *testing.T) {
	newVal := tests.PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Formatted() + 5
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapDepositFee(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Deposit.DepositFee.Raw(), eth.EthToWei(0.1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Deposit.DepositFee.SetRawValue(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Deposit.DepositFee.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapInflationIntervalRate(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Inflation.IntervalRate.Raw(), eth.EthToWei(1))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Inflation.IntervalRate.SetRawValue(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Inflation.IntervalRate.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapInflationIntervalStartTime(t *testing.T) {
	newVal := tests.PDaoDefaults.Inflation.StartTime.Formatted().Add(24 * time.Hour)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Inflation.StartTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Inflation.StartTime.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapSubmitWithdrawableEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Minipool.IsSubmitWithdrawableEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.IsSubmitWithdrawableEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBondReductionEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Minipool.IsBondReductionEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Minipool.IsBondReductionEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.IsBondReductionEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMinipoolLaunchTimeout(t *testing.T) {
	newVal := tests.PDaoDefaults.Minipool.LaunchTimeout.Formatted() + (24 * time.Hour)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Minipool.LaunchTimeout.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.LaunchTimeout.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapMaximumMinipoolCount(t *testing.T) {
	newVal := tests.PDaoDefaults.Minipool.MaximumCount.Formatted() + 1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Minipool.MaximumCount.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.MaximumCount.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapUserDistributeWindowStart(t *testing.T) {
	newVal := tests.PDaoDefaults.Minipool.UserDistributeWindowStart.Formatted() + (24 * time.Hour)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Minipool.UserDistributeWindowStart.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.UserDistributeWindowStart.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapUserDistributeWindowLength(t *testing.T) {
	newVal := tests.PDaoDefaults.Minipool.UserDistributeWindowLength.Formatted() + (24 * time.Hour)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Minipool.UserDistributeWindowLength.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Minipool.UserDistributeWindowLength.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapOracleDaoConsensusThreshold(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.OracleDaoConsensusThreshold.Formatted() + 0.15
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.OracleDaoConsensusThreshold.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.OracleDaoConsensusThreshold.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapSubmitBalancesEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitBalancesEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.IsSubmitBalancesEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.IsSubmitBalancesEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSubmitBalancesFrequency(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.SubmitBalancesFrequency.Formatted() + 100
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.SubmitBalancesFrequency.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.SubmitBalancesFrequency.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapSubmitPricesEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitPricesEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.IsSubmitPricesEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.IsSubmitPricesEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSubmitPricesFrequency(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.SubmitPricesFrequency.Formatted() + 100
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.SubmitPricesFrequency.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.SubmitPricesFrequency.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapMinimumNodeFee(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.MinimumNodeFee.Formatted() - 0.1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.MinimumNodeFee.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.MinimumNodeFee.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapTargetNodeFee(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.TargetNodeFee.Formatted() + 0.1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.TargetNodeFee.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.TargetNodeFee.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapMaximumNodeFee(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.MaximumNodeFee.Formatted() + 0.3
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.MaximumNodeFee.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.MaximumNodeFee.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapNodeFeeDemandRange(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.PDaoDefaults.Network.NodeFeeDemandRange.Get(), eth.EthToWei(100))
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.NodeFeeDemandRange.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.NodeFeeDemandRange.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapTargetRethCollateralRate(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.TargetRethCollateralRate.Formatted() + 0.1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.TargetRethCollateralRate.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.TargetRethCollateralRate.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapNodePenaltyThreshold(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.NodePenaltyThreshold.Formatted() + 0.15
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.NodePenaltyThreshold.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.NodePenaltyThreshold.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapPerPenaltyRate(t *testing.T) {
	newVal := tests.PDaoDefaults.Network.PerPenaltyRate.Formatted() + 0.1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.PerPenaltyRate.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.PerPenaltyRate.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapSubmitRewardsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Network.IsSubmitRewardsEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Network.IsSubmitRewardsEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Network.IsSubmitRewardsEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapNodeRegistrationEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsRegistrationEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Node.IsRegistrationEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.IsRegistrationEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapSmoothingPoolRegistrationEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Node.IsSmoothingPoolRegistrationEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.IsSmoothingPoolRegistrationEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapNodeDepositEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.IsDepositingEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Node.IsDepositingEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.IsDepositingEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapVacantMinipoolsEnabled(t *testing.T) {
	newVal := !tests.PDaoDefaults.Node.AreVacantMinipoolsEnabled.Get()
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Node.AreVacantMinipoolsEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.AreVacantMinipoolsEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapMinimumPerMinipoolStake(t *testing.T) {
	newVal := tests.PDaoDefaults.Node.MinimumPerMinipoolStake.Formatted() + 0.1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Node.MinimumPerMinipoolStake.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.MinimumPerMinipoolStake.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapMaximumPerMinipoolStake(t *testing.T) {
	newVal := tests.PDaoDefaults.Node.MaximumPerMinipoolStake.Formatted() + 0.1
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Node.MaximumPerMinipoolStake.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Node.MaximumPerMinipoolStake.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapRewardsIntervalTime(t *testing.T) {
	newVal := tests.PDaoDefaults.Rewards.IntervalTime.Formatted() + (24 * time.Hour)
	testPdaoParameterBootstrap(t, func(newSettings *protocol.ProtocolDaoSettings) {
		newSettings.Rewards.IntervalTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return pdaoMgr.Settings.Rewards.IntervalTime.Bootstrap(core.GetValueForUint256(newVal), opts)
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

	// Create new settings
	pdaoMgr, err := protocol.NewProtocolDaoManager(mgr.RocketPool)
	if err != nil {
		err = fmt.Errorf("error creating protocol DAO manager: %w", err)
		return
	}
	newPdaoSettings := *pdaoMgr.Settings
	newPdaoSettings.Auction.IsCreateLotEnabled.Set(!tests.PDaoDefaults.Auction.IsCreateLotEnabled.Get())
	newPdaoSettings.Auction.IsBidOnLotEnabled.Set(!tests.PDaoDefaults.Auction.IsBidOnLotEnabled.Get())
	newPdaoSettings.Auction.LotMinimumEthValue.Set(big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMinimumEthValue.Get(), eth.EthToWei(1)))
	newPdaoSettings.Auction.LotMaximumEthValue.Set(big.NewInt(0).Add(tests.PDaoDefaults.Auction.LotMaximumEthValue.Get(), eth.EthToWei(1)))
	newPdaoSettings.Auction.LotDuration.Set(tests.PDaoDefaults.Auction.LotDuration.Formatted() + 1)
	newPdaoSettings.Auction.LotStartingPriceRatio.Set(tests.PDaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.2)
	newPdaoSettings.Auction.LotReservePriceRatio.Set(tests.PDaoDefaults.Auction.LotReservePriceRatio.Formatted() - 0.1)
	newPdaoSettings.Deposit.IsDepositingEnabled.Set(!tests.PDaoDefaults.Deposit.IsDepositingEnabled.Get())
	newPdaoSettings.Deposit.AreDepositAssignmentsEnabled.Set(!tests.PDaoDefaults.Deposit.AreDepositAssignmentsEnabled.Get())
	newPdaoSettings.Deposit.MinimumDeposit.Set(big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MinimumDeposit.Get(), eth.EthToWei(0.01)))
	newPdaoSettings.Deposit.MaximumDepositPoolSize.Set(big.NewInt(0).Add(tests.PDaoDefaults.Deposit.MaximumDepositPoolSize.Get(), eth.EthToWei(100)))
	newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit.Set(tests.PDaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Formatted() + 10)
	newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Set(tests.PDaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Formatted() + 5)
	newPdaoSettings.Deposit.DepositFee.SetRawValue(big.NewInt(0).Add(tests.PDaoDefaults.Deposit.DepositFee.Raw(), eth.EthToWei(0.1)))
	newPdaoSettings.Inflation.IntervalRate.SetRawValue(big.NewInt(0).Add(tests.PDaoDefaults.Inflation.IntervalRate.Raw(), eth.EthToWei(1)))
	newPdaoSettings.Inflation.StartTime.Set(tests.PDaoDefaults.Inflation.StartTime.Formatted().Add(24 * time.Hour))
	newPdaoSettings.Minipool.IsSubmitWithdrawableEnabled.Set(!tests.PDaoDefaults.Minipool.IsSubmitWithdrawableEnabled.Get())
	newPdaoSettings.Minipool.IsBondReductionEnabled.Set(!tests.PDaoDefaults.Minipool.IsBondReductionEnabled.Get())
	newPdaoSettings.Minipool.LaunchTimeout.Set(tests.PDaoDefaults.Minipool.LaunchTimeout.Formatted() + (24 * time.Hour))
	newPdaoSettings.Minipool.MaximumCount.Set(tests.PDaoDefaults.Minipool.MaximumCount.Formatted() + 1)
	newPdaoSettings.Minipool.UserDistributeWindowStart.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowStart.Formatted() + (24 * time.Hour))
	newPdaoSettings.Minipool.UserDistributeWindowLength.Set(tests.PDaoDefaults.Minipool.UserDistributeWindowLength.Formatted() + (24 * time.Hour))
	newPdaoSettings.Network.OracleDaoConsensusThreshold.Set(tests.PDaoDefaults.Network.OracleDaoConsensusThreshold.Formatted() + 0.15)
	newPdaoSettings.Network.IsSubmitBalancesEnabled.Set(!tests.PDaoDefaults.Network.IsSubmitBalancesEnabled.Get())
	newPdaoSettings.Network.SubmitBalancesFrequency.Set(tests.PDaoDefaults.Network.SubmitBalancesFrequency.Formatted() + 100)
	newPdaoSettings.Network.IsSubmitPricesEnabled.Set(!tests.PDaoDefaults.Network.IsSubmitPricesEnabled.Get())
	newPdaoSettings.Network.SubmitPricesFrequency.Set(tests.PDaoDefaults.Network.SubmitPricesFrequency.Formatted() + 100)
	newPdaoSettings.Network.MinimumNodeFee.Set(tests.PDaoDefaults.Network.MinimumNodeFee.Formatted() - 0.1)
	newPdaoSettings.Network.TargetNodeFee.Set(tests.PDaoDefaults.Network.TargetNodeFee.Formatted() + 0.1)
	newPdaoSettings.Network.MaximumNodeFee.Set(tests.PDaoDefaults.Network.MaximumNodeFee.Formatted() + 0.3)
	newPdaoSettings.Network.NodeFeeDemandRange.Set(big.NewInt(0).Add(tests.PDaoDefaults.Network.NodeFeeDemandRange.Get(), eth.EthToWei(100)))
	newPdaoSettings.Network.TargetRethCollateralRate.Set(tests.PDaoDefaults.Network.TargetRethCollateralRate.Formatted() + 0.1)
	newPdaoSettings.Network.NodePenaltyThreshold.Set(tests.PDaoDefaults.Network.NodePenaltyThreshold.Formatted() + 0.15)
	newPdaoSettings.Network.PerPenaltyRate.Set(tests.PDaoDefaults.Network.PerPenaltyRate.Formatted() + 0.1)
	newPdaoSettings.Network.IsSubmitRewardsEnabled.Set(!tests.PDaoDefaults.Network.IsSubmitRewardsEnabled.Get())
	newPdaoSettings.Node.IsRegistrationEnabled.Set(!tests.PDaoDefaults.Node.IsRegistrationEnabled.Get())
	newPdaoSettings.Node.IsSmoothingPoolRegistrationEnabled.Set(!tests.PDaoDefaults.Node.IsSmoothingPoolRegistrationEnabled.Get())
	newPdaoSettings.Node.IsDepositingEnabled.Set(!tests.PDaoDefaults.Node.IsDepositingEnabled.Get())
	newPdaoSettings.Node.AreVacantMinipoolsEnabled.Set(!tests.PDaoDefaults.Node.AreVacantMinipoolsEnabled.Get())
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
			return pdaoMgr.Settings.Auction.IsCreateLotEnabled.Bootstrap(newPdaoSettings.Auction.IsCreateLotEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.IsBidOnLotEnabled.Bootstrap(newPdaoSettings.Auction.IsBidOnLotEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotMinimumEthValue.Bootstrap(newPdaoSettings.Auction.LotMinimumEthValue.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotMaximumEthValue.Bootstrap(newPdaoSettings.Auction.LotMaximumEthValue.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotDuration.Bootstrap(newPdaoSettings.Auction.LotDuration.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotStartingPriceRatio.Bootstrap(newPdaoSettings.Auction.LotStartingPriceRatio.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Auction.LotReservePriceRatio.Bootstrap(newPdaoSettings.Auction.LotReservePriceRatio.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.IsDepositingEnabled.Bootstrap(newPdaoSettings.Deposit.IsDepositingEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.AreDepositAssignmentsEnabled.Bootstrap(newPdaoSettings.Deposit.AreDepositAssignmentsEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.MinimumDeposit.Bootstrap(newPdaoSettings.Deposit.MinimumDeposit.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.MaximumDepositPoolSize.Bootstrap(newPdaoSettings.Deposit.MaximumDepositPoolSize.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.MaximumAssignmentsPerDeposit.Bootstrap(newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Bootstrap(newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Deposit.DepositFee.Bootstrap(newPdaoSettings.Deposit.DepositFee.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Inflation.IntervalRate.Bootstrap(newPdaoSettings.Inflation.IntervalRate.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Inflation.StartTime.Bootstrap(newPdaoSettings.Inflation.StartTime.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.IsSubmitWithdrawableEnabled.Bootstrap(newPdaoSettings.Minipool.IsSubmitWithdrawableEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.IsBondReductionEnabled.Bootstrap(newPdaoSettings.Minipool.IsBondReductionEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.LaunchTimeout.Bootstrap(newPdaoSettings.Minipool.LaunchTimeout.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.MaximumCount.Bootstrap(newPdaoSettings.Minipool.MaximumCount.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.UserDistributeWindowStart.Bootstrap(newPdaoSettings.Minipool.UserDistributeWindowStart.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Minipool.UserDistributeWindowLength.Bootstrap(newPdaoSettings.Minipool.UserDistributeWindowLength.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.OracleDaoConsensusThreshold.Bootstrap(newPdaoSettings.Network.OracleDaoConsensusThreshold.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.IsSubmitBalancesEnabled.Bootstrap(newPdaoSettings.Network.IsSubmitBalancesEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.SubmitPricesFrequency.Bootstrap(newPdaoSettings.Network.SubmitPricesFrequency.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.IsSubmitPricesEnabled.Bootstrap(newPdaoSettings.Network.IsSubmitPricesEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.SubmitPricesFrequency.Bootstrap(newPdaoSettings.Network.SubmitPricesFrequency.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.MinimumNodeFee.Bootstrap(newPdaoSettings.Network.MinimumNodeFee.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.TargetNodeFee.Bootstrap(newPdaoSettings.Network.TargetNodeFee.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.MaximumNodeFee.Bootstrap(newPdaoSettings.Network.MaximumNodeFee.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.NodeFeeDemandRange.Bootstrap(newPdaoSettings.Network.NodeFeeDemandRange.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.TargetRethCollateralRate.Bootstrap(newPdaoSettings.Network.TargetRethCollateralRate.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.NodePenaltyThreshold.Bootstrap(newPdaoSettings.Network.NodePenaltyThreshold.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.PerPenaltyRate.Bootstrap(newPdaoSettings.Network.PerPenaltyRate.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Network.IsSubmitRewardsEnabled.Bootstrap(newPdaoSettings.Network.IsSubmitRewardsEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.IsRegistrationEnabled.Bootstrap(newPdaoSettings.Node.IsRegistrationEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.IsSmoothingPoolRegistrationEnabled.Bootstrap(newPdaoSettings.Node.IsSmoothingPoolRegistrationEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.IsDepositingEnabled.Bootstrap(newPdaoSettings.Node.IsDepositingEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.AreVacantMinipoolsEnabled.Bootstrap(newPdaoSettings.Node.AreVacantMinipoolsEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.MinimumPerMinipoolStake.Bootstrap(newPdaoSettings.Node.MinimumPerMinipoolStake.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Node.MaximumPerMinipoolStake.Bootstrap(newPdaoSettings.Node.MaximumPerMinipoolStake.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdaoMgr.Settings.Rewards.IntervalTime.Bootstrap(newPdaoSettings.Rewards.IntervalTime.Raw(), opts)
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
		core.QueryAllFields(pdaoMgr.Settings, mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	settings_test.EnsureSameDetails(t.Fatalf, &newPdaoSettings, pdaoMgr.Settings)
	t.Log("New settings match expected settings")
}

func testPdaoParameterBootstrap(t *testing.T, setter func(*protocol.ProtocolDaoSettings), bootstrapper func() (*core.TransactionInfo, error)) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
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

	// Get the original settings
	pdaoMgr, err := protocol.NewProtocolDaoManager(mgr.RocketPool)
	if err != nil {
		err = fmt.Errorf("error creating protocol DAO manager: %w", err)
		return
	}
	settings := *pdaoMgr.Settings
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
		core.QueryAllFields(pdaoMgr.Settings, mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	settings_test.EnsureSameDetails(t.Fatalf, &settings, pdaoMgr.Settings)
	t.Log("New settings match expected settings")
}
