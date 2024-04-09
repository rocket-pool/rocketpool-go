package dao_test

import (
	"fmt"
	"testing"

	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/v2/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/v2/tests"
)

func Test_ChallengeAndKick(t *testing.T) {
	account, _ := prepChallenge(t)

	// Wait the challenge period
	secondsToWait := int(odaoMgr.Settings.Member.ChallengeWindow.Formatted().Seconds())
	err := mgr.IncreaseTime(secondsToWait)
	if err != nil {
		t.Fatalf("error waiting %s for challenge window: %s", odaoMgr.Settings.Member.ChallengeWindow.Formatted(), err.Error())
	}
	t.Logf("Time increased by %s", odaoMgr.Settings.Member.ChallengeWindow.Formatted())

	// Decide it
	err = rp.CreateAndWaitForTransaction(func() (*eth.TransactionInfo, error) {
		return odaoMgr.DecideChallenge(account.Address, odao1.Transactor)
	}, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error deciding challenge: %s", err.Error())
	}
	t.Logf("Challenge completed")

	// Get the oDAO member count
	err = rp.Query(nil, nil, odaoMgr.MemberCount)
	if err != nil {
		t.Fatalf("error getting oDAO member count: %s", err.Error())
	}
	count := odaoMgr.MemberCount.Formatted()
	t.Logf("oDAO now has %d members", count)

	// Make sure the node isn't on the oDAO anymore
	addresses, err := odaoMgr.GetMemberAddresses(count, nil)
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
	err := rp.CreateAndWaitForTransaction(func() (*eth.TransactionInfo, error) {
		return odaoMgr.DecideChallenge(account.Address, account.Transactor)
	}, true, account.Transactor)
	if err != nil {
		t.Fatalf("error deciding challenge: %s", err.Error())
	}
	t.Logf("Challenge completed")

	// Get the oDAO member count
	err = rp.Query(nil, nil, odaoMgr.MemberCount)
	if err != nil {
		t.Fatalf("error getting oDAO member count: %s", err.Error())
	}
	count := odaoMgr.MemberCount.Formatted()
	t.Logf("oDAO now has %d members", count)

	// Make sure the node is still on the oDAO
	addresses, err := odaoMgr.GetMemberAddresses(count, nil)
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
	err = rp.Query(nil, nil, member.IsChallenged)
	if err != nil {
		t.Fatalf("error getting contract state: %s", err.Error())
	}

	// Make sure it's not challenged anymore
	if member.IsChallenged.Get() {
		t.Fatalf("member is challenged, but should not be")
	}
	t.Logf("Challenge resolved!")
}

func prepChallenge(t *testing.T) (*tests.Account, *oracle.OracleDaoMember) {
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
	member, err := oracle.NewOracleDaoMember(rp, account.Address)
	if err != nil {
		t.Fatalf("error creating oDAO member binding for node %s: %s", account.Address.Hex(), err.Error())
	}
	t.Logf("Bootstrapped node %s onto the oDAO", account.Address.Hex())

	// Verify it's on the oDAO
	err = rp.Query(nil, nil, odaoMgr.MemberCount)
	if err != nil {
		t.Fatalf("error getting oDAO member count: %s", err.Error())
	}
	count := odaoMgr.MemberCount.Formatted()
	t.Logf("oDAO now has %d members", count)

	// Find it
	addresses, err := odaoMgr.GetMemberAddresses(count, nil)
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
	err = rp.CreateAndWaitForTransaction(func() (*eth.TransactionInfo, error) {
		return odaoMgr.MakeChallenge(account.Address, odao1.Transactor)
	}, true, odao1.Transactor)
	if err != nil {
		t.Fatalf("error challenging member %s: %s", account.Address.Hex(), err.Error())
	}
	t.Logf("Challenge issued")

	// Query some state
	err = rp.Query(nil, nil, member.IsChallenged, odaoMgr.Settings.Member.ChallengeWindow)
	if err != nil {
		t.Fatalf("error getting contract state: %s", err.Error())
	}

	// Make sure the challenge is visible
	if !member.IsChallenged.Get() {
		t.Fatalf("member is not challenged, but should be")
	}
	t.Logf("Challenge is visible")

	return account, member
}
