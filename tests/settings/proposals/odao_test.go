package proposals_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/settings"
	"github.com/rocket-pool/rocketpool-go/tests"
	settings_test "github.com/rocket-pool/rocketpool-go/tests/settings"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

func Test_ProposeChallengeCooldown(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Members.ChallengeCooldown.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeCooldown.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeChallengeCooldown(newVal, odao1.Transactor)
	})
}

func Test_ProposeChallengeCost(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Members.ChallengeCost, eth.EthToWei(1))
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeCost = newVal
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeChallengeCost(newVal, odao1.Transactor)
	})
}

func Test_ProposeChallengeWindow(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Members.ChallengeWindow.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.ChallengeWindow.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeChallengeWindow(newVal, odao1.Transactor)
	})
}

func Test_ProposeQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Members.Quorum.Formatted() + 0.15)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.Quorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeQuorum(newVal, odao1.Transactor)
	})
}

func Test_ProposeRplBond(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Members.RplBond, eth.EthToWei(1))
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.RplBond = newVal
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeRplBond(newVal, odao1.Transactor)
	})
}

func Test_ProposeUnbondedMinipoolMax(t *testing.T) {
	newVal := core.Parameter[uint64]{}
	newVal.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMax.Formatted() + 1)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.UnbondedMinipoolMax.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeUnbondedMinipoolMax(newVal, odao1.Transactor)
	})
}

func Test_ProposeUnbondedMinipoolMinFee(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Members.UnbondedMinipoolMinFee.Formatted() + 0.01)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Members.UnbondedMinipoolMinFee.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeUnbondedMinipoolMinFee(newVal, odao1.Transactor)
	})
}

func Test_ProposeBondReductionCancellationQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionCancellationQuorum.Formatted() + 0.15)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionCancellationQuorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeBondReductionCancellationQuorum(newVal, odao1.Transactor)
	})
}

func Test_ProposeBondReductionWindowLength(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionWindowLength.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionWindowLength.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeBondReductionWindowLength(newVal, odao1.Transactor)
	})
}

func Test_ProposeBondReductionWindowStart(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.BondReductionWindowStart.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.BondReductionWindowStart.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeBondReductionWindowStart(newVal, odao1.Transactor)
	})
}

func Test_ProposeScrubPenaltyEnabled(t *testing.T) {
	newVal := !tests.ODaoDefaults.Minipools.IsScrubPenaltyEnabled
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.IsScrubPenaltyEnabled = newVal
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeScrubPenaltyEnabled(newVal, odao1.Transactor)
	})
}

func Test_ProposePromotionScrubPeriod(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.PromotionScrubPeriod.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.PromotionScrubPeriod.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposePromotionScrubPeriod(newVal, odao1.Transactor)
	})
}

func Test_ProposeScrubPeriod(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipools.ScrubPeriod.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.ScrubPeriod.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeScrubPeriod(newVal, odao1.Transactor)
	})
}

func Test_ProposeScrubQuorum(t *testing.T) {
	newVal := core.Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipools.ScrubQuorum.Formatted() + 0.15)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Minipools.ScrubQuorum.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeScrubQuorum(newVal, odao1.Transactor)
	})
}

func Test_ProposeProposalActionTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.ActionTime.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.ActionTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeProposalActionTime(newVal, odao1.Transactor)
	})
}

func Test_ProposeProposalCooldownTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.CooldownTime.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.CooldownTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeProposalCooldownTime(newVal, odao1.Transactor)
	})
}

func Test_ProposeProposalExecuteTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.ExecuteTime.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.ExecuteTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeProposalExecuteTime(newVal, odao1.Transactor)
	})
}

func Test_ProposeVoteDelayTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.VoteDelayTime.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.VoteDelayTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeVoteDelayTime(newVal, odao1.Transactor)
	})
}

func Test_ProposeVoteTime(t *testing.T) {
	newVal := core.Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposals.VoteTime.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *settings.OracleDaoSettingsDetails) {
		newSettings.Proposals.VoteTime.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odao.ProposeVoteTime(newVal, odao1.Transactor)
	})
}

