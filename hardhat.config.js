require('@nomiclabs/hardhat-ethers')
require('@nomiclabs/hardhat-waffle')
require("@nomiclabs/hardhat-etherscan")
require('hardhat-deploy')
require('solidity-coverage')
require('hardhat-gas-reporter')
require('dotenv').config();

const { networks, etherscan } = require('./hardhat.networks');

module.exports = {
  etherscan,
  solidity: {
    version: "0.8.9",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  namedAccounts: {
    deployer: {
      "goerli": 0,
      "mainnet": 0,
      "mainnet_test": 0
    },
    rewardsHolder: {
      "goerli": "0xCe692F6fA86319Af43050fB7F09FDC43319F7612",
      "mainnet": "0x9F6e831c8F8939DC0C830C6e492e7cEf4f9C2F5f",
      "mainnet_test": 0
    },
    tokenContract: {
      "goerli": "0x3f16380656cAE45D3f80D8833682d2b606eD094A",
      "mainnet": "0xCdF7028ceAB81fA0C6971208e83fa7872994beE5"
    },
    owner: {
      "goerli": 0,
      "mainnet": "0x9F6e831c8F8939DC0C830C6e492e7cEf4f9C2F5f",
      "mainnet_test": 0
    }
  },
  networks,
};
