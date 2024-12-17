# trueblocks-node

A minimal indexing/monitoring trueblocks node.

## Installation

```bash
go install github.com/TrueBlocks/trueblocks-node/v4@latest
```

This will install the `trueblocks-node` binary in your `GOBIN` directory which is likely in your `$PATH`.

## Testing

```bash
trueblocks-node --version
```

If this doesn't work, check your `$PATH`. No, we won't explain what this means. Google it.

## Running

The easiest way to run this code is to copy `env.example` to `.env` in this folder and edit the values to match your system.

Then run:

```bash
trueblocks-node
```

If you'd rather not do that for whatever reason, try this...

### Alternative

Export two values in the environment and then run the binary:

```bash
export TB_NODE_DATADIR=<your-datadir>
export TB_NODE_MAINNETRPC=<your-rpc-provider>
trueblocks-node
```

These two values are required (even if you're indexing other chains).

For `TB_NODE_DATADIR`, you may use special values for `PWD`, `~`, and `HOME` which are replaced with the obvious values.

For `TB_NODE_MAINNETRPC` use the URL of valid (and fast) RPC provider such as those available from [BlockJoy endpoints](https://www.blockjoy.com/).

You may index other chains by exporting additional values in the environment. Please see the file called `env.example` in this folder for more information.

## Other environment variables

You may increase the logging level of the node by setting the `TB_LOGLEVEL` environment variable. Valid values are `debug`, `info`, `warn`, or `error`. The default is `info`. `debug` is the most verbose.

## Documentation

The documentation includes this README.md file.

Much more detailed documentation (which is derived from the source code), is [available here](https://pkg.go.dev/github.com/TrueBlocks/trueblocks-node/v4).

Internally, `trueblocks-node` uses both the [trueblocks-sdk](https://pkg.go.dev/github.com/TrueBlocks/trueblocks-sdk/v4) and [trueBlocks-core](https://trueblocks.io/chifra/introduction/). In some cases, documentation for these packages may be useful.

## Contributing

We love contributors. Please see information about our [workflow](https://github.com/TrueBlocks/trueblocks-core/blob/develop/docs/BRANCHING.md) before proceeding.

1. Fork this repository into your own repo.
2. Create a branch: `git checkout -b <branch_name>`.
3. Make changes to your local branch and commit them to your forked repo: `git commit -m '<commit_message>'`
4. Push back to the original branch: `git push origin TrueBlocks/trueblocks-core`
5. Create the pull request.

## Contact

If you have questions, comments, or complaints, please join the discussion on our discord server which is [linked from our website](https://trueblocks.io).

## List of Contributors

Thanks to the following people who have contributed to this project:

- [@tjayrush](https://github.com/tjayrush)
- [@dszlachta](https://github.com/dszlachta)
- many others
