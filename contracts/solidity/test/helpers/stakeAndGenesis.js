import generateTickets from './generateTickets';
import mineBlocks from './mineBlocks';
import {sign} from './signature';
import {bls} from './data';
import stakeDelegate from './stakeDelegate';

const operator1StakingWeight = 2000;
const operator2StakingWeight = 2000;
const operator3StakingWeight = 3000;

const minimumStake = web3.utils.toBN(200000);

// Function stakes first three accounts provided in the array of accounts and
// executes the entire genesis cycle registering group with bls.groupPubKey
// on the chain.
// 
// It expects three contracts to be passed:
// - contracts.operatorContract,
// - contracts.stakingContract,
// - contracts.token.
//
// This function should be usually used on the result of initContracts which
// initializes contracts up to the point when genesis should be performed.
export default async function stakeAndGenesis(accounts, contracts) {
    let operator1 = accounts[1];
    let operator2 = accounts[2];
    let operator3 = accounts[3];

    let operatorContract = contracts.operatorContract;
    let stakingContract = contracts.stakingContract;
    let token = contracts.token;

    let owner = accounts[0];

    await operatorContract.setMinimumStake(minimumStake);

    await stakeDelegate(stakingContract, token, owner, operator1, operator1, minimumStake.mul(web3.utils.toBN(operator1StakingWeight)));
    await stakeDelegate(stakingContract, token, owner, operator2, operator2, minimumStake.mul(web3.utils.toBN(operator2StakingWeight)));
    await stakeDelegate(stakingContract, token, owner, operator3, operator3, minimumStake.mul(web3.utils.toBN(operator3StakingWeight)));

    let groupSize = await operatorContract.groupSize();

    let tickets1 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator1, 2000);
    let tickets2 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator2, 2000);
    let tickets3 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator3, 3000);

    for(let i = 0; i < groupSize; i++) {
      await operatorContract.submitTicket(tickets1[i].value, operator1, tickets1[i].virtualStakerIndex, {from: operator1});
    }

    for(let i = 0; i < groupSize; i++) {
      await operatorContract.submitTicket(tickets2[i].value, operator2, tickets2[i].virtualStakerIndex, {from: operator2});
    }

    for(let i = 0; i < groupSize; i++) {
      await operatorContract.submitTicket(tickets3[i].value, operator3, tickets3[i].virtualStakerIndex, {from: operator3});
    }

    let ticketSubmissionStartBlock = (await operatorContract.getTicketSubmissionStartBlock()).toNumber();
    let submissionTimeout = (await operatorContract.ticketReactiveSubmissionTimeout()).toNumber();

    mineBlocks(submissionTimeout);

    let selectedParticipants = await operatorContract.selectedParticipants();

    let timeDKG = (await operatorContract.timeDKG()).toNumber();
    let resultPublicationTime = ticketSubmissionStartBlock + submissionTimeout + timeDKG;

    mineBlocks(resultPublicationTime);

    let disqualified = '0x0000000000000000000000000000000000000000';
    let inactive = '0x0000000000000000000000000000000000000000';
    let resultHash = web3.utils.soliditySha3(bls.groupPubKey, disqualified, inactive);

    let signingMemberIndices = [];
    let signatures = undefined;

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    await operatorContract.submitDkgResult(
      1, bls.groupPubKey, disqualified, inactive, signatures, signingMemberIndices,
      {from: selectedParticipants[0]}
    )
}
