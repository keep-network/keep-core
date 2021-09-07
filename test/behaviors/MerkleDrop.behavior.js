const { expectEvent, expectRevert } = require('@openzeppelin/test-helpers');

function claimedEvent (account, amount) {
    return { account, amount };
}

function shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234 (errorPrefix, wallets, findSortedIndex, is128version = false) {
    describe('Single drop for 4 wallets: [1, 2, 3, 4]', async function () {
        describe('First wallet', async function () {
            it('should succeed to claim 1 token', async function () {
                await expectEvent(
                    (
                        is128version
                            ? await this.drop.claim(this.salts[0], wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)])
                            : await this.drop.claim(wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)])
                    ),
                    'Claimed', claimedEvent(wallets[0], '1'),
                );
            });

            it('should fail to claim second time', async function () {
                if (is128version) {
                    await this.drop.claim(this.salts[0], wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)]);
                } else {
                    await this.drop.claim(wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)]);
                }

                await expectRevert(
                    (
                        is128version
                            ? this.drop.claim(this.salts[0], wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)])
                            : this.drop.claim(wallets[0], 1, this.root, this.proofs[findSortedIndex(this, 0)])
                    ),
                    `${errorPrefix}: Nothing to claim`,
                );
            });
        });

        describe('Second wallet', async function () {
            it('should succeed to claim', async function () {
                await expectEvent(
                    (
                        is128version
                            ? await this.drop.claim(this.salts[1], wallets[1], 2, this.root, this.proofs[findSortedIndex(this, 1)])
                            : await this.drop.claim(wallets[1], 2, this.root, this.proofs[findSortedIndex(this, 1)])
                    ),
                    'Claimed', claimedEvent(wallets[1], '2'),
                );
            });

            it('should fail to claim second time', async function () {
                if (is128version) {
                    await this.drop.claim(this.salts[1], wallets[1], 2, this.root, this.proofs[findSortedIndex(this, 1)]);
                } else {
                    await this.drop.claim(wallets[1], 2, this.root, this.proofs[findSortedIndex(this, 1)]);
                }

                await expectRevert(
                    (
                        is128version
                            ? this.drop.claim(this.salts[1], wallets[1], 2, this.root, this.proofs[findSortedIndex(this, 1)])
                            : this.drop.claim(wallets[1], 2, this.root, this.proofs[findSortedIndex(this, 1)])
                    ),
                    `${errorPrefix}: Nothing to claim`,
                );
            });
        });

        describe('Third wallet', async function () {
            it('should succeed to claim', async function () {
                await expectEvent(
                    (
                        is128version
                            ? await this.drop.claim(this.salts[2], wallets[2], 3, this.root, this.proofs[findSortedIndex(this, 2)])
                            : await this.drop.claim(wallets[2], 3, this.root, this.proofs[findSortedIndex(this, 2)])
                    ),
                    'Claimed', claimedEvent(wallets[2], '3'),
                );
            });

            it('should fail to claim second time', async function () {
                if (is128version) {
                    await this.drop.claim(this.salts[2], wallets[2], 3, this.root, this.proofs[findSortedIndex(this, 2)]);
                } else {
                    await this.drop.claim(wallets[2], 3, this.root, this.proofs[findSortedIndex(this, 2)]);
                }

                await expectRevert(
                    (
                        is128version
                            ? this.drop.claim(this.salts[2], wallets[2], 3, this.root, this.proofs[findSortedIndex(this, 2)])
                            : this.drop.claim(wallets[2], 3, this.root, this.proofs[findSortedIndex(this, 2)])
                    ),
                    `${errorPrefix}: Nothing to claim`,
                );
            });
        });

        describe('Forth wallet', async function () {
            it('should succeed to claim', async function () {
                await expectEvent(
                    (
                        is128version
                            ? await this.drop.claim(this.salts[3], wallets[3], 4, this.root, this.proofs[findSortedIndex(this, 3)])
                            : await this.drop.claim(wallets[3], 4, this.root, this.proofs[findSortedIndex(this, 3)])
                    ),
                    'Claimed', claimedEvent(wallets[3], '4'),
                );
            });

            it('should fail to claim 1 tokens after 4 tokens', async function () {
                await expectEvent(
                    (
                        is128version
                            ? await this.drop.claim(this.salts[3], wallets[3], 4, this.root, this.proofs[findSortedIndex(this, 3)])
                            : await this.drop.claim(wallets[3], 4, this.root, this.proofs[findSortedIndex(this, 3)])
                    ),
                    'Claimed', claimedEvent(wallets[3], '4'),
                );

                await expectRevert(
                    (
                        is128version
                            ? this.drop.claim(this.salts[3], wallets[3], 4, this.root, this.proofs[findSortedIndex(this, 3)])
                            : this.drop.claim(wallets[3], 4, this.root, this.proofs[findSortedIndex(this, 3)])
                    ),
                    `${errorPrefix}: Nothing to claim`,
                );
            });
        });
    });
}

module.exports = {
    shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234,
};
