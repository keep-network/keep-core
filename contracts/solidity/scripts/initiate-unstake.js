const TokenStaking = artifacts.require("TokenStaking");

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

  const accountToUnstake = accounts[1];
  const amountToUnstake = web3.utils.toBN(200000).mul(web3.utils.toBN(10**18));

  console.log('Using account: ' + accountToUnstake);
  console.log('Stake before:', (await tokenStaking.stakeBalanceOf(accountToUnstake)).toString());
  await tokenStaking.initiateUnstake(amountToUnstake, accountToUnstake).catch((err) => {
    console.log(`could not unstake: ${err}`);
  });
  console.log('Stake after:', (await tokenStaking.stakeBalanceOf(accountToUnstake)).toString());

  process.exit()
}