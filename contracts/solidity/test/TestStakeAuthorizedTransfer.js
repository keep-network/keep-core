import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');

contract('TestStakeAuthorizedTransfer', function(accounts) {

  let token, stakingProxy, stakingContract,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})
    await token.approveAndCall(stakingContract.address, 10000000, "", {from: account_one});
  });

  it("should be able to move staked tokens if authorized", async function() {

    // Fail to transfer if account not authorized.
    await exceptThrow(stakingContract.authorizedTransfer(account_one, 100, {from: account_two}));

    // Authorize an account that can do staked token transfers.
    await stakingContract.authorize(account_two);

    // Move staked tokens.
    await stakingContract.authorizedTransfer(account_one, 100, {from: account_two})

    assert.equal(await stakingContract.stakeBalanceOf(account_two), 100, "Staked tokens should be transferred.");

  });
});
