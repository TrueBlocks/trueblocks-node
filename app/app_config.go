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

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-node/v4/utils"
	"github.com/joho/godotenv"
)

func init() {
	if pwd, err := os.Getwd(); err == nil {
		if utils.FileExists(filepath.Join(pwd, ".env")) {
			if err = godotenv.Load(filepath.Join(pwd, ".env")); err != nil {
				fmt.Fprintf(os.Stderr, "Found .env, but could not read it\n")
			}
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
		if chainStr, targets, err = cleanChainString(chainStr); err != nil {
			return err
		}
	}
	a.Logger.Info("cleaned chain string", "chainStr", chainStr, "targets", targets)
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
			if err := a.tryConnect(chain, providerUrl, 5); err != nil {
				return err
			} else {
				a.Logger.Info("connected to RPC", "chain", chain, "providerUrl", providerUrl)
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

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "TB_") || strings.HasPrefix(env, "XDG_") {
			a.Logger.Info("environment", "value", env)
		}
	}

	configFn := filepath.Join(a.Config.ConfigPath, "trueBlocks.toml")
	if utils.FileExists(configFn) {
		a.Logger.Info("Using existing config", "configFile", configFn, "nChains", len(a.Config.ProviderMap))
		// check to make sure the config file has all the chains
		contents := file.AsciiFileToString(configFn)
		for chain := range a.Config.ProviderMap {
			search := "[chains." + chain + "]"
			if !strings.Contains(contents, search) {
				msg := fmt.Sprintf("config file {%s} does not contain {%s}", configFn, search)
				msg = colors.ColoredWith(msg, colors.Red)
				return errors.New(msg)
			}
		}
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

func (a *App) tryConnect(chain, providerUrl string, maxAttempts int) error {
	for i := 1; i <= maxAttempts; i++ {
		err := rpc.PingRpc(providerUrl)
		if err == nil {
			return nil
		} else {
			a.Logger.Warn("retrying RPC", "chain", chain, "provider", providerUrl)
			if i < maxAttempts {
				time.Sleep(1 * time.Second)
			}
		}
	}
	return fmt.Errorf("failed to connect to RPC (%s-%s) after %d attempts", chain, providerUrl, maxAttempts)
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
	ret := filepath.Clean(out)
	if strings.HasSuffix(ret, "/unchained") {
		ret = strings.ReplaceAll(ret, "/unchained", "")
	}
	return ret, nil
}

var configTmpl string = `[version]
  current = "v4.0.0"

[settings]
  cachePath = "{{.CachePath}}"
  defaultChain = "mainnet"
  indexPath = "{{.IndexPath}}"

[keys]
  [keys.etherscan]
    apiKey = ""

[chains]{{.TmplChain}}
`
