package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stateset/stateset-blockchain/x/liquidity/types"
)

func (k Keeper) ValidateMsgCreatePool(ctx sdk.Context, msg *types.MsgCreatePool) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	params := k.GetParams(ctx)
	var poolType types.PoolType

	// check poolType exist, get poolType from param
	if len(params.PoolTypes) >= int(msg.PoolTypeIndex) {
		poolType = params.PoolTypes[msg.PoolTypeIndex-1]
		if poolType.PoolTypeIndex != msg.PoolTypeIndex {
			return types.ErrPoolTypeNotExists
		}
	} else {
		return types.ErrPoolTypeNotExists
	}

	if poolType.MaxReserveCoinNum > types.MaxReserveCoinNum || types.MinReserveCoinNum > poolType.MinReserveCoinNum {
		return types.ErrNumOfReserveCoin
	}

	// TODO: duplicated with ValidateBasic
	if len(msg.ReserveCoinDenoms) != msg.DepositCoins.Len() {
		return types.ErrNumOfReserveCoin
	}

	if uint32(msg.DepositCoins.Len()) > poolType.MaxReserveCoinNum && poolType.MinReserveCoinNum > uint32(msg.DepositCoins.Len()) {
		return types.ErrNumOfReserveCoin
	}

	denomA, denomB := types.AlphabeticalDenomPair(msg.ReserveCoinDenoms[0], msg.ReserveCoinDenoms[1])
	if denomA != msg.ReserveCoinDenoms[0] || denomB != msg.ReserveCoinDenoms[1] {
		return types.ErrBadOrderingReserveCoin
	}

	if denomA == denomB {
		return types.ErrEqualDenom
	}

	poolKey := types.GetPoolKey(msg.ReserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)
	_, found := k.GetPoolByReserveAccIndex(ctx, reserveAcc)
	if found {
		return types.ErrPoolAlreadyExists
	}
	return nil
}

// Validate logic for Liquidity pool after set or before export
func (k Keeper) ValidatePool(ctx sdk.Context, pool *types.Pool) error {
	params := k.GetParams(ctx)
	var poolType types.PoolType

	// check poolType exist, get poolType from param
	if len(params.PoolTypes) >= int(pool.PoolTypeIndex) {
		poolType = params.PoolTypes[pool.PoolTypeIndex-1]
		if poolType.PoolTypeIndex != pool.PoolTypeIndex {
			return types.ErrPoolTypeNotExists
		}
	} else {
		return types.ErrPoolTypeNotExists
	}

	if poolType.MaxReserveCoinNum > types.MaxReserveCoinNum || types.MinReserveCoinNum > poolType.MinReserveCoinNum {
		return types.ErrNumOfReserveCoin
	}

	reserveCoins := k.GetReserveCoins(ctx, *pool)
	if uint32(reserveCoins.Len()) > poolType.MaxReserveCoinNum && poolType.MinReserveCoinNum > uint32(reserveCoins.Len()) {
		return types.ErrNumOfReserveCoin
	}

	if len(pool.ReserveCoinDenoms) != reserveCoins.Len() {
		return types.ErrNumOfReserveCoin
	}
	for i, denom := range pool.ReserveCoinDenoms {
		if denom != reserveCoins[i].Denom {
			return types.ErrInvalidDenom
		}
	}

	denomA, denomB := types.AlphabeticalDenomPair(pool.ReserveCoinDenoms[0], pool.ReserveCoinDenoms[1])
	if denomA != pool.ReserveCoinDenoms[0] || denomB != pool.ReserveCoinDenoms[1] {
		return types.ErrBadOrderingReserveCoin
	}

	poolKey := types.GetPoolKey(pool.ReserveCoinDenoms, pool.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)
	poolCoin := k.GetPoolCoinTotal(ctx, *pool)
	if poolCoin.Denom != types.GetPoolCoinDenom(reserveAcc) {
		return types.ErrBadPoolCoinDenom
	}

	_, found := k.GetPoolBatch(ctx, pool.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}

	return nil
}

