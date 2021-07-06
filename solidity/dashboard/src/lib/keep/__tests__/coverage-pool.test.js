import CoveragePoolV1 from "../coverage-pool"

const createMockedContract = (address) => ({
  makeCall: jest.fn(),
  address: address,
})

describe("Test CoveragePoolV1 lib", () => {
  let coveragePoolV1
  beforeEach(() => {
    const assetPoolContract = createMockedContract("0x9")
    const rewardPoolContract = createMockedContract("0x8")
    const covTokenContract = createMockedContract("0x7")
    const collateralToken = createMockedContract("0x6")

    coveragePoolV1 = new CoveragePoolV1(
      assetPoolContract,
      rewardPoolContract,
      covTokenContract,
      collateralToken
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

  it("should return estimated rewards", async () => {
    const shareOfPool = 0.4
    const spyOnTvl = jest
      .spyOn(coveragePoolV1, "totalValueLocked")
      .mockResolvedValue(150)
    const estimatedKeepBalance = 50
    // tvl * shareOfPool - estimatedBalance
    const epxectedResult = 10

    const result = await coveragePoolV1.estimatedRewards(
      shareOfPool,
      estimatedKeepBalance
    )

    expect(spyOnTvl).toHaveBeenCalled()
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
})
