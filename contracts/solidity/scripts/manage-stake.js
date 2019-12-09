/*
  This script provides friendly interface to manage stake for development purposes.

  To stake KEEP tokens, use 'stake' command and provide three parameters:
  - operator address
  - amount of KEEP to stake
  - KEEP owner address (we assume it's also magpie/beneficiary address)

    $ truffle exec scripts/manage-stake.js stake 0x524f2E0176350d950fA630D9A5a59A0a190DAf48 10000 0xFa3DA235947AaB49D439f3BcB46effD1a7237E32

  To initiate unstake, use 'initiate-unstake' command and provide two parameters:
  - operator address
  - amount of KEEP to unstake

    $ truffle exec scripts/manage-stake.js initiate-unstake 0x524f2E0176350d950fA630D9A5a59A0a190DAf48 100

  To finish unstaking, use 'finish-unstake' command and provide operator address
  as a parameter. Please bear in mind you may need to wait for the expected 
  withdrawal delay time to be able to finish unstaking.

    $ truffle exec scripts/manage-stake.js finish-unstake 0x524f2E0176350d950fA630D9A5a59A0a190DAf48

  To print information about the stake, use 'print-stake' command and provide
  operator address as a parameter.

    $ truffle exec scripts/manage-stake.js print-stake 0x524f2E0176350d950fA630D9A5a59A0a190DAf48
*/
const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("TokenStaking.sol");
const KeepRandomBeaconOperator = artifacts.require("KeepRandomBeaconOperator.sol");

module.exports = async function() {
    let keepToken;
    let tokenStaking;
    let keepRandomBeaconOperator;

    keepToken = await KeepToken.deployed();
    tokenStaking = await TokenStaking.deployed();
    keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();

    const operation = process.argv[4];
    const operator = process.argv[5];

    switch(operation) {
        case "stake":         
            await stake();
            break;
        case "initiate-unstake":               
            await initiateUnstake(); 
            break;
        case "finish-unstake":            
            await finishUnstake();
            break;
        case "print-stake":
            console.log(`Getting stake info for operator ${operator}`); 
            await printStakeInfo();  
            break;
        default:
            console.log(`Unrecognized operation ${operation}`);
    }

    async function stake() {
        const amountToStake = process.argv[6]; 
        const owner = process.argv[7];
        const magpie = owner;

        console.log(`Staking ${amountToStake} tokens from ${owner} to operator ${operator} using beneficiary ${magpie}`); 
       
        const delegation = '0x' + Buffer.concat([
            Buffer.from(magpie.substr(2), 'hex'),
            Buffer.from(operator.substr(2), 'hex')]).toString('hex');
    
        try {
            staked = await keepToken.approveAndCall(
                tokenStaking.address, 
                amountToStake,
                delegation
            );  
    
            if (staked) {
                console.log(`successfully staked KEEP to ${operator}`)
            } else {
                console.log(`failed to stake KEEP to ${operator}`) 
            }
        } catch (err) {
            console.log(`could not stake KEEP to ${operator}: ${err}`);
        }
    }

    async function initiateUnstake() {
        const amountToUnstake = process.argv[6];

        console.log(`Initiating unstake of ${amountToUnstake} KEEP from operator ${operator}`);  

        try {          
          await tokenStaking.initiateUnstake(amountToUnstake, operator);
        } catch (err) {
            console.log(err);
        }
    }

    async function finishUnstake() {
        console.log(`Finishing unstake from operator ${operator}`); 
        
        try {
          await tokenStaking.finishUnstake(operator);          
        } catch (err) {
            console.log(err);
        }
    }

    async function printStakeInfo() {
        try {
            let owner = await getOwner();
            if (owner == 0) {
                console.log('No KEEP tokens staked for this operator');
                return;
            }

            console.log(`KEEP owner:     ${owner.toString()}`);
            console.log(`KEEP unstaked:  ${(await getOwnerBalance()).toString()}`);            
            console.log(`KEEP staked:    ${(await getStakeBalance()).toString()}`);
            console.log(`Minimum stake:  ${(await getMinimumStake()).toString()}`);
        } catch (err) {
            console.log(err);
        }  
    }

    async function getOwner() {
        return tokenStaking.ownerOf(operator);
    }

    async function getStakeBalance() {
      return tokenStaking.balanceOf(operator)
    }

    async function getMinimumStake() {
        return keepRandomBeaconOperator.minimumStake();
    }

    async function getOwnerBalance() {
        let owner = await getOwner()
        return keepToken.balanceOf(owner);
    }

    process.exit();
}
