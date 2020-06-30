class Delegations {
  constructor() {}

  // getOwnerDelegations = async (ownerAddress) => {
  //   const ownerOperators = this.getOperatorsOf(ownerAddress)
  //   return await this.getDelegations(ownerOperators)
  // }

  // getDelegations = async (operatorAddresses) => {
  //   const delegations = []
  //   for (const operatorAddress of operatorAddresses) {
  //     const delegationInfo = await this.tokenStaking.makeCall(
  //       "getDelegationInfo",
  //       operatorAddress
  //     )
  //     const beneficiary = await this.tokenStaking.makeCall(
  //       "beneficiaryOf",
  //       operatorAddress
  //     )
  //     const authorizer = await this.tokenStaking.makeCall(
  //       "authorizerOf",
  //       operatorAddress
  //     )

  //     delegations.push({ ...delegationInfo, beneficiary, authorizer })
  //   }

  //   return delegations
  // }

  // getDelegationStatus(delegation) {
  //   const { amount, createdAt, undelegatedAt } = delegation
  //   let delegationStatus
  //   if (amount !== "0" && createdAt !== "0" && undelegatedAt !== "0") {
  //     // delegation undelegated
  //     delegationStatus = "UNDELEGATED"
  //   } else if (amount === "0" && createdAt !== "0" && undelegatedAt === "0") {
  //     // delegation canceled
  //     delegationStatus = "CANCELED"
  //   } else if (amount === "0" && createdAt !== "0" && undelegatedAt !== "0") {
  //     // delegation recovered
  //     delegationStatus = "RECOVERED"
  //   }

  //   return delegationStatus
  // }
}
