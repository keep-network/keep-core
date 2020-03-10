import mineBlocks from './helpers/mineBlocks';
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

  let token, registry, stakingContract;
    
  const ownerOne = accounts[0],
    ownerTwo = accounts[1],
    operatorOne = accounts[2],
    operatorTwo = accounts[3],
    magpie = accounts[4],
    authorizer = accounts[5],
    operatorContract = accounts[6];
    
  const initializationPeriod = 10;
  const undelegationPeriod = 30;
  const stakingAmount = web3.utils.toBN(10000000);

  before(async () => {
    token = await KeepToken.new();
    registry = await Registry.new();
    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod
    );

    await registry.approveOperatorContract(operatorContract);
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  async function delegate(operator) {
    let data = Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operator.substr(2), 'hex'),
      Buffer.from(authorizer.substr(2), 'hex')
    ]);
    
    return token.approveAndCall(
      stakingContract.address, stakingAmount, 
      '0x' + data.toString('hex'), 
      {from: ownerOne}
    );
  }

  it("should send tokens correctly", async () => {
    let amount = web3.utils.toBN(1000000000);

    let ownerOneStartingBalance = await token.balanceOf.call(ownerOne);
    let ownerTwoStartingBalance = await token.balanceOf.call(ownerTwo);

    await token.transfer(ownerTwo, amount, {from: ownerOne});

    let ownerOneEndingBalance = await token.balanceOf.call(ownerOne);
    let ownerTwoEndingBalance = await token.balanceOf.call(ownerTwo);

    expect(ownerOneEndingBalance).to.eq.BN(
      ownerOneStartingBalance.sub(amount), 
      "Amount wasn't correctly taken from the sender"
    )
    expect(ownerTwoEndingBalance).to.eq.BN(
      ownerTwoStartingBalance.add(amount), 
      "Amount wasn't correctly sent to the receiver"
    );
  });

  it("should update balances when delegating", async () => {
    let ownerStartBalance = await token.balanceOf.call(ownerOne);

    await delegate(operatorOne);
    
    let ownerEndBalance = await token.balanceOf.call(ownerOne);
    let operatorEndStakeBalance = await stakingContract.balanceOf.call(operatorOne);
    
    expect(ownerEndBalance).to.eq.BN(
      ownerStartBalance.sub(stakingAmount),
      "Staking amount should be transferred from owner balance"
    );
    expect(operatorEndStakeBalance).to.eq.BN(
      stakingAmount,
      "Staking amount should be added to the operator balance"
    ); 
  })

  it("should allow to delegate, undelegate, and recover stake", async () => {
    let ownerStartBalance = await token.balanceOf.call(ownerOne)

    await delegate(operatorOne);

    await mineBlocks(initializationPeriod);
    await stakingContract.undelegate(operatorOne, {from: operatorOne});
    await mineBlocks(undelegationPeriod);    
    await stakingContract.recoverStake(operatorOne);
        
    let ownerEndBalance = await token.balanceOf.call(ownerOne);
    let operatorEndStakeBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(ownerEndBalance).to.eq.BN(
      ownerStartBalance,
      "Staking amount should be transferred back to owner"
    );
    expect(operatorEndStakeBalance).to.eq.BN( 
      0, 
      "Staking amount should be removed from operator balance"
    );
  })

  it("should allow to cancel delegation right away", async () => {
    let ownerStartBalance = await token.balanceOf.call(ownerOne);

    await delegate(operatorOne);

    await stakingContract.cancelStake(operatorOne, {from: ownerOne});

    let ownerEndBalance = await token.balanceOf.call(ownerOne);
    let operatorEndStakeBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(ownerEndBalance).to.eq.BN(
      ownerStartBalance,
      "Staking amount should be transferred back to owner"
    );
    expect(operatorEndStakeBalance).to.eq.BN( 
      0, 
      "Staking amount should be removed from operator balance"
    );
  })

  it("should allow to cancel delegation just before initialization period is over", async () => {
    await delegate(operatorOne);

    await mineBlocks(initializationPeriod - 2)

    await stakingContract.cancelStake(operatorOne, {from: ownerOne})
  })

  it("should not allow to cancel delegation after initialization period is over", async () => {
    await delegate(operatorOne);

    await mineBlocks(initializationPeriod);

    await expectThrowWithMessage(
      stakingContract.cancelStake(operatorOne, {from: ownerOne}),
      "Initialization period is over"
    );
  })

  it("should not allow to recover stake before undelegation period is over", async () => {
    await delegate(operatorOne);

    await mineBlocks(initializationPeriod);
    await stakingContract.undelegate(operatorOne, {from: operatorOne});

    await expectThrowWithMessage(
      stakingContract.recoverStake(operatorOne),
      "Can not recover stake before undelgation period is over"
    )
  })

  it("should not allow to delegate to the same operator twice", async () => {
    await delegate(operatorOne);

    await expectThrowWithMessage(
      delegate(operatorOne),
      "Operator address is already in use."
    )
  })

  it("should not allow to delegate to the same operator even after recovering stake", async () => {
    await delegate(operatorOne);

    await mineBlocks(initializationPeriod);
    await stakingContract.undelegate(operatorOne, {from: operatorOne});
    await mineBlocks(undelegationPeriod);    
    await stakingContract.recoverStake(operatorOne);
        
    await expectThrowWithMessage(
      delegate(operatorOne),
      "Operator address is already in use."
    )
  })

  it("should allow to delegate to two different operators", async () => {
    let ownerStartBalance = await token.balanceOf.call(ownerOne)

    await delegate(operatorOne);
    await delegate(operatorTwo);

    let ownerEndBalance = await token.balanceOf.call(ownerOne);
    let operatorOneEndStakeBalance = await stakingContract.balanceOf.call(operatorOne);
    let operatorTwoEndStakeBalance = await stakingContract.balanceOf.call(operatorTwo);

    expect(ownerEndBalance).to.eq.BN(
      ownerStartBalance.sub(stakingAmount).sub(stakingAmount),
      "Staking amount should be transferred from owner balance"
    );
    expect(operatorOneEndStakeBalance).to.eq.BN(
      stakingAmount,
      "Staking amount should be added to the operator balance"
    );
    expect(operatorTwoEndStakeBalance).to.eq.BN(
      stakingAmount,
      "Staking amount should be added to the operator balance"
    );
  })

  it("should retain delegation info after recovering stake", async () => {
    await delegate(operatorOne)
    await mineBlocks(initializationPeriod)

    let delegationInfoBefore = await stakingContract.getDelegationInfo.call(operatorOne)
    
    await stakingContract.undelegate(operatorOne, {from: operatorOne})
    let undelegationBlock = await web3.eth.getBlockNumber()
    await mineBlocks(undelegationPeriod)
    await stakingContract.recoverStake(operatorOne)

    let delegationInfoAfter = await stakingContract.getDelegationInfo.call(operatorOne)

    expect(delegationInfoAfter.createdAt).to.eq.BN(
      delegationInfoBefore.createdAt,
      "Unexpected delegation creation block"
    )
    expect(delegationInfoAfter.amount).to.eq.BN(
      0,
      "Should have no delegated tokens"
    )
    expect(delegationInfoAfter.undelegatedAt).to.eq.BN(
      undelegationBlock,
      "Unexpected undelegation block"
    )
  })

  it("should retain delegation info after cancelling delegation", async () => {
    await delegate(operatorOne);

    let delegationInfoBefore = await stakingContract.getDelegationInfo.call(operatorOne)

    await stakingContract.cancelStake(operatorOne, {from: ownerOne});

    let delegationInfoAfter = await stakingContract.getDelegationInfo.call(operatorOne)

    expect(delegationInfoAfter.createdAt).to.eq.BN(
      delegationInfoBefore.createdAt,
      "Unexpected delegation creation block"
    )
    expect(delegationInfoAfter.amount).to.eq.BN(
      0,
      "Should have no delegated tokens"
    )
    expect(delegationInfoAfter.undelegatedAt).to.eq.BN(
      0,
      "Unexpected undelegation block"
    )
  })

  it("should report active stake after initialization period is over", async () => {
    await delegate(operatorOne)
    await stakingContract.authorizeOperatorContract(
      operatorOne, operatorContract, {from: authorizer}
    )

    await mineBlocks(initializationPeriod);

    let activeStake = await stakingContract.activeStake(operatorOne, operatorContract)

    expect(activeStake).to.eq.BN(
      stakingAmount,
      "Active stake should equal staked amount"
    )
  })

  it("should report no active stake before initialization period is over", async () => {
    await delegate(operatorOne)
    await stakingContract.authorizeOperatorContract(
      operatorOne, operatorContract, {from: authorizer}
    )

    let activeStake = await stakingContract.activeStake(operatorOne, operatorContract)

    expect(activeStake).to.eq.BN(
      0,
      "There should be no active stake"
    )
  })

  it("should report no active stake for not authorized operator contract", async () => {
    await delegate(operatorOne)
    await mineBlocks(initializationPeriod);

    let activeStake = await stakingContract.activeStake(operatorOne, operatorContract)

    expect(activeStake).to.eq.BN(
      0,
      "There should be no active stake"
    )
  })

  it("should report no active stake after cancelling delegation", async () => {
    await delegate(operatorOne);
    await stakingContract.authorizeOperatorContract(
      operatorOne, operatorContract, {from: authorizer}
    )

    await stakingContract.cancelStake(operatorOne, {from: ownerOne});

    let activeStake = await stakingContract.activeStake(operatorOne, operatorContract)

    expect(activeStake).to.eq.BN(
      0,
      "There should be no active stake"
    )
  })

  it("should report no active stake after recovering stake", async () => {
    await delegate(operatorOne);
    await stakingContract.authorizeOperatorContract(
      operatorOne, operatorContract, {from: authorizer}
    )

    await mineBlocks(initializationPeriod);
    await stakingContract.undelegate(operatorOne, {from: operatorOne});
    await mineBlocks(undelegationPeriod);    
    await stakingContract.recoverStake(operatorOne);
    
    let activeStake = await stakingContract.activeStake(operatorOne, operatorContract)

    expect(activeStake).to.eq.BN(
      0,
      "There should be no active stake"
    )
  })

  it("should report eligible stake after initialization period is over", async () => {
    await delegate(operatorOne)
    await stakingContract.authorizeOperatorContract(
      operatorOne, operatorContract, {from: authorizer}
    )

    await mineBlocks(initializationPeriod);

    let eligibleStake = await stakingContract.eligibleStake(operatorOne, operatorContract)

    expect(eligibleStake).to.eq.BN(
      stakingAmount,
      "Eligible stake should equal staked amount"
    )
  })

  it("should report no eligible stake before initialization period is over", async () => {
    await delegate(operatorOne)
    await stakingContract.authorizeOperatorContract(
      operatorOne, operatorContract, {from: authorizer}
    )

    let eligibleStake = await stakingContract.eligibleStake(operatorOne, operatorContract)

    expect(eligibleStake).to.eq.BN(
      0,
      "There should be no active stake"
    )
  })

  it("should report no eligible stake for not authorized operator contract", async () => {
    await delegate(operatorOne)

    await mineBlocks(initializationPeriod);

    let eligibleStake = await stakingContract.eligibleStake(operatorOne, operatorContract)

    expect(eligibleStake).to.eq.BN(
      0,
      "There should be no eligible stake"
    )
  })

  it("should report no eligible stake after cancelling delegation", async () => {
    await delegate(operatorOne);
    await stakingContract.authorizeOperatorContract(
      operatorOne, operatorContract, {from: authorizer}
    )

    await stakingContract.cancelStake(operatorOne, {from: ownerOne})

    let eligibleStake = await stakingContract.eligibleStake(operatorOne, operatorContract)

    expect(eligibleStake).to.eq.BN(
      0,
      "There should be no eligible stake"
    )
  })

  it("should report no eligible stake when undelegating", async () => {
    await delegate(operatorOne);
    await stakingContract.authorizeOperatorContract(
      operatorOne, operatorContract, {from: authorizer}
    )

    await mineBlocks(initializationPeriod);
    await stakingContract.undelegate(operatorOne, {from: operatorOne})

    await mineBlocks(1)

    let eligibleStake = await stakingContract.eligibleStake(operatorOne, operatorContract)

    expect(eligibleStake).to.eq.BN(
      0,
      "There should be no eligible stake"
    )
  })
});
