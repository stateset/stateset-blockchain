package invoice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CodeType is the error code type for the module
type CodeType = sdk.CodeType

// Invoice error types
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeInvalidBodyTooShort         CodeType = 101
	ErrorCodeInvalidBodyTooLong          CodeType = 102
	ErrorCodeInvalidID                   CodeType = 103
	ErrorCodeNotFound                    CodeType = 104
	ErrorCodeInvalidSType                CodeType = 105
	ErrorCodeInvoicesWithMarketNotFound CodeType = 106
	ErrorCodeInvalidSourceURL            CodeType = 107
	ErrorCodeMerchantJailed               CodeType = 108
	ErrorCodeAddressNotAuthorised        CodeType = 109
	ErrorCodeJSONParsing                 CodeType = 110
)

// ErrInvalidBodyTooShort throws an error on invalid invoice body
func ErrInvalidBodyTooShort(body string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBodyTooShort,
		"Invalid invoice body, too short: "+body)
}

// ErrInvalidBodyTooLong throws an error on invalid invoice body
func ErrInvalidBodyTooLong() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBodyTooLong,
		"Invalid invoice body, too long")
}

// ErrUnknownInvoice throws an error on an unknown invoice id
func ErrUnknownInvoice(id uint64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidID,
		fmt.Sprintf("Unknown invoice id: %d", id))
}

// ErrInvalidMarketID throws an error on invalid market id
func ErrInvalidMarketID(id string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeClaimsWithMarketNotFound,
		fmt.Sprintf("Invalid market id: %s", id))
}

// ErrInvalidSourceURL throws an error when a URL in invalid
func ErrInvalidSourceURL(url string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidSourceURL,
		"Invalid source URL: "+url)
}

// ErrMerchantJailed throws an error on jailed merchant
func ErrMerchantJailed(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeCreatorJailed,
		"Merchant cannot be jailed: "+addr.String())
}

// ErrAddressNotAuthorised throws an error when the address is not admin
func ErrAddressNotAuthorised() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeAddressNotAuthorised,
		"This address is not authorised to perform this action.")
}

// ErrJSONParse throws an error on failed JSON parsing
func ErrJSONParse(err error) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeJSONParsing,
		"JSON parsing error: "+err.Error())
}