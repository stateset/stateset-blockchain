package loan

import (
	"net/url"
	"time"

	app "github.com/stateset/stateset-blockchain/types"
	
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/gaskv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	log "github.com/tendermint/tendermint/libs/log"
)

// Keeper is the model object for the module
type Keeper struct {
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace

	accountKeeper   AccountKeeper
	marketKeeper market.Keeper
}

// NewKeeper creates a new account keeper
func NewKeeper(storeKey sdk.StoreKey, paramStore params.Subspace, codec *codec.Codec, accountKeeper AccountKeeper, marketKeeper market.Keeper) Keeper {
	return Keeper{
		storeKey,
		codec,
		paramStore.WithKeyTable(ParamKeyTable()),
		accountKeeper,
		marketKeeper,
	}
}

// SubmitLoan creates a new loan in the loan key-value store
func (k Keeper) SubmitLoan(ctx sdk.Context, body, loanID string,
	lender sdk.AccAddress, source url.URL) (loan Loan, err sdk.Error) {

	err = k.validateLength(ctx, body)
	if err != nil {
		return
	}
	jailed, err := k.loanKeeper.IsJailed(ctx, lender)
	if err != nil {
		return
	}
	if jailed {
		return loan, ErrCreatorJailed(lender)
	}
	market, err := k.marketKeeper.market(ctx, MarketID)
	if err != nil {
		return loan, ErrInvalidMarketID(MarketID.ID)
	}

	loanID, err := k.loanID(ctx)
	if err != nil {
		return
	}
	loan = NewLoan(loanID, invoiceID, accountId, body, creator, source,
		ctx.BlockHeader().Time,
	)

	// persist loan
	k.setLoan(ctx, loan)
	// increment loanID (primary key) for next loan
	k.setLoanID(ctx, loanID+1)

	// persist associations
	k.setControllerLoan(ctx, loan.Merchant, loanID)
	k.setProcessorLoan(ctx, loan.Lender, loanId)
	k.setCreatedTimeLoan(ctx, loan.CreatedTime, loanID)

	logger(ctx).Info("Submitted " + loan.String())

	return loan, nil
}


func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return gaskv.NewStore(ctx.MultiStore().GetKVStore(k.storeKey), ctx.GasMeter(), app.KVGasConfig())
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}