build:
	@rm -f *node ; go build -o trueblocks-node *.go ; ls -l

run:
	@make build
	@rm -f ~/Library/Application\ Support/TrueBlocks/trueBlocks.toml
	@TB_CHAINS_MAINNET_SCRAPEROUTPUT=/Users/jrush/Data/buidlguidl/ethereum_clients/trueBlocks \
		TB_SETTINGS_CACHEPATH="/Users/jrush/Data/buidlguidl/ethereum_clients/trueBlocks/database/cache" \
		TB_SETTINGS_INDEXPATH="/Users/jrush/Data/buidlguidl/ethereum_clients/trueBlocks/database/unchained" \
		TB_CHAINS_MAINNET_RPCPROVIDER="http://localhost:23456" \
		./trueblocks-node

run-full:
	@make build
	@rm -f ~/Library/Application\ Support/TrueBlocks/trueBlocks.toml
	@TB_CHAINS_MAINNET_SCRAPEROUTPUT=/Users/jrush/Data/buidlguidl/ethereum_clients/trueBlocks \
		TB_SETTINGS_CACHEPATH="/Users/jrush/Data/trueblocks/v1.0.0/cache/" \
		TB_SETTINGS_INDEXPATH="/Users/jrush/Data/trueblocks/v1.0.0/unchained" \
		TB_CHAINS_MAINNET_RPCPROVIDER="http://localhost:23456" \
		./trueblocks-node

# TB_SETTINGS_CACHEPATH="/Users/jrush/Data/trueblocks/v1.0.0/cache/" \
# TB_SETTINGS_INDEXPATH="/Users/jrush/Data/trueblocks/v1.0.0/unchained" \
