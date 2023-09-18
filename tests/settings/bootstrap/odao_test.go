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
	newVal.Set(tests.ODaoDefaults.Members.ChallengeCooldown.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeCooldown.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapChallengeCooldown(newVal, opts)
	})
}

func Test_BootstrapChallengeCost(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Members.ChallengeCost, eth.EthToWei(1))
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeCost = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapChallengeCost(newVal, opts)
	})
}

func Test_BootstrapChallengeWindow(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Members.ChallengeWindow.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeWindow = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapChallengeWindow(newVal, opts)
	})
}

func Test_BootstrapQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Members.Quorum.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.Quorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapQuorum(newVal, opts)
	})
}

func Test_BootstrapRplBond(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Members.RplBond, eth.EthToWei(1000))
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.RplBond = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapRplBond(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMax(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMax.Formatted() + 5)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.UnbondedMinipoolMax.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapUnbondedMinipoolMax(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMinFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMinFee.Formatted() + 0.1)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Members.UnbondedMinipoolMinFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapUnbondedMinipoolMinFee(newVal, opts)
	})
}

func Test_BootstrapBondReductionCancellationQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionCancellationQuorum.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionCancellationQuorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapBondReductionCancellationQuorum(newVal, opts)
	})
}

func Test_BootstrapBondReductionWindowLength(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionWindowLength.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionWindowLength.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapBondReductionWindowLength(newVal, opts)
	})
}

func Test_BootstrapBondReductionWindowStart(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionWindowStart.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionWindowStart.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapBondReductionWindowStart(newVal, opts)
	})
}

func Test_BootstrapScrubPenaltyEnabled(t *testing.T) {
	newVal := !tests.ODaoDefaults.Minipools.IsScrubPenaltyEnabled
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.IsScrubPenaltyEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapScrubPenaltyEnabled(newVal, opts)
	})
}

func Test_BootstrapPromotionScrubPeriod(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.PromotionScrubPeriod.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.PromotionScrubPeriod.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapPromotionScrubPeriod(newVal, opts)
	})
}

func Test_BootstrapScrubPeriod(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.ScrubPeriod.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.ScrubPeriod.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapScrubPeriod(newVal, opts)
	})
}

func Test_BootstrapScrubQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipools.ScrubQuorum.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipools.ScrubQuorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapScrubQuorum(newVal, opts)
	})
}

func Test_BootstrapProposalActionTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.ActionTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.ActionTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapProposalActionTime(newVal, opts)
	})
}

func Test_BootstrapProposalCooldownTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.CooldownTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.CooldownTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapProposalCooldownTime(newVal, opts)
	})
}

func Test_BootstrapProposalExecuteTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.ExecuteTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.ExecuteTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapProposalExecuteTime(newVal, opts)
	})
}

func Test_BootstrapVoteDelayTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.VoteDelayTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.VoteDelayTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapVoteDelayTime(newVal, opts)
	})
}

