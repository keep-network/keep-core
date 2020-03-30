import { duration, increaseTimeTo } from '../helpers/increaseTime';
import latestTime from '../helpers/latestTime';
import expectThrow from '../helpers/expectThrow';
import grantTokens from '../helpers/grantTokens';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const Registry = artifacts.require("./Registry.sol");

contract('TokenGrant', function(accounts) {

  let token, registry, grantContract, stakingContract,
    amount, unlockingDuration, start, cliff,
    grant_manager = accounts[0],
    account_two = accounts[1],
    grantee = accounts[2];

  before(async () => {
    token = await KeepToken.new();
    registry = await Registry.new();
    stakingContract = await TokenStaking.new(token.address, registry.address, duration.days(1), duration.days(30));
    grantContract = await TokenGrant.new(token.address);
    await grantContract.authorizeStakingContract(stakingContract.address);
    amount = web3.utils.toBN(100);
    unlockingDuration = duration.days(30);
    start = await latestTime();
    cliff = duration.days(0);
  });

  it("should grant tokens correctly", async function() {

    let amount = web3.utils.toBN(1000000000);
    let unlockingDuration = duration.days(30);
    let start = await latestTime();
    let cliff = duration.days(10);
    let revocable = true;

    // Starting balances
    let grant_manager_starting_balance = await token.balanceOf.call(grant_manager);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Grant tokens
    let id = await grantTokens(grantContract, token, amount, grant_manager, account_two, unlockingDuration, start, cliff, revocable);

    // Ending balances
    let grant_manager_ending_balance = await token.balanceOf.call(grant_manager);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    let account_two_grant_balance = await grantContract.balanceOf.call(account_two);

    assert.equal(grant_manager_ending_balance.eq(grant_manager_starting_balance.sub(amount)), true, "Amount should be transfered from sender balance");
    assert.equal(account_two_grant_balance.eq(amount), true, "Amount should be added to the grantee grant balance");
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance), true, "Grantee main balance should stay unchanged");

    // Should not be able to withdraw token grant (0 withdrawable amount)
    await expectThrow(grantContract.withdraw(id))

    // jump in time, third unlocking duration
    await increaseTimeTo(await latestTime()+unlockingDuration/3);

    // Should be able to withdraw token grant withdrawable amount
    await grantContract.withdraw(id)

    // should withdraw some of grant to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.lte(account_two_starting_balance.add(amount.div(web3.utils.toBN(2)))), true, 'Should withdraw some of the grant to the main balance')

    // jump in time, full unlocking duration
    await increaseTimeTo(await latestTime()+unlockingDuration);
    await grantContract.withdraw(id);

    // should withdraw full grant amount to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance.add(amount)), true, "Should withdraw full grant amount to the main balance");

    account_two_grant_balance = await grantContract.balanceOf.call(account_two);
    assert.equal(account_two_grant_balance, 0, "Grant amount should become 0");

  });

  it("token holder should be able to grant it's tokens to a grantee.", async function() {

    let grant_manager_starting_balance = await token.balanceOf.call(grant_manager);

    let id = await grantTokens(grantContract, token, amount, grant_manager, grantee, unlockingDuration, start, cliff, true);

    let grant_manager_ending_balance = await token.balanceOf.call(grant_manager);

    assert.equal(grant_manager_ending_balance.eq(grant_manager_starting_balance.sub(amount)), true, "Amount should be taken out from grant manager main balance.");
    assert.equal((await grantContract.balanceOf.call(grantee)).eq(amount), true, "Amount should be added to grantee's granted balance.");

    let grant = await grantContract.getGrant(id);
    assert.equal(grant[0].eq(amount), true, "Grant should maintain a record of the granted amount.");
    assert.equal(grant[1].isZero(), true, "Grant should have 0 amount withdrawn initially.");
    assert.equal(grant[2], false, "Grant should initially be undelegated.");
    assert.equal(grant[3], false, "Grant should not be marked as revoked initially.");

    let schedule = await grantContract.getGrantUnlockingSchedule(id);
    assert.equal(schedule[0], grant_manager, "Grant should maintain a record of the grant manager.");
    assert.equal(schedule[1].eq(web3.utils.toBN(unlockingDuration)), true, "Grant should have unlocking schedule duration.");
    assert.equal(schedule[2].eq(web3.utils.toBN(start)), true, "Grant should have start time.");
    assert.equal(schedule[3].eq(web3.utils.toBN(start).add(web3.utils.toBN(cliff))), true, "Grant should have unlocking schedule cliff duration.");

  });
});
