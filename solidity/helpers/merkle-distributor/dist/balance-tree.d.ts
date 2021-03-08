/// <reference types="node" />
import { BigNumber } from 'ethers';
export default class BalanceTree {
    private readonly tree;
    constructor(balances: {
        account: string;
        amount: BigNumber;
    }[]);
    static verifyProof(index: number | BigNumber, account: string, amount: BigNumber, proof: Buffer[], root: Buffer): boolean;
    static toNode(index: number | BigNumber, account: string, amount: BigNumber): Buffer;
    getHexRoot(): string;
    getProof(index: number | BigNumber, account: string, amount: BigNumber): string[];
}
