package main

import (
	"fmt"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v3"
)

// parseArgs makes sure there are no command line arguments (other than --version or --help)
func parseArgs() bool {
	if len(os.Args) < 2 {
		return true
	}

	if os.Args[1] == "--version" {
		fmt.Println("trueblocks-node " + sdk.Version())
	} else {
		if os.Args[1] != "--help" {
			fmt.Println(colors.Red+"Unknown option:", os.Args[1]+"\n", colors.Off)
		}
		fmt.Println(helpText)
	}

	return false
}

const helpText = `Usage: trueblocks-node <options>

Options:
---------
 --version: display the version string
 --help:    display this help text

Environment:
-------------
You MUST export the following value to the environment:

  TB_NODE_DATADIR:    A directory to store data (created if it does not exist)
  TB_NODE_MAINNETRPC: A valid RPC endpoint for Ethereum mainnet

You MAY also export these environment variables:

  TB_NODE_CHAIN:    The name of the chain to index if not "mainnet"
  TB_NODE_CHAINRPC: An RPC endpoint running that chain's RPC endpoint`
