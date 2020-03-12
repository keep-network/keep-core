import mineBlocks from './helpers/mineBlocks';
import { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import expectThrow from './helpers/expectThrow';
import expectThrowWithMessage from './helpers/expectThrowWithMessage'
import grantTokens from './helpers/grantTokens';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const Registry = artifacts.require("./Registry.sol");

contract.only('TokenGrant/Stake', function(accounts) {

  let token, registry, grantContract, stakingContract,
    id, amount,
    grantCreator = accounts[0],
    grantee_1 = accounts[1],
    operator_1 = accounts[2],
    beneficiary_1 = accounts[3],
    authorizer_1 = accounts[4];

  const initializationPeriod = 10;
  const undelegationPeriod = 30;

  beforeEach(async () => {
    token = await KeepToken.new();
    registry = await Registry.new();
    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod
    );
    grantContract = await TokenGrant.new(token.address, stakingContract.address);

    let vestingDuration = duration.days(60),
    start = await latestTime(),
    cliff = duration.days(10),
    revocable = false;
    amount = web3.utils.toBN(1000000000);

    // Grant tokens
    id = await grantTokens(
      grantContract,
      token, amount,
      grantCreator, grantee_1,
      vestingDuration,
      start, cliff, revocable
    );
  });


  it("should stake granted tokens correctly", async function() {
    let delegation = Buffer.concat([
      Buffer.from(beneficiary_1.substr(2), 'hex'),
      Buffer.from(operator_1.substr(2), 'hex'),
      Buffer.from(authorizer_1.substr(2), 'hex')
    ]);

    // should throw if stake granted tokens called by anyone except grant grantee
    await expectThrowWithMessage(
      grantContract.stake(id, stakingContract.address, amount, delegation),
      "Only grantee of the grant can stake it."
    );

    // stake granted tokens can be only called by grant grantee
    await grantContract.stake(
      id,
      stakingContract.address,
      amount,
      delegation,
      {from: grantee_1}
    );
    let operator_1_stake_balance = await stakingContract.balanceOf.call(operator_1);
    assert.equal(
      operator_1_stake_balance.eq(amount),
      true,
      "Should stake grant amount"
    );

    // should throw if undelegate called by anyone except grant grantee
    await expectThrowWithMessage(
      grantContract.undelegate(operator_1),
      "Only operator or grantee can undelegate."
    );

    // Undelegate granted tokens by grant grantee
    await grantContract.undelegate(operator_1, {from: grantee_1});

    // should not be able to recover stake before undelegation period is over
    await expectThrowWithMessage(
      grantContract.recoverStake(operator_1),
      "Can not recover stake before undelegation period is over."
    );

    // should not be able to withdraw grant as its still locked for staking
    await expectThrow(grantContract.withdraw(id));

    // jump in time over undelegation period
    await mineBlocks(undelegationPeriod);
    await grantContract.recoverStake(operator_1);
    operator_1_stake_balance = await stakingContract.balanceOf.call(operator_1);
    assert.equal(
      operator_1_stake_balance.isZero(),
      true,
      "Stake grant amount should be 0"
    );

    // jump in time to allow tokens to vest
    await increaseTimeTo(await latestTime()+duration.days(30));

    // should be able to withdraw 'withdrawable' granted amount as it's not locked for staking anymore
    await grantContract.withdraw(id);
    let grantee_1_ending_balance = await token.balanceOf.call(grantee_1);
    assert.equal(
      grantee_1_ending_balance.gte(amount.div(web3.utils.toBN(2))),
      true,
      "Should have some withdrawn grant amount"
    );

    // Get grant available balance after withdraw
    let grant = await grantContract.getGrant(id);
    let grantAmount = grant[0];
    let grantReleased = grant[1];
    let updatedGrantBalance = grantAmount.sub(grantReleased);

    // should be able to delegate stake to the same operator after finishing unstaking
    await grantContract.stake(
      id,
      stakingContract.address,
      updatedGrantBalance,
      delegation,
      {from: grantee_1}
    );
    operator_1_stake_balance = await stakingContract.balanceOf.call(operator_1);
    assert.equal(
      operator_1_stake_balance.eq(updatedGrantBalance),
      true,
      "Should stake grant amount"
    );
  });
});
