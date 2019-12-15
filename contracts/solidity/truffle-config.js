require('babel-register');
require('babel-polyfill');

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
      version: "0.5.14",
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  }
};
