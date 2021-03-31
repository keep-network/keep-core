require('babel-register')
require('babel-polyfill')
const HDWalletProvider = require('@truffle/hdwallet-provider')

module.exports = {
  networks: {
    local: {
      host: 'localhost',
      port: 8545,
      network_id: '*',
    },
    sov: {
      provider: function () {
        return new HDWalletProvider(
          '2d61b31f93df83e90e78b61943019f3d03fd9f31901359a0e065a4c896eee23d',
          'wss://testnet.sovryn.app/ws'
        )
      },
      gas: 6700000,
      gasPrice: 60000000,
      skipDryRun: false,
      network_id: '*',
      timeoutBlocks: 100,
      deploymentPollingInterval: 4000,
      disableConfirmationListener: true,
    },
    // local: {
    //   provider: function () {
    //     return new HDWalletProvider(
    //       '4526476adb17c8f751fc4e71e23c4f5e7e2013cd62417b63825cb6cde8a42627',
    //       'HTTP://127.0.0.1:8545'
    //     )
    //   },
    //   gas: 6700000,
    //   gasPrice: 80000000,
    //   skipDryRun: false,
    //   network_id: '*',
    //   timeoutBlocks: 50,
    //   deploymentPollingInterval: 1000,
    // },
    keep_dev: {
      provider: function () {
        return new HDWalletProvider(
          process.env.CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY,
          'http://localhost:8545'
        )
      },
      gas: 6721975,
      network_id: 1101,
    },

    keep_dev_vpn: {
      provider: function () {
        return new HDWalletProvider(
          process.env.CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY,
          'http://eth-tx-node.default.svc.cluster.local:8545'
        )
      },
      gas: 6721975,
      network_id: 1101,
    },

    ropsten: {
      provider: function () {
        return new HDWalletProvider(
          process.env.CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY,
          process.env.ETH_HOSTNAME
        )
      },
      gas: 8000000,
      network_id: 3,
      skipDryRun: true,
    },
  },

  mocha: {
    useColors: true,
    reporter: 'eth-gas-reporter',
    reporterOptions: {
      currency: 'USD',
      gasPrice: 21,
      showTimeSpent: true,
    },
  },

  compilers: {
    solc: {
      version: '0.5.17',
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
}