func Test_BootstrapVoteTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.VoteTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposals.VoteTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.BootstrapVoteTime(newVal, opts)
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
	newOdaoSettings.Members.ChallengeCooldown.Set(tests.ODaoDefaults.Members.ChallengeCooldown.Formatted() + time.Hour)
	newOdaoSettings.Members.ChallengeCost = big.NewInt(0).Add(tests.ODaoDefaults.Members.ChallengeCost, eth.EthToWei(1))
	newOdaoSettings.Members.ChallengeWindow.Set(tests.ODaoDefaults.Members.ChallengeWindow.Formatted() + time.Hour)
	newOdaoSettings.Members.Quorum.Set(tests.ODaoDefaults.Members.Quorum.Formatted() + 0.15)
	newOdaoSettings.Members.RplBond = big.NewInt(0).Add(tests.ODaoDefaults.Members.RplBond, eth.EthToWei(1000))
	newOdaoSettings.Members.UnbondedMinipoolMax.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMax.Formatted() + 5)
	newOdaoSettings.Members.UnbondedMinipoolMinFee.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMinFee.Formatted() + 0.1)
	newOdaoSettings.Minipools.BondReductionCancellationQuorum.Set(tests.ODaoDefaults.Minipools.BondReductionCancellationQuorum.Formatted() + 0.15)
	newOdaoSettings.Minipools.BondReductionWindowLength.Set(tests.ODaoDefaults.Minipools.BondReductionWindowLength.Formatted() + time.Hour)
	newOdaoSettings.Minipools.BondReductionWindowStart.Set(tests.ODaoDefaults.Minipools.BondReductionWindowStart.Formatted() + time.Hour)
	newOdaoSettings.Minipools.IsScrubPenaltyEnabled = !tests.ODaoDefaults.Minipools.IsScrubPenaltyEnabled
	newOdaoSettings.Minipools.PromotionScrubPeriod.Set(tests.ODaoDefaults.Minipools.PromotionScrubPeriod.Formatted() + time.Hour)
	newOdaoSettings.Minipools.ScrubPeriod.Set(tests.ODaoDefaults.Minipools.ScrubPeriod.Formatted() + time.Hour)
	newOdaoSettings.Minipools.ScrubQuorum.Set(tests.ODaoDefaults.Minipools.ScrubQuorum.Formatted() + 0.15)
	newOdaoSettings.Proposals.ActionTime.Set(tests.ODaoDefaults.Proposals.ActionTime.Formatted() + time.Hour)
	newOdaoSettings.Proposals.CooldownTime.Set(tests.ODaoDefaults.Proposals.CooldownTime.Formatted() + time.Hour)
	newOdaoSettings.Proposals.ExecuteTime.Set(tests.ODaoDefaults.Proposals.ExecuteTime.Formatted() + time.Hour)
	newOdaoSettings.Proposals.VoteDelayTime.Set(tests.ODaoDefaults.Proposals.VoteDelayTime.Formatted() + time.Hour)
	newOdaoSettings.Proposals.VoteTime.Set(tests.ODaoDefaults.Proposals.VoteTime.Formatted() + time.Hour)

	// Ensure they're all different from the default
	settings_test.EnsureDifferentDetails(t.Fatalf, &tests.ODaoDefaults, &newOdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	txInfos := []*core.TransactionInfo{}
	bootstrappers := []func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapChallengeCooldown(newOdaoSettings.Members.ChallengeCooldown, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapChallengeCost(newOdaoSettings.Members.ChallengeCost, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapChallengeWindow(newOdaoSettings.Members.ChallengeWindow, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapQuorum(newOdaoSettings.Members.Quorum, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapRplBond(newOdaoSettings.Members.RplBond, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapUnbondedMinipoolMax(newOdaoSettings.Members.UnbondedMinipoolMax, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapUnbondedMinipoolMinFee(newOdaoSettings.Members.UnbondedMinipoolMinFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapBondReductionCancellationQuorum(newOdaoSettings.Minipools.BondReductionCancellationQuorum, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapBondReductionWindowLength(newOdaoSettings.Minipools.BondReductionWindowLength, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapBondReductionWindowStart(newOdaoSettings.Minipools.BondReductionWindowStart, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapScrubPenaltyEnabled(newOdaoSettings.Minipools.IsScrubPenaltyEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapPromotionScrubPeriod(newOdaoSettings.Minipools.PromotionScrubPeriod, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapScrubPeriod(newOdaoSettings.Minipools.ScrubPeriod, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapScrubQuorum(newOdaoSettings.Minipools.ScrubQuorum, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapProposalActionTime(newOdaoSettings.Proposals.ActionTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapProposalCooldownTime(newOdaoSettings.Proposals.CooldownTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapProposalExecuteTime(newOdaoSettings.Proposals.ExecuteTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapVoteDelayTime(newOdaoSettings.Proposals.VoteDelayTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.BootstrapVoteTime(newOdaoSettings.Proposals.VoteTime, opts)
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
