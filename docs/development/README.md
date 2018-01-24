## Keep developer documentation

### Relay States

There is a set of threshold relay state diagrams auto-generated from this
repo's `docs` available at https://docs.keep.network/relay-states.pdf. The
images in the diagram, whose sources are at `img-src/*.tikz`, are also
available at `https://docs.keep.network/img/generated/*.png` (the filenames
are identical to their TikZ sources, with a `.png` suffix instead of
`.tikz`). These URLs are for the `master` version of the repo; non-`master`
branches are instead published to `https://docs.keep.network/<branch name>/`.

### [Getting started with `geth` on the test network](getting-started-ethereum.adoc)

A note-taken walkthrough of how to start from not having anything connected to
Ethereum to developing the basics of a smart contract that can emit events and
deploying it on the Rinkeby test network. Covers running `geth`, getting it
hooked into the Rinkeby testnet, getting some eth from the faucet, and
interacting with the JSON-RPC API. Also covers some basic solidity, compiling
it, and using JSON-RPC to install a contract and call it. Relatively low-level,
to provide some familiarity with how Ethereum works under the covers.

