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

const getDelegationExtraData = (account) => {
  return '0x' + Buffer.concat([
    Buffer.from(account.substr(2), 'hex'),
    Buffer.from(account.substr(2), 'hex'),
    Buffer.from(account.substr(2), 'hex')
  ]).toString('hex');
}

function GrantStakingStrategy(grantId, from) {
  return {
    stake: async (stakingContractAddress, amount, delegation, operator) => {
      const grantContract = await TokenGrant.deployed();

      await grantContract.stake(
        grantId,
        stakingContractAddress,
        amount,
        delegation,
        { from }
      ).catch((err) => {
        console.log(`could not stake KEEP tokens from a grant for ${operator}: ${err}`);
      });

      console.log(`successfully staked KEEP tokens from a grant (${grantId.toString()}) for ${operator}`)
    }
  }
}

function ManagedGrantStakingStrategy(managedGrantAddress, from) {
  return {
    stake: async (stakingContractAddress, amount, delegation, operator) => {
      const managedGrantContract = await ManagedGrant.at(managedGrantAddress);
      await managedGrantContract.stake(
        stakingContractAddress,
        amount,
        delegation,
        { from }
      ).catch((err) => {
        console.log(`could not stake KEEP tokens from a managed grant for ${operator}: ${err}`);
      });

      console.log(`successfully staked KEEP tokens from a managed grant (${(await managedGrantContract.grantId()).toString()}) for ${operator}`)
    }
  }
}

function OwnedTokensStakingStrategy(from) {
  return {
    stake: async (stakingContractAddress, amount, delegation, operator) => {
      const token = await KeepToken.deployed();
      const staked = await token.approveAndCall(
        stakingContractAddress,
        amount,
        delegation,
        { from }
      ).catch((err) => {
        console.log(`could not stake KEEP tokens for ${operator}: ${err}`);
      });

      if (staked) {
        console.log(`successfully staked KEEP tokens for account ${operator}`)
      }

    }
  }
}

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

    const stakingManager = async (numberOfAccounts, accountsOffset = 0, stakingStrategy) => {
      for(let i = accountsOffset; i < numberOfAccounts + accountsOffset; i++) {
        const delegation = getDelegationExtraData(accounts[i])
        const operator = accounts[i]
        const amount = formatAmount(200000, 18)
    
        await stakingStrategy.stake(tokenStaking.address, amount, delegation, operator)
    
        await stakingPortBackerContract.allowOperator(operator, { from: owner })
          .catch(err => console.log(`could not allowOperator for ${operator}`, err));
      }
    }

    await stakingManager(5, 5, new OwnedTokensStakingStrategy(owner))

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
    await stakingManager(5, 10, new GrantStakingStrategy(grantId, grantee))

    // Create managed grant
    await token.approveAndCall(managedGrantFactoryContract.address, formatAmount(12300000, 18), managedGrantExtraData, { from: owner })
    // Get the address of managed grant contract from an event.
    const managedGrant1Event = (await managedGrantFactoryContract.getPastEvents())[0]
    const managedGrant1Address = managedGrant1Event.args['grantAddress'];

    // delegate stake from a managed grant
    await stakingManager(5, 15, new ManagedGrantStakingStrategy(managedGrant1Address, grantee))
  } catch (err) {
    console.error('unexpected error:', err)
    process.exit(1)
  }

  process.exit();
};