func (k Keeper) CreatePool(ctx sdk.Context, msg *types.MsgCreatePool) error {
	if err := k.ValidateMsgCreatePool(ctx, msg); err != nil {
		return err
	}
	params := k.GetParams(ctx)
	poolKey := types.GetPoolKey(msg.ReserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)

	poolCreator := msg.GetPoolCreator()
	accPoolCreator := k.accountKeeper.GetAccount(ctx, poolCreator)
	poolCreatorBalances := k.bankKeeper.GetAllBalances(ctx, accPoolCreator.GetAddress())

	if !poolCreatorBalances.IsAllGTE(msg.DepositCoins) {
		return types.ErrInsufficientBalance
	}
	for _, coin := range msg.DepositCoins {
		if coin.Amount.LT(params.MinInitDepositToPool) {
			return types.ErrLessThanMinInitDeposit
		}
	}

	if !poolCreatorBalances.IsAllGTE(params.PoolCreationFee.Add(msg.DepositCoins...)) {
		return types.ErrInsufficientPoolCreationFee
	}

	denom1, denom2 := types.AlphabeticalDenomPair(msg.ReserveCoinDenoms[0], msg.ReserveCoinDenoms[1])
	reserveCoinDenoms := []string{denom1, denom2}

	PoolCoinDenom := types.GetPoolCoinDenom(reserveAcc)

	pool := types.Pool{
		//PoolId: will set on SetPoolAtomic
		PoolTypeIndex:         msg.PoolTypeIndex,
		ReserveCoinDenoms:     reserveCoinDenoms,
		ReserveAccountAddress: reserveAcc.String(),
		PoolCoinDenom:         PoolCoinDenom,
	}

	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	mintPoolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, params.InitPoolCoinMintAmount))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
		return err
	}

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	// TODO: write test case
	poolCreationFeePoolAcc := types.GetLiquidityModuleFeePoolAcc()
	inputs = append(inputs, banktypes.NewInput(poolCreator, params.PoolCreationFee))
	outputs = append(outputs, banktypes.NewOutput(poolCreationFeePoolAcc, params.PoolCreationFee))

	inputs = append(inputs, banktypes.NewInput(poolCreator, msg.DepositCoins))
	outputs = append(outputs, banktypes.NewOutput(reserveAcc, msg.DepositCoins))

	inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, mintPoolCoin))
	outputs = append(outputs, banktypes.NewOutput(poolCreator, mintPoolCoin))

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	pool = k.SetPoolAtomic(ctx, pool)
	batch := types.NewPoolBatch(pool.PoolId, 1)

	k.SetPoolBatch(ctx, batch)

	// TODO: remove result state check, debugging
	reserveCoins := k.GetReserveCoins(ctx, pool)
	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	logger := k.Logger(ctx)
	logger.Info("createPool", msg, "pool", pool, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	return nil
}

// Get reserve Coin from the liquidity pool
func (k Keeper) GetReserveCoins(ctx sdk.Context, pool types.Pool) (reserveCoins sdk.Coins) {
	reserveAcc := pool.GetReserveAccount()
	for _, denom := range pool.ReserveCoinDenoms {
		reserveCoins = reserveCoins.Add(k.bankKeeper.GetBalance(ctx, reserveAcc, denom))
	}
	return
}

// Get total supply of pool coin of the pool as sdk.Int
func (k Keeper) GetPoolCoinTotalSupply(ctx sdk.Context, pool types.Pool) sdk.Int {
	supply := k.bankKeeper.GetSupply(ctx)
	total := supply.GetTotal()
	return total.AmountOf(pool.PoolCoinDenom)
}

// Get total supply of pool coin of the pool as sdk.Coin
func (k Keeper) GetPoolCoinTotal(ctx sdk.Context, pool types.Pool) sdk.Coin {
	return sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, pool))
}

