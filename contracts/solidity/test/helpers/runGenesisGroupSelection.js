import generateTickets from './generateTickets';
import mineBlocks from './mineBlocks';
import {sign} from './signature';
import {bls} from './data';

export default async function runGenesisGroupSelection(
    operatorContract,
    operator1,
    operator2,
    operator3
) {
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
