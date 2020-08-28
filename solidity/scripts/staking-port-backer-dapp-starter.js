const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("../stubs/OldTokenStaking.sol");
const StakingPortBacker = artifacts.require("./StakingPortBacker.sol")
const GuaranteedMinimumStakingPolicy = artifacts.require("./GuaranteedMinimumStakingPolicy.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol")
const ManagedGrantFactory = artifacts.require("./ManagedGrantFactory.sol");
const ManagedGrant = artifacts.require("./ManagedGrant.sol");

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
    const policy = await GuaranteedMinimumStakingPolicy.deployed();
    const grantContract = await TokenGrant.deployed();
    const managedGrantFactoryContract = await ManagedGrantFactory.deployed();

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

    const grantee = owner
    const duration = 1;
    const start =  (await web3.eth.getBlock("latest")).timestamp;
    const cliff = 1;
    const revocable = false;
    const policyAddress = policy.address;

    const grantExtraData = web3.eth.abi.encodeParameters(
      ["address", "address", "uint256", "uint256", "uint256", "bool", "address"],
      [
        owner,
        grantee,
        duration,
        start,
        cliff,
        revocable,
        policyAddress,
      ]
    )
    
    const managedGrantExtraData = web3.eth.abi.encodeParameters(
      ["address", "uint256", "uint256", "uint256", "bool", "address"],
      [
        grantee,
        duration,
        start,
        cliff,
        revocable,
        policyAddress,
      ]
    )

    // create default grant
    await token.approveAndCall(grantContract.address, formatAmount(12300000, 18), grantExtraData, { from: owner })
    const tokenGrantCreatedEvent = (await grantContract.getPastEvents())[0]
    const grantId = tokenGrantCreatedEvent.args['id']
    // Stake from a default grant
    for(let i = 10; i < 15; i++) {
      let operator = accounts[i]
      let beneficiary = accounts[i]
      let authorizer = accounts[i]

      let delegation = '0x' + Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operator.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ]).toString('hex');

      staked = await grantContract.stake(
        grantId,
        tokenStaking.address,
        formatAmount(100000, 18),
        delegation,
        { from: grantee }
      ).catch((err) => {
        console.log(`could not stake KEEP tokens from a grant for ${operator}: ${err}`);
      });

      await stakingPortBackerContract.allowOperator(operator, { from: owner })
        .catch(err => console.log(`could not allowOperator for ${operator}`, err));
      console.log(`successfully staked KEEP tokens from grant for account ${operator}`)
    }

    // Create managed grant
    await token.approveAndCall(managedGrantFactoryContract.address, formatAmount(12300000, 18), managedGrantExtraData, { from: owner })
    // Get the address of managed grant contract from an event.
    const managedGrant1Event = (await managedGrantFactoryContract.getPastEvents())[0]
    const managedGrant1Address = managedGrant1Event.args['grantAddress'];
    const managedGrant1 = await ManagedGrant.at(managedGrant1Address);

    // delegate stake from a managed grant
    for(let i = 15; i < 20; i++) {
      let operator = accounts[i]
      let beneficiary = accounts[i]
      let authorizer = accounts[i]

      let delegation = '0x' + Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operator.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ]).toString('hex');

      staked = await managedGrant1.stake(tokenStaking.address, formatAmount(100000, 18), delegation, { from: grantee }).catch((err) => {
        console.log(`could not stake KEEP tokens from a managed grant for ${operator}: ${err}`);
      });

     
      await stakingPortBackerContract.allowOperator(operator, { from: owner })
        .catch(err => console.log(`could not allowOperator for ${operator}`, err));
      console.log(`successfully staked KEEP tokens from manage grant for account ${operator}`)
    }
  } catch (err) {
    console.error('unexpected error:', err)
    process.exit(1)
  }

  process.exit();
};
