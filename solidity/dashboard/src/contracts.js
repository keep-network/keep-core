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
]

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

function Deferred() {
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

export const ContractsDeferred = new Deferred()
export const Web3Deferred = new Deferred()

export const Web3Loadeed = Web3Deferred.promise
export const ContractsLoaded = ContractsDeferred.promise

export async function getContracts(web3) {
  const web3Contracts = {}
  for (const contractData of contracts) {
    const [options, jsonArtifact] = contractData

    web3Contracts[options.contractName] = await getContract(
      web3,
      jsonArtifact,
      options
    )
  }

  ContractsDeferred.resolve(web3Contracts)
  return web3Contracts
}

const getContract = async (web3, jsonArtifact, options) => {
  const { contractName, withDeployBlock } = options
  const address = getContractAddress(jsonArtifact)
  // const code = await web3.eth.getCode(address)

  // if (!isCodeValid(code)) throw Error("No contract at address")
  if (withDeployBlock) {
    CONTRACT_DEPLOY_BLOCK_NUMBER[contractName] = await contractDeployedAtBlock(
      web3,
      jsonArtifact
    )
  }
  return new web3.eth.Contract(jsonArtifact.abi, address)
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
