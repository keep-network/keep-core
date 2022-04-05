import "@keep-network/hardhat-helpers"
import "@keep-network/hardhat-local-networks-config"
import "@nomiclabs/hardhat-waffle"
import "@typechain/hardhat"
import "hardhat-deploy"
import "@tenderly/hardhat-tenderly"
import "hardhat-gas-reporter"
import "hardhat-contract-sizer"
import "hardhat-dependency-compiler"

import "./tasks"
import { task } from "hardhat/config"
import { TASK_TEST } from "hardhat/builtin-tasks/task-names"

import type { HardhatUserConfig } from "hardhat/config"

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
  operatorsCount: 100,
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
    settings: {
      outputSelection: {
        "*": {
          "*": ["storageLayout"],
        },
      },
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
      tags: ["local"],
      // we use higher gas price for tests to obtain more realistic results
      // for gas refund tests than when the default hardhat ~1 gwei gas price is
      // used
      gasPrice: 200000000000, // 200 gwei
    },
    development: {
      url: "http://localhost:8545",
      chainId: 1101,
      tags: ["local"],
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
      default: 1, // take the second account
      // mainnet: ""
    },
    governance: {
      default: 2,
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
    // deployments: {
    //   // For hardhat environment we can fork the mainnet, so we need to point it
    //   // to the contract artifacts.
    //   hardhat: process.env.FORKING_URL ? ["./external/mainnet"] : [],
    //   // For development environment we expect the local dependencies to be linked
    //   // with `yarn link` command.
    //   development: ["node_modules/@keep-network/keep-core/artifacts"],
    //   ropsten: ["node_modules/@keep-network/keep-core/artifacts"],
    //   mainnet: ["./external/mainnet"],
    // },
  },
  dependencyCompiler: {
    paths: [
      "@threshold-network/solidity-contracts/contracts/token/T.sol",
      "@threshold-network/solidity-contracts/contracts/staking/TokenStaking.sol",
      "@keep-network/random-beacon/contracts/api/IRandomBeacon.sol",
    ],
    keep: true,
  },
  contractSizer: {
    alphaSort: true,
    disambiguatePaths: false,
    runOnCompile: true,
    strict: true,
    except: ["TokenStaking$"],
  },
  mocha: {
    timeout: 60000,
  },
  typechain: {
    outDir: "typechain",
  },
}

task(TASK_TEST, "Runs mocha tests").setAction(async (args, hre, runSuper) => {
  // eslint-disable-next-line @typescript-eslint/no-var-requires,global-require
  const { constants } = require("./test/fixtures")

  if (testConfig.operatorsCount < constants.groupSize) {
    throw new Error(
      "not enough accounts predefined for configured group size: " +
        `expected group size: ${constants.groupSize} ` +
        `number of predefined accounts: ${testConfig.operatorsCount}`
    )
  }

  return runSuper(args)
})

export default config
