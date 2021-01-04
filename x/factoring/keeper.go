package factoring

import (
	"time"

	app "github.com/stateset/stateset-blockchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/gaskv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper is the model object for the package staking module
type Keeper struct {
	storeKey      sdk.StoreKey
	codec         *codec.Codec
	paramStore    params.Subspace
	codespace     sdk.CodespaceType
	bankKeeper    BankKeeper
	accountKeeper AccountKeeper
	invoiceKeeper InvoiceKeeper
	loanKeeper    LoanKeeper
	supplyKeeper  supply.Keeper
}

// NewKeeper creates a staking keeper.
func NewKeeper(codec *codec.Codec, storeKey sdk.StoreKey,
	accountKeeper AccountKeeper, bankKeeper BankKeeper, invoiceKeeper InvoiceKeeper, loanKeeper LoanKeeper, supplyKeeper supply.Keeper,
	paramStore params.Subspace,
	codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:      storeKey,
		codec:         codec,
		paramStore:    paramStore.WithKeyTable(ParamKeyTable()),
		codespace:     codespace,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		invoiceKeeper: invoiceKeeper,
		loanKeeper:    loanKeeper,
		supplyKeeper:  supplyKeeper,
	}
}

func (k Keeper) Arguments(ctx sdk.Context) []Argument {
	arguments := make([]Argument, 0)
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), ArgumentsKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var argument Argument
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &argument)
		arguments = append(arguments, argument)
	}
	return arguments
}

func (k Keeper) Stakes(ctx sdk.Context) []Stake {
	stakes := make([]Stake, 0)
	iterator := sdk.KVStorePrefixIterator(k.store(ctx), StakesKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var stake Stake
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &stake)
		stakes = append(stakes, stake)
	}
	return stakes
}

func (k Keeper) InvoiceLoans(ctx sdk.Context, invoiceID uint64) []Loan {
	loans := make([]Loan, 0)
	k.IterateInvoiceLoans(ctx, invoiceID, func(loan Loan) bool {
		loans = append(loans, loan)
		return false
	})
	return arguments
}

func (k Keeper) LoanStakes(ctx sdk.Context, loanID uint64) []Stake {
	stakes := make([]Stake, 0)
	k.IterateLoanStakes(ctx, loanID, func(stake Stake) bool {
		stakes = append(stakes, stake)
		return false
	})
	return stakes
}

func (k Keeper) MarketStakes(ctx sdk.Context, MarketID string) []Stake {
	stakes := make([]Stake, 0)
	k.IterateMarketStakes(ctx, MarketID, func(stake Stake) bool {
		stakes = append(stakes, stake)
		return false
	})
	return stakes
}

func (k Keeper) UserStakes(ctx sdk.Context, address sdk.AccAddress) []Stake {
	stakes := make([]Stake, 0)
	k.IterateUserStakes(ctx, address, func(stake Stake) bool {
		stakes = append(stakes, stake)
		return false
	})
	return stakes
}

func (k Keeper) UserMarketStakes(ctx sdk.Context, address sdk.AccAddress, MarketID string) []Stake {
	stakes := make([]Stake, 0)
	k.IterateUserMarketStakes(ctx, address, MarketID, func(stake Stake) bool {
		stakes = append(stakes, stake)
		return false
	})
	return stakes
}

func (k Keeper) UserLoans(ctx sdk.Context, address sdk.AccAddress) []Loan {
	loans := make([]Loan, 0)
	k.IterateUserLoans(ctx, address, func(loan Loan) bool {
		loans = append(loans, loan)
		return false
	})
	return arguments
}

