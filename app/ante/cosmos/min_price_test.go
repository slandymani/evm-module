package cosmos_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	cosmosante "github.com/slandymani/evm-module/app/ante/cosmos"
	"github.com/slandymani/evm-module/testutil"
	testutiltx "github.com/slandymani/evm-module/testutil/tx"
	"github.com/slandymani/evm-module/utils"
)

var execTypes = []struct {
	name      string
	isCheckTx bool
	simulate  bool
}{
	{"deliverTx", false, false},
	{"deliverTxSimulate", false, true},
}

func (suite *AnteTestSuite) TestMinGasPriceDecorator() {
	denom := utils.BaseDenom
	testMsg := banktypes.MsgSend{
		FromAddress: "evmm1tjdjfavsy956d25hvhs3p0nw9a7pfghqm0up92",
		ToAddress:   "evmm1hdr0lhv75vesvtndlh78ck4cez6esz8u2lk0hq",
		Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: denom}},
	}

	testCases := []struct {
		name                string
		malleate            func() sdk.Tx
		expPass             bool
		errMsg              string
		allowPassOnSimulate bool
	}{
		{
			"invalid cosmos tx type",
			func() sdk.Tx {
				return &testutiltx.InvalidTx{}
			},
			false,
			"invalid transaction type",
			false,
		},
		{
			"valid cosmos tx with MinGasPrices = 0, gasPrice = 0",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.ZeroDec()
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(0), denom, &testMsg)
				return txBuilder.GetTx()
			},
			true,
			"",
			false,
		},
		{
			"valid cosmos tx with MinGasPrices = 0, gasPrice > 0",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.ZeroDec()
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), denom, &testMsg)
				return txBuilder.GetTx()
			},
			true,
			"",
			false,
		},
		{
			"valid cosmos tx with MinGasPrices = 10, gasPrice = 10",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.NewDec(10)
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), denom, &testMsg)
				return txBuilder.GetTx()
			},
			true,
			"",
			false,
		},
		{
			"invalid cosmos tx with MinGasPrices = 10, gasPrice = 0",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.NewDec(10)
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(0), denom, &testMsg)
				return txBuilder.GetTx()
			},
			false,
			"provided fee < minimum global fee",
			true,
		},
		{
			"invalid cosmos tx with wrong denom",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.NewDec(10)
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), "stake", &testMsg)
				return txBuilder.GetTx()
			},
			false,
			"provided fee < minimum global fee",
			true,
		},
	}

	for _, et := range execTypes {
		for _, tc := range testCases {
			suite.Run(et.name+"_"+tc.name, func() {
				// s.SetupTest(et.isCheckTx)
				ctx := suite.ctx.WithIsReCheckTx(et.isCheckTx)
				dec := cosmosante.NewMinGasPriceDecorator(suite.app.FeeMarketKeeper, suite.app.EvmKeeper)
				_, err := dec.AnteHandle(ctx, tc.malleate(), et.simulate, testutil.NextFn)

				if tc.expPass || (et.simulate && tc.allowPassOnSimulate) {
					suite.Require().NoError(err, tc.name)
				} else {
					suite.Require().Error(err, tc.name)
					suite.Require().Contains(err.Error(), tc.errMsg, tc.name)
				}
			})
		}
	}
}
