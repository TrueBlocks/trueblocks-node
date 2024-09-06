# trueblocks-node

A minimal indexing/monitoring trueblocks node.

## Installation

```bash
go install github.com/TrueBlocks/trueblocks-node/v3@latest
```

This will install the `trueblocks-node` binary in your `GOBIN` directory.

## Testing

```bash
trueblocks-node --version
```

## Running

You MUST export two values in the environment before running the node.

```bash
export TB_NODE_DATADIR=<your-datadir>
export TB_NODE_MAINNETRPC=<your-rpc-provider>
trueblocks-node
```

These two values are always required. (Even if you're indexing another chain.)

## Other Options

There are other options and environment variable available that allow you to to index other chains and specify settings. Run

```[bash]
trueblocks-node --help
```

for more information.
