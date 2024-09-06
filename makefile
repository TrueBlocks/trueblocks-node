build:
	@rm -f *node ; go build -o trueblocks-node *.go


clean:
	@rm -fR data

run:
	@make build
	@./trueblocks-node

install:
	@make build
	@mv trueblocks-node ~/go/bin
