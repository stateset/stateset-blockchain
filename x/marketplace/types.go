package marketplace


import (
	"fmt"
	"time"
)

// Defines module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName
)

// Marketplace represents the state of a marketplace on Stateset
type Marketplace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
}

// Marketplace is a slice of marketplaces
type Marketplaces []Marketplace

// NewMarketplace creates a new Marktplace
func NewMarketplace(id, name, description string, createdTime time.Time) Marketplace {
	return Marketplace{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedTime: createdTime,
	}
}

func (c Marketplace) String() string {
	return fmt.Sprintf(`Marketplace:
   ID: 			    %s
   Name: 			%s
   Description:  	%s
   CreatedTime: 	%s`,
		c.ID, c.Name, c.Description, c.CreatedTime.String())
}