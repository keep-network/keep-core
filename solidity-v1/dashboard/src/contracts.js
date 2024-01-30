import KeepToken from "@keep-network/keep-core/artifacts/KeepToken.json"
import TokenStaking from "@keep-network/keep-core/artifacts/TokenStaking.json"
import TokenGrant from "@keep-network/keep-core/artifacts/TokenGrant.json"
import KeepRandomBeaconOperator from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperator.json"
import BondedECDSAKeepFactory from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeepFactory.json"
import TBTCSystem from "@keep-network/tbtc/artifacts/TBTCSystem.json"
import KeepBonding from "@keep-network/keep-ecdsa/artifacts/KeepBonding.json"
import GuaranteedMinimumStakingPolicy from "@keep-network/keep-core/artifacts/GuaranteedMinimumStakingPolicy.json"
import PermissiveStakingPolicy from "@keep-network/keep-core/artifacts/PermissiveStakingPolicy.json"
import KeepRandomBeaconOperatorStatistics from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperatorStatistics.json"
import ManagedGrant from "@keep-network/keep-core/artifacts/ManagedGrant.json"
import ManagedGrantFactory from "@keep-network/keep-core/artifacts/ManagedGrantFactory.json"
import TBTCToken from "@keep-network/tbtc/artifacts/TBTCToken.json"
import Deposit from "@keep-network/tbtc/artifacts/Deposit.json"
import BondedECDSAKeep from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeep.json"
import TokenStakingEscrow from "@keep-network/keep-core/artifacts/TokenStakingEscrow.json"
import StakingPortBacker from "@keep-network/keep-core/artifacts/StakingPortBacker.json"
import BeaconRewards from "@keep-network/keep-core/artifacts/BeaconRewards.json"
import ECDSARewardsDistributor from "@keep-network/keep-ecdsa/artifacts/ECDSARewardsDistributor.json"
import LPRewardsKEEPETH from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPETH.json"
import LPRewardsTBTCETH from "@keep-network/keep-ecdsa/artifacts/LPRewardsTBTCETH.json"
import LPRewardsKEEPTBTC from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPTBTC.json"
import LPRewardsTBTCSaddle from "@keep-network/keep-ecdsa/artifacts/LPRewardsTBTCSaddle.json"
import LPRewardsTBTCv2Saddle from "@keep-network/keep-ecdsa/artifacts/LPRewardsTBTCv2Saddle.json"
import LPRewardsTBTCv2SaddleV2 from "@keep-network/keep-ecdsa/artifacts/LPRewardsTBTCv2SaddleV2.json"
import KeepOnlyPool from "@keep-network/keep-core/artifacts/KeepVault.json"
import IERC20 from "@keep-network/keep-core/artifacts/IERC20.json"
import SaddleSwap from "./contracts-artifacts/SaddleSwap.json"
import SaddleTBTCMetaPool from "./contracts-artifacts/SaddleTBTCMetaPool.json"
import SaddleTBTCMetaPoolV2 from "./contracts-artifacts/SaddleTBTCMetaPoolV2.json"

import Web3 from "web3"

import {
  KEEP_TOKEN_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  OPERATOR_CONTRACT_NAME,
  KEEP_OPERATOR_STATISTICS_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
  TBTC_TOKEN_CONTRACT_NAME,
  TBTC_SYSTEM_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
  BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
  STAKING_PORT_BACKER_CONTRACT_NAME,
  OLD_TOKEN_STAKING_CONTRACT_NAME,
  LP_REWARDS_KEEP_ETH_CONTRACT_NAME,
  LP_REWARDS_TBTC_ETH_CONTRACT_NAME,
  LP_REWARDS_KEEP_TBTC_CONTRACT_NAME,
  LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME,
  KEEP_TOKEN_GEYSER_CONTRACT_NAME,
  LP_REWARDS_TBTCV2_SADDLE_CONTRACT_NAME,
  LP_REWARDS_TBTCV2_SADDLEV2_CONTRACT_NAME,
  ASSET_POOL_CONTRACT_NAME,
} from "./constants/constants"

