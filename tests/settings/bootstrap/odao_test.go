package bootstrap_test

import (
	"fmt"
	"math/big"
	"runtime/debug"
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
	newVal := tests.ODaoDefaults.Member.ChallengeCooldown.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.ChallengeCooldown.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeCooldown.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapChallengeCost(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Member.ChallengeCost.Get(), eth.EthToWei(1))
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.ChallengeCost.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeCost.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapChallengeWindow(t *testing.T) {
	newVal := tests.ODaoDefaults.Member.ChallengeWindow.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.ChallengeWindow.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeWindow.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapQuorum(t *testing.T) {
	newVal := tests.ODaoDefaults.Member.Quorum.Formatted() + 0.15
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.Quorum.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.Quorum.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapRplBond(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Member.RplBond.Get(), eth.EthToWei(1000))
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.RplBond.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.RplBond.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMax(t *testing.T) {
	newVal := tests.ODaoDefaults.Member.UnbondedMinipoolMax.Formatted() + 5
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.UnbondedMinipoolMax.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.UnbondedMinipoolMax.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapUnbondedMinipoolMinFee(t *testing.T) {
	newVal := tests.ODaoDefaults.Member.UnbondedMinipoolMinFee.Formatted() + 0.1
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.UnbondedMinipoolMinFee.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.UnbondedMinipoolMinFee.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapBondReductionCancellationQuorum(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.BondReductionCancellationQuorum.Formatted() + 0.15
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.BondReductionCancellationQuorum.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionCancellationQuorum.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapBondReductionWindowLength(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.BondReductionWindowLength.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.BondReductionWindowLength.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionWindowLength.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapBondReductionWindowStart(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.BondReductionWindowStart.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.BondReductionWindowStart.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionWindowStart.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapScrubPenaltyEnabled(t *testing.T) {
	newVal := !tests.ODaoDefaults.Minipool.IsScrubPenaltyEnabled.Get()
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.IsScrubPenaltyEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.IsScrubPenaltyEnabled.Bootstrap(newVal, opts)
	})
}

func Test_BootstrapPromotionScrubPeriod(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.PromotionScrubPeriod.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.PromotionScrubPeriod.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.PromotionScrubPeriod.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapScrubPeriod(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.ScrubPeriod.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.ScrubPeriod.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.ScrubPeriod.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapScrubQuorum(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.ScrubQuorum.Formatted() + 0.15
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.ScrubQuorum.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.ScrubQuorum.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapProposalActionTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.ActionTime.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.ActionTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.ActionTime.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapProposalCooldownTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.CooldownTime.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.CooldownTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.CooldownTime.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapProposalExecuteTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.ExecuteTime.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.ExecuteTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.ExecuteTime.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapVoteDelayTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.VoteDelayTime.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.VoteDelayTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.VoteDelayTime.Bootstrap(core.GetValueForUint256(newVal), opts)
	})
}

func Test_BootstrapVoteTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.VoteTime.Formatted() + time.Hour
	testOdaoParameterBootstrap(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.VoteTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.VoteTime.Bootstrap(core.GetValueForUint256(newVal), opts)
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
	odaoMgr, err := oracle.NewOracleDaoManager(mgr.RocketPool)
	if err != nil {
		t.Fatal("error creating oracle DAO manager: %w", err)
	}
	newOdaoSettings := *odaoMgr.Settings
	newOdaoSettings.Member.ChallengeCooldown.Set(tests.ODaoDefaults.Member.ChallengeCooldown.Formatted() + time.Hour)
	newOdaoSettings.Member.ChallengeCost.Set(big.NewInt(0).Add(tests.ODaoDefaults.Member.ChallengeCost.Get(), eth.EthToWei(1)))
	newOdaoSettings.Member.ChallengeWindow.Set(tests.ODaoDefaults.Member.ChallengeWindow.Formatted() + time.Hour)
	newOdaoSettings.Member.Quorum.Set(tests.ODaoDefaults.Member.Quorum.Formatted() + 0.15)
	newOdaoSettings.Member.RplBond.Set(big.NewInt(0).Add(tests.ODaoDefaults.Member.RplBond.Get(), eth.EthToWei(1000)))
	newOdaoSettings.Member.UnbondedMinipoolMax.Set(tests.ODaoDefaults.Member.UnbondedMinipoolMax.Formatted() + 5)
	newOdaoSettings.Member.UnbondedMinipoolMinFee.Set(tests.ODaoDefaults.Member.UnbondedMinipoolMinFee.Formatted() + 0.1)
	newOdaoSettings.Minipool.BondReductionCancellationQuorum.Set(tests.ODaoDefaults.Minipool.BondReductionCancellationQuorum.Formatted() + 0.15)
	newOdaoSettings.Minipool.BondReductionWindowLength.Set(tests.ODaoDefaults.Minipool.BondReductionWindowLength.Formatted() + time.Hour)
	newOdaoSettings.Minipool.BondReductionWindowStart.Set(tests.ODaoDefaults.Minipool.BondReductionWindowStart.Formatted() + time.Hour)
	newOdaoSettings.Minipool.IsScrubPenaltyEnabled.Set(!tests.ODaoDefaults.Minipool.IsScrubPenaltyEnabled.Get())
	newOdaoSettings.Minipool.PromotionScrubPeriod.Set(tests.ODaoDefaults.Minipool.PromotionScrubPeriod.Formatted() + time.Hour)
	newOdaoSettings.Minipool.ScrubPeriod.Set(tests.ODaoDefaults.Minipool.ScrubPeriod.Formatted() + time.Hour)
	newOdaoSettings.Minipool.ScrubQuorum.Set(tests.ODaoDefaults.Minipool.ScrubQuorum.Formatted() + 0.15)
	newOdaoSettings.Proposal.ActionTime.Set(tests.ODaoDefaults.Proposal.ActionTime.Formatted() + time.Hour)
	newOdaoSettings.Proposal.CooldownTime.Set(tests.ODaoDefaults.Proposal.CooldownTime.Formatted() + time.Hour)
	newOdaoSettings.Proposal.ExecuteTime.Set(tests.ODaoDefaults.Proposal.ExecuteTime.Formatted() + time.Hour)
	newOdaoSettings.Proposal.VoteDelayTime.Set(tests.ODaoDefaults.Proposal.VoteDelayTime.Formatted() + time.Hour)
	newOdaoSettings.Proposal.VoteTime.Set(tests.ODaoDefaults.Proposal.VoteTime.Formatted() + time.Hour)

	// Ensure they're all different from the default
	settings_test.EnsureDifferentDetails(t.Fatalf, &tests.ODaoDefaults, &newOdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	txInfos := []*core.TransactionInfo{}
	bootstrappers := []func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.ChallengeCooldown.Bootstrap(newOdaoSettings.Member.ChallengeCooldown.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.ChallengeCost.Bootstrap(newOdaoSettings.Member.ChallengeCost.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.ChallengeWindow.Bootstrap(newOdaoSettings.Member.ChallengeWindow.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.Quorum.Bootstrap(newOdaoSettings.Member.Quorum.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.RplBond.Bootstrap(newOdaoSettings.Member.RplBond.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.UnbondedMinipoolMax.Bootstrap(newOdaoSettings.Member.UnbondedMinipoolMax.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Member.UnbondedMinipoolMinFee.Bootstrap(newOdaoSettings.Member.UnbondedMinipoolMinFee.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.BondReductionCancellationQuorum.Bootstrap(newOdaoSettings.Minipool.BondReductionCancellationQuorum.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.BondReductionWindowLength.Bootstrap(newOdaoSettings.Minipool.BondReductionWindowLength.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.BondReductionWindowStart.Bootstrap(newOdaoSettings.Minipool.BondReductionWindowStart.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.IsScrubPenaltyEnabled.Bootstrap(newOdaoSettings.Minipool.IsScrubPenaltyEnabled.Get(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.PromotionScrubPeriod.Bootstrap(newOdaoSettings.Minipool.PromotionScrubPeriod.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.ScrubPeriod.Bootstrap(newOdaoSettings.Minipool.ScrubPeriod.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Minipool.ScrubQuorum.Bootstrap(newOdaoSettings.Minipool.ScrubQuorum.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.ActionTime.Bootstrap(newOdaoSettings.Proposal.ActionTime.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.CooldownTime.Bootstrap(newOdaoSettings.Proposal.CooldownTime.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.ExecuteTime.Bootstrap(newOdaoSettings.Proposal.ExecuteTime.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.VoteDelayTime.Bootstrap(newOdaoSettings.Proposal.VoteDelayTime.Raw(), opts)
		},
		func() (*core.TransactionInfo, error) {
			return odaoMgr.Settings.Proposal.VoteTime.Bootstrap(newOdaoSettings.Proposal.VoteTime.Raw(), opts)
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
		core.QueryAllFields(odaoMgr.Settings, mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	settings_test.EnsureSameDetails(t.Fatalf, &newOdaoSettings, odaoMgr.Settings)
	t.Log("New settings match expected settings")
}

func testOdaoParameterBootstrap(t *testing.T, setter func(*oracle.OracleDaoSettings), bootstrapper func() (*core.TransactionInfo, error)) {
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
	odaoMgr, err := oracle.NewOracleDaoManager(mgr.RocketPool)
	if err != nil {
		t.Fatal("error creating oracle DAO manager: %w", err)
	}
	settings := *odaoMgr.Settings
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
		core.QueryAllFields(odaoMgr.Settings, mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	settings_test.EnsureSameDetails(t.Fatalf, &settings, odaoMgr.Settings)
	t.Log("New settings match expected settings")
}