func (k Keeper) FactorInvoice(ctx sdk.Context, amount,
	factor sdk.AccAddress, invoiceID uint64, stakeType StakeType) (Loan, sdk.Error) {
	// only backing or challenge
	if !stakeType.ValidForLoan() {
		return Loan{}, ErrCodeInvalidStakeType(stakeType)
	}
	err := k.checkJailed(ctx, factor)
	if err != nil {
		return Loan{}, err
	}
	invoice, ok := k.invoiceKeeper.Invoice(ctx, invoiceID)
	if !ok {
		return Loan{}, ErrCodeUnknownInvoice(invoiceID)
	}

	loans := k.InvoiceLoans(ctx, invoiceID)
	count := 0
	for _, a := range loans {
		if a.Factor.Equals(factor) {
			count++
		}
	}
	p := k.GetParams(ctx)
	if count >= p.MaxLoansPerInvoice {
		return Invoice{}, ErrCodeMaxNumOfLoansReached(p.MaxLoansPerInvoice)
	}

	loanAmount := p.LoanCreationStake
	loanID, err := k.loanID(ctx)
	if err != nil {
		return Loan{}, err
	}
	loan := Loan{
		ID:           loanID,
		Factor:       factor,
		invoiceID:    invoiceID,
		MarketID: invoice.MarketID,
		Amount:       amount,
		StakeType:    stakeType,
		CreatedTime:  ctx.BlockHeader().Time,
		UpdatedTime:  ctx.BlockHeader().Time,
		TotalStake:   loanAmount,
		EditedTime:   ctx.BlockHeader().Time,
		Edited:       false,
	}
	_, err = k.newStake(ctx, loanAmount, factor, stakeType, loan.ID, invoice.MarketID)
	if err != nil {
		return Loan{}, err
	}

	k.setLoan(ctx, loan)
	k.setLoanID(ctx, loanID+1)
	k.setInvoiceLoan(ctx, invoiceID, loan.ID)
	k.setUserLoan(ctx, factor, loan.ID)

	if invoice.FirstLoanTime.Equal(time.Time{}) {
		err = k.invoiceKeeper.SetFirstLoanTime(ctx, invoiceID, ctx.BlockHeader().Time)
		if err != nil {
			return Loan{}, err
		}
	}

	switch {
	case stakeType == StakeFactoring:
		err := k.invoiceKeeper.AddFactoringStake(ctx, invoiceID, loanAmount)
		if err != nil {
			return Loan{}, err
		}
	case stakeType == StakeLending:
		err := k.invoiceKeeper.AddLendingStake(ctx, invoiceID, loanAmount)
		if err != nil {
			return Loan{}, err
		}
	}

	return loan, nil
}

func (k Keeper) Loan(ctx sdk.Context, loanID uint64) (Loan, bool) {
	loan := Loan{}
	bz := k.store(ctx).Get(loanKey(loanID))
	if bz == nil {
		return Loan{}, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &loan)
	return loan, true
}


func (k Keeper) SetStakeExpired(ctx sdk.Context, stakeID uint64) sdk.Error {
	stake, ok := k.Stake(ctx, stakeID)
	if !ok {
		return ErrCodeUnknownStake(stakeID)
	}
	stake.Expired = true
	k.setStake(ctx, stake)
	return nil
}

// AddAdmin adds a new admin
func (k Keeper) AddAdmin(ctx sdk.Context, admin, factor sdk.AccAddress) (err sdk.Error) {
	params := k.GetParams(ctx)

	// first admin can be added without any authorisation
	if len(params.StakingAdmins) > 0 && !k.isAdmin(ctx, factor) {
		err = ErrAddressNotAuthorised()
	}

	// if already present, don't add again
	for _, currentAdmin := range params.StakingAdmins {
		if currentAdmin.Equals(admin) {
			return
		}
	}

	params.StakingAdmins = append(params.StakingAdmins, admin)

	k.SetParams(ctx, params)

	return
}

// RemoveAdmin removes an admin
func (k Keeper) RemoveAdmin(ctx sdk.Context, admin, remover sdk.AccAddress) (err sdk.Error) {
	if !k.isAdmin(ctx, remover) {
		err = ErrAddressNotAuthorised()
	}

	params := k.GetParams(ctx)
	for i, currentAdmin := range params.StakingAdmins {
		if currentAdmin.Equals(admin) {
			params.StakingAdmins = append(params.StakingAdmins[:i], params.StakingAdmins[i+1:]...)
		}
	}

	k.SetParams(ctx, params)

	return
}

func (k Keeper) isAdmin(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, admin := range k.GetParams(ctx).StakingAdmins {
		if address.Equals(admin) {
			return true
		}
	}
	return false
}

func (k Keeper) setLoan(ctx sdk.Context, loan Loan) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(loan)
	k.store(ctx).Set(loanKey(loan.ID), bz)
}

var tierLimitsEarnedCoins = []sdk.Int{
	sdk.NewInt(app.Set * 10),
	sdk.NewInt(app.Set * 20),
	sdk.NewInt(app.Set * 30),
	sdk.NewInt(app.Set * 40),
	sdk.NewInt(app.Set * 50),
}

