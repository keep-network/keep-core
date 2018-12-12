import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
import abi from 'ethereumjs-abi';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const StakingManager = artifacts.require('./examples/StakingManagerExample.sol');


contract('TestStakeAuthorizedTransfer', function(accounts) {

  let token, stakingProxy, stakingContract, stakingManager,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})
    await token.approveAndCall(stakingContract.address, 10000000, "", {from: account_one});
    stakingManager = await StakingManager.new();
  });

  it("should not be able to move staked tokens if not authorized", async function() {
    await exceptThrow(stakingManager.slash(account_one, 100, {from: account_two}));
  });

  it("should be able to authorize a staking manager contract", async function() {
    let hashedManagerAddress = '0x' + abi.soliditySHA3(["address"], [stakingManager.address]).toString('hex');
    let signature = '0x' + web3.eth.sign(account_one, hashedManagerAddress).substr(2);
    await stakingContract.authorizeManager(stakingManager.address, signature);

    assert.equal(await stakingManager.manageableContractFor(account_one), stakingContract.address, "Address should match the one authorized by a staker.");
    assert.equal(await stakingContract.isManagerAuthorizedFor(account_one, stakingManager.address), true, "Address should match the one authorized by a staker.");

    // Staking manager can move staked tokens if authorized
    await stakingManager.slash(account_one, 100, {from: account_two});
    assert.equal(await stakingContract.stakeBalanceOf(stakingManager.address), 100, "Staked balance should be moved.");

  });
});
