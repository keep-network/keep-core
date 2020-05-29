const { accounts, contract } = require('@openzeppelin/test-environment')
const { time, expectRevert } = require("@openzeppelin/test-helpers")
const grantTokens = require('./helpers/grantTokens')
const { delegateStake, delegateStakeFromGrant, delegateStakeFromManagedGrant } = require('./helpers/delegateStake')
const { createSnapshot, restoreSnapshot } = require('./helpers/snapshot')
const assert = require('chai').assert

const KeepToken = contract.fromArtifact('KeepToken')
const TokenStaking = contract.fromArtifact('TokenStaking')
const TokenGrant = contract.fromArtifact('TokenGrant')
const KeepRegistry = contract.fromArtifact('KeepRegistry')
const PermissiveStakingPolicy = contract.fromArtifact('PermissiveStakingPolicy')
const ManagedGrantFactory = contract.fromArtifact('ManagedGrantFactory')
const ManagedGrant = contract.fromArtifact('ManagedGrant')
const RolesLookup = contract.fromArtifact('RolesLookup')
const RolesLookupStub = contract.fromArtifact('RolesLookupStub')

describe('RolesLookup', () => {
    
  const deployer = accounts[0],
    tokenOwner = accounts[1],
    nonTokenOwner = accounts[2],
    operator = accounts[3],
    nonOperator = accounts[4],
    beneficiary = accounts[5],
    authorizer = accounts[6],
    grantee = accounts[7]
    nonGrantee = accounts[8]
    
  const initializationPeriod = time.duration.seconds(0),
    undelegationPeriod = time.duration.seconds(0),
    grantUnlockingDuration = time.duration.seconds(0),
    grantStart = time.duration.seconds(0),
    grantCliff = time.duration.seconds(0),
    grantRevocable = true

  let token,
    tokenGrant,
    tokenStaking, 
    tokenGrantStakingPolicy,
    managedGrantFactory,
    lookup

  before(async () => {
    const registry = await KeepRegistry.new({from: deployer})
    token = await KeepToken.new({from: deployer})
    tokenGrant = await TokenGrant.new(token.address, {from: deployer})
    tokenStaking  = await TokenStaking.new(
      token.address, 
      registry.address, 
      initializationPeriod, 
      undelegationPeriod, 
      {from: deployer}
    )
    tokenGrantStakingPolicy = await PermissiveStakingPolicy.new()
    managedGrantFactory = await ManagedGrantFactory.new(
      token.address,
      tokenGrant.address,
      {from: deployer}
    )

    await tokenGrant.authorizeStakingContract(
      tokenStaking.address, {from: deployer}
    )

    const lookupLib = await RolesLookup.new({from: deployer})
    await RolesLookupStub.detectNetwork()
    await RolesLookupStub.link("RolesLookup", lookupLib.address)
    lookup = await RolesLookupStub.new(
      tokenStaking.address,
      tokenGrant.address,
      {from: deployer}
    )
  })

  describe("isTokenOwnerForOperator", async () => {
    before(async () => {
      await createSnapshot()
      const amount = await tokenStaking.minimumStake()
      await token.transfer(tokenOwner, amount, {from: deployer})
      await delegateStake(
        token, 
        tokenStaking, 
        tokenOwner, 
        operator, 
        beneficiary,
        authorizer, 
        amount,
        {from: tokenOwner}
      )
    })

    after(async () => {
      await restoreSnapshot()
    })

    it("returns true for token owner and its operator", async () => {
      assert.isTrue(await lookup.isTokenOwnerForOperator(
        tokenOwner, operator
      ))
    })

    it("returns false for incorrect operator", async () => {
      assert.isFalse(await lookup.isTokenOwnerForOperator(
        tokenOwner, nonOperator
      ))
    })

    it("returns false for non-token-owner", async () => {
      assert.isFalse(await lookup.isTokenOwnerForOperator(
        nonTokenOwner, operator
      ))  
    })
  })

  describe("isGranteeForOperator", async () => {
    before(async () => {
      await createSnapshot()
      const amount = await tokenStaking.minimumStake()
      let grantId = await grantTokens(
        tokenGrant,
        token,
        amount,
        deployer,
        grantee,
        grantUnlockingDuration,
        grantStart,
        grantCliff,
        grantRevocable,
        tokenGrantStakingPolicy.address,
        {from: deployer}
      )
      await delegateStakeFromGrant(
        tokenGrant,
        tokenStaking.address,
        grantee,
        operator,
        beneficiary,
        authorizer,
        amount,
        grantId
      )
    })

    after(async () => {
      await restoreSnapshot()
    })

    it("returns true for grantee and its operator", async () => {
      assert.isTrue(await lookup.isGranteeForOperator(
        grantee, operator
      ))
    })

    it("returns false for incorrect operator", async () => {
      assert.isFalse(await lookup.isGranteeForOperator(
        grantee, nonOperator
      ))
    })

    it("returns false for non-grantee", async () => {
      assert.isFalse(await lookup.isGranteeForOperator(
        nonGrantee, operator
      ))
    })
  })

  describe("isManagedGranteeForOperator", async () => {
    let managedGrantAddress

    before(async () => {
      await createSnapshot()
      const amount = await tokenStaking.minimumStake()
      await token.approve(
        managedGrantFactory.address, amount, {from: deployer}
      )
      managedGrantAddress = await managedGrantFactory.createManagedGrant.call(
        grantee,
        amount,
        grantUnlockingDuration,
        grantStart,
        grantCliff,
        grantRevocable,
        tokenGrantStakingPolicy.address,
        {from: deployer}
      )
      await managedGrantFactory.createManagedGrant(
        grantee,
        amount,
        grantUnlockingDuration,
        grantStart,
        grantCliff,
        grantRevocable,
        tokenGrantStakingPolicy.address,
        {from: deployer}
      )
      const managedGrant = await ManagedGrant.at(managedGrantAddress)
      await delegateStakeFromManagedGrant(
        managedGrant,
        tokenStaking.address,
        grantee,
        operator,
        beneficiary,
        authorizer,
        amount
      )
    })

    after(async () => {
      await restoreSnapshot()
    })

    it("returns true for grantee and its operator", async () => {
      assert.isTrue(await lookup.isManagedGranteeForOperator(
        grantee, operator, managedGrantAddress
      ))
    })

    it("returns false for incorrect operator", async () => {
      assert.isFalse(await lookup.isManagedGranteeForOperator(
        grantee, nonOperator, managedGrantAddress
      ))
    })

    it("reverts for non-grantee", async () => {
      await expectRevert(
        lookup.isManagedGranteeForOperator(
          nonGrantee, nonOperator, managedGrantAddress
        ),
        "Not a grantee of the provided contract"   
      )
    })
  })
})