# Environment provides the solidity directory as a potentially-relative path,
# which we resolve. Then we resolve the Solidity files in a contracts/ directory
# at that path.
solidity_dir := $(realpath ${SOLIDITY_DIR})

# Contracts solidity sources relative to the solidity_dir.
solidity_files_sources := \
	contracts/WalletRegistry.sol \
	node_modules/@keep-network/sortition-pools/contracts/SortitionPool.sol

solidity_dependencies_mapping := \
	openzeppelin-solidity/=${solidity_dir}/node_modules/openzeppelin-solidity/ \
	@openzeppelin/=${solidity_dir}/node_modules/@openzeppelin/ \
	@keep-network/=${solidity_dir}/node_modules/@keep-network/ \
	@threshold-network/=${solidity_dir}/node_modules/@threshold-network/ \
	@thesis/=${solidity_dir}/node_modules/@thesis/

include ../common/Makefile
