package keeper

import (
	"context"

	"github.com/stateset/stateset-blockchain/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DepositCollateral(goCtx context.Context, msg *types.MsgDepositCollateral) (*types.MsgDepositCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	
	// Get the borrower address
	borrower, _ := sdk.AccAddressFromBech32(loan.Creator)

	// Get the collateral as sdk.Coins
	collateral, err := sdk.ParseCoinsNormalized(loan.Collateral)
	if err != nil {
		panic(err)
	}

	// Use the module account as escrow account
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, borrower, types.ModuleName, collateral)
	if sdkError != nil {
		return nil, sdkError
	}

	// Add the loan to the keeper
	k.AppendLoan(
		ctx,
		loan,
	)

	return &types.MsgDepositCollateralResponse{}, nil
}