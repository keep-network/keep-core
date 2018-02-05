require('babel-register');
require('babel-polyfill');

module.exports = {
  networks: {
    development: {
      host: "localhost",
      port: 8545,
      network_id: "*", // Match any network id
      gas: 6000000
    },
    testnet: {
      host: "testnet.keep.network",
      port: 443,
      network_id: "1101",
      gas: 4712388
    }
  }
};
