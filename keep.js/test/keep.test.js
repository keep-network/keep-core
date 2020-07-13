import chai from "chai"
const { expect } = chai
import sinon from "sinon"
import KEEP, { contracts } from "../src/keep.js"
import ContractFactory, { ContractWrapper } from "../src/contract-wrapper.js"
import { TokenStakingConstants } from "../src/constants.js"

const operatorAddress = "0x0"
const beneficiaryAddress = "0x1"
const authorizerAddress = "0x2"
const ownerAddress = "0x3"
const TBTCSortitionPollMockAddress = "0x4"

describe("Keep initialization", () => {
  let config
  let web3Mock
  let networkId

  beforeEach(() => {
    web3Mock = { eth: () => {} }
    networkId = 1

    config = { web3: web3Mock, networkId }
  })

  afterEach(function () {
    sinon.restore()
  })

  it("should initialize keep.js", async () => {
    sinon
      .stub(ContractFactory, "createContractInstance")
      .callsFake(() => Promise.resolve(new ContractWrapper({}, 100)))
    sinon.stub(TokenStakingConstants, "initialize").callsFake(() =>
      Promise.resolve({
        minimumStake: "1",
        initializationPeriod: "2",
        undelegationPeriod: "3",
      })
    )

    const keep = await KEEP.initialize(config)

    for (const [, propertyName] of contracts) {
      expect(keep[propertyName] instanceof ContractWrapper).to.be.true
    }

    expect(keep.config).equals(config)
    expect(keep.tokenStakingConstants.minimumStake).to.eq("1")
    expect(keep.tokenStakingConstants.initializationPeriod).to.eq("2")
    expect(keep.tokenStakingConstants.undelegationPeriod).to.eq("3")
    sinon.reset
  })
})

