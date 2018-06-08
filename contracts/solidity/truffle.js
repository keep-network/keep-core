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
    develop: {
      host: "localhost",
      port: 9545,
      network_id: "*", // Match any network id
	  gasPrice: 1,
      gas: 6000000
    },
    testnet: {
      host: "https://testnet.keep.network",
      port: 443,
      network_id: "*",
      gas: 4712388
    }
  }
};
