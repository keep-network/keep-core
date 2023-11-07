import contracts, {
  RISK_MANAGER_V1_CONTRACT_NAME,
  SADDLE_SWAP_CONTRACT_NAME,
  SaddleSwapArtifact,
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
  COV_TOKEN_CONTRACT_NAME,
  ASSET_POOL_CONTRACT_NAME,
  REWARDS_POOL_CONTRACT_NAME,
  THRESHOLD_STAKING_CONTRACT_NAME,
  THRESHOLD_KEEP_STAKE_CONTRACT_NAME,
  SIMPLE_PRE_APPLICATION_CONTRACT_NAME,
} from "./contracts"
import CoveragePoolV1 from "./coverage-pool"
import { UniswapV2Exchange } from "./exchange-api"
import KeepToTStaking from "./keep-to-t-staking"
import { PRE } from "./constants"
import SimplePreApplicationAbi from "./contracts-artifacts/SimplePreApplication.json"

/** @typedef { import("../web3").Web3LibWrapper} Web3LibWrapper */
/** @typedef { import("../web3").BaseContract} BaseContract */
/** @typedef { import("./exchange-api").BaseExchange} BaseExchange */

class Keep {
  static initialize(web3, chainId) {
    const keep = new Keep(web3, chainId)
    keep.initializeContracts()
    keep.initializeServices()

    return keep
  }

  /**
   * @param {Web3LibWrapper} _web3 The web3 lib wrapper.
   * @param {number} chainId Id of the current chain.
   * @param {BaseExchange} exchangeService The exchange service that provides
   * data from the external exchanges.
   */
  constructor(_web3, chainId, exchangeService = null) {
    this.web3 = _web3
    this.exchangeService = exchangeService || new UniswapV2Exchange()
    this.chainId = chainId
  }

  /** @type {BaseContract} */
  [KEEP_TOKEN_CONTRACT_NAME];

  /** @type {BaseContract} */
  [TOKEN_STAKING_CONTRACT_NAME];

  /** @type {BaseContract} */
  [TOKEN_GRANT_CONTRACT_NAME];

  /** @type {BaseContract} */
  [OPERATOR_CONTRACT_NAME];

  /** @type {BaseContract} */
  [KEEP_OPERATOR_STATISTICS_CONTRACT_NAME];

  /** @type {BaseContract} */
  [MANAGED_GRANT_FACTORY_CONTRACT_NAME];

  /** @type {BaseContract} */
  [KEEP_BONDING_CONTRACT_NAME];

  /** @type {BaseContract} */
  [TBTC_TOKEN_CONTRACT_NAME];

  /** @type {BaseContract} */
  [TBTC_SYSTEM_CONTRACT_NAME];

  /** @type {BaseContract} */
  [TOKEN_STAKING_ESCROW_CONTRACT_NAME];

  /** @type {BaseContract} */
  [BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME];

  /** @type {BaseContract} */
  [STAKING_PORT_BACKER_CONTRACT_NAME];

  /** @type {BaseContract} */
  [OLD_TOKEN_STAKING_CONTRACT_NAME];

  /** @type {BaseContract} */
  [LP_REWARDS_KEEP_ETH_CONTRACT_NAME];

  /** @type {BaseContract} */
  [LP_REWARDS_TBTC_ETH_CONTRACT_NAME];

  /** @type {BaseContract} */
  [LP_REWARDS_KEEP_TBTC_CONTRACT_NAME];

  /** @type {BaseContract} */
  [LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME];

  /** @type {BaseContract} */
  [KEEP_TOKEN_GEYSER_CONTRACT_NAME];

  /** @type {BaseContract} */
  [SADDLE_SWAP_CONTRACT_NAME];

  /** @type {BaseContract} */
  [COV_TOKEN_CONTRACT_NAME];

  /** @type {BaseContract} */
  [ASSET_POOL_CONTRACT_NAME];

  /** @type {BaseContract} */
  [RISK_MANAGER_V1_CONTRACT_NAME];

  /** @type {BaseContract} */
  [REWARDS_POOL_CONTRACT_NAME];

  /** @type {BaseContract} */
  [THRESHOLD_STAKING_CONTRACT_NAME];

  /** @type {BaseContract} */
  [THRESHOLD_KEEP_STAKE_CONTRACT_NAME];

  /** @type {BaseContract} */
  [SIMPLE_PRE_APPLICATION_CONTRACT_NAME]

  initializeContracts = () => {
    const getDeploymentInfo = (artifact) => {
      const { networks } = artifact
      let deploymentTxnHash = null
      let address = null
      let deployedAtBlock = null
      if (networks) {
        // If the JSON artifact has the `networks` field, it means that an
        // artifact was built by Truffle.
        const networkId = Object.keys(networks)[0]
        address = networks[networkId].address
        deploymentTxnHash = networks[networkId].transactionHash
      } else {
        // Otherwise, it means an artifact was built by Hardhat.
        address = artifact.address
        deploymentTxnHash = artifact.transactionHash
        deployedAtBlock =
          artifact.receipt && artifact.receipt.blockNumber
            ? artifact.receipt.blockNumber
            : 1
      }
      return { address, deploymentTxnHash, deployedAtBlock }
    }

    for (const [contractName, { artifact }] of Object.entries(contracts)) {
      const { address, deploymentTxnHash, deployedAtBlock } =
        getDeploymentInfo(artifact)
      this[contractName] = this.web3.createContractInstance(
        artifact.abi,
        address,
        deploymentTxnHash,
        deployedAtBlock
      )
    }

    this.saddleSwapContract = this.web3.createContractInstance(
      SaddleSwapArtifact.abi,
      "0x4f6A43Ad7cba042606dECaCA730d4CE0A57ac62e",
      null,
      1
    )

    this._initializePREContract()
  }

  initializeServices = () => {
    this.coveragePoolV1 = new CoveragePoolV1(
      this.assetPoolContract,
      this.covTokenContract,
      this.keepTokenContract,
      this.rewardsPoolContract,
      this.riskManagerV1Contract,
      this.exchangeService,
      this.web3
    )
    this.keepToTStaking = new KeepToTStaking(
      this.thresholdStakingContract,
      this.thresholdKeepStakeContract,
      this.simplePREApplicationContract,
      this.web3
    )
  }

  _initializePREContract = () => {
    const preContractAddress = PRE.PRE_ADDRESSESS[this.chainId]
    const txHash =
      this.chainId === 1 ? PRE.MAINNET_PRE_DEPLOYMENT_TX_HASH : null
    const deploymentBlock =
      this.chainId === 1 ? PRE.MAINNET_PRE_DEPLOYMENT_BLOCK : 1

    this.simplePREApplicationContract = this.web3.createContractInstance(
      SimplePreApplicationAbi,
      preContractAddress,
      txHash,
      deploymentBlock
    )
  }

  setProvider = (provider) => {
    this.web3.setProvider(provider)

    for (const [contractName] of Object.entries(contracts)) {
      this[contractName].setProvider(provider)
    }
    this.saddleSwapContract.setProvider(provider)
  }

  /**
   * @param {string} acc Account address
   */
  set defaultAccount(acc) {
    this.web3.defaultAccount = acc

    for (const [contractName] of Object.entries(contracts)) {
      this[contractName].defaultAccount = acc
    }
  }
}

export default Keep
