import {sign} from './helpers/signature'
import mineBlocks from './helpers/mineBlocks'
import packTicket from './helpers/packTicket'
import generateTickets from './helpers/generateTickets'
import stakeDelegate from './helpers/stakeDelegate'
import {initContracts} from './helpers/initContracts'
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"
import {bls} from './helpers/data'


contract('KeepRandomBeaconOperator', function(accounts) {

  let resultPublicationTime, token, stakingContract, operatorContract, serviceContract, dkgPayment,
  owner = accounts[0], magpie = accounts[4], ticket,
  operator1 = accounts[0],
  operator2 = accounts[1],
  operator3 = accounts[2],
  selectedParticipants, signatures, signingMemberIndices = [],
  disqualified = '0x0000000000000000000000000000000000000000',
  inactive = '0x0000000000000000000000000000000000000000',
  groupPubKey = bls.groupPubKey,
  resultHash = web3.utils.soliditySha3(groupPubKey, disqualified, inactive)

  const groupSize = 20
  const groupThreshold = 15
  const minimumStake = web3.utils.toBN(200000)

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
    serviceContract = contracts.serviceContract

    operatorContract.setGroupSize(groupSize)
    operatorContract.setGroupThreshold(groupThreshold)
    operatorContract.setMinimumStake(minimumStake)

    await stakeDelegate(stakingContract, token, owner, operator1, magpie, minimumStake.mul(web3.utils.toBN(2000)))
    await stakeDelegate(stakingContract, token, owner, operator2, magpie, minimumStake.mul(web3.utils.toBN(2000)))
    await stakeDelegate(stakingContract, token, owner, operator3, magpie, minimumStake.mul(web3.utils.toBN(3000)))

    let tickets1 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator1, 2000)
    let tickets2 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator2, 2000)
    let tickets3 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator3, 3000)

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets1[i].valueHex, tickets1[i].virtualStakerIndex, operator1)
      await operatorContract.submitTicket(ticket, {from: operator1})
    }

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets2[i].valueHex, tickets2[i].virtualStakerIndex, operator2)
      await operatorContract.submitTicket(ticket, {from: operator2})
    }

    for(let i = 0; i < groupSize; i++) {
      ticket = packTicket(tickets3[i].valueHex, tickets3[i].virtualStakerIndex, operator3)
      await operatorContract.submitTicket(ticket, {from: operator3})
    }

    let ticketSubmissionStartBlock = (await operatorContract.getTicketSubmissionStartBlock()).toNumber()
    let timeoutChallenge = (await operatorContract.ticketSubmissionTimeout()).toNumber()
    let timeDKG = (await operatorContract.timeDKG()).toNumber()
    resultPublicationTime = ticketSubmissionStartBlock + timeoutChallenge + timeDKG

    selectedParticipants = await operatorContract.selectedParticipants()

    signingMemberIndices = []
    signatures = undefined

    for(let i = 0; i < selectedParticipants.length; i++) {
      let signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i+1)
      if (signatures == undefined) signatures = signature
      else signatures += signature.slice(2, signature.length)
    }

    let dkgGasEstimateCost = await operatorContract.dkgGasEstimate()
    let fluctuationMargin = await operatorContract.fluctuationMargin()
    let priceFeedEstimate = await serviceContract.priceFeedEstimate()
    let gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(fluctuationMargin).div(web3.utils.toBN(100)))
    dkgPayment = dkgGasEstimateCost.mul(gasPriceWithFluctuationMargin)

    // STEP 4 Make one group
    let currentBlock = await web3.eth.getBlockNumber()
    mineBlocks(resultPublicationTime - currentBlock)
    await operatorContract.submitDkgResult(1, groupPubKey, disqualified, inactive, signatures, signingMemberIndices, {from: selectedParticipants[0]})
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should be able to submit correct result as first member after DKG finished.", async function() {
    for (let i = 0; i < 10; i++) {
      // STEP 5 Fund dkgFeePool with an amount which will allow to trigger group selection after the next relay entry
      await serviceContract.fundDkgFeePool({value: dkgPayment})

      // STEP 6 Request a relay entry
      let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
      await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate})
      let tx = await operatorContract.relayEntry(bls.nextGroupSignature)

      // STEP 7 After the entry is submitted a group selection should start.
      // Simulate DKG fail
      mineBlocks(160)

      // STEP 8 After DKG timeout passes, fund dkgFeePool again with an amount which will allow to trigger group selection after the next relay entry. Make sure the pool is properly funded. Also, at this point dkgSubmitterReimbursementFee should have a value equal to a DKG fee
      await serviceContract.fundDkgFeePool({value: dkgPayment})

      // STEP 9 Request a relay entry
      await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate})

      // STEP 10 After the entry is submitted a group selection should be triggered.
      tx = await operatorContract.relayEntry(bls.nextGroupSignature)
      console.log('Relay entry ' + i + ' gas used:', tx.receipt.gasUsed)

      assert.equal(tx.logs[1].event, "GroupSelectionStarted", "Group selection should be triggered")
    }
  })
})
