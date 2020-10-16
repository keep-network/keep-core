import KeepToken from "@keep-network/keep-core/artifacts/KeepToken.json"
import TokenStaking from "@keep-network/keep-core/artifacts/TokenStaking.json"
import TokenGrant from "@keep-network/keep-core/artifacts/TokenGrant.json"
import KeepRandomBeaconOperator from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperator.json"
import BondedECDSAKeepFactory from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeepFactory.json"
import TBTCSystem from "@keep-network/tbtc/artifacts/TBTCSystem.json"
import TBTCConstants from "@keep-network/tbtc/artifacts/TBTCConstants.json"
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

import {
  KEEP_TOKEN_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  OPERATOR_CONTRACT_NAME,
  REGISTRY_CONTRACT_NAME,
  KEEP_OPERATOR_STATISTICS_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
  TBTC_TOKEN_CONTRACT_NAME,
  TBTC_SYSTEM_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
  BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
  STAKING_PORT_BACKER_CONTRACT_NAME,
  OLD_TOKEN_STAKING_CONTRACT_NAME,
  TBTC_CONSTANTS_CONTRACT_NAME
} from "./constants/constants"

export const CONTRACT_DEPLOY_BLOCK_NUMBER = {
  [KEEP_TOKEN_CONTRACT_NAME]: 0,
  [TOKEN_GRANT_CONTRACT_NAME]: 0,
  [OPERATOR_CONTRACT_NAME]: 0,
  [TOKEN_STAKING_CONTRACT_NAME]: 0,
  [REGISTRY_CONTRACT_NAME]: 0,
  [KEEP_OPERATOR_STATISTICS_CONTRACT_NAME]: 0,
  [MANAGED_GRANT_FACTORY_CONTRACT_NAME]: 0,
  [KEEP_BONDING_CONTRACT_NAME]: 0,
  [TBTC_TOKEN_CONTRACT_NAME]: 0,
  [TBTC_SYSTEM_CONTRACT_NAME]: 0,
  [TOKEN_STAKING_ESCROW_CONTRACT_NAME]: 0,
  [STAKING_PORT_BACKER_CONTRACT_NAME]: 0,
  [TBTC_CONSTANTS_CONTRACT_NAME]: 0,
}

const contracts = [
  [{ contractName: KEEP_TOKEN_CONTRACT_NAME }, KeepToken],
  [{ contractName: TOKEN_GRANT_CONTRACT_NAME }, TokenGrant],
  [
    { contractName: OPERATOR_CONTRACT_NAME, withDeployBlock: true },
    KeepRandomBeaconOperator,
  ],
  [
    { contractName: TOKEN_STAKING_CONTRACT_NAME, withDeployBlock: true },
    TokenStaking,
  ],
  [
    {
      contractName: KEEP_OPERATOR_STATISTICS_CONTRACT_NAME,
    },
    KeepRandomBeaconOperatorStatistics,
  ],
  [
    {
      contractName: MANAGED_GRANT_FACTORY_CONTRACT_NAME,
      withDeployBlock: true,
    },
    ManagedGrantFactory,
  ],
  [
    { contractName: KEEP_BONDING_CONTRACT_NAME, withDeployBlock: true },
    KeepBonding,
  ],
  [
    { contractName: TBTC_TOKEN_CONTRACT_NAME, withDeployBlock: true },
    TBTCToken,
  ],
  [
    { contractName: TBTC_SYSTEM_CONTRACT_NAME, withDeployBlock: true },
    TBTCSystem,
  ],
  [
    { contractName: TOKEN_STAKING_ESCROW_CONTRACT_NAME, withDeployBlock: true },
    TokenStakingEscrow,
  ],
  [
    { contractName: BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME },
    BondedECDSAKeepFactory,
  ],
  [
    { contractName: STAKING_PORT_BACKER_CONTRACT_NAME, withDeployBlock: true },
    StakingPortBacker,
  ],
]

export async function getTBTCConstantsContract(web3) {
  return getContract(web3, TBTCConstants, TBTC_CONSTANTS_CONTRACT_NAME)
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

  return {
    promise,
    reject,
    resolve,
  }
}

let ContractsDeferred = new Deferred()
let Web3Deferred = new Deferred()

export let Web3Loaded = Web3Deferred.promise
export let ContractsLoaded = ContractsDeferred.promise

export const resolveWeb3Deferred = (web3) => {
  Web3Deferred = new Deferred()
  Web3Deferred.resolve(web3)
  Web3Loaded = Web3Deferred.promise
}

export const resovleContractsDeferred = (contracts) => {
  ContractsDeferred = new Deferred()
  ContractsDeferred.resolve(contracts)
  ContractsLoaded = ContractsDeferred.promise
}

export async function getContracts(web3) {
  // This is a workaround for `web3.eth.net.getId()`, since on local machine
  // returns unexpected values(eg. 105) and after calling this method the web3
  // cannot call any other function eq `getTransaction` and doesnt throw errors.
  // This is probably problem with `WebocketProvider`, because if we replace it
  // with `RpcSubprovider` the result will be as expected (eg. for Mainnet
  // returns 1).
  const netIdDeferred = new Deferred()
  web3.currentProvider.sendAsync(
    {
      jsonrpc: "2.0",
      method: "net_version",
      params: [],
      id: new Date().getTime(),
    },
    (error, response) => {
      if (error) {
        netIdDeferred.reject(error)
      } else {
        netIdDeferred.resolve(response.result)
      }
    }
  )
  const netId = await netIdDeferred.promise
  if (netId.toString() !== getFirstNetworkIdFromArtifact()) {
    console.error(
      `network id: ${netId}; expected network id ${getFirstNetworkIdFromArtifact()}`
    )
    throw new Error("Please connect to the appropriate Ethereum network.")
  }

  const web3Contracts = {}
  for (const contractData of contracts) {
    const [options, jsonArtifact] = contractData

    web3Contracts[options.contractName] = await getContract(
      web3,
      jsonArtifact,
      options
    )
  }

  const oldTokenStakingArtifact = await getOldTokenStakingArtifact()
  web3Contracts[OLD_TOKEN_STAKING_CONTRACT_NAME] = await getContract(
    web3,
    oldTokenStakingArtifact,
    { contractName: OLD_TOKEN_STAKING_CONTRACT_NAME }
  )

  resovleContractsDeferred(web3Contracts)
  return web3Contracts
}

const getContract = async (web3, jsonArtifact, options) => {
  const { contractName, withDeployBlock } = options
  const address = getContractAddress(jsonArtifact)

  if (withDeployBlock) {
    CONTRACT_DEPLOY_BLOCK_NUMBER[contractName] = await contractDeployedAtBlock(
      web3,
      jsonArtifact
    )
  }
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
