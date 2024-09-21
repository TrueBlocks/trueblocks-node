package app

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// isValidChainString return false if the environment variable contains anything other then
// alphanumeric characters or comma, colon, parens, periods, or dashes
func isValidString(input string) bool {
	if strings.Contains(input, "@") {
		return false
	}
	regex := regexp.MustCompile(`^[a-zA-Z0-9,-_ ]+$`)
	return regex.MatchString(input)
}

// splitChainString first validates the chain string then splits it at the commas and trims each item.
func splitChainString(input string) ([]string, error) {
	if ok := isValidString(input); !ok {
		return nil, fmt.Errorf("invalid chain string: %s", input)
	}

	parts := strings.Split(input, ",")
	var cleanedParts []string
	var haveMap = map[string]bool{}
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		for _, r := range trimmedPart {
			if unicode.IsSpace(r) {
				return nil, fmt.Errorf("invalid chain string: part contains whitespace: %s", part)
			}
		}
		if len(trimmedPart) > 0 && !haveMap[trimmedPart] {
			cleanedParts = append(cleanedParts, trimmedPart)
			haveMap[trimmedPart] = true
		}
	}

	return cleanedParts, nil
}

// cleanChainString cleans up the chainStr...no spaces, move 'mainnet' to the front, add it if needed.
func cleanChainString(input string) (string, string, error) {
	if cleaned, err := splitChainString(input); err != nil {
		return "", "", err
	} else {
		chains := []string{"mainnet"}
		for _, chain := range cleaned {
			if chain == "mainnet" {
				continue
			}
			chains = append(chains, chain)
		}
		return strings.Join(chains, ","), strings.Join(cleaned, ","), nil
	}
}
