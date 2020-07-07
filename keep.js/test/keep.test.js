import chai from "chai"
const { expect } = chai
import sinon from "sinon"
import KEEP, { contracts } from "../src/keep.js"
import ContractFactory, { ContractWrapper } from "../src/contract-wrapper.js"
import { TokenStakingConstants } from "../src/constants.js"

const operatorAddress = "0x0"
const beneficiaryAddress = "0x1"
const authorizerAddress = "0x2"

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
})
