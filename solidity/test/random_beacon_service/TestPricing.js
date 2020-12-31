const stakeDelegate = require("../helpers/stakeDelegate")
const blsData = require("../helpers/data.js")
const { initContracts } = require("../helpers/initContracts")
const assert = require("chai").assert
const { contract, accounts, web3 } = require("@openzeppelin/test-environment")
const { time } = require("@openzeppelin/test-helpers")
const CallbackContract = contract.fromArtifact("CallbackContract")

describe("TestKeepRandomBeaconService/Pricing", function () {
  let token
  let stakingContract
  let operatorContract
  let serviceContract
  let callbackContract
  let entryFee
  let groupSize
  let group
  const owner = accounts[0]
  const requestor = accounts[1]
  const operator1 = accounts[2]
  const operator2 = accounts[3]
  const operator3 = accounts[4]
  const beneficiary1 = accounts[5]
  const beneficiary2 = accounts[6]
  const beneficiary3 = accounts[7]

  beforeEach(async () => {
    const contracts = await initContracts(
      contract.fromArtifact("TokenStaking"),
      contract.fromArtifact("KeepRandomBeaconService"),
      contract.fromArtifact("KeepRandomBeaconServiceImplV1"),
      contract.fromArtifact("KeepRandomBeaconOperatorServicePricingStub")
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    operatorContract = contracts.operatorContract
    serviceContract = contracts.serviceContract
    callbackContract = await CallbackContract.new()

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(blsData.groupPubKey)

    groupSize = web3.utils.toBN(3)
    group = await operatorContract.getGroupPublicKey(0)
    await operatorContract.setGroupMembers(group, [
      operator1,
      operator2,
      operator3,
    ])
    const minimumStake = await stakingContract.minimumStake()

    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator1,
      beneficiary1,
      operator1,
      minimumStake
    )
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator2,
      beneficiary2,
      operator2,
      minimumStake
    )
    await stakeDelegate(
      stakingContract,
      token,
      owner,
      operator3,
      beneficiary3,
      operator3,
      minimumStake
    )

    entryFee = await serviceContract.entryFeeBreakdown()
  })

  it("should successfully refund callback surplus for a lower submission gas price", async () => {
    const gasPriceCeiling = web3.utils.toBN(web3.utils.toWei("20", "gwei"))
    await operatorContract.setGasPriceCeiling(gasPriceCeiling)

    const callbackGas = web3.utils.toBN(
      await callbackContract.__beaconCallback.estimateGas(
        blsData.groupSignature
      )
    )
    const entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)

    await serviceContract.methods["requestRelayEntry(address,uint256)"](
      callbackContract.address,
      callbackGas,
      {
        value: entryFeeEstimate,
        from: requestor,
      }
    )

    const submissionGasPrice = web3.utils.toBN(web3.utils.toWei("5", "gwei"))
    const gasPriceDiff = gasPriceCeiling.sub(submissionGasPrice)

    const requestorBalance = await web3.eth.getBalance(requestor)
    await operatorContract.relayEntry(blsData.groupSignature, {
      gasPrice: submissionGasPrice,
    })
    const updatedRequestorBalance = await web3.eth.getBalance(requestor)

    const refund = web3.utils
      .toBN(updatedRequestorBalance)
      .sub(web3.utils.toBN(requestorBalance))

    const baseCallbackGas = await serviceContract.baseCallbackGas()
    const expectedSurplus = callbackGas.add(baseCallbackGas).mul(gasPriceDiff)

    assert.isTrue(
      expectedSurplus.eq(refund),
      "Callback gas surplus should be refunded to the requestor."
    )
  })

  it("should send group reward to each operator.", async function () {
    const entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    const tx = await serviceContract.methods[
      "requestRelayEntry(address,uint256)"
    ](callbackContract.address, 0, { value: entryFeeEstimate, from: requestor })

    const currentRequestStartBlock = web3.utils.toBN(tx.receipt.blockNumber)
    const relayEntryTimeout = await operatorContract.relayEntryTimeout()
    const deadlineBlock = currentRequestStartBlock.add(relayEntryTimeout)
    const entryReceivedBlock = currentRequestStartBlock.addn(1)
    const remainingBlocks = deadlineBlock.sub(entryReceivedBlock)
    const submissionWindow = deadlineBlock.sub(entryReceivedBlock)
    const decimalPoints = web3.utils.toBN(1e16)
    const delayFactor = remainingBlocks
      .mul(decimalPoints)
      .div(submissionWindow)
      .pow(web3.utils.toBN(2))
    const memberBaseReward = entryFee.groupProfitFee.div(groupSize)
    const expectedGroupMemberReward = memberBaseReward
      .mul(delayFactor)
      .div(decimalPoints.pow(web3.utils.toBN(2)))

    await operatorContract.relayEntry(blsData.groupSignature)

    assert.isTrue(
      delayFactor.eq(web3.utils.toBN(1e16).pow(web3.utils.toBN(2))),
      "Delay factor expected to be 1 * 1e16 ^ 2."
    )

    const groupMemberRewards = await operatorContract.getGroupMemberRewards(
      group
    )
    assert.isTrue(
      web3.utils
        .toBN(groupMemberRewards)
        .eq(web3.utils.toBN(expectedGroupMemberReward)),
      "Unexpected group member reward."
    )
  })

  it("should send part of the group reward to request subsidy pool based on the submission block.", async function () {
    // Example rewards breakdown:
    // entryVerificationGasEstimate: 1240000
    // groupCreationGasEstimate: 2260000
    // dkgContributionMargin: 10%
    // groupMemberBaseReward: 1050000000000000
    // groupSize: 5
    // entry fee estimate: 49230000000000000 wei
    // signing fee: 37200000000000000 wei
    // DKG fee: 6780000000000000 wei
    // relayEntryTimeout: 10 blocks
    // currentRequestStartBlock: 38
    // relay entry submission block: 44
    // decimals: 1e16
    // groupProfitFee: 42450000000000000 - 37200000000000000 = 5250000000000000 wei
    // memberBaseReward: 5250000000000000 / 5 = 1050000000000000 wei
    // entryTimeout: 38 + 10 = 48
    // delayFactor: ((48 - 44) * 1e16 / (10 - 1)) ^ 2 = 19753086419753082469135802469136
    // groupMemberDelayPenalty: 1050000000000000 * 80246913580246917530864197530864 / (1e16 ^ 2) = 842592592592592
    // groupMemberReward: 1050000000000000 * 19753086419753082469135802469136) / (1e16 ^ 2) = 207407407407407 wei
    // submitterExtraReward: 842592592592592 * 5 * 5 / 100 = 210648148148148 wei
    // submitterReward: 37200000000000000 + 210648148148148 = 37410648148148148 wei
    // subsidy = 5250000000000000 - 207407407407407 * 5 - 210648148148148 = 4002314814814817 wei

    const entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    const tx = await serviceContract.methods[
      "requestRelayEntry(address,uint256)"
    ](callbackContract.address, 0, { value: entryFeeEstimate, from: requestor })

    const currentRequestStartBlock = web3.utils.toBN(tx.receipt.blockNumber)
    const relayEntryTimeout = await operatorContract.relayEntryTimeout()
    const deadlineBlock = currentRequestStartBlock
      .add(relayEntryTimeout)
      .addn(1)
    const submissionStartBlock = currentRequestStartBlock.addn(1)
    const decimalPoints = web3.utils.toBN(1e16)

    await time.advanceBlockTo(
      web3.utils.toBN(await web3.eth.getBlockNumber()).addn(1)
    )

    const entryReceivedBlock = web3.utils
      .toBN(await web3.eth.getBlockNumber())
      .add(web3.utils.toBN(1)) // web3.eth.getBlockNumber is 1 block behind solidity 'block.number'.
    const remainingBlocks = deadlineBlock.sub(entryReceivedBlock)
    const submissionWindow = deadlineBlock.sub(submissionStartBlock)
    const delayFactor = remainingBlocks
      .mul(decimalPoints)
      .div(submissionWindow)
      .pow(web3.utils.toBN(2))

    const memberBaseReward = entryFee.groupProfitFee.div(groupSize)
    const expectedGroupMemberReward = memberBaseReward
      .mul(delayFactor)
      .div(decimalPoints.pow(web3.utils.toBN(2)))
    const expectedDelayPenalty = memberBaseReward.sub(
      memberBaseReward
        .mul(delayFactor)
        .div(decimalPoints.pow(web3.utils.toBN(2)))
    )
    const expectedSubmitterExtraReward = expectedDelayPenalty
      .mul(groupSize)
      .muln(5)
      .div(web3.utils.toBN(100))
    const requestSubsidy = entryFee.groupProfitFee
      .sub(expectedGroupMemberReward.mul(groupSize))
      .sub(expectedSubmitterExtraReward)

    const serviceContractBalance = web3.utils.toBN(
      await web3.eth.getBalance(serviceContract.address)
    )

    await operatorContract.relayEntry(blsData.groupSignature)

    const groupMemberRewards = await operatorContract.getGroupMemberRewards(
      group
    )
    assert.isTrue(
      groupMemberRewards.eq(expectedGroupMemberReward),
      "Unexpected group member reward."
    )
    assert.isTrue(
      serviceContractBalance
        .add(requestSubsidy)
        .eq(
          web3.utils.toBN(await web3.eth.getBalance(serviceContract.address))
        ),
      "Service contract should receive request subsidy."
    )
  })
})
