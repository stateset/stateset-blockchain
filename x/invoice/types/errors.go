package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/invoice module sentinel errors
var (
	ErrInvoiceNotFound      = sdkerrors.Register(ModuleName, 1, "invoice not found")
	ErrInvoiceAlreadyExist  = sdkerrors.Register(ModuleName, 2, "invoice already exist")
	ErrInvoiceLocked        = sdkerrors.Register(ModuleName, 3, "invoice is locked")
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 4, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 5, "invalid version")
)
