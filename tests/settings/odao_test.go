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

func Test_BootstrapChallengeCooldown(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Members.ChallengeCooldown.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeCooldown.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapChallengeCooldown(newVal, opts)
	})
}

func Test_BootstrapChallengeCost(t *testing.T) {
	newVal := big.NewInt(0).Add(odaoDefaults.Members.ChallengeCost, eth.EthToWei(1))
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeCost = newVal
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapChallengeCost(newVal, opts)
	})
}

func Test_BootstrapChallengeWindow(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Members.ChallengeWindow.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeWindow = newVal
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapChallengeWindow(newVal, opts)
	})
}

func Test_BootstrapQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(odaoDefaults.Members.Quorum.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.Quorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapQuorum(newVal, opts)
	})
}

func Test_BootstrapRplBond(t *testing.T) {
	newVal := big.NewInt(0).Add(odaoDefaults.Members.RplBond, eth.EthToWei(1000))
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.RplBond = newVal
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapRplBond(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMax(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(odaoDefaults.Members.UnbondedMinipoolMax.Formatted() + 5)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.UnbondedMinipoolMax.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapUnbondedMinipoolMax(newVal, opts)
	})
}

func Test_BootstrapUnbondedMinipoolMinFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(odaoDefaults.Members.UnbondedMinipoolMinFee.Formatted() + 0.1)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.UnbondedMinipoolMinFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapUnbondedMinipoolMinFee(newVal, opts)
	})
}

func Test_BootstrapBondReductionCancellationQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(odaoDefaults.Minipools.BondReductionCancellationQuorum.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionCancellationQuorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapBondReductionCancellationQuorum(newVal, opts)
	})
}

func Test_BootstrapBondReductionWindowLength(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Minipools.BondReductionWindowLength.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionWindowLength.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapBondReductionWindowLength(newVal, opts)
	})
}

func Test_BootstrapBondReductionWindowStart(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Minipools.BondReductionWindowStart.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionWindowStart.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapBondReductionWindowStart(newVal, opts)
	})
}

func Test_BootstrapScrubPenaltyEnabled(t *testing.T) {
	newVal := !odaoDefaults.Minipools.IsScrubPenaltyEnabled
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.IsScrubPenaltyEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapScrubPenaltyEnabled(newVal, opts)
	})
}

func Test_BootstrapPromotionScrubPeriod(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Minipools.PromotionScrubPeriod.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.PromotionScrubPeriod.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapPromotionScrubPeriod(newVal, opts)
	})
}

func Test_BootstrapScrubPeriod(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Minipools.ScrubPeriod.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.ScrubPeriod.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapScrubPeriod(newVal, opts)
	})
}

func Test_BootstrapScrubQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(odaoDefaults.Minipools.ScrubQuorum.Formatted() + 0.15)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.ScrubQuorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapScrubQuorum(newVal, opts)
	})
}

func Test_BootstrapProposalActionTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Proposals.ActionTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.ActionTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapProposalActionTime(newVal, opts)
	})
}

func Test_BootstrapProposalCooldownTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Proposals.CooldownTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.CooldownTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapProposalCooldownTime(newVal, opts)
	})
}

func Test_BootstrapProposalExecuteTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Proposals.ExecuteTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.ExecuteTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapProposalExecuteTime(newVal, opts)
	})
}

func Test_BootstrapVoteDelayTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Proposals.VoteDelayTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.VoteDelayTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapVoteDelayTime(newVal, opts)
	})
}

