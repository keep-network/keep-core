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
      host: "10.51.244.207",
      port: 8545,
      network_id: "*",
      gas: 4712388
    }
  }
};
