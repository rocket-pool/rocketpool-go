package proposals_test

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
)

func Test_ProposeChallengeCooldown(t *testing.T) {
	newVal := tests.ODaoDefaults.Member.ChallengeCooldown.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.ChallengeCooldown.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeCooldown.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeChallengeCost(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Member.ChallengeCost.Get(), eth.EthToWei(1))
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.ChallengeCost.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeCost.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeChallengeWindow(t *testing.T) {
	newVal := tests.ODaoDefaults.Member.ChallengeWindow.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.ChallengeWindow.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.ChallengeWindow.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeQuorum(t *testing.T) {
	newVal := tests.ODaoDefaults.Member.Quorum.Formatted() + 0.15
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.Quorum.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.Quorum.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeRplBond(t *testing.T) {
	newVal := big.NewInt(0).Add(tests.ODaoDefaults.Member.RplBond.Get(), eth.EthToWei(1))
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Member.RplBond.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Member.RplBond.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposeBondReductionCancellationQuorum(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.BondReductionCancellationQuorum.Formatted() + 0.15
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.BondReductionCancellationQuorum.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionCancellationQuorum.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeBondReductionWindowLength(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.BondReductionWindowLength.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.BondReductionWindowLength.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionWindowLength.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeBondReductionWindowStart(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.BondReductionWindowStart.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.BondReductionWindowStart.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.BondReductionWindowStart.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeScrubPenaltyEnabled(t *testing.T) {
	newVal := !tests.ODaoDefaults.Minipool.IsScrubPenaltyEnabled.Get()
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.IsScrubPenaltyEnabled.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.IsScrubPenaltyEnabled.ProposeSet(newVal, odao1.Transactor)
	})
}

func Test_ProposePromotionScrubPeriod(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.PromotionScrubPeriod.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.PromotionScrubPeriod.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.PromotionScrubPeriod.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeScrubPeriod(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.ScrubPeriod.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.ScrubPeriod.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.ScrubPeriod.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeScrubQuorum(t *testing.T) {
	newVal := tests.ODaoDefaults.Minipool.ScrubQuorum.Formatted() + 0.15
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Minipool.ScrubQuorum.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Minipool.ScrubQuorum.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeProposalActionTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.ActionTime.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.ActionTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.ActionTime.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeProposalCooldownTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.CooldownTime.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.CooldownTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.CooldownTime.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeProposalExecuteTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.ExecuteTime.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.ExecuteTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.ExecuteTime.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeVoteDelayTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.VoteDelayTime.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.VoteDelayTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.VoteDelayTime.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func Test_ProposeVoteTime(t *testing.T) {
	newVal := tests.ODaoDefaults.Proposal.VoteTime.Formatted() + time.Hour
	testOdaoParameterProposal(t, func(newSettings *oracle.OracleDaoSettings) {
		newSettings.Proposal.VoteTime.Set(newVal)
	}, func() (*core.TransactionInfo, error) {
		return odaoMgr.Settings.Proposal.VoteTime.ProposeSet(core.GetValueForUint256(newVal), odao1.Transactor)
	})
}

func testOdaoParameterProposal(t *testing.T, setter func(*oracle.OracleDaoSettings), proposer func() (*core.TransactionInfo, error)) {
	// Revert to the initialized state at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToInitialized()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to initialized snapshot: %w", err))
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

	// Make sure there aren't any proposals
	err = rp.Query(nil, nil, dpm.ProposalCount)
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
		core.QueryAllFields(odaoMgr.Settings, mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error querying all updated details: %s", err)
	}
	settings_test.EnsureSameDetails(t.Fatalf, &tests.ODaoDefaults, odaoMgr.Settings)
	t.Log("Settings match the defaults after proposal creation, ok")

	// Make sure the proposal count is good
	err = rp.Query(nil, nil, dpm.ProposalCount)
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
	t.Logf("Got prop with ID %d: %s", prop.ID, prop.Message.Get())

	// Skip enough time to allow voting
	voteDelayTime := settings.Proposal.VoteDelayTime.Formatted()
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
	err = rp.BatchCreateAndWaitForTransactions([]func() (*core.TransactionSubmission, error){
		func() (*core.TransactionSubmission, error) {
			return core.CreateTxSubmissionFromInfo(prop.VoteOn(true, odao2.Transactor))
		},
		func() (*core.TransactionSubmission, error) {
			return core.CreateTxSubmissionFromInfo(prop.Execute(odao2.Transactor))
		},
	}, false, odao2.Transactor)
	if err != nil {
		t.Fatalf("error voting on and executing proposal: %s", err.Error())
	}
	t.Log("Voted yes from node 2 and executed")

	// Get new values and make sure they match
	err = rp.Query(func(mc *batch.MultiCaller) error {
		core.QueryAllFields(odaoMgr.Settings, mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error querying all updated details: %s", err.Error())
	}
	settings_test.EnsureSameDetails(t.Fatalf, &settings, odaoMgr.Settings)
	t.Log("New settings match expected settings - proposal succeeded")
}
