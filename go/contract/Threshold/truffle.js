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
      host: "10.51.245.75",
      port: 8545,
      network_id: "*",
      gas: 4712388
, from: "0x6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
    }
  }
};