import KeepLib from "./lib/keep"
import { Web3jsWrapper } from "./lib/web3"
import { getChainId, getWsUrl } from "./connectors/utils.js"
import { ExplorerModeConnector } from "./connectors/explorer-mode-connector"

const CONTRACT_DEPLOYMENT_BLOCK_CACHE = {}

export const getContractDeploymentBlockNumber = async (
  contractName,
  web3Instance
) => {
  if (
    CONTRACT_DEPLOYMENT_BLOCK_CACHE?.[contractName] !== undefined &&
    CONTRACT_DEPLOYMENT_BLOCK_CACHE?.[contractName] !== null
  ) {
    return CONTRACT_DEPLOYMENT_BLOCK_CACHE[contractName]
  }
  const web3 = web3Instance || (await Web3Loaded)
  const blockNumber = await contractDeployedAtBlock(
    web3,
    contracts[contractName].artifact
  )

  CONTRACT_DEPLOYMENT_BLOCK_CACHE[contractName] = blockNumber

  return blockNumber
}

const contracts = {
  [KEEP_TOKEN_CONTRACT_NAME]: { artifact: KeepToken },
  [TOKEN_GRANT_CONTRACT_NAME]: { artifact: TokenGrant },
  [OPERATOR_CONTRACT_NAME]: {
    artifact: KeepRandomBeaconOperator,
    withDeployBlock: true,
  },
  [TOKEN_STAKING_CONTRACT_NAME]: {
    artifact: TokenStaking,
    withDeployBlock: true,
  },
  [KEEP_OPERATOR_STATISTICS_CONTRACT_NAME]: {
    artifact: KeepRandomBeaconOperatorStatistics,
  },
  [MANAGED_GRANT_FACTORY_CONTRACT_NAME]: {
    artifact: ManagedGrantFactory,
    withDeployBlock: true,
  },
  [KEEP_BONDING_CONTRACT_NAME]: {
    artifact: KeepBonding,
    withDeployBlock: true,
  },
  [TBTC_TOKEN_CONTRACT_NAME]: {
    artifact: TBTCToken,
    withDeployBlock: true,
  },
  [TBTC_SYSTEM_CONTRACT_NAME]: {
    artifact: TBTCSystem,
    withDeployBlock: true,
  },
  [TOKEN_STAKING_ESCROW_CONTRACT_NAME]: {
    artifact: TokenStakingEscrow,
    withDeployBlock: true,
  },
  [BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME]: {
    artifact: BondedECDSAKeepFactory,
  },
  [STAKING_PORT_BACKER_CONTRACT_NAME]: {
    artifact: StakingPortBacker,
    withDeployBlock: true,
  },
  beaconRewardsContract: {
    artifact: BeaconRewards,
  },
  ECDSARewardsDistributorContract: {
    artifact: ECDSARewardsDistributor,
    withDeployBlock: true,
  },
  [LP_REWARDS_KEEP_ETH_CONTRACT_NAME]: {
    artifact: LPRewardsKEEPETH,
    withDeployBlock: true,
  },
  [LP_REWARDS_TBTC_ETH_CONTRACT_NAME]: {
    artifact: LPRewardsTBTCETH,
    withDeployBlock: true,
  },
  [LP_REWARDS_KEEP_TBTC_CONTRACT_NAME]: {
    artifact: LPRewardsKEEPTBTC,
    withDeployBlock: true,
  },
  [LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME]: {
    artifact: LPRewardsTBTCSaddle,
    withDeployBlock: true,
  },
  [KEEP_TOKEN_GEYSER_CONTRACT_NAME]: {
    artifact: KeepOnlyPool,
    withDeployBlock: true,
  },
  [LP_REWARDS_TBTCV2_SADDLE_CONTRACT_NAME]: {
    artifact: LPRewardsTBTCv2Saddle,
  },
  [LP_REWARDS_TBTCV2_SADDLEV2_CONTRACT_NAME]: {
    artifact: LPRewardsTBTCv2SaddleV2,
  },
}

