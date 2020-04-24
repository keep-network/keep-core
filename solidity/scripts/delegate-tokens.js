const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");

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

module.exports = async function () {
  try {
    const accounts = await getAccounts();
    const token = await KeepToken.deployed();
    const tokenStaking = await TokenStaking.deployed();
    const operatorContract = await KeepRandomBeaconOperator.deployed();

    let owner = accounts[0]; // The address of an owner of the staked tokens.
    // accounts[1]...[4] Operators for owner delegated stake and receivers of the rewards.

    // Stake delegate tokens for first 5 accounts as operators,
    // including the first account where owner operating for themself.
    for (let i = 0; i < 5; i++) {
      let operator = accounts[i]
      let beneficiary = accounts[i] // The address where the rewards for participation are sent.
      let authorizer = accounts[i] // Authorizer authorizes operator contracts the staker operates on.

      // The owner provides to the contract a beneficiary address and the operator address. 
      let delegation = '0x' + Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operator.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ]).toString('hex');

      staked = await token.approveAndCall(
        tokenStaking.address,
        formatAmount(20000000, 18),
        delegation,
        { from: owner }
      ).catch((err) => {
        console.log(`could not stake KEEP tokens for ${operator}: ${err}`);
      });

      await tokenStaking.authorizeOperatorContract(operator, operatorContract.address, { from: authorizer });

      if (staked) {
        console.log(`successfully staked KEEP tokens for account ${operator}`)
      }
    }
  } catch (err) {
    console.error('unexpected error:', err)
    process.exit(1)
  }

  process.exit();
};
