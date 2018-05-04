const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");

function getAccounts() {
  return new Promise((resolve, reject) => {
      web3.eth.getAccounts((error, accounts) => {
        resolve(accounts);
      });
  });
};

module.exports = async function() {

  const accounts = await getAccounts();
  const token = await KeepToken.deployed();
  const stakingProxy = await StakingProxy.deployed();
  const tokenStaking = await TokenStaking.deployed();
  const tokenGrant = await TokenGrant.deployed();

  // Authorize contracts to work via proxy
  if (!await stakingProxy.isAuthorized(tokenStaking.address)) {
    stakingProxy.authorizeContract(tokenStaking.address);
  }
  if (!await stakingProxy.isAuthorized(tokenGrant.address)) {
    stakingProxy.authorizeContract(tokenGrant.address);
  }

  // Stake tokens
  await token.approveAndCall(TokenStaking.address, 1000000, "");

  // Grant tokens to a different account
  let amount = 1000000000;
  let vestingDuration = 86400*60;
  let start = web3.eth.getBlock('latest').timestamp;
  let cliff = 86400*10;
  let revocable = true;
  await token.approve(tokenGrant.address, amount);
  await tokenGrant.grant(amount, accounts[1], vestingDuration, start, cliff, revocable);
};
