const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');
const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol');
const fs = require('fs');

// MAKE SURE NONE OF THOSE ACCOUNTS IS A MINER ACCOUNT
const requestor = '0x146748a2b46b99ee1470b587bc9812ea59b79597';

const operators = [
    '0x65ea55c1f10491038425725dc00dffeab2a1e28a',
    '0x524f2e0176350d950fa630d9a5a59a0a190daf48',
    '0x3365d0ed0e526d3b1d8b417fc0fde5b1cef2f416',
    '0x7020a5556ba1ce5f92c81063a13d33512cf1305c'
];

const delay = 60000; //1 min in milliseconds


module.exports = async function() {
    const keepRandomBeaconService = await KeepRandomBeaconService.deployed();
    const contractService = await KeepRandomBeaconServiceImplV1.at(keepRandomBeaconService.address);
    const contractOperator = await KeepRandomBeaconOperator.deployed();

    let accounts = operators.slice();
    accounts.push(requestor);

    let count = 0;
    let requestorAccountBalance = await web3.eth.getBalance(requestor);
    let requestorPrevAccountBalance = 0;

    for (;;) {
        try {
            console.log("---------- count: " + count + " ----------\n");

            requestorPrevAccountBalance = requestorAccountBalance;

            const prevBalances = new Array(accounts.length);
            const prevRewards = new Array(accounts.length);

            for (let i = 0; i < accounts.length; i++) {
                prevBalances[i] = await web3.eth.getBalance(accounts[i]);
                prevRewards[i] = (await availableRewards(accounts[i], contractOperator)).toString();
            }

            let gasPrice = await web3.eth.getGasPrice();

            let callbackGas = 0;
            let entryFeeEstimate = await contractService.entryFeeEstimate(callbackGas);
            let tx = await contractService.methods['requestRelayEntry()'](
                {value: entryFeeEstimate, from: requestor}
            );

            await wait(delay);

            let txGasCost = tx.receipt.gasUsed;
            let txCost = web3.utils.toBN(gasPrice).mul(web3.utils.toBN(txGasCost));

            requestorAccountBalance = await web3.eth.getBalance(requestor);

            const requestorAccountBalanceChange = web3.utils.toBN(requestorAccountBalance).sub(web3.utils.toBN(requestorPrevAccountBalance)).toString();

            const pricingSummary = new PricingSummary(
                callbackGas,
                entryFeeEstimate.toString(),
                txCost.toString(),
                requestorAccountBalance,
                requestorAccountBalanceChange
            );

            console.log("Summary");
            console.table([pricingSummary]);
            console.log("\n");
            let file = pricingSummary.toString();

            const clientsTable = new Array(accounts.length);

            for (let i = 0; i < accounts.length; i++) {
                const address = accounts[i];
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
    const expiredGroupCount = (await contractOperator.getFirstActiveGroupIndex()).toNumber();
    const activeGroupCount = (await contractOperator.numberOfGroups()).toNumber();
    const totalGroupCount = expiredGroupCount + activeGroupCount;
    const groupsPublicKeys = new Array(totalGroupCount);

    for (let groupIndex = 0; groupIndex < totalGroupCount; groupIndex++) {
        groupsPublicKeys[groupIndex] = await contractOperator.getGroupPublicKey(groupIndex);
    }

    let accountRewards = web3.utils.toBN(0);
    for (let i = 0; i < groupsPublicKeys.length; i++) {
        const groupMembersCount = (await contractOperator.getGroupMemberIndices(groupsPublicKeys[i], account)).length;
        const groupMemberReward = await contractOperator.getGroupMemberRewards(groupsPublicKeys[i]);
        accountRewards = accountRewards.add(web3.utils.toBN(groupMembersCount).mul(groupMemberReward));
    }

    return accountRewards;
}

function PricingSummary(
    callbackGas,
    entryFeeEstimate,
    relayRequestTransactionCost,
    requestorAccountBalance,
    requestorAccountBalanceChange
) {
    this.callbackGas = callbackGas,
    this.entryFeeEstimate = entryFeeEstimate,
    this.relayRequestTransactionCost = relayRequestTransactionCost,
    this.requestorAccountBalance = requestorAccountBalance,
    this.requestorAccountBalanceChange = requestorAccountBalanceChange
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
        this.relayRequestTransactionCost + ', ' +
        this.requestorAccountBalance + ', ' +
        this.requestorAccountBalanceChange + ', ';
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