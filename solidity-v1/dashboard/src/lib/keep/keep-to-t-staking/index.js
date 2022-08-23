import BigNumber from "bignumber.js"

class KeepToTStaking {
  floatingPointDivisor = new BigNumber(10).pow(15)
  // You can check the ratio here: https://etherscan.io/address/0xE47c80e8c23f6B4A1aE41c34837a0599D5D16bb0#readContract
  keepToTExchangeRateInWei = 4783188631255016

  /**
   * @param {BaseContract} _thresholdStakingContract
   * @param {BaseContract} _thresholdKeepStakeContract
   * @param {BaseContract} _simplePREApplicationContract
   * @param {Web3LibWrapper} _web3
   */
  constructor(
    _thresholdStakingContract,
    _thresholdKeepStakeContract,
    _simplePREApplicationContract,
    _web3
  ) {
    this.thresholdStakingContract = _thresholdStakingContract
    this.thresholdKeepStakeContract = _thresholdKeepStakeContract
    this.simplePREApplicationContract = _simplePREApplicationContract
    this.web3 = _web3
  }

  getStakedEventsByOperator = async (operatorAddresses) => {
    return await this.thresholdStakingContract.getPastEvents("Staked", {
      stakingProvider: operatorAddresses,
    })
  }

  /**
   * Returns all tokens staked by the operator in T Network
   *
   * @param {string} operatorAddress
   * @return {Promise} Tokens staked in T Network as object:
   * {
   *    tStake: string,
   *    keepInTStake: string,
   *    nuInTStake: string,
   * }
   */
  tokensStakedInTNetwork = async (operatorAddress) => {
    return await this.thresholdStakingContract.makeCall(
      "stakes",
      operatorAddress
    )
  }

  /**
   * Checks if the operator has Keep tokens staked in T Network
   *
   * @param {string} operatorAddress
   * @return {Promise<boolean>} true if the operator has Keep tokens staked in
   * T Network and false otherwise
   */
  hasKeepTokensStakedInTNetwork = async (operatorAddress) => {
    const { keepInTStake } = await this.tokensStakedInTNetwork(operatorAddress)

    const amountOfKeepTokensStakedInBN = new BigNumber(keepInTStake)

    return !amountOfKeepTokensStakedInBN.eq("0")
  }

  getOperatorConfirmedEvents = async (operatorAddresses) => {
    return await this.simplePREApplicationContract.getPastEvents(
      "OperatorConfirmed",
      {
        stakingProvider: operatorAddresses,
      }
    )
  }

  resolveOwner = async (operatorAddress) => {
    return await this.thresholdKeepStakeContract.makeCall(
      "resolveOwner",
      operatorAddress
    )
  }

  toThresholdTokenAmount = (keepAmount) => {
    const amountInBN = new BigNumber(keepAmount)
    const wrappedRemainder = amountInBN.modulo(this.floatingPointDivisor)
    const convertibleAmount = amountInBN.minus(wrappedRemainder)

    return convertibleAmount
      .multipliedBy(this.keepToTExchangeRateInWei)
      .dividedBy(this.floatingPointDivisor)
      .toString()
  }

  fromThresholdTokenAmount = (tAmount) => {
    const amountInBN = new BigNumber(tAmount)
    const tRemainder = amountInBN.modulo(this.keepToTExchangeRateInWei)
    const convertibleAmount = amountInBN.minus(tRemainder)

    return convertibleAmount
      .multipliedBy(this.floatingPointDivisor)
      .dividedBy(this.keepToTExchangeRateInWei)
      .toString()
  }
}

export default KeepToTStaking
