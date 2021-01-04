package invoice

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keys for invoice store
// Items are stored with the following key: values
//
// - 0x00<invoiceID_Bytes>: Invoice_Bytes
// - 0x01: nextInvoiceID_Bytes
//
// - 0x10<marketID_Bytes><invoiceID_Bytes>: invoiceID_Bytes
// - 0x11<merchant_Bytes><invoiceID_Bytes>: invoiceID_Bytes
// - 0x12<createdTime_Bytes><invoiceID_Bytes>: invoiceID_Bytes
var (
	InvoicesKeyPrefix = []byte{0x00}
	InvoiceIDKey      = []byte{0x01}

	MaretplaceInvoicesPrefix   = []byte{0x10}
	MerchantInvoicesPrefix     = []byte{0x11}
	CreatedTimeInvoicesPrefix = []byte{0x12}
)

// key for getting a specific invoice from the store
func key(invoiceID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(invoiceID)
	return append(InvoicesKeyPrefix, bz...)
}

// marketInvoicesKey gets the first part of the market invoices key based on the marketID
func marketInvoicesKey(marketID string) []byte {
	return append(marketInvoicesPrefix, []byte(marketID)...)
}

// marketplceInvoiceKey key of a specific market <-> invoice association from the store
func marketInvoiceKey(marketID string, invoiceID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(invoiceID)
	return append(marketInvoicesKey(marketID), bz...)
}

func merchantInvoicesKey(merchant sdk.AccAddress) []byte {
	return append(MerchantInvoicesPrefix, merchant.Bytes()...)
}

func merchantInvoiceKey(merchant sdk.AccAddress) []byte {}, invoiceID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(invoiceID)
	return append(merchantInvoicesKey(merchant), bz...)
}

func createdTimeInvoicesKey(createdTime time.Time) []byte {
	return append(CreatedTimeInvoicesPrefix, sdk.FormatTimeBytes(createdTime)...)
}

func createdTimeInvoiceKey(createdTime time.Time, invoiceID uint64) []byte {
	bz := sdk.Uint64ToBigEndian(invoiceID)
	return append(createdTimeInvoicesKey(createdTime), bz...)
}