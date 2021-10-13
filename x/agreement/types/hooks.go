package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type AgreementHooks interface {
	AfterAgreementCreated(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64)
	AfterFinanceAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64, amount sdk.Coins)
	AfterCompleteAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64, amount sdk.Coins)
	AfterCancelAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64, amount sdk.Coins)
	AfterRenewAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64, amount sdk.Coins)
	AfterAmendAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64, amount sdk.Coins)
	AfterExpireAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64, amount sdk.Coins)
	AfterLockAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64, amount sdk.Coins)
}

var _ AgreementHooks = MultiAgreementHooks{}

// combine multiple gamm hooks, all hook functions are run in array sequence
type MultiAgreementHooks []AgreementHooks

// Creates hooks for the Agreement Module
func NewAgreementGammHooks(hooks ...AgreementHooks) MultiAgreementHooks {
	return hooks
}

func (h MultiAgreementHooks) AfterAgreementCreated(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64) {
	for i := range h {
		h[i].AfterAgreementCreated(ctx, sender, agreementId)
	}
}

func (h MultiAgreementHooks) AfterFinanceAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64) {
	for i := range h {
		h[i].AfterFinanceAgreement(ctx, sender, agreementId, amount)
	}
}

func (h MultiAgreementHooks) AfterCompleteAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId) {
	for i := range h {
		h[i].AfterCompleteAgreement(ctx, sender, poolId, amount)
	}
}

func (h MultiAgreementHooks) AfterLockAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64) {
	for i := range h {
		h[i].AfterLockAgreement(ctx, sender, poolId, input, output)
	}
}

func (h MultiAgreementHooks) AfterRenewAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64) {
	for i := range h {
		h[i].AfterRenewAgreement(ctx, sender, poolId, input, output)
	}
}

func (h MultiAgreementHooks) AfterAmendAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64) {
	for i := range h {
		h[i].AfterAmendAgreement(ctx, sender, poolId, input, output)
	}
}

func (h MultiAgreementHooks) AfterExpireAgreement(ctx sdk.Context, sender sdk.AccAddress, agreementId uint64) {
	for i := range h {
		h[i].AfterExpireAgreement(ctx, sender, poolId, input, output)
	}
}