import { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');

contract('TestTokenGrantsTransfer', function(accounts) {

  let token, grantContract, stakingProxy,
    amount, vestingDuration, start, cliff,
    sourceId,
    account_one = accounts[0],
    beneficiary = accounts[1],
    newBeneficiary = accounts[2];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    grantContract = await TokenGrant.new(token.address, stakingProxy.address, duration.days(30));
    amount = web3.utils.toBN(600);
    vestingDuration = duration.days(30);
    start = await latestTime();
    cliff = duration.days(0);

    await token.approve(grantContract.address, amount, {from: account_one});
    sourceId = (await grantContract.grant(amount, beneficiary, vestingDuration,
      start, cliff, true, {from: account_one})).logs[0].args.id.toNumber()
  });

  it("should be able to transfer token grants unvested amount.", async function() {

    let amountToTransfer = web3.utils.toBN(300);
    let targetId = (await grantContract.transfer(sourceId, newBeneficiary, amountToTransfer, {from: account_one})).logs[0].args.id.toNumber();

    let sourceGrant = await grantContract.getGrant(sourceId);
    assert.isTrue(sourceGrant[0].eq(amount.sub(amountToTransfer)), "Source grant should have correct amount removed.");

    let sourceGrantVestedAmount = await grantContract.grantedAmount(sourceId);
    assert.isTrue(sourceGrantVestedAmount.isZero(), "Source grant should not have vested amount.");

    let targetGrant = await grantContract.getGrant(targetId);
    assert.isTrue(targetGrant[0].eq(amountToTransfer), "Target grant should have correct amount added.");

    let targetGrantVestedAmount = await grantContract.grantedAmount(sourceId);
    assert.isTrue(targetGrantVestedAmount.isZero(), "Target grant should not have vested amount.");

  });

  it("should be able to transfer token grants from vested and unvested amounts.", async function() {

    let amountToTransfer = web3.utils.toBN(450); // 300 should be taken from the releasable balance and 150 from unvested.

    // jump in time, half of the amount should be vested and releasable.
    await increaseTimeTo(await latestTime()+duration.days(15));

    let sourceGrantUnreleasedAmount = await grantContract.unreleasedAmount(sourceId);
    assert.isTrue(sourceGrantUnreleasedAmount.eq(web3.utils.toBN(300)), "Source grant should have the vested amount as unreleased.");

    // Execute the transfer.
    let targetId = (await grantContract.transfer(sourceId, newBeneficiary, amountToTransfer, {from: account_one})).logs[0].args.id.toNumber();

    let sourceGrant = await grantContract.getGrant(sourceId);
    assert.isTrue(sourceGrant[0].eq(web3.utils.toBN(150)), "Source grant should have correct amount removed.");

    let sourceGrantVestedAmount = await grantContract.grantedAmount(sourceId);
    assert.isTrue(sourceGrantVestedAmount.isZero(), "Source grant should not have vested amount left.");

    let targetGrant = await grantContract.getGrant(targetId);
    assert.isTrue(targetGrant[0].eq(amountToTransfer), "Target grant should have correct amount added.");

    let targetGrantVestedAmount = await grantContract.grantedAmount(targetId);
    assert.isTrue(targetGrantVestedAmount.eq(web3.utils.toBN(300)), "Target grant should have the rest of the transferred amount as vested.");

    let targetGrantUnreleasedAmount = await grantContract.unreleasedAmount(targetId);
    assert.isTrue(targetGrantUnreleasedAmount.eq(web3.utils.toBN(300)), "Target grant should have the vested amount as releasable.");

  });

  it("should be able to transfer token grants from a partially released grants unvested balance.", async function() {


    // jump in time, half of the amount should be vested and releasable.
    await increaseTimeTo(await latestTime()+duration.days(15));

    // release available vested amount
    await grantContract.release(sourceId);
    let sourceGrantUnreleasedAmount = await grantContract.unreleasedAmount(sourceId);
    assert.isTrue(sourceGrantUnreleasedAmount.isZero(), "Source grant should not have any releasable amounts.");

    // Should fail to transfer since not enough amount left after grant release.
    let amountToTransfer = web3.utils.toBN(450);
    await exceptThrow(grantContract.transfer(sourceId, newBeneficiary, amountToTransfer, {from: account_one}));

    // Execute the transfer with lower amount.
    amountToTransfer = web3.utils.toBN(150);
    let targetId = (await grantContract.transfer(sourceId, newBeneficiary, amountToTransfer, {from: account_one})).logs[0].args.id.toNumber();

    let sourceGrant = await grantContract.getGrant(sourceId);
    assert.isTrue(sourceGrant[0].eq(web3.utils.toBN(150)), "Source grant should have correct amount removed.");

    let sourceGrantVestedAmount = await grantContract.grantedAmount(sourceId);
    assert.isTrue(sourceGrantVestedAmount.isZero(), "Source grant should not have vested amount left.");

    let targetGrant = await grantContract.getGrant(targetId);
    assert.isTrue(targetGrant[0].eq(amountToTransfer), "Target grant should have correct amount added.");

    let targetGrantVestedAmount = await grantContract.grantedAmount(targetId);
    assert.isTrue(targetGrantVestedAmount.isZero(), "Target grant should not have any vested amount.");

    let targetGrantUnreleasedAmount = await grantContract.unreleasedAmount(targetId);
    assert.isTrue(targetGrantUnreleasedAmount.isZero(), "Target grant should not have any releasable amounts.");

  });

});
