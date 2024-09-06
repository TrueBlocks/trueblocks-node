package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

type Config struct {
	ConfigFolder    string `json:"configFolder"`
	DefaultChain    string `json:"defaultChain"`
	MainnetProvider string `json:"mainnetProvider"`
	ChainProvider   string `json:"chainProvider"`
}

func (c *Config) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c *Config) CachePath() string {
	return filepath.Join(c.ConfigFolder, "cache")
}

func (c *Config) IndexPath() string {
	return filepath.Join(c.ConfigFolder, "unchained")
}

func (c *Config) IsMainnet() bool {
	return c.DefaultChain == "mainnet"
}

/*
TB_NODE_MAINNETRPC: A valid RPC endpoint for Ethereum mainnet
TB_NODE_CHAIN:    The name of the chain to index if not "mainnet"
TB_NODE_CHAINRPC: An RPC endpoint running that chain's RPC endpoint
*/

func (a *App) establishConfig() error {
	var ok bool
	if a.Config.ConfigFolder, ok = os.LookupEnv("TB_NODE_DATADIR"); !ok {
		return errors.New("TB_NODE_DATADIR is required in the environment")
	}
	os.Setenv("XDG_CONFIG_HOME", a.Config.ConfigFolder)

	if a.Config.MainnetProvider, ok = os.LookupEnv("TB_NODE_MAINNETRPC"); !ok {
		return errors.New("TB_NODE_MAINNETRPC is required in the environment")
	}

	if a.Config.DefaultChain, ok = os.LookupEnv("TB_NODE_CHAIN"); !ok || a.Config.DefaultChain == "mainnet" {
		a.Config.DefaultChain = "mainnet"
		a.Config.ChainProvider = a.Config.MainnetProvider
	} else {
		if a.Config.ChainProvider, ok = os.LookupEnv("TB_NODE_CHAINRPC"); !ok {
			return errors.New("if TB_NODE_CHAIN is not empty, TB_NODE_CHAINRPC is required in the environment")
		}
	}

	// Set the environment trueblocks-core needs
	chain := strings.ToUpper(a.Config.DefaultChain)
	os.Setenv("TB_SETTINGS_DEFAULTCHAIN", a.Config.DefaultChain)
	os.Setenv("TB_SETTINGS_INDEXPATH", a.Config.IndexPath())
	os.Setenv("TB_SETTINGS_CACHEPATH", a.Config.CachePath())
	os.Setenv("TB_CHAINS_MAINNET_RPCPROVIDER", a.Config.MainnetProvider)
	if !a.Config.IsMainnet() {
		os.Setenv("TB_CHAINS_"+chain+"_RPCPROVIDER", a.Config.ChainProvider)
	}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "TB_") {
			a.Logger.Info("", "env", env)
		}
	}

	if err := EstablishFolder(a.Config.ConfigFolder); err != nil {
		return err
	}
	mainnetConfig := filepath.Join(a.Config.ConfigFolder, "config", "mainnet")
	if err := EstablishFolder(mainnetConfig); err != nil {
		return err
	}
	chainConfig := filepath.Join(a.Config.ConfigFolder, "config", a.Config.DefaultChain)
	if err := EstablishFolder(chainConfig); err != nil {
		return err
	}

	configFn := filepath.Join(a.Config.ConfigFolder, "trueBlocks.toml")
	if FileExists(configFn) {
		a.Logger.Info("Using existing config", "configFile", configFn)
		return nil
	}

	tmpl, err := template.New("tmpl").Parse(configTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &a.Config)
	if err != nil {
		return err
	}

	file.StringToAsciiFile(configFn, buf.String())
	a.Logger.Info("Created config file", "configFile", configFn, "config", a.Config.String())

	return err
}

var configTmpl string = `[version]
  current = "v3.3.0-release"

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
    rpcProvider = "{{.MainnetProvider}}"
    symbol = "ETH"{{if ne .DefaultChain "mainnet"}}
  [chains.{{.DefaultChain}}]
    chain = "{{.DefaultChain}}"
    chainId = ""
    remoteExplorer = ""
    rpcProvider = "{{.ChainProvider}}"
    symbol = ""{{end}}
`
