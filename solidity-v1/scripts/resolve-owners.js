const TokenGrant = artifacts.require("./TokenGrant.sol")
const TokenStaking = artifacts.require("./TokenStaking.sol")
const ManagedGrantFactory = artifacts.require("./ManagedGrantFactory.sol")

module.exports = async function () {
  try {
    const tokenGrant = await TokenGrant.at(
      "0x175989c71fd023d580c65f5dc214002687ff88b7"
    )
    const tokenStaking = await TokenStaking.at(
      "0x1293a54e160d1cd7075487898d65266081a15458"
    )
    const managedGrantFactory = await ManagedGrantFactory.at(
      "0x43cf9e26857b188868051bdcfacedbb38531964e"
    )

    const allEventsOpts = { fromBlock: 0, toBlock: "latest" }
    const stakeDelegatedEvents = await tokenStaking.getPastEvents(
      "StakeDelegated",
      allEventsOpts
    )
    const managedGrantCreatedEvents = await managedGrantFactory.getPastEvents(
      "ManagedGrantCreated",
      allEventsOpts
    )

    for (let i = 0; i < stakeDelegatedEvents.length; i++) {
      const operator = stakeDelegatedEvents[i].args["operator"]

      const delegationInfo = await tokenStaking.getDelegationInfo.call(operator)

      // skip those who undelegated
      if (delegationInfo.undelegatedAt != 0) {
        continue
      }

      // skip those who canceled their delegation
      if (delegationInfo.amount == 0) {
        continue
      }

      const owner = await tokenStaking.ownerOf.call(operator)
      const grantStake = await tokenGrant.grantStakes(operator)

      if (grantStake == "0x0000000000000000000000000000000000000000") {
        // it is liquid token delegation
        console.log(`[LI], ${owner}, ${operator}`)
      } else {
        // it is grant delegation
        const grantStakeDetails = await tokenGrant.getGrantStakeDetails.call(
          operator
        )
        const grantId = grantStakeDetails.grantId
        const stakingContract = grantStakeDetails.stakingContract

        const grant = await tokenGrant.getGrant.call(grantId)
        const grantee = grant.grantee

        let grantType = "[SG]" // assume standard grant by default

        for (let j = 0; j < managedGrantCreatedEvents.length; j++) {
          if (grantee == managedGrantCreatedEvents[j].args["grantAddress"]) {
            // it is managed grant delegation
            grantType = "[MG]"
            break
          }
        }

        console.log(
            `${grantType}, ${owner}, ${operator}, ${grantId}, ${stakingContract}, ${grantee}`
        )
      }
    }
  } catch (err) {
    console.error("unexpected error:", err)
    process.exit(1)
  }

  process.exit()
}
