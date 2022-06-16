# Environment provides the solidity directory as a potentially-relative path,
# which we resolve. Then we resolve the Solidity files in a contracts/ directory
# at that path.
solidity_dir := $(realpath ${SOLIDITY_DIR})

# Define a path to the Solidity files, as it's not the default `solidity_dir/contracts`
# direcotry.
contracts_dir := $(realpath ${solidity_dir}/node_modules/@threshold-network/solidity-contracts/contracts)

# Contracts for which the bindings should be generated.
required_contracts := TokenStaking

solidity_dependencies_mapping := \
	@openzeppelin/=${solidity_dir}/node_modules/@openzeppelin/ \
	@thesis/=${solidity_dir}/node_modules/@thesis/

include ../../common/gen/Makefile
