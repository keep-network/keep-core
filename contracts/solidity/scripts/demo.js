const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
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
  const tokenStaking = await TokenStaking.deployed();
  const tokenGrant = await TokenGrant.deployed();

  let owner = accounts[0]; // The address of an owner of the staked tokens.
  let magpie = accounts[0]; // The address where the rewards for participation are sent.

  // Stake delegate tokens for first 5 accounts as operators,
  // including the first account where owner operating for themself.
  for(let i = 0; i < 5; i++) {
    let operator = accounts[i]

    let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator)).substr(2), 'hex');
    let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');

    staked = await token.approveAndCall(
      tokenStaking.address, 
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

  // Create a demo accounts with tokens but without any operators
  await token.transfer(accounts[5], formatAmount(100000, 18), {from: accounts[0]})

  // Grant tokens to the stake owner account
  let amount = formatAmount(70000, 18);
  let vestingDuration = web3.utils.toBN(86400).mul(web3.utils.toBN(60));
  let start = (await web3.eth.getBlock('latest')).timestamp;
  let cliff = web3.utils.toBN(86400).mul(web3.utils.toBN(10));
  let revocable = true;
  await token.approve(tokenGrant.address, amount, {from: accounts[0]});
  await tokenGrant.grant(amount, accounts[5], vestingDuration, start, cliff, revocable, {from: accounts[0]});

  // Grant tokens from the stake owner account
  amount = formatAmount(1000, 18);
  await token.approve(tokenGrant.address, amount, {from: accounts[5]});
  await tokenGrant.grant(amount, accounts[6], vestingDuration, start, cliff, revocable, {from: accounts[5]});

  process.exit();
};