// Get meta data of the pool, containing pool coin total supply, Reserved Coins, It used for result of queries
func (k Keeper) GetPoolMetaData(ctx sdk.Context, pool types.Pool) types.PoolMetadata {
	return types.PoolMetadata{
		PoolId:              pool.PoolId,
		PoolCoinTotalSupply: k.GetPoolCoinTotal(ctx, pool),
		ReserveCoins:        k.GetReserveCoins(ctx, pool),
	}
}

// Get Liquidity Pool Record
func (k Keeper) GetPoolRecord(ctx sdk.Context, pool types.Pool) (*types.PoolRecord, bool) {
	batch, found := k.GetPoolBatch(ctx, pool.PoolId)
	if !found {
		return nil, found
	}
	return &types.PoolRecord{
		Pool:         pool,
		PoolMetadata: k.GetPoolMetaData(ctx, pool),
		PoolBatch:    batch,
		BatchPoolDepositMsgs:  k.GetAllPoolBatchDepositMsgs(ctx, batch),
		BatchPoolWithdrawMsgs: k.GetAllPoolBatchWithdrawMsgs(ctx, batch),
		BatchPoolSwapMsgs:     k.GetAllPoolBatchSwapMsgs(ctx, batch),
	}, true
}
func (k Keeper) SetPoolRecord(ctx sdk.Context, record *types.PoolRecord) {
	k.SetPoolAtomic(ctx, record.Pool)
	//k.SetPool(ctx, record.Pool)
	//k.SetPoolByReserveAccIndex(ctx, record.Pool)
	k.GetNextBatchIndexWithUpdate(ctx, record.Pool.PoolId)
	record.PoolBatch.BeginHeight = ctx.BlockHeight()
	k.SetPoolBatch(ctx, record.PoolBatch)
	k.SetPoolBatchDepositMsgs(ctx, record.Pool.PoolId, record.BatchPoolDepositMsgs)
	k.SetPoolBatchWithdrawMsgs(ctx, record.Pool.PoolId, record.BatchPoolWithdrawMsgs)
	k.SetPoolBatchSwapMsgs(ctx, record.Pool.PoolId, record.BatchPoolSwapMsgs)
}

func (k Keeper) ValidateMsgDepositPool(ctx sdk.Context, msg types.MsgDepositToPool) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	if msg.DepositCoins.Len() != len(pool.ReserveCoinDenoms) {
		return types.ErrNumOfReserveCoin
	}

	// TODO: duplicated with ValidateBasic
	if uint32(msg.DepositCoins.Len()) > types.MaxReserveCoinNum ||
		types.MinReserveCoinNum > uint32(msg.DepositCoins.Len()) {
		return types.ErrNumOfReserveCoin
	}
	// TODO: validate msgIndex

	denomA, denomB := types.AlphabeticalDenomPair(msg.DepositCoins[0].Denom, msg.DepositCoins[1].Denom)
	if denomA != pool.ReserveCoinDenoms[0] || denomB != pool.ReserveCoinDenoms[1] {
		return types.ErrNotMatchedReserveCoin
	}
	return nil
}

