/*
trueBlocks-node serves a number of purposes:

- actively index one or more EVM blockchains depending on the configuration,
- pin new portions of the Unchained Index to IPFS,
- publish the Unchained Index's manifest hash to the smart contract,
- monitor transactions for a given set of user-defined addresses, and
- expose a REST Api to all of TrueBlocks' `chifra` commands,

The program requires a configuration file, `.env`, to be present in the same directory
from which it is run. An example file (whose contents are copied below) may be
found in the file called `env.example` at the root of this repo.

#
# Rename this file to .env in the folder from which you run trueblocks-node
#
# The full path of where you want the index data to reside.
#
# You may use special values for "~", "HOME", and "PWD" to represent the
# home directory or the current working directory. The default is "PWD/data".
TB_NODE_DATADIR="PWD/data"

# A comma separated list of chains to index (valid values are any combination of "mainnet",
# "sepolia", "gnosis", "optimism"). "mainnet" is required and will be added if not present.
TB_NODE_CHAINS="mainnet"
#TB_NODE_CHAINS="mainnet,sepolia"

# for each chain, an RPC provider for that chain is required. "MAINNET" is required.
TB_NODE_MAINNETRPC="<your-mainnet-eth-rpc>"
#TB_NODE_SEPOLIARPC="<your-sepolia-eth-rpc>"

TB_LOGLEVEL="Info" # or Warn, Error, Debug
*/
package main
