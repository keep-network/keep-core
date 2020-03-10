import mineBlocks from './helpers/mineBlocks';
import expectThrow from './helpers/expectThrow';
import expectThrowWithMessage from './helpers/expectThrowWithMessage'
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"

const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const Registry = artifacts.require("./Registry.sol");

const BN = web3.utils.BN

const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

contract('TokenStaking', function(accounts) {

  let token, registry, stakingContract,
    account_one = accounts[0],
    account_one_operator = accounts[1],
    account_one_magpie = accounts[2],
    account_one_authorizer = accounts[3],
    account_two = accounts[4];

  const initializationPeriod = 10;
  const undelegationPeriod = 30;

  const stakingAmount = web3.utils.toBN(10000000);

  before(async () => {
    token = await KeepToken.new();
    registry = await Registry.new();
    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod
    );
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  async function delegate() {
    let data = Buffer.concat([
      Buffer.from(account_one_magpie.substr(2), 'hex'),
      Buffer.from(account_one_operator.substr(2), 'hex'),
      Buffer.from(account_one_authorizer.substr(2), 'hex')
    ]);
    
    await token.approveAndCall(
      stakingContract.address, stakingAmount, 
      '0x' + data.toString('hex'), 
      {from: account_one}
    );
  }

  it("should send tokens correctly", async function() {
    let amount = web3.utils.toBN(1000000000);

    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Send tokens
    await token.transfer(account_two, amount, {from: account_one});

    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_two_ending_balance = await token.balanceOf.call(account_two);

    expect(account_one_ending_balance).to.eq.BN(
      account_one_starting_balance.sub(amount), 
      "Amount wasn't correctly taken from the sender"
    )
    expect(account_two_ending_balance).to.eq.BN(
      account_two_starting_balance.add(amount), 
      "Amount wasn't correctly sent to the receiver"
    );
  });

  it("should update balances when delegating", async () => {
    let ownerStartBalance = await token.balanceOf.call(account_one);

    await delegate();
    
    let ownerEndBalance = await token.balanceOf.call(account_one);
    let operatorEndStakeBalance = await stakingContract.balanceOf.call(account_one_operator);
    
    expect(ownerEndBalance).to.eq.BN(
      ownerStartBalance.sub(stakingAmount),
      "Staking amount should be transferred from owner balance"
    );
    expect(operatorEndStakeBalance).to.eq.BN(
      stakingAmount,
      "Staking amount should be added to the operator balance"
    ); 
  })

  it("should allow to cancel delegation", async () => {
    let ownerStartBalance = await token.balanceOf.call(account_one);

    await delegate();

    await stakingContract.cancelStake(account_one_operator, {from: account_one});

    let ownerEndBalance = await token.balanceOf.call(account_one);
    let operatorEndStakeBalance = await stakingContract.balanceOf.call(account_one_operator);

    expect(ownerEndBalance).to.eq.BN(
      ownerStartBalance,
      "Staking amount should be transferred back to owner"
    );
    expect(operatorEndStakeBalance).to.eq.BN( 
      0, 
      "Staking amount should be removed from operator balance"
    );
  })

  it("should allow to cancel stake before initialization period is over", async () => {
    await delegate()

    await mineBlocks(initializationPeriod - 2)

    await stakingContract.cancelStake(account_one_operator, {from: account_one})
  })

  it("should not allow to cancel stake after initialization period is over", async () => {
    await delegate();

    await mineBlocks(initializationPeriod);

    await expectThrowWithMessage(
      stakingContract.cancelStake(account_one_operator, {from: account_one}),
      "Initialization period is over"
    );
  })

  it("should not allow to recover stake before undelegation period is over", async () => {
    await delegate();

    await mineBlocks(initializationPeriod);
    await stakingContract.undelegate(account_one_operator, {from: account_one_operator});

    await expectThrowWithMessage(
      stakingContract.recoverStake(account_one_operator),
      "Can not recover stake before undelgation period is over"
    )
  })

  it("should stake delegate and undelegate tokens correctly", async function() {
    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);

    let data = Buffer.concat([
      Buffer.from(account_one_magpie.substr(2), 'hex'),
      Buffer.from(account_one_operator.substr(2), 'hex'),
      Buffer.from(account_one_authorizer.substr(2), 'hex')
    ]);

    // Stake tokens using approveAndCall pattern
    await token.approveAndCall(stakingContract.address, stakingAmount, '0x' + data.toString('hex'), {from: account_one});

    // jump in time, full initialization period
    await mineBlocks(initializationPeriod);

    // Can not cancel stake
    await expectThrow(stakingContract.cancelStake(account_one_operator, {from: account_one}));

    // Undelegate tokens as operator
    await stakingContract.undelegate(account_one_operator, {from: account_one_operator});

    // should not be able to recover stake
    await expectThrow(stakingContract.recoverStake(account_one_operator));

    // jump in time, full undelegation period
    await mineBlocks(undelegationPeriod);

    // should be able to recover stake
    await stakingContract.recoverStake(account_one_operator);

    // should fail cause there is no stake to recover
    await expectThrow(stakingContract.recoverStake(account_one_operator));

    // check balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_one_operator_stake_balance = await stakingContract.balanceOf.call(account_one_operator);

    expect(account_one_ending_balance).to.eq.BN(
      account_one_starting_balance, 
      "Staking amount should be transfered to sender balance"
    );
    expect(account_one_operator_stake_balance).to.eq.BN(
      0, 
      "Staking amount should be removed from sender staking balance"
    );

    // Starting balances
    account_one_starting_balance = await token.balanceOf.call(account_one);

    data = Buffer.concat([
      Buffer.from(account_one_magpie.substr(2), 'hex'),
      Buffer.from(account_one_operator.substr(2), 'hex'),
      Buffer.from(account_one_authorizer.substr(2), 'hex')
    ]);

    // Stake tokens using approveAndCall pattern
    await token.approveAndCall(stakingContract.address, stakingAmount, '0x' + data.toString('hex'), {from: account_one});

    // Ending balances
    account_one_ending_balance = await token.balanceOf.call(account_one);
    account_one_operator_stake_balance = await stakingContract.balanceOf.call(account_one_operator);

    expect(account_one_ending_balance).to.eq.BN(
      account_one_starting_balance.sub(stakingAmount), 
      "Staking amount should be transfered from sender balance for the second time"
    );
    expect(account_one_operator_stake_balance).to.eq.BN(
      stakingAmount, 
      "Staking amount should be added to the sender staking balance for the second time"
    );
  });
});
