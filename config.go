package main

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type Config struct {
	DefaultChain string `json:"defaultChain"`
	IndexPath    string `json:"indexPath"`
	CachePath    string `json:"cachePath"`
	RpcMainnet   string `json:"rpcMainnet"`
}

func loadConfigFromEnv() Config {
	var ok bool
	var ret Config
	if ret.DefaultChain, ok = os.LookupEnv("DEFAULT_CHAIN"); !ok {
		panic("DEFAULT_CHAIN is required")
	}
	if ret.IndexPath, ok = os.LookupEnv("INDEX_PATH"); !ok {
		panic("INDEX_PATH is required")
	}
	if ret.CachePath, ok = os.LookupEnv("CACHE_PATH"); !ok {
		panic("CACHE_PATH is required")
	}
	if ret.RpcMainnet, ok = os.LookupEnv("RPC_MAINNET"); !ok {
		panic("RPC_MAINNET is required")
	}
	return ret
}

func establishConfig() {
	if homeDir, err := os.UserHomeDir(); err != nil {
		panic(err)
	} else {
		configPath := filepath.Join(homeDir, "Library/Application Support/TrueBlocks")
		if err := EstablishFolder((configPath)); err != nil {
			panic(err)
		}

		configFilename := filepath.Join(configPath, "trueBlocks.toml")
		if FileExists(configFilename) {
			logger.Info("Config file already exists. Skipping creation.")
		}

		t, err := template.New("personTemplate").Parse(configTmpl)
		if err != nil {
			panic(err)
		}

		config := loadConfigFromEnv()

		var buf bytes.Buffer
		err = t.Execute(&buf, config)
		if err != nil {
			panic(err)
		}

		file.StringToAsciiFile(configFilename, buf.String())
		logger.Info("Created config file at: ", configFilename)
	}
}

var configTmpl string = `[version]
  current = "v3.1.3-release"

[settings]
  cachePath = "{{.CachePath}}"
  defaultChain = "{{.DefaultChain}}"
  indexPath = "{{.IndexPath}}"

[keys]
  [keys.etherscan]
    apiKey = ""

[chains]
  [chains.mainnet]
    chain = "mainnet"
    chainId = "1"
    remoteExplorer = "https://etherscan.io"
    rpcProvider = "{{.RpcMainnet}}"
    symbol = "ETH"
`

// EstablishFolder creates folders given a list of folders
func EstablishFolder(rootPath string) error {
	_, err := os.Stat(rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(rootPath, 0755)
			if err != nil {
				return err
			}
		} else {
			// If there's an error other than not exist...we fail
			return err
		}
	}
	return nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
