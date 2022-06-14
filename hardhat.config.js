require('@nomiclabs/hardhat-ethers')
require('@nomiclabs/hardhat-waffle')
require('hardhat-deploy')
require('solidity-coverage')
require('hardhat-gas-reporter')
require('dotenv').config();

module.exports = {
  solidity: '0.8.9',
  // namedAccounts: {
  //   deployer: {
  //     default: 0,
  //     ropsten: process.env.CONTRACT_OWNER_ACCOUNT_PUBLIC_KEY
  //   }
  // },
  // networks: {
  //   ropsten: {
  //     url: process.env.ROPSTEN_API_URL,
  //     chainId: 3,
  //     accounts: [process.env.CONTRACT_OWNER_ACCOUNT_PRIV_KEY]
  //   }
  // },
  // external: {
  //   contracts: [
  //     {
  //       artifacts: 'node_modules/@threshold-network/solidity-contracts/artifacts',
  //     }
  //   ]
  // }
}
