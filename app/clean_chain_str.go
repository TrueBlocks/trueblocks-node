package app

import (
	"strings"
)

// cleanChainString cleans up the chainStr...no spaces, move 'mainnet' to the front, add it if needed.
func cleanChainString(in string) (string, string) {
	scrapeTargets := strings.ReplaceAll(in, " ", "")
	scrapeTargets = strings.Trim(scrapeTargets, ",")
	scrapeTargets = strings.ReplaceAll(scrapeTargets, ",,", ",")

	chains := scrapeTargets
	if strings.Contains(chains, "mainnet") {
		chains = strings.ReplaceAll(chains, "mainnet", "")
		chains = strings.ReplaceAll(chains, ",,", ",")
	}
	chains = strings.Trim("mainnet,"+chains, ",")

	return chains, scrapeTargets
}
