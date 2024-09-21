package app

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-node/v3/utils"
	"github.com/joho/godotenv"
)

func init() {
	if pwd, err := os.Getwd(); err == nil {
		if err = godotenv.Load(filepath.Join(pwd, ".env")); err != nil {
			fmt.Fprintf(os.Stderr, "Found .env, but could not read it\n")
		}
	}
}

// EstablishConfig either reads an existing configuration file or creates it if it doesn't exist.
func (a *App) EstablishConfig() error {
	var ok bool
	var err error
	if a.Config.ConfigPath, ok = os.LookupEnv("TB_NODE_DATADIR"); !ok {
		return errors.New("TB_NODE_DATADIR is required in the environment")
	} else {
		if a.Config.ConfigPath, err = cleanDataPath(a.Config.ConfigPath); err != nil {
			return err
		}
	}
	a.Logger.Info("Using data directory", "dataDir", a.Config.ConfigPath)

	var targets string
	chainStr, ok := os.LookupEnv("TB_NODE_CHAINS")
	if !ok {
		chainStr, targets = "mainnet", "mainnet"
	} else {
		chainStr, targets = cleanChainString(chainStr)
	}
	a.Logger.Debug("cleaned chain string", "chainStr", chainStr, "targets", targets)
	a.Config.Targets = strings.Split(targets, ",")

	chains := strings.Split(chainStr, ",")
	for _, chain := range chains {
		key := "TB_NODE_" + strings.ToUpper(chain) + "RPC"
		if providerUrl, ok := os.LookupEnv(key); !ok {
			msg := fmt.Sprintf("%s is required in the environment (implied by TB_NODE_CHAINS=%s)", key, chainStr)
			return errors.New(msg)
		} else {
			providerUrl = strings.Trim(providerUrl, "/")
			if !isValidURL(providerUrl) {
				return fmt.Errorf("invalid URL for %s: %s", key, providerUrl)
			}
			if err := a.tryConnect(providerUrl, 5); err != nil {
				return err
			} else {
				a.Logger.Info("connected to RPC", "providerUrl", providerUrl)
			}
			a.Config.ProviderMap[chain] = providerUrl
		}
	}

	// // Set the environment trueblocks-core needs
	os.Setenv("XDG_CONFIG_HOME", a.Config.ConfigPath)
	os.Setenv("TB_SETTINGS_DEFAULTCHAIN", "mainnet")
	os.Setenv("TB_SETTINGS_INDEXPATH", a.Config.IndexPath())
	os.Setenv("TB_SETTINGS_CACHEPATH", a.Config.CachePath())
	for chain, providerUrl := range a.Config.ProviderMap {
		envKey := "TB_CHAINS_" + strings.ToUpper(chain) + "_RPCPROVIDER"
		os.Setenv(envKey, providerUrl)
	}

	for i, env := range os.Environ() {
		if strings.HasPrefix(env, "TB_") || strings.HasPrefix(env, "XDG_") {
			a.Logger.Debug(fmt.Sprintf("Env[%d]:", i), "value", env)
		}
	}

	configFn := filepath.Join(a.Config.ConfigPath, "trueBlocks.toml")
	if utils.FileExists(configFn) {
		a.Logger.Info("Using existing config", "configFile", configFn, "nChains", len(a.Config.ProviderMap))
		return nil
	}

	if err := utils.EstablishFolder(a.Config.ConfigPath); err != nil {
		return err
	}
	for _, chain := range chains {
		chainConfig := filepath.Join(a.Config.ConfigPath, "config", chain)
		if err := utils.EstablishFolder(chainConfig); err != nil {
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
	a.Logger.Info("Created config file", "configFile", configFn, "nChains", len(a.Config.ProviderMap))

	return nil
}

func (a *App) tryConnect(providerUrl string, maxAttempts int) error {
	for i := 1; i <= maxAttempts; i++ {
		err := utils.PingRpc(providerUrl)
		if err == nil {
			return nil
		} else {
			a.Logger.Warn("retrying RPC", "provider", providerUrl)
			if i < maxAttempts {
				time.Sleep(1 * time.Second)
			}
		}
	}
	return fmt.Errorf("failed to connect to RPC (%s) after %d attempts", providerUrl, maxAttempts)
}

func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
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

[chains]{{.TmplChain}}
`
