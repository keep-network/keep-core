const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

module.exports = async function() {

  const callbackContract = await CallbackContract.deployed();
  
  try {
    let randomBeacon = await callbackContract.lastEntry();
    console.log("Random beacon: ", randomBeacon.toString())

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
