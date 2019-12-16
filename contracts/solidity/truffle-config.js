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
      host: "localhost",
      port: 8545,
      network_id: "*",
      from: "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB"
    },
    keep_dev_vpn: {
      host: "eth-tx-node.default.svc.cluster.local",
      port: 8545,
      network_id: "*",
      from: "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB"
    },
    keep_test: {
      host: "localhost",
      port: 8545,
      network_id: "*",
      from: "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB"
    },
    ropsten: {
      provider: function() {
        return new HDWalletProvider("EBAE221D3C6A4707B1B00927CE9DD6F866DC426658842CE3CFF5EBDAC2BF6000", "https://ropsten.infura.io/v3/59fb36a36fa4474b890c13dd30038be5")
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
      version: "0.5.4",
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  }
};
