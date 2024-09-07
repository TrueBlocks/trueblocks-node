build: *.go
	@rm -f *node ; go build -o trueblocks-node *.go


clean:
	@rm -fR data

run:
	@make build
	@./trueblocks-node --init

install:
	@make build
	@mv trueblocks-node ~/go/bin

test:
	@go test ./...

