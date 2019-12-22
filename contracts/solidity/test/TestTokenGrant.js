import { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import expectThrow from './helpers/expectThrow';
import grantTokens from './helpers/grantTokens';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const RegistryKeeper = artifacts.require("./RegistryKeeper.sol");

contract('TestTokenGrant', function(accounts) {

  let token, registryKeeper, grantContract, stakingContract,
    amount, vestingDuration, start, cliff,
    grant_manager = accounts[0],
    account_two = accounts[1],
    grantee = accounts[2];

  before(async () => {
    token = await KeepToken.new();
    registryKeeper = await RegistryKeeper.new();
    stakingContract = await TokenStaking.new(token.address, registryKeeper.address, duration.days(30));
    grantContract = await TokenGrant.new(token.address, stakingContract.address);
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
    let grant_manager_starting_balance = await token.balanceOf.call(grant_manager);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Grant tokens
    let id = await grantTokens(grantContract, token, amount, grant_manager, account_two, vestingDuration, start, cliff, revocable);

    // Ending balances
    let grant_manager_ending_balance = await token.balanceOf.call(grant_manager);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    let account_two_grant_balance = await grantContract.balanceOf.call(account_two);

    assert.equal(grant_manager_ending_balance.eq(grant_manager_starting_balance.sub(amount)), true, "Amount should be transfered from sender balance");
    assert.equal(account_two_grant_balance.eq(amount), true, "Amount should be added to the grantee grant balance");
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance), true, "Grantee main balance should stay unchanged");

    // Should not be able to withdraw token grant (0 withdrawable amount)
    await expectThrow(grantContract.withdraw(id))

    // jump in time, third vesting duration
    await increaseTimeTo(await latestTime()+vestingDuration/3);

    // Should be able to withdraw token grant withdrawable amount
    await grantContract.withdraw(id)

    // should withdraw some of grant to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.lte(account_two_starting_balance.add(amount.div(web3.utils.toBN(2)))), true, 'Should withdraw some of the grant to the main balance')

    // jump in time, full vesting duration
    await increaseTimeTo(await latestTime()+vestingDuration);
    await grantContract.withdraw(id);

    // should withdraw full grant amount to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance.add(amount)), true, "Should withdraw full grant amount to the main balance");

    account_two_grant_balance = await grantContract.balanceOf.call(account_two);
    assert.equal(account_two_grant_balance, 0, "Grant amount should become 0");

  });

  it("token holder should be able to grant it's tokens to a grantee.", async function() {

    let grant_manager_starting_balance = await token.balanceOf.call(grant_manager);

    let id = await grantTokens(grantContract, token, amount, grant_manager, grantee, vestingDuration, start, cliff, true);

    let grant_manager_ending_balance = await token.balanceOf.call(grant_manager);

    assert.equal(grant_manager_ending_balance.eq(grant_manager_starting_balance.sub(amount)), true, "Amount should be taken out from grant manager main balance.");
    assert.equal((await grantContract.balanceOf.call(grantee)).eq(amount), true, "Amount should be added to grantee's granted balance.");

    let grant = await grantContract.getGrant(id);
    assert.equal(grant[0].eq(amount), true, "Grant should maintain a record of the granted amount.");
    assert.equal(grant[1].isZero(), true, "Grant should have 0 amount withdrawn initially.");
    assert.equal(grant[2], false, "Grant should initially be unstaked.");
    assert.equal(grant[3], false, "Grant should not be marked as revoked initially.");

    let schedule = await grantContract.getGrantVestingSchedule(id);
    assert.equal(schedule[0], grant_manager, "Grant should maintain a record of the grant manager.");
    assert.equal(schedule[1].eq(web3.utils.toBN(vestingDuration)), true, "Grant should have vesting schedule duration.");
    assert.equal(schedule[2].eq(web3.utils.toBN(start)), true, "Grant should have start time.");
    assert.equal(schedule[3].eq(web3.utils.toBN(start).add(web3.utils.toBN(cliff))), true, "Grant should have vesting schedule cliff duration.");

  });

  it("should not be able to revoke token grant.", async function() {

    // Create non revocable token grant.
    let id = await grantTokens(grantContract, token, amount, grant_manager, grantee, vestingDuration, start, cliff, false);
    await expectThrow(grantContract.revoke(id));

  });

  it("should be able to revoke revocable token grant as grant manager.", async function() {

    let grant_manager_starting_balance = await token.balanceOf.call(grant_manager);
    let grantee_starting_balance = await grantContract.balanceOf.call(grantee);

    // Create revocable token grant.
    let id = await grantTokens(grantContract, token, amount, grant_manager, grantee, vestingDuration, start, cliff, true);

    let grant_manager_ending_balance = await token.balanceOf.call(grant_manager);

    assert.equal(grant_manager_ending_balance.eq(grant_manager_starting_balance.sub(amount)), true, "Amount should be taken out from grant manager main balance.");
    assert.equal((await grantContract.balanceOf.call(grantee)).eq(grantee_starting_balance.add(amount)), true, "Amount should be added to grantee's granted balance.");

    await grantContract.revoke(id);

    grantee_starting_balance = await grantContract.balanceOf.call(grantee);
    grant_manager_starting_balance = await token.balanceOf.call(grant_manager);

    let withdrawable = await grantContract.withdrawable(id)
    assert.equal((await token.balanceOf.call(grant_manager)).eq(grant_manager_starting_balance.add(amount).sub(withdrawable)), true, "Amount should be added to grant manager main balance.");
    assert.equal((await grantContract.balanceOf.call(grantee)).eq(grantee_starting_balance.sub(amount).add(withdrawable)), true, "Amount should be taken out from grantee's granted balance.");
  });

  it("should be able to revoke the grant but no amount is refunded since duration of the vesting is over.", async function() {

    let grant_manager_starting_balance = await token.balanceOf.call(grant_manager);
    let grantee_starting_balance = await grantContract.balanceOf.call(grantee);

    // Create revocable token grant with 0 duration.
    let id = await grantTokens(grantContract, token, amount, grant_manager, grantee, duration.days(0), start, cliff, true);

    let grant_manager_ending_balance = await token.balanceOf.call(grant_manager);

    assert.equal(grant_manager_ending_balance.eq(grant_manager_starting_balance.sub(amount)), true, "Amount should be taken out from grant manager main balance.");
    assert.equal((await grantContract.balanceOf.call(grantee)).eq(grantee_starting_balance.add(amount)), true, "Amount should be added to grantee's granted balance.");

    grantee_starting_balance = await grantContract.balanceOf.call(grantee);
    await grantContract.revoke(id);

    assert.equal((await token.balanceOf.call(grant_manager)).eq(grant_manager_starting_balance.sub(amount)), true, "No amount to be returned to grant manager since vesting duration is over.");
    assert.equal((await grantContract.balanceOf.call(grantee)).eq(grantee_starting_balance), true, "Amount should stay at grantee's grant balance.");
  });

});
