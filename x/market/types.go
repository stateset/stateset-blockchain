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

// market represents the state of a market on Stateset
type market struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
}

// market is a slice of markets
type markets []market

// Newmarket creates a new Marktplace
func Newmarket(id, name, description string, createdTime time.Time) market {
	return market{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedTime: createdTime,
	}
}

func (c market) String() string {
	return fmt.Sprintf(`market:
   ID: 			    %s
   Name: 			%s
   Description:  	%s
   CreatedTime: 	%s`,
		c.ID, c.Name, c.Description, c.CreatedTime.String())
}