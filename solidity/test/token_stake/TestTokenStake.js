const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

// Depending on test network increaseTimeTo can be inconsistent and add
// extra time. As a workaround we subtract timeRoundMargin in all cases
// that test times before initialization/undelegation periods end.
const timeRoundMargin = time.duration.minutes(1)

const KeepToken = contract.fromArtifact('KeepToken');
const MinimumStakeSchedule = contract.fromArtifact('MinimumStakeSchedule');
const TokenGrant = contract.fromArtifact('TokenGrant');
const TokenStaking = contract.fromArtifact('TokenStaking');
const GrantStakingInfo = contract.fromArtifact('GrantStakingInfo');
const TokenStakingEscrow = contract.fromArtifact('TokenStakingEscrow');
const KeepRegistry = contract.fromArtifact("KeepRegistry");

describe('TokenStaking', function() {

  let token, registry, stakingContract, stakingAmount, minimumStake;
    
  const ownerOne = accounts[0],
    ownerTwo = accounts[1],
    operatorOne = accounts[2],
    operatorTwo = accounts[3],
    beneficiary = accounts[4],
    authorizer = accounts[5],
    operatorContract = accounts[6];

  const initializationPeriod = time.duration.minutes(10);
  const undelegationPeriod = time.duration.minutes(30);

  before(async () => {
    token = await KeepToken.new({from: accounts[0]});
    tokenGrant = await TokenGrant.new(token.address,  {from: accounts[0]});
    registry = await KeepRegistry.new({from: accounts[0]});
    stakingEscrow = await TokenStakingEscrow.new(
      token.address, 
      tokenGrant.address, 
      {from: accounts[0]}
    );
    await TokenStaking.detectNetwork();
    await TokenStaking.link(
      'MinimumStakeSchedule', 
      (await MinimumStakeSchedule.new({from: accounts[0]})).address
    );
    await TokenStaking.link(
      'GrantStakingInfo', 
      (await GrantStakingInfo.new({from: accounts[0]})).address
    );
    stakingContract = await TokenStaking.new(
      token.address,
      tokenGrant.address,
      stakingEscrow.address,
      registry.address,
      initializationPeriod,
      undelegationPeriod,
      {from: accounts[0]}
    );
    await stakingEscrow.transferOwnership(stakingContract.address, {from: accounts[0]});
    await registry.approveOperatorContract(operatorContract, {from: accounts[0]});

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
      Buffer.from(beneficiary.substr(2), 'hex'),
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
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

    await time.increaseTo(createdAt.add(initializationPeriod).addn(1))

    await stakingContract.undelegate(operatorOne, {from: ownerOne});
    await time.increase(undelegationPeriod.addn(1));
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
  
      await expectRevert(
        delegate(operatorOne, stakingAmount),
        "Operator already in use"
      )
    })
  
    it("should not allow to delegate to the same operator even after recovering stake", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: ownerOne});
      await time.increase(undelegationPeriod.addn(1));
      await stakingContract.recoverStake(operatorOne);
          
      await expectRevert(
        delegate(operatorOne, stakingAmount),
        "Operator already in use"
      )
    })
  
    it("should not allow to delegate less than the minimum stake", async () => {    
      await expectRevert(
        delegate(operatorOne, minimumStake.subn(1)),
        "Value must be greater than the minimum stake"
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
  
      await expectRevert(
        stakingContract.cancelStake(operatorOne, {from: operatorTwo}),
        "Unauthorized"
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
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).sub(timeRoundMargin))
  
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
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
      await expectRevert(
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
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: operatorOne})
      // ok, no revert
    })

    it("should not allow third party to undelegate", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await expectRevert(
        stakingContract.undelegate(operatorOne, {from: operatorTwo}),
        "Unauthorized"
      )
    })

    it("should permit undelegating at the time when initialization " + 
    "period passed", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: operatorOne})
      // ok, no revert
    })

    it("should not permit undelegating at the time before initialization " + 
    "period passed", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await time.increaseTo(createdAt.add(initializationPeriod).sub(timeRoundMargin))
      await expectRevert(
        stakingContract.undelegate(operatorOne, {from: operatorOne}),
        "Cannot undelegate in initialization period"
      )
    })

    it("should let the operator undelegate earlier", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      let currentTime = await time.latest()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime.addn(20),
        {from: operatorOne}
      )

      await stakingContract.undelegate(operatorOne, {from: operatorOne})
      // ok, no revert
    })

    it("should let the owner postpone undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: operatorOne})

      await stakingContract.undelegate(
        operatorOne,
        {from: ownerOne}
      )
      // ok, no revert
    })

    it("should not let the operator postpone undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: operatorOne})

      await expectRevert(
        stakingContract.undelegate(operatorOne, {from: operatorOne}),
        "Only the owner may postpone undelegation"
      )
    })
  })

  describe("undelegateAt", async () => {
    it("should let operator undelegate", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))

      let currentTime = await time.latest()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime.addn(1),
        {from: operatorOne}
      )
      // ok, no revert
    })

    it("should not allow third party to undelegate", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))

      let currentTime = await time.latest()

      await expectRevert(
        stakingContract.undelegateAt(
          operatorOne, currentTime.addn(10),
          {from: operatorTwo}
        ),
        "Unauthorized"
      )
    })

    it("should permit undelegating at the current time", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))

      let currentTime = await time.latest()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime.addn(1),
        {from: operatorOne}
      )
      // ok, no revert
    })

    it("should permit undelegating at the time when initialization " +
    "period passed", async () => {
      await delegate(operatorOne, stakingAmount)

      let currentTime = await time.latest()
      await stakingContract.undelegateAt(
        operatorOne, currentTime.add(initializationPeriod).addn(1),
        {from: operatorOne}
      )
      // ok, no revert
    })

    it("should not permit undelegating at the time before initialization " + 
    "period passed", async () => {
      await delegate(operatorOne, stakingAmount)

      let currentTime = await time.latest()
      await expectRevert(
        stakingContract.undelegateAt(
          operatorOne, currentTime.add(initializationPeriod).sub(timeRoundMargin),
          {from: operatorOne}
        ),
        "Cannot undelegate in initialization period"
      )
    })

    it("should not permit undelegating in the past", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))

      let currentTime = await time.latest()

      await expectRevert(
        stakingContract.undelegateAt(
          operatorOne, currentTime - 1,
          {from: operatorOne}
        ),
        "Undelegation timestamp in the past"
      )
    })

    it("should let the operator undelegate earlier", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      let currentTime = await time.latest()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime.addn(20),
        {from: operatorOne}
      )

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime.addn(1),
        {from: operatorOne}
      )
      // ok, no revert
    })

    it("should let the owner postpone undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: operatorOne})

      let currentTime = await time.latest()

      await stakingContract.undelegateAt(
        operatorOne,
        currentTime.addn(1),
        {from: ownerOne}
      )
      // ok, no revert
    })

    it("should not let the operator postpone undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: operatorOne})

      let currentTime = await time.latest()

      await expectRevert(
        stakingContract.undelegateAt(
          operatorOne, currentTime.addn(1),
          {from: operatorOne}
        ),
        "Only the owner may postpone undelegation"
      )
    })
  })

  describe("recoverStake", async () => {
    it("should not allow to recover stake without undelegating first", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).add(undelegationPeriod))
  
      await expectRevert(
        stakingContract.recoverStake(operatorOne),
        "Can not recover without first undelegating"
      )
    })

    it("should not allow to recover stake before undelegation period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      tx = await stakingContract.undelegate(operatorOne, {from: ownerOne});
      let undelegatedAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp);
      await time.increaseTo(undelegatedAt.add(undelegationPeriod).sub(timeRoundMargin));
  
      await expectRevert(
        stakingContract.recoverStake(operatorOne),
        "Can not recover before undelegation period is over"
      )
    })

    it("should retain delegation info", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
      let delegationInfoBefore = await stakingContract.getDelegationInfo.call(operatorOne)
      
      await stakingContract.undelegate(operatorOne, {from: ownerOne})
      let blockNumber = await web3.eth.getBlockNumber()
      let undelegationBlock = await web3.eth.getBlock(blockNumber)
      
      await time.increase(undelegationPeriod.addn(1));
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
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
      let activeStake = await stakingContract.activeStake.call(operatorOne, operatorContract)
  
      expect(activeStake).to.eq.BN(
        stakingAmount,
        "Active stake should equal staked amount"
      )
    })
  
    it("should report no active stake before initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await time.increaseTo(createdAt.add(initializationPeriod).sub(timeRoundMargin))
  
      let activeStake = await stakingContract.activeStake.call(operatorOne, operatorContract)
  
      expect(activeStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  
    it("should report no active stake for not authorized operator contract", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
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
  
    it("should report no active stake after undelegation is finished", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: ownerOne});
      await time.increase(undelegationPeriod.addn(1));

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
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
  
      expect(eligibleStake).to.eq.BN(
        stakingAmount,
        "Eligible stake should equal staked amount"
      )
    })
  
    it("should report no eligible stake before initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
      await time.increaseTo(createdAt.add(initializationPeriod).sub(timeRoundMargin))
      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
  
      expect(eligibleStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  
    it("should report no eligible stake for not authorized operator contract", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
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
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: ownerOne})
  
      await time.increase(1);

      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
  
      expect(eligibleStake).to.eq.BN(
        0,
        "There should be no eligible stake"
      )
    })

    it("should report eligible stake for future undelegation", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      const delegationTime = time.duration.seconds(10)

      let currentTime = await time.latest()
      let undelegateAt = currentTime.add(initializationPeriod).add(delegationTime)
      await stakingContract.undelegateAt(
        operatorOne, 
        undelegateAt,
        {from: ownerOne}
      );

      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
      expect(eligibleStake).to.eq.BN(
        stakingAmount,
        "Eligible stake should equal staked amount"
      )

      await time.increaseTo(undelegateAt.subn(1));
      eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
      expect(eligibleStake).to.eq.BN(
        stakingAmount,
        "Eligible stake should equal staked amount"
      )
    })

    it("should report no eligible stake for passed future undelegation", async () => {
      await delegate(operatorOne, stakingAmount)
      await stakingContract.authorizeOperatorContract(
        operatorOne, operatorContract, {from: authorizer}
      )
  
      const delegationTime = time.duration.seconds(10)

      let currentTime = await time.latest()
      let undelegateAt = currentTime.add(initializationPeriod).add(delegationTime)

      await stakingContract.undelegateAt(
        operatorOne,
        undelegateAt,
        {from: ownerOne}
      );

      await time.increaseTo(undelegateAt.addn(1))

      let eligibleStake = await stakingContract.eligibleStake.call(operatorOne, operatorContract)
      expect(eligibleStake).to.eq.BN(
        0,
        "There should be no active stake"
      )
    })
  })
});
