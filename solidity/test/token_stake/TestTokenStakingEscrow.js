const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, expectEvent, time} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot');

const {grantTokens, grantTokensToManagedGrant} = require('../helpers/grantTokens');

const KeepToken = contract.fromArtifact('KeepToken')
const TokenGrant = contract.fromArtifact('TokenGrant')
const PermissiveStakingPolicy = contract.fromArtifact('PermissiveStakingPolicy')
const ManagedGrantFactory = contract.fromArtifact('ManagedGrantFactory')
const ManagedGrant = contract.fromArtifact('ManagedGrant')
const TokenStakingEscrow = contract.fromArtifact('TokenStakingEscrow')
const MigrationEscrowStub = contract.fromArtifact('MigrationEscrowStub')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('TokenStakingEscrow', () => {
  
  const deployer = accounts[0],
    grantManager = accounts[1],
    grantee = accounts[2],
    operator = accounts[3],
    operator2 = accounts[4],
    thirdParty = accounts[5],
    tokenStaking = accounts[6]

  let grantedAmount, grantStart, grantUnlockingDuration,
  grantId, managedGrantId, managedGrant

  let token, tokenGrant, permissivePolicy, managedGrantFactory, escrow

  before(async () => {
    token = await KeepToken.new({from: deployer})
    await token.transfer(tokenStaking, 100000, {from: deployer})
    await token.transfer(grantManager, 100000, {from: deployer})

    tokenGrant = await TokenGrant.new(token.address, {from: deployer})
    permissivePolicy = await PermissiveStakingPolicy.new()
    managedGrantFactory = await ManagedGrantFactory.new(
      token.address,
      tokenGrant.address,
      {from: deployer}
    );
    
    escrow = await TokenStakingEscrow.new(
      token.address, 
      tokenGrant.address,
      {from: deployer}
    )

    await escrow.transferOwnership(tokenStaking, {from: deployer})

    grantedAmount = 10000
    grantStart = await time.latest()
    grantCliff = time.duration.days(5)
    grantUnlockingDuration = time.duration.days(30)
    
    grantId = await grantTokens(
      tokenGrant, 
      token, 
      grantedAmount, 
      grantManager, 
      grantee, 
      grantUnlockingDuration,
      grantStart,
      grantCliff,
      true,
      permissivePolicy.address
    )

    const managedGrantAddress = await grantTokensToManagedGrant(
      managedGrantFactory,
      token,
      grantedAmount,
      grantManager,
      grantee,
      grantUnlockingDuration,
      grantStart,
      grantCliff,
      false,
      permissivePolicy.address,
    )
    managedGrant = await ManagedGrant.at(managedGrantAddress)
    managedGrantId = (await managedGrant.grantId()).toNumber()
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe('receiveApproval', async () => {
    it('reverts for unknown token', async () => {
      let anotherToken = await KeepToken.new({from: deployer})
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )

      await expectRevert(
        anotherToken.approveAndCall(
            escrow.address, grantedAmount, data, {from: tokenStaking}
        ),
        "Not a KEEP token"
      )
    })

    it('reverts for corrupted extraData', async () => {
      const corruptedData = web3.eth.abi.encodeParameters(
        ['address'], [operator]
      )

      await expectRevert(
        token.approveAndCall(
            escrow.address, grantedAmount, corruptedData, {from: tokenStaking}
        ),
        "Unexpected data length"
      )
    })

    it('reverts for unknown grant', async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, 999]
      )

      await expectRevert(
        token.approveAndCall(
            escrow.address, grantedAmount, data, {from: tokenStaking}
        ),
        "Grant with this ID does not exist"
      )
    })

    it('deposits KEEP', async () => {    
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )
      await token.approveAndCall(
        escrow.address, grantedAmount, data, {from: tokenStaking}
      )

      const deposited = await escrow.depositedAmount(operator)
      expect(deposited).to.eq.BN(grantedAmount)
    })

    it('accepts deposits only from the owner', async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )

      await expectRevert(
        token.approveAndCall(
          escrow.address, grantedAmount, data, {from: operator}
        ),
        "Only owner can deposit"
      )
      await expectRevert(
        token.approveAndCall(
          escrow.address, grantedAmount, data, {from: grantee}
        ),
        "Only owner can deposit"
      )
      await expectRevert(
        token.approveAndCall(
          escrow.address, grantedAmount, data, {from: deployer}
        ),
        "Only owner can deposit"
      )
    })
  })

  describe('depositedAmount', async () => {
    it('returns 0 for unknown operator', async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )
      await token.approveAndCall(
        escrow.address, grantedAmount, data, {from: tokenStaking}
      )

      const deposited = await escrow.depositedAmount(grantee)
      expect(deposited).to.eq.BN(0)   
    })
  })

  describe('withdrawable', async () => {
    const depositedAmount = 1000
    beforeEach(async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )
      await token.approveAndCall(
        escrow.address, depositedAmount, data, {from: tokenStaking}
      )
    })

    it('returns 0 for unknown operator', async () => {
      const withdrawable = await escrow.withdrawable(grantee)
      expect(withdrawable).to.eq.BN(0) 
    })

    it('returns 0 just before the cliff', async () => {
      await time.increaseTo(
        // 1 minute before the cliff ends
        grantStart.add(grantCliff).sub(time.duration.minutes(1))
      )
      const withdrawable = await escrow.withdrawable(operator)
      expect(withdrawable).to.eq.BN(0) 
    })

    it('returns unlocked amount just after the cliff', async () => {
      await time.increaseTo(grantStart.add(grantCliff))
      const withdrawable = await escrow.withdrawable(operator)
      expect(withdrawable).to.eq.BN(166) // (1000 / 30) * 5 = 166
    })

    it('returns unlocked amount in the middle of unlocking period', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      const withdrawable = await escrow.withdrawable(operator)
      expect(withdrawable).to.eq.BN(500) // (1000 / 30) * 15 = 500 
    })

    it('returns whole deposited amount after it unlocked', async () => {
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      const withdrawable = await escrow.withdrawable(operator)
      expect(withdrawable).to.eq.BN(depositedAmount) 
    })

    it('returns 0 just after the cliff if all unlocked withdrawn', async () => {
      await time.increaseTo(grantStart.add(grantCliff))
      await escrow.withdraw(operator, {from: grantee})
      const withdrawable = await escrow.withdrawable(operator)
      expect(withdrawable).to.eq.BN(0)
    })

    it('returns remaining unlocked, non-withdrawn amount', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await escrow.withdraw(operator, {from: grantee}) // withdraws 500
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      const withdrawable = await escrow.withdrawable(operator)
      expect(withdrawable).to.eq.BN(500) // the remaining 500
    })

    it('returns 0 for revoked grant', async () => {
      await time.increaseTo(grantStart.add(grantCliff))
      await tokenGrant.revoke(grantId, {from: grantManager})
      const withdrawable = await escrow.withdrawable(operator)
      expect(withdrawable).to.eq.BN(0)
    })
  })

  describe('withdraw', async () => {
    const depositedAmount = 2000
    beforeEach(async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )
      await token.approveAndCall(
        escrow.address, depositedAmount, data, {from: tokenStaking}
      )
    })

    it('can be called by grantee', async () => {
      await escrow.withdraw(operator, {from: grantee})
      // ok, no reverts
    })

    it('can be called by operator', async () => {
      await escrow.withdraw(operator, {from: operator})
      // ok, no reverts
    })

    it('can not be called by third-party', async () => {
      await expectRevert(
        escrow.withdraw(operator, {from: thirdParty}),
        "Only grantee or operator can withdraw" 
      )
    })

    it('can not be called by deployer', async () => {
      await expectRevert(
        escrow.withdraw(operator, {from: deployer}),
        "Only grantee or operator can withdraw" 
      )
    })

    it('withdraws entire unlocked amount just after the cliff', async () => {
      await time.increaseTo(grantStart.add(grantCliff))
      await escrow.withdraw(operator, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(333) // (2000 / 30) * 5 = 333
    })

    it('withdraws entire unlocked amount in the middle of unlocking period', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await escrow.withdraw(operator, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(1000) // (2000 / 30) * 15 = 1000  
    })

    it('withdraws entire unlocked amount after the whole unlocking period', async () => {
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      await escrow.withdraw(operator, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(depositedAmount)
    })

    it('allows to withdraw in multiple rounds', async () => {
      await time.increaseTo(grantStart.add(grantCliff))
      await escrow.withdraw(operator, {from: grantee})

      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await escrow.withdraw(operator, {from: grantee})

      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      await escrow.withdraw(operator, {from: grantee})

      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(depositedAmount)
    })

    it('withdraws entire deposit for fully unlocked, revoked grant', async () => {
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      await tokenGrant.revoke(grantId, {from: grantManager})
      await escrow.withdraw(operator, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(depositedAmount)
    })

    it('withdraws nothing if already withdrawn', async () => {
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      await escrow.withdraw(operator, {from: grantee})
      await escrow.withdraw(operator, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(depositedAmount)
    })

    it('emits an event', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      const receipt = await escrow.withdraw(operator, {from: grantee})

      await expectEvent(receipt, 'DepositWithdrawn', {
        operator: operator,
        grantee: grantee,
        amount: web3.utils.toBN(1000)// (2000 / 30) * 15 = 1000 
      })
    })
    
    it('can not be called for managed grant', async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator2, managedGrantId]
      )
      await token.approveAndCall(
        escrow.address, 600, data, {from: tokenStaking}
      )

      await expectRevert(
          escrow.withdraw(operator2, {from: operator2}),
          "Can not be called for managed grant"
      );
    })
  })

  describe('withdrawToManagedGrantee', async () => {
    const depositedAmount = 2000
    beforeEach(async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator2, managedGrantId]
      )
      await token.approveAndCall(
        escrow.address, depositedAmount, data, {from: tokenStaking}
      )
    })

    it('can be called by grantee', async () => {
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
      // ok, no reverts
    })
  
    it('can be called by operator', async () => {
      await escrow.withdrawToManagedGrantee(operator2, {from: operator2})
      // ok, no reverts
    })

    it('can not be called by third-party', async () => {
      await expectRevert(
        escrow.withdrawToManagedGrantee(operator2, {from: thirdParty}),
        "Only grantee or operator can withdraw" 
      )
    })

    it('can not be called by deployer', async () => {
      await expectRevert(
        escrow.withdrawToManagedGrantee(operator2, {from: deployer}),
        "Only grantee or operator can withdraw" 
      )
    })

    it('withdraws entire unlocked amount just after the cliff', async () => {
      await time.increaseTo(grantStart.add(grantCliff))
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(333) // (2000 / 30) * 5 = 333
    })
  
    it('withdraws entire unlocked amount in the middle of unlocking period', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(1000) // (2000 / 30) * 15 = 1000  
    })
  
    it('withdraws entire unlocked amount after the whole unlocking period', async () => {
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(depositedAmount)
    })

    it('allows to withdraw in multiple rounds', async () => {
      await time.increaseTo(grantStart.add(grantCliff))
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
  
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
  
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
  
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(depositedAmount)
    })

    it('withdraws nothing if already withdrawn', async () => {
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
      await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
      const balance = await token.balanceOf(grantee);
      expect(balance).to.eq.BN(depositedAmount) 
    })

    it('emits an event', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      const receipt = await escrow.withdrawToManagedGrantee(operator2, {from: grantee})
  
      await expectEvent(receipt, 'DepositWithdrawn', {
        operator: operator2,
        grantee: grantee,
        amount: web3.utils.toBN(1000)// (2000 / 30) * 15 = 1000 
      })
    })
  })

  describe('withdrawRevoked', async () => {
    const depositedAmount = 10000
    beforeEach(async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )
      await token.approveAndCall(
        escrow.address, depositedAmount, data, {from: tokenStaking}
      )
    })

    it('can be called by grant manager', async () => {
      await tokenGrant.revoke(grantId, {from: grantManager})
      await escrow.withdrawRevoked(operator, {from: grantManager})
      // ok, no reverts
    })

    it('can be called by grantee', async () => {
      await tokenGrant.revoke(grantId, {from: grantManager})
      await escrow.withdrawRevoked(operator, {from: grantee})
      // ok, no reverts
    })

    it('can be called by operator', async () => {
      await tokenGrant.revoke(grantId, {from: grantManager})
      await escrow.withdrawRevoked(operator, {from: operator})
      // ok, no reverts
    })

    it('can be called by third party', async () => {
      await tokenGrant.revoke(grantId, {from: grantManager})
      await escrow.withdrawRevoked(operator, {from: thirdParty})
      // ok, no reverts
    })

    it('can not be called for non-revoked grant', async () => {
      await expectRevert(
        escrow.withdrawRevoked(operator, {from: grantManager}),
        "No revoked tokens to withdraw"
      )
    })

    it('withdraws part of deposited amount if something has been withdrawn before', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await escrow.withdraw(operator, {from: operator}) // (1000 / 30) * 15 = 5000
      await tokenGrant.revoke(grantId, {from: grantManager})

      const balanceBefore = await token.balanceOf(grantManager)
      await escrow.withdrawRevoked(operator, {from: grantManager})
      const balanceAfter = await token.balanceOf(grantManager)

      const diff = balanceAfter.sub(balanceBefore)
      expect(diff).to.eq.BN(5000) // 10000 - 5000 = 5000
    })

    it('withdraws entire deposited amount if nothing has been withdrawn before', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await tokenGrant.revoke(grantId, {from: grantManager})

      const balanceBefore = await token.balanceOf(grantManager)
      await escrow.withdrawRevoked(operator, {from: grantManager})
      const balanceAfter = await token.balanceOf(grantManager)

      const diff = balanceAfter.sub(balanceBefore)
      expect(diff).to.eq.BN(depositedAmount)
    })

    it('withdraws nothing if already withdrawn', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await tokenGrant.revoke(grantId, {from: grantManager})

      await escrow.withdrawRevoked(operator, {from: grantManager})

      const balanceBefore = await token.balanceOf(grantManager)
      await escrow.withdrawRevoked(operator, {from: grantManager})
      const balanceAfter = await token.balanceOf(grantManager)

      const diff = balanceAfter.sub(balanceBefore)
      expect(diff).to.eq.BN(0) 
    })

    it('reverts if the entire grant unlocked', async () => {
      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      await escrow.withdraw(operator, {from: operator})

      await expectRevert(
        escrow.withdrawRevoked(operator, {from: grantManager}),
        "No revoked tokens to withdraw"
      )
    })

    it('emits an event', async () => {
      await tokenGrant.revoke(grantId, {from: grantManager})
      const receipt = await escrow.withdrawRevoked(operator, {from: grantManager})

      await expectEvent(receipt, 'RevokedDepositWithdrawn', {
        operator: operator,
        grantManager: grantManager,
        amount: web3.utils.toBN(depositedAmount)
      })
    })
  })
  
  describe('depositWithdrawnAmount', async () => {
    const depositedAmount = 3000

    beforeEach(async () => {
      const data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )
      await token.approveAndCall(
        escrow.address, depositedAmount, data, {from: tokenStaking}
      )
    })

    it('returns 0 for unknown operator', async () => {
      const withdrawn = await escrow.depositWithdrawnAmount(grantee)
      expect(withdrawn).to.eq.BN(0)   
    })

    it('returns 0 if nothing has been withdrawn', async () => {
      const withdrawn = await escrow.depositWithdrawnAmount(operator)
      expect(withdrawn).to.eq.BN(0)   
    })
  
    it('returns withdrawn amount in the middle of unlocking period', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await escrow.withdraw(operator, {from: grantee})
      const withdrawn = await escrow.depositWithdrawnAmount(operator)
      expect(withdrawn).to.eq.BN(1500) // (3000 / 30) * 15 = 1500  
    })
  
    it('returns withdrawn amount at the end of unlocking period', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(grantUnlockingDuration)))
      await escrow.withdraw(operator, {from: grantee})
      const withdrawn = await escrow.depositWithdrawnAmount(operator)
      expect(withdrawn).to.eq.BN(depositedAmount)
    })
  })

  describe('migrate', async () => {
    const depositedAmount = 3000
    let anotherEscrow

    beforeEach(async () => {
      let data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator, grantId]
      )
      await token.approveAndCall(
        escrow.address, depositedAmount, data, {from: tokenStaking}
      )

      data = web3.eth.abi.encodeParameters(
        ['address', 'uint256'], [operator2, managedGrantId]
      )
      await token.approveAndCall(
        escrow.address, depositedAmount, data, {from: tokenStaking}
      )

      anotherEscrow = await MigrationEscrowStub.new({from: deployer})
      await escrow.authorizeEscrow(anotherEscrow.address, {from: grantManager})
    })

    it('fails for not authorized escrow', async () => {
      await expectRevert(
        escrow.migrate(operator, thirdParty, {from: grantee}),
        "Escrow not authorized"
      )
    })

    it('can be called by grantee', async () => {
      await escrow.migrate(operator, anotherEscrow.address, {from: grantee})
      // on, no revert
    })

    it('can be called by grantee of managed grant', async () => {
      await escrow.migrate(operator2, anotherEscrow.address, {from: grantee})
      // on, no revert
    })

    it('can not be called by operator', async () => {
      await expectRevert(
        escrow.migrate(operator, anotherEscrow.address, {from: operator}),
        "Not authorized"
      )
    })

    it('can not be called by grant manager', async () => {
      await expectRevert(
        escrow.migrate(operator, anotherEscrow.address, {from: grantManager}),
        "Not authorized"
      )
    })

    it('can not be called by third party', async () => {
      await expectRevert(
        escrow.migrate(operator, anotherEscrow.address, {from: thirdParty}),
        "Not authorized"
      )
    })

    it('moves all tokens to another escrow', async () => {
      await escrow.migrate(operator, anotherEscrow.address, {from: grantee})

      await time.increaseTo(grantStart.add(grantUnlockingDuration))
      expect(await escrow.withdrawable(operator)).to.eq.BN(0)
      expect(await anotherEscrow.depositedAmount(operator)).to.eq.BN(depositedAmount)
    })

    it('move the rest of tokens to another escrow', async () => {
      await time.increaseTo(grantStart.add(time.duration.days(15)))
      await escrow.withdraw(operator, {from: grantee}) // (3000 / 30) * 15 = 1500
      
      await time.increaseTo(grantStart.add(grantUnlockingDuration))

      await escrow.migrate(operator, anotherEscrow.address, {from: grantee})

      expect(await escrow.withdrawable(operator)).to.eq.BN(0)
      expect(await anotherEscrow.depositedAmount(operator)).to.eq.BN(1500) // 3000 - 1500
    })
  })
})