const blsData = require("../helpers/data.js")
const initContracts = require('../helpers/initContracts')
const assert = require('chai').assert
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { contract, accounts, web3 } = require("@openzeppelin/test-environment")
const { expectRevert, time } = require("@openzeppelin/test-helpers")
const stakeDelegate = require('../helpers/stakeDelegate')
const BLS = contract.fromArtifact('BLS');

describe('KeepRandomBeaconOperator/Slashing', function () {
  let token, stakingContract, serviceContract, operatorContract, minimumStake, largeStake, entryFeeEstimate, groupIndex,
    registry, bls,
    owner = accounts[0],
    operator1 = accounts[1],
    operator2 = accounts[2],
    operator3 = accounts[3],
    tattletale = accounts[4],
    authorizer = accounts[5],
    anotherOperatorContract = accounts[6],
    registryKeeper = accounts[7];

  before(async () => {

    let contracts = await initContracts(
      contract.fromArtifact('KeepToken'),
      contract.fromArtifact('TokenStakingStub'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorStub')
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    serviceContract = contracts.serviceContract
    operatorContract = contracts.operatorContract
    registry = contracts.registry
    bls = await BLS.new()

    groupIndex = 0
    await operatorContract.registerNewGroup(blsData.groupPubKey)
    await operatorContract.setGroupMembers(blsData.groupPubKey, [operator1, operator2, operator3])

    minimumStake = await stakingContract.minimumStake()
    largeStake = minimumStake.muln(2)
    await stakeDelegate(stakingContract, token, owner, operator1, owner, authorizer, largeStake)
    await stakeDelegate(stakingContract, token, owner, operator2, owner, authorizer, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator3, owner, authorizer, minimumStake)
    await stakingContract.authorizeOperatorContract(operator1, operatorContract.address, { from: authorizer })
    await stakingContract.authorizeOperatorContract(operator2, operatorContract.address, { from: authorizer })
    await stakingContract.authorizeOperatorContract(operator3, operatorContract.address, { from: authorizer })

    time.increase((await stakingContract.initializationPeriod()).addn(1));

    entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({ value: entryFeeEstimate, from: accounts[0] })

    await registry.setRegistryKeeper(registryKeeper, { from: accounts[0] })

    await registry.approveOperatorContract(anotherOperatorContract, { from: registryKeeper })
    await stakingContract.authorizeOperatorContract(operator1, anotherOperatorContract, { from: authorizer })
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should be able to report unauthorized signing", async () => {
    let tattletaleSignature = await bls.sign(tattletale, blsData.secretKey);

    await operatorContract.reportUnauthorizedSigning(
      groupIndex,
      tattletaleSignature,
      { from: tattletale }
    )

    assert.isTrue((await stakingContract.balanceOf(operator1)).eq(largeStake.sub(minimumStake)), "Unexpected operator 1 balance")
    assert.isTrue((await stakingContract.balanceOf(operator2)).isZero(), "Unexpected operator 2 balance")
    assert.isTrue((await stakingContract.balanceOf(operator3)).isZero(), "Unexpected operator 3 balance")

    // Expecting 5% of all the seized tokens
    let expectedTattletaleReward = minimumStake.muln(3).muln(5).divn(100)
    assert.isTrue((await token.balanceOf(tattletale)).eq(expectedTattletaleReward), "Unexpected tattletale balance")

    // Group should be terminated, expecting total number of groups to become 0
    await expectRevert(
      serviceContract.methods['requestRelayEntry()']({ value: entryFeeEstimate, from: accounts[0] }),
      "Total number of groups must be greater than zero."
    )
  })

  it("should ignore invalid report of unauthorized signing", async () => {
    await expectRevert(
      operatorContract.reportUnauthorizedSigning(
        groupIndex,
        blsData.nextGroupSignature, // Wrong signature
        { from: tattletale }
      ),
      "Group terminated or sig invalid"
    )
    // Transaction reverted no changes are applied.
  })

  it("should be able to report failure to produce entry after relay entry timeout", async () => {
    let operator1balance = await stakingContract.balanceOf(operator1)
    let operator2balance = await stakingContract.balanceOf(operator2)
    let operator3balance = await stakingContract.balanceOf(operator3)

    await expectRevert(
      operatorContract.reportRelayEntryTimeout({ from: tattletale }),
      "Entry did not time out."
    )

    await time.advanceBlockTo(web3.utils.toBN(20).addn(await web3.eth.getBlockNumber()));
    await operatorContract.reportRelayEntryTimeout({ from: tattletale })

    assert.isTrue((await stakingContract.balanceOf(operator1)).eq(operator1balance.sub(minimumStake)), "Unexpected operator 1 balance")
    assert.isTrue((await stakingContract.balanceOf(operator2)).eq(operator2balance.sub(minimumStake)), "Unexpected operator 2 balance")
    assert.isTrue((await stakingContract.balanceOf(operator3)).eq(operator3balance.sub(minimumStake)), "Unexpected operator 3 balance")

    // Expecting 5% of all the seized tokens with reward adjustment of (20 / 64) = 31%
    let expectedTattletaleReward = minimumStake.muln(3).muln(5).divn(100).muln(31).divn(100)
    assert.isTrue((await token.balanceOf(tattletale)).eq(expectedTattletaleReward), "Unexpected tattletale balance")
  })
})
