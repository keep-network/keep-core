const { BN, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');

const {
    shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234
} = require('./MerkleDrop.behavior');

const TokenMock = artifacts.require('TokenMock');
const CumulativeMerkleDrop128 = artifacts.require('CumulativeMerkleDrop128');

function keccak128 (input) {
    return keccak256(input).slice(16, 32);
}

async function makeDrop (token, drop, wallets, amounts, deposit) {
    const elements = wallets.map((w, i) => w + web3.utils.padLeft(web3.utils.toHex(amounts[i]), 64).substr(2));
    const leafs = elements.map(keccak128);
    const tree = new MerkleTree(leafs, keccak128);
    const root = tree.getHexRoot();
    const proofs = leafs
        .map(tree.getHexProof, tree)
        .map(proof => '0x' + proof.map(p => p.substr(2)).join(''));

    await drop.setMerkleRoot(root);
    await token.mint(drop.address, deposit);

    return { leafs, root, proofs };
}

contract('CumulativeMerkleDrop128', async function ([_, w1, w2, w3, w4]) {
    beforeEach(async function () {
        this.token = await TokenMock.new('1INCH Token', '1INCH');
        this.drop = await CumulativeMerkleDrop128.new(this.token.address);
    });

    it('Benchmark 30000 wallets (merkle tree height 15)', async function () {
        const wallets = Array(30000).fill().map((_, i) => '0x' + (new BN(w1)).addn(i).toString('hex'));
        const amounts = Array(30000).fill().map((_, i) => i + 1);

        const { leafs, root, proofs } = await makeDrop(this.token, this.drop, wallets, amounts, 1000000);
        this.leafs = leafs;
        this.root = root;
        this.proofs = proofs;

        console.log(this.leafs[0], this.proofs[0]);
        await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
        await this.drop.contract.methods.applyProof2(0, this.leafs[0], this.proofs[0]).send({ from: _ });
        expect(await this.drop.applyProof(0, this.leafs[0], this.proofs[0])).to.be.equal(this.root);
        expect(await this.drop.applyProof2(0, this.leafs[0], this.proofs[0])).to.be.equal(this.root);
    });

    describe('Single drop for 4 wallets: [1, 2, 3, 4]', async function () {
        beforeEach(async function () {
            const { leafs, root, proofs } = await makeDrop(this.token, this.drop, [w1, w2, w3, w4], [1, 2, 3, 4], 10);
            this.leafs = leafs;
            this.root = root;
            this.proofs = proofs;
        });

        shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234('CMD', [w1, w2, w3, w4]);
    });

    describe('Double drop for 4 wallets: [1, 2, 3, 4] + [2, 3, 4, 5] = [3, 5, 7, 9]', async function () {
        async function firstDrop (token, drop) {
            return makeDrop(token, drop, [w1, w2, w3, w4], [1, 2, 3, 4], 1 + 2 + 3 + 4);
        }

        async function secondDrop (token, drop) {
            return await makeDrop(token, drop, [w1, w2, w3, w4], [3, 5, 7, 9], 2 + 3 + 4 + 5);
        }

        beforeEach(async function () {
            const { leafs, root, proofs } = await firstDrop(this.token, this.drop);
            this.leafs = leafs;
            this.root = root;
            this.proofs = proofs;
        });

        describe('First wallet', async function () {
            it('should success to claim 1 token, second drop and claim 2 tokens', async function () {
                await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(0, w1, 1, 1, this.root, this.proofs[0]),
                    'Claimed', '0', w1, '1',
                );

                const { leafs, root, proofs } = await secondDrop(this.token, this.drop);
                this.leafs = leafs;
                this.root = root;
                this.proofs = proofs;

                await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(0, w1, 2, 3, this.root, this.proofs[0]),
                    'Claimed', '0', w1, '2',
                );
            });

            it('should success to claim 3 tokens after second drop', async function () {
                const { leafs, root, proofs } = await secondDrop(this.token, this.drop);
                this.leafs = leafs;
                this.root = root;
                this.proofs = proofs;

                await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(0, w1, 3, 3, this.root, this.proofs[0]),
                    'Claimed', '0', w1, '3',
                );
            });
        });
    });
});
