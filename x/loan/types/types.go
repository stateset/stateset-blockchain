package loan

import (
	"fmt"
	"net/url"
	"time"

	app "github.com/stateset/stateset-blockchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	StoreKey          = ModuleName
	DefaultParamspace = ModuleName
)

// Loan stores data about a loan
type Loan struct {
	Amount	     	  uint64		 `json:"amount"`
	Fee     		  uint64 		 `json:"fee"`   
	Collateral 		  uint64 		 `json:"collateral"`
	Deadline 		  uint64 	     `json:"deadline"`
	State      	  	  string 		 `json:"state"` 
	Borrower   	  	  sdk.Address 	 `json:"borrower"`
}
