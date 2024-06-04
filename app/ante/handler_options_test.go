package ante_test

import (
	"github.com/slandymani/evm-module/app"
	"github.com/slandymani/evm-module/app/ante"
	ethante "github.com/slandymani/evm-module/app/ante/evm"
	"github.com/slandymani/evm-module/encoding"
	"github.com/slandymani/evm-module/types"
)

func (suite *AnteTestSuite) TestValidateHandlerOptions() {
	cases := []struct {
		name    string
		options ante.HandlerOptions
		expPass bool
	}{
		{
			"fail - empty options",
			ante.HandlerOptions{},
			false,
		},
		{
			"fail - empty account keeper",
			ante.HandlerOptions{
				Cdc:           suite.app.AppCodec(),
				AccountKeeper: nil,
			},
			false,
		},
		{
			"fail - empty bank keeper",
			ante.HandlerOptions{
				Cdc:           suite.app.AppCodec(),
				AccountKeeper: suite.app.AccountKeeper,
				BankKeeper:    nil,
			},
			false,
		},
		{
			"fail - empty distribution keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: nil,

				IBCKeeper: nil,
			},
			false,
		},
		{
			"fail - empty IBC keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,

				IBCKeeper: nil,
			},
			false,
		},
		{
			"fail - empty staking keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				StakingKeeper: nil,
			},
			false,
		},
		{
			"fail - empty fee market keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,

				StakingKeeper:   &suite.app.StakingKeeper,
				FeeMarketKeeper: nil,
			},
			false,
		},
		{
			"fail - empty EVM keeper",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				StakingKeeper:      &suite.app.StakingKeeper,
				FeeMarketKeeper:    suite.app.FeeMarketKeeper,
				EvmKeeper:          nil,
			},
			false,
		},
		{
			"fail - empty signature gas consumer",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				StakingKeeper:      &suite.app.StakingKeeper,
				FeeMarketKeeper:    suite.app.FeeMarketKeeper,
				EvmKeeper:          suite.app.EvmKeeper,
				SigGasConsumer:     nil,
			},
			false,
		},
		{
			"fail - empty signature mode handler",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				StakingKeeper:      &suite.app.StakingKeeper,
				FeeMarketKeeper:    suite.app.FeeMarketKeeper,
				EvmKeeper:          suite.app.EvmKeeper,
				SigGasConsumer:     ante.SigVerificationGasConsumer,
				SignModeHandler:    nil,
			},
			false,
		},
		{
			"fail - empty tx fee checker",
			ante.HandlerOptions{
				Cdc:                suite.app.AppCodec(),
				AccountKeeper:      suite.app.AccountKeeper,
				BankKeeper:         suite.app.BankKeeper,
				DistributionKeeper: suite.app.DistrKeeper,
				StakingKeeper:      &suite.app.StakingKeeper,
				FeeMarketKeeper:    suite.app.FeeMarketKeeper,
				EvmKeeper:          suite.app.EvmKeeper,
				SigGasConsumer:     ante.SigVerificationGasConsumer,
				SignModeHandler:    suite.app.GetTxConfig().SignModeHandler(),
				TxFeeChecker:       nil,
			},
			false,
		},
		{
			"success - default app options",
			ante.HandlerOptions{
				Cdc:                    suite.app.AppCodec(),
				AccountKeeper:          suite.app.AccountKeeper,
				BankKeeper:             suite.app.BankKeeper,
				DistributionKeeper:     suite.app.DistrKeeper,
				ExtensionOptionChecker: types.HasDynamicFeeExtensionOption,
				EvmKeeper:              suite.app.EvmKeeper,
				StakingKeeper:          &suite.app.StakingKeeper,
				FeegrantKeeper:         suite.app.FeeGrantKeeper,
				FeeMarketKeeper:        suite.app.FeeMarketKeeper,
				SignModeHandler:        encoding.MakeConfig(app.ModuleBasics).TxConfig.SignModeHandler(),
				SigGasConsumer:         ante.SigVerificationGasConsumer,
				MaxTxGasWanted:         40000000,
				TxFeeChecker:           ethante.NewDynamicFeeChecker(suite.app.EvmKeeper),
			},
			true,
		},
	}

	for _, tc := range cases {
		err := tc.options.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
