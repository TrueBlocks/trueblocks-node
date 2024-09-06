package main

import (
	"fmt"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v3"
)

// parseArgs makes sure there are no command line arguments (other than --version or --help)
func (a *App) parseArgs() bool {
	if len(os.Args) < 2 {
		return true
	}

	for _, arg := range os.Args[1:] {
		if arg == "--monitor" {
			a.Monitor = "on"
			return true
		} else {
			if arg == "--version" {
				fmt.Println("trueblocks-node " + sdk.Version())
				return false
			} else {
				if arg != "--help" {
					fmt.Println(colors.Red+"Unknown option:", os.Args[1]+"\n", colors.Off)
				}
				fmt.Println(helpText)
				return false
			}
		}
	}
	return false
}

const helpText = `Usage: trueblocks-node <options>

Options:
---------
 --sync [all*|blooms]: sync to unchained index with blooms only or full index (default: all)
 --api [on*|off]:      enable/disable API server (default: on)
 --monitor [on|off*]:  enable/disable address monitoring (default: off)
 --sleep int:          the number of seconds to sleep between updates (default: 30)
 --version:            display the version string
 --help:               display this help text

Environment:
-------------
You MUST export the following values to the environment:

  TB_NODE_DATADIR:    A directory to store the indexer's data (required, created if necessary)
  TB_NODE_MAINNETRPC: A valid RPC endpoint for Ethereum mainnet (required)

You MAY also export these environment variables:

  TB_NODE_CHAINS:     A comma-separated list of chains to index ("mainnet" is added if not present)
  TB_NODE_<CHAIN>RPC: For each CHAIN in the TB_NODE_CHAINS list, a valid RPC endpoint
                      (example: TB_NODE_SEPOLIARPC=http://localhost:8548)
					  
You may put these values in a .env file in the current folder. See env.example.`
