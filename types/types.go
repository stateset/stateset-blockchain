package types

import (
	stypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AppName is the name of the Cosmos app
	AppName = "Stateset"
	// StakeDenom is the name of the main staking currency
	StakeDenom = "ustate"
	// Hostname is the address the app's HTTP server will bind to
	Hostname = "0.0.0.0"
	// Portname is the port the app's HTTP server will bind to
	Portname = "1337"
)

// Coin units
const (
	PARO = 1
)

const (
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = "stateset"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = "statesetpub"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = "statesetvaloper"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = "statesetvaloperpub"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = "statesetvalcons"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = "statesetvalconspub"
)

// InitialStake is an `sdk.Coins` representing the balance a new user is granted upon registration
var InitialStake = sdk.Coin{Amount: sdk.NewInt(330 * PARO), Denom: StakeDenom}

// NewStatesetCoin returns the desired amount
func NewStatesetCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(StakeDenom, amount*PARO)
}

// MsgResult is the default success response for a chain request
type MsgResult struct {
	ID int64 `json:"id"`
}

// StakeNotificationResult defines data for a Stake push notification
type StakeNotificationResult struct {
	MsgResult
	StateID int64          `json:"state_id"`
	From    sdk.AccAddress `json:"from,omitempty"`
	To      sdk.AccAddress `json:"to,omitempty"`
	Amount  sdk.Coin       `json:"amount"`
	Cred    *sdk.Coin      `json:"cred,omitempty"`
}

// Stake represents a lender with the amount Stakeed.
type Stake struct {
	Address sdk.AccAddress
	Amount  sdk.Coin
}

// CompletedStateset defines a stateset result.
type CompletedStateset struct {
	ID                          int64                       `json:"id"`
	Merchant                    sdk.AccAddress              `json:"merchant"`
	Depositors                  []Stake                     `json:"depositors"`
	Borrowers                   []Stake                     `json:"borrowers"`
	StakeDistributionResults    StakeDistributionResults    `json:"Stake_distribution_results"`
	InterestDistributionResults InterestDistributionResults `json:"interest_distribution_results"`
}

// CompletedStatesetNotificationResult defines the notification result of completed a Stateset in a new Block.
type CompletedStatesetNotificationResult struct {
	Statesets []CompletedStateset `json:"statesets"`
}

// StakeReward represents the amount of the invoice Stakeed by a user.
type StakeReward struct {
	Account sdk.AccAddress `json:"account"`
	Amount  sdk.Coin       `json:"amount"`
}

// StakeDistributionResultsType indicates who wins the pool.
type StakeDistributionResultsType int64

// Distribution result constants
const (
	DistributionMajorityNotReached StakeDistributionResultsType = iota
)

// StakeDistributionResults contains how the Stake was distributed after the Purchase Order or Invoice is paid.
type StakeDistributionResults struct {
	Type        StakeDistributionResultsType `json:"type"`
	TotalAmount sdk.Coin                     `json:"total_amount"`
	Rewards     []StakeReward                `json:"rewards"`
}

// Interest represents the amount of interest earned by Stakeing the Purchase Order or Invoice
type Interest struct {
	Account sdk.AccAddress `json:"account"`
	Amount  sdk.Coin       `json:"amount"`
	Rate    sdk.Int        `json:"rate"`
}

// InterestDistributionResults contains how the interest was applied after a stateset completes.
type InterestDistributionResults struct {
	TotalAmount sdk.Coin   `json:"total_amount"`
	Interests   []Interest `json:"interests"`
}

func KVGasConfig() stypes.GasConfig {
	return stypes.GasConfig{
		HasCost:          0,
		DeleteCost:       0,
		ReadCostFlat:     0,
		ReadCostPerByte:  0,
		WriteCostFlat:    0,
		WriteCostPerByte: 0,
		IterNextCostFlat: 0,
	}
}
