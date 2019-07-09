import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');

contract('TestTokenGrant', function(accounts) {

  let token, grantContract, stakingContract,
    amount, vestingDuration, start, cliff,
    account_one = accounts[0],
    account_two = accounts[1],
    beneficiary = accounts[2];

  before(async () => {
    token = await KeepToken.new();
    stakingContract = await TokenStaking.new(token.address, duration.days(30));
    grantContract = await TokenGrant.new(token.address, stakingContract.address, duration.days(30));
    amount = web3.utils.toBN(100);
    vestingDuration = duration.days(30);
    start = await latestTime();
    cliff = duration.days(0);
  });

  it("should grant tokens correctly", async function() {

    let amount = web3.utils.toBN(1000000000);
    let vestingDuration = duration.days(30);
    let start = await latestTime();
    let cliff = duration.days(10);
    let revocable = true;

    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Grant tokens
    await token.approve(grantContract.address, amount, {from: account_one});
    let id = await grantContract.grant(amount, account_two, vestingDuration, 
      start, cliff, revocable, {from: account_one}).then((result)=>{
      // Look for CreatedTokenGrant event in transaction receipt and get vesting id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "CreatedTokenGrant") {
          return log.args.id.toNumber();
        }
      }
    })

    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    let account_two_grant_balance = await grantContract.totalBalanceOf.call(account_two);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount should be transfered from sender balance");
    assert.equal(account_two_grant_balance.eq(amount), true, "Amount should be added to the beneficiary grant balance");
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance), true, "Beneficiary main balance should stay unchanged");

    // Should not be able to release token grant (0 unreleased amount)
    await exceptThrow(grantContract.release(id))

    // jump in time, third vesting duration
    await increaseTimeTo(await latestTime()+vestingDuration/3);

    // Should be able to release token grant unreleased amount
    await grantContract.release(id)

    // should release some of grant to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.lte(account_two_starting_balance.add(amount.div(web3.utils.toBN(2)))), true, 'Should release some of the grant to the main balance')

    // jump in time, full vesting duration
    await increaseTimeTo(await latestTime()+vestingDuration);
    await grantContract.release(id);

    // should release full grant amount to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance.add(amount)), true, "Should release full grant amount to the main balance");

    account_two_grant_balance = await grantContract.totalBalanceOf.call(account_two);
    assert.equal(account_two_grant_balance, 0, "Grant amount should become 0");

  });

  it("token holder should be able to grant it's tokens to a beneficiary.", async function() {

    let account_one_starting_balance = await token.balanceOf.call(account_one);

    await token.approve(grantContract.address, amount, {from: account_one});
    let id = (await grantContract.grant(amount, beneficiary, vestingDuration,
      start, cliff, true, {from: account_one})).logs[0].args.id.toNumber()

    let account_one_ending_balance = await token.balanceOf.call(account_one);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount should be taken out from grant creator main balance.");
    assert.equal((await grantContract.totalBalanceOf.call(beneficiary)).eq(amount), true, "Amount should be added to beneficiary's granted balance.");

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
    let beneficiary_starting_balance = await grantContract.totalBalanceOf.call(beneficiary);

    // Create revocable token grant.
    await token.approve(grantContract.address, amount, {from: account_one});
    let id = (await grantContract.grant(amount, beneficiary, vestingDuration,
      start, cliff, true, {from: account_one})).logs[0].args.id.toNumber()

    let account_one_ending_balance = await token.balanceOf.call(account_one);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount should be taken out from grant creator main balance.");
    assert.equal((await grantContract.totalBalanceOf.call(beneficiary)).eq(beneficiary_starting_balance.add(amount)), true, "Amount should be added to beneficiary's granted balance.");

    await grantContract.revoke(id);

    beneficiary_starting_balance = await grantContract.totalBalanceOf.call(beneficiary);
    account_one_starting_balance = await token.balanceOf.call(account_one);
    
    let unreleasedAmount = await grantContract.unreleasedAmount(id)
    assert.equal((await token.balanceOf.call(account_one)).eq(account_one_starting_balance.add(amount).sub(unreleasedAmount)), true, "Amount should be added to grant creator main balance.");
    assert.equal((await grantContract.totalBalanceOf.call(beneficiary)).eq(beneficiary_starting_balance.sub(amount).add(unreleasedAmount)), true, "Amount should be taken out from beneficiary's granted balance.");
  });

  it("should be able to revoke the grant but no amount is refunded since duration of the vesting is over.", async function() {

    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let beneficiary_starting_balance = await grantContract.totalBalanceOf.call(beneficiary);

    // Create revocable token grant with 0 duration.
    await token.approve(grantContract.address, amount, {from: account_one});
    let id = (await grantContract.grant(amount, beneficiary, duration.days(0),
      start, cliff, true, {from: account_one})).logs[0].args.id.toNumber()

    let account_one_ending_balance = await token.balanceOf.call(account_one);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount should be taken out from grant creator main balance.");
    assert.equal((await grantContract.totalBalanceOf.call(beneficiary)).eq(beneficiary_starting_balance.add(amount)), true, "Amount should be added to beneficiary's granted balance.");

    beneficiary_starting_balance = await grantContract.totalBalanceOf.call(beneficiary);
    await grantContract.revoke(id);

    assert.equal((await token.balanceOf.call(account_one)).eq(account_one_starting_balance.sub(amount)), true, "No amount to be returned to grant creator since vesting duration is over.");
    assert.equal((await grantContract.totalBalanceOf.call(beneficiary)).eq(beneficiary_starting_balance), true, "Amount should stay at beneficiary's grant balance.");
  });

});
