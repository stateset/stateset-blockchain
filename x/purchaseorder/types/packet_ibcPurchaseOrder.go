package types

// ValidateBasic is used for validating the packet
func (p IbcPurchaseOrderPacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p IbcPurchaseOrderPacketData) GetBytes() ([]byte, error) {
	var modulePacket PurchaseorderPacketData

	modulePacket.Packet = &PurchaseorderPacketData_IbcPurchaseOrderPacket{&p}

	return modulePacket.Marshal()
}
