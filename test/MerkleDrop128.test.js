const { BN } = require('@openzeppelin/test-helpers');
const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');
const { toBN, generateSalt } = require('./helpers/utils');
const ethSigUtil = require('eth-sig-util');

const {
    shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234,
} = require('./behaviors/MerkleDrop.behavior');

const {
    shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234,
} = require('./behaviors/CumulativeMerkleDrop.behavior');

const TokenMock = artifacts.require('TokenMock');
const MerkleDrop128 = artifacts.require('MerkleDrop128');

function keccak128 (input) {
    return keccak256(input).slice(0, 16);
}

async function makeDrop (token, accountWithDropValues, deposit) {
    const elements = accountWithDropValues.map((w) => '0x' + w.account.substr(2) + toBN(w.amount).toString(16, 64));
    const hashedElements = elements.map(keccak128).map(x => MerkleTree.bufferToHex(x));
    const tree = new MerkleTree(elements, keccak128, { hashLeaves: true, sort: true });
    const root = tree.getHexRoot();
    const leaves = tree.getHexLeaves();
    const proofs = leaves
        .map(tree.getHexProof, tree)
        .map(proof => '0x' + proof.map(p => p.substr(2)).join(''));

    const drop = await MerkleDrop128.new(token.address, root, tree.getDepth());
    await token.mint(drop.address, deposit);

    return { hashedElements, leaves, root, proofs, drop };
}

contract('MerkleDrop128', async function ([addr1, w1, w2, w3, w4]) {
    const wallets = [w1, w2, w3, w4];

    function findSortedIndex (self, i) {
        return self.leaves.indexOf(self.hashedElements[i]);
    }

    beforeEach(async function () {
        this.token = await TokenMock.new('1INCH Token', '1INCH');
        await Promise.all(wallets.map(w => this.token.mint(w, 1)));
    });

    it.only('Benchmark 30000 wallets (merkle tree height 15)', async function () {
        const accountWithDropValues = [
            {
                account: addr1,
                amount: 1,
            },
            {
                account: w1,
                amount: 1,
            },
        ];
        const { hashedElements, leaves, root, proofs, salts, drop } = await makeDrop(this.token, accountWithDropValues, 1000000);
        this.hashedElements = hashedElements;
        this.leaves = leaves;
        this.root = root;
        this.proofs = proofs;
        const data = "hello";
        const signature = ethSigUtil.personalSign(Buffer.from('ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80', 'hex'), { data });
        await drop.claim(addr1, addr1, 1, this.proofs[findSortedIndex(this, 0)], signature);
    });
});
