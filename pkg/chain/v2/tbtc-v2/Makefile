
# Environment provides the solidity directory as a potentially-relative path,
# which we resolve. Then we resolve the Solidity files in a contracts/ directory
# at that path.
solidity_dir := $(realpath ${SOLIDITY_DIR})

# Contracts solidity sources relative to the solidity_dir.
solidity_files_sources := \
	node_modules/@keep-network/tbtc-v2/contracts/bridge/Bridge.sol

solidity_dependencies_mapping := \
	@openzeppelin/=${solidity_dir}/node_modules/@openzeppelin/ \
	@keep-network/=${solidity_dir}/node_modules/@keep-network/ \
	@thesis/=${solidity_dir}/node_modules/@thesis/

include ../common/Makefile
