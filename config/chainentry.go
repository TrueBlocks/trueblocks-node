package config

import (
	"bytes"
	"encoding/json"
	"html/template"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

// TmplEntry represents the configuration of a single chain in the configurion file's template.
type TmplEntry struct {
	Chain          string `json:"chain"`
	ChainId        string `json:"chainId"`
	LocalExplorer  string `json:"localExplorer"`
	RemoteExplorer string `json:"remoteExplorer"`
	RpcProvider    string `json:"rpcProvider"`
	Symbol         string `json:"symbol"`
}

// TmplChain is used used to establish the app's config file (with a GoLang template) if it doesn't
// exist. Note this is a shortened form of the trueblocks-core's configuration file which will be
// used if it already exists.
func (c *Config) TmplChain() string {
	dataFn := filepath.Join("chains.json")
	chainData := file.AsciiFileToString(dataFn)
	if !file.FileExists(dataFn) || len(chainData) == 0 {
		chainData = `{
  "mainnet": {
    "chain": "mainnet",
    "chainId": "1",
    "remoteExplorer": "https://etherscan.io",
    "symbol": "ETH"
  }
}
`
	}

	chains := make(map[string]TmplEntry)
	if err := json.Unmarshal([]byte(chainData), &chains); err != nil {
		return err.Error()
	}

	tmpl, err := template.New("chainConfigTmpl").Parse(`  [chains.{{.Chain}}]
    chain = "{{.Chain}}"
    chainId = "{{.ChainId}}"
    localExplorer = "{{.LocalExplorer}}"
    remoteExplorer = "{{.RemoteExplorer}}"
    rpcProvider = "{{.RpcProvider}}"
    symbol = "{{.Symbol}}"`)
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
