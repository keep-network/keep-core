const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, expectEvent, time} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot')
const {initTokenStaking} = require('../helpers/initContracts')

const {grantTokens, grantTokensToManagedGrant} = require('../helpers/grantTokens');
const {
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

describe('TokenStaking/StakingGrant', () => {

    const deployer = accounts[0],
    grantManager = accounts[1],
    grantee = accounts[2],
    managedGrantee = accounts[3],
    operatorOne = accounts[4],
    operatorTwo = accounts[5],
    operatorThree = accounts[6],
    operatorFour = accounts[7],
    beneficiary = accounts[8],
    authorizer = accounts[9],
    thirdParty = accounts[10]

    const initializationPeriod = time.duration.hours(6),
      grantUnlockingDuration = time.duration.years(2),
      grantCliff = time.duration.seconds(0),
      grantRevocable = true

    let undelegationPeriod

    let token, tokenGrant, tokenStakingEscrow, tokenStaking

    let grantStart, grantId, managedGrantId, managedGrant, 
      grantedAmount, delegatedAmount

    before(async () => {
      grantStart = await time.latest()

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
      
      const stakingContracts = await initTokenStaking(
        token.address,
        tokenGrant.address,
        registry.address,
        initializationPeriod,
        contract.fromArtifact('TokenStakingEscrow'),
        contract.fromArtifact('TokenStakingStub')
      )
      tokenStaking = stakingContracts.tokenStaking;
      tokenStakingEscrow = stakingContracts.tokenStakingEscrow;

      await tokenGrant.authorizeStakingContract(tokenStaking.address, {
        from: grantManager,
      })

      undelegationPeriod = await tokenStaking.undelegationPeriod()
      
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
          'Not authorized'
        )
      })

      it('should not let grantee cancel delegation of another grantee', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorTwo, {from: grantee}),
          'Not authorized'
        ) 
      })

      it('should not let managed grantee cancel delegation of another grantee', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorOne, {from: managedGrantee}),
          'Not authorized'
        )     
      })

      it('should not let third party cancel delegation', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorOne, {from: thirdParty}),
          'Not authorized'
        )
        await expectRevert(
          tokenStaking.cancelStake(operatorTwo, {from: thirdParty}),
          'Not authorized'
        )    
      })

      it('should not let grant manager cancel delegation of non-revoked grant', async () => {
        await expectRevert(
          tokenStaking.cancelStake(operatorOne, {from: grantManager}),
          'Not authorized'
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

      it('transfers tokens to escrow', async () => {
        await tokenStaking.cancelStake(operatorOne, {from: grantee})
        await tokenStaking.cancelStake(operatorTwo, {from: managedGrantee})

        let deposited = await tokenStakingEscrow.depositedAmount(operatorOne)
        expect(deposited).to.eq.BN(delegatedAmount)

        deposited = await tokenStakingEscrow.depositedAmount(operatorTwo)
        expect(deposited).to.eq.BN(delegatedAmount)
      })

      it('fails if already cancelled', async () => {
        await tokenStaking.cancelStake(operatorOne, {from: grantee})
        await expectRevert(
          tokenStaking.cancelStake(operatorOne, {from: grantee}),
          "Stake for the operator already deposited in the escrow"
        )

        await tokenStaking.cancelStake(operatorTwo, {from: managedGrantee})
        await expectRevert(
          tokenStaking.cancelStake(operatorTwo, {from: managedGrantee}),
          "Stake for the operator already deposited in the escrow"
        )
      })
    })

    describe('undelegate', async () => {
      beforeEach(async () => {
        await time.increase(initializationPeriod.addn(1))
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
          'Not authorized'
        )   
      })

      it('should not let grantee undelegate for another grantee', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorTwo, {from: grantee}),
          'Not authorized'
        )     
      })

      it('should not let managed grantee undelegate for another grantee', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorOne, {from: managedGrantee}),
          'Not authorized'
        )    
      })

      it('should not let third party undelegate', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorOne, {from: thirdParty}),
          'Not authorized'
        )
        await expectRevert(
          tokenStaking.undelegate(operatorTwo, {from: thirdParty}),
          'Not authorized'
        )    
      })

      it('should not let grant manager undelegate non-revoked grant', async () => {
        await expectRevert(
          tokenStaking.undelegate(operatorOne, {from: grantManager}),
          'Not authorized'
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
      it('transfers tokens to escrow', async () => {
        await time.increase(initializationPeriod.addn(1))

        await tokenStaking.undelegate(operatorOne, {from: operatorOne})
        await tokenStaking.undelegate(operatorTwo, {from: operatorTwo})

        await time.increase(undelegationPeriod.addn(1))

        await tokenStaking.recoverStake(operatorOne, {from: thirdParty})
        await tokenStaking.recoverStake(operatorTwo, {from: thirdParty})

        let deposited = await tokenStakingEscrow.depositedAmount(operatorOne)
        expect(deposited).to.eq.BN(delegatedAmount)

        deposited = await tokenStakingEscrow.depositedAmount(operatorTwo)
        expect(deposited).to.eq.BN(delegatedAmount)
      })

      it('fails if already recovered', async () => {
        await time.increase(initializationPeriod.addn(1))

        await tokenStaking.undelegate(operatorOne, {from: operatorOne})
        await tokenStaking.undelegate(operatorTwo, {from: operatorTwo})

        await time.increase(undelegationPeriod.addn(1))

        await tokenStaking.recoverStake(operatorOne, {from: thirdParty})
        await expectRevert(
          tokenStaking.recoverStake(operatorOne, {from: thirdParty}),
          "Stake for the operator already deposited in the escrow"
        )

        await tokenStaking.recoverStake(operatorTwo, {from: thirdParty})
        await expectRevert(
          tokenStaking.recoverStake(operatorTwo, {from: thirdParty}),
          "Stake for the operator already deposited in the escrow"
        )
      })
    })

    describe('redelegate from escrow', async () => {
      const data3 = Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operatorThree.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ])
      const data4 = Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operatorFour.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ])

      beforeEach(async () => {
        await tokenStaking.cancelStake(operatorOne, {from: operatorOne})

        await time.increase(initializationPeriod.addn(1))
        await tokenStaking.undelegate(operatorTwo, {from: operatorTwo})
        await time.increase(undelegationPeriod.addn(1))
        await tokenStaking.recoverStake(operatorTwo, {from: thirdParty})
      })
      
      it('can be done by grantee', async () => {
        await tokenStakingEscrow.redelegate(
          operatorOne, delegatedAmount, data3, {from: grantee}
        )
        // ok, no revert
      })

      it('can be done by managed grantee', async () => {
        await tokenStakingEscrow.redelegate(
          operatorTwo, delegatedAmount, data3, {from: managedGrantee}
        )
        // ok, no revert
      })

      it('can not be done by operator', async () => {
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne, delegatedAmount, data3, {from: operatorOne}
          ),
          "Not authorized"
        )
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne, delegatedAmount, data3, {from: operatorThree}
          ),
          "Not authorized"
        )
      })

      it('can not be done by grant manager', async () => {
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne, delegatedAmount, data3, {from: grantManager}
          ),
          "Not authorized"
        )
      })

      it('can not be done by third party', async () => {
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne, delegatedAmount, data3, {from: thirdParty}
          ),
          "Not authorized"
        )
      })

      it('redelegates token to a new operator', async () => {
        const redelegatedAmount = delegatedAmount

        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )

        const delegationInfo = await tokenStaking.getDelegationInfo(operatorThree)
        expect(delegationInfo.amount).to.eq.BN(redelegatedAmount)

        const availableAmount = await tokenStakingEscrow.availableAmount(operatorOne)
        expect(availableAmount).to.eq.BN(0)
      })

      it('redelegates tokens to more than one operator', async () => {
        const redelegatedAmount = delegatedAmount.divn(2)

        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )
        let availableAmount = await tokenStakingEscrow.availableAmount(operatorOne)
        expect(availableAmount).to.eq.BN(redelegatedAmount)
        let delegationInfo = await tokenStaking.getDelegationInfo(operatorThree)
        expect(delegationInfo.amount).to.eq.BN(redelegatedAmount)
      
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data4, {from: grantee}
        )
        availableAmount = await tokenStakingEscrow.availableAmount(operatorOne)
        expect(availableAmount).to.eq.BN(0)
        delegationInfo = await tokenStaking.getDelegationInfo(operatorFour)
        expect(delegationInfo.amount).to.eq.BN(redelegatedAmount)
      })

      it('fails for revoked grant', async () => {
        await tokenGrant.revoke(grantId, {from: grantManager})
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne, delegatedAmount, data3, {from: grantee}
          ),
          "Grant revoked"
        )
      })

      it("fails when trying to redelegate to operator with cancelled stake", async () => {
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne,
            delegatedAmount,
            Buffer.concat([
              Buffer.from(beneficiary.substr(2), 'hex'),
              Buffer.from(operatorOne.substr(2), 'hex'),
              Buffer.from(authorizer.substr(2), 'hex')
            ]),
            {from: grantee}
          ),
          "Redelegating to previously used operator is not allowed"
        )
      })

      it('fails when trying to redelegate to operator with undelegated stake', async () => {
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorTwo,
            delegatedAmount,
            Buffer.concat([
              Buffer.from(beneficiary.substr(2), 'hex'),
              Buffer.from(operatorTwo.substr(2), 'hex'),
              Buffer.from(authorizer.substr(2), 'hex')
            ]),
            {from: managedGrantee}
          ),
          "Redelegating to previously used operator is not allowed"
        )
      })

      it('fails when trying to redelegate more then deposited', async () => {
        let redelegatedAmount = delegatedAmount.addn(1)
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne, redelegatedAmount, data3, {from: grantee}
          ),
          "Insufficient balance"
        )
      })

      it('fails when trying to redelegate more than remaining', async () => {
        let redelegatedAmount = delegatedAmount.divn(2)
        await tokenStakingEscrow.redelegate(
          operatorOne, delegatedAmount, data3, {from: grantee}
        )

        redelegatedAmount = redelegatedAmount.addn(1)
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne, redelegatedAmount, data4, {from: grantee}
          ),
          "Insufficient balance"
        )
      })

      it('redelegates not yet withdrawn tokens', async () => {
        await time.increaseTo(grantStart.add(time.duration.years(1)))
        // 2 000 000 undelegated to escrow for 2-years grant.
        // One year passed, so 50% of tokens, 1 000 000, is withdrawable
        // from the escrow. Let's withdraw them.
        await tokenStakingEscrow.withdraw(operatorOne, {from: grantee})

        // And now, let's redelegate the remaining 1 000 000 KEEP...

        // 1000000000000000000000000 
        const redelegatedAmount = web3.utils.toWei('1000000')

        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )

        const availableAmount = await tokenStakingEscrow.availableAmount(operatorOne)
        expect(availableAmount).to.eq.BN(0)
      })

      it('fails when trying to redelegate withdrawn tokens', async () => {
        await time.increaseTo(grantStart.add(time.duration.years(1)))
        // 2 000 000 undelegated to escrow for 2-years grant.
        // One year passed, so 50% of tokens, 1 000 000, is withdrawable
        // from the escrow. Let's withdraw them.
        await tokenStakingEscrow.withdraw(operatorOne, {from: grantee})

        // And now, let's try to redelegate 1 000 000 KEEP + 1e-18 KEEP...

        let redelegatedAmount = web3.utils.toWei('1000000')
        // 1000000000000000000000000 + 1 =
        // 1000000000000000000000001
        redelegatedAmount = web3.utils.toBN(redelegatedAmount).addn(1)
        
        await expectRevert(
          tokenStakingEscrow.redelegate(
            operatorOne, redelegatedAmount, data3, {from: grantee}
          ),
          "Insufficient balance"
        )
      })

      it('allows to withdraw not redelegated tokens', async () => {
        await time.increaseTo(grantStart.add(time.duration.years(1)))
        // 2 000 000 undelegated to escrow for 2-years grant.
        // One year passed, so 50% of tokens, 1 000 000, is withdrawable
        // from the escrow. Let's withdraw them.
        await tokenStakingEscrow.withdraw(operatorOne, {from: grantee})

        // And now, let's redelegate 1 000 000 - 1 KEEP ...
        const redelegatedAmount = web3.utils.toWei('999999')
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )

        let availableAmount = await tokenStakingEscrow.availableAmount(operatorOne)
        expect(availableAmount).to.eq.BN(web3.utils.toWei('1'))

        // Finally, we need to wait until the remaining 1 KEEP becomes
        // withdrawable and withdraw it.
        await time.increaseTo(grantStart.add(time.duration.years(2)))
        const withdrawable = await tokenStakingEscrow.withdrawable(operatorOne)
        expect(withdrawable).to.eq.BN(web3.utils.toWei('1'))
        await tokenStakingEscrow.withdraw(operatorOne, {from: grantee})
        // ok, no revert

        availableAmount = await tokenStakingEscrow.availableAmount(operatorOne)
        expect(availableAmount).to.eq.BN(0)
      })
  
      it('can be cancelled by grantee', async () => {
        const redelegatedAmount = delegatedAmount
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )
        await tokenStaking.cancelStake(operatorThree, {from: grantee})
        // ok, no reverts
      })

      it('can be cancelled by managed grantee', async () => {
        const redelegatedAmount = delegatedAmount
        await tokenStakingEscrow.redelegate(
          operatorTwo, redelegatedAmount, data3, {from: managedGrantee}
        )
        await tokenStaking.cancelStake(operatorThree, {from: managedGrantee})
        // ok, no reverts 
      })

      it('can be cancelled by operator', async () => {
        const redelegatedAmount = delegatedAmount
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )
        await tokenStaking.cancelStake(operatorThree, {from: operatorThree})
        // ok, no reverts  
      })

      it('can not be cancelled by the previous operator', async () => {
        const redelegatedAmount = delegatedAmount
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )
        await expectRevert(
          tokenStaking.cancelStake(operatorThree, {from: operatorOne}),
          "Not authorized"
        )
      })

      it('can be undelegated by grantee', async () => {
        const redelegatedAmount = delegatedAmount
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )
        await time.increase(initializationPeriod.addn(1))

        await tokenStaking.undelegate(operatorThree, {from: grantee})
        // ok, no reverts 
      })

      it('can be undelegated by managed grantee', async () => {
        const redelegatedAmount = delegatedAmount
        await tokenStakingEscrow.redelegate(
          operatorTwo, redelegatedAmount, data3, {from: managedGrantee}
        )
        await time.increase(initializationPeriod.addn(1))

        await tokenStaking.undelegate(operatorThree, {from: managedGrantee})
        // ok, no reverts  
      })

      it('can be undelegated by operator', async () => {
        const redelegatedAmount = delegatedAmount
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )
        await time.increase(initializationPeriod.addn(1))

        await tokenStaking.undelegate(operatorThree, {from: operatorThree})
        // ok, no reverts 
      })

      it('can not be undelegated by the previous operator', async () => {
        const redelegatedAmount = delegatedAmount
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )
        await time.increase(initializationPeriod.addn(1))

        await expectRevert(
          tokenStaking.undelegate(operatorThree, {from: operatorOne}),
          "Not authorized"
        )
      })

      it('lands back in the escrow when undelegated and recovered', async () => {
        const redelegatedAmount = delegatedAmount.subn(10)
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )

        await time.increase(initializationPeriod.addn(1))
        await tokenStaking.undelegate(operatorThree, {from: operatorThree})

        await time.increase(undelegationPeriod.addn(1))
        await tokenStaking.recoverStake(operatorThree)

        expect(
          await tokenStakingEscrow.depositedAmount(operatorThree)
        ).to.eq.BN(redelegatedAmount)
        expect(
          await tokenStakingEscrow.depositGrantId(operatorThree)
        ).to.eq.BN(grantId)
        expect(
          await tokenStakingEscrow.depositWithdrawnAmount(operatorThree)
        ).to.eq.BN(0)
        expect(
          await tokenStakingEscrow.depositRedelegatedAmount(operatorThree)
        ).to.eq.BN(0)
      })

      it('lands back in the escrow when cancelled', async () => {
        const redelegatedAmount = delegatedAmount.subn(10)
        await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )

        await tokenStaking.cancelStake(operatorThree, {from: operatorThree})

        expect(
          await tokenStakingEscrow.depositedAmount(operatorThree)
        ).to.eq.BN(redelegatedAmount)
        expect(
          await tokenStakingEscrow.depositGrantId(operatorThree)
        ).to.eq.BN(grantId)
        expect(
          await tokenStakingEscrow.depositWithdrawnAmount(operatorThree)
        ).to.eq.BN(0)
        expect(
          await tokenStakingEscrow.depositRedelegatedAmount(operatorThree)
        ).to.eq.BN(0)
      })

      it('emits an event', async () => {
        const redelegatedAmount = delegatedAmount
        const receipt = await tokenStakingEscrow.redelegate(
          operatorOne, redelegatedAmount, data3, {from: grantee}
        )
        
        await expectEvent(receipt, 'DepositRedelegated', {
          previousOperator: operatorOne,
          newOperator: operatorThree,
          grantId: grantId.toString(),
          amount: redelegatedAmount.toString()
        })
      })
    })
})