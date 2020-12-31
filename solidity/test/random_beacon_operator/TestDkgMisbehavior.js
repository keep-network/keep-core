const blsData = require("../helpers/data")
const { contract, web3, accounts } = require("@openzeppelin/test-environment")
const { time } = require("@openzeppelin/test-helpers")
const assert = require("chai").assert
const { initContracts } = require("../helpers/initContracts")
const sign = require("../helpers/signature")
const stakeDelegate = require("../helpers/stakeDelegate")
const packTicket = require("../helpers/packTicket")
const generateTickets = require("../helpers/generateTickets")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")

describe("KeepRandomBeaconOperator/DkgMisbehavior", function () {
  let token
  let stakingContract
  let operatorContract
  const owner = accounts[0]
  const operator1 = accounts[1]
  const operator2 = accounts[2]
  const operator3 = accounts[3]
  const operator4 = accounts[4]
  const operator5 = accounts[5]
  const authorizer = owner
  let selectedParticipants
  let signatures
  let signingMemberIndices = []
  const misbehaved = "0x0305" // disqualified operator with selected member index 3 and inactive with 5
  const groupPubKey = blsData.groupPubKey
  const resultHash = web3.utils.soliditySha3(groupPubKey, misbehaved)

  before(async () => {
    const contracts = await initContracts(
      contract.fromArtifact("TokenStaking"),
      contract.fromArtifact("KeepRandomBeaconService"),
      contract.fromArtifact("KeepRandomBeaconServiceImplV1"),
      contract.fromArtifact("KeepRandomBeaconOperatorMisbehaviorStub")
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    operatorContract = contracts.operatorContract

    const minimumStake = await stakingContract.minimumStake()
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator1,
      owner,
      authorizer,
      minimumStake,
      { from: owner }
    )
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator2,
      owner,
      authorizer,
      minimumStake,
      { from: owner }
    )
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator3,
      owner,
      authorizer,
      minimumStake,
      { from: owner }
    )
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator4,
      owner,
      authorizer,
      minimumStake,
      { from: owner }
    )
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator5,
      owner,
      authorizer,
      minimumStake,
      { from: owner }
    )

    await stakingContract.authorizeOperatorContract(
      operator1,
      operatorContract.address,
      { from: authorizer }
    )
    await stakingContract.authorizeOperatorContract(
      operator2,
      operatorContract.address,
      { from: authorizer }
    )
    await stakingContract.authorizeOperatorContract(
      operator3,
      operatorContract.address,
      { from: authorizer }
    )
    await stakingContract.authorizeOperatorContract(
      operator4,
      operatorContract.address,
      { from: authorizer }
    )
    await stakingContract.authorizeOperatorContract(
      operator5,
      operatorContract.address,
      { from: authorizer }
    )

    time.increase((await stakingContract.initializationPeriod()).addn(1))

    const groupSelectionRelayEntry = await operatorContract.getGroupSelectionRelayEntry()
    const tickets1 = generateTickets(groupSelectionRelayEntry, operator1, 1)
    const tickets2 = generateTickets(groupSelectionRelayEntry, operator2, 1)
    const tickets3 = generateTickets(groupSelectionRelayEntry, operator3, 1)
    const tickets4 = generateTickets(groupSelectionRelayEntry, operator4, 1)
    const tickets5 = generateTickets(groupSelectionRelayEntry, operator5, 1)

    await operatorContract.submitTicket(
      packTicket(
        tickets1[0].valueHex,
        tickets1[0].virtualStakerIndex,
        operator1
      ),
      { from: operator1 }
    )

    await operatorContract.submitTicket(
      packTicket(
        tickets2[0].valueHex,
        tickets2[0].virtualStakerIndex,
        operator2
      ),
      { from: operator2 }
    )

    await operatorContract.submitTicket(
      packTicket(
        tickets3[0].valueHex,
        tickets3[0].virtualStakerIndex,
        operator3
      ),
      { from: operator3 }
    )

    await operatorContract.submitTicket(
      packTicket(
        tickets4[0].valueHex,
        tickets4[0].virtualStakerIndex,
        operator4
      ),
      { from: operator4 }
    )

    await operatorContract.submitTicket(
      packTicket(
        tickets5[0].valueHex,
        tickets5[0].virtualStakerIndex,
        operator5
      ),
      { from: operator5 }
    )

    const ticketSubmissionStartBlock = await operatorContract.getTicketSubmissionStartBlock()
    const timeoutChallenge = await operatorContract.ticketSubmissionTimeout()
    const timeDKG = await operatorContract.timeDKG()
    const resultPublicationTime = ticketSubmissionStartBlock
      .add(timeoutChallenge)
      .add(timeDKG)

    await time.advanceBlockTo(resultPublicationTime)

    selectedParticipants = await operatorContract.selectedParticipants()

    signingMemberIndices = []
    signatures = undefined

    for (let i = 0; i < selectedParticipants.length; i++) {
      const signature = await sign(resultHash, selectedParticipants[i])
      signingMemberIndices.push(i + 1)
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

  it("should be able to save group members based on misbehaved data", async () => {
    await operatorContract.submitDkgResult(
      1,
      groupPubKey,
      misbehaved,
      signatures,
      signingMemberIndices,
      { from: selectedParticipants[0] }
    )
    const registeredMembers = await operatorContract.getGroupMembers(
      groupPubKey
    )
    assert.isTrue(
      registeredMembers.indexOf(selectedParticipants[0]) != -1,
      "Member should be registered"
    )
    assert.isTrue(
      registeredMembers.indexOf(selectedParticipants[1]) != -1,
      "Member should be registered"
    )
    assert.isTrue(
      registeredMembers.indexOf(selectedParticipants[2]) == -1,
      "Member should not be registered"
    )
    assert.isTrue(
      registeredMembers.indexOf(selectedParticipants[3]) != -1,
      "Member should be registered"
    )
    assert.isTrue(
      registeredMembers.indexOf(selectedParticipants[4]) == -1,
      "Member should not be registered"
    )
  })
})
