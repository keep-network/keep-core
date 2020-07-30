export class StakeOwnedStartegy {
  constructor(contract) {
    this.keepTokenContract = contract
  }

  stake(stakingContractAddress, amount, delegationExtraData) {
    return this.keepTokenContract.sendTransaction(
      "approveAndCall",
      stakingContractAddress,
      amount,
      delegationExtraData
    )
  }
}

export class StakeGrantStrategy {
  constructor(contract, tokenGrantId) {
    this.tokenGrantContract = contract
    this.tokenGrantId = tokenGrantId
  }

  stake(stakingContractAddress, amount, delegationExtraData) {
    return this.tokenGrantContract.sendTransaction(
      "stake",
      this.tokenGrantId,
      stakingContractAddress,
      amount,
      delegationExtraData
    )
  }
}

export class StakeMangedGrantStrategy {
  constructor(contract) {
    this.managedGrantContract = contract
  }

  stake(stakingContractAddress, amount, delegationExtraData) {
    return this.managedGrantContract.sendTransaction(
      "stake",
      stakingContractAddress,
      amount,
      delegationExtraData
    )
  }
}

export class StakingManager {
  static async stake(data, stakingStrategy) {
    const {
      stakingContractAddress,
      amount,
      beneficiaryAddress,
      operatorAddress,
      authorizerAddress,
    } = data
    const extraData =
      "0x" +
      Buffer.concat([
        Buffer.from(beneficiaryAddress.substr(2), "hex"),
        Buffer.from(operatorAddress.substr(2), "hex"),
        Buffer.from(authorizerAddress.substr(2), "hex"),
      ]).toString("hex")

    await stakingStrategy.stake(stakingContractAddress, amount, extraData)
  }
}
