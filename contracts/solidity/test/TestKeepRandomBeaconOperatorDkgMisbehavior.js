import {initContracts} from './helpers/initContracts'
import {sign} from './helpers/signature'
import mineBlocks from './helpers/mineBlocks'
import stakeDelegate from './helpers/stakeDelegate'
import packTicket from './helpers/packTicket'
import generateTickets from './helpers/generateTickets'
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"

contract('KeepRandomBeaconOperator', function(accounts) {
  let token, stakingContract, operatorContract, minimumStake,
    owner = accounts[0],
    operator1 = accounts[1],
    operator2 = accounts[2],
    operator3 = accounts[3],
    operator4 = accounts[4],
    operator5 = accounts[5],
    selectedParticipants, signatures, signingMemberIndices = [],
    disqualified = '0x03', // disqualified operator3
    inactive = '0x05', // inactive operator5
    groupPubKey = '0x1000000000000000000000000000000000000000000000000000000000000000',
    resultHash = web3.utils.soliditySha3(groupPubKey, disqualified, inactive)

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    operatorContract = contracts.operatorContract
    operatorContract.setGroupSize(5)
    operatorContract.setGroupThreshold(3)
    minimumStake = await operatorContract.minimumStake()

    await stakeDelegate(stakingContract, token, owner, operator1, owner, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator2, owner, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator3, owner, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator4, owner, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator5, owner, owner, minimumStake)

    let tickets1 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator1, 1)
    let tickets2 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator2, 1)
    let tickets3 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator3, 1)
    let tickets4 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator4, 1)
    let tickets5 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator5, 1)

    await operatorContract.submitTicket(
      packTicket(tickets1[0].valueHex, tickets1[0].virtualStakerIndex, operator1),
      {from: operator1}
    )

    await operatorContract.submitTicket(
      packTicket(tickets2[0].valueHex, tickets2[0].virtualStakerIndex, operator2),
      {from: operator2}
    )

    await operatorContract.submitTicket(
      packTicket(tickets3[0].valueHex, tickets3[0].virtualStakerIndex, operator3),
      {from: operator3}
    )

    await operatorContract.submitTicket(
      packTicket(tickets4[0].valueHex, tickets4[0].virtualStakerIndex, operator4),
      {from: operator4}
    )

    await operatorContract.submitTicket(
      packTicket(tickets5[0].valueHex, tickets5[0].virtualStakerIndex, operator5),
      {from: operator5}
    )

    let ticketSubmissionStartBlock = (await operatorContract.getTicketSubmissionStartBlock()).toNumber()
    let timeoutChallenge = (await operatorContract.ticketSubmissionTimeout()).toNumber()
    let timeDKG = (await operatorContract.timeDKG()).toNumber()
    let resultPublicationTime = ticketSubmissionStartBlock + timeoutChallenge + timeDKG

    let currentBlock = await web3.eth.getBlockNumber()
    mineBlocks(resultPublicationTime - currentBlock)

    selectedParticipants = await operatorContract.selectedParticipants()

    signingMemberIndices = []
    signatures = undefined

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i+1)
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length)
    }
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should be able to save group members based on disqualified and inactive data", async () => {
    await operatorContract.submitDkgResult(1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, {from: operator1})
    let registeredMembers = await operatorContract.getGroupMembers(groupPubKey)
    assert.isTrue(registeredMembers.indexOf(operator1) == 0, "Member should be registered")
    assert.isTrue(registeredMembers.indexOf(operator2) == 1, "Member should be registered")
    assert.isTrue(registeredMembers.indexOf(operator3) == -1, "Member should not be registered")
    assert.isTrue(registeredMembers.indexOf(operator4) == 3, "Member should be registered")
    assert.isTrue(registeredMembers.indexOf(operator5) == -1, "Member should not be registered")
  })
})
