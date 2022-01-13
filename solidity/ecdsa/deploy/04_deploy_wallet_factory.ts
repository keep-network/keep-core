import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()

  const SortitionPool = await deployments.get("SortitionPool")
  const T = await deployments.get("T")
  let TokenStaking = await deployments.get("TokenStaking")
  const DKGValidator = await deployments.get("DKGValidator")
  const MasterWallet = await deployments.get("MasterWallet")

  const DKG = await deployments.deploy("DKG", {
    from: deployer,
    log: true,
  })

  // TODO: StakingStub contract should be replaced by actual implementation of
  // TokenStaking contract, as soon as integration is implemented.
  if (deployments.getNetworkName() === "hardhat") {
    console.log("deploying StakingStub contract instead of TokenStaking")
    TokenStaking = await deployments.deploy("StakingStub", {
      from: deployer,
    })
  }

  const WalletFactory = await deployments.deploy("WalletFactory", {
    contract:
      deployments.getNetworkName() === "hardhat"
        ? "WalletFactoryStub"
        : undefined,
    from: deployer,
    args: [
      SortitionPool.address,
      T.address,
      TokenStaking.address,
      DKGValidator.address,
      MasterWallet.address,
    ],
    libraries: { DKG: DKG.address },
    log: true,
  })

  await helpers.ownable.transferOwnership(
    "SortitionPool",
    WalletFactory.address,
    deployer
  )

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "WalletFactory",
      address: WalletFactory.address,
    })
  }
}

export default func

func.tags = ["WalletFactory"]
func.dependencies = [
  "SortitionPool",
  "T",
  "TokenStaking",
  "DKGValidator",
  "MasterWallet",
]
