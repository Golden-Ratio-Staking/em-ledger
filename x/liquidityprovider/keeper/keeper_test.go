// This software is Copyright (c) 2019 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"

	appcodec "github.com/e-money/em-ledger/app/codec"
	"github.com/e-money/em-ledger/x/liquidityprovider/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	accAddr1 = sdk.AccAddress([]byte("account1"))

	defaultMintable = sdk.NewCoins(
		sdk.NewCoin("eeur", sdk.NewIntWithDecimal(1000, 2)),
	)

	initialBalance = sdk.NewCoins(
		sdk.NewCoin("eeur", sdk.NewIntWithDecimal(50, 2)),
		sdk.NewCoin("ejpy", sdk.NewInt(250)),
	)
)

func TestCreateAndMint(t *testing.T) {
	ctx, ak, sk, bk, keeper := createTestComponents(t, initialBalance)

	assert.Equal(t, initialBalance, sk.GetSupply(ctx).GetTotal())

	acc := accAddr1
	account := ak.NewAccountWithAddress(ctx, acc)
	bk.SetBalances(ctx, account.GetAddress(), initialBalance)
	ak.SetAccount(ctx, account)

	// Turn account into a LP
	keeper.CreateLiquidityProvider(ctx, acc, defaultMintable)
	account = ak.GetAccount(ctx, acc)

	assert.IsType(t, &types.LiquidityProviderAccount{}, account)

	toMint := sdk.NewCoins(sdk.NewCoin("eeur", sdk.NewIntWithDecimal(500, 2)))
	keeper.MintTokens(ctx, acc, toMint)

	assert.Equal(t, initialBalance.Add(toMint...), bk.GetAllBalances(ctx, account.GetAddress()))
	assert.Equal(t, initialBalance.Add(toMint...), sk.GetSupply(ctx).GetTotal())

	// Ensure that mintable amount available has been correspondingly reduced
	lpAcc := keeper.GetLiquidityProviderAccount(ctx, acc)
	assert.Equal(t, defaultMintable.Sub(toMint), lpAcc.Mintable)
}

func TestMintTooMuch(t *testing.T) {
	ctx, ak, sk, bk, keeper := createTestComponents(t, initialBalance)

	acc := accAddr1
	account := ak.NewAccountWithAddress(ctx, acc)
	bk.SetBalances(ctx, account.GetAddress(), initialBalance)
	ak.SetAccount(ctx, account)

	// Turn account into a LP
	keeper.CreateLiquidityProvider(ctx, acc, defaultMintable)
	account = ak.GetAccount(ctx, acc)

	toMint := sdk.NewCoins(sdk.NewCoin("eeur", sdk.NewIntWithDecimal(5000, 2)))
	keeper.MintTokens(ctx, acc, toMint)

	assert.Equal(t, initialBalance, bk.GetAllBalances(ctx, account.GetAddress()))
	assert.Equal(t, initialBalance, sk.GetSupply(ctx).GetTotal())

	// Ensure that the mintable amount of the account has not been modified by failed attempt to mint.
	lpAcc := keeper.GetLiquidityProviderAccount(ctx, acc)
	assert.Equal(t, defaultMintable, lpAcc.Mintable)
}

func TestMintMultipleDenoms(t *testing.T) {
	ctx, ak, sk, bk, keeper := createTestComponents(t, initialBalance)

	jpy := sdk.NewCoin("ejpy", sdk.NewInt(1000000))
	extendedMintable := defaultMintable.Add(jpy)

	acc := accAddr1
	account := ak.NewAccountWithAddress(ctx, acc)
	bk.SetBalances(ctx, account.GetAddress(), initialBalance)
	ak.SetAccount(ctx, account)

	// Turn account into a LP
	keeper.CreateLiquidityProvider(ctx, acc, extendedMintable)
	account = ak.GetAccount(ctx, acc)

	toMint := sdk.NewCoins(
		sdk.NewCoin("eeur", sdk.NewIntWithDecimal(500, 2)),
		sdk.NewCoin("ejpy", sdk.NewInt(500000)),
	)

	keeper.MintTokens(ctx, acc, toMint)

	assert.Equal(t, initialBalance.Add(toMint...), bk.GetAllBalances(ctx, account.GetAddress()))
	assert.Equal(t, initialBalance.Add(toMint...), sk.GetSupply(ctx).GetTotal())

	// Ensure that mintable amount available has been correspondingly reduced
	lpAcc := keeper.GetLiquidityProviderAccount(ctx, acc)
	assert.Equal(t, extendedMintable.Sub(toMint), lpAcc.Mintable)
}

func TestMintWithoutLPAccount(t *testing.T) {
	ctx, ak, sk, bk, keeper := createTestComponents(t, initialBalance)

	acc := accAddr1
	account := ak.NewAccountWithAddress(ctx, acc)
	bk.SetBalances(ctx, account.GetAddress(), initialBalance)
	ak.SetAccount(ctx, account)

	toMint := sdk.NewCoins(sdk.NewCoin("eeur", sdk.NewIntWithDecimal(500, 2)))
	keeper.MintTokens(ctx, acc, toMint)

	assert.IsType(t, &auth.BaseAccount{}, account)
	assert.Equal(t, initialBalance, sk.GetSupply(ctx).GetTotal())
	assert.Equal(t, initialBalance, bk.GetAllBalances(ctx, account.GetAddress()))
}

func TestCreateAndRevoke(t *testing.T) {
	ctx, ak, _, bk, keeper := createTestComponents(t, initialBalance)
	acc := accAddr1

	account := ak.NewAccountWithAddress(ctx, acc)
	bk.SetBalances(ctx, account.GetAddress(), initialBalance)
	ak.SetAccount(ctx, account)

	// Turn account into a LP
	keeper.CreateLiquidityProvider(ctx, acc, defaultMintable)
	account = ak.GetAccount(ctx, acc)

	assert.IsType(t, &types.LiquidityProviderAccount{}, account)

	keeper.RevokeLiquidityProviderAccount(ctx, account)
	account = ak.GetAccount(ctx, acc)
	assert.IsType(t, &auth.BaseAccount{}, account)
}

func TestAccountNotFound(t *testing.T) {
	ctx, ak, _, _, keeper := createTestComponents(t, initialBalance)

	acc := accAddr1
	keeper.CreateLiquidityProvider(ctx, acc, defaultMintable)
	assert.Nil(t, ak.GetAccount(ctx, acc))
}

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()

	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return
}

func createTestComponents(t *testing.T, initialSupply sdk.Coins) (sdk.Context, auth.AccountKeeper, supply.Keeper, bank.Keeper, Keeper) {
	cdc := makeTestCodec()

	appCodec := appcodec.NewAppCodec(cdc)

	var (
		keyAcc     = sdk.NewKVStoreKey(auth.StoreKey)
		keyBank    = sdk.NewKVStoreKey(bank.StoreKey)
		keyParams  = sdk.NewKVStoreKey(params.StoreKey)
		keySupply  = sdk.NewKVStoreKey(supply.StoreKey)
		tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "supply-chain"}, true, log.NewNopLogger())

	maccPerms := map[string][]string{
		types.ModuleName: {supply.Minter},
	}

	pk := params.NewKeeper(appCodec, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(appCodec, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(appCodec, keyBank, ak, pk.Subspace(bank.DefaultParamspace), make(map[string]bool))
	sk := supply.NewKeeper(appCodec, keySupply, ak, bk, maccPerms)

	// Empty supply
	sk.SetSupply(ctx, supply.NewSupply(initialSupply))

	keeper := NewKeeper(ak, bk, sk)

	return ctx, ak, sk, bk, keeper
}
