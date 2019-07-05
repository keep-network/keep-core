const TokenStaking = artifacts.require("TokenStaking.sol");
const KeepRandomBeaconOperator = artifacts.require("KeepRandomBeaconOperator.sol");

function getAccounts() {
  return new Promise((resolve, reject) => {
      web3.eth.getAccounts((error, accounts) => {
        resolve(accounts);
      });
  });
};

module.exports = async function () {
  const accounts = await getAccounts();
  const tokenStaking = await TokenStaking.deployed();
  const keepRandomBeacon = await KeepRandomBeaconOperator.deployed();

  const accountToUnstake = accounts[1];

  const amountToUnstake = web3.utils.toBN(200000).mul(web3.utils.toBN(10**18));

  console.log('Using account:      ', accountToUnstake);
  console.log('Stake before:       ', (await tokenStaking.stakeBalanceOf(accountToUnstake)).toString());
  console.log('Amount to unstake:  ', amountToUnstake.toString(10));
  await tokenStaking.initiateUnstake(amountToUnstake, accountToUnstake).catch((err) => {
    console.log(`could not unstake: ${err}`);
  });
  console.log('Stake after:        ', (await tokenStaking.stakeBalanceOf(accountToUnstake)).toString());

  const hasMinStake = await keepRandomBeacon.hasMinimumStake(accountToUnstake).catch((err) => {
      console.log(`could not check for minimum stake: ${err}`);
  });
  console.log('Has minimum stake?: ', hasMinStake);

  process.exit()
}
