package app

import (
	"fmt"
	"os"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
)

// ParseArgs parses the command line options and sets the app's configuration accordingly. See README.md or run trueblocks-node --help.
func (a *App) ParseArgs() (bool, error) {
	if len(os.Args) < 2 {
		return true, nil
	}

	hasValue := func(i int) bool {
		return i+1 < len(os.Args) && os.Args[i+1][0] != '-'
	}

	validateMode := func(value string) (InitMode, error) {
		mode := InitMode(value)
		switch mode {
		case All, Blooms, None:
			return mode, nil
		}
		return mode, fmt.Errorf("invalid value: " + value)
	}

	validateOnOff := func(value string) (OnOff, error) {
		mode := OnOff(value)
		switch mode {
		case On, Off:
			return mode, nil
		}
		return mode, fmt.Errorf("invalid value: " + value)
	}

	validateSleep := func(value string) (int, error) {
		sleep := 1
		if _, err := fmt.Sscanf(value, "%d", &sleep); err != nil || sleep < 1 {
			return 1, fmt.Errorf("invalid value for --sleep: " + value)
		}
		return sleep, nil
	}

	var err error
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "--init":
			if !hasValue(i) {
				return true, fmt.Errorf("missing argument for --init")
			} else {
				if a.InitMode, err = validateMode(os.Args[i+1]); err != nil {
					return true, err
				}
				i++
			}
		case "--api":
			if !hasValue(i) {
				return true, fmt.Errorf("missing argument for --init")
			} else {
				if a.Api, err = validateOnOff(os.Args[i+1]); err != nil {
					return true, err
				}
				i++
			}
		case "--monitor":
			if !hasValue(i) {
				return true, fmt.Errorf("missing argument for --monitor")
			} else {
				if a.Monitor, err = validateOnOff(os.Args[i+1]); err != nil {
					return true, err
				}
				i++
			}
		case "--sleep":
			if !hasValue(i) {
				return true, fmt.Errorf("missing argument for --sleep")
			} else {
				if a.Sleep, err = validateSleep(os.Args[i+1]); err != nil {
					return true, err
				}
				i++
			}
		case "--version":
			fmt.Println("trueblocks-node " + sdk.Version())
			return false, nil
		default:
			if arg != "--help" {
				err = fmt.Errorf("%s%s", "unknown option:"+os.Args[1]+"\n", helpText)
			} else {
				err = fmt.Errorf("%s", helpText)
			}
			return true, err
		}
	}
	return true, nil
}

const helpText = `Usage: trueblocks-node <options>

Options:
---------
 --init [all*|blooms|none]: download from the unchained index smart contract (default: all)
 --api [on*|off]:           enable/disable API server (default: on)
 --monitor [on|off*]:       enable/disable address monitoring (currently disabled, default: off)
 --sleep int:               the number of seconds to sleep between updates (default: 30)
 --version:                 display the version string
 --help:                    display this help text

Environment:
-------------
You MUST export the following values to the environment:

  TB_NODE_DATADIR:    A directory to store the indexer's data (required, created if necessary)
  TB_NODE_MAINNETRPC: A valid RPC endpoint for Ethereum mainnet (required)

You MAY also export these environment variables:

  TB_NODE_CHAINS:     A comma-separated list of chains to index (default: "mainnet")
  TB_NODE_<CHAIN>RPC: For each CHAIN in the TB_NODE_CHAINS list, a valid RPC endpoint
                      (example: TB_NODE_SEPOLIARPC=http://localhost:8548)
					  
You may put these values in a .env file in the current folder. See env.example.`
