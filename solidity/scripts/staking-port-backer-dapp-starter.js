const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("../stubs/OldTokenStaking.sol");
const StakingPortBacker = artifacts.require("./StakingPortBacker.sol")

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
    const stakingPortBackerContract = await StakingPortBacker.deployed();

    let owner = accounts[0];
    
    for (let i = 5; i < 10; i++) {
      let operator = accounts[i]
      let beneficiary = accounts[i]
      let authorizer = accounts[i]

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

      if (staked) {
        await stakingPortBackerContract.allowOperator(operator, { from: owner })
          .catch(err => console.log(`could not allowOperator for ${operator}`, err));
        console.log(`successfully staked KEEP tokens for account ${operator}`)
      }
    }
    const allTokens = await token.balanceOf(owner)
    await token.transfer(stakingPortBackerContract.address, allTokens.divn(2), {from: owner})
  } catch (err) {
    console.error('unexpected error:', err)
    process.exit(1)
  }

  process.exit();
};
