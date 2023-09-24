package bootstrap_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/tests"
	settings_test "github.com/rocket-pool/rocketpool-go/tests/settings"
)

var (
	mgr     *tests.TestManager
	rp      *rocketpool.RocketPool
	pdaoMgr *protocol.ProtocolDaoManager
	odaoMgr *oracle.OracleDaoManager
	opts    *bind.TransactOpts
)

func TestMain(m *testing.M) {
	// Make the test manager
	var err error
	mgr, err = tests.NewTestManager()
	if err != nil {
		log.Fatal(fmt.Sprintf("error getting test manager: %s", err.Error()))
	}
	rp = mgr.RocketPool

	// Make the pDAO / oDAO bindings
	pdaoMgr, err = protocol.NewProtocolDaoManager(rp)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating pdao manager binding: %w", err))
	}
	odaoMgr, err = oracle.NewOracleDaoManager(rp)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating odao manager binding: %w", err))
	}

	// Create the default values for the pDAO / oDAO settings as a reference
	err = tests.CreateDefaults(mgr)
	if err != nil {
		log.Fatal("error creating defaults: %w", err)
	}

	// Use the owner account for bootstrapping things
	opts = mgr.OwnerAccount.Transactor

	// Get all of the current settings
	err = rp.Query(func(mc *batch.MultiCaller) error {
		core.QueryAllFields(odaoMgr.Settings, mc)
		core.QueryAllFields(pdaoMgr.Settings, mc)
		return nil
	}, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("error querying all initial details: %w", err))
	}

	// Verify details
	settings_test.EnsureSameDetails(log.Fatalf, &tests.ODaoDefaults, odaoMgr.Settings)
	settings_test.EnsureSameDetails(log.Fatalf, &tests.PDaoDefaults, pdaoMgr.Settings)
	log.Printf("Stock details are correct!")
	log.Println("Defaults:")
	printOdao(&tests.ODaoDefaults)
	log.Println()
	log.Println("Chain:")
	printOdao(odaoMgr.Settings)

	// Run tests
	os.Exit(m.Run())
}

func printOdao(settings *oracle.OracleDaoSettings) {
	log.Println("Member:")
	log.Printf("\tQuorum: %s\n", settings.Member.Quorum.Raw().String())
	log.Printf("\tRplBond: %s\n", settings.Member.RplBond.Get().String())
	log.Printf("\tUnbondedMinipoolMax: %s\n", settings.Member.UnbondedMinipoolMax.Raw().String())
	log.Printf("\tUnbondedMinipoolMinFee: %s\n", settings.Member.UnbondedMinipoolMinFee.Raw().String())
	log.Printf("\tChallengeCooldown: %s\n", settings.Member.ChallengeCooldown.Raw().String())
	log.Printf("\tChallengeWindow: %s\n", settings.Member.ChallengeWindow.Raw().String())
	log.Printf("\tChallengeCost: %s\n", settings.Member.ChallengeCost.Get().String())
	log.Println()
	log.Println("Minipool:")
	log.Printf("\tScrubPeriod: %s\n", settings.Minipool.ScrubPeriod.Raw().String())
	log.Printf("\tScrubQuorum: %s\n", settings.Minipool.ScrubQuorum.Raw().String())
	log.Printf("\tPromotionScrubPeriod: %s\n", settings.Minipool.PromotionScrubPeriod.Raw().String())
	log.Printf("\tIsScrubPenaltyEnabled: %t\n", settings.Minipool.IsScrubPenaltyEnabled.Get())
	log.Printf("\tBondReductionWindowStart: %s\n", settings.Minipool.BondReductionWindowStart.Raw().String())
	log.Printf("\tBondReductionWindowLength: %s\n", settings.Minipool.BondReductionWindowLength.Raw().String())
	log.Printf("\tBondReductionCancellationQuorum: %s\n", settings.Minipool.BondReductionCancellationQuorum.Raw().String())
	log.Println()
	log.Println("Proposal:")
	log.Printf("\tCooldownTime: %s\n", settings.Proposal.CooldownTime.Raw().String())
	log.Printf("\tVoteTime: %s\n", settings.Proposal.VoteTime.Raw().String())
	log.Printf("\tVoteDelayTime: %s\n", settings.Proposal.VoteDelayTime.Raw().String())
	log.Printf("\tExecuteTime: %s\n", settings.Proposal.ExecuteTime.Raw().String())
	log.Printf("\tActionTime: %s\n", settings.Proposal.ActionTime.Raw().String())
}
