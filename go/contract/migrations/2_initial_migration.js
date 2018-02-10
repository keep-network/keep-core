// var ProofOfExistence1 = artifacts.require("./ProofOfExistence1.sol");
// var Token = artifacts.require("./token.sol");
// var Greeter = artifacts.require("./greeter.sol");
// var GenRequestID = artifacts.require("./GenRequestID.sol");
var KStart = artifacts.require("./KStart.sol");

// console.log ( "Token=", Token );

module.exports = function(deployer) {
	// deployer.deploy(ProofOfExistence1);
	// deployer.deploy(Token);
	// deployer.deploy(Greeter);
	// deployer.deploy(GenRequestID).then(function () {
    //        deployer.deploy(KStart);
    //    });
	// deployer.deploy(GenRequestID);
    deployer.deploy(KStart);
};
