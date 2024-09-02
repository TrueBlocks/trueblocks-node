package main

import (
	"bytes"
	"encoding/json"
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
	OutputPath   string `json:"outputPath"`
}

func (c *Config) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func loadConfigFromEnv() Config {
	var ok bool
	var ret Config
	if ret.DefaultChain, ok = os.LookupEnv("TB_SETTINGS_DEFAULTCHAIN"); !ok {
		ret.DefaultChain = "mainnet"
		// panic("TB_SETTINGS_DEFAULTCHAIN is required in the environment")
	}
	if ret.IndexPath, ok = os.LookupEnv("TB_SETTINGS_INDEXPATH"); !ok {
		panic("TB_SETTINGS_INDEXPATH is required in the environment")
	}
	if ret.CachePath, ok = os.LookupEnv("TB_SETTINGS_CACHEPATH"); !ok {
		panic("TB_SETTINGS_CACHEPATH is required in the environment")
	}
	if ret.RpcMainnet, ok = os.LookupEnv("TB_CHAINS_MAINNET_RPCPROVIDER"); !ok {
		panic("TB_CHAINS_MAINNET_RPCPROVIDER is required in the environment")
	}
	if ret.OutputPath, ok = os.LookupEnv("TB_CHAINS_MAINNET_SCRAPEROUTPUT"); !ok {
		ret.OutputPath, _ = os.Getwd()
		// panic("TB_CHAINS_MAINNET_SCRAPEROUTPUT is required in the environment")
	}

	return ret
}

func establishConfig() (*Config, error) {
	if homeDir, err := os.UserHomeDir(); err != nil {
		return nil, err

	} else {
		configPath := filepath.Join(homeDir, "Library/Application Support/TrueBlocks")
		if err := EstablishFolder((configPath)); err != nil {
			return nil, err
		}

		config := loadConfigFromEnv()

		configFilename := filepath.Join(configPath, "trueBlocks.toml")
		if FileExists(configFilename) {
			logger.Info("Config file already exists. Using existing config.")
			return &config, nil
		}

		t, err := template.New("theTemplate").Parse(configTmpl)
		if err != nil {
			return &config, err
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, config)
		if err != nil {
			return &config, err
		}

		file.StringToAsciiFile(configFilename, buf.String())
		logger.Info("Created config file at: ", configFilename)
		return &config, err
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
