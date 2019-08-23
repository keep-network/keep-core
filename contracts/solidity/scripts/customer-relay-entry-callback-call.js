const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

// seed value for a relay entry
const seed = web3.utils.toBN('31415926535897932384626433832795028841971693993751058209749445923078164062862')

module.exports = async function() {
  const keepRandomBeaconService = await KeepRandomBeaconService.deployed()
  const contractService = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address)
  const callbackContract = await CallbackContract.deployed();
  
  try {
    // estimates the gas used by a callback function when the method would be executed on the chain
    // a customer must provide a callback contract with a callback function
    let callbackGas = await callbackContract.callback.estimateGas(seed);
    // minimum payment for a relay entry request in Wei. Needs to include callbackGas to 
    // cover estimated gas for his callback.
    let minimumPayment = await contractService.minimumPayment(callbackGas)
    // account(s) should have enough ether to cover the cost of a relay entry
    let accounts = await web3.eth.getAccounts();
    
    // creates a transaction object and call requestRelayEntry for a new relay entry which would be stored
    // under callbackContract.address
    await contractService.methods['requestRelayEntry(uint256,address,string,uint256)'](
      seed,
      callbackContract.address,
      "callback(uint256)",
      callbackGas,
      {value: minimumPayment, from: accounts[0]}
    );
  } catch(error) {
    console.error('Request failed with', error)
  }

  process.exit()
}

/* 
 * Callback contract example
 *
contract CallbackContract {

  uint256 internal _lastEntry; // result of a relay entry request (beacon)

  function callback(uint256 requestResponse) public {
      _lastEntry = requestResponse;
  }

  function lastEntry() public view returns (uint256){
      return _lastEntry;
  }
}
*/
