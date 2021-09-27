package types

// Stateset Purchase Order IBC events
const (
	EventTypeTimeout = "timeout"
	EventTypeIbcPurchaseOrderPacket = "ibcPurchaseOrder_packet"

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"
)
