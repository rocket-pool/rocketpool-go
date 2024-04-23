package v10_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
	"github.com/rocket-pool/rocketpool-go/v2/tests"
)

var (
	mgr *tests.TestManager
	rp  *rocketpool.RocketPool

	// oDAO accounts
	odao1 *tests.Account
	odao2 *tests.Account
	odao3 *tests.Account

	// Node accounts
	node1 *tests.Account
	node2 *tests.Account
	node3 *tests.Account
	node4 *tests.Account
	node5 *tests.Account
	node6 *tests.Account
	node7 *tests.Account
	node8 *tests.Account

	// Non node accounts (withdrawal addresses)
	node2Primary  *tests.Account
	node3Rpl      *tests.Account
	node4Primary  *tests.Account
	node4Rpl      *tests.Account
	multiReceiver *tests.Account

	accountNames map[common.Address]string
)

func TestMain(m *testing.M) {
	// Make the test manager
	var err error
	mgr, err = tests.NewTestManager()
	if err != nil {
		log.Fatal(fmt.Sprintf("error getting test manager: %s", err.Error()))
	}
	rp = mgr.RocketPool

	// Initialize the network
	err = mgr.InitializeDeployment()
	if err != nil {
		fail("error initializing deployment: %s", err.Error())
	}

	// Assign accounts
	odao1 = mgr.NonOwnerAccounts[0]
	odao2 = mgr.NonOwnerAccounts[1]
	odao3 = mgr.NonOwnerAccounts[2]

	node1 = mgr.NonOwnerAccounts[3]
	node2 = mgr.NonOwnerAccounts[4]
	node3 = mgr.NonOwnerAccounts[5]
	node4 = mgr.NonOwnerAccounts[6]

	node2Primary = mgr.NonOwnerAccounts[7]
	node3Rpl = mgr.NonOwnerAccounts[8]
	node4Primary = mgr.NonOwnerAccounts[9]
	node4Rpl = mgr.NonOwnerAccounts[10]

	node5 = mgr.NonOwnerAccounts[11]
	node6 = mgr.NonOwnerAccounts[12]
	node7 = mgr.NonOwnerAccounts[13]
	node8 = mgr.NonOwnerAccounts[14]
	multiReceiver = mgr.NonOwnerAccounts[15]

	accountNames = map[common.Address]string{
		odao1.Address:         "oDAO 1",
		odao2.Address:         "oDAO 2",
		odao3.Address:         "oDAO 3",
		node1.Address:         "Node 1",
		node2.Address:         "Node 2",
		node3.Address:         "Node 3",
		node4.Address:         "Node 4",
		node2Primary.Address:  "Node 2 Primary",
		node3Rpl.Address:      "Node 3 RPL",
		node4Primary.Address:  "Node 4 Primary",
		node4Rpl.Address:      "Node 4 RPL",
		node5.Address:         "Node 5",
		node6.Address:         "Node 6",
		node7.Address:         "Node 7",
		node8.Address:         "Node 8",
		multiReceiver.Address: "Multi Receiver",
	}

	// Run tests
	code := m.Run()

	// Revert to the baseline after testing is done
	cleanup()

	// Done
	os.Exit(code)
}

func fail(format string, args ...any) {
	log.Printf(format, args...)
	cleanup()
	os.Exit(1)
}

func cleanup() {
	err := mgr.RevertToBaseline()
	if err != nil {
		log.Fatalf("error reverting to baseline snapshot: %s\nPlease restart Hardhat as the state will now be corrupted for other tests", err.Error())
	}
}
