package dao_test

import (
	"fmt"
	"testing"

	batchquery "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/tests"
)

func Test_ChallengeAndKick(t *testing.T) {
	account, _ := prepChallenge(t)

	// Wait the challenge period
	secondsToWait := int(odao.Details.Members.ChallengeWindow.Formatted().Seconds())
	err := mgr.IncreaseTime(secondsToWait)
	if err != nil {
		t.Fatalf("error waiting %s for challenge window: %s", odao.Details.Members.ChallengeWindow.Formatted(), err.Error())
	}
	t.Logf("Time increased by %s", odao.Details.Members.ChallengeWindow.Formatted())

	// Decide it
	err = rp.CreateAndWaitForTransaction(func() (*core.TransactionInfo, error) {
		return dnta.DecideChallenge(account.Address, odao1.Transactor)
	}, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error deciding challenge: %s", err.Error())
	}
	t.Logf("Challenge completed")

	// Get the oDAO member count
	err = rp.Query(func(mc *batchquery.MultiCaller) error {
		dnt.GetMemberCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting oDAO member count: %s", err.Error())
	}
	count := dnt.Details.MemberCount.Formatted()
	t.Logf("oDAO now has %d members", count)

	// Make sure the node isn't on the oDAO anymore
	addresses, err := dnt.GetMemberAddresses(count, nil)
	if err != nil {
		t.Fatalf("error getting oDAO member addresses: %s", err.Error())
	}
	for i, address := range addresses {
		if address == account.Address {
			t.Fatalf("node %s was still found in the oDAO addresses with index %d", account.Address.Hex(), i)
		}
	}
	t.Logf("Node %s is no longer on the oDAO", account.Address.Hex())
}

func Test_ChallengeResolve(t *testing.T) {
	account, member := prepChallenge(t)

	// Respond with the new account
	err := rp.CreateAndWaitForTransaction(func() (*core.TransactionInfo, error) {
		return dnta.DecideChallenge(account.Address, account.Transactor)
	}, true, account.Transactor)
	if err != nil {
		t.Fatalf("error deciding challenge: %s", err.Error())
	}
	t.Logf("Challenge completed")

	// Get the oDAO member count
	err = rp.Query(func(mc *batchquery.MultiCaller) error {
		dnt.GetMemberCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting oDAO member count: %s", err.Error())
	}
	count := dnt.Details.MemberCount.Formatted()
	t.Logf("oDAO now has %d members", count)

	// Make sure the node is still on the oDAO
	addresses, err := dnt.GetMemberAddresses(count, nil)
	if err != nil {
		t.Fatalf("error getting oDAO member addresses: %s", err.Error())
	}
	found := false
	for i, address := range addresses {
		if address == account.Address {
			t.Logf("Found member %s with index %d on the oDAO", account.Address.Hex(), i)
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("member %s was not found in the oDAO addresses", account.Address.Hex())
	}

	// Query some state
	err = rp.Query(func(mc *batchquery.MultiCaller) error {
		member.GetIsChallenged(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting contract state: %s", err.Error())
	}

	// Make sure it's not challenged anymore
	if member.Details.IsChallenged {
		t.Fatalf("member is challenged, but should not be")
	}
	t.Logf("Challenge resolved!")
}

func prepChallenge(t *testing.T) (*tests.Account, *trustednode.OracleDaoMember) {
	// Revert to the initialized state at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToInitialized()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to initialized snapshot: %w", err))
		}
	})

	// Get a 4th member to join the oDAO
	account := mgr.NonOwnerAccounts[3]
	_, err := tests.BootstrapNodeToOdao(rp, mgr.OwnerAccount, account, "Etc/UTC", "The Sacrifice", "rocketpool.net")
	if err != nil {
		t.Fatalf("error bootstrapping node %s to the oDAO: %s", account.Address.Hex(), err.Error())
	}
	member, err := trustednode.NewOracleDaoMember(rp, account.Address)
	if err != nil {
		t.Fatalf("error creating oDAO member binding for node %s: %s", account.Address.Hex(), err.Error())
	}
	t.Logf("Bootstrapped node %s onto the oDAO", account.Address.Hex())

	// Verify it's on the oDAO
	err = rp.Query(func(mc *batchquery.MultiCaller) error {
		dnt.GetMemberCount(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting oDAO member count: %s", err.Error())
	}
	count := dnt.Details.MemberCount.Formatted()
	t.Logf("oDAO now has %d members", count)

	// Find it
	addresses, err := dnt.GetMemberAddresses(count, nil)
	if err != nil {
		t.Fatalf("error getting member addresses: %s", err.Error())
	}
	found := false
	for i, address := range addresses {
		if address == account.Address {
			t.Logf("Found member %s with index %d on the oDAO", account.Address.Hex(), i)
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("member %s was not found in the oDAO addresses", account.Address.Hex())
	}

	// Issue a challenge to it
	err = rp.CreateAndWaitForTransaction(func() (*core.TransactionInfo, error) {
		return dnta.MakeChallenge(account.Address, odao1.Transactor)
	}, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error challenging member %s: %s", account.Address.Hex(), err.Error())
	}
	t.Logf("Challenge issued")

	// Query some state
	err = rp.Query(func(mc *batchquery.MultiCaller) error {
		member.GetIsChallenged(mc)
		odao.GetChallengeWindow(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("error getting contract state: %s", err.Error())
	}

	// Make sure the challenge is visible
	if !member.Details.IsChallenged {
		t.Fatalf("member is not challenged, but should be")
	}
	t.Logf("Challenge is visible")

	return account, member
}
