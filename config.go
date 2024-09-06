package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

type Config struct {
	ConfigPath  string            `json:"configPath"`
	ProviderMap map[string]string `json:"providers"` // chain to provider
}

func (c *Config) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c *Config) CachePath() string {
	return filepath.Join(c.ConfigPath, "cache")
}

func (c *Config) IndexPath() string {
	return filepath.Join(c.ConfigPath, "unchained")
}

type ChainConfig struct {
	Chain          string `json:"chain"`
	ChainId        string `json:"chainId"`
	LocalExplorer  string `json:"localExplorer"`
	RemoteExplorer string `json:"remoteExplorer"`
	RpcProvider    string `json:"rpcProvider"`
	Symbol         string `json:"symbol"`
}

func (c *Config) ChainConfigs() string {
	chainData := file.AsciiFileToString("chains.json")
	if len(chainData) == 0 {
		chainData = `{
	"mainnet": {
		"chain": "mainnet",
		"chainId": "1",
		"remoteExplorer": "https://etherscan.io",
		"symbol": "ETH"
	},
	"sepolia": {
		"chain": "sepolia",
		"chainId": "11155111",
		"remoteExplorer": "https://sepolia.otterscan.io/",
		"symbol": "ETH"
	},
	"optimism": {
		"chain": "optimism",
		"chainId": "10",
		"remoteExplorer": "https://optimistic.etherscan.io",
		"symbol": "ETH"
	},
	"gnosis": {
		"chain": "gnosis",
		"chainId": "100",
		"remoteExplorer": "https://gnosisscan.io/",
		"symbol": "xDAI"
	}
}`
	}

	chains := make(map[string]ChainConfig)
	if err := json.Unmarshal([]byte(chainData), &chains); err != nil {
		return err.Error()
	}

	tmpl, err := template.New("chainConfigTmpl").Parse(chainConfigTmpl)
	if err != nil {
		return err.Error()
	}

	ret := []string{}
	for chain, provider := range c.ProviderMap {
		if chainConfig, ok := chains[chain]; ok {
			chainConfig.RpcProvider = provider

			var buf bytes.Buffer
			if err = tmpl.Execute(&buf, &chainConfig); err != nil {
				return err.Error()
			}

			ret = append(ret, buf.String())
		} else {
			ret = append(ret, "  # "+chain+" is not supported")
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return strings.Compare(ret[i], ret[j]) < 0
	})

	return "\n" + strings.Join(ret, "\n")
}

var chainConfigTmpl = `  [chains.{{.Chain}}]
    chain = "{{.Chain}}"
    chainId = "{{.ChainId}}"
    localExplorer = "{{.LocalExplorer}}"
    remoteExplorer = "{{.RemoteExplorer}}"
    rpcProvider = "{{.RpcProvider}}"
    symbol = "{{.Symbol}}"`
