package types

// Stateset Invoice IBC events
const (
	EventTypeTimeout = "timeout"
	EventTypeIbcPurchaseOrderPacket = "ibcInvoice_packet"

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"
)
