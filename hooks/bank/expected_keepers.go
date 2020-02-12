// This software is Copyright (c) 2019-2020 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package bank

import (
	"github.com/e-money/em-ledger/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RestrictedKeeper interface {
	GetRestrictedDenoms(sdk.Context) types.RestrictedDenoms
}
