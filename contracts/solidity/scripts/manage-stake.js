/*
  This script provides friendly interface to manage stake for development purposes.

  To stake KEEP tokens, use 'stake' command and provide three parameters:
  - operator address
  - amount of KEEP to stake
  - KEEP owner address (we assume it's also magpie/beneficiary/authorizer address)

    $ truffle exec scripts/manage-stake.js stake 0x524f2E0176350d950fA630D9A5a59A0a190DAf48 10000 0xFa3DA235947AaB49D439f3BcB46effD1a7237E32

  To undelegate, use 'undelegate' command and provide operator address as a parameter

    $ truffle exec scripts/manage-stake.js undelegate 0x524f2E0176350d950fA630D9A5a59A0a190DAf48

  To recover stake, use 'recover-stake' command and provide operator address
  as a parameter. Please bear in mind you may need to wait for the expected 
  undelegation period time to be able to recover stake.

    $ truffle exec scripts/manage-stake.js recover-stake 0x524f2E0176350d950fA630D9A5a59A0a190DAf48

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
        case "undelegate":               
            await undelegate(); 
            break;
        case "recover-stake":            
            await recoverStake();
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
        const authorizer = owner;

        console.log(`Staking ${amountToStake} tokens from ${owner} to operator ${operator} using beneficiary ${magpie} and authorizer ${authorizer}`);
       
        const delegation = '0x' + Buffer.concat([
            Buffer.from(magpie.substr(2), 'hex'),
            Buffer.from(operator.substr(2), 'hex'),
            Buffer.from(authorizer.substr(2), 'hex')
        ]).toString('hex');
    
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

        try {
            await tokenStaking.authorizeOperatorContract(operator, keepRandomBeaconOperator.address, {from: authorizer});
        } catch (err) {
            console.log(`could not authorize operator contract for ${operator}: ${err}`);
        }
    }

    async function undelegate() {
        console.log(`Undelegate stake from operator ${operator}`);

        try {          
          await tokenStaking.undelegate(operator);
        } catch (err) {
            console.log(err);
        }
    }

    async function recoverStake() {
        console.log(`Recover stake from operator ${operator}`); 
        
        try {
          await tokenStaking.recoverStake(operator);          
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

            console.log(`KEEP owner:             ${owner.toString()}`);
            console.log(`KEEP tokens available:  ${(await getOwnerBalance()).toString()}`);            
            console.log(`KEEP eligible stake:    ${(await getEligibleStake()).toString()}`);
            console.log(`KEEP active stake:      ${(await getActiveStake()).toString()}`);
            console.log(`Minimum stake:          ${(await getMinimumStake()).toString()}`);
        } catch (err) {
            console.log(err);
        }  
    }

    async function getOwner() {
        return tokenStaking.ownerOf(operator);
    }

    async function getEligibleStake() {
      return tokenStaking.eligibleStake(operator, keepRandomBeaconOperator.address);
    }

    async function getActiveStake() {
        return tokenStaking.activeStake(operator, keepRandomBeaconOperator.address);
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
