package bank

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for bank module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSendState:
			return handleMsgSendState(ctx, keeper, msg)
		case MsgIssueState:
			return handleMsgIssueState(ctx, keeper, msg)
		case MsgBurnState:
			return handleMsgBurnState(ctx, keeper, msg)
		case MsgUpdateParams:
			return handleMsgUpdateParams(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSendState(ctx sdk.Context, keeper Keeper, msg MsgSendState) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	err := keeper.send(ctx, msg.Sender, msg.Recipient, msg.Amount)
	if err != nil {
		fmt.Println("error", err)
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgUpdateParams(ctx sdk.Context, k Keeper, msg MsgUpdateParams) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	err := k.UpdateParams(ctx, msg.Updates, msg.UpdatedFields)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := json.Marshal(true)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}