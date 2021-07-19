import CoveragePoolV1 from "../coverage-pool"
import { Token } from "../../../utils/token.utils"
import { RewardsPoolArtifact } from "../contracts"
import { APYCalculator } from "../helper"

jest.mock("../contracts", () => ({
  RewardsPoolArtifact: {
    abi: [],
  },
}))
jest.mock("../helper", () => ({
  APYCalculator: {
    calculatePoolRewardRate: jest.fn(),
    calculateAPY: jest.fn(),
  },
}))

const createMockedContract = (address) => ({
  makeCall: jest.fn(),
  address: address,
  getPastEvents: jest.fn(),
})

describe("Test CoveragePoolV1 lib", () => {
  /** @type {CoveragePoolV1} */
  let coveragePoolV1
  beforeEach(() => {
    const assetPoolContract = createMockedContract("0x9")
    const covTokenContract = createMockedContract("0x7")
    const collateralToken = createMockedContract("0x6")
    const exchangeService = {
      getKeepTokenPriceInUSD: jest.fn(),
    }
    const web3 = {
      createContractInstance: jest.fn(),
    }

    coveragePoolV1 = new CoveragePoolV1(
      assetPoolContract,
      covTokenContract,
      collateralToken,
      exchangeService,
      web3
    )
  })

  it("should return the total supply of the cov token", async () => {
    const mockedResult = "1000"
    coveragePoolV1.covTokenContract.makeCall.mockResolvedValue(mockedResult)

    const result = await coveragePoolV1.covTotalSupply()

    expect(coveragePoolV1.covTokenContract.makeCall).toHaveBeenCalledWith(
      "totalSupply"
    )
    expect(result).toEqual(mockedResult)
  })

  it("should return the balance of te cov tokens for the provided address", async () => {
    const mockedResult = "1000"
    const mockedAddress = "0x0"
    coveragePoolV1.covTokenContract.makeCall.mockResolvedValue(mockedResult)

    const result = await coveragePoolV1.covBalanceOf(mockedAddress)

    expect(coveragePoolV1.covTokenContract.makeCall).toHaveBeenCalledWith(
      "balanceOf",
      mockedAddress
    )
    expect(result).toEqual(mockedResult)
  })

  it("should return estimated rewards corretly if currently deposited collateral token balance is greater than 0", async () => {
    const shareOfPool = 0.4
    const address = "0x123"
    const spyOnTvl = jest
      .spyOn(coveragePoolV1, "totalValueLocked")
      .mockResolvedValue(Token.fromTokenUnit(150).toString())

    const mockedToEvents = [
      { returnValues: { value: Token.fromTokenUnit(25).toString() } },
      { returnValues: { value: Token.fromTokenUnit(25).toString() } },
    ]

    const mockedFromEvents = []

    const spyOnGetPastTransferEvent = jest
      .spyOn(coveragePoolV1.collateralToken, "getPastEvents")
      .mockResolvedValueOnce(mockedToEvents)
      .mockResolvedValueOnce(mockedFromEvents)

    const epxectedResult = Token.fromTokenUnit(10).toString()

    const result = await coveragePoolV1.estimatedRewards(address, shareOfPool)

    expect(spyOnTvl).toHaveBeenCalled()
    expect(spyOnGetPastTransferEvent).toHaveBeenNthCalledWith(1, "Transfer", {
      from: address,
      to: coveragePoolV1.assetPoolContract.address,
    })
    expect(spyOnGetPastTransferEvent).toHaveBeenNthCalledWith(2, "Transfer", {
      from: coveragePoolV1.assetPoolContract.address,
      to: address,
    })
    expect(result).toEqual(epxectedResult.toString())
  })

  it("the estimated reward balance equlas 0 if the share of pool equlas 0 ", async () => {
    const shareOfPool = 0
    const address = "0x123"
    const spyOnTvl = jest.spyOn(coveragePoolV1, "totalValueLocked")

    const spyOnGetPastTransferEvent = jest.spyOn(
      coveragePoolV1.collateralToken,
      "getPastEvents"
    )

    const epxectedResult = "0"
    const result = await coveragePoolV1.estimatedRewards(address, shareOfPool)

    expect(spyOnTvl).not.toHaveBeenCalled()
    expect(spyOnGetPastTransferEvent).not.toHaveBeenCalled()
    expect(result).toEqual(epxectedResult.toString())
  })

  it("should return the share of pool", async () => {
    const covBalanceOf = 35
    const covTotalSupply = 100
    const expectedResult = covBalanceOf / covTotalSupply

    const result = await coveragePoolV1.shareOfPool(
      covTotalSupply,
      covBalanceOf
    )

    expect(result).toEqual(expectedResult.toString())
  })

  it("should return the total value locked", async () => {
    const mockedResult = "100"
    coveragePoolV1.assetPoolContract.makeCall.mockResolvedValue(mockedResult)

    const result = await coveragePoolV1.totalValueLocked()

    expect(coveragePoolV1.assetPoolContract.makeCall).toHaveBeenCalledWith(
      "totalValue"
    )
    expect(result).toEqual(mockedResult)
  })

  it("should return the estimated collateral token balance", async () => {
    const mockedCollateralBalance = 100
    coveragePoolV1.collateralToken.makeCall.mockResolvedValue(
      mockedCollateralBalance
    )
    const shareOfPool = 0.35

    const result = await coveragePoolV1.estimatedCollateralTokenBalance(
      shareOfPool
    )

    expect(coveragePoolV1.collateralToken.makeCall).toHaveBeenCalledWith(
      "balanceOf",
      coveragePoolV1.assetPoolContract.address
    )
    expect(result).toEqual((mockedCollateralBalance * shareOfPool).toString())
  })

  it("should return the asset pool collateral token balance", async () => {
    const mockedBalance = Token.fromTokenUnit("300").toString()
    const spy = jest
      .spyOn(coveragePoolV1.collateralToken, "makeCall")
      .mockResolvedValue(mockedBalance)

    const result = await coveragePoolV1.assetPoolCollateralTokenBalance()

    expect(spy).toHaveBeenCalledWith(
      "balanceOf",
      coveragePoolV1.assetPoolContract.address
    )
    expect(result).toEqual(mockedBalance)
  })

  describe

  it("should return the reward pool contract", async () => {
    const mockedRewardPoolContractAddress = "0x123"
    const spyOnRewardsPool = jest
      .spyOn(coveragePoolV1.assetPoolContract, "makeCall")
      .mockResolvedValue(mockedRewardPoolContractAddress)

    const mockedRewardPoolContract = createMockedContract(
      mockedRewardPoolContractAddress
    )
    const spyOnCreateContract = jest
      .spyOn(coveragePoolV1.web3, "createContractInstance")
      .mockReturnValue(mockedRewardPoolContract)

    const result = await coveragePoolV1.getRewardPoolContract()
    const result2 = await coveragePoolV1.getRewardPoolContract()

    expect(spyOnRewardsPool).toHaveBeenCalledWith("rewardsPool")
    expect(spyOnRewardsPool).toHaveBeenCalledTimes(1)
    expect(spyOnCreateContract).toHaveBeenCalledWith(
      RewardsPoolArtifact.abi,
      mockedRewardPoolContractAddress,
      coveragePoolV1.assetPoolContract.deploymentTxnHash,
      coveragePoolV1.assetPoolContract.deployedAtBlock
    )
    expect(spyOnCreateContract).toHaveBeenCalledTimes(1)
    expect(result).toEqual(mockedRewardPoolContract)
    expect(result).toEqual(result2)
    expect(coveragePoolV1._rewardPoolContract).toEqual(result)
  })

  it("should return the reward rate of the rewards pool contract", async () => {
    const mockedRewardsPoolContract = createMockedContract("0x1234")
    const spyOnRewardPoolContract = jest
      .spyOn(coveragePoolV1, "getRewardPoolContract")
      .mockResolvedValue(mockedRewardsPoolContract)

    const rewardRate = Token.fromTokenUnit("1000").toString()
    const spyOnMakeCall = jest
      .spyOn(mockedRewardsPoolContract, "makeCall")
      .mockResolvedValue(rewardRate)

    const result = await coveragePoolV1.rewardPoolRewardRate()

    expect(spyOnRewardPoolContract).toHaveBeenCalled()
    expect(spyOnMakeCall).toHaveBeenCalledWith("rewardRate")
    expect(result).toEqual(rewardRate)
  })

  it("should return the reward pool per week", async () => {
    const rewardRate = Token.fromTokenUnit("1000")
    const spyOnRewardRate = jest
      .spyOn(coveragePoolV1, "rewardPoolRewardRate")
      .mockResolvedValue(rewardRate.toString())

    const result = await coveragePoolV1.rewardPoolPerWeek()

    expect(spyOnRewardRate).toHaveBeenCalled()
    expect(result).toEqual(Token.toTokenUnit(rewardRate).multipliedBy(604800))
  })

  it("should calculate apy correctly", async () => {
    const mockedTotalSupply = Token.fromTokenUnit(100)
    const spyOnTotalSupply = jest
      .spyOn(coveragePoolV1, "assetPoolCollateralTokenBalance")
      .mockResolvedValue(mockedTotalSupply)

    const mockedRewardPoolPerWeek = "150"
    const spyOnRewarPoolPerWeek = jest
      .spyOn(coveragePoolV1, "rewardPoolPerWeek")
      .mockResolvedValue(mockedRewardPoolPerWeek)

    const mockedPriceInUSD = 0.5
    const spyOnGetPriceInUSD = jest
      .spyOn(coveragePoolV1.exchangeService, "getKeepTokenPriceInUSD")
      .mockResolvedValue(mockedPriceInUSD)

    const mockedPoolRewardRate = 0.2
    const spyOnPoolRewardRate = jest
      .spyOn(APYCalculator, "calculatePoolRewardRate")
      .mockReturnValue(mockedPoolRewardRate)

    const mockedAPY = 0.99
    const spyOnAPY = jest
      .spyOn(APYCalculator, "calculateAPY")
      .mockReturnValue(mockedAPY)

    const result = await coveragePoolV1.apy()

    expect(spyOnTotalSupply).toHaveBeenCalled()
    expect(spyOnRewarPoolPerWeek).toHaveBeenCalled()
    expect(spyOnGetPriceInUSD).toHaveBeenCalled()
    expect(spyOnPoolRewardRate).toHaveBeenCalledWith(
      mockedPriceInUSD,
      mockedRewardPoolPerWeek,
      Token.toTokenUnit(mockedTotalSupply).multipliedBy(mockedPriceInUSD)
    )
    expect(spyOnAPY).toHaveBeenCalledWith(mockedPoolRewardRate)
    expect(result).toEqual(mockedAPY.toString())
  })

  it("should return the total allocated rewards", async () => {
    const mockedRewardsPoolContract = createMockedContract("0x0123")
    const spyOnGetRewardsPoolContract = jest
      .spyOn(coveragePoolV1, "getRewardPoolContract")
      .mockResolvedValue(mockedRewardsPoolContract)

    const mockedEvents = [
      {
        returnValues: {
          amount: Token.fromTokenUnit(30).toString(),
        },
      },
      {
        returnValues: {
          amount: Token.fromTokenUnit(30).toString(),
        },
      },
    ]
    const spyOnGetPastEvents = jest
      .spyOn(mockedRewardsPoolContract, "getPastEvents")
      .mockResolvedValue(mockedEvents)

    const result = await coveragePoolV1.totalAllocatedRewards()

    expect(spyOnGetRewardsPoolContract).toHaveBeenCalled()
    expect(spyOnGetPastEvents).toHaveBeenCalledWith("RewardToppedUp")
    expect(result.toString()).toEqual(Token.fromTokenUnit(60).toString())
  })
})
