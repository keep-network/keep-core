require('babel-register');
require('babel-polyfill');

module.exports = {
  networks: {
    local: {
      host: "localhost",
      port: 8545,
      network_id: "*"
    },
    testnet: {
      host: "https://testnet.keep.network",
      port: 443,
      network_id: "*",
      gas: 4712388
    }
  },
  compilers: {
    solc: {
      version: "0.5.4"
    }
  }
};
