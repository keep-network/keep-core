const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol');

module.exports = async function() {
    try {
        const contract = await KeepRandomBeaconOperator.deployed();

        for (let block = await web3.eth.getBlockNumber(); ; block++) {
            while (block > await web3.eth.getBlockNumber()) {
                await wait(1000);
            }

            const events = await contract.getPastEvents('allEvents', {
                fromBlock: block,
                toBlock: block
            });

            for (let i = 0; i < events.length; i++) {
                const event = events[i];
                const transactionHash = event.transactionHash;
                const gasUsed = (await web3.eth.getTransactionReceipt(transactionHash)).gasUsed;
                console.log(`Event [${event.event}] with transaction [${transactionHash}] - gas used: ${parseInt(gasUsed, 16)}`);
            }
        }

    } catch(error) {
        console.log(error)
    }
};

function wait(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}