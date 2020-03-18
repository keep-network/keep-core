import { duration, increaseTimeTo } from '../helpers/increaseTime';
import latestTime from '../helpers/latestTime';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage'
import grantTokens from '../helpers/grantTokens';
import { createSnapshot, restoreSnapshot } from '../helpers/snapshot'

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const Registry = artifacts.require("./Registry.sol");

contract('TokenGrant/Revoke', function(accounts) {

  let tokenContract, registryContract, grantContract, stakingContract;

  const tokenOwner = accounts[0],
    grantee = accounts[1];

  let grantId;
  let grantStart;
  const grantAmount = web3.utils.toBN(1000000000);
  const grantRevocable = true;
  const grantDuration = duration.seconds(60);;
  const grantCliff = duration.seconds(1);
    
  const initializationPeriod = 10;
  const undelegationPeriod = 30;

  before(async () => {
    tokenContract = await KeepToken.new();
    registryContract = await Registry.new();
    stakingContract = await TokenStaking.new(
      tokenContract.address, 
      registryContract.address, 
      initializationPeriod, 
      undelegationPeriod
    );
    grantContract = await TokenGrant.new(tokenContract.address);
    
    await grantContract.authorizeStakingContract(stakingContract.address);

    grantStart = await latestTime();

    grantId = await grantTokens(
      grantContract, 
      tokenContract,
      grantAmount,
      tokenOwner, 
      grantee, 
      grantDuration, 
      grantStart, 
      grantCliff, 
      grantRevocable,
    );
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should allow to revoke grant", async () => {
    const grantManagerKeepBalanceBefore = await tokenContract.balanceOf(tokenOwner);
    await increaseTimeTo(grantStart + duration.seconds(30));
    const withdrawable = await grantContract.withdrawable(grantId);
    const refund = grantAmount.sub(withdrawable);
    
    await grantContract.revoke(grantId, { from: tokenOwner });

    const withdrawableAfter = await grantContract.withdrawable(grantId);
    const grantDetails = await grantContract.getGrant(grantId);
    const grantManagerKeepBalanceAfter = await tokenContract.balanceOf(tokenOwner);
    const unlockedAmount = await grantContract.unlockedAmount(grantId);

    expect(grantManagerKeepBalanceAfter).to.eq.BN(
      grantManagerKeepBalanceBefore.add(refund),
      "The grant manager KEEP balance should be updated"
    );
    expect(grantDetails.revokedAt).to.be.gt.BN(
      0,
      "revokedAt should be greater than zero"
    );
    expect(withdrawableAfter.add(refund)).to.eq.BN(
      grantAmount,
      "Should be equal to the total grant amount"
    );
    expect(grantDetails.revokedAmount).to.eq.BN(
      grantAmount.sub(unlockedAmount),
      "Revoked amount should be equal to the subtraction grant amount and unlocked amount"
    );
    expect(grantDetails.revokedAmount).to.eq.BN(
      refund,
      "Revoked amount should be equal to returned amount to the grant creator"
    )
  })

  it("should not allow to revoke grant if sender is not a grant manager", async () => {
    await expectThrowWithMessage(
      grantContract.revoke(grantId, { from: grantee }),
      "Only grant manager can revoke."
    );
  })

  it("should not allow to revoke grant if the grant is non revocable", async () => {
    const nonRevocableGrantId= await grantTokens(
        grantContract, 
        tokenContract,
        grantAmount,
        tokenOwner, 
        grantee, 
        grantDuration, 
        grantStart, 
        grantCliff, 
        false,
    );
    
    await expectThrowWithMessage(
      grantContract.revoke(nonRevocableGrantId, { from: tokenOwner }),
      "Grant must be revocable in the first place."
    );
  })

  it("should not allow to revoke grant multiple times", async () => {
    await grantContract.revoke(grantId, { from: tokenOwner });
  
    await expectThrowWithMessage(
      grantContract.revoke(grantId, { from: tokenOwner }),
      "Grant must not be already revoked."
    );
  })

  it("should be able to revoke the grant but no amount is refunded since duration of the unlocking is over.", async () => {
    const grantDuration = web3.utils.toBN(0);
    const grantCliff = web3.utils.toBN(0);

    const grantManagerKeepBalance = await tokenContract.balanceOf(tokenOwner);
    
    const fullyUnlockedGrantId = await grantTokens(
      grantContract, 
      tokenContract,
      grantAmount,
      tokenOwner, 
      grantee, 
      grantDuration, 
      grantStart, 
      grantCliff, 
      grantRevocable,
      );
      
    const granteeGrantBalanceBefore = await grantContract.balanceOf.call(grantee);
    const grantManagerKeepBalanceBeforeRevoke= await tokenContract.balanceOf(tokenOwner);

    await grantContract.revoke(fullyUnlockedGrantId, { from: tokenOwner });
    
    const grantManagerKeepBalanceAfterRevoke = await tokenContract.balanceOf(tokenOwner);
    const granteeGrantBalanceAfter = await grantContract.balanceOf.call(grantee);

    expect(grantManagerKeepBalanceBeforeRevoke).to.eq.BN(
      grantManagerKeepBalance.sub(grantAmount),
      "Amount should be taken out from grant manager main balance"
    );
    expect(granteeGrantBalanceAfter).to.eq.BN(
      granteeGrantBalanceBefore,
      "Amount should stay at grantee's grant balance"
    );
    expect(grantManagerKeepBalanceAfterRevoke).to.eq.BN(
      grantManagerKeepBalanceBeforeRevoke,
      "No amount to be returned to grant manager since unlocking duration is over"
    );
  });
});
