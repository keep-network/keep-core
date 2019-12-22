import { sign } from './helpers/signature';
import { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import expectThrow from './helpers/expectThrow';
import grantTokens from './helpers/grantTokens';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const RegistryKeeper = artifacts.require("./RegistryKeeper.sol");

contract('TestTokenGrantStake', function(accounts) {

  let token, registryKeeper, grantContract, stakingContract,
    id, amount,
    account_one = accounts[0],
    account_two = accounts[3],
    account_two_operator = accounts[4],
    account_two_magpie = accounts[5],
    account_two_authorizer = accounts[6];

  beforeEach(async () => {
    token = await KeepToken.new();
    registryKeeper = await RegistryKeeper.new();
    stakingContract = await TokenStaking.new(token.address, registryKeeper.address, duration.days(30));
    grantContract = await TokenGrant.new(token.address, stakingContract.address);

    let vestingDuration = duration.days(60),
    start = await latestTime(),
    cliff = duration.days(10),
    revocable = false;
    amount = web3.utils.toBN(1000000000);

    // Grant tokens
    id = await grantTokens(grantContract, token, amount, account_one, account_two, vestingDuration, start, cliff, revocable);
  });


  it("should stake granted tokens correctly", async function() {

    let delegation = Buffer.concat([
      Buffer.from(account_two_magpie.substr(2), 'hex'),
      Buffer.from(account_two_operator.substr(2), 'hex'),
      Buffer.from(account_two_authorizer.substr(2), 'hex')
    ]);

    // should throw if stake granted tokens called by anyone except grant grantee
    await expectThrow(grantContract.stake(id, stakingContract.address, amount, delegation));

    // stake granted tokens can be only called by grant grantee
    await grantContract.stake(id, stakingContract.address, amount, delegation, {from: account_two});
    let account_two_operator_stake_balance = await stakingContract.balanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.eq(amount), true, "Should stake grant amount");

    // should throw if initiate unstake called by anyone except grant grantee
    await expectThrow(grantContract.initiateUnstake(account_two_operator));

    // Initiate unstake of granted tokens by grant grantee
    await grantContract.initiateUnstake(account_two_operator, {from: account_two});

    // should not be able to finish unstake before withdrawal delay is over
    await expectThrow(grantContract.finishUnstake(account_two_operator));

    // should not be able to withdraw grant as its still locked for staking
    await expectThrow(grantContract.withdraw(id));

    // jump in time over withdrawal delay
    await increaseTimeTo(await latestTime()+duration.days(30));
    await grantContract.finishUnstake(account_two_operator);
    account_two_operator_stake_balance = await stakingContract.balanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.isZero(), true, "Stake grant amount should be 0");

    // should be able to withdraw 'withdrawable' granted amount as it's not locked for staking anymore
    await grantContract.withdraw(id);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.gte(amount.div(web3.utils.toBN(2))), true, "Should have some withdrawn grant amount");

    // Get grant available balance after withdraw
    let grant = await grantContract.getGrant(id);
    let grantAmount = grant[0];
    let grantReleased = grant[1];
    let updatedGrantBalance = grantAmount.sub(grantReleased);

    // should be able to delegate stake to the same operator after finishing unstaking
    await grantContract.stake(id, stakingContract.address, updatedGrantBalance, delegation, {from: account_two});
    account_two_operator_stake_balance = await stakingContract.balanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.eq(updatedGrantBalance), true, "Should stake grant amount");

  });
});