describe("KEEP.js functions", () => {
  let keep
  const sandbox = sinon.createSandbox()

  before(async () => {
    sandbox
      .stub(ContractFactory, "createContractInstance")
      .callsFake(() => Promise.resolve(new ContractWrapper()))

    sandbox.stub(TokenStakingConstants, "initialize").callsFake(() =>
      Promise.resolve({
        minimumStake: "1",
        initializationPeriod: "2",
        undelegationPeriod: "3",
      })
    )

    keep = await KEEP.initialize({})
  })

  afterEach(() => {
    sandbox.restore()
  })

  it("should return beneficiary of the provided operator", async () => {
    const stub = sandbox
      .stub(keep.tokenStakingContract, "makeCall")
      .returns(beneficiaryAddress)

    const beneficiary = await keep.beneficiaryOf(operatorAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).to.equal("beneficiaryOf")
    expect(args[1]).equal(operatorAddress)
    expect(beneficiary).equal(beneficiaryAddress)
  })

  it("should return authorizer of the provided operator", async () => {
    const stub = sandbox
      .stub(keep.tokenStakingContract, "makeCall")
      .returns(authorizerAddress)

    const authorizer = await keep.authorizerOf(operatorAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).to.equal("authorizerOf")
    expect(args[1]).equal(operatorAddress)
    expect(authorizer).equal(authorizerAddress)
  })

  it("should return owner of the provided operator", async () => {
    const stub = sandbox
      .stub(keep.tokenStakingContract, "makeCall")
      .returns(ownerAddress)

    const owner = await keep.ownerOf(operatorAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).to.equal("ownerOf")
    expect(args[1]).to.equal(operatorAddress)
    expect(owner).to.deep.equal(ownerAddress)
  })

  it("should return operators of the provided owner", async () => {
    const stub = sandbox
      .stub(keep.tokenStakingContract, "makeCall")
      .returns([operatorAddress])

    const operators = await keep.operatorsOf(ownerAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).to.equal("operatorsOf")
    expect(args[1]).eq(ownerAddress)
    expect(operators).to.deep.equal([operatorAddress])
  })

  it("should return operators of beneficiary", async () => {
    const stub = sandbox
      .stub(keep.tokenStakingContract, "getPastEvents")
      .returns([{ returnValues: { operator: operatorAddress } }])

    const operators = await keep.operatorsOfBeneficiary(beneficiaryAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).to.equal("Staked")
    expect(args[1]).to.deep.eq({ beneficiary: beneficiaryAddress })
    expect(operators).to.deep.equal([operatorAddress])
  })

  it("should return operators of authorizer", async () => {
    const stub = sandbox
      .stub(keep.tokenStakingContract, "getPastEvents")
      .returns([{ returnValues: { operator: operatorAddress } }])

    const operators = await keep.operatorsOfAuthorizer(authorizerAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).to.equal("Staked")
    expect(args[1]).to.deep.eq({ authorizer: authorizerAddress })
    expect(operators).to.deep.equal([operatorAddress])
  })

  it("should return delegation info", async () => {
    const mockReturnData = {
      amount: "100",
      createdAt: "1234",
      undelegatedAt: "5678",
    }
    const stub = sandbox
      .stub(keep.tokenStakingContract, "makeCall")
      .returns(mockReturnData)

    const delegationInfo = await keep.getDelegationInfo(operatorAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).equal("getDelegationInfo")
    expect(args[1]).eq(operatorAddress)
    expect(delegationInfo).to.deep.equal(mockReturnData)
  })

  it("should return gelegations for the provided operators", async () => {
    const operatorAddress1 = "0x001"
    const operatorAddress2 = "0x002"

    const delegationInfo1 = {
      amount: "100",
      createdAt: "1234",
      undelegatedAt: "5678",
    }
    const delegationInfo2 = {
      amount: "200",
      createdAt: "300",
      undelegatedAt: "400",
    }
    const expectedRetrunData = [
      {
        ...delegationInfo1,
        beneficiary: beneficiaryAddress,
        authorizer: authorizerAddress,
        operator: operatorAddress1,
      },
      {
        ...delegationInfo2,
        beneficiary: beneficiaryAddress,
        authorizer: authorizerAddress,
        operator: operatorAddress2,
      },
    ]
    const getDelegationInfoStub = sandbox
      .stub(keep, "getDelegationInfo")
      .onFirstCall()
      .returns(delegationInfo1)
      .onSecondCall()
      .returns(delegationInfo2)

    const beneficiaryOfStub = sandbox
      .stub(keep, "beneficiaryOf")
      .returns(beneficiaryAddress)

    const authorizerOfStub = sandbox
      .stub(keep, "authorizerOf")
      .returns(authorizerAddress)

    const delegations = await keep.getDelegations([
      operatorAddress1,
      operatorAddress2,
    ])

    expect(getDelegationInfoStub.calledTwice).to.be.true
    expect(getDelegationInfoStub.getCall(0).args[0]).equal(operatorAddress1)
    expect(getDelegationInfoStub.getCall(1).args[0]).equal(operatorAddress2)

    expect(beneficiaryOfStub.calledTwice).to.be.true
    expect(beneficiaryOfStub.getCall(0).args[0]).equal(operatorAddress1)
    expect(beneficiaryOfStub.getCall(1).args[0]).equal(operatorAddress2)

    expect(authorizerOfStub.calledTwice).to.be.true
    expect(authorizerOfStub.getCall(0).args[0]).equal(operatorAddress1)
    expect(authorizerOfStub.getCall(1).args[0]).equal(operatorAddress2)

    expect(delegations).to.deep.equal(expectedRetrunData)
  })

  it("should authorize KEEP Random Beacon Operator Contract", () => {
    sandbox
      .stub(keep, "keepRandomBeaconOperatorContract")
      .value({ address: "0x123456789" })

    const stub = sandbox
      .stub(keep.tokenStakingContract, "sendTransaction")
      .returns({})

    keep.authorizeKeepRandomBeaconOperatorContract(operatorAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).equal("authorizeOperatorContract")
    expect(args[1]).eq(operatorAddress)
    expect(args[2]).eq("0x123456789")
  })

  it("should check if operator contract has access to the staked token balance of the provided operator", async () => {
    sandbox
      .stub(keep, "keepRandomBeaconOperatorContract")
      .value({ address: "0x123456789" })

    const stub = sandbox
      .stub(keep.tokenStakingContract, "makeCall")
      .returns(true)

    await keep.isAuthorizedForKeepRandomBeacon(operatorAddress)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).equal("isAuthorizedForOperator")
    expect(args[1]).eq(operatorAddress)
    expect(args[2]).eq("0x123456789")
  })

  it("should return withdrawn rewards for the provided beneficiary", async () => {
    const mockData = [
      {
        returnValues: {
          beneficiary: beneficiaryAddress,
          operator: operatorAddress,
          amount: "100",
          groupIndex: "1",
        },
      },
    ]
    const stub = sandbox
      .stub(keep.keepRandomBeaconOperatorContract, "getPastEvents")
      .returns(mockData)

    const withdrawnRewards = await keep.getWithdrawnRewardsForBeneficiary(
      beneficiaryAddress
    )

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).equal("GroupMemberRewardsWithdrawn")
    expect(args[1]).to.deep.eq({ beneficiary: beneficiaryAddress })
    expect(withdrawnRewards).to.deep.eq(mockData)
  })

  it("should call withdrawGroupMemberRewards from KEEP Random Beacon Operator contract", () => {
    const groupIndex = 1
    const stub = sandbox.stub(
      keep.keepRandomBeaconOperatorContract,
      "sendTransaction"
    )

    keep.withdrawGroupMemberRewards(operatorAddress, groupIndex)

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).equal("withdrawGroupMemberRewards")
    expect(args[1]).to.deep.eq(operatorAddress)
    expect(args[2]).to.deep.eq(groupIndex)
  })

  it("should return all created groups", async () => {
    const stub = sandbox.stub(
      keep.keepRandomBeaconOperatorContract,
      "getPastEvents"
    )

    await keep.getAllCreatedGroups()

    expect(stub.calledOnce).to.be.true
    const args = stub.getCall(0).args
    expect(args[0]).equal("DkgResultSubmittedEvent")
  })

  it("should find all available beacon rewards for beneficiary", async () => {
    // given
    const mockGroups = [
      { returnValues: { groupPubKey: "1" } },
      { returnValues: { groupPubKey: "2" } },
    ]
    const operators = ["0x1234", "0x5678"]
    const getAllCreatedGroupsStub = sandbox
      .stub(keep, "getAllCreatedGroups")
      .returns(mockGroups)

    const operatorsOfBeneficiaryStub = sandbox
      .stub(keep, "operatorsOfBeneficiary")
      .returns(operators)

    const awaitingRewardsStub = sandbox
      .stub(keep.keepRandomBeaconOperatorStatisticsContract, "makeCall")
      .withArgs("awaitingRewards")
      .onFirstCall()
      .returns("1")
      .onSecondCall()
      .returns("0")
      .onThirdCall()
      .returns("2")
      .onCall(3)
      .returns("3")

    const keepRandomBeaconOperatorContractStub = sandbox
      .stub(keep.keepRandomBeaconOperatorContract, "makeCall")
      .withArgs("isStaleGroup")
      .onFirstCall()
      .returns(true)

    keepRandomBeaconOperatorContractStub.withArgs("isStaleGroup").returns(false)

    keepRandomBeaconOperatorContractStub
      .withArgs("isGroupTerminated")
      .returns(true)

    // when
    const rewards = await keep.findKeepRandomBeaconRewardsForBeneficiary(
      beneficiaryAddress
    )

    // then
    expect(getAllCreatedGroupsStub.calledOnce).to.be.true
    expect(operatorsOfBeneficiaryStub.calledOnce).to.be.true
    expect(operatorsOfBeneficiaryStub.getCall(0).args[0]).eq(beneficiaryAddress)
    expect(awaitingRewardsStub.callCount).to.eq(4)
    expect(keepRandomBeaconOperatorContractStub.callCount).to.eq(2)

    expect(
      keepRandomBeaconOperatorContractStub.calledWithExactly(
        "isStaleGroup",
        mockGroups[0].returnValues.groupPubKey
      )
    )
    expect(
      keepRandomBeaconOperatorContractStub.neverCalledWith(
        "isGroupTerminated",
        mockGroups[0].returnValues.groupPubKey
      )
    ).to.be.true

    expect(
      keepRandomBeaconOperatorContractStub.calledWithExactly(
        "isStaleGroup",
        mockGroups[1].returnValues.groupPubKey
      )
    )
    expect(
      keepRandomBeaconOperatorContractStub.calledWithExactly(
        "isGroupTerminated",
        mockGroups[1].returnValues.groupPubKey
      )
    )

    expect(rewards.length).to.eq(3)
    const firstReward = rewards[0]
    const secondReward = rewards[1]
    const thirdReward = rewards[2]
    expect(firstReward).to.deep.eq({
      groupIndex: "0",
      groupPublicKey: "1",
      isTerminated: false,
      isStale: true,
      operatorAddress: "0x1234",
      reward: "1",
    })
    expect(secondReward).to.deep.eq({
      groupIndex: "1",
      groupPublicKey: "2",
      isTerminated: true,
      isStale: false,
      operatorAddress: "0x1234",
      reward: "2",
    })
    expect(thirdReward).to.deep.eq({
      groupIndex: "1",
      groupPublicKey: "2",
      isTerminated: true,
      isStale: false,
      operatorAddress: "0x5678",
      reward: "3",
    })
  })

  it("should return empty data when no slashed tokens events", async () => {
    const tokenStakingStub = sandbox.stub(
      keep.tokenStakingContract,
      "getPastEvents"
    )

    tokenStakingStub
      .withArgs("TokensSlashed")
      .returns([])
      .withArgs("TokensSeized")
      .returns([])

    const slashedTokens = await keep.getSlashedTokens(operatorAddress)

    expect(
      tokenStakingStub.calledWithExactly("TokensSlashed", {
        operator: operatorAddress,
      })
    ).to.be.true
    expect(
      tokenStakingStub.calledWithExactly("TokensSeized", {
        operator: operatorAddress,
      })
    ).to.be.true

    expect(slashedTokens.length).eq(0)
  })

  it("should return slashed tokens correctly", async () => {
    const tokenStakingStub = sandbox.stub(
      keep.tokenStakingContract,
      "getPastEvents"
    )
    const keepRandomBeaconOperatorContractStub = sandbox.stub(
      keep.keepRandomBeaconOperatorContract,
      "getPastEvents"
    )

    tokenStakingStub
      .withArgs("TokensSlashed")
      .returns([
        {
          transactionHash: "1",
          returnValues: { operator: operatorAddress, amount: "1" },
        },
        {
          transactionHash: "1",
          returnValues: { operator: operatorAddress, amount: "2" },
        },
        {
          transactionHash: "2",
          returnValues: { operator: operatorAddress, amount: "1" },
        },
      ])
      .withArgs("TokensSeized")
      .returns([
        {
          transactionHash: "3",
          returnValues: { operator: operatorAddress, amount: "0" },
        },
      ])

    keepRandomBeaconOperatorContractStub
      .withArgs("UnauthorizedSigningReported")
      .returns([
        { transactionHash: "1", returnValues: { groupIndex: "1" } },
        { transactionHash: "2", returnValues: { groupIndex: "1" } },
      ])
      .withArgs("RelayEntryTimeoutReported")
      .returns([{ transactionHash: "3", returnValues: { groupIndex: "1" } }])

    const slashedTokens = await keep.getSlashedTokens(operatorAddress)

    expect(
      tokenStakingStub.calledWithExactly("TokensSlashed", {
        operator: operatorAddress,
      })
    ).to.be.true
    expect(
      tokenStakingStub.calledWithExactly("TokensSeized", {
        operator: operatorAddress,
      })
    ).to.be.true

    expect(
      keepRandomBeaconOperatorContractStub.calledWithExactly(
        "UnauthorizedSigningReported"
      )
    ).to.be.true
    expect(
      keepRandomBeaconOperatorContractStub.calledWithExactly(
        "RelayEntryTimeoutReported"
      )
    ).to.be.true

    expect(slashedTokens.length).eq(2)
    expect(slashedTokens[0].groupIndex).eq("1")
    expect(slashedTokens[0].amount.toString()).eq("3")

    expect(slashedTokens[1].groupIndex).eq("1")
    expect(slashedTokens[1].amount.toString()).eq("1")
  })

  it("should return withdrawn tbtc rewards to the provided beneficiary", async () => {
    const depositAddress1 = "0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca"
    const depositAddress2 = "0x2993ac0a73f1270973DF507F0b94622b45aBF47C"
    const mockTransferEvents = [
      {
        returnValues: {
          from: depositAddress1,
          to: beneficiaryAddress,
          value: "100",
        },
      },
      {
        returnValues: {
          from: depositAddress2,
          to: beneficiaryAddress,
          value: "200",
        },
      },
      // Transfer not from Deposit contract
      {
        returnValues: {
          from: "0x8738b323dF0eb841467996920Eb1eF0599C21395",
          to: beneficiaryAddress,
          value: "300",
        },
      },
    ]
    const mocktDepositCreatedEvents = [
      { returnValues: { _depositContractAddress: depositAddress1 } },
      { returnValues: { _depositContractAddress: depositAddress2 } },
    ]

    const tbtcTokenContractStub = sandbox
      .stub(keep.tbtcTokenContract, "getPastEvents")
      .returns(mockTransferEvents)

    const tbtcSystemContractStub = sandbox
      .stub(keep.tbtcSystemContract, "getPastEvents")
      .returns(mocktDepositCreatedEvents)

    const tbtcRewards = await keep.getWithdrawnTBTCRewards(beneficiaryAddress)

    expect(
      tbtcTokenContractStub.calledWithExactly("Transfer", {
        to: beneficiaryAddress,
      })
    ).to.be.true

    expect(
      tbtcSystemContractStub.calledWithExactly("Created", {
        _depositContractAddress: mockTransferEvents.map(
          (_) => _.returnValues.from
        ),
      })
    ).to.be.true
    expect(tbtcRewards.length).eq(2)
  })

  it("should return a operators for deposit", async () => {
    const depositAddressMock = "0x1"
    const keepAddressMock = "0x2"
    const makeCallStub = sandbox
      .stub()
      .withArgs("getKeepAddress")
      .returns(keepAddressMock)
    sandbox
      .stub(ContractFactory, "new")
      .callsFake(() => ({ makeCall: makeCallStub }))

    const ethContractStub = sandbox.stub()
    sandbox
      .stub(keep, "config")
      .value({ web3: { eth: { Contract: ethContractStub } } })

    await keep.getOperatorsForDeposit(depositAddressMock)

    expect(ethContractStub.callCount).eq(2)
    expect(ethContractStub.getCall(0).args[1]).eq(depositAddressMock)
    expect(ethContractStub.getCall(1).args[1]).eq(keepAddressMock)
    expect(makeCallStub.getCall(0).calledWithExactly("getKeepAddress")).to.be
      .true
    expect(makeCallStub.getCall(1).calledWithExactly("getMembers")).to.be.true
  })

  it("should create grants", () => {
    const tokenGrantContractMock = { address: "0x0011" }
    const extraDataMock = "extraData"
    const encodeParametersStub = sinon.stub().returns(extraDataMock)
    sandbox.stub(keep, "config").value({
      web3: { eth: { abi: { encodeParameters: encodeParametersStub } } },
    })
    sandbox.stub(keep, "tokenGrantContract").value(tokenGrantContractMock)

    const keepTokenContractStub = sandbox.stub(
      keep.keepTokenContract,
      "sendTransaction"
    )

    const mockData = {
      grantManager: "0x123",
      grantee: "0x456",
      duration: 1,
      start: 1,
      cliffDuration: 1,
      revocable: true,
      stakingPolicyAddress: "0x789",
      amount: "1000000000",
    }

    const {
      grantManager,
      grantee,
      duration,
      start,
      cliffDuration,
      revocable,
      stakingPolicyAddress,
    } = mockData

    keep.createGrant(mockData)

    expect(
      encodeParametersStub.calledWithExactly(
        [
          "address",
          "address",
          "uint256",
          "uint256",
          "uint256",
          "bool",
          "address",
        ],
        [
          grantManager,
          grantee,
          duration,
          start,
          cliffDuration,
          revocable,
          stakingPolicyAddress,
        ]
      )
    ).to.be.true

    expect(
      keepTokenContractStub.calledWithExactly(
        "approveAndCall",
        tokenGrantContractMock.address,
        mockData.amount,
        extraDataMock
      )
    ).to.be.true
  })

  it("should get managed grant addresses for the provided grantee", async () => {
    // given
    const granteeAddress = "0x2A489EacBf4de172B4018D2b4a405F05C400f530"
    const granteeAddress2 = "0xeA56cDa2dfE567EF5ac9F5aDa6cF62036D298050"
    const grantAddress1 = "0x9bB49eEB07FDFE8EC10c49e6668d7E55F902BcF7"
    const grantAddress2 = "0x942e505ce0194b42dc87D06e9a6dc2dcF8f77733"
    const managedGrantCreatedEventsMock = [
      { returnValues: { grantAddress: grantAddress1 } },
      { returnValues: { grantAddress: grantAddress2 } },
    ]
    const makeCallManagedGrantStub = sandbox
      .stub()
      .onFirstCall()
      .returns(granteeAddress)
      .onSecondCall()
      .returns(granteeAddress2)

    sandbox
      .stub(ContractFactory, "new")
      .callsFake(() => ({ makeCall: makeCallManagedGrantStub }))

    const ethContractStub = sandbox.stub()
    sandbox
      .stub(keep, "config")
      .value({ web3: { eth: { Contract: ethContractStub } } })

    const manageGrantFactoryStub = sandbox
      .stub(keep.managedGrantFactoryContract, "getPastEvents")
      .returns(managedGrantCreatedEventsMock)

    // when
    const managedGrantAddresses = await keep.getGranteeManagedGrantAddresses(
      granteeAddress
    )

    // then
    expect(manageGrantFactoryStub.calledWithExactly("ManagedGrantCreated")).to
      .be.true
    expect(managedGrantAddresses.length).eq(1)
    expect(makeCallManagedGrantStub.callCount).eq(2)
    expect(ethContractStub.callCount).eq(2)
    expect(ethContractStub.getCall(0).args[1]).eq(grantAddress1)
    expect(ethContractStub.getCall(1).args[1]).eq(grantAddress2)
  })

  it("should return grant details by id", async () => {
    // given
    const grantId = "1"
    const getGrantMockData = {
      amount: "200",
      withdrawn: "100",
      staked: "100",
      revokedAmount: "0",
      revokedAt: "1",
      grantee: operatorAddress,
    }
    const getGrantUnlockingScheduleMockData = {
      grantManager: ownerAddress,
      duration: "1",
      start: "1",
      cliff: "1",
      policy: "0x0",
    }
    const unlcokedAmountMockData = "200"
    const withdrawableMockData = "300"
    const availableToStakeMockData = "400"
    const stub = sandbox
      .stub(keep.tokenGrantContract, "makeCall")
      .onFirstCall()
      .returns(getGrantMockData)
      .onSecondCall()
      .returns(getGrantUnlockingScheduleMockData)
      .onThirdCall()
      .returns(unlcokedAmountMockData)
      .onCall(3)
      .returns(withdrawableMockData)
      .onCall(4)
      .returns(availableToStakeMockData)

    // when
    const grantDetails = await keep.getGrantDetails(grantId)

    // then
    expect(stub.callCount).eq(5)
    expect(stub.getCall(0).calledWithExactly("getGrant", grantId)).to.be.true
    expect(
      stub.getCall(1).calledWithExactly("getGrantUnlockingSchedule", grantId)
    ).to.be.true
    expect(stub.getCall(2).calledWithExactly("unlockedAmount", grantId)).to.be
      .true
    expect(stub.getCall(3).calledWithExactly("withdrawable", grantId)).to.be
      .true
    expect(stub.getCall(4).calledWithExactly("availableToStake", grantId)).to.be
      .true
    expect(grantDetails).to.deep.eq({
      id: grantId,
      unlocked: unlcokedAmountMockData,
      released: getGrantMockData.withdrawn,
      readyToRelease: withdrawableMockData,
      availableToStake: availableToStakeMockData,
      ...getGrantUnlockingScheduleMockData,
      ...getGrantMockData,
    })
  })

  it("should authorize BondedECDSAKeepFactory contract", () => {
    const bondedECDSAKeepFactoryMockAddress = "0x0"
    sandbox
      .stub(keep, "bondedECDSAKeepFactoryContract")
      .value({ address: bondedECDSAKeepFactoryMockAddress })

    const stub = sandbox.stub(keep.tokenStakingContract, "sendTransaction")

    keep.authorizeBondedECDSAKeepFactory(operatorAddress)

    expect(
      stub.calledWithExactly(
        "authorizeOperatorContract",
        operatorAddress,
        bondedECDSAKeepFactoryMockAddress
      )
    ).to.be.true
  })
})
