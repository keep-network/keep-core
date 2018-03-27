## Keep developer documentation

### Building

Currently the easiest way to build is using the `Dockerfile` at the root of the
repository. A simple `docker build` should get you a functioning container.

If you want to build natively, you'll need to clone `keep-network/bn` from
GitHub and run `make install`. To successfully build `bn`, you'll need to have
libgmp (with headers) and openssl (also with headers) installed, as well as
the LLVM developer tools. On macOS, you can `brew install gmp openssl llvm` to
install all of these. Note that `llvm` requires some additional env variable
work that the formula will explain when you install it.

Once you've installed `bn`, you can run `dep ensure` in the `go/` directory of
this repository and then you are ready to build.

#### Protobufs

In addition to installing `bn` and `dep`, you'll also need to install the protobuf compiler.
On OSX, this will be `brew install protobuf` (requirement: need `homebrew installed first`).

Lastly, you'll need to get the protoc-gen-gogo toolchain:

    ```
    go get github.com/gogo/protobuf/proto
    go get github.com/gogo/protobuf/jsonpb
    go get github.com/gogo/protobuf/protoc-gen-gogo
    go get github.com/gogo/protobuf/gogoproto
    ```

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

