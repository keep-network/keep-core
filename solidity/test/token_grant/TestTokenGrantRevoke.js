const { delegateStakeFromGrant } = require('../helpers/delegateStake')
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const {initTokenStaking} = require('../helpers/initContracts')
const {grantTokens} = require('../helpers/grantTokens');
const { createSnapshot, restoreSnapshot } = require('../helpers/snapshot');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const KeepToken = contract.fromArtifact('KeepToken');
const TokenGrant = contract.fromArtifact('TokenGrant');
const KeepRegistry = contract.fromArtifact("KeepRegistry");
const GuaranteedMinimumStakingPolicy = contract.fromArtifact("GuaranteedMinimumStakingPolicy");

describe('TokenGrant/Revoke', function() {

  let tokenContract, registryContract, grantContract, 
    stakingContract, stakingEscrowContract, minimumPolicy;

  const tokenOwner = accounts[0],
        grantee = accounts[1],
        beneficiary = accounts[2],
        authorizer = accounts[3],
        operator = accounts[4];

  let grantId;
  let grantStart;
  let grantAmount;
  const grantRevocable = true;
  const grantDuration = time.duration.minutes(60);
  const grantCliff = time.duration.minutes(1);

  const initializationPeriod = time.duration.minutes(10);
  let undelegationPeriod;

  let minimumStake;

  before(async () => {
    tokenContract = await KeepToken.new( {from: accounts[0]});
    grantContract = await TokenGrant.new(tokenContract.address,  {from: accounts[0]});
    registryContract = await KeepRegistry.new( {from: accounts[0]});
    const stakingContracts = await initTokenStaking(
      tokenContract.address,
      grantContract.address,
      registryContract.address,
      initializationPeriod,
      contract.fromArtifact('TokenStakingEscrow'),
      contract.fromArtifact('TokenStaking')
    )
    stakingContract = stakingContracts.tokenStaking
    stakingEscrowContract = stakingContracts.tokenStakingEscrow
    undelegationPeriod = await stakingContract.undelegationPeriod()
    minimumStake = await stakingContract.minimumStake();
    grantAmount = minimumStake.muln(10);

    await grantContract.authorizeStakingContract(stakingContract.address, {from: accounts[0]});

    minimumPolicy = await GuaranteedMinimumStakingPolicy.new(stakingContract.address);

    grantStart = await time.latest();

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
      minimumPolicy.address,
      {from: accounts[0]}
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
    await time.increaseTo(grantStart.add(time.duration.minutes(30)));

    const tx = await grantContract.revoke(grantId, { from: tokenOwner });
    const revokedAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
    const withdrawableAtRevokedTimestamp = grantAmount.mul(revokedAt.sub(grantStart)).div(grantDuration);

    const refund = grantAmount.sub(withdrawableAtRevokedTimestamp);

    const withdrawableAfter = await grantContract.withdrawable(grantId);

    await grantContract.withdrawRevoked(grantId, { from: tokenOwner });

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
    await expectRevert(
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
        minimumPolicy.address
    );
    
    await expectRevert(
      grantContract.revoke(nonRevocableGrantId, { from: tokenOwner }),
      "Grant must be revocable in the first place."
    );
  })

  it("should not allow to revoke grant multiple times", async () => {
    await grantContract.revoke(grantId, { from: tokenOwner });
  
    await expectRevert(
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
      minimumPolicy.address
      );
      
    const granteeGrantBalanceBefore = await grantContract.balanceOf.call(grantee);
    const grantManagerKeepBalanceBeforeRevoke= await tokenContract.balanceOf(tokenOwner);

    await grantContract.revoke(fullyUnlockedGrantId, { from: tokenOwner });
    await expectRevert(
      grantContract.withdrawRevoked(fullyUnlockedGrantId, { from: tokenOwner }),
      "All revoked tokens withdrawn"
    );
    
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

  it("should not be able to withdraw revoked tokens locked as stakes", async () => {
    await delegateStakeFromGrant(
      grantContract,
      stakingContract.address,
      grantee,
      operator,
      beneficiary,
      authorizer,
      minimumStake,
      grantId
    );
    await grantContract.revoke(grantId, { from: tokenOwner });

    const grantDetails = await grantContract.getGrant(grantId);
    const revokedAmount = grantDetails[3];
    const stakedAmount = grantDetails[2];

    expect(revokedAmount).to.eq.BN(grantAmount);
    expect(stakedAmount).to.eq.BN(minimumStake, "Minimum stake should be staked");

    const grantManagerKeepBalanceBeforeWithdraw = await tokenContract.balanceOf(tokenOwner);
    await grantContract.withdrawRevoked(grantId, { from: tokenOwner });
    const grantManagerKeepBalanceAfterWithdraw = await tokenContract.balanceOf(tokenOwner);
    expect(grantManagerKeepBalanceAfterWithdraw).to.eq.BN(
      grantManagerKeepBalanceBeforeWithdraw.add(grantAmount).sub(minimumStake),
      "The staked amount should be subtracted from the withdrawn amount"
    );
  });

  it("should be able to force stake cancellation and withdraw revoked tokens", async () => {
    await delegateStakeFromGrant(
      grantContract,
      stakingContract.address,
      grantee,
      operator,
      beneficiary,
      authorizer,
      minimumStake,
      grantId
    );
    await grantContract.revoke(grantId, { from: tokenOwner });

    const grantDetails = await grantContract.getGrant(grantId);
    const revokedAmount = grantDetails[3];
    const stakedAmount = grantDetails[2];

    expect(revokedAmount).to.eq.BN(grantAmount);
    expect(stakedAmount).to.eq.BN(minimumStake, "Minimum stake should be staked");

    const grantManagerKeepBalanceBeforeWithdraw = await tokenContract.balanceOf(tokenOwner);
    await grantContract.withdrawRevoked(grantId, { from: tokenOwner });
    const grantManagerKeepBalanceMidWithdraw = await tokenContract.balanceOf(tokenOwner);
    await grantContract.cancelRevokedStake(operator, { from: tokenOwner });
    await stakingEscrowContract.withdrawRevoked(operator, { from: tokenOwner });
    const grantManagerKeepBalanceAfterWithdraw = await tokenContract.balanceOf(tokenOwner);

    expect(grantManagerKeepBalanceMidWithdraw).to.eq.BN(
      grantManagerKeepBalanceBeforeWithdraw.add(grantAmount).sub(minimumStake),
      "The staked amount should be subtracted from the withdrawn amount"
    );
    expect(grantManagerKeepBalanceAfterWithdraw).to.eq.BN(
      grantManagerKeepBalanceMidWithdraw.add(minimumStake),
      "The staked amount should be withdrawn now"
    );
  });

  it("should be able to force undelegation and withdraw returned revoked tokens", async () => {
    await delegateStakeFromGrant(
      grantContract,
      stakingContract.address,
      grantee,
      operator,
      beneficiary,
      authorizer,
      minimumStake,
      grantId
    );
    await grantContract.revoke(grantId, { from: tokenOwner });

    const grantDetails = await grantContract.getGrant(grantId);
    const revokedAmount = grantDetails[3];
    const stakedAmount = grantDetails[2];

    expect(revokedAmount).to.eq.BN(grantAmount);
    expect(stakedAmount).to.eq.BN(minimumStake, "Minimum stake should be staked");

    const grantManagerKeepBalanceBeforeWithdraw = await tokenContract.balanceOf(tokenOwner);
    await grantContract.withdrawRevoked(grantId, { from: tokenOwner });
    const grantManagerKeepBalanceMidWithdraw = await tokenContract.balanceOf(tokenOwner);
    await time.increase(initializationPeriod.add(time.duration.minutes(5)));
    await grantContract.undelegateRevoked(operator, { from: tokenOwner });
    await time.increase(undelegationPeriod.add(time.duration.minutes(5)));
    await grantContract.recoverStake(operator, { from: tokenOwner });

    await stakingEscrowContract.withdrawRevoked(operator, { from: tokenOwner });
    const grantManagerKeepBalanceAfterWithdraw = await tokenContract.balanceOf(tokenOwner);

    expect(grantManagerKeepBalanceMidWithdraw).to.eq.BN(
      grantManagerKeepBalanceBeforeWithdraw.add(grantAmount).sub(minimumStake),
      "The staked amount should be subtracted from the withdrawn amount"
    );
    expect(grantManagerKeepBalanceAfterWithdraw).to.eq.BN(
      grantManagerKeepBalanceMidWithdraw.add(minimumStake),
      "The staked amount should be withdrawn now"
    );
  });
});
