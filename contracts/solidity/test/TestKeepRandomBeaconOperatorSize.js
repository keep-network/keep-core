contract('KeepRandomBeaconOperator', (_) => {

  // The purpose of this test is to warn us when the operator contract bytecode
  // gets too large. We want to keep some safety margin on the operator contract
  // to be able to implement required fixes quickly, if needed.
  //
  // Bear in mind that in Solidity, the maximum size of a contract is restricted 
  // to 24 KB by EIP 170
  it("should not have its bytecode too large", () => {
    let KeepRandomBeaconOperator = artifacts.require("KeepRandomBeaconOperator.sol");

    return KeepRandomBeaconOperator.deployed().then((instance) => {
      let bytecode = instance.constructor._json.bytecode;
      let deployedBytecode = instance.constructor._json.deployedBytecode;
      
      let bytecodeSize = bytecode.length / 2; // size in bytes
      let deployedBytecodeSize = deployedBytecode.length / 2; // size in bytes

      const maxSafeBytecodeSize = 22482

      console.log(
        "KeepRandomBeaconOperator size of bytecode in bytes = ", 
        bytecodeSize
      );
      console.log(
        "KeepRandomBeaconOperator size of deployed bytecode in bytes = ", 
        deployedBytecodeSize
      );
      console.log(
        "KeepRandomBeaconOperator initialization and constructor code in bytes = ", 
        bytecodeSize - deployedBytecodeSize
      );

      assert.isBelow(
        bytecodeSize, 
        maxSafeBytecodeSize,
        "KeepRandomBeaconOperator bytecode is getting too large"
      )
    });  
  });
});