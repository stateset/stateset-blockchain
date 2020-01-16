package contact

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

// Contact stores data about an contact
type Contact struct {
	contactID         uint64         `json:"contactId"`
	firstName		  string		 `json:"firstName"`
	lastName		  string		 `json:"lastName"`
	email	 		  string 		 `json:"email"`
	phone 			  string 		 `json:"phone"`
	controller		  sdk.AccAddress `json:"controller"`
	processor		  sdk.AccAddress `json:"processor"`
	CreatedTime       time.Time      `json:"created_time"`
}

// Contacts is an array of contacts
type Contacts []Contact

// NewContact creates a new contact object
func NewContact(contactId uint64, firstName string, lastName string, email string, phone string, controller sdk.AccAddress, processor sdk.AccAddress, createdTime time.Time) Contact {
	return Contact{
		ContactID:       contactId,
		FirstName:		 firstName,
		LastName:		 lastName,
		Email: 			 email,
		Phone:			 phone,
		Controller:	     controller,
		Processor: 	     processor,
		CreatedTime:     createdTime,
	}
}

func (c Contact) String() string {
	return fmt.Sprintf(`Contact %d:
  ContactID:    %s
  FirstName:	%s
  LastName:     %s
  Email:  		%s
  Phone:		%s
  Controller:	%s
  Processor:    %s`,
		c.ContactID, c.FirstName, c.LastName, c.Email, c.Phone, c.Controller.String(), c.Processor.String(), a.CreatedTime.String())
}