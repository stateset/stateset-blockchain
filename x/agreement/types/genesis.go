package types

import (
	"fmt"

	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,
		// this line is used by starport scaffolding # genesis/types/default
		TimedoutAgreementList: []*TimedoutAgreement{},
		SentAgreementList:     []*SentAgreement{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated ID in timedoutAgreement
	timedoutAgreementIdMap := make(map[uint64]bool)

	for _, elem := range gs.TimedoutAgreementList {
		if _, ok := timedoutAgreementIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for timedoutAgreement")
		}
		timedoutAgreementIdMap[elem.Id] = true
	}
	// Check for duplicated ID in sentAgreement
	sentAgreementIdMap := make(map[uint64]bool)

	for _, elem := range gs.SentAgreementList {
		if _, ok := sentAgreementIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for sentAgreement")
		}
		sentAgreementIdMap[elem.Id] = true
	}

	return nil
}