func (k Keeper) DepositPool(ctx sdk.Context, msg types.BatchPoolDepositMsg) error {
	msg.Executed = true
	k.SetPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg)

	if err := k.ValidateMsgDepositPool(ctx, *msg.Msg); err != nil {
		return err
	}

	depositCoins := msg.Msg.DepositCoins.Sort()

	pool, found := k.GetPool(ctx, msg.Msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	reserveAcc := pool.GetReserveAccount()
	depositor := msg.Msg.GetDepositor()
	params := k.GetParams(ctx)

	reserveCoins := k.GetReserveCoins(ctx, pool)

	// case of reserve coins has run out, ReinitializePool
	if reserveCoins.IsZero() {
		mintPoolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, params.InitPoolCoinMintAmount))
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
			return err
		}
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, msg.Msg.DepositCoins))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, msg.Msg.DepositCoins))

		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, mintPoolCoin))
		outputs = append(outputs, banktypes.NewOutput(depositor, mintPoolCoin))

		// execute multi-send
		if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
			return err
		}

		msg.Succeeded = true
		msg.ToBeDeleted = true
		k.SetPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg)

		// TODO: remove result state check, debugging
		reserveCoins := k.GetReserveCoins(ctx, pool)
		lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
		logger := k.Logger(ctx)
		logger.Info("ReinitializePool", msg, "pool", pool, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
		return nil
	} else if reserveCoins.Len() != msg.Msg.DepositCoins.Len() {
		return types.ErrNumOfReserveCoin
	}

	reserveCoins.Sort()

	coinA := depositCoins[0]
	coinB := depositCoins[1]

	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	depositableAmount := coinB.Amount.ToDec().Mul(lastReserveRatio).TruncateInt()
	depositableAmountA := coinA.Amount
	depositableAmountB := coinB.Amount

	if coinA.Amount.LT(depositableAmount) {
		depositableAmountB = coinA.Amount.ToDec().Quo(lastReserveRatio).TruncateInt()
		refundAmtB := coinB.Amount.Sub(depositableAmountB)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableAmountB))))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableAmountB))))
		// refund
		if refundAmtB.IsPositive() {
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
			outputs = append(outputs, banktypes.NewOutput(depositor, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
		}
	} else if coinA.Amount.GT(depositableAmount) {
		depositableAmountA = coinB.Amount.ToDec().Mul(lastReserveRatio).TruncateInt()
		refundAmtA := coinA.Amount.Sub(depositableAmountA)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableAmountA))))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableAmountA))))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinB)))
		// refund
		if refundAmtA.IsPositive() {
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))))
			outputs = append(outputs, banktypes.NewOutput(depositor, sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))))
		}
	} else {
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinB)))
	}

	// calculate pool token mint amount
	poolCoinAmt := k.GetPoolCoinTotalSupply(ctx, pool).Mul(depositableAmountA).Quo(reserveCoins[0].Amount) // TODO: coinA after executed ?
	mintPoolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt))

	// mint pool token to Depositor
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
		panic(err)
	}
	inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, mintPoolCoin))
	outputs = append(outputs, banktypes.NewOutput(depositor, mintPoolCoin))

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	msg.Succeeded = true
	msg.ToBeDeleted = true
	k.SetPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg)
	// TODO: add events for batch result, each err cases

	// TODO: remove result state check, debugging
	reserveCoins = k.GetReserveCoins(ctx, pool)
	lastReserveRatio = sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	logger := k.Logger(ctx)
	logger.Info("deposit", msg, "pool", pool, "inputs", inputs, "outputs", outputs, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	return nil
}

func (k Keeper) ValidateMsgWithdrawPool(ctx sdk.Context, msg types.MsgWithdrawFromPool) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	if msg.PoolCoin.Denom != pool.PoolCoinDenom {
		return types.ErrBadPoolCoinDenom
	}

	poolCoinTotalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	if msg.PoolCoin.Amount.GT(poolCoinTotalSupply) {
		return types.ErrBadPoolCoinAmount
	}
	return nil
}

func (k Keeper) ValidateMsgSwap(ctx sdk.Context, msg types.MsgSwap) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	denomA, denomB := types.AlphabeticalDenomPair(msg.OfferCoin.Denom, msg.DemandCoinDenom)
	if denomA != pool.ReserveCoinDenoms[0] || denomB != pool.ReserveCoinDenoms[1] {
		return types.ErrNotMatchedReserveCoin
	}

	// can not exceed max order ratio  of reserve coins that can be ordered at a order
	reserveCoinAmt := k.GetReserveCoins(ctx, pool).AmountOf(msg.OfferCoin.Denom)
	maximumOrderableAmt := reserveCoinAmt.ToDec().Mul(types.GetMaxOrderRatio()).TruncateInt()
	if msg.OfferCoin.Amount.GT(maximumOrderableAmt) {
		return types.ErrExceededMaxOrderable
	}
	return nil
}

