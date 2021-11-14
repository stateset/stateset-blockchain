package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/purchase order module sentinel errors
var (
	ErrPurchaseOrderNotFound      = sdkerrors.Register(ModuleName, 1, "purchase order not found")
	ErrPurchaseOrderAlreadyExist  = sdkerrors.Register(ModuleName, 2, "purchase order already exist")
	ErrPurchaseOrderLocked        = sdkerrors.Register(ModuleName, 3, "purchase order is locked")
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 4, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 5, "invalid version")
	ErrWrongPurchaseOrderState = sdkerrors.Register(ModuleName, 6, "invalide purchase order state")
)
