export const KEEP_TOKEN_CONTRACT_NAME = "token"
export const TOKEN_STAKING_CONTRACT_NAME = "stakingContract"
export const TOKEN_GRANT_CONTRACT_NAME = "grantContract"
export const OPERATOR_CONTRACT_NAME = "keepRandomBeaconOperatorContract"
export const REGISTRY_CONTRACT_NAME = "registryContract"
export const KEEP_OPERATOR_STATISTICS_CONTRACT_NAME =
  "keepRandomBeaconOperatorStatistics"
export const MANAGED_GRANT_FACTORY_CONTRACT_NAME = "managedGrantFactoryContract"
export const BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME =
  "bondedEcdsaKeepFactoryContract"
export const KEEP_BONDING_CONTRACT_NAME = "keepBondingContract"
export const TBTC_TOKEN_CONTRACT_NAME = "tbtcTokenContract"
export const TBTC_SYSTEM_CONTRACT_NAME = "tbtcSystemContract"
export const TOKEN_STAKING_ESCROW_CONTRACT_NAME = "tokenStakingEscrow"
export const OLD_TOKEN_STAKING_CONTRACT_NAME = "oldTokenStakingContract"
export const STAKING_PORT_BACKER_CONTRACT_NAME = "stakingPortBackerContract"
export const TBTC_CONSTANTS_CONTRACT_NAME = "tbtcConstantsContract"

export const PENDING_STATUS = "PENDING"
export const COMPLETE_STATUS = "COMPLETE"

export const WALLETS = {
  METAMASK: { label: "MetaMask" },
  TREZOR: { label: "Trezor" },
  LEDGER: { label: "Ledger" },
  COINBASE: { label: "Coinbase" },
}

export const REWARD_STATUS = {
  AVAILABLE: "AVAILABLE",
  WITHDRAWN: "WITHDRAWN",
  ACCUMULATING: "ACCUMULATING",
}

export const SIGNING_GROUP_STATUS = {
  COMPLETED: "Completed",
  TERMINATED: "Terminated",
  ACTIVE: "Active work",
}

/**
 * Copied from `@keep-network/tbtc/artifacts/DepositStates.json`, since metadata is located in ast as opposed to the abi.
 * Can't read with web3.js
 */
export const DEPOSIT_STATES = [
  {
    id: 3300,
    name: "START",
    depositStatusId: 0,
  },
  {
    id: 3301,
    name: "AWAITING_SIGNER_SETUP",
    depositStatusId: 1,
  },
  {
    id: 3302,
    name: "AWAITING_BTC_FUNDING_PROOF",
    depositStatusId: 2,
  },
  {
    id: 3303,
    name: "FAILED_SETUP",
    depositStatusId: 3,
  },
  {
    id: 3304,
    name: "ACTIVE",
    depositStatusId: 4,
  },
  {
    id: 3305,
    name: "AWAITING_WITHDRAWAL_SIGNATURE",
    depositStatusId: 5,
  },
  {
    id: 3306,
    name: "AWAITING_WITHDRAWAL_PROOF",
    depositStatusId: 6,
  },
  {
    id: 3307,
    name: "REDEEMED",
    depositStatusId: 7,
  },
  {
    id: 3308,
    name: "COURTESY_CALL",
    depositStatusId: 8,
  },
  {
    id: 3309,
    name: "FRAUD_LIQUIDATION_IN_PROGRESS",
    depositStatusId: 9,
  },
  {
    id: 3310,
    name: "LIQUIDATION_IN_PROGRESS",
    depositStatusId: 10,
  },
  {
    id: 3311,
    name: "LIQUIDATED",
    depositStatusId: 11,
  },
]
