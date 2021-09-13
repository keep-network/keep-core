const { expectEvent, expectRevert } = require('@openzeppelin/test-helpers');

function claimedEvent (account, amount) {
    return { account, amount };
}

function shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234 (errorPrefix, _, wallets, findSortedIndex, makeFirstDrop, makeSecondDrop, is128version = false) {
    describe('First wallet checks', async function () {
        beforeEach(async function () {
            await makeFirstDrop(this);
        });

        it('should success to claim 1 token, second drop and claim 2 tokens', async function () {
            await expectEvent(
                (
                    is128version
                        ? await this.drop.claim(this.salts[0], wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)])
                        : await this.drop.claim(wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                'Claimed', claimedEvent(wallets[0], '1'),
            );

            await makeSecondDrop(this);

            await expectEvent(
                (
                    is128version
                        ? await this.drop.claim(this.salts[0], wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                        : await this.drop.claim(wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                'Claimed', claimedEvent(wallets[0], '2'),
            );

            await expectRevert(
                (
                    is128version
                        ? this.drop.claim(this.salts[0], wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                        : this.drop.claim(wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                `${errorPrefix}: Nothing to claim`,
            );
        });

        it('should success to claim 1 token, second drop and claim 2 tokens twice', async function () {
            await expectEvent(
                (
                    is128version
                        ? await this.drop.claim(this.salts[0], wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)])
                        : await this.drop.claim(wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                'Claimed', claimedEvent(wallets[0], '1'),
            );

            await makeSecondDrop(this);

            await expectEvent(
                (
                    is128version
                        ? await this.drop.claim(this.salts[0], wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                        : await this.drop.claim(wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                'Claimed', claimedEvent(wallets[0], '2'),
            );

            await expectRevert(
                (
                    is128version
                        ? this.drop.claim(this.salts[0], wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                        : this.drop.claim(wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                `${errorPrefix}: Nothing to claim`,
            );
        });

        it('should success to claim all 3 tokens after second drop', async function () {
            await makeSecondDrop(this);

            await expectEvent(
                (
                    is128version
                        ? await this.drop.claim(this.salts[0], wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                        : await this.drop.claim(wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                'Claimed', claimedEvent(wallets[0], '3'),
            );
        });

        it('should fail to claim after succelfful claim of all 3 tokens after second drop', async function () {
            await makeSecondDrop(this);

            await expectEvent(
                (
                    is128version
                        ? await this.drop.claim(this.salts[0], wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                        : await this.drop.claim(wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                'Claimed', claimedEvent(wallets[0], '3'),
            );

            await expectRevert(
                (
                    is128version
                        ? this.drop.claim(this.salts[0], wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                        : this.drop.claim(wallets[0], 3, this.root, this.proofs[findSortedIndex(this, 0)])
                ),
                `${errorPrefix}: Nothing to claim`,
            );
        });
    });
}

module.exports = {
    shouldBehaveLikeCumulativeMerkleDropFor4WalletsWithBalances1234,
};
