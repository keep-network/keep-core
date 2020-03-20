import generateTickets from './generateTickets';
import mineBlocks from './mineBlocks';
import packTicket from './packTicket';
import {sign} from './signature';
import {bls} from './data';
import stakeDelegate from './stakeDelegate';


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
    let beneficiary1 = accounts[4];
    let beneficiary2 = accounts[5];
    let beneficiary3 = accounts[6];

    let operatorContract = contracts.operatorContract;
    let stakingContract = contracts.stakingContract;
    let token = contracts.token;
    let ticket;

    const operator1StakingWeight = 100;
    const operator2StakingWeight = 200;
    const operator3StakingWeight = 300;

    let owner = accounts[0];
    let authorizer = accounts[0];

    let minimumStake = await stakingContract.minimumStake()

    await stakeDelegate(stakingContract, token, owner, operator1, beneficiary1, authorizer, minimumStake.muln(operator1StakingWeight));
    await stakeDelegate(stakingContract, token, owner, operator2, beneficiary2, authorizer, minimumStake.muln(operator2StakingWeight));
    await stakeDelegate(stakingContract, token, owner, operator3, beneficiary3, authorizer, minimumStake.muln(operator3StakingWeight));

    await stakingContract.authorizeOperatorContract(operator1, operatorContract.address, {from: authorizer})
    await stakingContract.authorizeOperatorContract(operator2, operatorContract.address, {from: authorizer})
    await stakingContract.authorizeOperatorContract(operator3, operatorContract.address, {from: authorizer})

    let groupSize = await operatorContract.groupSize();

    let tickets1 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator1, operator1StakingWeight);
    let tickets2 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator2, operator2StakingWeight);
    let tickets3 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator3, operator3StakingWeight);

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets1[i].valueHex, tickets1[i].virtualStakerIndex, operator1);
      await operatorContract.submitTicket(ticket, {from: operator1});
    }

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets2[i].valueHex, tickets2[i].virtualStakerIndex, operator2);
      await operatorContract.submitTicket(ticket, {from: operator2});
    }

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets3[i].valueHex, tickets3[i].virtualStakerIndex, operator3);
      await operatorContract.submitTicket(ticket, {from: operator3});
    }

    let ticketSubmissionStartBlock = (await operatorContract.getTicketSubmissionStartBlock()).toNumber();
    let submissionTimeout = (await operatorContract.ticketSubmissionTimeout()).toNumber();

    mineBlocks(submissionTimeout);

    let selectedParticipants = await operatorContract.selectedParticipants();

    let timeDKG = (await operatorContract.timeDKG()).toNumber();
    let resultPublicationTime = ticketSubmissionStartBlock + submissionTimeout + timeDKG;

    mineBlocks(resultPublicationTime);

    let misbehaved = '0x';
    let resultHash = web3.utils.soliditySha3(bls.groupPubKey, misbehaved);

    let signingMemberIndices = [];
    let signatures = undefined;

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await sign(resultHash, selectedParticipants[i]);
      signingMemberIndices.push(i+1);
      if (signatures === undefined) signatures = signature
      else signatures += signature.slice(2, signature.length);
    }

    await operatorContract.submitDkgResult(
      1, bls.groupPubKey, misbehaved, signatures, signingMemberIndices,
      {from: selectedParticipants[0]}
    );
}
