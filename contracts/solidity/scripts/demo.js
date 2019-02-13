const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");

function formatAmount(amount, decimals) {
  return amount * (10 ** decimals)
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
  const tokenStaking = await TokenStaking.deployed();
  const tokenGrant = await TokenGrant.deployed();

  // Authorize contracts to work via proxy
  if (!await stakingProxy.isAuthorized(tokenStaking.address)) {
    stakingProxy.authorizeContract(tokenStaking.address);
  }
  if (!await stakingProxy.isAuthorized(tokenGrant.address)) {
    stakingProxy.authorizeContract(tokenGrant.address);
  }

  // KEEP tokens has been transfered to the first account on the list 
  // during the KEEP token contract deployment. Here, we stake those tokens.
  let staked = await token.approveAndCall(
    tokenStaking.address, 
    formatAmount(1000000, 18), 
    "", 
    {from: accounts[0]}
  ).catch((err) => {
    console.log(`could not stake KEEP tokens for ${accounts[0]}: ${err}`);
  });

  if (staked) {
    console.log(`successfully staked KEEP tokens for account ${accounts[0]}`)
  }

  // Transfer KEEP tokens to all other accounts and stake them.
  for(let i = 1; i < accounts.length; i++) {    
    let account = accounts[i]

    await token.transfer(
      account, 
      formatAmount(1000000, 18)
    ).catch((err) => { 
      console.log(`could not transfer KEEP tokens for ${account}: ${err}`); 
    });

    staked = await token.approveAndCall(
      tokenStaking.address, 
      formatAmount(1000000, 18),
      "", 
      {from: account}
    ).catch((err) => {
      console.log(`could not stake KEEP tokens for ${account}: ${err}`);
    });

    if (staked) {
      console.log(`successfully staked KEEP tokens for account ${account}`)
    }
  }

  // Grant tokens to the second account
  let amount = formatAmount(70000, 18);
  let vestingDuration = 86400*60;
  let start = web3.eth.getBlock('latest').timestamp;
  let cliff = 86400*10;
  let revocable = true;
  await token.transfer(accounts[1], formatAmount(70000,18))
  await token.approve(tokenGrant.address, amount);
  await tokenGrant.grant(amount, accounts[1], vestingDuration, start, cliff, revocable);

  // Grant tokens from the second account
  amount = formatAmount(1000, 18);
  await token.approve(tokenGrant.address, amount, {from: accounts[1]});
  await tokenGrant.grant(amount, accounts[0], vestingDuration, start, cliff, revocable, {from: accounts[1]});
};
