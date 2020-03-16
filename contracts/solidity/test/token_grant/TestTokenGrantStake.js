import mineBlocks from '../helpers/mineBlocks';
import { duration, increaseTimeTo } from '../helpers/increaseTime';
import latestTime from '../helpers/latestTime';
import expectThrow from '../helpers/expectThrow';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage'
import grantTokens from '../helpers/grantTokens';
import { createSnapshot, restoreSnapshot } from '../helpers/snapshot'

function formatAmount(amount) {
  return web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(18)))
}

export const wait = (ms) => {
  return new Promise((resolve) => {
    return setTimeout(resolve, ms)
  })
}

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const Registry = artifacts.require("./Registry.sol");

contract('TokenGrant/Stake', function(accounts) {

  let tokenContract, registryContract, grantContract, stakingContract;

  const tokenOwner = accounts[0],
    grantee = accounts[1],
    operatorOne = accounts[2],
    operatorTwo = accounts[3],
    magpie = accounts[4],
    authorizer = accounts[5];

  let grantId;
  let grantStart;

  const grantAmount = web3.utils.toBN(1000000000),
    grantVestingDuration = duration.days(60),
    grantCliff = duration.days(10),
    grantRevocable = false;

  let shortGrantId;
  const shortGrantDuration = 60;
  const shortGrantCliff = 1;
    
  const initializationPeriod = 10;
  const undelegationPeriod = 30;

  before(async () => {
    tokenContract = await KeepToken.new();
    registryContract = await Registry.new();
    stakingContract = await TokenStaking.new(
      tokenContract.address, 
      registryContract.address, 
      initializationPeriod, 
      undelegationPeriod
    );
    grantContract = await TokenGrant.new(
      tokenContract.address, 
      stakingContract.address
    );

    grantStart = await latestTime();

    // Grant tokens
    grantId = await grantTokens(
      grantContract, 
      tokenContract, 
      grantAmount, 
      tokenOwner, 
      grantee, 
      grantVestingDuration, 
      grantStart, 
      grantCliff, 
      grantRevocable
    );

    shortGrantId = await grantTokens(
      grantContract, 
      tokenContract,
      grantAmount,
      tokenOwner, 
      grantee, 
      shortGrantDuration, 
      grantStart, 
      shortGrantCliff, 
      false,
    );
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  async function delegate(grantee, operator, amount, grantId = grantId) {
    let delegation = Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operator.substr(2), 'hex'),
      Buffer.from(authorizer.substr(2), 'hex')
    ]);

    return grantContract.stake(
      grantId, 
      stakingContract.address, 
      amount, 
      delegation, 
      {from: grantee}
    );
  }

  it("should update balances when delegating", async () => {
    let amountToDelegate = web3.utils.toBN(20000);
    let remaining = grantAmount.sub(amountToDelegate)

    await delegate(grantee, operatorOne, amountToDelegate);

    let availableForStaking = await grantContract.availableToStake.call(grantId)
    let operatorBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(availableForStaking).to.eq.BN(
      remaining,
      "All granted tokens delegated, should be nothing more available"
    )
    expect(operatorBalance).to.eq.BN(
      amountToDelegate, 
      "Staking amount should be added to the operator balance"
    );
  })

  it("should allow to delegate, undelegate, and recover grant", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await mineBlocks(initializationPeriod);
    await grantContract.undelegate(operatorOne, {from: grantee});
    await mineBlocks(undelegationPeriod);
    await grantContract.recoverStake(operatorOne);

    let availableForStaking = await grantContract.availableToStake.call(grantId)
    let operatorBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(availableForStaking).to.eq.BN(
      grantAmount,
      "All granted tokens should be again available for staking"
    )
    expect(operatorBalance).to.eq.BN(
      0, 
      "Staking amount should be removed from operator balance"
    );
  })

  it("should allow to cancel delegation right away", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await grantContract.cancelStake(operatorOne, {from: grantee});

    let availableForStaking = await grantContract.availableToStake.call(grantId)
    let operatorBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(availableForStaking).to.eq.BN(
      grantAmount,
      "All granted tokens should be again available for staking"
    )
    expect(operatorBalance).to.eq.BN(
      0, 
      "Staking amount should be removed from operator balance"
    );
  })

  it("should allow to cancel delegation just before initialization period is over", async () => {
    await delegate(grantee, operatorOne, grantAmount);
    
    await mineBlocks(initializationPeriod - 1);

    await grantContract.cancelStake(operatorOne, {from: grantee});

    let availableForStaking = await grantContract.availableToStake.call(grantId)
    let operatorBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(availableForStaking).to.eq.BN(
      grantAmount,
      "All granted tokens should be again available for staking"
    )
    expect(operatorBalance).to.eq.BN(
      0, 
      "Staking amount should be removed from operator balance"
    );
  })

  it("should not allow to cancel delegation after initialization period is over", async () => {
    await delegate(grantee, operatorOne, grantAmount);
    
    await mineBlocks(initializationPeriod);

    await expectThrowWithMessage(
      grantContract.cancelStake(operatorOne, {from: grantee}),
      "Initialization period is over"
    );
  })

  it("should not allow to recover stake before undelegation period is over", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await mineBlocks(initializationPeriod);
    await grantContract.undelegate(operatorOne, {from: grantee});

    await mineBlocks(undelegationPeriod - 1);

    await expectThrowWithMessage(
      stakingContract.recoverStake(operatorOne),
      "Can not recover stake before undelegation period is over"
    )
  })

  it("should not allow to delegate to the same operator twice", async () => {
    let amountToDelegate = web3.utils.toBN(20000);
    await delegate(grantee, operatorOne, amountToDelegate);

    await expectThrowWithMessage(
      delegate(grantee, operatorOne, amountToDelegate),
      "Operator address is already in use"
    )
  })

  it("should not allow to delegate to the same operator even after recovering stake", async () => {
    await delegate(grantee, operatorOne, grantAmount);
    await mineBlocks(initializationPeriod);
    await grantContract.undelegate(operatorOne, {from: grantee});
    await mineBlocks(undelegationPeriod);
    await grantContract.recoverStake(operatorOne, {from: grantee});

    await expectThrowWithMessage(
      delegate(grantee, operatorOne, grantAmount),
      "Operator address is already in use."
    )
  })

  it("should allow to delegate to two different operators", async () => {
    let amountToDelegate = web3.utils.toBN(20000);

    await delegate(grantee, operatorOne, amountToDelegate);
    await delegate(grantee, operatorTwo, amountToDelegate);

    let availableForStaking = await grantContract.availableToStake.call(grantId)
    let operatorOneBalance = await stakingContract.balanceOf.call(operatorOne);
    let operatorTwoBalance = await stakingContract.balanceOf.call(operatorTwo);

    expect(availableForStaking).to.eq.BN(
      grantAmount.sub(amountToDelegate).sub(amountToDelegate),
      "All granted tokens delegated, should be nothing more available"
    )
    expect(operatorOneBalance).to.eq.BN(
      amountToDelegate, 
      "Staking amount should be added to the operator balance"
    );  
    expect(operatorTwoBalance).to.eq.BN(
      amountToDelegate, 
      "Staking amount should be added to the operator balance"
    );  
  })

  it("should not allow anyone but grantee to stake", async () => {
    await expectThrowWithMessage(
      delegate(operatorOne, operatorOne, grantAmount),
      "Only grantee of the grant can stake it."
    );
  })

  it("should let operator cancel delegation", async () => {
    await delegate(grantee, operatorOne, grantAmount);
  
    await grantContract.cancelStake(operatorOne, {from: operatorOne})
    // ok, no exception
  })

  it("should not allow third party to cancel delegation", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await expectThrowWithMessage(
      grantContract.cancelStake(operatorOne, {from: operatorTwo}),
      "Only operator or grantee can cancel the delegation."
    );
  })

  it("should let operator undelegate", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await mineBlocks(initializationPeriod);
    await grantContract.undelegate(operatorOne, {from: operatorOne})
    // ok, no exceptions
  })

  it("should not allow third party to undelegate", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await mineBlocks(initializationPeriod);
    await expectThrowWithMessage(
      grantContract.undelegate(operatorOne, {from: operatorTwo}),
      "Only operator or grantee can undelegate"
    )
  })

  it("should recover tokens recovered outside the grant contract", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await mineBlocks(initializationPeriod);
    await grantContract.undelegate(operatorOne, {from: grantee});
    await mineBlocks(undelegationPeriod);
    await stakingContract.recoverStake(operatorOne);
    let availablePre = await grantContract.availableToStake(grantId);

    expect(availablePre).to.eq.BN(
      0,
      "Staked tokens should be displaced"
    );

    await grantContract.recoverStake(operatorOne);
    let availablePost = await grantContract.availableToStake(grantId);

    expect(availablePost).to.eq.BN(
      grantAmount,
      "Staked tokens should be recovered safely"
    );
  })

  it("should allow to wihtdraw some tokens", async () => {
    await increaseTimeTo(grantStart + shortGrantDuration - 30)

    const withdrawable = await grantContract.withdrawable(shortGrantId)
    const granteeTokenGrantBalance = await grantContract.balanceOf(grantee)
    await grantContract.withdraw(shortGrantId)
    const granteeTokenGrantBalancePost = await grantContract.balanceOf(grantee)

    const granteeTokenBalance = await tokenContract.balanceOf(grantee)
    const gratDetails = await grantContract.getGrant(shortGrantId)
    
    expect(withdrawable).to.be.gt.BN(
      0,
      "Should allow to withdraw more than 0"
    )
    expect(withdrawable).to.be.lt.BN(
      grantAmount,
      `Should allow to withdraw less than ${grantAmount.toString()}`
    )
    expect(granteeTokenBalance).to.eq.BN(
      gratDetails.withdrawn,
      "Grantee KEEP token balance should be euqlas to the grant withdrawn amount"
    )
    expect(granteeTokenGrantBalance.sub(granteeTokenGrantBalancePost)).to.eq.BN(
      gratDetails.withdrawn,
      "Grantee token grant balance should be updated"
    )
  })

  it("should allow to wihtdraw the whole grant amount ", async () => {
    await increaseTimeTo(grantStart + shortGrantDuration)

    const withdrawable = await grantContract.withdrawable(shortGrantId)
    const granteeTokenGrantBalance = await grantContract.balanceOf(grantee)
    await grantContract.withdraw(shortGrantId)
    const withdrawablePost = await grantContract.withdrawable(shortGrantId)
    const granteeTokenGrantBalancePost = await grantContract.balanceOf(grantee)

    const granteeTokenBalance = await tokenContract.balanceOf(grantee)
    const gratDetails = await grantContract.getGrant(shortGrantId)

    expect(withdrawable).to.eq.BN(
      grantAmount,
      "The withdrawable amount should be equals to the whole grant amount"
    )
    expect(granteeTokenBalance).to.eq.BN(
      grantAmount,
      "Grantee KEEP token balance should be euqlas to the grant amount"
    )
    expect(withdrawablePost).to.eq.BN(
      0,
      "The withdrawable amount should be equals to 0, when the whole grant amount has been withdrawn"
    )
    expect(granteeTokenGrantBalance.sub(grantAmount)).to.eq.BN(
      granteeTokenGrantBalancePost,
      "Grantee token grant balance should be updated"
    )
    expect(gratDetails.withdrawn).to.eq.BN(
      grantAmount,
      "The grant withdrawan amount should be updated"
    )
  })

  it.only("should not allow to withdraw tokens", async () => {
    await increaseTimeTo(grantStart + shortGrantDuration)
    const withdrawable = await grantContract.withdrawable(shortGrantId)
    await delegate(grantee, operatorOne, grantAmount, shortGrantId)
    const withdrawableAfterStake = await grantContract.withdrawable(shortGrantId)

    await expectThrowWithMessage(
      grantContract.withdraw(shortGrantId),
      "Grant available to withdraw amount should be greater than zero."
    )
    expect(withdrawable).to.eq.BN(
      grantAmount,
      "The withdrawable amount should be equals to the whole grant amount"
    )
    expect(withdrawableAfterStake).to.eq.BN(
      0,
      "The withdrawable amount should be equals to 0"
    )
  })
});
