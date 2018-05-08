## Keep developer documentation

### Getting Set Up

If you're on macOS, install Homebrew and run `scripts/macos-setup.sh`. Note
that if you don't have Homebrew or you're not on macOS, the below information
details what you'll need. The script additionally sets up pre-commit hooks.

### Building

Currently the easiest way to build is using the `Dockerfile` at the root of the
repository. A simple `docker build` should get you a functioning container.

If you want to build natively, there are a few prereqs you'll need to go through.
First, you'll need to clone `keep-network/bn` from GitHub and run `make
install`. To successfully build `bn`, you'll need to have `libgmp` (with
headers) and `openssl` (also with headers) installed, as well as the LLVM
developer tools. On macOS, you can `brew install gmp openssl llvm` to install
all of these. Note that `llvm` requires some additional env variable work that
the formula will explain when you install it.

You'll also need [`dep`](https://github.com/golang/dep#installation), the Go
dependency manager we use.

Lastly, you'll need the [protobuf compiler](https://developers.google.com/protocol-buffers/docs/downloads).
You'll also need to install the `protoc-gen-gogoslick` toolchain, which you can
install using `go get`:

```
go get -u github.com/gogo/protobuf/protoc-gen-gogoslick
```

Finally, you can run `dep ensure` in the root directory of this repository and
you'll be ready to build!

### Code Style

Go code generally follows common community style. We try to track with core Go
practices for the most part, including formatting using `go-imports` and
linting using `go-vet` and `go-lint` and keeping an eye on the collection of
[Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
Two major deviations worth calling out:

 - We do *not* prefix commit messages with the packages touched by the commit.
   The commit includes diffs, diffs include paths, paths imply packages. We
   consider this unnecessary and noisy.
 - We *discourage* single-letter variable names and related extra-shortness,
   with exceptions for external packages (we use the package name irrespective
   of our own practices, for the most part), the `err` variable, and iteration
   indices. Short variable names produce diffs that are more difficult to
   analyze quickly, and generally result in lower clarity for less experienced
   developers. We consider this an antipattern, and the additional typed
   characters to be comparatively very cheap.

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

