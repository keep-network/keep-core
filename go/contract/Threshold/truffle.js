require('babel-register');
require('babel-polyfill');

module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",
      port: 9545,
      network_id: "*", // Match any network id
      gas: 6000000
, from: "0x627306090abab3a6e1400e9345bc60c78a8bef57"
    },
    testnet: {
      host: "10.51.244.207",
      port: 8545,
      network_id: "*",
      gas: 4712388
, from: "0x93d3299712e81aeb05feb28d8571ca0ed5c08c56"
    }
  }
};
