package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/joho/godotenv"
)

func (a *App) establishConfig() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = godotenv.Load(filepath.Join(pwd, ".env"))
	if err != nil {
		a.Logger.Warn(err.Error())
		a.Logger.Info("No .env file found, proceed anyway...")
	}

	var ok bool
	if a.Config.ConfigPath, ok = os.LookupEnv("TB_NODE_DATADIR"); !ok {
		return errors.New("TB_NODE_DATADIR is required in the environment")
	} else {
		if a.Config.ConfigPath, err = cleanDataPath(a.Config.ConfigPath); err != nil {
			return err
		}
	}

	chainStr, ok := os.LookupEnv("TB_NODE_CHAINS")
	if !ok {
		chainStr = "mainnet"
	} else {
		chainStr = cleanChainString(chainStr)
	}

	chains := strings.Split(chainStr, ",")
	for _, chain := range chains {
		key := "TB_NODE_" + strings.ToUpper(chain) + "RPC"
		if provider, ok := os.LookupEnv(key); !ok {
			msg := fmt.Sprintf("%s is required in the environment (implied by TB_NODE_CHAINS=%s)", key, chainStr)
			return errors.New(msg)
		} else {
			// TODO: We should hit healthcheck here to make sure the provider is real
			a.Config.ProviderMap[chain] = provider
		}
	}

	// // Set the environment trueblocks-core needs
	os.Setenv("XDG_CONFIG_HOME", a.Config.ConfigPath)
	os.Setenv("TB_SETTINGS_DEFAULTCHAIN", "mainnet")
	os.Setenv("TB_SETTINGS_INDEXPATH", a.Config.IndexPath())
	os.Setenv("TB_SETTINGS_CACHEPATH", a.Config.CachePath())
	for chain, provider := range a.Config.ProviderMap {
		envKey := "TB_CHAINS_" + strings.ToUpper(chain) + "_RPCPROVIDER"
		os.Setenv(envKey, provider)
	}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "TB_") || strings.HasPrefix(env, "XDG_") {
			a.Logger.Debug("", "env", env)
		}
	}

	configFn := filepath.Join(a.Config.ConfigPath, "trueBlocks.toml")
	// if FileExists(configFn) {
	// 	a.Logger.Info("Using existing config", "configFile", configFn)
	// 	return nil
	// }

	if err := EstablishFolder(a.Config.ConfigPath); err != nil {
		return err
	}
	for _, chain := range chains {
		chainConfig := filepath.Join(a.Config.ConfigPath, "config", chain)
		if err := EstablishFolder(chainConfig); err != nil {
			return err
		}
	}

	tmpl, err := template.New("tmpl").Parse(configTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, &a.Config); err != nil {
		return err
	}

	file.StringToAsciiFile(configFn, buf.String())
	a.Logger.Info("Created config file", "configFile", configFn, "config", file.AsciiFileToString(configFn))

	return nil
}

// cleanChainString cleans up the chainStr...no spaces, move 'mainnet' to the front, add it if needed.
func cleanChainString(in string) string {
	out := strings.ReplaceAll(in, " ", "")
	out = strings.ReplaceAll(out, "mainnet", "")
	out = strings.ReplaceAll(out, ",,", ",")
	ret := strings.Trim("mainnet,"+strings.Trim(out, ","), ",")
	return ret
}

// cleanDataPath cleans up the data path, replacing PWD, ~, and HOME with the appropriate values
func cleanDataPath(in string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return in, err
	}
	out := strings.ReplaceAll(in, "PWD", pwd)

	home, err := os.UserHomeDir()
	if err != nil {
		return in, err
	}
	out = strings.ReplaceAll(out, "~", home)
	out = strings.ReplaceAll(out, "HOME", home)
	return filepath.Clean(out), nil
}

var configTmpl string = `[version]
  current = "v3.3.0-release"

[settings]
  cachePath = "{{.CachePath}}"
  defaultChain = "mainnet"
  indexPath = "{{.IndexPath}}"

[keys]
  [keys.etherscan]
    apiKey = ""

[chains]{{.ChainConfigs}}
`
