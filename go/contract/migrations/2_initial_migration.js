var ProofOfExistence1 = artifacts.require("./ProofOfExistence1.sol");
var Token = artifacts.require("./token.sol");

// console.log ( "Token=", Token );

module.exports = function(deployer) {
  deployer.deploy(ProofOfExistence1);
  deployer.deploy(Token);
};
