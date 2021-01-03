package auth

import (
	"github.com/stateset/stateset-blockchain/x/auth/types"
)

var (
	// functions aliases
	RegisterCodec            = types.RegisterCodec
	RegisterAccountTypeCodec = types.RegisterAccountTypeCodec

	// variable aliases
	ModuleCdc = types.ModuleCdc
)