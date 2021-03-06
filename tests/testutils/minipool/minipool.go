package minipool

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/eth"

	"github.com/rocket-pool/rocketpool-go/tests/testutils/accounts"
	nodeutils "github.com/rocket-pool/rocketpool-go/tests/testutils/node"
	"github.com/rocket-pool/rocketpool-go/tests/testutils/validator"
)

// Minipool created event
type minipoolCreated struct {
    Minipool common.Address
    Node common.Address
    Time *big.Int
}


// Create a minipool
func CreateMinipool(rp *rocketpool.RocketPool, ownerAccount, nodeAccount *accounts.Account, depositAmount *big.Int) (*minipool.Minipool, error) {

    // Mint & stake RPL required for mininpool
    rplRequired, err := GetMinipoolRPLRequired(rp)
    if err != nil { return nil, err }
    if err := nodeutils.StakeRPL(rp, ownerAccount, nodeAccount, rplRequired); err != nil { return nil, err }

    // Make node deposit
    opts := nodeAccount.GetTransactor()
    opts.Value = depositAmount
    hash, err := node.Deposit(rp, 0, opts)
    if err != nil { return nil, err }
    txReceipt, err := utils.WaitForTransaction(rp.Client, hash)
    if err != nil { return nil, err }

    // Get minipool manager contract
    rocketMinipoolManager, err := rp.GetContract("rocketMinipoolManager")
    if err != nil { return nil, err }

    // Get created minipool address
    minipoolCreatedEvents, err := rocketMinipoolManager.GetTransactionEvents(txReceipt, "MinipoolCreated", minipoolCreated{})
    if err != nil || len(minipoolCreatedEvents) == 0 {
        return nil, errors.New("Could not get minipool created event")
    }
    minipoolAddress := minipoolCreatedEvents[0].(minipoolCreated).Minipool

    // Return minipool instance
    return minipool.NewMinipool(rp, minipoolAddress)

}


// Stake a minipool
func StakeMinipool(rp *rocketpool.RocketPool, mp *minipool.Minipool, nodeAccount *accounts.Account) error {

    // Get validator & deposit data
    validatorPubkey, err := validator.GetValidatorPubkey()
    if err != nil { return err }
    withdrawalCredentials, err := mp.GetWithdrawalCredentials(nil)
    if err != nil { return err }
    validatorSignature, err := validator.GetValidatorSignature()
    if err != nil { return err }
    depositDataRoot, err := validator.GetDepositDataRoot(validatorPubkey, withdrawalCredentials, validatorSignature)
    if err != nil { return err }

    // Stake minipool & return
    _, err = mp.Stake(validatorPubkey, validatorSignature, depositDataRoot, nodeAccount.GetTransactor())
    return err

}


// Get the RPL required per minipool
func GetMinipoolRPLRequired(rp *rocketpool.RocketPool) (*big.Int, error) {

    // Get data
    depositUserAmount, err := protocol.GetMinipoolHalfDepositUserAmount(rp, nil)
    if err != nil { return nil, err }
    minimumPerMinipoolStake, err := protocol.GetMinimumPerMinipoolStake(rp, nil)
    if err != nil { return nil, err }
    rplPrice, err := network.GetRPLPrice(rp, nil)
    if err != nil { return nil, err }

    // Calculate and return RPL required
    var tmp big.Int
    var rplRequired big.Int
    tmp.Mul(depositUserAmount, eth.EthToWei(minimumPerMinipoolStake))
    rplRequired.Quo(&tmp, rplPrice)
    return &rplRequired, nil

}