func Test_BootstrapVoteTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(odaoDefaults.Proposals.VoteTime.Formatted() + time.Hour)
	testOdaoParameterBootstrap(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.VoteTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.BootstrapVoteTime(newVal, opts)
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
	newOdaoSettings := settings.OracleDaoSettingsDetails{}
	newOdaoSettings.Members.ChallengeCooldown.Set(odaoDefaults.Members.ChallengeCooldown.Formatted() + time.Hour)
	newOdaoSettings.Members.ChallengeCost = big.NewInt(0).Add(odaoDefaults.Members.ChallengeCost, eth.EthToWei(1))
	newOdaoSettings.Members.ChallengeWindow.Set(odaoDefaults.Members.ChallengeWindow.Formatted() + time.Hour)
	newOdaoSettings.Members.Quorum.Set(odaoDefaults.Members.Quorum.Formatted() + 0.15)
	newOdaoSettings.Members.RplBond = big.NewInt(0).Add(odaoDefaults.Members.RplBond, eth.EthToWei(1000))
	newOdaoSettings.Members.UnbondedMinipoolMax.Set(odaoDefaults.Members.UnbondedMinipoolMax.Formatted() + 5)
	newOdaoSettings.Members.UnbondedMinipoolMinFee.Set(odaoDefaults.Members.UnbondedMinipoolMinFee.Formatted() + 0.1)
	newOdaoSettings.Minipools.BondReductionCancellationQuorum.Set(odaoDefaults.Minipools.BondReductionCancellationQuorum.Formatted() + 0.15)
	newOdaoSettings.Minipools.BondReductionWindowLength.Set(odaoDefaults.Minipools.BondReductionWindowLength.Formatted() + time.Hour)
	newOdaoSettings.Minipools.BondReductionWindowStart.Set(odaoDefaults.Minipools.BondReductionWindowStart.Formatted() + time.Hour)
	newOdaoSettings.Minipools.IsScrubPenaltyEnabled = !odaoDefaults.Minipools.IsScrubPenaltyEnabled
	newOdaoSettings.Minipools.PromotionScrubPeriod.Set(odaoDefaults.Minipools.PromotionScrubPeriod.Formatted() + time.Hour)
	newOdaoSettings.Minipools.ScrubPeriod.Set(odaoDefaults.Minipools.ScrubPeriod.Formatted() + time.Hour)
	newOdaoSettings.Minipools.ScrubQuorum.Set(odaoDefaults.Minipools.ScrubQuorum.Formatted() + 0.15)
	newOdaoSettings.Proposals.ActionTime.Set(odaoDefaults.Proposals.ActionTime.Formatted() + time.Hour)
	newOdaoSettings.Proposals.CooldownTime.Set(odaoDefaults.Proposals.CooldownTime.Formatted() + time.Hour)
	newOdaoSettings.Proposals.ExecuteTime.Set(odaoDefaults.Proposals.ExecuteTime.Formatted() + time.Hour)
	newOdaoSettings.Proposals.VoteDelayTime.Set(odaoDefaults.Proposals.VoteDelayTime.Formatted() + time.Hour)
	newOdaoSettings.Proposals.VoteTime.Set(odaoDefaults.Proposals.VoteTime.Formatted() + time.Hour)

	// Ensure they're all different from the default
	EnsureDifferentDetails(t.Fatalf, &odaoDefaults, &newOdaoSettings)
	t.Log("Updated details all differ from original details")

	// Set the new settings
	txInfos := []*core.TransactionInfo{}
	bootstrappers := []func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapChallengeCooldown(newOdaoSettings.Members.ChallengeCooldown, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapChallengeCost(newOdaoSettings.Members.ChallengeCost, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapChallengeWindow(newOdaoSettings.Members.ChallengeWindow, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapQuorum(newOdaoSettings.Members.Quorum, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapRplBond(newOdaoSettings.Members.RplBond, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapUnbondedMinipoolMax(newOdaoSettings.Members.UnbondedMinipoolMax, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapUnbondedMinipoolMinFee(newOdaoSettings.Members.UnbondedMinipoolMinFee, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapBondReductionCancellationQuorum(newOdaoSettings.Minipools.BondReductionCancellationQuorum, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapBondReductionWindowLength(newOdaoSettings.Minipools.BondReductionWindowLength, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapBondReductionWindowStart(newOdaoSettings.Minipools.BondReductionWindowStart, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapScrubPenaltyEnabled(newOdaoSettings.Minipools.IsScrubPenaltyEnabled, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapPromotionScrubPeriod(newOdaoSettings.Minipools.PromotionScrubPeriod, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapScrubPeriod(newOdaoSettings.Minipools.ScrubPeriod, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapScrubQuorum(newOdaoSettings.Minipools.ScrubQuorum, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapProposalActionTime(newOdaoSettings.Proposals.ActionTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapProposalCooldownTime(newOdaoSettings.Proposals.CooldownTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapProposalExecuteTime(newOdaoSettings.Proposals.ExecuteTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapVoteDelayTime(newOdaoSettings.Proposals.VoteDelayTime, opts)
		},
		func() (*core.TransactionInfo, error) {
			return odao.BootstrapVoteTime(newOdaoSettings.Proposals.VoteTime, opts)
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
		odao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	EnsureSameDetails(t.Fatalf, &newOdaoSettings, &odao.Details)
	t.Log("New settings match expected settings")
}

func testOdaoParameterBootstrap(t *testing.T, setter func(*settings.OracleDaoSettingsDetails), bootstrapper func() (*core.TransactionInfo, error)) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Get the original settings
	var settings settings.OracleDaoSettingsDetails
	Clone(t, &odaoDefaults, &settings)
	pass := EnsureSameDetails(t.Errorf, &odaoDefaults, &settings)
	if !pass {
		t.Fatalf("Details differed unexpectedly, can't continue")
	}
	t.Log("Cloned default settings")

	// Set the new setting
	setter(&settings)
	pass = EnsureSameDetails(t.Logf, &odaoDefaults, &settings)
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
		odao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all updated details: %w", err))
	}
	EnsureSameDetails(t.Fatalf, &settings, &odao.Details)
	t.Log("New settings match expected settings")
}
