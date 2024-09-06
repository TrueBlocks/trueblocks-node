build:
	@rm -f *node ; go build -o trueblocks-node *.go

test:
	@make build
	@rm -fR Data trueBlocks.toml
	@TB_NODE_DATADIR=(pwd)/Data \
		TB_NODE_MAINNETRPC=http://localhost:23456 \
		TB_NODE_CHAIN=sepolia \
		TB_NODE_CHAINRPC=http://localhost:36963 \
		./trueblocks-node --version
