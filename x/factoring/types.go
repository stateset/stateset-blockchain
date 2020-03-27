package staking

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stateset/stateset-blockchain/x/bank"
)

// Defines staking module constants
const (
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName

	EventTypeInterestRewardPaid = "interest-reward-paid"
	AttributeKeyExpiredStakes   = "expired-stakes"

	EventTypeFactorLimitIncreased  = "factor-limit-increased"
	AttributeKeyFactorLimitUpgrade = "factor-limit-upgrade"

	UserFactorsPoolName = "user_factors_tokens_pool"
)

type FactorType byte

func (t FactorType) String() string {
	if int(t) >= len(FactorTypeName) {
		return "Unknown"
	}
	return FactorTypeName[t]
}

const (
	FactorBacking FactorType = iota
	FactorChallenge
)

var FactorTypeName = []string{
	FactorBacking:   "FactorBacking",
	FactorChallenge: "FactorChallenge"
}

var bankTransactionMappings = []TransactionType{
	FactorBacking:   TransactionBacking,
	FactorChallenge: TransactionChallenge
}

func (t FactorType) BankTransactionType() bank.TransactionType {
	if int(t) >= len(bankTransactionMappings) {
		panic("invalid factor type")
	}
	return bankTransactionMappings[t]
}

func (t FactorType) ValidForLoan() bool {
	return t.oneOf([]FactorType{FactorBacking, FactorChallenge})
}

func (t FactorType) ValidForUpvote() bool {
	return t.oneOf([]FactorType{FactorBacking, FactorChallenge})
}

func (t FactorType) Valid() bool {
	return t.oneOf([]FactorType{FactorBacking, FactorChallenge})
}

func (t FactorType) oneOf(types []FactorType) bool {
	for _, tType := range types {
		if tType == t {
			return true
		}
	}
	return false
}

type Factor struct {
	ID          uint64         `json:"id"`
	InvoiceID   uint64         `json:"invoice_id"`
	LoanID      uint64         `json:"loan_id"`
	MarketplaceID string       `json:"marketplace_id"`
	Type        FactorType      `json:"type"`
	Amount      sdk.Coin       `json:"amount"`
	Factor      sdk.AccAddress `json:"factor"`
	CreatedTime time.Time      `json:"created_time"`
	EndTime     time.Time      `json:"end_time"`
	Expired     bool           `json:"expired"`
	Result      *RewardResult  `json:"result,omitempty"`
}

func (f Factor) String() string {
	return fmt.Sprintf(`Factor %d:
  InvoiceID: %d
  LoanID: %d
  Amount: %s
  Factor: %s`,
		f.ID, f.InvoiceID, f.LoanID, f.Amount.String(), s.Factor.String())
}

type Loans struct {
	ID                uint64         `json:"id"`
	Lender            sdk.AccAddress `json:"lender"`
	InvoiceID         uint64         `json:"invoice_id"`
	MarketplaceID     string         `json:"marketplace_id"`
	FactorType        FactorType     `json:"factor_type"`
	TotalFactor       sdk.Coin       `json:"total_factor"`
	LoanNumber        string         `json:"loanNumber`
	LoanName		  string		 `json:"loanName"`
	Description    	  string		 `json:"description"`
	LoanAmount        sdk.Coin       `json:"loanAmount"`
	PaidAmount 		  sdk.Coin		 `json:"paidAmount"`
	AmountRemaining   sdk.Coin       `json:"amountRemaining"`
	Subtotal	      sdk.Coin		 `json:"subtotal"`
	Total			  sdk.Coin 	     `json:"total"`
	DueDate			  time.Time 	 `json:"dueDate"`
	Paid			  bool			 `json:"paid"`
	Active 	          bool           `json:"active"`
	CreatedTime       time.Time      `json:"created_time"`
	CreatedTime       time.Time      `json:"created_time"`
	UpdatedTime       time.Time      `json:"updated_time"`
	EditedTime        time.Time      `json:"edited_time"`
	Edited            bool           `json:"edited"`
}

type FactorLimitUpgrade struct {
	Address     sdk.AccAddress `json:"address"`
	NewLimit    int            `json:"new_limit"`
	Earned      sdk.Coin       `json:"earned"`
}