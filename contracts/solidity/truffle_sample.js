require('babel-register');
require('babel-polyfill');

module.exports = {
  networks: {
    development: {
      host: "localhost",
      port: 8545,
      network_id: "*"
    },
    testnet: {
      host: "https://testnet.keep.network",
      port: 443,
      network_id: "*",
      gas: 4712388
    },
    keep_dev: {
      host: "localhost",
      port: 8545,
      network_id: "*",
      gas: 4712388,
      from: "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB"
    }
  },
  compilers: {
    solc: {
      version: "0.5.4"
    }
  }
};