func testOdaoParameterProposal(t *testing.T, setter func(*settings.OracleDaoSettingsDetails), proposer func() (*core.TransactionInfo, error)) {
	// Revert to the initialized state at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToInitialized()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to initialized snapshot: %w", err))
		}
	})

	// Get the original settings
	var settings settings.OracleDaoSettingsDetails
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

	// Make sure there aren't any proposals
	err := rp.Query(func(mc *batch.MultiCaller) error {
		dpm.GetProposalCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting proposal count: %s", err.Error())
	}
	if dpm.Details.ProposalCount.Formatted() != 0 {
		t.Fatalf("expected 0 proposals but count was %d", dpm.Details.ProposalCount.Formatted())
	}

	// Run the proposer
	err = rp.CreateAndWaitForTransaction(proposer, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error submitting proposal: %s", err.Error())
	}

	// Make sure the actual network settings haven't changed
	err = rp.Query(func(mc *batch.MultiCaller) error {
		odao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error querying all updated details: %s", err)
	}
	settings_test.EnsureSameDetails(t.Fatalf, &tests.ODaoDefaults, &odao.Details)
	t.Log("Settings match the defaults after proposal creation, ok")

	// Make sure the proposal count is good
	err = rp.Query(func(mc *batch.MultiCaller) error {
		dpm.GetProposalCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting proposal count: %s", err.Error())
	}
	propCount := dpm.Details.ProposalCount.Formatted()
	if propCount != 1 {
		t.Fatalf("expected 1 proposal but count was %d", propCount)
	}
	t.Logf("Prop count = %d, ok", propCount)

	// Get the proposal
	pdaoProps, odaoProps, err := dpm.GetProposals(rp, propCount, true, nil)
	if err != nil {
		t.Fatalf("error getting proposals: %s", err.Error())
	}
	if len(pdaoProps) != 0 && len(odaoProps) != 1 {
		t.Fatalf("expected 0 pDAO prop and 1 oDAO prop but counts were %d and %d", len(pdaoProps), len(odaoProps))
	}
	prop := odaoProps[0]
	t.Logf("Got prop with ID %d: %s", prop.Details.ID.Formatted(), prop.Details.Message)

	// Skip enough time to allow voting
	voteDelayTime := settings.Proposals.VoteDelayTime.Formatted()
	waitSeconds := int(voteDelayTime.Seconds())
	err = mgr.IncreaseTime(waitSeconds)
	if err != nil {
		t.Fatalf("error skipping time by %d seconds: %s", waitSeconds, err.Error())
	}
	t.Logf("Skipped forward %d seconds (%s) to allow voting", waitSeconds, voteDelayTime)

	// Vote yay from node 1
	err = rp.CreateAndWaitForTransaction(func() (*core.TransactionInfo, error) {
		return op.VoteOnProposal(prop.Details.ID.Formatted(), true, odao1.Transactor)
	}, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error voting on proposal: %s", err.Error())
	}
	t.Log("Voted yes from node 1")

	// Vote yay from node 2 and execute it
	err = rp.BatchCreateAndWaitForTransactions([]func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return op.VoteOnProposal(prop.Details.ID.Formatted(), true, odao2.Transactor)
		},
		func() (*core.TransactionInfo, error) {
			return op.ExecuteProposal(prop.Details.ID.Formatted(), odao2.Transactor)
		},
	}, false, odao2.Transactor)
	if err != nil {
		t.Fatalf("error voting on and executing proposal: %s", err.Error())
	}
	t.Log("Voted yes from node 2 and executed")

	// Get new values and make sure they match
	err = rp.Query(func(mc *batch.MultiCaller) error {
		odao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error querying all updated details: %s", err.Error())
	}
	settings_test.EnsureSameDetails(t.Fatalf, &settings, &odao.Details)
	t.Log("New settings match expected settings - proposal succeeded")
}
