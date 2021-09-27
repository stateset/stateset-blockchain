package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

type (
	Keeper struct {
		cdc           codec.Marshaler
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey

		paramSpace paramtypes.Subspace
		hooks      types.PurchaseOrderHooks
	
		// keepers
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper

		channelKeeper types.ChannelKeeper
		portKeeper    types.PortKeeper
		scopedKeeper  types.ScopedKeeper
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	scopedKeeper types.ScopedKeeper,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		scopedKeeper:  scopedKeeper,
	}
}


func (k *Keeper) createFinanceEvent(ctx sdk.Context, sender sdk.AccAddress, purchaseOrderId uint64, input sdk.Coins) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtPurchaseOrderFinanced,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyPurchaseOrderId, strconv.FormatUint(purchaseOrderId, 10)),
			sdk.NewAttribute(types.AttributeKeyTokensIn, input.String()),
		),
	})
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
