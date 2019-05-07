const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const Staking = artifacts.require("./Staking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");

function formatAmount(amount, decimals) {
  return web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals)))
}

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
  const staking = await Staking.deployed();
  const tokenGrant = await TokenGrant.deployed();

  // Authorize contracts to work via proxy
  if (!await stakingProxy.isAuthorized(staking.address)) {
    stakingProxy.authorizeContract(staking.address);
  }
  if (!await stakingProxy.isAuthorized(tokenGrant.address)) {
    stakingProxy.authorizeContract(tokenGrant.address);
  }

  let owner = accounts[0]; // The address of an owner of the staked tokens.
  let magpie = accounts[0]; // The address where the rewards for participation are sent.

  // Stake delegate tokens for each account as an operator,
  // including the first account where owner operating for themself.
  for(let i = 0; i < accounts.length; i++) {
    let operator = accounts[i]

    let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator)).substr(2), 'hex');
    let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');

    staked = await token.approveAndCall(
      staking.address, 
      formatAmount(1000000, 18),
      delegation,
      {from: owner}
    ).catch((err) => {
      console.log(`could not stake KEEP tokens for ${operator}: ${err}`);
    });

    if (staked) {
      console.log(`successfully staked KEEP tokens for account ${operator}`)
    }
  }

  // Grant tokens to the second account
  let amount = formatAmount(70000, 18);
  let vestingDuration = web3.utils.toBN(86400).mul(web3.utils.toBN(60));
  let start = (await web3.eth.getBlock('latest')).timestamp;
  let cliff = web3.utils.toBN(86400).mul(web3.utils.toBN(10));
  let revocable = true;
  await token.transfer(accounts[1], formatAmount(70000,18))
  await token.approve(tokenGrant.address, amount);
  await tokenGrant.grant(amount, accounts[1], vestingDuration, start, cliff, revocable);

  // Grant tokens from the second account
  amount = formatAmount(1000, 18);
  await token.approve(tokenGrant.address, amount, {from: accounts[1]});
  await tokenGrant.grant(amount, accounts[0], vestingDuration, start, cliff, revocable, {from: accounts[1]});

  process.exit();
};
