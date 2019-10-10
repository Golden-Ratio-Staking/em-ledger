package keeper

import (
	"emoney/x/liquidityprovider/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	authKeeper   auth.AccountKeeper
	supplyKeeper supply.Keeper
}

func NewKeeper(ak auth.AccountKeeper, sk supply.Keeper) Keeper {
	return Keeper{
		authKeeper:   ak,
		supplyKeeper: sk,
	}
}

func (k Keeper) CreateLiquidityProvider(ctx sdk.Context, address sdk.AccAddress, credit sdk.Coins) {
	logger := k.Logger(ctx)

	account := k.authKeeper.GetAccount(ctx, address)
	if account == nil {
		logger.Info("Account not found", "account", address)
		return
	}
	lpAcc := types.NewLiquidityProviderAccount(account, credit)
	k.authKeeper.SetAccount(ctx, lpAcc)

	logger.Info("Created liquidity provider account.", "account", lpAcc.GetAddress())
}

func (k Keeper) MintTokensFromCredit(ctx sdk.Context, liquidityProvider sdk.AccAddress, amount sdk.Coins) {
	logger := k.Logger(ctx)

	account := k.GetLiquidityProviderAccount(ctx, liquidityProvider)
	if account == nil {
		return
	}

	updatedCredit, anyNegative := account.Credit.SafeSub(amount)
	if anyNegative {
		logger.Debug(fmt.Sprintf("Insufficient credit for minting operation"), "requested", amount, "available", account.Credit)
		fmt.Println(" *** Insufficient credit for minting", amount, account.Credit)
		return
	}

	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, amount)
	if err != nil {
		panic(err)
	}

	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, liquidityProvider, amount)
	if err != nil {
		panic(err)
	}

	account = k.GetLiquidityProviderAccount(ctx, liquidityProvider)
	account.Credit = updatedCredit
	k.SetLiquidityProviderAccount(ctx, account)
}

func (k Keeper) SetLiquidityProviderAccount(ctx sdk.Context, account *types.LiquidityProviderAccount) {
	k.authKeeper.SetAccount(ctx, account)
}

func (k Keeper) RevokeLiquidityProviderAccount(ctx sdk.Context, account auth.Account) bool {
	if lpAcc, isLpAcc := account.(*types.LiquidityProviderAccount); isLpAcc {
		account = lpAcc.Account
		k.authKeeper.SetAccount(ctx, account)
		return true
	}

	return false
}

func (k Keeper) GetLiquidityProviderAccount(ctx sdk.Context, liquidityProvider sdk.AccAddress) *types.LiquidityProviderAccount {
	logger := k.Logger(ctx)

	a := k.authKeeper.GetAccount(ctx, liquidityProvider)
	account, ok := a.(*types.LiquidityProviderAccount)
	if !ok {
		logger.Debug(fmt.Sprintf("Account is not a liquidity provider"), "address", liquidityProvider)
		return nil
	}

	return account
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}