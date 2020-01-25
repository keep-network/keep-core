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
  // accounts[1]...[4] Operators for owner delegated stake and receivers of the rewards.

  // Token Grants demo accounts
  let grantee = accounts[0];
  let grantManager = accounts[5];
  let granteeOperator = accounts[6];

  // Stake delegate tokens for first 5 accounts as operators,
  // including the first account where owner operating for themself.
  for(let i = 0; i < 5; i++) {
    let operator = accounts[i]
    let magpie = accounts[i] // The address where the rewards for participation are sent.
    let authorizer = accounts[i] // Authorizer authorizes operator contracts the staker operates on.

    // The owner provides to the contract a magpie address and the operator address. 
    let delegation = '0x' + Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operator.substr(2), 'hex'),
      Buffer.from(authorizer.substr(2), 'hex')
    ]).toString('hex');

    staked = await token.approveAndCall(
      tokenStaking.address, 
      formatAmount(20000000, 18),
      delegation,
      {from: owner}
    ).catch((err) => {
      console.log(`could not stake KEEP tokens for ${operator}: ${err}`);
    });

    if (staked) {
      console.log(`successfully staked KEEP tokens for account ${operator}`)
    }
  }

  // Make sure grant manager has some tokens to be able to create a grant.
  await token.transfer(grantManager, formatAmount(100000, 18), {from: owner})

  // Grant tokens to grantee.
  let amount = formatAmount(70000, 18);
  let vestingDuration = web3.utils.toBN(86400).mul(web3.utils.toBN(60));
  let start = (await web3.eth.getBlock('latest')).timestamp;
  let cliff = web3.utils.toBN(86400).mul(web3.utils.toBN(10));
  let revocable = false; // Can not stake revocable token grants. More info in RFC14 

  await token.approveAndCall(
    tokenGrant.address,
    amount,
    Buffer.concat([
      Buffer.from(grantee.substr(2), 'hex'),
      web3.utils.toBN(vestingDuration).toBuffer('be', 32),
      web3.utils.toBN(start).toBuffer('be', 32),
      web3.utils.toBN(cliff).toBuffer('be', 32),
      Buffer.from(revocable ? "01" : "00", 'hex'),
    ]),
    {from: grantManager}
  )
  let grantId = (await tokenGrant.getPastEvents())[0].args[0].toNumber()

  // Operator must sign grantee and grant contract address since grant contract becomes the owner during staking.
  let signature1 = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(tokenGrant.address), granteeOperator)).substr(2), 'hex');
  let signature2 = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(grantee), granteeOperator)).substr(2), 'hex');
  let delegation = Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature1, signature2]);
  await tokenGrant.stake(grantId, tokenStaking.address, 10, delegation, {from: grantee});

  process.exit();
};
