const { BN, expectEvent } = require('@openzeppelin/test-helpers');
const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');

// const { gasspectEVM } = require('./helpers/profileEVM');

const {
    shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234,
} = require('./MerkleDrop.behavior');

const TokenMock = artifacts.require('TokenMock');
const CumulativeMerkleDrop = artifacts.require('CumulativeMerkleDrop');

function claimedEvent(account, amount) {
    return { account, amount };
}

async function makeDrop (token, drop, wallets, amounts, deposit) {
    const elements = wallets.map((w, i) => w + web3.utils.padLeft(web3.utils.toHex(amounts[i]), 64).substr(2));
    const hashedElements = elements.map(keccak256).map(x => MerkleTree.bufferToHex(x));
    const tree = new MerkleTree(elements, keccak256, { hashLeaves: true, sortPairs: true });
    const root = tree.getHexRoot();
    const leaves = tree.getHexLeaves();
    const proofs = leaves.map(tree.getHexProof, tree);

    await drop.setMerkleRoot(root);
    await token.mint(drop.address, deposit);

    return { hashedElements, leaves, root, proofs };
}

contract('CumulativeMerkleDrop', async function ([_, w1, w2, w3, w4]) {
    const wallets = [w1, w2, w3, w4];

    function findSortedIndex (self, i) {
        return self.leaves.indexOf(self.hashedElements[i]);
    }

    beforeEach(async function () {
        this.token = await TokenMock.new('1INCH Token', '1INCH');
        this.drop = await CumulativeMerkleDrop.new(this.token.address);
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

        await this.drop.contract.methods.verify(this.proofs[0], this.root, this.leaves[0]).send({ from: _ });
        await this.drop.contract.methods.verify2(this.proofs[0], this.root, this.leaves[0]).send({ from: _ });
        expect(await this.drop.verify(this.proofs[0], this.root, this.leaves[0])).to.be.true;
        expect(await this.drop.verify2(this.proofs[0], this.root, this.leaves[0])).to.be.true;
        await this.drop.claim(accounts[findSortedIndex(this, 0)], 1, this.root, this.proofs[0]);
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
        async function firstDrop (token, drop) {
            return makeDrop(token, drop, [w1, w2, w3, w4], [1, 2, 3, 4], 1 + 2 + 3 + 4);
        }

        async function secondDrop (token, drop) {
            return makeDrop(token, drop, [w1, w2, w3, w4], [3, 5, 7, 9], 2 + 3 + 4 + 5);
        }

        beforeEach(async function () {
            const { hashedElements, leaves, root, proofs } = await firstDrop(this.token, this.drop);
            this.hashedElements = hashedElements;
            this.leaves = leaves;
            this.root = root;
            this.proofs = proofs;
        });

        describe('First wallet', async function () {
            it('should success to claim 1 token, second drop and claim 2 tokens', async function () {
                await this.drop.contract.methods.verify(this.proofs[0], this.root, this.leaves[0]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(wallets[findSortedIndex(this, 0)], 1, this.root, this.proofs[0]),
                    'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '1')
                );

                const { hashedElements, leaves, root, proofs } = await secondDrop(this.token, this.drop);
                this.hashedElements = hashedElements;
                this.leaves = leaves;
                this.root = root;
                this.proofs = proofs;

                await this.drop.contract.methods.verify(this.proofs[0], this.root, this.leaves[0]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                    'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '2')
                );
            });

            it('should success to claim all 3 tokens after second drop', async function () {
                const { hashedElements, leaves, root, proofs } = await secondDrop(this.token, this.drop);
                this.hashedElements = hashedElements;
                this.leaves = leaves;
                this.root = root;
                this.proofs = proofs;

                await this.drop.contract.methods.verify(this.proofs[0], this.root, this.leaves[0]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                    'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '3')
                );
            });
        });
    });
});
