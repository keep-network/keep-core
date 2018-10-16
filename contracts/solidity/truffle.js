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
    }
    , test3: {
      host: "192.168.0.199",
      port: 8545,
      network_id: "11011"
      , gas: 4712388
      , from: "0x023e291a99d21c944a871adcc44561a58f99bdbc"
    }
    , oldTest3: {
      host: "192.168.0.199",
      port: 8545,
      network_id: "*"
      , from: "0x6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
    }
  }
};
