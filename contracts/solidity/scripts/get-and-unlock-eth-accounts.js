/* Here we assume that the passphrase for unlocking all the accounts on
some private testnet is the same.  This is intended for use with
truffle.  Example:

truffle exec ./scripts/get-and-unlock-eth-accounts.js \
 http://eth-tx-node.default.svc.cluster.local:8545 \
 my-sweet-passphrase \
 --network keep_dev
*/

function getAccounts() {
  return new Promise((resolve, reject) => {
      web3.eth.getAccounts((error, accounts) => {
        resolve(accounts);
      });
  });
};

module.exports = async function() {

  const Web3 = require('web3');
  // We take arg 4 because "truffle exec /path/to/script" take up positions 1-3
  const web3 = new Web3(new Web3.providers.HttpProvider(process.argv[4]));

  let accounts = await getAccounts()

  console.log(`Total accounts: ${accounts.length}`)
  console.log(`---------------------------------`)

  for(let i = 1; i < accounts.length; i++) {
    let account = accounts[i]

    try {
      console.log(`\nUnlocking account: ${account}`)
      // We take arg 5 because "truffle exec /path/to/script eth-host" take up positions 1-4
      web3.personal.unlockAccount(account, `${process.argv[5]}`, 150000);
      console.log(`Account unlocked!`)
      console.log(`\n---------------------------------`)
    }
    catch(error) {
      console.log(`\nAccount: ${account} not unlocked!`)
      console.error(error)
      console.log(`\n---------------------------------`)
    }
  }
};