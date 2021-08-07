const { expectEvent, expectRevert } = require('@openzeppelin/test-helpers');

function shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234 (errorPrefix, [w1, w2, w3, w4]) {
    describe('Single drop for 4 wallets: [1, 2, 3, 4]', async function () {
        describe('First wallet', async function () {
            it('should fail to claim 2 tokens', async function () {
                await expectRevert(
                    this.drop.claim(0, w1, 2, 1, this.root, this.proofs[0]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should succeed to claim 1 token', async function () {
                await expectEvent(
                    await this.drop.claim(0, w1, 1, 1, this.root, this.proofs[0]),
                    'Claimed', '0', w1, '1',
                );
            });

            it('should fail to claim 1 token after 1 token', async function () {
                await this.drop.claim(0, w1, 1, 1, this.root, this.proofs[0]);

                await expectRevert(
                    this.drop.claim(0, w1, 1, 1, this.root, this.proofs[0]),
                    `${errorPrefix}: Drop already claimed`,
                );
            });

            it('should fail to claim 2 tokens after 1 token', async function () {
                await expectRevert(
                    this.drop.claim(0, w1, 2, 1, this.root, this.proofs[0]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });
        });

        describe('Second wallet', async function () {
            it('should fail to claim 3 tokens', async function () {
                await expectRevert(
                    this.drop.claim(1, w2, 3, 2, this.root, this.proofs[1]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should succeed to claim 1 + 1 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(1, w2, 1, 2, this.root, this.proofs[1]),
                    'Claimed', '1', w2, '1',
                );

                await expectEvent(
                    await this.drop.claim(1, w2, 1, 2, this.root, this.proofs[1]),
                    'Claimed', '1', w2, '1',
                );
            });

            it('should fail to claim 2 tokens after 1 token', async function () {
                await this.drop.claim(1, w2, 1, 2, this.root, this.proofs[1]);

                await expectRevert(
                    this.drop.claim(1, w2, 2, 2, this.root, this.proofs[1]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should succeed to claim 2 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(1, w2, 2, 2, this.root, this.proofs[1]),
                    'Claimed', '1', w2, '2',
                );
            });

            it('should fail to claim 1 token after 2 tokens', async function () {
                await this.drop.claim(1, w2, 2, 2, this.root, this.proofs[1]);

                await expectRevert(
                    this.drop.claim(1, w2, 1, 2, this.root, this.proofs[1]),
                    `${errorPrefix}: Drop already claimed`,
                );
            });
        });

        describe('Third wallet', async function () {
            it('should fail to claim 4 tokens', async function () {
                await expectRevert(
                    this.drop.claim(2, w3, 4, 3, this.root, this.proofs[2]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should succeed to claim 3 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(2, w3, 3, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '3',
                );
            });

            it('should succeed to claim 1 + 2 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(2, w3, 1, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '1',
                );

                await expectEvent(
                    await this.drop.claim(2, w3, 2, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '2',
                );
            });

            it('should succeed to claim 2 + 1 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(2, w3, 2, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '2',
                );

                await expectEvent(
                    await this.drop.claim(2, w3, 1, 3, this.root, this.proofs[2]),
                    'Claimed', '2', w3, '1',
                );
            });

            it('should fail to claim 3 tokens after 1 tokens', async function () {
                await this.drop.claim(2, w3, 1, 3, this.root, this.proofs[2]);

                await expectRevert(
                    this.drop.claim(2, w3, 3, 3, this.root, this.proofs[2]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should fail to claim 2 tokens after 2 tokens', async function () {
                await this.drop.claim(2, w3, 2, 3, this.root, this.proofs[2]);

                await expectRevert(
                    this.drop.claim(2, w3, 2, 3, this.root, this.proofs[2]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should fail to claim 1 tokens after 3 tokens', async function () {
                await this.drop.claim(2, w3, 3, 3, this.root, this.proofs[2]);

                await expectRevert(
                    this.drop.claim(2, w3, 1, 3, this.root, this.proofs[2]),
                    `${errorPrefix}: Drop already claimed`,
                );
            });
        });

        describe('Forth wallet', async function () {
            it('should fail to claim 5 tokens', async function () {
                await expectRevert(
                    this.drop.claim(3, w4, 5, 4, this.root, this.proofs[3]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should succeed to claim 4 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 4, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '3',
                );
            });

            it('should succeed to claim 1 + 3 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1',
                );

                await expectEvent(
                    await this.drop.claim(3, w4, 3, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '3',
                );
            });

            it('should succeed to claim 2 + 2 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '2',
                );

                await expectEvent(
                    await this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '2',
                );
            });

            it('should succeed to claim 3 + 1 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 3, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '3',
                );

                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1',
                );
            });

            it('should succeed to claim 1 + 2 + 1 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1',
                );

                await expectEvent(
                    await this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '2',
                );

                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1',
                );
            });

            it('should fail to claim 1 tokens after 4 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 4, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '4',
                );

                await expectRevert(
                    this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    `${errorPrefix}: Drop already claimed`,
                );
            });

            it('should fail to claim 2 tokens after 3 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 3, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '3',
                );

                await expectRevert(
                    this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should fail to claim 3 tokens after 2 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 2, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '2',
                );

                await expectRevert(
                    this.drop.claim(3, w4, 3, 4, this.root, this.proofs[3]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });

            it('should fail to claim 4 tokens after 1 tokens', async function () {
                await expectEvent(
                    await this.drop.claim(3, w4, 1, 4, this.root, this.proofs[3]),
                    'Claimed', '3', w4, '1',
                );

                await expectRevert(
                    this.drop.claim(3, w4, 4, 4, this.root, this.proofs[3]),
                    `${errorPrefix}: Claiming amount is too high`,
                );
            });
        });
    });
}

module.exports = {
    shouldBehaveLikeMerkleDropFor4WalletsWithBalances1234,
};
