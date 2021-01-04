package market


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

// Market represents the state of a market on Stateset
type Market struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
}

// Market is a slice of markets
type Markets []Market

// NewMarket creates a new Markets
func NewMarket(id, name, description string, createdTime time.Time) Market {
	return Market{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedTime: createdTime,
	}
}

func (c Market) String() string {
	return fmt.Sprintf(`Market:
   ID: 			    %s
   Name: 			%s
   Description:  	%s
   CreatedTime: 	%s`,
		c.ID, c.Name, c.Description, c.CreatedTime.String())
}