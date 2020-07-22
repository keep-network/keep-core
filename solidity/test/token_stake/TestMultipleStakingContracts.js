const {delegateStake, delegateStakeFromGrant} = require('../helpers/delegateStake')
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {time} = require("@openzeppelin/test-helpers")
const {initTokenStaking} = require('../helpers/initContracts')
const {grantTokens} = require('../helpers/grantTokens');
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot');

const KeepToken = contract.fromArtifact('KeepToken')
const KeepRegistry = contract.fromArtifact('KeepRegistry')
const TokenGrant = contract.fromArtifact('TokenGrant')
const PermissiveStakingPolicy = contract.fromArtifact('PermissiveStakingPolicy')

const BN = web3.utils.BN
const chai = require("chai")
chai.use(require("bn-chai")(BN))
const expect = chai.expect

describe('TokenStake/MultipleStakingContracts', () => {

  const deployer = accounts[0],
    grantManager = accounts[1],
    grantee = accounts[2],
    tokenOwner = accounts[3],
    operator = accounts[4],
    beneficiary = accounts[5],
    authorizer = accounts[6]

  const initializationPeriod = time.duration.hours(6),
    grantUnlockingDuration = time.duration.years(2),
    grantCliff = time.duration.seconds(0),
    grantRevocable = true

  let token, tokenGrant
  let oldTokenStaking, oldTokenStakingEscrow, newTokenStaking, newTokenStakingEscrow

  let grantId, grantedAmount, delegatedAmount

  before(async () => {
    //
    // Deploy KEEP token contract.
    // Transfer 50% of all tokens to grant manager and 10% of tokens
    // to token owner (account delegating liquid tokens in tests).
    //
    token = await KeepToken.new({from: deployer})
    const allTokens = await token.balanceOf(deployer)
    await token.transfer(grantManager, allTokens.divn(2), {from: deployer})
    await token.transfer(tokenOwner, allTokens.divn(10), {from: deployer})
    
    //
    // Deploy TokenGrant, KeepRegistry
    //
    const registry = await KeepRegistry.new({from: deployer})
    tokenGrant = await TokenGrant.new(token.address, {from: deployer})

    //
    // Deploy two instances of TokenStaking and TokenStakingEscrow contracts.
    // One of them will act as a previous version and another as a new version.
    //
    const oldStakingContracts = await initTokenStaking(
      token.address,
      tokenGrant.address,
      registry.address,
      initializationPeriod,
      contract.fromArtifact('TokenStakingEscrow'),
      contract.fromArtifact('TokenStaking')
    )
    oldTokenStaking = oldStakingContracts.tokenStaking
    oldTokenStakingEscrow = oldStakingContracts.tokenStakingEscrow

    const newStakingContracts = await initTokenStaking(
      token.address,
      tokenGrant.address,
      registry.address,
      initializationPeriod,
      contract.fromArtifact('TokenStakingEscrow'),
      contract.fromArtifact('TokenStaking')
    )
    newTokenStaking = newStakingContracts.tokenStaking
    newTokenStakingEscrow = newStakingContracts.tokenStakingEscrow

    await tokenGrant.authorizeStakingContract(oldTokenStaking.address, {
      from: grantManager,
    })
    await tokenGrant.authorizeStakingContract(newTokenStaking.address, {
      from: grantManager,
    })

    //
    // Grant tokens to grantee with a standard grant.
    //
    const minimumStake = await newTokenStaking.minimumStake()
    grantedAmount = minimumStake.muln(40)
    delegatedAmount = minimumStake.muln(10)
    const grantStart = await time.latest()
    const permissivePolicy = await PermissiveStakingPolicy.new()
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
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("when staking tokens with multiple staking contracts", async () => {
    it("should let to reuse operator address between contracts", async () => {
      await delegateStake(
        token, 
        oldTokenStaking, 
        tokenOwner,
        operator,
        beneficiary,
        authorizer,
        delegatedAmount
      )
      await delegateStake(
        token, 
        newTokenStaking, 
        tokenOwner,
        operator,
        beneficiary,
        authorizer,
        delegatedAmount
      )

      expect(await oldTokenStaking.balanceOf(operator)).to.eq.BN(delegatedAmount)
      expect(await newTokenStaking.balanceOf(operator)).to.eq.BN(delegatedAmount)
    })
  })

  describe("when staking granted tokens with multiple staking contracts", async () => {
    it("should let to reuse operator address between contracts", async () => {
      await delegateStakeFromGrant(
        tokenGrant,
        oldTokenStaking.address,
        grantee,
        operator,
        beneficiary,
        authorizer,
        delegatedAmount,
        grantId
      )
      await delegateStakeFromGrant(
        tokenGrant,
        newTokenStaking.address,
        grantee,
        operator,
        beneficiary,
        authorizer,
        delegatedAmount,
        grantId
      )

      expect(await oldTokenStaking.balanceOf(operator)).to.eq.BN(delegatedAmount)
      expect(await newTokenStaking.balanceOf(operator)).to.eq.BN(delegatedAmount)
    })
  })

  describe("when staking liquid and granted tokens with multiple staking contracts", async () => {
    it("should let to reuse operator address between contracts", async () => {
      await delegateStake(
        token, 
        oldTokenStaking, 
        tokenOwner,
        operator,
        beneficiary,
        authorizer,
        delegatedAmount
      )
      await delegateStakeFromGrant(
        tokenGrant,
        newTokenStaking.address,
        grantee,
        operator,
        beneficiary,
        authorizer,
        delegatedAmount,
        grantId
      )
  
      expect(await oldTokenStaking.balanceOf(operator)).to.eq.BN(delegatedAmount)
      expect(await newTokenStaking.balanceOf(operator)).to.eq.BN(delegatedAmount)  
    })

    it("should correctly identify the source of tokens", async () => {
      // Staking to the old contract from a grant, old staking contract captures
      // grant ID.
      await delegateStakeFromGrant(
        tokenGrant,
        oldTokenStaking.address,
        grantee,
        operator,
        beneficiary,
        authorizer,
        delegatedAmount,
        grantId
      )

      // Staking to a new contract using liquid tokens and the same operator.
      // The new contract should not capture grant ID even though the operator
      // is the same and TokenGrant has a TokenGrantStake for this operator.
      await delegateStake(
        token, 
        newTokenStaking, 
        tokenOwner,
        operator,
        beneficiary,
        authorizer,
        delegatedAmount
      )

      await time.increase(initializationPeriod.addn(1))
      await oldTokenStaking.undelegate(operator, {from: grantee})
      await newTokenStaking.undelegate(operator, {from: tokenOwner})
      
      const undelegationPeriod = await newTokenStaking.undelegationPeriod()
      await time.increase(undelegationPeriod.addn(1))

      const tokenOwnerBalanceBefore = await token.balanceOf(tokenOwner)
      const granteeBalanceBefore = await token.balanceOf(grantee)
      await oldTokenStaking.recoverStake(operator, {from: grantee})
      await newTokenStaking.recoverStake(operator, {from: tokenOwner})
      const tokenOwnerBalanceAfter = await token.balanceOf(tokenOwner)
      const granteeBalanceAfter = await token.balanceOf(grantee)

      // Grantee tokens should go to the escrow.
      // Token owner tokens should return back to him.
      const granteeReturned = granteeBalanceAfter.sub(granteeBalanceBefore)
      const tokenOwnerReturned = tokenOwnerBalanceAfter.sub(tokenOwnerBalanceBefore)
      const oldEscrowDeposited = await oldTokenStakingEscrow.depositedAmount(operator)
      const newEscrowDeposited = await newTokenStakingEscrow.depositedAmount(operator)

      expect(granteeReturned).to.eq.BN(0)
      expect(tokenOwnerReturned).to.eq.BN(delegatedAmount)
      expect(oldEscrowDeposited).to.eq.BN(delegatedAmount)
      expect(newEscrowDeposited).to.eq.BN(0)
    })
  })
})