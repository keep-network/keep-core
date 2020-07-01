import chai from "chai"
const { expect } = chai
import sinon from "sinon"
import {
  StakeOwnedStartegy,
  StakeGrantStrategy,
  StakeMangedGrantStrategy,
  StakingManager,
} from "../src/staking.js"

describe("Staking strategies test", () => {
  let conctractWrapperMock
  let stakingContractAddressMock
  let amountMock
  let delegationDataMock

  beforeEach(() => {
    conctractWrapperMock = { sendTransaction: sinon.spy() }
    stakingContractAddressMock = "0x0"
    amountMock = 10000
    delegationDataMock = "0x12345"
  })

  it("should stake via KEEPToken contract", () => {
    const stakeOwned = new StakeOwnedStartegy(conctractWrapperMock)

    stakeOwned.stake(stakingContractAddressMock, amountMock, delegationDataMock)

    expect(
      conctractWrapperMock.sendTransaction.calledOnceWithExactly(
        "approveAndCall",
        stakingContractAddressMock,
        amountMock,
        delegationDataMock
      )
    ).to.be.true
  })

  it("should stake via TokenGrant contract", () => {
    const tokenGrantId = 1
    const stakeGrant = new StakeGrantStrategy(
      conctractWrapperMock,
      tokenGrantId
    )

    stakeGrant.stake(stakingContractAddressMock, amountMock, delegationDataMock)

    expect(
      conctractWrapperMock.sendTransaction.calledOnceWithExactly(
        "stake",
        tokenGrantId,
        stakingContractAddressMock,
        amountMock,
        delegationDataMock
      )
    ).to.be.true
  })

  it("should stake via ManagedGrant contract", () => {
    const stakeGrant = new StakeMangedGrantStrategy(conctractWrapperMock)

    stakeGrant.stake(stakingContractAddressMock, amountMock, delegationDataMock)

    expect(
      conctractWrapperMock.sendTransaction.calledOnceWithExactly(
        "stake",
        stakingContractAddressMock,
        amountMock,
        delegationDataMock
      )
    ).to.be.true
  })
})

describe("StakingManager test", () => {
  it("shuld call stake from the staking starategy", () => {
    const stakeOwnedstarategyMock = { stake: sinon.spy() }
    const mockData = {
      stakingContractAddress: "0x0",
      amount: 100,
      beneficiaryAddress: "0x1",
      operatorAddress: "0x2",
      authorizerAddress: "0x3",
    }
    const expectedExtraData =
      "0x" +
      Buffer.concat([
        Buffer.from(mockData.beneficiaryAddress.substr(2), "hex"),
        Buffer.from(mockData.operatorAddress.substr(2), "hex"),
        Buffer.from(mockData.authorizerAddress.substr(2), "hex"),
      ]).toString("hex")

    StakingManager.stake(mockData, stakeOwnedstarategyMock)

    expect(
      stakeOwnedstarategyMock.stake.calledOnceWithExactly(
        mockData.stakingContractAddress,
        mockData.amount,
        expectedExtraData
      )
    ).to.be.true
  })
})
