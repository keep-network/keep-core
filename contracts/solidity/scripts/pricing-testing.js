const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

// seed value for a relay entry
const seed = web3.utils.toBN('31415926535897932384626433832795028841971693993751058209749445923078164062862')

module.exports = async function() {
  const keepRandomBeaconService = await KeepRandomBeaconService.deployed()
  const contractService = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address)
  const callbackContract = await CallbackContract.deployed();
  const delay = 600000 //10 min in milliseconds
  const numberOfKeepNodes = 4;
  
  let accounts = await web3.eth.getAccounts();
  let accountFrom = accounts[4]
  let count = 0
  let currentAccountBalance = await web3.eth.getBalance(accountFrom)
  let prevAccountBalance = 0;

  try {
    for (;;) {
      console.log("--- count ---", count)
      
      let callbackGas = await callbackContract.callback.estimateGas(seed)
      let recommendedPayment = await contractService.minimumPayment(callbackGas)
      prevAccountBalance = currentAccountBalance;
      
      await contractService.methods['requestRelayEntry(uint256,address,string,uint256)'](
        seed,
        callbackContract.address,
        "callback(uint256)",
        callbackGas,
        {value: recommendedPayment, from: accountFrom}
      );

      console.log("estimated gas: ", callbackGas.toString())
      console.log("recommended payment: ", recommendedPayment.toString())

      currentAccountBalance = await web3.eth.getBalance(accountFrom)
      console.log("Requestor account: " + accountFrom + " has a balance of: " + currentAccountBalance)
      
      let total = web3.utils.toBN(prevAccountBalance).sub(web3.utils.toBN(currentAccountBalance)).toString()
      console.log("Total payment for a relay entry request: ", total)

      // Balances of keep nodes accounts after a relay entry generation.
      for (i = 0; i < numberOfKeepNodes; i++) {
        let accountBalance = await web3.eth.getBalance(accounts[i])
        console.log("account balance for node: " + i + " is: " + accountBalance)
      }

      count++
        
      wait(delay);  
    }
  } catch(error) {
    console.error('Request failed with', error)
  }

  process.exit()
}


function wait(ms){
  var start = new Date().getTime();
  var end = start;
  while(end < start + ms) {
    end = new Date().getTime();
 }
}
