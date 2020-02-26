// This software is Copyright (c) 2019 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNegativeMintableBalance     = sdkerrors.Register(ModuleName, 1, "")
	ErrNotLiquidityProvider        = sdkerrors.Register(ModuleName, 2, "")
	ErrDoesNotControlDenomination  = sdkerrors.Register(ModuleName, 3, "")
	ErrDenominationAlreadyAssigned = sdkerrors.Register(ModuleName, 4, "")
	ErrIssuerNotFound              = sdkerrors.Register(ModuleName, 5, "")
	ErrNegativeInflation           = sdkerrors.Register(ModuleName, 6, "")
	ErrNotAnIssuer                 = sdkerrors.Register(ModuleName, 7, "")
)

//const (
//	DefaultCodespace sdk.CodespaceType = "iss"
//
//	CodeNegativeMintable      sdk.CodeType = 1
//	CodeNotLiquidityProvider  sdk.CodeType = 2
//	CodeDuplicateDenomination sdk.CodeType = 3
//	CodeIssuerNotFound        sdk.CodeType = 4
//	CodeNegativeInflation     sdk.CodeType = 5
//	CodeDoesNotControlDenom   sdk.CodeType = 6
//	CodeNotAnIssuer           sdk.CodeType = 7
//)
//
//func ErrNotAnIssuer(address sdk.AccAddress) sdk.Error {
//	return sdk.NewError(DefaultCodespace, CodeNotAnIssuer, fmt.Sprintf("%v is not an issuer", address))
//}
//
//func ErrNegativeMintableBalance(lp sdk.AccAddress) sdk.Error {
//	return sdk.NewError(DefaultCodespace, CodeNegativeMintable, fmt.Sprintf("mintable balance decrease would become negative for %d", lp))
//}
//
//func ErrNotLiquidityProvider(lp sdk.AccAddress) sdk.Error {
//	return sdk.NewError(DefaultCodespace, CodeNotLiquidityProvider, fmt.Sprint("account is not a liquidity provider:", lp))
//}
//
//func ErrDoesNotControlDenomination(denom string) sdk.Error {
//	return sdk.NewError(DefaultCodespace, CodeDoesNotControlDenom, fmt.Sprintf("issuer does not control denomination %v", denom))
//}
//
//func ErrDenominationAlreadyAssigned() sdk.Error {
//	return sdk.NewError(DefaultCodespace, CodeDuplicateDenomination, "denomination is already under control of an issuer")
//}
//
//func ErrIssuerNotFound(issuer sdk.AccAddress) sdk.Error {
//	return sdk.NewError(DefaultCodespace, CodeIssuerNotFound, fmt.Sprintf("unable to find issuer %v", issuer))
//}
//
//func ErrNegativeInflation() sdk.Error {
//	return sdk.NewError(DefaultCodespace, CodeNegativeInflation, "cannot set negative inflation")
//}
