import increaseTime, {duration, increaseTimeTo} from '../helpers/increaseTime';
import latestTime from '../helpers/latestTime';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage'
import {createSnapshot, restoreSnapshot} from "../helpers/snapshot"

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

// Depending on test network increaseTimeTo can be inconsistent and add
// extra time. As a workaround we subtract timeRoundMargin in all cases
// that test times before initialization/undelegation periods end.
const timeRoundMargin = duration.minutes(1)

const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const Registry = artifacts.require("./Registry.sol");

contract('TokenStaking', function(accounts) {

  let token, registry, stakingContract, stakingAmount, minimumStake;
    
  const ownerOne = accounts[0],
    ownerTwo = accounts[1],
    operatorOne = accounts[2],
    operatorTwo = accounts[3],
    magpie = accounts[4],
    authorizer = accounts[5],
    operatorContract = accounts[6];

  const initializationPeriod = duration.minutes(10);
  const undelegationPeriod = duration.minutes(30);

  before(async () => {
    token = await KeepToken.new();
    registry = await Registry.new();
    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod
    );

    await registry.approveOperatorContract(operatorContract);

    minimumStake = await stakingContract.minimumStake();
    stakingAmount = minimumStake.muln(20);
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  async function delegate(operator, amount) {
    let data = Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operator.substr(2), 'hex'),
      Buffer.from(authorizer.substr(2), 'hex')
    ]);
    
    return token.approveAndCall(
      stakingContract.address, amount, 
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

  it("should allow to delegate, undelegate, and recover stake", async () => {
    let ownerStartBalance = await token.balanceOf.call(ownerOne)

    let tx = await delegate(operatorOne, stakingAmount)
    let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

    await increaseTimeTo(createdAt + initializationPeriod + 1)

    await stakingContract.undelegate(operatorOne, {from: ownerOne});
    await increaseTime(undelegationPeriod + 1);
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

  describe("delegate", async () => {
    it("should update balances", async () => {
      let ownerStartBalance = await token.balanceOf.call(ownerOne);
  
      await delegate(operatorOne, stakingAmount);
      
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

    it("should not allow to delegate to the same operator twice", async () => {
      await delegate(operatorOne, stakingAmount)
  
      await expectThrowWithMessage(
        delegate(operatorOne, stakingAmount),
        "Operator address is already in use."
      )
    })
  
    it("should not allow to delegate to the same operator even after recovering stake", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
  
      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: ownerOne});
      await increaseTime(undelegationPeriod + 1);
      await stakingContract.recoverStake(operatorOne);
          
      await expectThrowWithMessage(
        delegate(operatorOne, stakingAmount),
        "Operator address is already in use."
      )
    })
  
    it("should not allow to delegate less than the minimum stake", async () => {    
      await expectThrowWithMessage(
        delegate(operatorOne, minimumStake.subn(1)),
        "Tokens amount must be greater than the minimum stake"
      )
    })
  
    it("should allow to delegate the minimum stake", async () => {    
      await delegate(operatorOne, minimumStake)
      // ok, no reverts
    })
  
    it("should allow to delegate to two different operators", async () => {
      let ownerStartBalance = await token.balanceOf.call(ownerOne)
  
      await delegate(operatorOne, stakingAmount);
      await delegate(operatorTwo, stakingAmount);
  
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
  })

  describe("cancelStake", async () => {
    it("should let operator cancel delegation", async () => {
      await delegate(operatorOne, stakingAmount)
  
      await stakingContract.cancelStake(operatorOne, {from: operatorOne})
      // ok, no revert
    })
  
    it("should not allow third party to cancel delegation", async () => {
      await delegate(operatorOne, stakingAmount)
  
      await expectThrowWithMessage(
        stakingContract.cancelStake(operatorOne, {from: operatorTwo}),
        "Only operator or the owner of the stake can cancel the delegation"
      )
    })

    it("should allow to cancel delegation right away", async () => {
      let ownerStartBalance = await token.balanceOf.call(ownerOne);
  
      await delegate(operatorOne, stakingAmount);
  
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
      let ownerStartBalance = await token.balanceOf.call(ownerOne);
      
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
  
      await increaseTimeTo(createdAt + initializationPeriod - timeRoundMargin)
  
      await stakingContract.cancelStake(operatorOne, {from: ownerOne})
  
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
  
    it("should not allow to cancel delegation after initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
  
      await increaseTimeTo(createdAt + initializationPeriod + 1)
  
      await expectThrowWithMessage(
        stakingContract.cancelStake(operatorOne, {from: ownerOne}),
        "Initialization period is over"
      );
    })

    it("should retain delegation info", async () => {
      await delegate(operatorOne, stakingAmount);
  
      let delegationInfoBefore = await stakingContract.getDelegationInfo.call(operatorOne)
  
      await stakingContract.cancelStake(operatorOne, {from: ownerOne});
  
      let delegationInfoAfter = await stakingContract.getDelegationInfo.call(operatorOne)
  
      expect(delegationInfoAfter.createdAt).to.eq.BN(
        delegationInfoBefore.createdAt,
        "Unexpected delegation creation time"
      )
      expect(delegationInfoAfter.amount).to.eq.BN(
        0,
        "Should have no delegated tokens"
      )
      expect(delegationInfoAfter.undelegatedAt).to.eq.BN(
        0,
        "Unexpected undelegation time"
      )
    })
  })

  describe("undelegate", async () => {
    it("should let operator undelegate", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: operatorOne})
      // ok, no revert
    })

    it("should not allow third party to undelegate", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await expectThrowWithMessage(
        stakingContract.undelegate(operatorOne, {from: operatorTwo}),
        "Only operator or the owner of the stake can undelegate"
      )
    })

    it("should permit undelegating at the time when initialization " + 
    "period passed", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: operatorOne})
      // ok, no revert
    })

    it("should not permit undelegating at the time before initialization " + 
    "period passed", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod - timeRoundMargin)
      await expectThrowWithMessage(
        stakingContract.undelegate(operatorOne, {from: operatorOne}),
        "Cannot undelegate in initialization period, use cancelStake instead"
      )
    })

    it("should let the operator undelegate earlier", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      let currentTime = await latestTime()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime + 20,
        {from: operatorOne}
      )

      await stakingContract.undelegate(operatorOne, {from: operatorOne})
      // ok, no revert
    })

    it("should let the owner postpone undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: operatorOne})

      await stakingContract.undelegate(
        operatorOne,
        {from: ownerOne}
      )
      // ok, no revert
    })

    it("should not let the operator postpone undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: operatorOne})

      await expectThrowWithMessage(
        stakingContract.undelegate(operatorOne, {from: operatorOne}),
        "Only the owner may postpone previously set undelegation"
      )
    })
  })

  describe("undelegateAt", async () => {
    it("should let operator undelegate", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)

      let currentTime = await latestTime()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime + 1,
        {from: operatorOne}
      )
      // ok, no revert
    })

    it("should not allow third party to undelegate", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)

      let currentTime = await latestTime()

      await expectThrowWithMessage(
        stakingContract.undelegateAt(
          operatorOne, currentTime + 10,
          {from: operatorTwo}
        ),
        "Only operator or the owner of the stake can undelegate"
      )
    })

    it("should permit undelegating at the current time", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)

      let currentTime = await latestTime()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime + 1,
        {from: operatorOne}
      )
      // ok, no revert
    })

    it("should permit undelegating at the time when initialization " +
    "period passed", async () => {
      await delegate(operatorOne, stakingAmount)

      let currentTime = await latestTime()
      await stakingContract.undelegateAt(
        operatorOne, currentTime + initializationPeriod + 1,
        {from: operatorOne}
      )
      // ok, no revert
    })

    it("should not permit undelegating at the time before initialization " + 
    "period passed", async () => {
      await delegate(operatorOne, stakingAmount)

      let currentTime = await latestTime()
      await expectThrowWithMessage(
        stakingContract.undelegateAt(
          operatorOne, currentTime + initializationPeriod - timeRoundMargin,
          {from: operatorOne}
        ),
        "Cannot undelegate in initialization period, use cancelStake instead"
      )
    })

    it("should not permit undelegating in the past", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)

      let currentTime = await latestTime()

      await expectThrowWithMessage(
        stakingContract.undelegateAt(
          operatorOne, currentTime - 1,
          {from: operatorOne}
        ),
        "May not set undelegation timestamp in the past"
      )
    })

    it("should let the operator undelegate earlier", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      let currentTime = await latestTime()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime + 20,
        {from: operatorOne}
      )

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime + 1,
        {from: operatorOne}
      )
      // ok, no revert
    })

    it("should let the owner postpone undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: operatorOne})

      let currentTime = await latestTime()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime + 1,
        {from: ownerOne}
      )
      // ok, no revert
    })

    it("should not let the operator postpone undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: operatorOne})

      let currentTime = await latestTime()

      await expectThrowWithMessage(
        stakingContract.undelegateAt(
          operatorOne, currentTime + 1,
          {from: operatorOne}
        ),
        "Only the owner may postpone previously set undelegation"
      )
    })
  })

  describe("recoverStake", async () => {
    it("should not allow to recover stake before undelegation period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
  
      await increaseTimeTo(createdAt + initializationPeriod + 1)
      tx = await stakingContract.undelegate(operatorOne, {from: ownerOne});
      let undelegatedAt = await web3.eth.getBlock(tx.receipt.blockNumber);
      await increaseTimeTo(undelegatedAt.timestamp + undelegationPeriod - timeRoundMargin);
  
      await expectThrowWithMessage(
        stakingContract.recoverStake(operatorOne),
        "Can not recover stake before undelegation period is over"
      )
    })

    it("should retain delegation info", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await increaseTimeTo(createdAt + initializationPeriod + 1)
  
      let delegationInfoBefore = await stakingContract.getDelegationInfo.call(operatorOne)
      
      await stakingContract.undelegate(operatorOne, {from: ownerOne})
      let blockNumber = await web3.eth.getBlockNumber()
      let undelegationBlock = await web3.eth.getBlock(blockNumber)
    
      await increaseTime(undelegationPeriod + 1)
      await stakingContract.recoverStake(operatorOne)
  
      let delegationInfoAfter = await stakingContract.getDelegationInfo.call(operatorOne)
  
      expect(delegationInfoAfter.createdAt).to.eq.BN(
        delegationInfoBefore.createdAt,
        "Unexpected delegation creation time"
      )
      expect(delegationInfoAfter.amount).to.eq.BN(
        0,
        "Should have no delegated tokens"
      )
      expect(delegationInfoAfter.undelegatedAt).to.eq.BN(
        undelegationBlock.timestamp,
        "Unexpected undelegation time"
      )
    })
  })

  describe("activeStake", async () => {
    it("should report active stake after initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await increaseTimeTo(createdAt + initializationPeriod + 1)
  
      let activeStake = await stakingContract.activeStake.call(operatorOne, operatorContract)
  
      expect(activeStake).to.eq.BN(
        stakingAmount,
        "Active stake should equal staked amount"
      )
    })
  
    it("should report no active stake before initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await increaseTimeTo(createdAt + initializationPeriod - timeRoundMargin)
  
      let activeStake = await stakingContract.activeStake.call(operatorOne, operatorContract)
  
      expect(activeStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  
    it("should report no active stake for not authorized operator contract", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await increaseTimeTo(createdAt + initializationPeriod + 1)
  
      let activeStake = await stakingContract.activeStake.call(operatorOne, operatorContract)
  
      expect(activeStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  
    it("should report no active stake after cancelling delegation", async () => {
      await delegate(operatorOne, stakingAmount);
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await stakingContract.cancelStake(operatorOne, {from: ownerOne});
  
      let activeStake = await stakingContract.activeStake.call(operatorOne, operatorContract)
  
      expect(activeStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  
    it("should report no active stake after recovering stake", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: ownerOne});
      await increaseTime(undelegationPeriod + 1);    
      await stakingContract.recoverStake(operatorOne);
      
      let activeStake = await stakingContract.activeStake.call(operatorOne, operatorContract)
  
      expect(activeStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  })
  
  describe("eligibleStake", async () => {
    it("should report eligible stake after initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await increaseTimeTo(createdAt + initializationPeriod + 1)
  
      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
  
      expect(eligibleStake).to.eq.BN(
        stakingAmount,
        "Eligible stake should equal staked amount"
      )
    })
  
    it("should report no eligible stake before initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
      await increaseTimeTo(createdAt + initializationPeriod - timeRoundMargin)
      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
  
      expect(eligibleStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  
    it("should report no eligible stake for not authorized operator contract", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
  
      await increaseTimeTo(createdAt + initializationPeriod + 1)
  
      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
  
      expect(eligibleStake).to.eq.BN(
        0,
        "There should be no eligible stake"
      )
    })
  
    it("should report no eligible stake after cancelling delegation", async () => {
      await delegate(operatorOne, stakingAmount);
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await stakingContract.cancelStake(operatorOne, {from: ownerOne})
  
      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
  
      expect(eligibleStake).to.eq.BN(
        0,
        "There should be no eligible stake"
      )
    })
  
    it("should report no eligible stake when undelegating", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await increaseTimeTo(createdAt + initializationPeriod + 1)
      await stakingContract.undelegate(operatorOne, {from: ownerOne})
  
      await increaseTime(1);

      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
  
      expect(eligibleStake).to.eq.BN(
        0,
        "There should be no eligible stake"
      )
    })

    it("should report eligible stake for future undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      const delegationTime = 10

      let currentTime = await latestTime()
      await stakingContract.undelegateAt(
        operatorOne, 
        currentTime + initializationPeriod + delegationTime, 
        {from: ownerOne}
      );

      await increaseTimeTo(createdAt + initializationPeriod + 1)
      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
      expect(eligibleStake).to.eq.BN(
        stakingAmount,
        "Eligible stake should equal staked amount"
      )

      await increaseTime(delegationTime - 1);
      eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
      expect(eligibleStake).to.eq.BN(
        stakingAmount,
        "Eligible stake should equal staked amount"
      )
    })

    it("should report no eligible stake for passed future undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      const delegationTime = 10

      let currentTime = await latestTime()
      await stakingContract.undelegateAt(
        operatorOne, 
        currentTime + initializationPeriod + delegationTime, 
        {from: ownerOne}
      );

      await increaseTimeTo(createdAt + initializationPeriod + delegationTime + 1)

      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
      expect(eligibleStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  })
});
