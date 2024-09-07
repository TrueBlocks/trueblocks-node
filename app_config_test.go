package main

import (
	"testing"
)

// func TestMain(t *testing.T) {
// 	fmt.Println(len(os.Args), os.Args)
// 	os.Args = os.Args[:1]
// 	fmt.Println(len(os.Args), os.Args)
// 	main()
// }

func TestCleanChainString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"mainnet", "mainnet"},
		{"  mainnet  ", "mainnet"},
		{"mainnet,", "mainnet"},
		{"  mainnet  ,", "mainnet"},
		{"testnet", "mainnet,testnet"},
		{"mainnet, testnet", "mainnet,testnet"},
		{"  mainnet, testnet  ", "mainnet,testnet"},
		{"  mainnet, , testnet  ", "mainnet,testnet"},
		{"mainnet,testnet,mainnet", "mainnet,testnet"},
		{"mainnet,testnet,,othernet", "mainnet,testnet,othernet"},
		{"mainnet,mainnet", "mainnet"},
		{" ,mainnet, , , ", "mainnet"},
		{"  ,, ", "mainnet"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, _ := cleanChainString(tt.input)
			if result != tt.expected {
				t.Errorf("cleanChainString(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
