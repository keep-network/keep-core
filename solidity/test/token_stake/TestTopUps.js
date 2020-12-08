const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const {initTokenStaking} = require('../helpers/initContracts')
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot');

const {grantTokens, grantTokensToManagedGrant} = require('../helpers/grantTokens');
const {
    delegateStake,
    delegateStakeFromGrant,
    delegateStakeFromManagedGrant,
  } = require('../helpers/delegateStake')

const KeepToken = contract.fromArtifact('KeepToken')
const KeepRegistry = contract.fromArtifact('KeepRegistry')
const TokenGrant = contract.fromArtifact('TokenGrant')
const PermissiveStakingPolicy = contract.fromArtifact('PermissiveStakingPolicy')
const ManagedGrantFactory = contract.fromArtifact('ManagedGrantFactory')
const ManagedGrant = contract.fromArtifact('ManagedGrant')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('TokenStaking/TopUps', () => {

  const deployer = accounts[0],
    grantManager = accounts[1],
    grantee = accounts[2],
    managedGrantee = accounts[3],
    tokenOwner = accounts[4],
    operatorOne = accounts[5],
    operatorTwo = accounts[6],
    operatorThree = accounts[7],
    operatorFour = accounts[8],
    beneficiary = accounts[9],
    authorizer = accounts[10],
    thirdParty = accounts[11],
    operatorContract = accounts[12]

  const initializationPeriod = time.duration.days(10),
    grantStart = time.duration.seconds(0),
    grantUnlockingDuration = time.duration.years(100),
    grantCliff = time.duration.seconds(0),
    grantRevocable = true

  let undelegationPeriod;

  let token, tokenGrant, tokenStakingEscrow, tokenStaking

  let grantId, grant2Id, managedGrant, grantedAmount, delegatedAmount

  before(async () => {
    //
    // Deploy KEEP token contract.
    // Transfer 50% of all tokens to grant manager and 50% of tokens
    // to token owner (account delegating liquid tokens in tests).
    //
    token = await KeepToken.new({from: deployer})
    const allTokens = await token.balanceOf(deployer)
    await token.transfer(grantManager, allTokens.divn(2), {from: deployer})
    await token.transfer(tokenOwner, allTokens.divn(2), {from: deployer})

    //
    // Deploy TokenGrant, ManagedGrantFactory, KeepRegistry, TokenStaking,
    // and TokenStakingEscrow. 
    // Authorize TokenStaking contract in TokenGrant contract.
    //
    tokenGrant = await TokenGrant.new(token.address, {from: deployer})
    const permissivePolicy = await PermissiveStakingPolicy.new()
    const managedGrantFactory = await ManagedGrantFactory.new(
      token.address,
      tokenGrant.address,
      {from: deployer}
    ) 
    const registry = await KeepRegistry.new({from: deployer})
    const stakingContracts = await initTokenStaking(
      token.address,
      tokenGrant.address,
      registry.address,
      initializationPeriod,
      contract.fromArtifact('TokenStakingEscrow'),
      contract.fromArtifact('TokenStaking')
    )
    tokenStaking = stakingContracts.tokenStaking
    tokenStakingEscrow = stakingContracts.tokenStakingEscrow
    await tokenGrant.authorizeStakingContract(tokenStaking.address, {
      from: grantManager,
    })

    undelegationPeriod = await tokenStaking.undelegationPeriod()

    //
    // Create three grants:
    // - two separate grants goes to grantee,
    // - one grant goes to managed grantee. 
    // 
    const minimumStake = await tokenStaking.minimumStake();
    grantedAmount = minimumStake.muln(40);
    delegatedAmount = minimumStake.muln(10);

    grantId = await grantTokens(
      tokenGrant, 
      token, 
      grantedAmount, 
      grantManager, 
      grantee, 
      grantUnlockingDuration,
      grantStart,
      grantCliff,
      grantRevocable,
      permissivePolicy.address
    )
    grant2Id = await grantTokens(
      tokenGrant, 
      token, 
      grantedAmount, 
      grantManager, 
      grantee, 
      grantUnlockingDuration,
      grantStart,
      grantCliff,
      grantRevocable,
      permissivePolicy.address
    )
    const managedGrantAddress = await grantTokensToManagedGrant(
      managedGrantFactory,
      token,
      grantedAmount,
      grantManager,
      managedGrantee,
      grantUnlockingDuration,
      grantStart,
      grantCliff,
      false,
      permissivePolicy.address,
    )
    managedGrant = await ManagedGrant.at(managedGrantAddress)

    //
    // Delegate stakes:
    // - operatorOne receives delegation from liquid tokens of tokenOwner
    // - operatorTwo receives delegation from granted tokens of grantee
    //   from the first grant
    // - operatorThree receives delegation from granted tokens of managed grantee.
    //
    await delegateStake(
      token, 
      tokenStaking, 
      tokenOwner,
      operatorOne,
      beneficiary,
      authorizer,
      delegatedAmount
    )
    await delegateStakeFromGrant(
      tokenGrant,
      tokenStaking.address,
      grantee,
      operatorTwo,
      beneficiary,
      authorizer,
      delegatedAmount,
      grantId
    )
    await delegateStakeFromManagedGrant(
      managedGrant,
      tokenStaking.address,
      managedGrantee,
      operatorThree,
      beneficiary,
      authorizer,
      delegatedAmount
    ) 

    //
    // Approve operator contract in the registry, authorize operator contract
    // for all three operators.
    //
    await registry.approveOperatorContract(
      operatorContract,
      {from: deployer}
    )
    await tokenStaking.authorizeOperatorContract(
      operatorOne,
      operatorContract,
      {from: authorizer}
    )
    await tokenStaking.authorizeOperatorContract(
      operatorTwo,
      operatorContract,
      {from: authorizer}
    )
    await tokenStaking.authorizeOperatorContract(
      operatorThree,
      operatorContract,
      {from: authorizer}
    )
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("delegated liquid tokens top-ups", async () => {
    async function initiateTopUp(value) {
      return delegateStake(
        token, tokenStaking, tokenOwner, operatorOne,
        beneficiary, authorizer, value
      )
    }

    it("can not be done when stake is undelegating", async () => {
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.undelegate(operatorOne, {from: tokenOwner})
      time.increase(1) // we need the block timestamp to increase
      await expectRevert(
        initiateTopUp(delegatedAmount),
        "Stake undelegated"
      ) 
    })

    it("can not be done when stake is recovered", async () => {
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.undelegate(operatorOne, {from: tokenOwner})
      await time.increase(undelegationPeriod.addn(1))
      await tokenStaking.recoverStake(operatorOne)
      time.increase(1) // we need the block timestamp to increase
      await expectRevert(
        initiateTopUp(delegatedAmount),
        "Stake undelegated"
      )
    })

    it("can not be done by another token owner", async () => {
      await token.transfer(thirdParty, delegatedAmount, {from: tokenOwner})
      await time.increase(initializationPeriod.addn(1))
      await expectRevert(
        delegateStake(
          token, tokenStaking, thirdParty, operatorOne,
          beneficiary, authorizer, delegatedAmount
        ),
        "Not the same owner"
      )
    })

    it("can not be done from a token grant", async () => {
      await time.increase(initializationPeriod.addn(1))
      await expectRevert(
        delegateStakeFromGrant(
          tokenGrant,
          tokenStaking.address,
          grantee,
          operatorOne,
          beneficiary,
          authorizer,
          delegatedAmount,
          grantId
        ),
        "Must not be from a grant"
      )   
    })

    it("can be done in one step during initialization period", async () => {
      // half of the initialization period passed
      await time.increase(initializationPeriod.divn(2))

      await initiateTopUp(delegatedAmount)

      const delegationInfo = await tokenStaking.getDelegationInfo(operatorOne)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount.muln(2))
      expect(delegationInfo.createdAt).to.eq.BN(await time.latest())
      expect(delegationInfo.undelegatedAt).to.eq.BN(0)
    })

    it("fails to commit when done in one step during initialization period", async () => {
      // half of the initialization period passed
      await time.increase(initializationPeriod.divn(2))

      // We are still in the initialization period, top-up is done in one step
      // and there is nothing to commit.
      await initiateTopUp(delegatedAmount) 
      await expectRevert(
        tokenStaking.commitTopUp(operatorOne),
        "No top up to commit"
      )
    })

    it("fails for a zero-value top-up when doing in one step", async () => {
      await expectRevert(
        initiateTopUp(0),
        "Top-up value must be greater than zero"
      )
    })

    it("fails for a zero-value top-up when doing in two steps", async () => {
      await time.increase(initializationPeriod.addn(1))
      await expectRevert(
        initiateTopUp(0),
        "Top-up value must be greater than zero"
      )
    })

    it("does not increase stake before committed", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod)

      let currentStake = await tokenStaking.activeStake(
        operatorOne,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount)
    })

    it("increases stake once committed", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorOne)

      const currentStake = await tokenStaking.activeStake(
        operatorOne,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(2))
    })

    it("can not be committed before initialization period is over", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await expectRevert(
        tokenStaking.commitTopUp(operatorOne),
        "Stake is initializing"
      )
    })

    it("returns full amount with committed top-ups to owner after recovering stake", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorOne)
      await tokenStaking.undelegate(operatorOne, {from: tokenOwner})
      await time.increase(undelegationPeriod.addn(1))

      const before = await token.balanceOf(tokenOwner) 
      await tokenStaking.recoverStake(operatorOne, {from: tokenOwner})
      const after = await token.balanceOf(tokenOwner)

      expect(after.sub(before)).to.eq.BN(delegatedAmount.muln(2))
    })

    it("returns uncommitted top-ups to owner after recovering stake", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await tokenStaking.undelegate(operatorOne, {from: tokenOwner})
      await time.increase(undelegationPeriod.addn(1))

      const before = await token.balanceOf(tokenOwner) 
      await tokenStaking.recoverStake(operatorOne, {from: tokenOwner})
      const after = await token.balanceOf(tokenOwner)

      expect(after.sub(before)).to.eq.BN(delegatedAmount.muln(2))
    })

    it("fails to commit if not first initialized", async () => {
      await time.increase(initializationPeriod.addn(1))
      await expectRevert(
        tokenStaking.commitTopUp(operatorOne),
        "No top up to commit"
      )
    })
    
    it("can be done multiple times", async () => {
      // The first top-up.
      // Half of the initialization period passed, top-up should  be processed
      // immediately but also reset the initialization period.
      await time.increase(initializationPeriod.divn(2))
      await initiateTopUp(delegatedAmount)
    
      await time.increase(initializationPeriod.addn(1))
      let currentStake = await tokenStaking.activeStake(
        operatorOne,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(2))

      // The second top-up.
      // Initialization period passed so it has to be done in two steps.
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorOne)

      currentStake = await tokenStaking.activeStake(
        operatorOne,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(3))


      // And yet one top-up...
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorOne)

      currentStake = await tokenStaking.activeStake(
        operatorOne,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(4))
    })

    it("can be done with less than the minimum stake", async () => {
      const topUpAmount = (await tokenStaking.minimumStake()).subn(1)

      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(topUpAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorOne)

      const currentStake = await tokenStaking.activeStake(
        operatorOne,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.add(topUpAmount))
    })
  })

  describe("delegated grant top-ups", async () => {
    async function initiateTopUp(value) {
      await delegateStakeFromGrant(
        tokenGrant,
        tokenStaking.address,
        grantee,
        operatorTwo,
        beneficiary,
        authorizer,
        value,
        grantId
      ) 
    }

    it("can not be done when stake is undelegated", async () => {
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.undelegate(operatorTwo, {from: grantee})
      time.increase(1) // we need the block timestamp to increase
      await expectRevert(
        initiateTopUp(delegatedAmount),
        "Stake undelegated"
      )
    })

    it("can not be done when stake is recovered", async () => {
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.undelegate(operatorTwo, {from: grantee})
      await time.increase(undelegationPeriod.addn(1))
      await tokenStaking.recoverStake(operatorTwo)
      time.increase(1) // we need the block timestamp to increase
      await expectRevert(
        initiateTopUp(delegatedAmount),
        "Stake undelegated"
      )
    })

    it("can not be done from another grant", async () => {
      await time.increase(initializationPeriod.addn(1))

      await expectRevert(
        delegateStakeFromGrant(
          tokenGrant, tokenStaking.address, grantee, operatorTwo,
          beneficiary, authorizer, delegatedAmount, grant2Id
        ),
        "Not the same grant"
      ) 
    })

    it("can not be done from liquid tokens", async () => {
      await time.increase(initializationPeriod.addn(1))

      await token.transfer(grantee, grantedAmount, {from: grantManager})
      await expectRevert(
        delegateStake(
          token, tokenStaking, grantee, operatorTwo,
          beneficiary, authorizer, delegatedAmount
        ),
        "Must be from a grant"
      );
    })

    it("can be done in one step during initialization period", async () => {
      // half of the initialization period passed
      await time.increase(initializationPeriod.divn(2))

      await initiateTopUp(delegatedAmount)

      const delegationInfo = await tokenStaking.getDelegationInfo(operatorTwo)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount.muln(2))
      expect(delegationInfo.createdAt).to.eq.BN(await time.latest())
      expect(delegationInfo.undelegatedAt).to.eq.BN(0)
    })

    it("fails for a zero-value top-up when doing in one step", async () => {
      await expectRevert(
        initiateTopUp(0),
        "Top-up value must be greater than zero"
      )
    })

    it("fails for a zero-value top-up when doing in two steps", async () => {
      await time.increase(initializationPeriod.addn(1))
      await expectRevert(
        initiateTopUp(0),
        "Top-up value must be greater than zero"
      )
    })

    it("does not increase stake before committed", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod)

      let currentStake = await tokenStaking.activeStake(
        operatorTwo,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount)
    })

    it("increases stake once committed", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorTwo)

      const currentStake = await tokenStaking.activeStake(
        operatorTwo,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(2))
    })

    it("can not be committed before initialization period is over", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await expectRevert(
        tokenStaking.commitTopUp(operatorTwo),
        "Stake is initializing"
      )
    })

    it("returns full amount with committed top-ups to escrow after recovering stake", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorTwo)
      await tokenStaking.undelegate(operatorTwo, {from: grantee})
      await time.increase(undelegationPeriod.addn(1))

      const before = await tokenStakingEscrow.depositedAmount(operatorTwo)
      await tokenStaking.recoverStake(operatorTwo, {from: grantee})
      const after = await tokenStakingEscrow.depositedAmount(operatorTwo)

      expect(after.sub(before)).to.eq.BN(delegatedAmount.muln(2))  
    })

    it("returns uncommitted top-ups to escrow after recovering stake", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await tokenStaking.undelegate(operatorTwo, {from: grantee})
      await time.increase(undelegationPeriod.addn(1))

      const before = await tokenStakingEscrow.depositedAmount(operatorTwo)
      await tokenStaking.recoverStake(operatorTwo, {from: grantee})
      const after = await tokenStakingEscrow.depositedAmount(operatorTwo)

      expect(after.sub(before)).to.eq.BN(delegatedAmount.muln(2))  
    })

    it("fails to commit if not first initialized", async () => {
      await time.increase(initializationPeriod.addn(1))
      await expectRevert(
        tokenStaking.commitTopUp(operatorTwo),
        "No top up to commit"
      )
    })
    
    it("can be done multiple times", async () => {
      // The first top-up.
      // Half of the initialization period passed, top-up should  be processed
      // immediately but also reset the initialization period.
      await time.increase(initializationPeriod.divn(2))
      await initiateTopUp(delegatedAmount)
    
      await time.increase(initializationPeriod.addn(1))
      let currentStake = await tokenStaking.activeStake(
        operatorTwo,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(2))

      // The second top-up.
      // Initialization period passed so it has to be done in two steps.
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorTwo)

      currentStake = await tokenStaking.activeStake(
        operatorTwo,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(3))


      // And yet one top-up...
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorTwo)

      currentStake = await tokenStaking.activeStake(
        operatorTwo,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(4))
    })

    it("can be done with less than the minimum stake", async () => {
      const topUpAmount = (await tokenStaking.minimumStake()).subn(1)

      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(topUpAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorTwo)

      const currentStake = await tokenStaking.activeStake(
        operatorTwo,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.add(topUpAmount))
    })
  })

  describe("delegated managed grant top-ups", async () => {
    async function initiateTopUp(value) {
      await delegateStakeFromManagedGrant(
        managedGrant,
        tokenStaking.address,
        managedGrantee,
        operatorThree,
        beneficiary,
        authorizer,
        value
      ) 
    }

    it("can not be done when stake is undelegating", async () => {
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.undelegate(operatorThree, {from: managedGrantee})
      time.increase(1) // we need the block timestamp to increase
      await expectRevert(
        initiateTopUp(delegatedAmount),
        "Stake undelegated"
      )
    })

    it("can not be done when stake is recovered", async () => {
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.undelegate(operatorThree, {from: managedGrantee})
      await time.increase(undelegationPeriod.addn(1))
      await tokenStaking.recoverStake(operatorThree)
      time.increase(1) // we need the block timestamp to increase
      await expectRevert(
        initiateTopUp(delegatedAmount),
        "Stake undelegated"
      )
    })

    it("can not be done from another grant", async () => {
      await time.increase(initializationPeriod.addn(1))

      await expectRevert(
        delegateStakeFromGrant(
          tokenGrant, tokenStaking.address, grantee, operatorThree,
          beneficiary, authorizer, delegatedAmount, grantId
        ),
        "Not the same grant"
      )   
    })

    it("can not be done from liquid tokens", async () => {
      await time.increase(initializationPeriod.addn(1))

      await token.transfer(managedGrantee, grantedAmount, {from: grantManager})
      await expectRevert(
        delegateStake(
          token, tokenStaking, managedGrantee, operatorTwo,
          beneficiary, authorizer, delegatedAmount
        ),
        "Must be from a grant"
      );
    })

    it("can be done in one step during initialization period", async () => {
      // half of the initialization period passed
      await time.increase(initializationPeriod.divn(2))

      await initiateTopUp(delegatedAmount)

      const delegationInfo = await tokenStaking.getDelegationInfo(operatorThree)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount.muln(2))
      expect(delegationInfo.createdAt).to.eq.BN(await time.latest())
      expect(delegationInfo.undelegatedAt).to.eq.BN(0)
    })

    it("fails for a zero-value top-up when doing in one step", async () => {
      await expectRevert(
        initiateTopUp(0),
        "Top-up value must be greater than zero"
      )
    })

    it("fails for a zero-value top-up when doing in two steps", async () => {
      await time.increase(initializationPeriod.addn(1))
      await expectRevert(
        initiateTopUp(0),
        "Top-up value must be greater than zero"
      )
    })

    it("does not increase stake before committed", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod)

      let currentStake = await tokenStaking.activeStake(
        operatorThree,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount)
    })

    it("increases stake once committed", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorThree)

      const currentStake = await tokenStaking.activeStake(
        operatorThree,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(2))
    })

    it("can not be committed before initialization period is over", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await expectRevert(
        tokenStaking.commitTopUp(operatorThree),
        "Stake is initializing"
      )
    })

    it("returns full amount with committed top-ups to escrow after recovering stake", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorThree)
      await tokenStaking.undelegate(operatorThree, {from: managedGrantee})
      await time.increase(undelegationPeriod.addn(1))

      const before = await tokenStakingEscrow.depositedAmount(operatorThree)
      await tokenStaking.recoverStake(operatorThree, {from: managedGrantee})
      const after = await tokenStakingEscrow.depositedAmount(operatorThree)

      expect(after.sub(before)).to.eq.BN(delegatedAmount.muln(2))  
    })

    it("returns uncommitted top-ups to escrow after recovering stake", async () => {
      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(delegatedAmount)
      await tokenStaking.undelegate(operatorThree, {from: managedGrantee})
      await time.increase(undelegationPeriod.addn(1))

      const before = await tokenStakingEscrow.depositedAmount(operatorThree)
      await tokenStaking.recoverStake(operatorThree, {from: managedGrantee})
      const after = await tokenStakingEscrow.depositedAmount(operatorThree)

      expect(after.sub(before)).to.eq.BN(delegatedAmount.muln(2))
    })

    it("fails to commit if not first initialized", async () => {
      await time.increase(initializationPeriod.addn(1))
      await expectRevert(
        tokenStaking.commitTopUp(operatorThree),
        "No top up to commit"
      )
    })
    
    it("can be done multiple times", async () => {
      // The first top-up.
      // Half of the initialization period passed, top-up should  be processed
      // immediately but also reset the initialization period.
      await time.increase(initializationPeriod.divn(2))
      await initiateTopUp(delegatedAmount)
    
      await time.increase(initializationPeriod.addn(1))
      let currentStake = await tokenStaking.activeStake(
        operatorThree,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(2))

      // The second top-up.
      // Initialization period passed so it has to be done in two steps.
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorThree)

      currentStake = await tokenStaking.activeStake(
        operatorThree,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(3))


      // And yet one top-up...
      await initiateTopUp(delegatedAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorThree)

      currentStake = await tokenStaking.activeStake(
        operatorThree,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.muln(4))
    })

    it("can be done with less than the minimum stake", async () => {
      const topUpAmount = (await tokenStaking.minimumStake()).subn(1)

      await time.increase(initializationPeriod.addn(1))
      await initiateTopUp(topUpAmount)
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.commitTopUp(operatorThree)

      const currentStake = await tokenStaking.activeStake(
        operatorThree,
        operatorContract
      )
      expect(currentStake).to.eq.BN(delegatedAmount.add(topUpAmount))
    })
  })

  describe("escrow redelegation top-ups", async () => {
    beforeEach(async () => {
      await time.increase(initializationPeriod.addn(1))
      await tokenStaking.undelegate(operatorTwo, {from: grantee})
      await tokenStaking.undelegate(operatorThree, {from: managedGrantee})
      await time.increase(undelegationPeriod.addn(1))
      await tokenStaking.recoverStake(operatorTwo)
      await tokenStaking.recoverStake(operatorThree)
    })

    it("can be done for a grant", async () => {
      await delegateStakeFromGrant(
        tokenGrant,
        tokenStaking.address,
        grantee,
        operatorFour,
        beneficiary,
        authorizer,
        delegatedAmount,
        grantId
      )
      await time.increase(initializationPeriod.addn(1))

      let delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount)

      const data = Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operatorFour.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ])
      await tokenStakingEscrow.redelegate(
        operatorTwo, delegatedAmount, data, {from: grantee}
      )

      await time.increase(initializationPeriod.addn(1))
      delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount)
  
      await tokenStaking.commitTopUp(operatorFour, {from: grantee})
        
      delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount.muln(2))
    })

    it("can be done for a managed grant", async () => {
      await delegateStakeFromManagedGrant(
        managedGrant,
        tokenStaking.address,
        managedGrantee,
        operatorFour,
        beneficiary,
        authorizer,
        delegatedAmount
      ) 
      await time.increase(initializationPeriod.addn(1))

      let delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount)

      const data = Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operatorFour.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ])
      await tokenStakingEscrow.redelegate(
        operatorThree, delegatedAmount, data, {from: managedGrantee}
      )

      await time.increase(initializationPeriod.addn(1))
      delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount)
  
      await tokenStaking.commitTopUp(operatorFour, {from: managedGrantee})
        
      delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount.muln(2))
    })

    it("can be done with less than the minimum stake", async () => {
      await delegateStakeFromGrant(
        tokenGrant,
        tokenStaking.address,
        grantee,
        operatorFour,
        beneficiary,
        authorizer,
        delegatedAmount,
        grantId
      )

      await time.increase(initializationPeriod.addn(1))

      let delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount)

      const topUpAmount = (await tokenStaking.minimumStake()).subn(1)
      const data = Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operatorFour.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ])
      await tokenStakingEscrow.redelegate(
        operatorTwo, topUpAmount, data, {from: grantee}
      )

      await time.increase(initializationPeriod.addn(1))  
      await tokenStaking.commitTopUp(operatorFour, {from: grantee})
        
      delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
      expect(delegationInfo.amount).to.eq.BN(delegatedAmount.add(topUpAmount))
    })
  })
}) 
