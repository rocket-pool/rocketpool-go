package proposals_test

import (
	"fmt"
	"testing"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/settings"
	"github.com/rocket-pool/rocketpool-go/tests"
	settings_test "github.com/rocket-pool/rocketpool-go/tests/settings"
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
		dp.GetProposalCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting proposal count: %s", err.Error())
	}
	if dp.Details.ProposalCount.Formatted() != 0 {
		t.Fatalf("expected 0 proposals but count was %d", dp.Details.ProposalCount.Formatted())
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
		dp.GetProposalCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting proposal count: %s", err.Error())
	}
	propCount := dp.Details.ProposalCount.Formatted()
	if propCount != 1 {
		t.Fatalf("expected 1 proposal but count was %d", propCount)
	}
	t.Logf("Prop count = %d, ok", propCount)

	// Get the proposal
	pdaoProps, odaoProps, err := dp.GetProposals(rp, nil, propCount)
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
		return dntp.VoteOnProposal(prop.Details.ID.Formatted(), true, odao1.Transactor)
	}, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error voting on proposal: %s", err.Error())
	}
	t.Log("Voted yes from node 1")

	// Vote yay from node 2 and execute it
	err = rp.BatchCreateAndWaitForTransactions([]func() (*core.TransactionInfo, error){
		func() (*core.TransactionInfo, error) {
			return dntp.VoteOnProposal(prop.Details.ID.Formatted(), true, odao2.Transactor)
		},
		func() (*core.TransactionInfo, error) {
			return dntp.ExecuteProposal(prop.Details.ID.Formatted(), odao2.Transactor)
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
