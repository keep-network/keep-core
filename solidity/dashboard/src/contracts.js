import KeepToken from "@keep-network/keep-core/artifacts/KeepToken.json"
import TokenStaking from "@keep-network/keep-core/artifacts/TokenStaking.json"
import TokenGrant from "@keep-network/keep-core/artifacts/TokenGrant.json"
import KeepRandomBeaconOperator from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperator.json"
import BondedECDSAKeepFactory from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeepFactory.json"
import TBTCSystem from "@keep-network/tbtc/artifacts/TBTCSystem.json"
import KeepRegistry from "@keep-network/keep-core/artifacts/KeepRegistry.json"
import GuaranteedMinimumStakingPolicy from "@keep-network/keep-core/artifacts/GuaranteedMinimumStakingPolicy.json"
import PermissiveStakingPolicy from "@keep-network/keep-core/artifacts/PermissiveStakingPolicy.json"
import KeepRandomBeaconOperatorStatistics from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperatorStatistics.json"
import ManagedGrant from "@keep-network/keep-core/artifacts/ManagedGrant.json"
import ManagedGrantFactory from "@keep-network/keep-core/artifacts/ManagedGrantFactory.json"
import {
  KEEP_TOKEN_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  OPERATOR_CONTRACT_NAME,
  REGISTRY_CONTRACT_NAME,
  KEEP_OPERATOR_STATISTICS_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
} from "./constants/constants"

export const CONTRACT_DEPLOY_BLOCK_NUMBER = {
  [KEEP_TOKEN_CONTRACT_NAME]: 0,
  [TOKEN_GRANT_CONTRACT_NAME]: 0,
  [OPERATOR_CONTRACT_NAME]: 0,
  [TOKEN_STAKING_CONTRACT_NAME]: 0,
  [REGISTRY_CONTRACT_NAME]: 0,
  [KEEP_OPERATOR_STATISTICS_CONTRACT_NAME]: 0,
  [MANAGED_GRANT_FACTORY_CONTRACT_NAME]: 0,
}

export async function getKeepToken(web3) {
  return getContract(web3, KeepToken, KEEP_TOKEN_CONTRACT_NAME)
}

export async function getTokenStaking(web3) {
  return getContract(web3, TokenStaking, TOKEN_STAKING_CONTRACT_NAME)
}

export async function getTokenGrant(web3) {
  return getContract(web3, TokenGrant, TOKEN_GRANT_CONTRACT_NAME)
}

export async function getKeepRandomBeaconOperator(web3) {
  return getContract(web3, KeepRandomBeaconOperator, OPERATOR_CONTRACT_NAME)
}

export async function getRegistry(web3) {
  return getContract(web3, KeepRegistry, REGISTRY_CONTRACT_NAME)
}

export async function getKeepRandomBeaconOperatorStatistics(web3) {
  return getContract(
    web3,
    KeepRandomBeaconOperatorStatistics,
    KEEP_OPERATOR_STATISTICS_CONTRACT_NAME
  )
}

export async function getManagedGrantFactory(web3) {
  return getContract(web3, ManagedGrantFactory)
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

export async function getContracts(web3) {
  const contracts = await Promise.all([
    getKeepToken(web3),
    getTokenGrant(web3),
    getTokenStaking(web3),
    getKeepRandomBeaconOperator(web3),
    getRegistry(web3),
    getKeepRandomBeaconOperatorStatistics(web3),
    getManagedGrantFactory(web3),
  ])

  return {
    token: contracts[0],
    grantContract: contracts[1],
    stakingContract: contracts[2],
    keepRandomBeaconOperatorContract: contracts[3],
    registryContract: contracts[4],
    keepRandomBeaconOperatorStatistics: contracts[5],
    managedGrantFactoryContract: contracts[6],
  }
}

async function getContract(web3, contract, contractName) {
  const address = getContractAddress(contract)
  const code = await web3.eth.getCode(address)

  if (!isCodeValid(code)) throw Error("No contract at address")
  CONTRACT_DEPLOY_BLOCK_NUMBER[contractName] = await contractDeployedAtBlock(
    web3,
    contract
  )
  return new web3.eth.Contract(contract.abi, address)
}

export function isCodeValid(code) {
  return code && code !== "0x0" && code !== "0x"
}

function getTransactionHashOfContractDeploy({ networks }) {
  return networks[Object.keys(networks)[0]].transactionHash
}

function getContractAddress({ networks }) {
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
  return new web3.eth.Contract(ManagedGrant.abi, address)
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