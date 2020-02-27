// This software is Copyright (c) 2019 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrAccountDoesNotExist = sdkerrors.Register(ModuleName, 1, "account does not exist")
)

//const (
//	DefaultCodespace sdk.CodespaceType = "lp"
//
//	CodeAccountDoesNotExist sdk.CodeType = 1
//)
//
//func ErrAccountDoesNotExist(address sdk.AccAddress) sdk.Error {
//	return sdk.NewError(DefaultCodespace, CodeAccountDoesNotExist, fmt.Sprintf("account %v does not exist", address))
//}
