import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');

contract('KeepToken', function(accounts) {

  let account_one = accounts[0];
  let account_two = accounts[1];

  it("should send tokens correctly", async function() {
    let token = await KeepToken.deployed();
    let amount = 1000000000;

    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Send tokens
    await token.transfer(account_two, amount, {from: account_one});

    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_two_ending_balance = await token.balanceOf.call(account_two);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber() - amount, "Amount wasn't correctly taken from the sender");
    assert.equal(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber() + amount, "Amount wasn't correctly sent to the receiver");

  });

  it("should vest tokens correctly", async function() {
    let token = await KeepToken.deployed();
    let vestingAmount = 1000000000;
    let vestingDuration = duration.days(30);
    let vestingStart = latestTime();
    let vestingCliff = duration.days(10);
    let vestingRevocable = true;
    
    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Vest tokens
    let vestingId = await token.vest(vestingAmount, account_two, vestingDuration, vestingStart, vestingCliff, vestingRevocable, {from: account_one}).then((result)=>{
      // Look for NewVesting event in transaction receipt and get vesting id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "NewVesting") {
          return log.args.id.toNumber();
        }
      }
    })
    
    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    let account_two_vesting_balance = await token.vestingBalanceOf.call(account_two);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber() - vestingAmount, "Vesting amount should be transfered from sender balance");
    assert.equal(account_two_vesting_balance.toNumber(), vestingAmount, "Vesting amount should be added to the beneficiary vesting balance");
    assert.equal(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber(), "Beneficiary main balance should stay unchanged");
    
    // jump in time, third vesting duration
    await increaseTimeTo(latestTime()+vestingDuration/3);
    await token.releaseVesting(vestingId);
    
    // should release some of vesting to the main balance
    account_two_ending_balance= await token.balanceOf.call(account_two);
    assert.isBelow(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber() + vestingAmount/2, 'Should release some of the vesting to the main balance')

    // jump in time, full vesting duration
    await increaseTimeTo(latestTime()+vestingDuration);
    await token.releaseVesting(vestingId);

    // should release full vesting amount to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber() + vestingAmount, "Should release full vesting amount to the main balance");

    account_two_vesting_balance = await token.vestingBalanceOf.call(account_two);
    assert.equal(account_two_vesting_balance.toNumber(), 0, "Vesting amount should become 0");
    
  });
});