var tierLimitsStakeAmounts = []sdk.Int{
	sdk.NewInt(app.Set * 1000),
	sdk.NewInt(app.Set * 1500),
	sdk.NewInt(app.Set * 2000),
	sdk.NewInt(app.Set * 2500),
	sdk.NewInt(app.Set * 3000),
}

var defaultStakeLimit = sdk.NewInt(app.Set * 500)
var defaultMinimumBalance = sdk.NewInt(app.Set * 50)

func (k Keeper) checkStakeThreshold(ctx sdk.Context, address sdk.AccAddress, amount sdk.Int) sdk.Error {
	balance := k.bankKeeper.GetCoins(ctx, address).AmountOf(app.StakeDenom)
	if balance.IsZero() {
		return sdk.ErrInsufficientFunds("Insufficient coins")
	}
	p := k.GetParams(ctx)
	period := p.Period

	staked := sdk.NewInt(0)
	fromDate := ctx.BlockHeader().Time.Add(time.Duration(-1) * period)
	k.IterateAfterCreatedTimeUserStakes(ctx, address,
		fromDate, func(stake Stake) bool {
			// only account for non expired since expired would already have refunded the stake
			if stake.Expired {
				return false
			}
			staked = staked.Add(stake.Amount.Amount)
			return false
		},
	)
	if balance.Sub(amount).LT(defaultMinimumBalance) {
		return ErrCodeMinBalance()
	}

	switch totalEarned := k.TotalEarnedCoins(ctx, address); {
	// if total earned >= 50
	case totalEarned.GTE(tierLimitsEarnedCoins[4]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[4]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	// if total earned >= 40
	case totalEarned.GTE(tierLimitsEarnedCoins[3]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[3]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	// if total earned >= 30
	case totalEarned.GTE(tierLimitsEarnedCoins[2]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[2]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	// if total earned >= 20
	case totalEarned.GTE(tierLimitsEarnedCoins[1]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[1]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	// if total earned >= 10
	case totalEarned.GTE(tierLimitsEarnedCoins[0]):
		if staked.Add(amount).GT(tierLimitsStakeAmounts[0]) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	default:
		if staked.Add(amount).GT(defaultStakeLimit) {
			return ErrCodeMaxAmountStakingReached()
		}
		return nil
	}
}

func (k Keeper) TotalEarnedCoins(ctx sdk.Context, factor sdk.AccAddress) sdk.Int {
	earnedCoins := k.getEarnedCoins(ctx, factor)
	total := sdk.NewInt(0)
	for _, e := range earnedCoins {
		total = total.Add(e.Amount)
	}
	return total
}

func (k Keeper) newStake(ctx sdk.Context, amount sdk.Coin, factor sdk.AccAddress,
	stakeType StakeType, loanID uint64, MarketID string) (Stake, sdk.Error) {
	if !stakeType.Valid() {
		return Stake{}, ErrCodeInvalidStakeType(stakeType)
	}
	err := k.checkStakeThreshold(ctx, lender, amount.Amount)
	if err != nil {
		return Stake{}, err
	}
	period := k.GetParams(ctx).Period
	stakeID, err := k.stakeID(ctx)
	if err != nil {
		return Stake{}, err
	}

	_, err = k.bankKeeper.SubtractCoin(ctx, creator, amount,
		argumentID, stakeType.BankTransactionType(), WithMarkerplaceID(MarketID),
		ToModuleAccount(UserStakesPoolName),
	)
	if err != nil {
		return Stake{}, err
	}

	stake := Stake{
		ID:          stakeID,
		LoanID:  	 loanID,
		MarketID: MarketID,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(period),
		Factor:      factor,
		Amount:      amount,
		Type:        stakeType,
	}
	k.setStake(ctx, stake)
	k.setStakeID(ctx, stakeID+1)
	k.InsertActiveStakeQueue(ctx, stakeID, stake.EndTime)
	k.setLoanStake(ctx, loanID, stake.ID)
	k.setUserStake(ctx, factor, stake.CreatedTime, stake.ID)
	k.setMarketStake(ctx, MarketID, stake.ID)
	k.setUserMarketStake(ctx, stake.Factor, MarketID, stakeID)
	return stake, nil
}

func (k Keeper) Stake(ctx sdk.Context, stakeID uint64) (Stake, bool) {
	stake := Stake{}
	bz := k.store(ctx).Get(stakeKey(stakeID))
	if bz == nil {
		return stake, false
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &stake)
	return stake, true
}

func (k Keeper) setStake(ctx sdk.Context, stake Stake) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stake)
	k.store(ctx).Set(stakeKey(stake.ID), bz)
}

func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return gaskv.NewStore(ctx.MultiStore().GetKVStore(k.storeKey), ctx.GasMeter(), app.KVGasConfig())
}

func (k Keeper) setStakeID(ctx sdk.Context, stakeID uint64) {
	k.setID(ctx, StakeIDKey, stakeID)
}

func (k Keeper) setLoanID(ctx sdk.Context, loanID uint64) {
	k.setID(ctx, LoanIDKey, loanID)
}

func (k Keeper) setID(ctx sdk.Context, key []byte, length uint64) {
	b := k.codec.MustMarshalBinaryBare(length)
	k.store(ctx).Set(key, b)
}

func (k Keeper) getEarnedCoins(ctx sdk.Context, user sdk.AccAddress) sdk.Coins {
	earnedCoins := sdk.Coins{}
	bz := k.store(ctx).Get(userEarnedCoinsKey(user))
	if bz == nil {
		return sdk.NewCoins()
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &earnedCoins)
	return earnedCoins
}

func (k Keeper) setEarnedCoins(ctx sdk.Context, user sdk.AccAddress, earnedCoins sdk.Coins) {
	b := k.codec.MustMarshalBinaryLengthPrefixed(earnedCoins)
	k.store(ctx).Set(userEarnedCoinsKey(user), b)
}

func (k Keeper) addEarnedCoin(ctx sdk.Context, user sdk.AccAddress, MarketID string, amount sdk.Int) {
	earnedCoins := k.getEarnedCoins(ctx, user)
	earnedCoins = earnedCoins.Add(sdk.NewCoins(sdk.NewCoin(MarketID, amount)))
	k.setEarnedCoins(ctx, user, earnedCoins)
}

func (k Keeper) SubtractEarnedCoin(ctx sdk.Context, user sdk.AccAddress, MarketID string, amount sdk.Int) {
	earnedCoins := k.getEarnedCoins(ctx, user)
	earnedCoins = earnedCoins.Sub(sdk.NewCoins(sdk.NewCoin(MarketID, amount)))
	k.setEarnedCoins(ctx, user, earnedCoins)
}

func (k Keeper) stakeID(ctx sdk.Context) (uint64, sdk.Error) {
	id, err := k.getID(ctx, StakeIDKey)
	if err != nil {
		return 0, ErrCodeUnknownStake(id)
	}
	return id, nil
}

func (k Keeper) loanID(ctx sdk.Context) (uint64, sdk.Error) {
	id, err := k.getID(ctx, LoanIDKey)
	if err != nil {
		return 0, ErrCodeUnknownLoan(id)
	}
	return id, nil
}

func (k Keeper) getID(ctx sdk.Context, key []byte) (uint64, sdk.Error) {
	var id uint64
	b := k.store(ctx).Get(key)
	if b == nil {
		return 0, sdk.ErrInternal("unknown id")
	}
	k.codec.MustUnmarshalBinaryBare(b, &id)
	return id, nil
}

// InsertActiveStakeQueue inserts a stakeID into the active stake queue at endTime
func (k Keeper) InsertActiveStakeQueue(ctx sdk.Context, stakeID uint64, endTime time.Time) {
	bz := k.codec.MustMarshalBinaryLengthPrefixed(stakeID)
	k.store(ctx).Set(activeStakeQueueKey(stakeID, endTime), bz)
}

// RemoveFromActiveStakeQueue removes a stakeID from the Active Stake Queue
func (k Keeper) RemoveFromActiveStakeQueue(ctx sdk.Context, stakeID uint64, endTime time.Time) {
	k.store(ctx).Delete(activeStakeQueueKey(stakeID, endTime))
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", ModuleName)
}

func (k Keeper) UsersEarnings(ctx sdk.Context) []UserEarnedCoins {
	userEarnedCoins := make([]UserEarnedCoins, 0)
	k.IterateUserEarnedCoins(ctx, func(address sdk.AccAddress, coins sdk.Coins) bool {
		userEarnedCoins = append(userEarnedCoins, UserEarnedCoins{
			Address: address,
			Coins:   coins,
		})
		return false
	})
	return userEarnedCoins
}