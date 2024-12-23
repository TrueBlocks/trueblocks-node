package app

import (
	"errors"
	"fmt"
	"os"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
)

var (
	ErrMissingArgument = errors.New("missing argument")
	ErrInvalidValue    = errors.New("invalid value")
)

// ParseArgs parses the command line options and sets the app's configuration accordingly. See README.md or run trueblocks-node --help.
func (a *App) ParseArgs() (bool, error) {
	if len(os.Args) < 2 {
		return true, nil
	}

	hasValue := func(i int) bool {
		return i+1 < len(os.Args) && os.Args[i+1][0] != '-'
	}

	handleInit := func(i int) (int, error) {
		if hasValue(i) {
			if mode, err := validateMode(os.Args[i+1]); err == nil {
				a.InitMode = mode
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --init: %w", err)
			}
		}
		return i, fmt.Errorf("%w for --init", ErrMissingArgument)
	}

	handleScrape := func(i int) (int, error) {
		if hasValue(i) {
			if mode, err := validateOnOff(os.Args[i+1]); err == nil {
				a.Scrape = mode
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --scrape: %w", err)
			}
		}
		return i, fmt.Errorf("%w for --scrape", ErrMissingArgument)
	}

	handleApi := func(i int) (int, error) {
		if hasValue(i) {
			if mode, err := validateOnOff(os.Args[i+1]); err == nil {
				a.Api = mode
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --api: %w", err)
			}
		}
		return i, fmt.Errorf("%w for --api", ErrMissingArgument)
	}

	handleIpfs := func(i int) (int, error) {
		if hasValue(i) {
			if mode, err := validateOnOff(os.Args[i+1]); err == nil {
				a.Ipfs = mode
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --ipfs: %w", err)
			}
		}
		return i, fmt.Errorf("%w for --ipfs", ErrMissingArgument)
	}

	handleMonitor := func(i int) (int, error) {
		if hasValue(i) {
			if mode, err := validateOnOff(os.Args[i+1]); err == nil {
				a.Monitor = mode
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --monitor: %w", err)
			}
		}
		return i, fmt.Errorf("%w for --monitor", ErrMissingArgument)
	}

	handleSleep := func(i int) (int, error) {
		if hasValue(i) {
			if sleep, err := validateSleep(os.Args[i+1]); err == nil {
				a.Sleep = sleep
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --sleep: %w", err)
			}
		}
		return i, fmt.Errorf("%w for --sleep", ErrMissingArgument)
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		var err error
		switch arg {
		case "--init":
			i, err = handleInit(i)
		case "--scrape":
			i, err = handleScrape(i)
		case "--api":
			i, err = handleApi(i)
		case "--ipfs":
			i, err = handleIpfs(i)
		case "--monitor":
			i, err = handleMonitor(i)
		case "--sleep":
			i, err = handleSleep(i)
		case "--version":
			a.Logger.Info("trueblocks-node " + sdk.Version())
			return false, nil
		default:
			if arg != "--help" {
				return true, fmt.Errorf("unknown option:%s\n%s", os.Args[i], helpText)
			}
			fmt.Printf("%s\n", helpText)
			return false, nil
		}
		if err != nil {
			return true, err
		}
	}
	return true, nil
}

func validateEnum[T ~string](value T, validOptions []T, name string) (T, error) {
	for _, option := range validOptions {
		if value == option {
			return value, nil
		}
	}
	return value, fmt.Errorf("invalid value for %s: %s", name, value)
}

func validateMode(value string) (InitMode, error) {
	return validateEnum(InitMode(value), []InitMode{All, Blooms, None}, "mode")
}

func validateOnOff(value string) (OnOff, error) {
	return validateEnum(OnOff(value), []OnOff{On, Off}, "onOff")
}

func validateSleep(value string) (int, error) {
	var sleep int
	if _, err := fmt.Sscanf(value, "%d", &sleep); err != nil || sleep < 1 {
		return 1, fmt.Errorf("invalid value for sleep: %s", value)
	}
	return sleep, nil
}

const helpText = `Usage: trueblocks-node <options>

Options:
---------
 --init     [all*|blooms|none]   download from the unchained index smart contract (default: all)
 --scrape   [on*|off]            enable/disable the Unchained Index scraper (default: off)
 --api      [on*|off]            enable/disable API server (default: on)
 --ipfs     [on*|off]            enable/disable IPFS daemon (default: off)
 --monitor  [on|off*]            enable/disable address monitoring (currently disabled, default: off)
 --sleep    int                  the number of seconds to sleep between updates (default: 30)
 --version                       display the version string
 --help                          display this help text

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
