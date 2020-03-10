package account

import (
	"fmt"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateAccount:
			return handleMsgCreateAccount(ctx, keeper, msg)
		case MsgEditAccount:
			return handleMsgEditAccount(ctx, keeper, msg)
		case MsgAddAdmin:
			return handleMsgAddAdmin(ctx, keeper, msg)
		case MsgRemoveAdmin:
			return handleMsgRemoveAdmin(ctx, keeper, msg)
		case MsgUpdateParams:
			return handleMsgUpdateParams(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateAccount(ctx sdk.Context, keeper Keeper, msg MsgCreateAccount) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	// parse url from string
	sourceURL, urlError := url.Parse(msg.Source)
	if urlError != nil {
		return ErrInvalidSourceURL(msg.Source).Result()
	}

	account, err := keeper.SubmitAccount(ctx, msg.Body, msg.MarkerplaceID, msg.Merchant, *sourceURL)
	if err != nil {
		return err.Result()
	}

	res, codecErr := ModuleCodec.MarshalJSON(invoice)
	if codecErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", codecErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

func handleMsgEditAccount(ctx sdk.Context, keeper Keeper, msg MsgEditAccount) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	account, err := keeper.EditAccount(ctx, msg.ID, msg.Body, msg.Editor)
	if err != nil {
		return err.Result()
	}

	res, codecErr := ModuleCodec.MarshalJSON(invoice)
	if codecErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", codecErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

func handleMsgAddAdmin(ctx sdk.Context, k Keeper, msg MsgAddAdmin) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	err := k.AddAdmin(ctx, msg.Admin, msg.Creator)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := ModuleCodec.MarshalJSON(true)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

func handleMsgRemoveAdmin(ctx sdk.Context, k Keeper, msg MsgRemoveAdmin) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	err := k.RemoveAdmin(ctx, msg.Admin, msg.Remover)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := ModuleCodec.MarshalJSON(true)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

func handleMsgUpdateParams(ctx sdk.Context, k Keeper, msg MsgUpdateParams) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	err := k.UpdateParams(ctx, msg.Updater, msg.Updates, msg.UpdatedFields)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := ModuleCodec.MarshalJSON(true)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}