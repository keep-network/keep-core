# Environment provides the solidity directory as a potentially-relative path,
# which we resolve. Then we resolve the Solidity files in a contracts/ directory
# at that path.
solidity_dir := $(realpath ${SOLIDITY_DIR})

# Contracts for which the bindings should be generated.
required_contracts := RandomBeacon SortitionPool

solidity_dependencies_mapping := \
	@openzeppelin/=${solidity_dir}/node_modules/@openzeppelin/ \
	@keep-network/=${solidity_dir}/node_modules/@keep-network/  \
	@threshold-network/=${solidity_dir}/node_modules/@threshold-network/  \
	@thesis/=${solidity_dir}/node_modules/@thesis/

include ../../common/gen/Makefile
