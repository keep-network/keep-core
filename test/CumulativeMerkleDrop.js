const { expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');

const { profileEVM, gasspectEVM } = require('./helpers/profileEVM');

const TokenMock = artifacts.require('TokenMock');
const CumulativeMerkleDrop = artifacts.require('CumulativeMerkleDrop');

async function makeDrop(token, drop, wallets, amounts, deposit) {
    const elements = wallets.map((w, i) => w + web3.utils.padLeft(web3.utils.toHex(amounts[i]), 64).substr(2));
    const leafs = elements.map(keccak256);
    const tree = new MerkleTree(leafs, keccak256);
    const root = tree.getHexRoot();
    const proofs = leafs.map(tree.getHexProof, tree);

    await drop.setMerkleRoot(root);
    await token.mint(drop.address, deposit);

    return { leafs, root, proofs };
}
contract('CumulativeMerkleDrop', async function ([_, w1, w2, w3, w4]) {
    beforeEach(async function () {
        this.token = await TokenMock.new("1INCH Token", "1INCH");
        this.drop = await CumulativeMerkleDrop.new(this.token.address);
    });

    it('Benchmark 30000 wallets (merkle tree height 15)', async function () {
        const wallets = Array(30000).fill().map((_, i) => w1);
        const amounts = Array(30000).fill().map((_, i) => i + 1);

        const { leafs, root, proofs } = await makeDrop(this.token, this.drop, wallets, amounts, 1000000);
        this.leafs = leafs;
        this.root = root;
        this.proofs = proofs;

        await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
        await this.drop.contract.methods.applyProof2(0, this.leafs[0], this.proofs[0]).send({ from: _ });
    });

    describe('Single drop for 4 wallets: [1, 2, 3, 4]', async function () {
        beforeEach(async function () {
            const { leafs, root, proofs } = await makeDrop(this.token, this.drop, [w1, w2, w3, w4], [1, 2, 3, 4], 10);
            this.leafs = leafs;
            this.root = root;
            this.proofs = proofs;
        });

        describe('First wallet', async function () {
            it('should fail to claim 2 tokens', async function () {
                await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                await this.drop.contract.methods.applyProof2(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(0, w1, 2, 1, this.root, this.proofs[0]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should succeed to claim 1 token', async function () {
                await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                const receipt = await expectEvent(
                    await this.drop.claim(0, w1, 1, 1, this.root, this.proofs[0]),
                    'Claimed', '0', w1, '1'
                );
                // gasspectEVM(receipt.transactionHash);
            });

            it('should fail to claim 1 token after 1 token', async function () {
                await this.drop.claim(0, w1, 1, 1, this.root, this.proofs[0]);

                await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(0, w1, 1, 1, this.root, this.proofs[0]),
                    'CMD: Drop already claimed'
                );
            });

            it('should fail to claim 2 tokens after 1 token', async function () {
                await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(0, w1, 2, 1, this.root, this.proofs[0]),
                    'CMD: Claiming amount is too high'
                );
            });
        });

        describe('Second wallet', async function () {
            it('should fail to claim 3 tokens', async function () {
                await this.drop.contract.methods.applyProof(1, this.leafs[1], this.proofs[1]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(1, w2, 3, 2, this.root, this.proofs[1]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should succeed to claim 1 + 1 tokens', async function () {
                await this.drop.contract.methods.applyProof(1, this.leafs[1], this.proofs[1]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(1, w2, 1, 2, this.root, this.proofs[1]),
                    'Claimed', '1', w2, '1'
                );

                await this.drop.contract.methods.applyProof(1, this.leafs[1], this.proofs[1]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(1, w2, 1, 2, this.root, this.proofs[1]),
                    'Claimed', '1', w2, '1'
                );
            });

            it('should fail to claim 2 tokens after 1 token', async function () {
                await this.drop.claim(1, w2, 1, 2, this.root, this.proofs[1]);

                await this.drop.contract.methods.applyProof(1, this.leafs[1], this.proofs[1]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(1, w2, 2, 2, this.root, this.proofs[1]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should succeed to claim 2 tokens', async function () {
                await this.drop.contract.methods.applyProof(1, this.leafs[1], this.proofs[1]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(1, w2, 2, 2, this.root, this.proofs[1]),
                    'Claimed', '1', w2, '2'
                );
            });

            it('should fail to claim 1 token after 2 tokens', async function () {
                await this.drop.claim(1, w2, 2, 2, this.root, this.proofs[1]),

                await this.drop.contract.methods.applyProof(1, this.leafs[1], this.proofs[1]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(1, w2, 1, 2, this.root, this.proofs[1]),
                    'CMD: Drop already claimed'
                );
            });
        });

        describe('Third wallet', async function () {
            it('should fail to claim 4 tokens', async function () {
                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(2, w3, 4, 3, this.root, this.proofs[2]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should succeed to claim 3 tokens', async function () {
                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(2, w3, 3, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '3'
                );
            });

            it('should succeed to claim 1 + 2 tokens', async function () {
                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(2, w3, 1, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '1'
                );

                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(2, w3, 2, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '2'
                );
            });

            it('should succeed to claim 2 + 1 tokens', async function () {
                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(2, w3, 2, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '2'
                );

                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(2, w3, 1, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '1'
                );
            });

            it('should fail to claim 3 tokens after 1 tokens', async function () {
                await this.drop.claim(2, w3, 1, 3, this.root, this.proofs[2]);

                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(2, w3, 3, 3, this.root, this.proofs[2]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should fail to claim 2 tokens after 2 tokens', async function () {
                await this.drop.claim(2, w3, 2, 3, this.root, this.proofs[2]);

                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(2, w3, 2, 3, this.root, this.proofs[2]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should fail to claim 1 tokens after 3 tokens', async function () {
                await this.drop.claim(2, w3, 3, 3, this.root, this.proofs[2]);

                await this.drop.contract.methods.applyProof(2, this.leafs[2], this.proofs[2]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(2, w3, 1, 3, this.root, this.proofs[2]),
                    'CMD: Drop already claimed'
                );
            });
        });

        describe('Forth wallet', async function () {
            it('should fail to claim 5 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(3, w4, 5, 4, this.root, this.proofs[3]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should succeed to claim 4 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 4, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '3'
                );
            });

            it('should succeed to claim 1 + 3 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 3, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '3'
                );
            });

            it('should succeed to claim 2 + 2 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '2'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '2'
                );
            });

            it('should succeed to claim 3 + 1 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 3, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '3'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1'
                );
            });

            it('should succeed to claim 1 + 2 + 1 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '2'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1'
                );
            });

            it('should fail to claim 1 tokens after 4 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 4, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '4'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'CMD: Drop already claimed'
                );
            });

            it('should fail to claim 2 tokens after 3 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 3, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '3'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should fail to claim 3 tokens after 2 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '2'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(3, w4, 3, 4, this.root, this.proofs[3]),
                    'CMD: Claiming amount is too high'
                );
            });

            it('should fail to claim 4 tokens after 1 tokens', async function () {
                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1'
                );

                await this.drop.contract.methods.applyProof(3, this.leafs[3], this.proofs[3]).send({ from: _ });
                await expectRevert(
                    this.drop.claim(3, w4, 4, 4, this.root, this.proofs[3]),
                    'CMD: Claiming amount is too high'
                );
            });
        });
    });

    describe('Double drop for 4 wallets: [1, 2, 3, 4] + [2, 3, 4, 5] = [3, 5, 7, 9]', async function () {
        async function firstDrop(token, drop) {
            return makeDrop(token, drop, [w1, w2, w3, w4], [1, 2, 3, 4], 1 + 2 + 3 + 4);
        }

        async function secondDrop(token, drop) {
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
                    'Claimed', '0', w1, '1'
                );

                const { leafs, root, proofs } = await secondDrop(this.token, this.drop);
                this.leafs = leafs;
                this.root = root;
                this.proofs = proofs;

                await this.drop.contract.methods.applyProof(0, this.leafs[0], this.proofs[0]).send({ from: _ });
                await expectEvent(
                    await this.drop.claim(0, w1, 2, 3, this.root, this.proofs[0]),
                    'Claimed', '0', w1, '2'
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
                    'Claimed', '0', w1, '3'
                );
            });
        });
    });
});
