npm_package_name=@keep-network/tbtc-v2

# Contracts for which the bindings should be generated.
required_contracts := Bridge LightRelay WalletCoordinator

# There is a bug in the currently used abigen version (v1.10.19) that makes it
# re-declaring structs used by multiple contracts
# https://github.com/ethereum/go-ethereum/issues/24627. This is a problem
# for us because both Bridge and WalletCoordinator contracts use BitcoinTx.Info
# struct which is then re-declared in the same package once abigen does its work.
# An ultimate solution would be upgrading go-ethereum (thus abigen too) to v1.11 but
# that version contains some breaking changes that make the upgrade non-trivial.
# As a short-term workaround, we use some Makefile shenanigans to slightly rename
# the conflicting BitcoinTxInfo struct in the WalletCoordinator output file.
# We use perl for that purpose as sed is not cross-platform and works a bit
# differently on GNU and BSD.
#
# TODO: Remove once go-ethereum is upgraded to v1.11.
define after_abi_hook
	$(eval type := $(1))
	$(if $(filter $(type),WalletCoordinator),$(call fix_wallet_coordinator_collision))
endef
define fix_wallet_coordinator_collision
	@perl -pi -e s,BitcoinTxInfo,BitcoinTxInfo2,g ./abi/WalletCoordinator.go
endef

include ../../common/gen/Makefile
