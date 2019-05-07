import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');

contract('TestTokenGrants', function(accounts) {

  let token, grantContract, stakingProxy,
    amount, vestingDuration, start, cliff,
    account_one = accounts[0],
    beneficiary = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    grantContract = await TokenGrant.new(token.address);
    amount = web3.utils.toBN(100);
    vestingDuration = duration.days(30);
    start = await latestTime();
    cliff = duration.days(0);
  });


  it("token holder should be able to grant it's tokens to a beneficiary.", async function() {

    let account_one_starting_balance = await token.balanceOf.call(account_one);

    await token.approve(grantContract.address, amount, {from: account_one});
    let id = (await grantContract.grant(amount, beneficiary, vestingDuration,
      start, cliff, true, {from: account_one})).logs[0].args.id.toNumber()

    let account_one_ending_balance = await token.balanceOf.call(account_one);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount should be taken out from grant creator main balance.");
    assert.equal((await grantContract.balanceOf.call(beneficiary)).eq(amount), true, "Amount should be added to beneficiary's granted balance.");

    let grant = await grantContract.getGrant(id);
    assert.equal(grant[0].eq(amount), true, "Grant should maintain a record of the granted amount.");
    assert.equal(grant[1].isZero(), true, "Grant should have 0 amount released initially.");
    assert.equal(grant[2], false, "Grant should initially be unstaked.");
    assert.equal(grant[3], false, "Grant should not be marked as revoked initially.");

    let schedule = await grantContract.getGrantVestingSchedule(id);
    assert.equal(schedule[0], account_one, "Grant should maintain a record of the creator.");
    assert.equal(schedule[1].eq(web3.utils.toBN(vestingDuration)), true, "Grant should have vesting schedule duration.");
    assert.equal(schedule[2].eq(web3.utils.toBN(start)), true, "Grant should have start time.");
    assert.equal(schedule[3].eq(web3.utils.toBN(start).add(web3.utils.toBN(cliff))), true, "Grant should have vesting schedule cliff duration.");

  });

  it("should not be able to revoke token grant.", async function() {

    // Create non revocable token grant.
    await token.approve(grantContract.address, amount, {from: account_one});
    let id = (await grantContract.grant(amount, beneficiary, vestingDuration,
      start, cliff, false, {from: account_one})).logs[0].args.id.toNumber()

    await exceptThrow(grantContract.revoke(id));

  });

  it("should be able to revoke revocable token grant as grant owner.", async function() {

    let account_one_starting_balance = await token.balanceOf.call(account_one);

    // Create revocable token grant.
    await token.approve(grantContract.address, amount, {from: account_one});
    let id = (await grantContract.grant(amount, beneficiary, vestingDuration,
      start, cliff, true, {from: account_one})).logs[0].args.id.toNumber()

    let account_one_ending_balance = await token.balanceOf.call(account_one);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount should be taken out from grant creator main balance.");
    assert.equal((await grantContract.balanceOf.call(beneficiary)).eq(amount), true, "Amount should be added to beneficiary's granted balance.");

    await grantContract.revoke(id);

    assert.equal((await token.balanceOf.call(account_one)).eq(account_one_starting_balance), true, "Amount should be added to grant creator main balance.");
    assert.equal((await grantContract.balanceOf.call(beneficiary)).isZero(), true, "Amount should be taken out from beneficiary's granted balance.");
  });

  it("should be able to revoke the grant but no amount is refunded since duration of the vesting is over.", async function() {

    let account_one_starting_balance = await token.balanceOf.call(account_one);

    // Create revocable token grant with 0 duration.
    await token.approve(grantContract.address, amount, {from: account_one});
    let id = (await grantContract.grant(amount, beneficiary, duration.days(0),
      start, cliff, true, {from: account_one})).logs[0].args.id.toNumber()

    let account_one_ending_balance = await token.balanceOf.call(account_one);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount should be taken out from grant creator main balance.");
    assert.equal((await grantContract.balanceOf.call(beneficiary)).eq(amount), true, "Amount should be added to beneficiary's granted balance.");

    await grantContract.revoke(id);

    assert.equal((await token.balanceOf.call(account_one)).eq(account_one_starting_balance.sub(amount)), true, "No amount to be returned to grant creator since vesting duration is over.");
    assert.equal((await grantContract.balanceOf.call(beneficiary)).eq(amount), true, "Amount should stay at beneficiary's grant balance.");
  });

});
