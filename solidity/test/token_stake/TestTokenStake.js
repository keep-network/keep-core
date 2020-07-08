const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot')
const {initTokenStaking} = require('../helpers/initContracts')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

// Depending on test network increaseTimeTo can be inconsistent and add
// extra time. As a workaround we subtract timeRoundMargin in all cases
// that test times before initialization/undelegation periods end.
const timeRoundMargin = time.duration.minutes(1)

const KeepToken = contract.fromArtifact('KeepToken');
const TokenGrant = contract.fromArtifact('TokenGrant');
const KeepRegistry = contract.fromArtifact("KeepRegistry");

describe('TokenStaking', function() {

  let token, registry, stakingContract, stakingAmount, minimumStake;
    
  const owner = accounts[0],
    operatorOne = accounts[2],
    operatorTwo = accounts[3],
    beneficiary = accounts[4],
    authorizer = accounts[5],
    operatorContract = accounts[6];

  const initializationPeriod = time.duration.minutes(10);
  let undelegationPeriod;

  before(async () => {
    token = await KeepToken.new({from: accounts[0]});
    tokenGrant = await TokenGrant.new(token.address,  {from: accounts[0]});
    registry = await KeepRegistry.new({from: accounts[0]});
    const stakingContracts = await initTokenStaking(
      token.address,
      tokenGrant.address,
      registry.address,
      initializationPeriod,
      contract.fromArtifact('TokenStakingEscrow'),
      contract.fromArtifact('TokenStaking')
    )
    stakingContract = stakingContracts.tokenStaking;

    undelegationPeriod = await stakingContract.undelegationPeriod()

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
      {from: owner}
    );
  }

  describe("delegate", async () => {
    it("should update balances", async () => {
      let ownerStartBalance = await token.balanceOf.call(owner);
  
      await delegate(operatorOne, stakingAmount);
      
      let ownerEndBalance = await token.balanceOf.call(owner);
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

    it("should not allow to delegate to the same operator after recovering stake", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      await stakingContract.undelegate(operatorOne, {from: owner});
      await time.increase(undelegationPeriod.addn(1));
      await stakingContract.recoverStake(operatorOne);
          
      await expectRevert(
        delegate(operatorOne, stakingAmount),
        "Operator undelegated"
      )
    })

    it("should not allow to delegate less than the minimum stake", async () => {    
      await expectRevert(
        delegate(operatorOne, minimumStake.subn(1)),
        "Less than the minimum stake"
      )
    })
  
    it("should allow to delegate the minimum stake", async () => {    
      await delegate(operatorOne, minimumStake)
      // ok, no reverts
    })
  
    it("should allow to delegate to two different operators", async () => {
      let ownerStartBalance = await token.balanceOf.call(owner)
  
      await delegate(operatorOne, stakingAmount);
      await delegate(operatorTwo, stakingAmount);
  
      let ownerEndBalance = await token.balanceOf.call(owner);
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
        "Not authorized"
      )
    })

    it("should allow to cancel delegation right away", async () => {
      await delegate(operatorOne, stakingAmount);
  
      await stakingContract.cancelStake(operatorOne, {from: owner});
      // ok, no revert
    })
  
    it("should allow to cancel delegation just before initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).sub(timeRoundMargin))
  
      await stakingContract.cancelStake(operatorOne, {from: owner})
      // ok, no revert
    })
  
    it("should not allow to cancel delegation after initialization period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
      await expectRevert(
        stakingContract.cancelStake(operatorOne, {from: owner}),
        "Initialized stake"
      );
    })

    it("should transfer tokens back to the owner", async () => {
      let ownerStartBalance = await token.balanceOf.call(owner);
  
      await delegate(operatorOne, stakingAmount);
  
      await stakingContract.cancelStake(operatorOne, {from: owner});
  
      let ownerEndBalance = await token.balanceOf.call(owner);
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

    it("should retain delegation info", async () => {
      await delegate(operatorOne, stakingAmount);
  
      let delegationInfoBefore = await stakingContract.getDelegationInfo.call(operatorOne)
  
      await stakingContract.cancelStake(operatorOne, {from: owner});
  
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
        "Not authorized"
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
        "Invalid timestamp"
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
        {from: owner}
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
        "Operator may not postpone undelegation"
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
        "Not authorized"
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
        "Invalid timestamp"
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
        "Invalid timestamp"
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
        {from: owner}
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
        "Operator may not postpone undelegation"
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
        "Not undelegated"
      )
    })

    it("should not allow to recover stake before undelegation period is over", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
      tx = await stakingContract.undelegate(operatorOne, {from: owner});
      let undelegatedAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp);
      await time.increaseTo(undelegatedAt.add(undelegationPeriod).sub(timeRoundMargin));
  
      await expectRevert(
        stakingContract.recoverStake(operatorOne),
        "Still undelegating"
      )
    })

    it("should transfer tokens back to the owner", async () => {
      let ownerStartBalance = await token.balanceOf.call(owner)

      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
  
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
      await stakingContract.undelegate(operatorOne, {from: owner});
      await time.increase(undelegationPeriod.addn(1));
      await stakingContract.recoverStake(operatorOne);
          
      let ownerEndBalance = await token.balanceOf.call(owner);
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

    it("should retain delegation info", async () => {
      let tx = await delegate(operatorOne, stakingAmount)
      let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
      await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
  
      let delegationInfoBefore = await stakingContract.getDelegationInfo.call(operatorOne)
      
      await stakingContract.undelegate(operatorOne, {from: owner})
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
  
      await stakingContract.cancelStake(operatorOne, {from: owner});
  
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
      await stakingContract.undelegate(operatorOne, {from: owner});
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
  
      await stakingContract.cancelStake(operatorOne, {from: owner})
  
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
      await stakingContract.undelegate(operatorOne, {from: owner})
  
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
        {from: owner}
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
        {from: owner}
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
