package app

import (
	"fmt"
	"testing"
)

func TestCleanChainString(t *testing.T) {
	tests := []struct {
		input           string
		expectedChains  string
		expectedTargets string
		expectedError   bool
	}{
		{"mainnet", "mainnet", "mainnet", false},
		{"  mainnet  ", "mainnet", "mainnet", false},
		{"gnosis#", "", "", true},
		{"mainnet,testnet", "mainnet,testnet", "mainnet,testnet", false},
		{"testnet , mainnet", "mainnet,testnet", "testnet,mainnet", false},
		{"mainnet,testnet,mainnet", "mainnet,testnet", "mainnet,testnet", false},
		{"othernet,mainnet,testnet,,othernet", "mainnet,othernet,testnet", "othernet,mainnet,testnet", false},
		{" ,mainnet, , , ", "mainnet", "mainnet", false},
		{"  ,, ", "mainnet", "", false},
		{"mainnet-testnet_other", "mainnet,mainnet-testnet_other", "mainnet-testnet_other", false},
		{"@invalid", "", "", true},          // Invalid due to special character
		{"mainnet,test#net", "", "", true},  // Invalid due to special character
		{"mainnet, test net", "", "", true}, // Invalid due to internal space within a field
		{"mainnet,,test!net", "", "", true}, // Invalid due to special character in the field
		{"mainnet&chain", "", "", true},     // Invalid due to special character '@'
		{"_testnet", "mainnet,_testnet", "_testnet", false},
		{"-testnet", "mainnet,-testnet", "-testnet", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			chains, targets, err := cleanChainString(tt.input)
			fmt.Printf("[%s]-[%s]-[%s]-[%v]\n", tt.input, chains, targets, err)
			if (err != nil) != tt.expectedError {
				t.Errorf("cleanChainString(%q) error = %v; want error %v", tt.input, err, tt.expectedError)
			}

			if chains != tt.expectedChains {
				t.Errorf("cleanChainString(%q) = %q; want %q for chains", tt.input, chains, tt.expectedChains)
			}

			if targets != tt.expectedTargets {
				t.Errorf("cleanChainString(%q) = %q; want %q for targets", tt.input, targets, tt.expectedTargets)
			}
		})
	}
}