export async function getKeepTokenContractDeployerAddress(web3) {
  const deployTransactionHash = getTransactionHashOfContractDeploy(KeepToken)
  const transaction = await web3.eth.getTransaction(deployTransactionHash)

  return transaction.from
}

async function contractDeployedAtBlock(web3, contract) {
  const deployTransactionHash = getTransactionHashOfContractDeploy(contract)
  const transaction = await web3.eth.getTransaction(deployTransactionHash)

  return transaction.blockNumber.toString()
}

export function Deferred() {
  let resolve
  let reject

  const promise = new Promise((res, rej) => {
    resolve = res
    reject = rej
  })

  let isPending = true
  let isRejected = false
  let isFulfilled = false

  const extraPromise = promise
    .then((fn) => {
      isFulfilled = true
      isPending = false
      return fn
    })
    .catch((e) => {
      isRejected = true
      isPending = false
      throw e
    })

  extraPromise.isFulfilled = () => isFulfilled
  extraPromise.isRejected = () => isRejected
  extraPromise.isPending = () => isPending

  return {
    promise: extraPromise,
    reject,
    resolve,
  }
}

const ContractsDeferred = new Deferred()
const Web3Deferred = new Deferred()

/** @type {KeepLib} */
export const Keep = KeepLib.initialize(
  new Web3jsWrapper(new Web3(getWsUrl())),
  getChainId()
)

export const KeepExplorerMode = KeepLib.initialize(
  new Web3jsWrapper(new Web3(getWsUrl())),
  getChainId()
)

export const Web3Loaded = Web3Deferred.promise
export const ContractsLoaded = ContractsDeferred.promise

export const resolveWeb3Deferred = async (web3) => {
  if (Web3Loaded.isFulfilled()) {
    const existingWeb3 = await Web3Loaded
    existingWeb3.setProvider(web3.currentProvider)
    existingWeb3.eth.defaultAccount = web3.eth.defaultAccount
  } else {
    Web3Deferred.resolve(web3)
  }
}

export const resovleContractsDeferred = (contracts) => {
  ContractsDeferred.resolve(contracts)
}

export async function getContracts(web3, netId) {
  if (netId.toString() !== getFirstNetworkIdFromArtifact()) {
    console.error(
      `network id: ${netId}; expected network id ${getFirstNetworkIdFromArtifact()}`
    )
    throw new Error("Please connect to the appropriate Ethereum network.")
  }

  // if (ContractsLoaded.isFulfilled()) {
  //   console.log("if fulfilled")
  //   const existingContracts = await ContractsLoaded
  //   for (const contractInstance of Object.values(existingContracts)) {
  //     contractInstance.options.from = web3.eth.defaultAccount
  //   }

  //   return contracts
  // }
  // console.log("else not fulfilled")

  // const web3Contracts = {}
  // for (const [contractName, options] of Object.entries(contracts)) {
  //   options.contractName = contractName
  //   web3Contracts[contractName] = await getContract(
  //     web3,
  //     options.artifact,
  //     options
  //   )
  // }

  // const oldTokenStakingArtifact = await getOldTokenStakingArtifact()
  // web3Contracts[OLD_TOKEN_STAKING_CONTRACT_NAME] = await getContract(
  //   web3,
  //   oldTokenStakingArtifact,
  //   { contractName: OLD_TOKEN_STAKING_CONTRACT_NAME }
  // )

  // resovleContractsDeferred(web3Contracts)

  const web3Contracts = {}
  for (const [contractName, options] of Object.entries(contracts)) {
    options.contractName = contractName
    web3Contracts[contractName] = await getContract(
      web3,
      options.artifact,
      options
    )
  }

  const oldTokenStakingArtifact = await getOldTokenStakingArtifact()
  web3Contracts[OLD_TOKEN_STAKING_CONTRACT_NAME] = await getContract(
    web3,
    oldTokenStakingArtifact,
    { contractName: OLD_TOKEN_STAKING_CONTRACT_NAME }
  )

  console.log("heheszek: ", web3Contracts[OLD_TOKEN_STAKING_CONTRACT_NAME])

  resovleContractsDeferred(web3Contracts)
  return web3Contracts
}

