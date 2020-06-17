require('babel-register');
require('babel-polyfill');
const HDWalletProvider = require("@truffle/hdwallet-provider");

module.exports = {
  networks: {
    local: {
      host: "localhost",
      port: 8545,
      network_id: "*"
    },
    keep_dev: {
      provider: function() {
        return new HDWalletProvider(process.env.CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY, "http://localhost:8545")
      },
      gas: 6721975,
      network_id: 1101
    },

    keep_dev_vpn: {
      provider: function() {
        return new HDWalletProvider(process.env.CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY, "http://eth-tx-node.default.svc.cluster.local:8545")
      },
      gas: 6721975,
      network_id: 1101
    },

    ropsten: {
      provider: function() {
        return new HDWalletProvider(process.env.CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY, process.env.ETH_HOSTNAME)
      },
      gas: 6721975,
      network_id: 3
    }
  },

  mocha: {
    useColors: true,
    reporter: 'eth-gas-reporter',
    reporterOptions: {
      currency: 'USD',
      gasPrice: 21,
      showTimeSpent: true
    }
  },

  compilers: {
    solc: {
      version: "0.5.17",
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  }
};
