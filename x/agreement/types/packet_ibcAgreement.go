package types

// ValidateBasic is used for validating the packet
func (p IbcAgreementPacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p IbcAgreementPacketData) GetBytes() ([]byte, error) {
	var modulePacket AgreementPacketData

	modulePacket.Packet = &AgreementPacketData_IbcAgreementPacket{&p}

	return modulePacket.Marshal()
}
