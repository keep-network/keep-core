const MyTestContract = artifacts.require('MyTestContractOperator.sol');

module.exports = async function () {

    const myTestContract = await MyTestContract.at("0x15Ae06754d7E10Ed3D5b23bb934D9c633517577a");

  async function callFunctionOne() {
    try {
        console.log("f1 gas estimated: " + (await myTestContract.functionOne.estimateGas()));

        const startBlockNumber = await web3.eth.getBlock('latest').number        

        let tx = await myTestContract.functionOne();

        console.log("f1 gas used: " + tx.receipt.gasUsed);
    
        const eventList = await myTestContract.getPastEvents('MyEvent', {
            fromBlock: startBlockNumber,
            toBlock: 'latest',
        });
    

        console.log("f1 value: " + eventList[0].returnValues.value);
    
      } catch (err) {
          console.log(err)
      }
  }

  async function callFunctionTwo() {
    try {
        console.log("f2 gas estimated: " + (await myTestContract.functionTwo.estimateGas()));

        const startBlockNumber = await web3.eth.getBlock('latest').number
    
        let tx = await myTestContract.functionTwo();

        console.log("f2 gas used: " + tx.receipt.gasUsed);
    
        const eventList = await myTestContract.getPastEvents('MyEvent', {
            fromBlock: startBlockNumber,
            toBlock: 'latest',
        });
    

        console.log("f2 value: " + eventList[0].returnValues.value);
    
      } catch (err) {
          console.log(err)
      }
  }

  async function callFunctionThree() {
    try {
        console.log("f3 gas estimated: " + (await myTestContract.functionThree.estimateGas()));

        const startBlockNumber = await web3.eth.getBlock('latest').number
    
        let tx = await myTestContract.functionThree();

        console.log("f3 gas used: " + tx.receipt.gasUsed);
    
        const eventList = await myTestContract.getPastEvents('MyEvent', {
            fromBlock: startBlockNumber,
            toBlock: 'latest',
        });
    

        console.log("f3 value: " + eventList[0].returnValues.value);
    
      } catch (err) {
          console.log(err)
      }
  }

  async function callFunctionFour() {
    try {
        console.log("f4 gas estimated: " + (await myTestContract.functionFour.estimateGas()));

        const startBlockNumber = await web3.eth.getBlock('latest').number
    
        let tx = await myTestContract.functionFour();

        console.log('f4 tx hash: ' + tx.receipt.transactionHash);
        console.log("f4 gas used: " + tx.receipt.gasUsed);
    
        const eventList = await myTestContract.getPastEvents('MyEvent', {
            fromBlock: startBlockNumber,
            toBlock: 'latest',
        });
    

        console.log("f4 value: " + eventList[0].returnValues.value);
    
      } catch (err) {
          console.log(err)
      }
  }

  await callFunctionOne();
  await callFunctionTwo();
  await callFunctionThree();
  await callFunctionFour();
}
