const { expectEvent, expectRevert } = require('@openzeppelin/test-helpers');

function claimedEvent (account, amount) {
    return { account, amount };
}

function shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234 (errorPrefix, _, wallets, findSortedIndex, makeFirstDrop, makeSecondDrop) {
    describe('First wallet checks', async function () {
        beforeEach(async function () {
            await makeFirstDrop(this);
        });

        it('should success to claim 1 token, second drop and claim 2 tokens', async function () {
            await expectEvent(
                await this.drop.claim(wallets[findSortedIndex(this, 0)], 1, this.root, this.proofs[0]),
                'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '1'),
            );

            await makeSecondDrop(this);

            await expectEvent(
                await this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '2'),
            );

            await expectRevert(
                this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                `${errorPrefix}: Nothing to claim`,
            );
        });

        it('should success to claim 1 token, second drop and claim 2 tokens twice', async function () {
            await expectEvent(
                await this.drop.claim(wallets[findSortedIndex(this, 0)], 1, this.root, this.proofs[0]),
                'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '1'),
            );

            await makeSecondDrop(this);

            await expectEvent(
                await this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '2'),
            );

            await expectRevert(
                this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                `${errorPrefix}: Nothing to claim`,
            );
        });

        it('should success to claim all 3 tokens after second drop', async function () {
            await makeSecondDrop(this);

            await expectEvent(
                await this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '3'),
            );
        });

        it('should fail to claim after succelfful claim of all 3 tokens after second drop', async function () {
            await makeSecondDrop(this);

            await expectEvent(
                await this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '3'),
            );

            await expectRevert(
                this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                `${errorPrefix}: Nothing to claim`,
            );
        });
    });
}

module.exports = {
    shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234,
};
