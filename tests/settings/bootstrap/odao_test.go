package bootstrap_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/tests"
	settings_test "github.com/rocket-pool/rocketpool-go/tests/settings"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"golang.org/x/sync/errgroup"
)

func Test_BootstrapChallengeCooldown(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Members.ChallengeCooldown.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeCooldown.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Members.ChallengeCooldown.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapChallengeCost(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Members.ChallengeCost.Value, eth.EthToWei(1))
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeCost.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Members.ChallengeCost.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapChallengeWindow(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Members.ChallengeWindow.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeWindow.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Members.ChallengeWindow.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Members.Quorum.Value.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.Quorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Members.Quorum.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapRplBond(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Members.RplBond.Value, eth.EthToWei(1000))
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.RplBond.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Members.RplBond.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMax(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMax.Value.Formatted() + 5)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.UnbondedMinipoolMax.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Members.UnbondedMinipoolMax.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMinFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMinFee.Value.Formatted() + 0.1)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.UnbondedMinipoolMinFee.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Members.UnbondedMinipoolMinFee.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBondReductionCancellationQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionCancellationQuorum.Value.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionCancellationQuorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipools.BondReductionCancellationQuorum.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBondReductionWindowLength(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionWindowLength.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionWindowLength.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipools.BondReductionWindowLength.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBondReductionWindowStart(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionWindowStart.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionWindowStart.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipools.BondReductionWindowStart.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapScrubPenaltyEnabled(t *testing.T) {
	newVal := !tests.ODaoDefaults.Minipools.IsScrubPenaltyEnabled.Value
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.IsScrubPenaltyEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipools.IsScrubPenaltyEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapPromotionScrubPeriod(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.PromotionScrubPeriod.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.PromotionScrubPeriod.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipools.PromotionScrubPeriod.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapScrubPeriod(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.ScrubPeriod.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.ScrubPeriod.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipools.ScrubPeriod.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapScrubQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipools.ScrubQuorum.Value.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.ScrubQuorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipools.ScrubQuorum.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapProposalActionTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.ActionTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.ActionTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposals.ActionTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapProposalCooldownTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.CooldownTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.CooldownTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposals.CooldownTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapProposalExecuteTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.ExecuteTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.ExecuteTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposals.ExecuteTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapVoteDelayTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.VoteDelayTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.VoteDelayTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposals.VoteDelayTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapVoteTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.VoteTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.VoteTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposals.VoteTime.Bootstrap(newVal, opts)
	})
}

func Test_AllODaoBoostrapFunctions(t *testing.T) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Create new settings
	newOdaoSettings := oracle.OracleDaoSettingsDetails{}
	newOdaoSettings.Members.ChallengeCooldown.Value.Set(tests.ODaoDefaults.Members.ChallengeCooldown.Value.Formatted() + time.Hour)
	newOdaoSettings.Members.ChallengeCost.Value = big.NewInt(0).Add(tests.ODaoDefaults.Members.ChallengeCost.Value, eth.EthToWei(1))
	newOdaoSettings.Members.ChallengeWindow.Value.Set(tests.ODaoDefaults.Members.ChallengeWindow.Value.Formatted() + time.Hour)
	newOdaoSettings.Members.Quorum.Value.Set(tests.ODaoDefaults.Members.Quorum.Value.Formatted() + 0.15)
	newOdaoSettings.Members.RplBond.Value = big.NewInt(0).Add(tests.ODaoDefaults.Members.RplBond.Value, eth.EthToWei(1000))
	newOdaoSettings.Members.UnbondedMinipoolMax.Value.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMax.Value.Formatted() + 5)
	newOdaoSettings.Members.UnbondedMinipoolMinFee.Value.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMinFee.Value.Formatted() + 0.1)
	newOdaoSettings.Minipools.BondReductionCancellationQuorum.Value.Set(tests.ODaoDefaults.Minipools.BondReductionCancellationQuorum.Value.Formatted() + 0.15)
	newOdaoSettings.Minipools.BondReductionWindowLength.Value.Set(tests.ODaoDefaults.Minipools.BondReductionWindowLength.Value.Formatted() + time.Hour)
	newOdaoSettings.Minipools.BondReductionWindowStart.Value.Set(tests.ODaoDefaults.Minipools.BondReductionWindowStart.Value.Formatted() + time.Hour)
	newOdaoSettings.Minipools.IsScrubPenaltyEnabled.Value = !tests.ODaoDefaults.Minipools.IsScrubPenaltyEnabled.Value
	newOdaoSettings.Minipools.PromotionScrubPeriod.Value.Set(tests.ODaoDefaults.Minipools.PromotionScrubPeriod.Value.Formatted() + time.Hour)
	newOdaoSettings.Minipools.ScrubPeriod.Value.Set(tests.ODaoDefaults.Minipools.ScrubPeriod.Value.Formatted() + time.Hour)
	newOdaoSettings.Minipools.ScrubQuorum.Value.Set(tests.ODaoDefaults.Minipools.ScrubQuorum.Value.Formatted() + 0.15)
	newOdaoSettings.Proposals.ActionTime.Value.Set(tests.ODaoDefaults.Proposals.ActionTime.Value.Formatted() + time.Hour)
	newOdaoSettings.Proposals.CooldownTime.Value.Set(tests.ODaoDefaults.Proposals.CooldownTime.Value.Formatted() + time.Hour)
	newOdaoSettings.Proposals.ExecuteTime.Value.Set(tests.ODaoDefaults.Proposals.ExecuteTime.Value.Formatted() + time.Hour)
	newOdaoSettings.Proposals.VoteDelayTime.Value.Set(tests.ODaoDefaults.Proposals.VoteDelayTime.Value.Formatted() + time.Hour)
	newOdaoSettings.Proposals.VoteTime.Value.Set(tests.ODaoDefaults.Proposals.VoteTime.Value.Formatted() + time.Hour)

	// Ensure they're all different from the default
	settings_test.EnsureDifferentDetails(t.Fatalf, &tests.ODaoDefaults, &newOdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	txInfos := []*core.TransactionInfo{}
	bootstrappers := []func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Members.ChallengeCooldown.Bootstrap(newOdaoSettings.Members.ChallengeCooldown.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Members.ChallengeCost.Bootstrap(newOdaoSettings.Members.ChallengeCost.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Members.ChallengeWindow.Bootstrap(newOdaoSettings.Members.ChallengeWindow.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Members.Quorum.Bootstrap(newOdaoSettings.Members.Quorum.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Members.RplBond.Bootstrap(newOdaoSettings.Members.RplBond.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Members.UnbondedMinipoolMax.Bootstrap(newOdaoSettings.Members.UnbondedMinipoolMax.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Members.UnbondedMinipoolMinFee.Bootstrap(newOdaoSettings.Members.UnbondedMinipoolMinFee.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipools.BondReductionCancellationQuorum.Bootstrap(newOdaoSettings.Minipools.BondReductionCancellationQuorum.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipools.BondReductionWindowLength.Bootstrap(newOdaoSettings.Minipools.BondReductionWindowLength.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipools.BondReductionWindowStart.Bootstrap(newOdaoSettings.Minipools.BondReductionWindowStart.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipools.IsScrubPenaltyEnabled.Bootstrap(newOdaoSettings.Minipools.IsScrubPenaltyEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipools.PromotionScrubPeriod.Bootstrap(newOdaoSettings.Minipools.PromotionScrubPeriod.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipools.ScrubPeriod.Bootstrap(newOdaoSettings.Minipools.ScrubPeriod.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipools.ScrubQuorum.Bootstrap(newOdaoSettings.Minipools.ScrubQuorum.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposals.ActionTime.Bootstrap(newOdaoSettings.Proposals.ActionTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposals.CooldownTime.Bootstrap(newOdaoSettings.Proposals.CooldownTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposals.ExecuteTime.Bootstrap(newOdaoSettings.Proposals.ExecuteTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposals.VoteDelayTime.Bootstrap(newOdaoSettings.Proposals.VoteDelayTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposals.VoteTime.Bootstrap(newOdaoSettings.Proposals.VoteTime.Value, opts)
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
		odaoMgr.Settings.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	settings_test.EnsureSameDetails(t.Fatalf, &newOdaoSettings, odaoMgr.Settings.OracleDaoSettingsDetails)
	t.Log("New settings match expected settings")
}

func testOdaoParameterBootstrap(t *testing.T, setter func(*oracle.OracleDaoSettingsDetails), bootstrapper func() (*core.TransactionInfo, error)) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Get the original settings
	var settings oracle.OracleDaoSettingsDetails
	settings_test.Clone(t, &tests.ODaoDefaults, &settings)
	pass := settings_test.EnsureSameDetails(t.Errorf, &tests.ODaoDefaults, &settings)
	if !pass {
		t.Fatalf("Details differed unexpectedly, can't continue")
	}
	t.Log("Cloned default settings")

	// Set the new setting
	setter(&settings)
	pass = settings_test.EnsureSameDetails(t.Logf, &tests.ODaoDefaults, &settings)
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
		odaoMgr.Settings.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	settings_test.EnsureSameDetails(t.Fatalf, &settings, odaoMgr.Settings.OracleDaoSettingsDetails)
	t.Log("New settings match expected settings")
}
