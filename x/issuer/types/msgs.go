// This software is Copyright (c) 2019 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/e-money/em-ledger/util"
)

var (
	_ sdk.Msg = MsgIncreaseMintable{}
	_ sdk.Msg = MsgDecreaseMintable{}
	_ sdk.Msg = MsgRevokeLiquidityProvider{}
	_ sdk.Msg = MsgSetInflation{}
)

type (
	MsgIncreaseMintable struct {
		Issuer            sdk.AccAddress `json:"issuer" yaml:"issuer"`
		LiquidityProvider sdk.AccAddress `json:"liquidity_provider" yaml:"liquidity_provider"`
		MintableIncrease  sdk.Coins      `json:"amount" yaml:"amount"`
	}

	MsgDecreaseMintable struct {
		Issuer            sdk.AccAddress `json:"issuer" yaml:"issuer"`
		LiquidityProvider sdk.AccAddress `json:"liquidity_provider" yaml:"liquidity_provider"`
		MintableDecrease  sdk.Coins      `json:"amount" yaml:"amount"`
	}

	MsgRevokeLiquidityProvider struct {
		Issuer            sdk.AccAddress `json:"issuer" yaml:"issuer"`
		LiquidityProvider sdk.AccAddress `json:"liquidity_provider" yaml:"liquidity_provider"`
	}

	MsgSetInflation struct {
		Issuer        sdk.AccAddress `json:"issuer" yaml:"issuer"`
		Denom         string         `json:"denom" yaml:"denom"`
		InflationRate sdk.Dec        `json:"inflation_rate" yaml:"inflation_rate"`
	}
)

func (msg MsgSetInflation) Route() string { return ModuleName }

func (msg MsgSetInflation) Type() string { return "setInflation" }

func (msg MsgSetInflation) ValidateBasic() sdk.Error {
	if msg.InflationRate.IsNegative() {
		return ErrNegativeInflation()
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if !util.ValidateDenom(msg.Denom) {
		return ErrInvalidDenomination(msg.Denom)
	}

	return nil
}

func (msg MsgSetInflation) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetInflation) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}

func (msg MsgRevokeLiquidityProvider) Route() string { return ModuleName }

func (msg MsgRevokeLiquidityProvider) Type() string { return "revokeLiquidityProvider" }

func (msg MsgRevokeLiquidityProvider) ValidateBasic() sdk.Error {
	if msg.LiquidityProvider.Empty() {
		return sdk.ErrInvalidAddress("missing liquidity provider address")
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	return nil
}

func (msg MsgRevokeLiquidityProvider) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRevokeLiquidityProvider) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}

func (msg MsgDecreaseMintable) Route() string { return ModuleName }

func (msg MsgDecreaseMintable) Type() string { return "decreaseMintable" }

func (msg MsgDecreaseMintable) ValidateBasic() sdk.Error {
	if msg.LiquidityProvider.Empty() {
		return sdk.ErrInvalidAddress("missing liquidity provider address")
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if !msg.MintableDecrease.IsValid() {
		return sdk.ErrInvalidCoins("requested decrease is invalid: " + msg.MintableDecrease.String())
	}

	return nil
}

func (msg MsgDecreaseMintable) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDecreaseMintable) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}

func (msg MsgIncreaseMintable) Route() string { return ModuleName }

func (msg MsgIncreaseMintable) Type() string { return "increaseMintable" }

func (msg MsgIncreaseMintable) ValidateBasic() sdk.Error {
	if msg.LiquidityProvider.Empty() {
		return sdk.ErrInvalidAddress("missing liquidity provider address")
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if !msg.MintableIncrease.IsValid() {
		return sdk.ErrInvalidCoins("mintable increase is invalid: " + msg.MintableIncrease.String())
	}

	return nil
}

func (msg MsgIncreaseMintable) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIncreaseMintable) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}
