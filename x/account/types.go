package account

import (
	"fmt"
	"net/url"
	"time"

	app "github.com/stateset/stateset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	StoreKey          = ModuleName
	DefaultParamspace = ModuleName
)

// Account stores data about an account
type Account struct {
	accountID         uint64         `json:"accountId"`
	accountName		  string		 `json:"accountName"`
	accountType       string         `json:"accountType"`
	industry		  string		 `json:"industry"`
	phone			  string 		 `json:"industry"`
	controller		  sdk.AccAddress `json:"controller"`
	processor		  sdk.AccAddress `json:"processor"`
	CreatedTime       time.Time      `json:"created_time"`
}

// Accounts is an array of accounts
type Accounts []Account

// NewAccount creates a new account object
func NewAccount(accountId uint64, accountName string, accountType string, industry: string, phone: string, controller: sdk.AccAddress, processor: sdk.AccAddress, createdTime time.Time) Account {
	return Account{
		AccountID:       accountId,
		AccountName:     accountName,
		AccountType:	 accountType,
		Industry:	     industry,
		Phone: 			 phone,
		Controller:	     controller,
		Processor: 	     processor,
		CreatedTime:     createdTime,
	}
}

func (a Account) String() string {
	return fmt.Sprintf(`Account %d:
  AccountID:    %s
  AccountName:	%s
  AccountType:  %s
  Industry:     %s
  Phone:		%s
  Controller:	%s
  Processor:    %s`,
		a.AccountID, a.AccountName, a.AccountType, a.Industry, a.Phone, a.Controller.String(), a.Processor.String(), a.CreatedTime.String())
}