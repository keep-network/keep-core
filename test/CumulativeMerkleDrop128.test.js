const { BN } = require('@openzeppelin/test-helpers');
const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');

const {
    shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234,
} = require('./behaviors/MerkleDrop.behavior');

const {
    shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234,
} = require('./behaviors/CumulativeMerkleDrop.behavior');

const TokenMock = artifacts.require('TokenMock');
const CumulativeMerkleDrop128 = artifacts.require('CumulativeMerkleDrop128');

function keccak128 (input) {
    return keccak256(input).slice(0, 16);
}

async function makeDrop (token, drop, wallets, amounts, deposit) {
    const elements = wallets.map((w, i) => w + web3.utils.padLeft(web3.utils.toHex(amounts[i]), 64).substr(2));
    const hashedElements = elements.map(keccak128).map(x => MerkleTree.bufferToHex(x));
    const tree = new MerkleTree(elements, keccak128, { hashLeaves: true, sort: true });
    const root = tree.getHexRoot();
    const leaves = tree.getHexLeaves();
    const proofs = leaves
        .map(tree.getHexProof, tree)
        .map(proof => '0x' + proof.map(p => p.substr(2)).join(''));

    await drop.setMerkleRoot(root);
    await token.mint(drop.address, deposit);

    return { hashedElements, leaves, root, proofs };
}

contract('CumulativeMerkleDrop128', async function ([_, w1, w2, w3, w4]) {
    const wallets = [w1, w2, w3, w4];

    function findSortedIndex (self, i) {
        return self.leaves.indexOf(self.hashedElements[i]);
    }

    beforeEach(async function () {
        this.token = await TokenMock.new('1INCH Token', '1INCH');
        this.drop = await CumulativeMerkleDrop128.new(this.token.address);
        await Promise.all(wallets.map(w => this.token.mint(w, 1)));
    });

    it('Benchmark 30000 wallets (merkle tree height 15)', async function () {
        const accounts = Array(30000).fill().map((_, i) => '0x' + (new BN(w1.substr(2), 16)).addn(i).toString('hex'));
        const amounts = Array(30000).fill().map((_, i) => i + 1);

        const { hashedElements, leaves, root, proofs } = await makeDrop(this.token, this.drop, accounts, amounts, 1000000);
        this.hashedElements = hashedElements;
        this.leaves = leaves;
        this.root = root;
        this.proofs = proofs;

        if (this.drop.contract.methods.verify) {
            await this.drop.contract.methods.verify(this.proofs[findSortedIndex(this, 0)], this.root, this.leaves[findSortedIndex(this, 0)]).send({ from: _ });
            expect(await this.drop.verify(this.proofs[findSortedIndex(this, 0)], this.root, this.leaves[findSortedIndex(this, 0)])).to.be.true;
        }
        // await this.drop.contract.methods.verifyAsm(this.proofs[findSortedIndex(this, 0)], this.root, this.leaves[findSortedIndex(this, 0)]).send({ from: _ });
        // expect(await this.drop.verifyAsm(this.proofs[findSortedIndex(this, 0)], this.root, this.leaves[findSortedIndex(this, 0)])).to.be.true;
        await this.drop.claim(accounts[0], 1, this.root, this.proofs[findSortedIndex(this, 0)]);
    });

    describe('Single drop for 4 wallets: [1, 2, 3, 4]', async function () {
        beforeEach(async function () {
            const { hashedElements, leaves, root, proofs } = await makeDrop(this.token, this.drop, [w1, w2, w3, w4], [1, 2, 3, 4], 10);
            this.hashedElements = hashedElements;
            this.leaves = leaves;
            this.root = root;
            this.proofs = proofs;
        });

        shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234('CMD', [w1, w2, w3, w4], findSortedIndex);
    });

    describe('Double drop for 4 wallets: [1, 2, 3, 4] + [2, 3, 4, 5] = [3, 5, 7, 9]', async function () {
        async function makeFirstDrop (self) {
            const { hashedElements, leaves, root, proofs } = await makeDrop(self.token, self.drop, [w1, w2, w3, w4], [1, 2, 3, 4], 1 + 2 + 3 + 4);
            self.hashedElements = hashedElements;
            self.leaves = leaves;
            self.root = root;
            self.proofs = proofs;
        }

        async function makeSecondDrop (self) {
            const { hashedElements, leaves, root, proofs } = await makeDrop(self.token, self.drop, [w1, w2, w3, w4], [3, 5, 7, 9], 2 + 3 + 4 + 5);
            self.hashedElements = hashedElements;
            self.leaves = leaves;
            self.root = root;
            self.proofs = proofs;
        }

        shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234('CMD', _, [w1, w2, w3, w4], findSortedIndex, makeFirstDrop, makeSecondDrop);
    });
});
