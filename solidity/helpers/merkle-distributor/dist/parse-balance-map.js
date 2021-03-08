"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.parseBalanceMap = void 0;
const ethers_1 = require("ethers");
const balance_tree_1 = require("./balance-tree");
const { isAddress, getAddress } = ethers_1.utils;
function parseBalanceMap(balances) {
    // if balances are in an old format, process them
    const balancesInNewFormat = Array.isArray(balances)
        ? balances
        : Object.keys(balances).map((account) => ({
            address: account,
            earnings: `0x${balances[account].toString(16)}`,
            reasons: "",
        }));
    const dataByAddress = balancesInNewFormat.reduce((memo, { address: account, earnings, reasons }) => {
        if (!isAddress(account)) {
            throw new Error(`Found invalid address: ${account}`);
        }
        const parsed = getAddress(account);
        if (memo[parsed])
            throw new Error(`Duplicate address: ${parsed}`);
        const parsedNum = ethers_1.BigNumber.from(earnings);
        if (parsedNum.lte(0))
            throw new Error(`Invalid amount for account: ${account}`);
        const flags = {
            isSOCKS: reasons.includes("socks"),
            isLP: reasons.includes("lp"),
            isUser: reasons.includes("user"),
        };
        memo[parsed] = Object.assign({ amount: parsedNum }, (reasons === "" ? {} : { flags }));
        return memo;
    }, {});
    const sortedAddresses = Object.keys(dataByAddress).sort();
    // construct a tree
    const tree = new balance_tree_1.default(sortedAddresses.map((address) => ({
        account: address,
        amount: dataByAddress[address].amount,
    })));
    // generate claims
    const claims = sortedAddresses.reduce((memo, address, index) => {
        const { amount, flags } = dataByAddress[address];
        memo[address] = Object.assign({ index, amount: amount.toHexString(), proof: tree.getProof(index, address, amount) }, (flags ? { flags } : {}));
        return memo;
    }, {});
    const tokenTotal = sortedAddresses.reduce((memo, key) => memo.add(dataByAddress[key].amount), ethers_1.BigNumber.from(0));
    return {
        merkleRoot: tree.getHexRoot(),
        tokenTotal: tokenTotal.toHexString(),
        claims,
    };
}
exports.parseBalanceMap = parseBalanceMap;
