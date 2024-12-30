package app

import (
	"testing"
)

func TestSplitChainString(t *testing.T) {
	tests := []struct {
		input       string
		expected    []string
		expectError bool
	}{
		{"mainnet,testnet", []string{"mainnet", "testnet"}, false},
		{"testnet , mainnet", []string{"testnet", "mainnet"}, false},
		{"mainnet, mainnet-testnet", []string{"mainnet", "mainnet-testnet"}, false},
		{" ,mainnet, , , ", []string{"mainnet"}, false},
		{"mainnet@", nil, true},
		{"@@@@", nil, true},
		{"mainnet,test@", nil, true},
		{"mainnet,,test!net", nil, true}, // Invalid special character
		{"mainnet, test net", nil, true}, // Internal space in a chain
		{"", nil, true},                  // Empty input
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := splitChainString(test.input)
			if (err != nil) != test.expectError {
				t.Errorf("splitChainString(%q) error = %v; want error %v", test.input, err, test.expectError)
			}
			if !test.expectError && !equalStringSlices(result, test.expected) {
				t.Errorf("splitChainString(%q) = %v; want %v", test.input, result, test.expected)
			}
		})
	}
}

func TestCleanChainString(t *testing.T) {
	tests := []struct {
		input           string
		expectedChains  string
		expectedTargets string
		expectError     bool
	}{
		{"mainnet", "mainnet", "mainnet", false},
		{"testnet,mainnet", "mainnet,testnet", "testnet,mainnet", false},
		{"testnet,othernet", "mainnet,testnet,othernet", "testnet,othernet", false},
		{" ,mainnet, , , ", "mainnet", "mainnet", false},
		{"mainnet,testnet,mainnet", "mainnet,testnet", "mainnet,testnet", false},
		{"othernet,,testnet,mainnet", "mainnet,othernet,testnet", "othernet,testnet,mainnet", false},
		{"gnosis#", "", "", true},
		{"mainnet,test#net", "", "", true},
		{"mainnet, test net", "", "", true},
		{"mainnet,,test!net", "", "", true},
		{"mainnet&chain", "", "", true},
		{"_testnet", "mainnet,_testnet", "_testnet", false},
		{"-testnet", "mainnet,-testnet", "-testnet", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			chains, targets, err := cleanChainString(test.input)
			if (err != nil) != test.expectError {
				t.Errorf("cleanChainString(%q) error = %v; want error %v", test.input, err, test.expectError)
			}

			if chains != test.expectedChains {
				t.Errorf("cleanChainString(%q) = %q; want %q for chains", test.input, chains, test.expectedChains)
			}

			if targets != test.expectedTargets {
				t.Errorf("cleanChainString(%q) = %q; want %q for targets", test.input, targets, test.expectedTargets)
			}
		})
	}
}

// Helper function to compare two slices of strings
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
