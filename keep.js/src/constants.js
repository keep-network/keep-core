export class TokenStakingConstants {
  static async initialize(tokenStakingContract) {
    const toCall = new Map([
      ["minimumStake", "minimumStake"],
      ["undelegationPeriod", "undelegationPeriod"],
      ["initializationPeriod", "initializationPeriod"],
    ])

    const constants = {}
    for (const [methodName, propertyName] of toCall) {
      constants[propertyName] = await tokenStakingContract.makeCall(methodName)
    }

    return new TokenStakingConstants(constants)
  }

  constructor(constants) {
    Object.assign(this, constants)

    this.minimumStake
    this.undelegationPeriod
    this.initializationPeriod
  }
}
