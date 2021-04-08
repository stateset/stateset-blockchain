package types

import (
	"fmt"
)

const (
	ModuleName = "agreement"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

var (
	AgreementPrefix = []byte("stateset_agreement")
	GlobalAgreementNumber  = []byte("stateset_agreement_number")
)

func GetStatesetShareDenom(statesetlId uint64) string {
	return fmt.Sprintf("/stateset/%d", statesetId)
}