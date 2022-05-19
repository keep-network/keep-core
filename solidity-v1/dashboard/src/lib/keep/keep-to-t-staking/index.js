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
