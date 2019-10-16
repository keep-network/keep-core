const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');
var fs = require('fs')

// seed value for a relay entry
const seed = web3.utils.toBN('31415926535897932384626433832795028841971693993751058209749445923078164062862')

module.exports = async function() {
  const keepRandomBeaconService = await KeepRandomBeaconService.deployed()
  const contractService = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address)
  const callbackContract = await CallbackContract.deployed();
  // const delay = 600000 //10 min in milliseconds
  // const delay = 120000 //2 min in milliseconds
  const delay = 360000 //6 min in milliseconds
  
  let accounts = await web3.eth.getAccounts();
  let requestor = accounts[4]
  let count = 0
  let requestorAccountBalance = await web3.eth.getBalance(requestor)
  let requestorPrevAccountBalance = 0;

  for (;;) {
    try {
      console.log("--- count:", count)

      let callbackGas = await callbackContract.callback.estimateGas(seed)
      let entryFeeEstimate = await contractService.entryFeeEstimate(callbackGas)
      requestorPrevAccountBalance = requestorAccountBalance;

      let prevKeep1Balance = await web3.eth.getBalance(accounts[0])
      let prevKeep2Balance = await web3.eth.getBalance(accounts[1])
      let prevKeep3Balance = await web3.eth.getBalance(accounts[2])
      let prevKeep4Balance = await web3.eth.getBalance(accounts[3])
      
      await contractService.methods['requestRelayEntry(uint256,address,string,uint256)'](
        seed,
        callbackContract.address,
        "callback(uint256)",
        callbackGas,
        {value: entryFeeEstimate, from: requestor}
      );

      wait(delay); 

      requestorAccountBalance = await web3.eth.getBalance(requestor)
      
      let total = web3.utils.toBN(requestorPrevAccountBalance).sub(web3.utils.toBN(requestorAccountBalance)).toString()

      // Balances of keep nodes accounts after a relay entry generation.
      let keep1Balance = await web3.eth.getBalance(accounts[0])
      let keep1Earned = web3.utils.toBN(keep1Balance).sub(web3.utils.toBN(prevKeep1Balance)).toString()

      let keep2Balance = await web3.eth.getBalance(accounts[1])
      let keep2Earned = web3.utils.toBN(keep2Balance).sub(web3.utils.toBN(prevKeep2Balance)).toString()

      let keep3Balance = await web3.eth.getBalance(accounts[2])
      let keep3Earned = web3.utils.toBN(keep3Balance).sub(web3.utils.toBN(prevKeep3Balance)).toString()

      let keep4Balance = await web3.eth.getBalance(accounts[3])
      let keep4Earned = web3.utils.toBN(keep4Balance).sub(web3.utils.toBN(prevKeep4Balance)).toString()
      
      entryFeeEstimate = entryFeeEstimate.toString()
      let pricing = new Pricing(callbackGas, entryFeeEstimate, requestorAccountBalance, total,
        keep1Balance, keep1Earned, keep2Balance, keep2Earned,
        keep3Balance, keep3Earned, keep4Balance, keep4Earned)

      console.table([pricing])

      // Write data in 'pricing.txt' . 
      fs.appendFile("pricing.txt", pricing.toString(), (err) => {
        if (err) console.log(err);
      });

      count++
        
    } catch(error) {
      console.error('Request failed with', error)
    }
  }
}

function Pricing(callbackGas, entryFeeEstimate, requestorAccountBalance, totalForRelayEntry,
  keep1Balance, keep1Earned, keep2Balance, keep2Earned, 
  keep3Balance, keep3Earned, keep4Balance, keep4Earned) {
    this.callbackGas = callbackGas,
    this.entryFeeEstimate = entryFeeEstimate,
    this.requestorAccountBalance = requestorAccountBalance,
    this.totalForRelayEntry = totalForRelayEntry,
    this.keep1Balance = keep1Balance,
    this.keep1Earned = keep1Earned,
    this.keep2Balance = keep2Balance,
    this.keep2Earned = keep2Earned,
    this.keep3Balance = keep3Balance,
    this.keep3Earned = keep3Earned,
    this.keep4Balance = keep4Balance,
    this.keep4Earned = keep4Earned
  }

Pricing.prototype.toString = function pricingToString() {
  return '' + this.callbackGas + ', ' + this.entryFeeEstimate + ', ' + this.requestorAccountBalance + ', ' + this.totalForRelayEntry + ', '
  + this.keep1Balance + ', ' + this.keep1Earned + ', ' + this.keep2Balance + ', ' + this.keep2Earned
  + this.keep3Balance + ', ' + this.keep3Earned + ', ' + this.keep4Balance + ', ' + this.keep4Earned + '\n';
}


function wait(ms){
  var start = new Date().getTime();
  var end = start;
  while(end < start + ms) {
    end = new Date().getTime();
 }
}
