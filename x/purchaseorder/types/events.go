package types

// event types
const (
	TypeEvtPurchaseOrderCreated   = "purchase_order_created"
	TypeEvtPurchaseOrderUpdated  = "purchase_order_updated"
	TypeEvtPurchaseOrderDeleted = "purchase_order_deleted"
	TypeEvtPurchaseOrderCompleted = "purchase_order_completed"
	TypeEvtPurchaseOrderDeleted = "purchase_order_deleted"
	TypeEvtPurchaseOrderCanceled   = "purchase_order_canceled"
	TypeEvtPurchaseOrderLocked   = "purchase_order_locked"
	TypeEvtPurchaseOrderFinanced   = "purchase_order_financed"

	AttributeValueCategory = ModuleName
	AttributeKeyPurchaseOrderId     = "purchase_order_id"
	AttributeKeySwapFee    = "swap_fee"
	AttributeKeyTokensIn   = "tokens_in"
	AttributeKeyTokensOut  = "tokens_out"
)