build:
	@rm -f node ; go build -o node *.go ; ls -l

run:
	@make build
	@rm -f ~/Library/Application\ Support/TrueBlocks/trueBlocks.toml
	@DEFAULT_CHAIN=mainnet \
		CACHE_PATH="/Users/jrush/Data/buidlguidl/ethereum_clients/trueBlocks/databases/cache" \
		INDEX_PATH="/Users/jrush/Data/buidlguidl/ethereum_clients/trueBlocks/databases/unchained" \
		RPC_MAINNET="http://localhost:23456" \
		./node
