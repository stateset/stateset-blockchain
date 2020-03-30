package marketplace

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CodeType is the error code type for the module
type CodeType = sdk.CodeType

// Marketplace error types
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeInvalidBodyTooShort         CodeType = 101
	ErrorCodeInvalidBodyTooLong          CodeType = 102
	ErrorCodeInvalidID                   CodeType = 103
	ErrorCodeNotFound                    CodeType = 104
	ErrorCodeInvalidSType                CodeType = 105
	ErrorCodeInvalidSourceURL            CodeType = 106
	ErrorCodeAddressNotAuthorised        CodeType = 107
	ErrorCodeJSONParsing                 CodeType = 108
	ErrMarketplaceNotFound				 CodeType = 109
)

// ErrInvalidBodyTooShort throws an error on invalid claim body
func ErrInvalidBodyTooShort(body string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBodyTooShort,
		"Invalid marketplace body, too short: "+body)
}

// ErrInvalidBodyTooLong throws an error on invalid claim body
func ErrInvalidBodyTooLong() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBodyTooLong,
		"Invalid marketplace body, too long")
}

// ErrUnknownMarketplace throws an error on an unknown marketplace id
func ErrUnknownMarketplace(id uint64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidID,
		fmt.Sprintf("Unknown marketplace id: %d", id))
}

// ErrUnknownMarketplace throws an error on an unknown marketplace id
func ErrMarketplaceNotFound(id uint64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidID,
		fmt.Sprintf("Unknown marketplace id: %d", id))
}

// ErrInvalidMarketplaceID throws an error on invalid marketplace id
func ErrInvalidMarketplaceID(id string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeClaimsWithMarketplaceNotFound,
		fmt.Sprintf("Invalid marketplace id: %s", id))
}

// ErrInvalidSourceURL throws an error when a URL in invalid
func ErrInvalidSourceURL(url string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidSourceURL,
		"Invalid source URL: "+url)
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