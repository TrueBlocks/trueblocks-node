build:
	@rm -f node ; go build -o node *.go ; ls -l

run:
	@make build
	@TB_SETTINGS_CACHEPATH="/Users/jrush/Data/buidlguidl/ethereum_clients/trueBlocks/database/cache" \
		TB_SETTINGS_INDEXPATH="/Users/jrush/Data/buidlguidl/ethereum_clients/trueBlocks/database/unchained" \
		TB_CHAINS_MAINNET_RPCPROVIDER="http://localhost:23456" \
		./node

# @rm -f ~/Library/Application\ Support/TrueBlocks/trueBlocks.toml
# TB_SETTINGS_DEFAULTCHAIN=mainnet \
# TB_SETTINGS_CACHEPATH="/Users/jrush/Data/trueblocks/v1.0.0/cache/" \
# TB_SETTINGS_INDEXPATH="/Users/jrush/Data/trueblocks/v1.0.0/unchained" \
