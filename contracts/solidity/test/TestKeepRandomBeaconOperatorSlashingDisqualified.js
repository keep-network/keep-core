import {initContracts} from './helpers/initContracts'
import {sign} from './helpers/signature'
import mineBlocks from './helpers/mineBlocks'
import stakeDelegate from './helpers/stakeDelegate'
import packTicket from './helpers/packTicket'
import generateTickets from './helpers/generateTickets'
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"

contract('KeepRandomBeaconOperator', function(accounts) {
  let token, stakingContract, operatorContract, minimumStake, largeStake,
    owner = accounts[0],
    operator1 = accounts[1],
    operator2 = accounts[2],
    operator3 = accounts[3],
    selectedParticipants, signatures, signingMemberIndices = [],
    disqualified = '0x000001', // disqualified operator3
    inactive = '0x000000',
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

    minimumStake = await operatorContract.minimumStake()

    await stakeDelegate(stakingContract, token, owner, operator1, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator2, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator3, owner, minimumStake)

    operatorContract.setGroupSize(3)
    operatorContract.setGroupThreshold(2)
    let tickets1 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator1, 1)
    let tickets2 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator2, 1)
    let tickets3 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator3, 1)

    await operatorContract.submitTicket(
      packTicket(tickets1[0].valueHex, tickets1[0].virtualStakerIndex, operator1),
      {from: operator1}
    )

    await operatorContract.submitTicket(
      packTicket(tickets2[0].valueHex, tickets2[0].virtualStakerIndex, operator2),
      {from: operator2}
    )

    await operatorContract.submitTicket(
      packTicket(tickets3[0].valueHex, tickets1[0].virtualStakerIndex, operator3),
      {from: operator3}
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

  it("should be able to seize disqualified member stake and receive tattletale reward.", async function() {
    await operatorContract.submitDkgResult(1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, {from: operator1})
    assert.isTrue((await stakingContract.balanceOf(operator1)).eq(minimumStake),"Unexpected operator 1 balance")
    assert.isTrue((await stakingContract.balanceOf(operator2)).eq(minimumStake), "Unexpected operator 2 balance")

    // Expecting seized minimumStake
    assert.isTrue((await stakingContract.balanceOf(operator3)).isZero(), "Unexpected operator 3 balance")

    // Expecting 5% of the seized tokens
    let expectedTattletaleReward = minimumStake.muln(5).divn(100)
    assert.isTrue((await token.balanceOf(operator1)).eq(expectedTattletaleReward), "Unexpected tattletale balance")
  })
})