func (k Keeper) WithdrawPool(ctx sdk.Context, msg types.BatchPoolWithdrawMsg) error {
	msg.Executed = true
	k.SetPoolBatchWithdrawMsg(ctx, msg.Msg.PoolId, msg)

	if err := k.ValidateMsgWithdrawPool(ctx, *msg.Msg); err != nil {
		return err
	}
	// TODO: validate reserveCoin balance

	poolCoin := msg.Msg.PoolCoin
	poolCoins := sdk.NewCoins(poolCoin)

	pool, found := k.GetPool(ctx, msg.Msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	totalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	reserveAcc := pool.GetReserveAccount()
	withdrawer := msg.Msg.GetWithdrawer()

	for _, reserveCoin := range reserveCoins {
		withdrawAmt := reserveCoin.Amount.Mul(poolCoin.Amount).Quo(totalSupply)
		inputs = append(inputs, banktypes.NewInput(reserveAcc,
			sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
		outputs = append(outputs, banktypes.NewOutput(withdrawer,
			sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
	}

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	// burn the escrowed pool coins
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, poolCoins); err != nil {
		panic(err)
	}

	msg.Succeeded = true
	msg.ToBeDeleted = true
	k.SetPoolBatchWithdrawMsg(ctx, msg.Msg.PoolId, msg)
	// TODO: add events for batch result, each err cases

	// TODO: remove result state check, debugging
	reserveCoins = k.GetReserveCoins(ctx, pool)

	var lastReserveRatio sdk.Dec
	if reserveCoins.Empty() {
		lastReserveRatio = sdk.ZeroDec()
	} else {
		lastReserveRatio = sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	}
	logger := k.Logger(ctx)
	logger.Info("withdraw", msg, "pool", pool, "inputs", inputs, "outputs", outputs, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	return nil
}

func (k Keeper) RefundDepositPool(ctx sdk.Context, batchMsg types.BatchPoolDepositMsg) error {
	batchMsg, _ = k.GetPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !batchMsg.Executed || batchMsg.Succeeded {
		panic("can't refund not executed or or already succeed msg")
	}
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetDepositor(), batchMsg.Msg.DepositCoins)
	if err != nil {
		panic(err)
	}
	batchMsg.ToBeDeleted = true
	k.SetPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg)
	//k.DeletePoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	// TODO: not delete now, set toDelete, executed, succeed fail,  delete on next block beginblock
	return err
}

func (k Keeper) RefundWithdrawPool(ctx sdk.Context, batchMsg types.BatchPoolWithdrawMsg) error {
	batchMsg, _ = k.GetPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !batchMsg.Executed || batchMsg.Succeeded {
		panic("can't refund not executed or already succeed msg")
	}
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetWithdrawer(), sdk.NewCoins(batchMsg.Msg.PoolCoin))
	if err != nil {
		panic(err)
	}
	batchMsg.ToBeDeleted = true
	k.SetPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg)
	// not delete now, set toDelete, executed, succeed fail,  delete on next block beginblock
	//k.DeletePoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

// execute transact, refund, expire, send coins with escrow, update state by TransactAndRefundSwapPool
func (k Keeper) TransactAndRefundSwapPool(ctx sdk.Context, batchMsgs []*types.BatchPoolSwapMsg,
	matchResultMap map[uint64]types.MatchResult, pool types.Pool) error {

	var inputs []banktypes.Input
	var outputs []banktypes.Output
	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	poolReserveAcc := pool.GetReserveAccount()
	for _, batchMsg := range batchMsgs {
		// TODO: make Validate function to batchMsg
		if !batchMsg.Executed && batchMsg.Succeeded {
			panic("can't refund not executed with succeed msg")
		}
		if pool.PoolId != batchMsg.Msg.PoolId {
			panic("broken msg pool consistency")
		}

		// full matched, fractional matched
		if msgAfter, ok := matchResultMap[batchMsg.MsgIndex]; ok {
			if batchMsg.MsgIndex != msgAfter.OrderMsgIndex {
				panic("broken msg consistency")
			}

			if (*msgAfter.BatchMsg) != (*batchMsg) {
				panic("broken msg consistency")
			}
			if msgAfter.TransactedCoinAmt.Sub(msgAfter.FeeAmt).IsNegative() ||
				msgAfter.FeeAmt.GT(msgAfter.TransactedCoinAmt) {
				panic("fee over offer")
			}

			// fractional match, but expired order case
			if batchMsg.RemainingOfferCoin.IsPositive() {
				// not to delete, but expired case
				if !batchMsg.ToBeDeleted && batchMsg.OrderExpiryHeight <= ctx.BlockHeight() {
					panic("impossible case")
					// TODO: set to Delete
				} else if !batchMsg.ToBeDeleted && batchMsg.OrderExpiryHeight > ctx.BlockHeight() {
					// fractional matched, to be remaining order, not refund, only transact fractional exchange amt
					// Add transacted coins to multisend
					// TODO: coverage
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.Sub(msgAfter.FeeAmt)))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.Sub(msgAfter.FeeAmt)))))
					inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt))))
					outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt))))

					// Add swap fee to multisend
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.FeeAmt))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.FeeAmt))))

					batchMsg.Succeeded = true

				} else if batchMsg.ToBeDeleted || batchMsg.OrderExpiryHeight == ctx.BlockHeight() {
					// fractional matched, but expired order, transact with refund remaining offer coin

					// Add transacted coins to multisend
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.Sub(msgAfter.FeeAmt)))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.Sub(msgAfter.FeeAmt)))))
					inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt))))
					outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt))))

					// Add swap fee to multisend
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.FeeAmt))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.FeeAmt))))

					// refund remaining coins
					if input, output, err := k.ReleaseEscrowForMultiSend(batchMsg.Msg.GetSwapRequester(),
						sdk.NewCoins(batchMsg.RemainingOfferCoin)); err != nil {
						panic(err)
					} else {
						inputs = append(inputs, input)
						outputs = append(outputs, output)
					}
					batchMsg.Succeeded = true
					batchMsg.ToBeDeleted = true
				} else {
					panic("impossible case ")
					// TODO: set to Delete
				}
			} else if batchMsg.RemainingOfferCoin.IsZero() {
				// full matched case, Add transacted coins to multisend
				inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.Sub(msgAfter.FeeAmt)))))
				outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.Sub(msgAfter.FeeAmt)))))
				inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt))))
				outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
					sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt))))

				// Add swap fee to multisend
				inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.FeeAmt))))
				outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.FeeAmt))))

				batchMsg.Succeeded = true
				batchMsg.ToBeDeleted = true
			} else {
				panic("impossible case")
			}

		} else {
			// not matched, remaining
			if !batchMsg.ToBeDeleted && batchMsg.OrderExpiryHeight > ctx.BlockHeight() {
				// have fractional matching history, not matched and expired, remaining refund
				// refund remaining coins
				// TODO: coverage
				if input, output, err := k.ReleaseEscrowForMultiSend(batchMsg.Msg.GetSwapRequester(),
					sdk.NewCoins(batchMsg.RemainingOfferCoin)); err != nil {
					panic(err)
				} else {
					inputs = append(inputs, input)
					outputs = append(outputs, output)
				}

				batchMsg.Succeeded = false
				batchMsg.ToBeDeleted = true

			} else if batchMsg.ToBeDeleted && batchMsg.OrderExpiryHeight == ctx.BlockHeight() {
				// not matched and expired, remaining refund
				// refund remaining coins
				if input, output, err := k.ReleaseEscrowForMultiSend(batchMsg.Msg.GetSwapRequester(),
					sdk.NewCoins(batchMsg.RemainingOfferCoin)); err != nil {
					panic(err)
				} else {
					inputs = append(inputs, input)
					outputs = append(outputs, output)
				}

				batchMsg.Succeeded = false
				batchMsg.ToBeDeleted = true

			} else {
				panic("impossible case")
			}
		}
	}
	// remove zero coins
	newI := 0
	for _, i := range inputs {
		if i.Coins == nil || i.Coins.Empty() {
		} else {
			inputs[newI] = i
			newI++
		}
		if !i.Coins.IsValid() {
			i.Coins = sdk.NewCoins(i.Coins...) // for sanitizeCoins, remove zero coin
		}
	}
	inputs = inputs[:newI]
	newI = 0
	for _, i := range outputs {
		if i.Coins == nil || i.Coins.Empty() {
		} else {
			outputs[newI] = i
			newI++
		}
		if !i.Coins.IsValid() {
			i.Coins = sdk.NewCoins(i.Coins...) // for sanitizeCoins, remove zero coin
		}
	}
	outputs = outputs[:newI]
	if err := banktypes.ValidateInputsOutputs(inputs, outputs); err != nil {
		return err
	}
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}
	k.SetPoolBatchSwapMsgPointers(ctx, pool.PoolId, batchMsgs)
	return nil
}

