package types

import (
	"fmt"
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	StoreKey          = ModuleName
	DefaultParamspace = ModuleName
)

// PurchaseOrder stores data about an purchase order
type PurchaseOrder struct {
	purchaseOrderID         uint64         `json:"purchaseOrderId"`
	purchaseOrderNumber     string         `json:"purchaseOrderNumber`
	purchaseOrderName		  string		 `json:"purchaseOrderName"`
	description     string		 `json:"description"`
	subtotal	      sdk.Coin		 `json:"subtotal"`
	total			  sdk.Coin 	     `json:"total"`
	purchaser			  sdk.AccAddress `json:"purchaser"`
	vendor      sdk.AccAddress `json:"vendor"`
	financer      sdk.AccAddress `json:"financer"`
	purchaseDate			  time.Time 	 `json:"purchaseDate"`
	deliveryDate			  time.Time 	 `json:"deliveryDate"`
	paid			  bool			 `json:"paid"`
	active 	          bool           `json:"active"`
	CreatedTime       time.Time      `json:"created_time"`
}

// Purchase Orders is an array of purchaseorders
type PurchaseOrders []PurchaseOrder

// NewPurchaseOrder creates a new purchase order object
func NewPurchaseOrder(purchaseOrderId uint64, purchaseOrderNumber string, purchaseOrderName string, description string, subtotal sdk.Coin, total sdk.Coin, purchaser sdk.AccAddress, vendor sdk.AccAddress, financer sdk.AccAddress, purchaseDate time.Time, deliveryDate time.Time, periodEndDate time.Time, paid bool, active bool, createdTime time.Time) PurchaseOrder {
	return PurchaseOrder{
		PurchaseOrderID:       purchaseorderId,
		PurchaseOrderNumber:   purchaseorderNumber,
		PurchaseOrderName:     purchaseorderName,
		Description:   description,
		Subtotal:        subtotal,
		Total: 			 total,
		Purchaser:		     purchaser,
		Vendor:	 vendor,
		Financer:	 financer,
		PurchaseDate: purchaseDate,
		DeliveryDate:   deliveryDate,
		Paid:			 paid,
		Active: 		 active,
		CreatedTime:     createdTime,
	}
}

func (i PurchaseOrder) String() string {
	return fmt.Sprintf(`PurchaseOrder %d:
		PurchaseOrderID:       %s,
		PurchaseOrderNumber:   %s,
		PurchaseOrderName:     %s,
		Description:   %s,
		Subtotal:        %s,
		Total: 			 %s,
		Purchaser:		     %s,
		Vendor:	 %s,
		Financer:	 %s,
		PurchaseDate: %s,
		DeliveryDate:   %s,
		Paid:			 %s,
		Active: 		 %s,
		CreatedTime:     %s`,
		i.PurchaseOrderID, i.PurchaseOrderNumber, i.PurchaseOrderName, i.Description, i.Subtotal.String(), i.Total.String(), i.Purchaser.String(), i.Vendor.String(), i.Financer.String(), i.PurchaseDate.String(), i.DeliveryDate.String(), i.Paid.String(), i.Active.String(), a.CreatedTime.String())
}

type deletePurchaseorderRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Creator string `json:"creator"`
}