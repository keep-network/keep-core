import "@keep-network/hardhat-local-networks-config"
import "@keep-network/hardhat-helpers"
import "@nomiclabs/hardhat-ethers"
import "hardhat-deploy"
import "@tenderly/hardhat-tenderly"
import "@nomiclabs/hardhat-waffle"
import "hardhat-gas-reporter"
import "hardhat-contract-sizer"
import "@typechain/hardhat"
import "hardhat-dependency-compiler"
import "./tasks"

import { task } from "hardhat/config"

import type { HardhatUserConfig } from "hardhat/config"

const thresholdSolidityCompilerConfig = {
  version: "0.8.9",
  settings: {
    optimizer: {
      enabled: true,
      runs: 10,
    },
  },
}

// Configuration for testing environment.
export const testConfig = {
  // How many accounts we expect to define for non-staking related signers, e.g.
  // deployer, thirdParty, governance.
  // It is used as an offset for getting accounts for operators and stakes registration.
  nonStakingAccountsCount: 10,

  // How many roles do we need to define for staking, i.e. stakeOwner, stakingProvider,
  // operator, beneficiary, authorizer.
  stakingRolesCount: 5,

  // Number of operators to register. Should be at least the same as group size.
  operatorsCount: 64,
}

const config: HardhatUserConfig = {
  solidity: {
    compilers: [
      {
        version: "0.8.9",
        settings: {
          optimizer: {
            enabled: true,
          },
        },
      },
    ],
    overrides: {
      "@threshold-network/solidity-contracts/contracts/staking/TokenStaking.sol":
        thresholdSolidityCompilerConfig,
    },
  },
  paths: {
    artifacts: "./build",
  },
  networks: {
    hardhat: {
      forking: {
        // forking is enabled only if FORKING_URL env is provided
        enabled: !!process.env.FORKING_URL,
        // URL should point to a node with archival data (Alchemy recommended)
        url: process.env.FORKING_URL || "",
        // latest block is taken if FORKING_BLOCK env is not provided
        blockNumber: process.env.FORKING_BLOCK
          ? parseInt(process.env.FORKING_BLOCK, 10)
          : undefined,
      },
      accounts: {
        // Number of accounts that should be predefined on the testing environment.
        count:
          testConfig.nonStakingAccountsCount +
          testConfig.stakingRolesCount * testConfig.operatorsCount,
      },
      // we use higher gas price for tests to obtain more realistic results
      // for gas refund tests than when the default hardhat ~1 gwei gas price is
      // used
      gasPrice: 200000000000, // 200 gwei
      // Ignore contract size on deployment to hardhat network, to be able to
      // deploy stub contracts in tests.
      allowUnlimitedContractSize: process.env.TEST_USE_STUBS_BEACON === "true",
      tags: ["allowStubs"],
    },
    development: {
      url: "http://localhost:8545",
      chainId: 1101,
      tags: ["allowStubs"],
    },
    ropsten: {
      url: process.env.CHAIN_API_URL || "",
      chainId: 3,
      accounts: process.env.CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY
        ? [process.env.CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY]
        : undefined,
      tags: ["tenderly"],
    },
  },
  // // Define local networks configuration file path to load networks from the file.
  // localNetworksConfig: "./.hardhat/networks.ts",
  tenderly: {
    username: "thesis",
    project: "",
  },
  namedAccounts: {
    deployer: {
      default: 0, // take the first account as deployer
    },
    governance: {
      default: 1,
      // mainnet: ""
    },
  },
  external: {
    contracts: [
      {
        artifacts:
          "node_modules/@threshold-network/solidity-contracts/export/artifacts",
        deploy:
          "node_modules/@threshold-network/solidity-contracts/export/deploy",
      },
    ],
    deployments: {
      // For hardhat environment we can fork the mainnet, so we need to point it
      // to the contract artifacts.
      hardhat: process.env.FORKING_URL ? ["./external/mainnet"] : [],
      // For development environment we expect the local dependencies to be linked
      // with `yarn link` command.
      development: [
        "node_modules/@keep-network/keep-core/artifacts",
        "node_modules/@threshold-network/solidity-contracts",
      ],
      ropsten: ["node_modules/@keep-network/keep-core/artifacts"],
      mainnet: ["./external/mainnet"],
    },
  },
  dependencyCompiler: {
    paths: [
      "@threshold-network/solidity-contracts/contracts/token/T.sol",
      "@threshold-network/solidity-contracts/contracts/staking/TokenStaking.sol",
    ],
    keep: true,
  },
  contractSizer: {
    alphaSort: true,
    disambiguatePaths: false,
    runOnCompile: true,
    strict: true,
    except: ["^contracts/test"],
  },
  mocha: {
    timeout: 60000,
  },
  typechain: {
    outDir: "typechain",
  },
  gasReporter: {
    currency: "USD",
    coinmarketcap: process.env.COINMARKETCAP_API_KEY,
  },
}

task("check-accounts-count", "Checks accounts count").setAction(async () => {
  // eslint-disable-next-line @typescript-eslint/no-var-requires,global-require
  const { constants } = require("./test/fixtures")

  if (testConfig.operatorsCount < constants.groupSize) {
    throw new Error(
      "not enough accounts predefined for configured group size: " +
        `expected group size: ${constants.groupSize} ` +
        `number of predefined accounts: ${testConfig.operatorsCount}`
    )
  }
})

export default config
