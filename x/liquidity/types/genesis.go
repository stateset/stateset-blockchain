package types

// NewGenesisState is the constructor function for GenesisState
func NewGenesisState(params Params, pools []Pool) *GenesisState {
	return &GenesisState{
		Params:      params,
		Pools: pools,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), []Pool{})
}

// ValidateGenesis - placeholder function
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	for _, p := range data.Pools {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate Liquidity Pool after init or after export
func (p PoolRecord) Validate() error {
	if p.PoolBatch.DepositMsgIndex == 0 ||
		(len(p.DepositMsgStates) > 0 && p.PoolBatch.DepositMsgIndex != p.DepositMsgStates[len(p.DepositMsgStates)-1].MsgIndex+1) {
		return ErrBadBatchMsgIndex
	}
	if p.PoolBatch.WithdrawMsgIndex == 0 ||
		(len(p.WithdrawMsgStates) != 0 && p.PoolBatch.WithdrawMsgIndex != p.WithdrawMsgStates[len(p.WithdrawMsgStates)-1].MsgIndex+1) {
		return ErrBadBatchMsgIndex
	}
	if p.PoolBatch.SwapMsgIndex == 0 ||
		(len(p.SwapMsgStates) != 0 && p.PoolBatch.SwapMsgIndex != p.SwapMsgStates[len(p.SwapMsgStates)-1].MsgIndex+1) {
		return ErrBadBatchMsgIndex
	}

	return nil
}