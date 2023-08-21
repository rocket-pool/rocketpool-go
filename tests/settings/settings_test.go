package settings_test

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/settings"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"golang.org/x/sync/errgroup"
)

func TestBoostrapFunctions(t *testing.T) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Initializers
	odao, err := settings.NewOracleDaoSettings(rp)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating odao settings binding: %w", err))
	}
	pdao, err := settings.NewProtocolDaoSettings(rp)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating pdao settings binding: %w", err))
	}
	err = createDefaults(mgr)
	if err != nil {
		t.Fatal("error creating defaults: %w", err)
	}

	// Get all of the current settings
	err = rp.Query(func(mc *batch.MultiCaller) error {
		odao.GetAllDetails(mc)
		pdao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all initial details: %w", err))
	}

	// Verify details
	EnsureSameDetails(t, &odaoDefaults, &odao.Details)
	EnsureSameDetails(t, &pdaoDefaults, &pdao.Details)
	t.Log("Original values match contract initial values")

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
	newPdaoSettings.Deposit.MaximumAssignmentsPerDeposit.Set(pdaoDefaults.Deposit.MaximumAssignmentsPerDeposit.Formatted() * 2)
	newPdaoSettings.Deposit.MaximumSocialisedAssignmentsPerDeposit.Set(pdaoDefaults.Deposit.MaximumSocialisedAssignmentsPerDeposit.Formatted() * 2)
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
	EnsureDifferentDetails(t, &pdaoDefaults, &newPdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	opts := mgr.OwnerAccount.Transactor
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
			return pdao.BootstrapMinipoolLaunchTimeout(newPdaoSettings.Minipool.LaunchTimeout, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapBondReductionEnabled(newPdaoSettings.Minipool.IsBondReductionEnabled, opts)
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
			return pdao.BootstrapNodePenaltyThreshold(newPdaoSettings.Network.NodePenaltyThreshold, opts)
		},
		func() (*core.TransactionInfo, error) {
			return pdao.BootstrapPerPenaltyRate(newPdaoSettings.Network.PerPenaltyRate, opts)
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
	EnsureSameDetails(t, &newPdaoSettings, &pdao.Details)
	t.Log("New settings match expected settings")
}

// Compares two details structs to ensure their fields all have the same values
func EnsureSameDetails[objType any](t *testing.T, expected *objType, actual *objType) {
	expectedVal := reflect.ValueOf(expected).Elem()
	actualVal := reflect.ValueOf(actual).Elem()
	compareImpl(t, expectedVal, actualVal, expectedVal.Type().Name(), true)
}

// Compares two details structs to ensure their fields all have different values
func EnsureDifferentDetails[objType any](t *testing.T, expected *objType, actual *objType) {
	expectedVal := reflect.ValueOf(expected).Elem()
	actualVal := reflect.ValueOf(actual).Elem()
	compareImpl(t, expectedVal, actualVal, expectedVal.Type().Name(), false)
}

func testPdaoBootstrap[objType bool | *big.Int](t *testing.T, originalSettings *settings.ProtocolDaoSettingsDetails, newValue objType, bootstrapper func() (*core.TransactionInfo, error)) {
	// Get the original setting

	// Set the new setting

}

// Detail comparison implementation
func compareImpl(t *testing.T, expected reflect.Value, actual reflect.Value, header string, checkIfEqual bool) {
	refType := expected.Type()
	fieldCount := refType.NumField()

	for i := 0; i < fieldCount; i++ {
		field := refType.Field(i)
		childExpected := expected.Field(i)
		childActual := actual.Field(i)

		// Try casting to parameters first
		expectedParam, isIParameter := childExpected.Interface().(core.IParameter)
		expectedUint8Param, isIUint8Parameter := childExpected.Interface().(core.IUint8Parameter)

		passedCheck := true
		if isIParameter {
			// Handle parameters
			actualParam := childActual.Interface().(core.IParameter)
			if expectedParam.GetRawValue() == nil {
				t.Errorf("field %s.%s of type %s - expected was nil", header, field.Name, field.Type.Name())
			} else if actualParam.GetRawValue() == nil {
				t.Errorf("field %s.%s of type %s - actual was nil", header, field.Name, field.Type.Name())
			} else {
				if checkIfEqual {
					passedCheck = expectedParam.GetRawValue().Cmp(actualParam.GetRawValue()) == 0
				} else {
					passedCheck = expectedParam.GetRawValue().Cmp(actualParam.GetRawValue()) != 0
				}
			}
		} else if isIUint8Parameter {
			// Handle uint8 parameters
			actualUint8Param := childActual.Interface().(core.IUint8Parameter)
			if checkIfEqual {
				passedCheck = expectedUint8Param.GetRawValue() == actualUint8Param.GetRawValue()
			} else {
				passedCheck = expectedUint8Param.GetRawValue() != actualUint8Param.GetRawValue()
			}
		} else if field.Type.Kind() == reflect.Struct {
			// Handle other nested structs
			compareImpl(t, childExpected, childActual, fmt.Sprintf("%s.%s", header, field.Name), checkIfEqual)
			continue
		} else {
			// Handle primitives
			switch expectedVal := childExpected.Interface().(type) {
			case *big.Int:
				actualVal := childActual.Interface().(*big.Int)
				if expectedVal == nil {
					t.Errorf("field %s.%s (big.Int) - expected was nil", header, field.Name)
				} else if actualVal == nil {
					t.Errorf("field %s.%s (big.Int) - actual was nil", header, field.Name)
				} else {
					if checkIfEqual {
						passedCheck = expectedVal.Cmp(actualVal) == 0
					} else {
						passedCheck = expectedVal.Cmp(actualVal) != 0
					}
				}
			case bool:
				if checkIfEqual {
					passedCheck = expectedVal == childActual.Interface().(bool)
				} else {
					passedCheck = expectedVal != childActual.Interface().(bool)
				}
			default:
				t.Fatalf("unexpected type %s in field %s.%s", field.Type.Name(), header, field.Name)
			}
		}

		if !passedCheck {
			if checkIfEqual {
				t.Errorf("%s.%s differed; expected %v but got %v", header, field.Name, childExpected.Interface(), childActual.Interface())
			} else {
				t.Errorf("%s.%s was the same; expected not %v but got %v", header, field.Name, childExpected.Interface(), childActual.Interface())
			}
		}
	}
}
