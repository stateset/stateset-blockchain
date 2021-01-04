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
	DefaultParamspace = ModuleName
)

// Market represents the state of a Market on Stateset
type Market struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
}

// Market is a slice of markets
type Markets []Market

// NewMarket creates a new Markets
func NewMarket(id uint64, name string, description string, createdTime time.Time) Market {
	return Market{
		MarketID:    marketId,
		Name:        name,
		Description: description,
		CreatedTime: createdTime,
	}
}

func (c Market) String() string {
	return fmt.Sprintf(`Market:
   MarketID: 			    %s
   Name: 			%s
   Description:  	%s
   CreatedTime: 	%s`,
		c.MarketID, c.Name, c.Description, c.CreatedTime.String())
}