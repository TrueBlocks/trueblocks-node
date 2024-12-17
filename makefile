#-------------------------------------------------
bin=../bin

#-------------------------------------------------
exec=node
dest=$(bin)/$(exec)

#-------------------------------------------------
all:
	@make app

every:
	@cd ../build ; make ; cd -
	@make app

app: *.go app/*.go config/*.go
	@echo Building trueblocks-node...
	@mkdir -p $(bin)
	@go build -o $(dest) *.go

update:
	@go get "github.com/TrueBlocks/trueblocks-sdk/v4@latest"
	@go get github.com/TrueBlocks/trueblocks-core/src/apps/chifra@latest

install:
	@make build
	@mv trueblocks-node ~/go/bin

test:
	@go test ./...

#-------------------------------------------------
clean:
	-@$(RM) -f $(dest)
