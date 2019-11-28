const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');
const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol');
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');
const fs = require('fs');

// seed value for a relay entry
const seed = web3.utils.toBN('31415926535897932384626433832795028841971693993751058209749445923078164062862');

module.exports = async function() {
    const keepRandomBeaconService = await KeepRandomBeaconService.deployed();
    const contractService = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address);
    const contractOperator = await KeepRandomBeaconOperator.deployed();
    const callbackContract = await CallbackContract.deployed();
    const delay = 600000; //10 min in milliseconds
    const accountsCount = 4;
    const accounts = await web3.eth.getAccounts();
    const requestor = accounts[0];

    let count = 0;
    let requestorAccountBalance = await web3.eth.getBalance(requestor);
    let requestorPrevAccountBalance = 0;

    for (;;) {
        try {
            console.log("---------- count: " + count + " ----------\n");

            let callbackGas = await callbackContract.callback.estimateGas(seed);
            let entryFeeEstimate = await contractService.entryFeeEstimate(callbackGas);
            requestorPrevAccountBalance = requestorAccountBalance;

            const prevBalances = new Array(accountsCount);
            const prevRewards = new Array(accountsCount);

            for (let i = 0; i < accountsCount; i++) {
                prevBalances[i] = await web3.eth.getBalance(accounts[i+1]);
                prevRewards[i] = (await availableRewards(accounts[i+1], contractOperator)).toString();
            }

            await contractService.methods['requestRelayEntry(uint256,address,string,uint256)'](
                seed,
                callbackContract.address,
                "callback(uint256)",
                callbackGas,
                {value: entryFeeEstimate, from: requestor}
            );

            await wait(delay);

            requestorAccountBalance = await web3.eth.getBalance(requestor);

            const total = web3.utils.toBN(requestorPrevAccountBalance).sub(web3.utils.toBN(requestorAccountBalance)).toString();

            const pricingSummary = new PricingSummary(
                callbackGas,
                entryFeeEstimate.toString(),
                requestorAccountBalance,
                total
            );

            console.log("Summary");
            console.table([pricingSummary]);
            console.log("\n");
            let file = pricingSummary.toString();

            const clientsTable = new Array(accountsCount);

            for (let i = 0; i < accountsCount; i++) {
                const address = accounts[i+1];
                const balance = await web3.eth.getBalance(address);
                const balanceChange = web3.utils.toBN(balance).sub(web3.utils.toBN(prevBalances[i])).toString();

                const reward = (await availableRewards(address, contractOperator)).toString();
                const rewardChange = web3.utils.toBN(reward).sub(web3.utils.toBN(prevRewards[i])).toString();

                const pricingClient = new PricingClient(
                    address,
                    balance,
                    balanceChange,
                    reward,
                    rewardChange,
                );

                clientsTable[i] = pricingClient;
                file = file + pricingClient.toString();
            }

            console.log("Clients");
            console.table(clientsTable);
            console.log("\n");

            // Write data in 'pricing.txt' .
            fs.appendFile("pricing.txt", file + '\n', (err) => {
                if (err) console.log(err);
            });

            count++
        } catch(error) {
            console.error('Request failed with', error)
        }
    }
};

async function availableRewards(account, contractOperator) {
    const numberOfGroups = (await contractOperator.totalNumberOfGroups()).toNumber();
    const groupsPublicKeys = new Array(numberOfGroups);

    for (let groupIndex = 0; groupIndex < numberOfGroups; groupIndex++) {
        groupsPublicKeys[groupIndex] = await contractOperator.getGroupPublicKey(groupIndex);
    }

    let accountRewards = web3.utils.toBN(0);
    for (let i = 0; i < groupsPublicKeys.length; i++) {
        if (await contractOperator.isStaleGroup(groupsPublicKeys[i])) {
            const groupMembersCount = (await contractOperator.getGroupMemberIndices(groupsPublicKeys[i], account)).length;
            const groupMemberReward = await contractOperator.getGroupMemberRewards(groupsPublicKeys[i]);
            accountRewards = accountRewards.add(web3.utils.toBN(groupMembersCount).mul(groupMemberReward));
        }
    }

    return accountRewards;
}

function PricingSummary(
    callbackGas,
    entryFeeEstimate,
    requestorAccountBalance,
    totalForRelayEntry
) {
    this.callbackGas = callbackGas,
    this.entryFeeEstimate = entryFeeEstimate,
    this.requestorAccountBalance = requestorAccountBalance,
    this.totalForRelayEntry = totalForRelayEntry
}

function PricingClient(address, balance, balanceChange, reward, rewardChange) {
    this.address = address,
    this.balance = balance,
    this.balanceChange = balanceChange,
    this.reward = reward,
    this.rewardChange = rewardChange
}

PricingSummary.prototype.toString = function pricingSummaryToString() {
    return '' + this.callbackGas + ', ' +
        this.entryFeeEstimate + ', ' +
        this.requestorAccountBalance + ', ' +
        this.totalForRelayEntry + ', ';
};

PricingClient.prototype.toString = function pricingClientToString() {
    return '' + this.address + ', ' +
        this.balance + ', ' +
        this.balanceChange + ', ' +
        this.reward + ', ' +
        this.rewardChange + ', ';
};

function wait(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}