const { expectEvent, expectRevert } = require('@openzeppelin/test-helpers');

function claimedEvent(account, amount) {
    return { account, amount };
}

function shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234 (errorPrefix, _, wallets, findSortedIndex, makeFirstDrop, makeSecondDrop) {
    describe('Double drop for 4 wallets: [1, 2, 3, 4] + [2, 3, 4, 5] = [3, 5, 7, 9]', async function () {
        beforeEach(async function () {
            await makeFirstDrop(this);
        });

        describe('First wallet', async function () {
            it('should success to claim 1 token, second drop and claim 2 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(wallets[findSortedIndex(this, 0)], 1, this.root, this.proofs[0]),
                    'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '1')
                );

                await makeSecondDrop(this);

                await expectEvent(
                    await this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                    'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '2')
                );
            });

            it('should success to claim all 3 tokens after second drop', async function () {
                await makeSecondDrop(this);

                await expectEvent(
                    await this.drop.claim(wallets[findSortedIndex(this, 0)], 3, this.root, this.proofs[0]),
                    'Claimed', claimedEvent(wallets[findSortedIndex(this, 0)], '3')
                );
            });
        });
    });
}

module.exports = {
    shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234,
};
