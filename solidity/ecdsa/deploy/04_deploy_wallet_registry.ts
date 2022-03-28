import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()

  const SortitionPool = await deployments.get("SortitionPool")
  const TokenStaking = await deployments.get("TokenStaking")
  const ReimbursementPool = await deployments.get("ReimbursementPool")
  const EcdsaDkgValidator = await deployments.get("EcdsaDkgValidator")

  // TODO: RandomBeaconStub contract should be replaced by actual implementation of
  // RandomBeacon contract, once @keep-network/random-beacon hardhat deployments
  // scripts are implemented.
  console.log("deploying RandomBeaconStub contract instead of RandomBeacon")
  const RandomBeacon = await deployments.deploy("RandomBeaconStub", {
    from: deployer,
    log: true,
  })

  const EcdsaDkg = await deployments.deploy("EcdsaDkg", {
    from: deployer,
    log: true,
  })

  const WalletRegistry = await deployments.deploy("WalletRegistry", {
    contract:
      deployments.getNetworkName() === "hardhat"
        ? "WalletRegistryStub"
        : undefined,
    from: deployer,
    args: [
      SortitionPool.address,
      TokenStaking.address,
      EcdsaDkgValidator.address,
      RandomBeacon.address,
      ReimbursementPool.address,
    ],
    libraries: { EcdsaDkg: EcdsaDkg.address },
    log: true,
  })

  await helpers.ownable.transferOwnership(
    "SortitionPool",
    WalletRegistry.address,
    deployer
  )

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "WalletRegistry",
      address: WalletRegistry.address,
    })
  }
}

export default func

func.tags = ["WalletRegistry"]
func.dependencies = ["SortitionPool", "TokenStaking", "EcdsaDkgValidator"]