const getContract = async (web3, jsonArtifact) => {
  const address = getContractAddress(jsonArtifact)

  return createWeb3ContractInstance(web3, jsonArtifact.abi, address)
}

const createWeb3ContractInstance = (web3, abi, address) => {
  const contract = new web3.eth.Contract(abi, address)
  contract.options.from = web3.eth.defaultAccount
  contract.options.handleRevert = true

  return contract
}

export function isCodeValid(code) {
  return code && code !== "0x0" && code !== "0x"
}

function getTransactionHashOfContractDeploy({ networks }) {
  return networks[Object.keys(networks)[0]].transactionHash
}

export function getContractAddress({ networks }) {
  return networks[Object.keys(networks)[0]].address
}

// The artifacts from @keep-network/keep-core for a given build only support a single network id
export function getFirstNetworkIdFromArtifact() {
  return Object.keys(KeepToken.networks)[0]
}

export function getPermissiveStakingPolicyContractAddress() {
  return getContractAddress(PermissiveStakingPolicy)
}

export function getGuaranteedMinimumStakingPolicyContractAddress() {
  return getContractAddress(GuaranteedMinimumStakingPolicy)
}

export function createManagedGrantContractInstance(web3, address) {
  return createWeb3ContractInstance(web3, ManagedGrant.abi, address)
}

export function createDepositContractInstance(web3, address) {
  return createWeb3ContractInstance(web3, Deposit.abi, address)
}

export function createBondedECDSAKeepContractInstance(web3, address) {
  return createWeb3ContractInstance(web3, BondedECDSAKeep.abi, address)
}

export function getKeepRandomBeaconOperatorAddress() {
  return getContractAddress(KeepRandomBeaconOperator)
}

export function getBondedECDSAKeepFactoryAddress() {
  return getContractAddress(BondedECDSAKeepFactory)
}

export function getTBTCSystemAddress() {
  return getContractAddress(TBTCSystem)
}

const getOldTokenStakingArtifact = async () => {
  if (getFirstNetworkIdFromArtifact() === "1") {
    // Mainnet network ID.
    // Against mainnet, we want to use TokenStaking artifact
    // from 1.1.2 version at `0x6D1140a8c8e6Fac242652F0a5A8171b898c67600` address.
    return (await import("./old-contracts-artifacts/TokenStaking.json")).default
  }

  // For local, Ropsten and keep-dev network we want to use
  // the mocked old `TokenStaking` contract from `@keep-network/keep-core` package.
  return (
    await import("@keep-network/keep-core/artifacts/OldTokenStaking.json")
  ).default
}

export const createERC20Contract = (web3, address) => {
  return createWeb3ContractInstance(web3, IERC20.abi, address)
}

export const initializeWeb3 = (provider) => {
  return new Web3(provider)
}

export const createLPRewardsContract = async (web3, contractName) => {
  const { artifact } = contracts[contractName]
  return await getContract(web3, artifact, {})
}

export function createSaddleSwapContract(web3) {
  return createWeb3ContractInstance(
    web3,
    SaddleSwap.abi,
    "0x4f6A43Ad7cba042606dECaCA730d4CE0A57ac62e"
  )
}

export function createSaddleTBTCMetaPool(web3) {
  return createWeb3ContractInstance(
    web3,
    SaddleTBTCMetaPool.abi,
    // https://etherscan.io/address/0xf74ebe6e5586275dc4CeD78F5DBEF31B1EfbE7a5#code
    "0xf74ebe6e5586275dc4CeD78F5DBEF31B1EfbE7a5"
  )
}

export function createSaddleTBTCMetaPoolV2(web3) {
  return createWeb3ContractInstance(
    web3,
    SaddleTBTCMetaPoolV2.abi,
    // https://etherscan.io/address/0xA0b4a2667dD60d5CdD7EcFF1084F0CeB8dD84326#code
    "0xA0b4a2667dD60d5CdD7EcFF1084F0CeB8dD84326"
  )
}
