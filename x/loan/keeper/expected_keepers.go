package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	// Methods imported from bank should be defined here
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// DidKeeper defines the expected interface needed to add dids.
type DidKeeper interface {
	AddDid(ctx sdk.Context, did) error
	GetDid(ctx sdk.Context, did) error
	AddCredentials(ctx sdk.Context, did, credentials) error
}
