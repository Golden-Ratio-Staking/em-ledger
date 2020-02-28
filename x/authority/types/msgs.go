// This software is Copyright (c) 2019 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = MsgCreateIssuer{}
	_ sdk.Msg = MsgDestroyIssuer{}
	_ sdk.Msg = MsgSetGasPrices{}
)

type (
	MsgCreateIssuer struct {
		Authority     sdk.AccAddress `json:"authority" yaml:"authority"`
		Issuer        sdk.AccAddress `json:"issuer" yaml:"issuer"`
		Denominations []string       `json:"denoms" yaml:"denoms"`
	}
	MsgDestroyIssuer struct {
		Authority sdk.AccAddress `json:"authority" yaml:"authority"`
		Issuer    sdk.AccAddress `json:"issuer" yaml:"issuer"`
	}

	MsgSetGasPrices struct {
		Authority sdk.AccAddress `json:"authority" yaml:"authority"`
		GasPrices sdk.DecCoins   `json:"gas_prices" yaml:"gas_prices"`
	}
)

func (msg MsgDestroyIssuer) Type() string { return "destroyIssuer" }

func (msg MsgCreateIssuer) Type() string { return "createIssuer" }

func (msg MsgSetGasPrices) Type() string { return "setgasprices" }

func (msg MsgDestroyIssuer) ValidateBasic() sdk.Error {
	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if msg.Authority.Empty() {
		return sdk.ErrInvalidAddress("missing authority address")
	}

	return nil
}

func (msg MsgCreateIssuer) ValidateBasic() sdk.Error {
	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if msg.Authority.Empty() {
		return sdk.ErrInvalidAddress("missing authority address")
	}

	if len(msg.Denominations) == 0 {
		return ErrNoDenomsSpecified()
	}

	return nil
}

func (msg MsgSetGasPrices) ValidateBasic() sdk.Error {
	if msg.Authority.Empty() {
		return sdk.ErrInvalidAddress("missing authority address")
	}

	if !msg.GasPrices.IsValid() {
		return sdk.ErrInvalidCoins(msg.GasPrices.String())
	}

	return nil
}

func (msg MsgDestroyIssuer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg MsgCreateIssuer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg MsgSetGasPrices) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg MsgDestroyIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetGasPrices) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDestroyIssuer) Route() string { return ModuleName }

func (msg MsgCreateIssuer) Route() string { return ModuleName }

func (msg MsgSetGasPrices) Route() string { return ModuleName }
