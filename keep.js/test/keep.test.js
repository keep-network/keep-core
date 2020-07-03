import chai from "chai"
const { expect } = chai
import sinon from "sinon"
import KEEP, { contracts } from "../src/keep.js"
import ContractFactory, { ContractWrapper } from "../src/contract-wrapper.js"
import { TokenStakingConstants } from "../src/constants.js"
import { ContractWrapperMock } from "./mocks.js"

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

  before(async () => {
    sinon
      .stub(ContractFactory, "createContractInstance")
      .callsFake(() => Promise.resolve(ContractWrapperMock))

    sinon.stub(TokenStakingConstants, "initialize").callsFake(() =>
      Promise.resolve({
        minimumStake: "1",
        initializationPeriod: "2",
        undelegationPeriod: "3",
      })
    )

    keep = await KEEP.initialize({})
  })

  it("should beneficiary", async () => {
    const mockAddress = "0x0"
    await keep.beneficiaryOf(mockAddress)

    expect(
      keep.tokenStakingContract.makeCall.calledOnceWithExactly(
        "beneficiaryOf",
        mockAddress
      )
    ).to.be.true
  })
})