//func (k Keeper) GetPoolMetaData(ctx sdk.Context, pool types.Pool) *types.PoolMetadata {
//	totalSupply := sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, pool))
//	reserveCoin := k.GetReserveCoins(ctx, pool).Sort()
//	return &types.PoolMetadata{PoolId: pool.PoolId, PoolCoinTotalSupply: totalSupply, ReserveCoins: reserveCoin}
//}

func (k Keeper) ValidatePoolMetadata(ctx sdk.Context, pool *types.Pool, metaData *types.PoolMetadata) error {
	if !metaData.ReserveCoins.Sort().IsEqual(metaData.ReserveCoins) {
		return types.ErrBadOrderingReserveCoin
	}
	if err := metaData.ReserveCoins.Validate(); err != nil {
		return err
	}
	if !metaData.ReserveCoins.IsEqual(k.GetReserveCoins(ctx, *pool)) {
		return types.ErrNumOfReserveCoin
	}
	if !metaData.PoolCoinTotalSupply.IsEqual(sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, *pool))) {
		return types.ErrBadPoolCoinAmount
	}
	return nil
}

// Validate Liquidity Pool Record after init or after export
func (k Keeper) ValidatePoolRecord(ctx sdk.Context, record *types.PoolRecord) error {
	params := k.GetParams(ctx)
	if err := params.Validate(); err != nil {
		return err
	}

	// validate liquidity pool
	if err := k.ValidatePool(ctx, &record.Pool); err != nil {
		return err
	}

	// validate metadata
	if err := k.ValidatePoolMetadata(ctx, &record.Pool, &record.PoolMetadata); err != nil {
		return err
	}

	// validate each msgs with batch state

	if len(record.BatchPoolDepositMsgs) != 0 && record.PoolBatch.DepositMsgIndex != record.BatchPoolDepositMsgs[len(record.BatchPoolDepositMsgs)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}
	if len(record.BatchPoolWithdrawMsgs) != 0 && record.PoolBatch.WithdrawMsgIndex != record.BatchPoolWithdrawMsgs[len(record.BatchPoolWithdrawMsgs)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}
	if len(record.BatchPoolSwapMsgs) != 0 && record.PoolBatch.SwapMsgIndex != record.BatchPoolSwapMsgs[len(record.BatchPoolSwapMsgs)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}

	// TODO: add verify of escrow amount and poolcoin amount with compare to remaining msgs
	return nil
}