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
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Member.ChallengeCooldown.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.ChallengeCooldown.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeCooldown.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapChallengeCost(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Member.ChallengeCost.Value, eth.EthToWei(1))
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.ChallengeCost.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeCost.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapChallengeWindow(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Member.ChallengeWindow.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.ChallengeWindow.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeWindow.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapQuorum(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Member.Quorum.Value.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.Quorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.Quorum.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapRplBond(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Member.RplBond.Value, eth.EthToWei(1000))
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.RplBond.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.RplBond.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMax(t *testing.T) {
	newVal := core.Uint256Parameter[uint64]{}
	newVal.Set(tests.ODaoDefaults.Member.UnbondedMinipoolMax.Value.Formatted() + 5)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.UnbondedMinipoolMax.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.UnbondedMinipoolMax.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMinFee(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Member.UnbondedMinipoolMinFee.Value.Formatted() + 0.1)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.UnbondedMinipoolMinFee.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.UnbondedMinipoolMinFee.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBondReductionCancellationQuorum(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipool.BondReductionCancellationQuorum.Value.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.BondReductionCancellationQuorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionCancellationQuorum.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBondReductionWindowLength(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipool.BondReductionWindowLength.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.BondReductionWindowLength.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionWindowLength.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapBondReductionWindowStart(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipool.BondReductionWindowStart.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.BondReductionWindowStart.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionWindowStart.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapScrubPenaltyEnabled(t *testing.T) {
	newVal := !tests.ODaoDefaults.Minipool.IsScrubPenaltyEnabled.Value
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.IsScrubPenaltyEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.IsScrubPenaltyEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapPromotionScrubPeriod(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipool.PromotionScrubPeriod.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.PromotionScrubPeriod.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.PromotionScrubPeriod.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapScrubPeriod(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipool.ScrubPeriod.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.ScrubPeriod.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.ScrubPeriod.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapScrubQuorum(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipool.ScrubQuorum.Value.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.ScrubQuorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.ScrubQuorum.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapProposalActionTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.ActionTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.ActionTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.ActionTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapProposalCooldownTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.CooldownTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.CooldownTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.CooldownTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapProposalExecuteTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.ExecuteTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.ExecuteTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.ExecuteTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapVoteDelayTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.VoteDelayTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.VoteDelayTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.VoteDelayTime.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapVoteTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.VoteTime.Value.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.VoteTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.VoteTime.Bootstrap(newVal, opts)
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
	odaoMgr, err := oracle.NewOracleDaoManager(mgr.RocketPool)
	if err != nil {
		t.Fatal("error creating oracle DAO manager: %w", err)
	}
	newOdaoSettings := *odaoMgr.Settings.OracleDaoSettingsDetails
	newOdaoSettings.Member.ChallengeCooldown.Value.Set(tests.ODaoDefaults.Member.ChallengeCooldown.Value.Formatted() + time.Hour)
	newOdaoSettings.Member.ChallengeCost.Value = big.NewInt(0).Add(tests.ODaoDefaults.Member.ChallengeCost.Value, eth.EthToWei(1))
	newOdaoSettings.Member.ChallengeWindow.Value.Set(tests.ODaoDefaults.Member.ChallengeWindow.Value.Formatted() + time.Hour)
	newOdaoSettings.Member.Quorum.Value.Set(tests.ODaoDefaults.Member.Quorum.Value.Formatted() + 0.15)
	newOdaoSettings.Member.RplBond.Value = big.NewInt(0).Add(tests.ODaoDefaults.Member.RplBond.Value, eth.EthToWei(1000))
	newOdaoSettings.Member.UnbondedMinipoolMax.Value.Set(tests.ODaoDefaults.Member.UnbondedMinipoolMax.Value.Formatted() + 5)
	newOdaoSettings.Member.UnbondedMinipoolMinFee.Value.Set(tests.ODaoDefaults.Member.UnbondedMinipoolMinFee.Value.Formatted() + 0.1)
	newOdaoSettings.Minipool.BondReductionCancellationQuorum.Value.Set(tests.ODaoDefaults.Minipool.BondReductionCancellationQuorum.Value.Formatted() + 0.15)
	newOdaoSettings.Minipool.BondReductionWindowLength.Value.Set(tests.ODaoDefaults.Minipool.BondReductionWindowLength.Value.Formatted() + time.Hour)
	newOdaoSettings.Minipool.BondReductionWindowStart.Value.Set(tests.ODaoDefaults.Minipool.BondReductionWindowStart.Value.Formatted() + time.Hour)
	newOdaoSettings.Minipool.IsScrubPenaltyEnabled.Value = !tests.ODaoDefaults.Minipool.IsScrubPenaltyEnabled.Value
	newOdaoSettings.Minipool.PromotionScrubPeriod.Value.Set(tests.ODaoDefaults.Minipool.PromotionScrubPeriod.Value.Formatted() + time.Hour)
	newOdaoSettings.Minipool.ScrubPeriod.Value.Set(tests.ODaoDefaults.Minipool.ScrubPeriod.Value.Formatted() + time.Hour)
	newOdaoSettings.Minipool.ScrubQuorum.Value.Set(tests.ODaoDefaults.Minipool.ScrubQuorum.Value.Formatted() + 0.15)
	newOdaoSettings.Proposal.ActionTime.Value.Set(tests.ODaoDefaults.Proposal.ActionTime.Value.Formatted() + time.Hour)
	newOdaoSettings.Proposal.CooldownTime.Value.Set(tests.ODaoDefaults.Proposal.CooldownTime.Value.Formatted() + time.Hour)
	newOdaoSettings.Proposal.ExecuteTime.Value.Set(tests.ODaoDefaults.Proposal.ExecuteTime.Value.Formatted() + time.Hour)
	newOdaoSettings.Proposal.VoteDelayTime.Value.Set(tests.ODaoDefaults.Proposal.VoteDelayTime.Value.Formatted() + time.Hour)
	newOdaoSettings.Proposal.VoteTime.Value.Set(tests.ODaoDefaults.Proposal.VoteTime.Value.Formatted() + time.Hour)

	// Ensure they're all different from the default
	settings_test.EnsureDifferentDetails(t.Fatalf, &tests.ODaoDefaults, &newOdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	txInfos := []*core.TransactionInfo{}
	bootstrappers := []func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.ChallengeCooldown.Bootstrap(newOdaoSettings.Member.ChallengeCooldown.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.ChallengeCost.Bootstrap(newOdaoSettings.Member.ChallengeCost.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.ChallengeWindow.Bootstrap(newOdaoSettings.Member.ChallengeWindow.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.Quorum.Bootstrap(newOdaoSettings.Member.Quorum.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.RplBond.Bootstrap(newOdaoSettings.Member.RplBond.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.UnbondedMinipoolMax.Bootstrap(newOdaoSettings.Member.UnbondedMinipoolMax.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.UnbondedMinipoolMinFee.Bootstrap(newOdaoSettings.Member.UnbondedMinipoolMinFee.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.BondReductionCancellationQuorum.Bootstrap(newOdaoSettings.Minipool.BondReductionCancellationQuorum.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.BondReductionWindowLength.Bootstrap(newOdaoSettings.Minipool.BondReductionWindowLength.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.BondReductionWindowStart.Bootstrap(newOdaoSettings.Minipool.BondReductionWindowStart.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.IsScrubPenaltyEnabled.Bootstrap(newOdaoSettings.Minipool.IsScrubPenaltyEnabled.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.PromotionScrubPeriod.Bootstrap(newOdaoSettings.Minipool.PromotionScrubPeriod.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.ScrubPeriod.Bootstrap(newOdaoSettings.Minipool.ScrubPeriod.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.ScrubQuorum.Bootstrap(newOdaoSettings.Minipool.ScrubQuorum.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.ActionTime.Bootstrap(newOdaoSettings.Proposal.ActionTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.CooldownTime.Bootstrap(newOdaoSettings.Proposal.CooldownTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.ExecuteTime.Bootstrap(newOdaoSettings.Proposal.ExecuteTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.VoteDelayTime.Bootstrap(newOdaoSettings.Proposal.VoteDelayTime.Value, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.VoteTime.Bootstrap(newOdaoSettings.Proposal.VoteTime.Value, opts)
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
	odaoMgr, err := oracle.NewOracleDaoManager(mgr.RocketPool)
	if err != nil {
		t.Fatal("error creating oracle DAO manager: %w", err)
	}
	settings := *odaoMgr.Settings.OracleDaoSettingsDetails
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
