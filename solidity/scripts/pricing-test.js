const KeepRandomBeaconServiceImplV1 = artifacts.require("KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService.sol');
const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol');
const KeepRandomBeaconOperatorStatistics = artifacts.require('KeepRandomBeaconOperatorStatistics.sol');
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
    const contractStatistics = await KeepRandomBeaconOperatorStatistics.deployed();

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
                prevRewards[i] = (await availableRewards(accounts[i], contractOperator, contractStatistics)).toString();
            }


            let rewardsSum = web3.utils.toBN(0);
            for (let i = 0; i < accounts.length; i++) {
                rewardsSum = rewardsSum.add(web3.utils.toBN(prevRewards[i]));
            }

            const serviceContractBalance = await web3.eth.getBalance(contractService.address);
            const dkgFeePool = await contractService.dkgFeePool();
            const requestSubsidyFeePool = await contractService.requestSubsidyFeePool();
            const dkgContributionMargin = await contractService.dkgContributionMargin();

            const serviceContractSummary = new ServiceContractSummary(
                serviceContractBalance,
                dkgFeePool.toString(),
                requestSubsidyFeePool.toString(),
                dkgContributionMargin.toString(),
                dkgFeePool.add(requestSubsidyFeePool).toString() === serviceContractBalance
            )

            console.log("Service Contract Summary (before request)");
            console.table([serviceContractSummary]);
            console.log("\n");

            const operatorContractBalance = await web3.eth.getBalance(contractOperator.address);
            const dkgSubmitterReimbursementFee = await contractOperator.dkgSubmitterReimbursementFee();

            const operatorContractSummary = new OperatorContractSummary(
                operatorContractBalance,
                rewardsSum.toString(),
                dkgSubmitterReimbursementFee.toString(),
                rewardsSum.add(dkgSubmitterReimbursementFee).toString() === operatorContractBalance
            )

            console.log("Operator Contract Summary (before request)");
            console.table([operatorContractSummary]);
            console.log("\n");

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
                requestorAccountBalanceChange,
                operatorContractBalance,
                rewardsSum.toString()
            );

            console.log("Pricing Summary");
            console.table([pricingSummary]);
            console.log("\n");
            let file = pricingSummary.toString();

            const clientsTable = new Array(accounts.length);

            for (let i = 0; i < accounts.length; i++) {
                const address = accounts[i];
                const balance = await web3.eth.getBalance(address);
                const balanceChange = web3.utils.toBN(balance).sub(web3.utils.toBN(prevBalances[i])).toString();

                const reward = (await availableRewards(address, contractOperator, contractStatistics)).toString();
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

            console.log("Clients Summary");
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

async function availableRewards(account, contractOperator, contractStatistics) {
    const expiredGroupCount = (await contractOperator.getFirstActiveGroupIndex()).toNumber();
    const activeGroupCount = (await contractOperator.numberOfGroups()).toNumber();
    const totalGroupCount = expiredGroupCount + activeGroupCount;
    const groupsPublicKeys = new Array(totalGroupCount);

    for (let groupIndex = 0; groupIndex < totalGroupCount; groupIndex++) {
        groupsPublicKeys[groupIndex] = await contractOperator.getGroupPublicKey(groupIndex);
    }

    let accountRewards = web3.utils.toBN(0);
    for (let i = 0; i < groupsPublicKeys.length; i++) {
        const groupMembersCount = await contractStatistics.countGroupMembership(groupsPublicKeys[i], account);
        const groupMemberReward = await contractOperator.getGroupMemberRewards(groupsPublicKeys[i]);
        accountRewards = accountRewards.add(groupMembersCount.mul(groupMemberReward));
    }

    return accountRewards;
}

function ServiceContractSummary(
    balance,
    dkgFeePool,
    requestSubsidyFeePool,
    dkgContributionMargin,
    hasCorrectBalance
) {
    this.balance = balance,
    this.dkgFeePool = dkgFeePool
    this.requestSubsidyFeePool = requestSubsidyFeePool
    this.dkgContributionMargin = dkgContributionMargin
    this.hasCorrectBalance = hasCorrectBalance
}

function OperatorContractSummary(
    balance,
    sumOfRewards,
    dkgSubmitterReimbursementFee,
    hasCorrectBalance
) {
    this.balance = balance
    this.sumOfRewards = sumOfRewards
    this.dkgSubmitterReimbursementFee = dkgSubmitterReimbursementFee
    this.hasCorrectBalance = hasCorrectBalance
}

function PricingSummary(
    callbackGas,
    entryFeeEstimate,
    relayRequestTransactionCost,
    requestorAccountBalance,
    requestorAccountBalanceChange,
    operatorContractBalance,
    sumOfRewards
) {
    this.callbackGas = callbackGas
    this.entryFeeEstimate = entryFeeEstimate
    this.relayRequestTransactionCost = relayRequestTransactionCost
    this.requestorAccountBalance = requestorAccountBalance
    this.requestorAccountBalanceChange = requestorAccountBalanceChange
    this.operatorContractBalance = operatorContractBalance
    this.sumOfRewards = sumOfRewards
}

function PricingClient(address, balance, balanceChange, reward, rewardChange) {
    this.address = address,
    this.balance = balance,
    this.balanceChange = balanceChange,
    this.reward = reward,
    this.rewardChange = rewardChange
}

ServiceContractSummary.prototype.toString = function serviceContractSummaryToString() {
    return '' + this.balance + ', ' + 
        this.dkgFeePool + ', ' + 
        this.requestSubsidyFeePool + ', ' + 
        this.dkgContributionMargin + ', ' + 
        this.hasCorrectBalance + ', ';
}

OperatorContractSummary.prototype.toString = function operatorContractSummaryToString() {
    return '' + this.balance + ', ' +
        this.sumOfRewards + ', ' + 
        this.dkgSubmitterReimbursementFee + ', ' +
        this.hasCorrectBalance + ', ';
}

PricingSummary.prototype.toString = function pricingSummaryToString() {
    return '' + this.callbackGas + ', ' +
        this.entryFeeEstimate + ', ' +
        this.relayRequestTransactionCost + ', ' +
        this.requestorAccountBalance + ', ' +
        this.requestorAccountBalanceChange + ', ' +
        this.operatorContractBalance + ', ' + 
        this.sumOfRewards + ', ';
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
