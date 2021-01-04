package loan

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CodeType is the error code type for the module
type CodeType = sdk.CodeType

// Loan error types
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeInvalidBodyTooShort         CodeType = 101
	ErrorCodeInvalidBodyTooLong          CodeType = 102
	ErrorCodeInvalidID                   CodeType = 103
	ErrorCodeNotFound                    CodeType = 104
	ErrorCodeInvalidSType                CodeType = 105
	ErrorCodeLoansWithMarketNotFound CodeType = 106
	ErrorCodeInvalidSourceURL            CodeType = 107
	ErrorCodeLenderJailed               CodeType = 108
	ErrorCodeAddressNotAuthorised        CodeType = 109
	ErrorCodeJSONParsing                 CodeType = 110
)

// ErrInvalidBodyTooShort throws an error on invalid loan body
func ErrInvalidBodyTooShort(body string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBodyTooShort,
		"Invalid loan body, too short: "+body)
}

// ErrInvalidBodyTooLong throws an error on invalid loan body
func ErrInvalidBodyTooLong() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBodyTooLong,
		"Invalid loan body, too long")
}

// ErrUnknownLoan throws an error on an unknown loan id
func ErrUnknownLoan(id uint64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidID,
		fmt.Sprintf("Unknown loan id: %d", id))
}

// ErrInvalidMarketID throws an error on invalid market id
func ErrInvalidMarketID(id string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeLoansWithMarketNotFound,
		fmt.Sprintf("Invalid market id: %s", id))
}

// ErrInvalidSourceURL throws an error when a URL in invalid
func ErrInvalidSourceURL(url string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidSourceURL,
		"Invalid source URL: "+url)
}

// ErrLenderJailed throws an error on jailed lender
func ErrLenderJailed(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeCreatorJailed,
		"Creator cannot be jailed: "+addr.String())
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