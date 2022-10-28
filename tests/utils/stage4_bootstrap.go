package utils

import (
	"log"
	"time"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	"github.com/rocket-pool/rocketpool-go/tests/testutils/accounts"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// Bootstrap all of the parameters to mimic Stage 4 so the unit tests work correctly
func Stage4Bootstrap(rp *rocketpool.RocketPool, ownerAccount *accounts.Account) {

	opts := ownerAccount.GetTransactor()

	_, err := protocol.BootstrapDepositEnabled(rp, true, opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapAssignDepositsEnabled(rp, true, opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapMaximumDepositPoolSize(rp, eth.EthToWei(1000), opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapNodeRegistrationEnabled(rp, true, opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapNodeDepositEnabled(rp, true, opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapMinipoolSubmitWithdrawableEnabled(rp, true, opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapMinimumNodeFee(rp, 0.05, opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapTargetNodeFee(rp, 0.1, opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapMaximumNodeFee(rp, 0.2, opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapNodeFeeDemandRange(rp, eth.EthToWei(1000), opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = protocol.BootstrapInflationStartTime(rp,
		uint64(time.Now().Unix()+(60*60*24*14)), opts)
	if err != nil {
		log.Fatal(err)
	}

}
