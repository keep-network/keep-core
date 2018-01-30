//module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!
//};
//module.exports = {
//  networks: {
//    development: {
//      host: "192.168.0.139",
//      port: 8545,
//      network_id: "58342" // "*" to match any network id
//, from: "0x9980ecddef53089390136fde20feb7e03125c441"
//, gas: 4700000 
//, gasPrice: 12000
//    }
//  }
//};
module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",
      port: 9545,
      network_id: "*" // "*" to match any network id
, from: "0x627306090abab3a6e1400e9345bc60c78a8bef57"
, gas: 4700000 
    }
  }
};
