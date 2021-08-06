const { expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');

const { profileEVM, gasspectEVM } = require('./helpers/profileEVM');

const TokenMock = artifacts.require('TokenMock');
const CumulativeMerkleDrop = artifacts.require('CumulativeMerkleDrop');

contract('CumulativeMerkleDrop', async function ([_, w1, w2, w3, w4]) {
    beforeEach(async function () {
        this.token = await TokenMock.new("1INCH Token", "1INCH");
        this.drop = await CumulativeMerkleDrop.new(this.token.address);
    });

    it('should be ok', async function () {
        const elements = [
            w1 + '0000000000000000000000000000000000000000000000000000000000000001',
            w2 + '0000000000000000000000000000000000000000000000000000000000000002',
            w3 + '0000000000000000000000000000000000000000000000000000000000000003',
            w4 + '0000000000000000000000000000000000000000000000000000000000000004',
        ];
        const leafs = elements.map(keccak256);
        const merkleTree = new MerkleTree(leafs, keccak256);
        const merkleRoot = merkleTree.getHexRoot();
        const proofs = leafs.map(merkleTree.getHexProof, merkleTree);

        await this.token.mint(this.drop.address, 10);
        await this.drop.setMerkleRoot(merkleRoot);

        // First element (balance: 1)
        {
            // 2 tokens failed
            await this.drop.contract.methods.applyProof(0, leafs[0], proofs[0]).send({ from: _ });
            await expectRevert(
                this.drop.claim(0, w1, 2, 1, merkleRoot, proofs[0]),
                'CMD: Claiming amount is too high'
            );

            // 1 token success
            await this.drop.contract.methods.applyProof(0, leafs[0], proofs[0]).send({ from: _ });
            const receipt = await expectEvent(
                await this.drop.claim(0, w1, 1, 1, merkleRoot, proofs[0]),
                'Claimed', '0', w1, '1'
            );
            gasspectEVM(receipt.transactionHash);

            // 1 tokens failed
            await this.drop.contract.methods.applyProof(0, leafs[0], proofs[0]).send({ from: _ });
            await expectRevert(
                this.drop.claim(0, w1, 1, 1, merkleRoot, proofs[0]),
                'CMD: Drop already claimed'
            );
        }

        // Second element (balance: 2)
        {
            // 1st token success
            await this.drop.contract.methods.applyProof(1, leafs[1], proofs[1]).send({ from: _ });
            await expectEvent(
                await this.drop.claim(1, w2, 1, 2, merkleRoot, proofs[1]),
                'Claimed', '1', w2, '1'
            );

            // 3rd and 4th token failed
            await this.drop.contract.methods.applyProof(1, leafs[1], proofs[1]).send({ from: _ });
            await expectRevert(
                this.drop.claim(1, w2, 2, 2, merkleRoot, proofs[1]),
                'CMD: Claiming amount is too high'
            );

            // 2nd token success
            await this.drop.contract.methods.applyProof(1, leafs[1], proofs[1]).send({ from: _ });
            await expectEvent(
                await this.drop.claim(1, w2, 1, 2, merkleRoot, proofs[1]),
                'Claimed', '1', w2, '1'
            );

            // 3rd token failed
            await this.drop.contract.methods.applyProof(1, leafs[1], proofs[1]).send({ from: _ });
            await expectRevert(
                this.drop.claim(1, w2, 1, 2, merkleRoot, proofs[1]),
                'CMD: Drop already claimed'
            );
        }
    });
});
