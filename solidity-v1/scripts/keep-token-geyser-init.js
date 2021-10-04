const KeepVault = artifacts.require("./geyser/KeepVault.sol")
const KeepToken = artifacts.require("./KeepToken.sol")
const BatchedPhasedEscrow = artifacts.require("./BatchedPhasedEscrow")
const KeepTokenGeyserRewardsEscrowBeneficiary = artifacts.require(
  "./KeepTokenGeyserRewardsEscrowBeneficiary"
)

module.exports = async function () {
  try {
    const accounts = await web3.eth.getAccounts()
    const keepToken = await KeepToken.deployed()
    const tokenGeyser = await KeepVault.deployed()
    const rewardsAmount = web3.utils.toWei("100000", "ether")

    const owner = accounts[0]

    const initialEscrowBalance = web3.utils.toWei("500000", "ether") // 500k KEEP

    const escrow = await BatchedPhasedEscrow.new(keepToken.address, {
      from: owner,
    })

    // Configure escrow beneficiary.
    const escrowBeneficiary = await KeepTokenGeyserRewardsEscrowBeneficiary.new(
      keepToken.address,
      tokenGeyser.address,
      {
        from: owner,
      }
    )

    await escrowBeneficiary.transferOwnership(escrow.address, {
      from: owner,
    })

    await escrow.approveBeneficiary(escrowBeneficiary.address, {
      from: owner,
    })

    await tokenGeyser.setRewardDistribution(escrowBeneficiary.address, {
      from: owner,
    })

    await keepToken.approveAndCall(escrow.address, initialEscrowBalance, [], {
      from: owner,
    })

    // Initiate withdraw.
    await escrow.batchedWithdraw([escrowBeneficiary.address], [rewardsAmount], {
      from: owner,
    })
  } catch (err) {
    console.error("unexpected error:", err)
    process.exit(1)
  }

  process.exit()
}
