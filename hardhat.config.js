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
  solidity: '0.8.9',
  namedAccounts: {
    deployer: {
      "ropsten": 0
    },
    rewardsHolder: {
      "ropsten": "0xCe692F6fA86319Af43050fB7F09FDC43319F7612",
    },
    tokenContract: {
      "ropsten": "0x8562d01c9C0F1A8173360E48F50F6b9879c98Dc6",
    },
  },
  networks,
};
