build: *.go
	@rm -f *node ; go build -o trueblocks-node *.go


clean:
	@rm -fR data

update:
	@go get github.com/TrueBlocks/trueblocks-sdk/v3@latest
	@go get github.com/TrueBlocks/trueblocks-core/src/apps/chifra@latest

run:
	@make build
	@./trueblocks-node --init all

install:
	@make build
	@mv trueblocks-node ~/go/bin

test:
	@go test ./...
