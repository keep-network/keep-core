# Environment provides the solidity directory as a potentially-relative path,
# which we resolve. Then we resolve the Solidity files in a contracts/ directory
# at that path.
solidity_dir := $(realpath ${SOLIDITY_DIR})

# Contracts solidity sources relative to the solidity_dir.
solidity_files_sources := \
	contracts/RandomBeacon.sol \
	node_modules/@keep-network/sortition-pools/contracts/SortitionPool.sol

# Solidity files prefixed with the solidity directory path.
solidity_files := $(addprefix ${solidity_dir}/,${solidity_files_sources}) 

# Bare Solidity filenames without .sol or Solidity directory prefix.
contract_stems := $(notdir $(basename ${solidity_files}))
contract_files := $(addprefix contract/,$(addsuffix .go,${contract_stems}))

solidity_dependencies_mapping := \
	@openzeppelin/=${solidity_dir}/node_modules/@openzeppelin/ \
	@keep-network/=${solidity_dir}/node_modules/@keep-network/  \
	@threshold-network/=${solidity_dir}/node_modules/@threshold-network/  \
	@thesis/=${solidity_dir}/node_modules/@thesis/

include ../common/Makefile
