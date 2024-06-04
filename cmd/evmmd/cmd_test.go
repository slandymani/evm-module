//go:build ignore

package main_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/stretchr/testify/require"

	"github.com/slandymani/evm-module/app"
	evmmd "github.com/slandymani/evm-module/cmd/evmmd"
	"github.com/slandymani/evm-module/utils"
)

func TestInitCmd(t *testing.T) {
	rootCmd, _ := evmmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",      // Test the init cmd
		"evmm-test", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
		fmt.Sprintf("--%s=%s", flags.FlagChainID, utils.TestEdge2ChainID+"-3"),
	})

	err := svrcmd.Execute(rootCmd, "evmmd", app.DefaultNodeHome)
	require.NoError(t, err)
}

func TestAddKeyLedgerCmd(t *testing.T) {
	rootCmd, _ := evmmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"keys",
		"add",
		"dev0",
		fmt.Sprintf("--%s", flags.FlagUseLedger),
	})

	err := svrcmd.Execute(rootCmd, "evmmd", app.DefaultNodeHome)
	require.Error(t, err)
}
