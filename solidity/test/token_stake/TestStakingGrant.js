const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, expectEvent, time} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot');

const {grantTokens, grantTokensToManagedGrant} = require('../helpers/grantTokens');
const {
    delegateStakeFromGrant,
    delegateStakeFromManagedGrant,
  } = require('../helpers/delegateStake')

const KeepToken = contract.fromArtifact('KeepToken')
const KeepRegistry = contract.fromArtifact('KeepRegistry')
const MinimumStakeSchedule = contract.fromArtifact('MinimumStakeSchedule')
const GrantStaking = contract.fromArtifact('GrantStaking')
const TokenStaking = contract.fromArtifact('TokenStaking')
const TokenGrant = contract.fromArtifact('TokenGrant')
const PermissiveStakingPolicy = contract.fromArtifact('PermissiveStakingPolicy')
const ManagedGrantFactory = contract.fromArtifact('ManagedGrantFactory')
const ManagedGrant = contract.fromArtifact('ManagedGrant')
const TokenStakingEscrow = contract.fromArtifact('TokenStakingEscrow')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('TokenStaking/StakingGrant', () => {

    const deployer = accounts[0],
    grantManager = accounts[1],
    grantee = accounts[2],
    managedGrantee = accounts[3],
    operatorOne = accounts[4],
    operatorTwo = accounts[5],
    beneficiary = accounts[6],
    authorizer = accounts[7],
    thirdParty = accounts[8]

    const initializationPeriod = time.duration.seconds(10),
      undelegationPeriod = time.duration.seconds(10),
      grantUnlockingDuration = time.duration.seconds(0),
      grantStart = time.duration.seconds(0),
      grantCliff = time.duration.seconds(0),
      grantRevocable = true

    let token, tokenGrant, tokenStakingEscrow, tokenStaking

    let grantId, managedGrantId, managedGrant, grantedAmount, delegatedAmount

    before(async () => {
      const registry = await KeepRegistry.new({from: deployer})
      token = await KeepToken.new({from: deployer})
      const allTokens = await token.balanceOf(deployer)
      await token.transfer(grantManager, allTokens, {from: deployer})
  
      tokenGrant = await TokenGrant.new(token.address, {from: deployer})
      const permissivePolicy = await PermissiveStakingPolicy.new()
      const managedGrantFactory = await ManagedGrantFactory.new(
        token.address,
        tokenGrant.address,
        {from: deployer}
      )
      
      tokenStakingEscrow = await TokenStakingEscrow.new(
        token.address, 
        tokenGrant.address,
        {from: deployer}
      )

      await TokenStaking.detectNetwork()
      await TokenStaking.link(
        'MinimumStakeSchedule', 
        (await MinimumStakeSchedule.new({from: deployer})).address
      )
      await TokenStaking.link(
        'GrantStaking', 
        (await GrantStaking.new({from: deployer})).address
      )
      tokenStaking = await TokenStaking.new(
        token.address,
        tokenGrant.address,
        tokenStakingEscrow.address,
        registry.address,
        initializationPeriod,
        undelegationPeriod,
        {from: deployer}
      )
      await tokenStakingEscrow.transferOwnership(
        tokenStaking.address, 
        {from: deployer}
      )
      await tokenGrant.authorizeStakingContract(tokenStaking.address, {
        from: grantManager,
      })
      
      const minimumStake = await tokenStaking.minimumStake();
      grantedAmount = minimumStake.muln(40);
      delegatedAmount = minimumStake.muln(20);

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
      await delegateStakeFromGrant(
        tokenGrant,
        tokenStaking.address,
        grantee,
        operatorOne,
        beneficiary,
        authorizer,
        delegatedAmount,
        grantId
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
        grantRevocable,
        permissivePolicy.address,
      )
      managedGrant = await ManagedGrant.at(managedGrantAddress)
      managedGrantId = (await managedGrant.grantId()).toNumber()
      await delegateStakeFromManagedGrant(
        managedGrant,
        tokenStaking.address,
        managedGrantee,
        operatorTwo,
        beneficiary,
        authorizer,
        delegatedAmount
      )
    })

    beforeEach(async () => {
      await createSnapshot()
    })
    
    afterEach(async () => {
      await restoreSnapshot()
    })

    describe('cancelStake', async () => {
      it('should let operator cancel delegation', async () => {
        await tokenStaking.cancelStake(operatorOne, {from: operatorOne})
        await tokenStaking.cancelStake(operatorTwo, {from: operatorTwo})
        // ok, no revert
      })

      it('should let grantee cancel delegation', async () => {
        await tokenStaking.cancelStake(operatorOne, {from: grantee})
        // ok, no revert
      })

      it('should let managed grantee cancel delegation', async () => {
        await tokenStaking.cancelStake(operatorTwo, {from: managedGrantee})
        // ok, no revert
      })

      it('should let grantee cancel delegation via TokenGrant', async () => {
        await tokenGrant.cancelStake(operatorOne, {from: grantee})
        // ok, no revert
      })

      it('should let managed grantee cancel delegation via ManagedGrant', async () => {
        await managedGrant.cancelStake(operatorTwo, {from: managedGrantee})
        // ok, no revert
      })

      it('should not let operator cancel delegation for another operator', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorOne, {from: operatorTwo}),
          'Unauthorized'
        )
      })

      it('should not let grantee cancel delegation of another grantee', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorTwo, {from: grantee}),
          'Unauthorized'
        ) 
      })

      it('should not let managed grantee cancel delegation of another grantee', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorOne, {from: managedGrantee}),
          'Unauthorized'
        )     
      })

      it('should not let third party cancel delegation', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorOne, {from: thirdParty}),
          'Unauthorized'
        )
        await expectRevert(
          tokenStaking.cancelStake(operatorTwo, {from: thirdParty}),
          'Unauthorized'
        )    
      })

      it('should not let grant manager cancel delegation of non-revoked grant', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorOne, {from: grantManager}),
          'Unauthorized'
        )  
      })

      it('should let grant manager cancel delegation of revoked grant', async () => {
        await tokenGrant.revoke(grantId, {from: grantManager})
        await tokenGrant.revoke(managedGrantId, {from: grantManager})

        await tokenStaking.cancelStake(operatorOne, {from: grantManager})
        // ok, no revert

        await tokenStaking.cancelStake(operatorTwo, {from: grantManager})
        // ok, no revert
      })
    })

    describe('undelegate', async () => {
      before(async () => {
        await time.increase(initializationPeriod)
      })

      it('should let operator undelegate', async () => {
        await tokenStaking.undelegate(operatorOne, {from: operatorOne})
        await tokenStaking.undelegate(operatorTwo, {from: operatorTwo})
        // ok, no revert    
      })

      it('should let grantee undelegate', async () => {
        await tokenStaking.undelegate(operatorOne, {from: grantee})
        // ok, no revert  
      })

      it('should let managed grantee undelegate', async () => {
        await tokenStaking.undelegate(operatorTwo, {from: managedGrantee})
        // ok, no revert  
      })

      it('should let grantee undelegate via TokenGrant', async () => {
        await tokenGrant.undelegate(operatorOne, {from: grantee})
        // ok, no revert  
      })

      it('should let managed grantee undelegate via ManagedGrant', async () => {
        await managedGrant.undelegate(operatorTwo, {from: managedGrantee})
        // ok, no revert  
      })

      it('should not let operator undelegate for another operator', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorOne, {from: operatorTwo}),
          'Unauthorized'
        )   
      })

      it('should not let grantee undelegate for another grantee', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorTwo, {from: grantee}),
          'Unauthorized'
        )     
      })

      it('should not let managed grantee undelegate for another grantee', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorOne, {from: managedGrantee}),
          'Unauthorized'
        )    
      })

      it('should not let third party undelegate', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorOne, {from: thirdParty}),
          'Unauthorized'
        )
        await expectRevert(
          tokenStaking.undelegate(operatorTwo, {from: thirdParty}),
          'Unauthorized'
        )    
      })

      it('should not let grant manager undelegate non-revoked grant', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorOne, {from: grantManager}),
          'Unauthorized'
        ) 
      })

      it('should let grant manager undelegate revoked grant', async () => {
        await tokenGrant.revoke(grantId, {from: grantManager})
        await tokenGrant.revoke(managedGrantId, {from: grantManager})

        await tokenStaking.undelegate(operatorOne, {from: grantManager})
        // ok, no revert

        await tokenStaking.undelegate(operatorTwo, {from: grantManager})
        // ok, no revert
      })
    })

    describe('recoverStake', async () => {
      before(async () => {
        await time.increase(initializationPeriod)
        await tokenStaking.undelegate(operatorOne, {from: operatorOne})
        await tokenStaking.undelegate(operatorTwo, {from: operatorTwo})
        await time.increase(undelegationPeriod.addn(1))
      })

      it('transfers granted and undelegated tokens to escrow', async () => {
        await tokenStaking.recoverStake(operatorOne, {from: thirdParty})
        await tokenStaking.recoverStake(operatorTwo, {from: thirdParty})

        let deposited = await tokenStakingEscrow.depositedAmount(operatorOne)
        expect(deposited).to.eq.BN(delegatedAmount)

        deposited = await tokenStakingEscrow.depositedAmount(operatorTwo)
        expect(deposited).to.eq.BN(delegatedAmount)
      })
    })
})