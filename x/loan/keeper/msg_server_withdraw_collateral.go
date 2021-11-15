package keeper

import (
	"context"
	"fmt"

	"github.com/stateset/stateset-blockchain/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) WithdrawCollateral(goCtx context.Context, msg *types.MsgWithdrawCollateral) (*types.MsgWithdrawCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	if loan.Borrower != msg.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Cannot cancel: not the borrower")
	}

	if loan.State != "active" {
		return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}

	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral)

	loan.State = "withdrawed"

	k.SetLoan(ctx, loan)

	return &types.MsgCancelLoanResponse{}, nil
}