package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChainDescriptors(t *testing.T) {
	tests := []struct {
		name           string
		providerMap    map[string]string
		chainData      string
		expectedOutput string
	}{
		{
			name: "Valid Config with Supported Chains",
			providerMap: map[string]string{
				"mainnet": "http://localhost:8545",
			},
			chainData: `{
  "mainnet": {
    "chain": "mainnet",
    "chainId": "1",
    "remoteExplorer": "https://etherscan.io",
    "symbol": "ETH"
  }
}`,
			expectedOutput: `
  [chains.mainnet]
    chain = "mainnet"
    chainId = "1"
    localExplorer = ""
    remoteExplorer = "https://etherscan.io"
    rpcProvider = "http://localhost:8545"
    symbol = "ETH"
`,
		},
		{
			name: "Unsupported Chain",
			providerMap: map[string]string{
				"unknownchain": "http://localhost:8545",
			},
			chainData: `{
  "mainnet": {
    "chain": "mainnet",
    "chainId": "1",
    "remoteExplorer": "https://etherscan.io",
    "symbol": "ETH"
  }
}`,
			expectedOutput: `
  # unknownchain is not supported
`,
		},
		{
			name: "Empty Chain Data",
			providerMap: map[string]string{
				"mainnet": "http://localhost:8545",
			},
			chainData: "",
			expectedOutput: `
  [chains.mainnet]
    chain = "mainnet"
    chainId = "1"
    localExplorer = ""
    remoteExplorer = "https://etherscan.io"
    rpcProvider = "http://localhost:8545"
    symbol = "ETH"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for the test
			tempDir := t.TempDir()

			// Write the chain data to a temporary file
			tempFile := filepath.Join(tempDir, "chains.json")
			if tt.chainData != "" {
				if err := os.WriteFile(tempFile, []byte(tt.chainData), 0644); err != nil {
					t.Fatalf("failed to write to temp file: %v", err)
				}
			}

			// Create the Config with the temporary directory as ConfigPath
			config := Config{
				ConfigPath:  tempDir,
				ProviderMap: tt.providerMap,
			}

			// Call ChainDescriptors and compare the output
			output := strings.TrimSpace(config.ChainDescriptors())
			expected := strings.TrimSpace(tt.expectedOutput)

			if output != expected {
				t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
			}
		})
	}
}
