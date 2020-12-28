import Web3 from "web3"
import ProviderEngine from "web3-provider-engine"
import Subproviders from "@0x/subproviders"

// keep-core 1.1.2
import KeepRandomBeaconOperatorJson from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperator.json"

const engine = new ProviderEngine({pollingInterval: 1000})

engine.addProvider(
  // For address 0x94c0bD91530DF78ff41b9a7F9bCdd7E4730C7748.
  new Subproviders.PrivateKeyWalletSubprovider(
    "df5c8ed97b9d60ef0043fdc4f918804930c1647a7afade6bbf50b5f32226b2af"
  )
)
engine.addProvider(
  new Subproviders.RPCSubprovider(
    "https://mainnet.infura.io/v3/9b853f4554184e36ab15c027b1a6fa45"
  )
)

const web3 = new Web3(engine)

const keepRandomBeaconOperatorAbi = KeepRandomBeaconOperatorJson.abi
const keepRandomBeaconOperatorAddress =
  "0x70F2202D85a4F0Cad36e978976f84E982920A624"
const operator = new web3.eth.Contract(
  keepRandomBeaconOperatorAbi,
  keepRandomBeaconOperatorAddress
)

engine.start()

async function run() {
  const numberOfGroups = await operator.methods.numberOfGroups().call()
  console.log("numberOfGroups: ", numberOfGroups)

  for (let groupIndex = 0; groupIndex < numberOfGroups; groupIndex++) {
    const groupPubKey = await operator.methods
      .getGroupPublicKey(groupIndex)
      .call()
    const isGroupStale = await operator.methods.isStaleGroup(groupPubKey).call()

    if (isGroupStale) {
      const groupMembers = await operator.methods
        .getGroupMembers(groupPubKey)
        .call()

      const uniqueMembersInGroup = new Set()
      groupMembers.forEach((member) => {
        uniqueMembersInGroup.add(member)
      })

      console.log(`withdrawing rewards for group public key: ${groupPubKey}..`)
      uniqueMembersInGroup.forEach((memberAddress) => {
        console.log(`withdrawing rewards for member: ${memberAddress}..`)
        try {
          // await operator.methods.withdrawGroupMemberRewards(memberAddress, groupIndex).call()
        } catch (err) {
          console.log(
            `error occured while withdrawing rewards for ${memberAddress}..`,
            err
          )
        }
      })
      console.log(`\n`)
    }
  }
}

run()
  .then(() => {
    console.log("Withdrawal of rewards completed successfully")

    process.exit(0)
  })
  .catch((error) => {
    console.error("Withdrawal of rewards errored out: ", error)

    process.exit(1)
  })
