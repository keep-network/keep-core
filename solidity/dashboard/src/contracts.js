import KeepToken from "@keep-network/keep-core/artifacts/KeepToken.json"
import TokenStaking from "@keep-network/keep-core/artifacts/TokenStaking.json"
import TokenGrant from "@keep-network/keep-core/artifacts/TokenGrant.json"
import KeepRandomBeaconOperator from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperator.json"
import BondedECDSAKeepFactory from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeepFactory.json"
import TBTCSystem from "@keep-network/tbtc/artifacts/TBTCSystem.json"
import KeepBonding from "@keep-network/keep-ecdsa/artifacts/KeepBonding.json"
import KeepRegistry from "@keep-network/keep-core/artifacts/KeepRegistry.json"
import GuaranteedMinimumStakingPolicy from "@keep-network/keep-core/artifacts/GuaranteedMinimumStakingPolicy.json"
import PermissiveStakingPolicy from "@keep-network/keep-core/artifacts/PermissiveStakingPolicy.json"
import KeepRandomBeaconOperatorStatistics from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperatorStatistics.json"
import ManagedGrant from "@keep-network/keep-core/artifacts/ManagedGrant.json"
import ManagedGrantFactory from "@keep-network/keep-core/artifacts/ManagedGrantFactory.json"
import TBTCToken from "@keep-network/tbtc/artifacts/TBTCToken.json"
import Deposit from "@keep-network/tbtc/artifacts/Deposit.json"
import BondedECDSAKeep from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeep.json"
import {
  KEEP_TOKEN_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  OPERATOR_CONTRACT_NAME,
  REGISTRY_CONTRACT_NAME,
  KEEP_OPERATOR_STATISTICS_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
  BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
  TBTC_TOKEN_CONTRACT_NAME,
  TBTC_SYSTEM_CONTRACT_NAME,
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

export async function getBondedEcdsaKeepFactoryContract(web3) {
  return getContract(
    web3,
    BondedECDSAKeepFactory,
    BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME
  )
}

export async function getKeepBondingContract(web3) {
  return getContract(web3, KeepBonding, KEEP_BONDING_CONTRACT_NAME)
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

export async function getTBTCTokenContract(web3) {
  return getContract(web3, TBTCToken, TBTC_TOKEN_CONTRACT_NAME)
}

export async function getTBTCSystemContract(web3) {
  return getContract(web3, TBTCSystem, TBTC_SYSTEM_CONTRACT_NAME)
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

export let contracts

export async function getContracts(web3) {
  const web3Contracts = await Promise.all([
    getKeepToken(web3),
    getTokenGrant(web3),
    getTokenStaking(web3),
    getKeepRandomBeaconOperator(web3),
    getRegistry(web3),
    getKeepRandomBeaconOperatorStatistics(web3),
    getManagedGrantFactory(web3),
    getBondedEcdsaKeepFactoryContract(web3),
    getKeepBondingContract(web3),
    getTBTCTokenContract(web3),
    getTBTCSystemContract(web3),
  ])

  web3.eth.Contract.handleRevert = true
  web3.eth.Contract.defaultAddress = web3.eth.defaultAccount

  contracts = {
    token: web3Contracts[0],
    grantContract: web3Contracts[1],
    stakingContract: web3Contracts[2],
    keepRandomBeaconOperatorContract: web3Contracts[3],
    registryContract: web3Contracts[4],
    keepRandomBeaconOperatorStatistics: web3Contracts[5],
    managedGrantFactoryContract: web3Contracts[6],
    bondedEcdsaKeepFactoryContract: web3Contracts[7],
    keepBondingContract: web3Contracts[8],
    tbtcTokenContract: web3Contracts[9],
    tbtcSystemContract: web3Contracts[10],
  }

  return contracts
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
  return new web3.eth.Contract(ManagedGrant.abi, address)
}

export function createDepositContractInstance(web3, address) {
  return new web3.eth.Contract(Deposit.abi, address)
}

export function createBondedECDSAKeepContractInstance(web3, address) {
  return new web3.eth.Contract(BondedECDSAKeep.abi, address)
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
