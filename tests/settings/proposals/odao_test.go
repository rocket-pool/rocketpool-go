package proposals_test

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
)

func Test_ProposeChallengeCooldown(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Member.ChallengeCooldown.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.ChallengeCooldown.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeCooldown.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeChallengeCost(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Member.ChallengeCost.Value, eth.EthToWei(1))
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.ChallengeCost.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeCost.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeChallengeWindow(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Member.ChallengeWindow.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.ChallengeWindow.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeWindow.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeQuorum(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Member.Quorum.Value.Formatted() + 0.15)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.Quorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.Quorum.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeRplBond(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Member.RplBond.Value, eth.EthToWei(1))
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.RplBond.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.RplBond.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeUnbondedMinipoolMax(t *testing.T) {
	newVal := core.Uint256Parameter[uint64]{}
	newVal.Set(tests.ODaoDefaults.Member.UnbondedMinipoolMax.Value.Formatted() + 1)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.UnbondedMinipoolMax.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.UnbondedMinipoolMax.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeUnbondedMinipoolMinFee(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Member.UnbondedMinipoolMinFee.Value.Formatted() + 0.01)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Member.UnbondedMinipoolMinFee.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.UnbondedMinipoolMinFee.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeBondReductionCancellationQuorum(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipool.BondReductionCancellationQuorum.Value.Formatted() + 0.15)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.BondReductionCancellationQuorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionCancellationQuorum.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeBondReductionWindowLength(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipool.BondReductionWindowLength.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.BondReductionWindowLength.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionWindowLength.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeBondReductionWindowStart(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipool.BondReductionWindowStart.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.BondReductionWindowStart.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionWindowStart.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeScrubPenaltyEnabled(t *testing.T) {
	newVal := !tests.ODaoDefaults.Minipool.IsScrubPenaltyEnabled.Value
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.IsScrubPenaltyEnabled.Value = newVal
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.IsScrubPenaltyEnabled.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposePromotionScrubPeriod(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipool.PromotionScrubPeriod.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.PromotionScrubPeriod.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.PromotionScrubPeriod.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeScrubPeriod(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Minipool.ScrubPeriod.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.ScrubPeriod.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.ScrubPeriod.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeScrubQuorum(t *testing.T) {
	newVal := core.Uint256Parameter[float64]{}
	newVal.Set(tests.ODaoDefaults.Minipool.ScrubQuorum.Value.Formatted() + 0.15)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Minipool.ScrubQuorum.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.ScrubQuorum.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeProposalActionTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.ActionTime.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.ActionTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.ActionTime.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeProposalCooldownTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.CooldownTime.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.CooldownTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.CooldownTime.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeProposalExecuteTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.ExecuteTime.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.ExecuteTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.ExecuteTime.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeVoteDelayTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.VoteDelayTime.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.VoteDelayTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.VoteDelayTime.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeVoteTime(t *testing.T) {
	newVal := core.Uint256Parameter[time.Duration]{}
	newVal.Set(tests.ODaoDefaults.Proposal.VoteTime.Value.Formatted() + time.Hour)
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettingsDetails) {
		newSettings.Proposal.VoteTime.Value.SetRawValue(newVal.RawValue)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.VoteTime.ProposeSet(newVal, odao1.Transactor)
	})
}

func testOdaoParameterProposal(t *testing.T, setter func(*oracle.OracleDaoSettingsDetails), proposer func() (*core.TransactionInfo, error)) {
	// Revert to the initialized state at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToInitialized()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to initialized snapshot: %w", err))
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

	// Make sure there aren't any proposals
	err = rp.Query(func(mc *batch.MultiCaller) error {
		dpm.GetProposalCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting proposal count: %s", err.Error())
	}
	if dpm.ProposalCount.Formatted() != 0 {
		t.Fatalf("expected 0 proposals but count was %d", dpm.ProposalCount.Formatted())
	}

	// Run the proposer
	err = rp.CreateAndWaitForTransaction(proposer, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error submitting proposal: %s", err.Error())
	}

	// Make sure the actual network settings haven't changed
	err = rp.Query(func(mc *batch.MultiCaller) error {
		odaoMgr.Settings.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error querying all updated details: %s", err)
	}
	settings_test.EnsureSameDetails(t.Fatalf, &tests.ODaoDefaults, odaoMgr.Settings.OracleDaoSettingsDetails)
	t.Log("Settings match the defaults after proposal creation, ok")

	// Make sure the proposal count is good
	err = rp.Query(func(mc *batch.MultiCaller) error {
		dpm.GetProposalCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting proposal count: %s", err.Error())
	}
	propCount := dpm.ProposalCount.Formatted()
	if propCount != 1 {
		t.Fatalf("expected 1 proposal but count was %d", propCount)
	}
	t.Logf("Prop count = %d, ok", propCount)

	// Get the proposal
	pdaoProps, odaoProps, err := dpm.GetProposals(propCount, true, nil)
	if err != nil {
		t.Fatalf("error getting proposals: %s", err.Error())
	}
	if len(pdaoProps) != 0 && len(odaoProps) != 1 {
		t.Fatalf("expected 0 pDAO prop and 1 oDAO prop but counts were %d and %d", len(pdaoProps), len(odaoProps))
	}
	prop := odaoProps[0]
	t.Logf("Got prop with ID %d: %s", prop.ID.Formatted(), prop.Message)

	// Skip enough time to allow voting
	voteDelayTime := settings.Proposal.VoteDelayTime.Value.Formatted()
	waitSeconds := int(voteDelayTime.Seconds())
	err = mgr.IncreaseTime(waitSeconds)
	if err != nil {
		t.Fatalf("error skipping time by %d seconds: %s", waitSeconds, err.Error())
	}
	t.Logf("Skipped forward %d seconds (%s) to allow voting", waitSeconds, voteDelayTime)

	// Vote yay from node 1
	err = rp.CreateAndWaitForTransaction(func() (*core.TransactionInfo, error) {
		return prop.VoteOn(true, odao1.Transactor)
	}, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error voting on proposal: %s", err.Error())
	}
	t.Log("Voted yes from node 1")

	// Vote yay from node 2 and execute it
	err = rp.BatchCreateAndWaitForTransactions([]func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return prop.VoteOn(true, odao2.Transactor)
		},
		func() (*core.TransactionInfo, error) {
			return prop.Execute(odao2.Transactor)
		},
	}, false, odao2.Transactor)
	if err != nil {
		t.Fatalf("error voting on and executing proposal: %s", err.Error())
	}
	t.Log("Voted yes from node 2 and executed")

	// Get new values and make sure they match
	err = rp.Query(func(mc *batch.MultiCaller) error {
		odaoMgr.Settings.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error querying all updated details: %s", err.Error())
	}
	settings_test.EnsureSameDetails(t.Fatalf, &settings, odaoMgr.Settings.OracleDaoSettingsDetails)
	t.Log("New settings match expected settings - proposal succeeded")
}
