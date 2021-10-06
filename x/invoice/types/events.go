package types

// event types
const (
	TypeEvtInvoiceCreated   = "invoice_created"
	TypeEvtInvoiceUpdated  = "invoice_updated"
	TypeEvtInvoiceDeleted = "invoice_deleted"
	TypeEvtInvoiceCompleted = "invoice_completed"
	TypeEvtInvoiceDeleted = "invoice_deleted"
	TypeEvtInvoiceCanceled   = "invoice_canceled"
	TypeEvtInvoiceLocked   = "invoice_locked"
	TypeEvtInvoiceFactored   = "invoice_factored"
	TypeEvtInvoicePaid   = "invoice_paid"


	AttributeValueCategory = ModuleName
	AttributeKeyInvoiceId     = "invoice_id"
	AttributeKeySwapFee    = "swap_fee"
	AttributeKeyTokensIn   = "tokens_in"
	AttributeKeyTokensOut  = "tokens_out"
)