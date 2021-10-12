const { BN } = require('@openzeppelin/test-helpers');
const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');
const { toBN, generateSalt } = require('./helpers/utils');

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

async function makeDrop (token, wallets, amounts, deposit) {
    const salts = wallets.map(_ => generateSalt());
    const elements = wallets.map((w, i) => salts[i] + w.substr(2) + toBN(amounts[i]).toString(16, 64));
    const hashedElements = elements.map(keccak128).map(x => MerkleTree.bufferToHex(x));
    const tree = new MerkleTree(elements, keccak128, { hashLeaves: true, sort: true });
    const root = tree.getHexRoot();
    const leaves = tree.getHexLeaves();
    const proofs = leaves
        .map(tree.getHexProof, tree)
        .map(proof => '0x' + proof.map(p => p.substr(2)).join(''));

    const drop = await MerkleDrop128.new(token.address, root, tree.getDepth());
    await token.mint(drop.address, deposit);

    return { hashedElements, leaves, root, proofs, salts, drop };
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
        const accounts = Array(10).fill().map((_, i) => '0x' + (new BN(w1.substr(2), 16)).addn(i).toString('hex'));
        const amounts = Array(10).fill().map((_, i) => i + 1);
        const { hashedElements, leaves, root, proofs, salts, drop } = await makeDrop(this.token, accounts, amounts, 1000000);
        this.hashedElements = hashedElements;
        this.leaves = leaves;
        this.root = root;
        this.proofs = proofs;

        if (drop.contract.methods.verify) {
            await drop.contract.methods.verify(this.proofs[findSortedIndex(this, 0)], this.root, this.leaves[findSortedIndex(this, 0)]).send({ from: addr1 });
            for (let i = 0; i < 10; i++) {
                const callResult = await drop.verify(this.proofs[findSortedIndex(this, i)], this.root, this.leaves[findSortedIndex(this, i)]);
                expect(callResult.valid).to.be.true;
                //expect(callResult.index.toString()).to.be.bignumber.eq(findSortedIndex(this, i).toString());
            }
        }
        //await this.drop.claim(salts[0], accounts[0], 1, this.root, this.proofs[findSortedIndex(this, 0)]);
    });
});
