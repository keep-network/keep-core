import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()

  const T = await deployments.get("T")
  const TokenStaking = await deployments.get("TokenStaking")
  const ReimbursementPool = await deployments.get("ReimbursementPool")
  const BeaconSortitionPool = await deployments.get("BeaconSortitionPool")

  const BLS = await deployments.deploy("BLS", {
    from: deployer,
    log: true,
  })

  const BeaconAuthorization = await deployments.deploy("BeaconAuthorization", {
    from: deployer,
    log: true,
  })

  const BeaconDkg = await deployments.deploy("BeaconDkg", {
    from: deployer,
    log: true,
  })

  const BeaconInactivity = await deployments.deploy("BeaconInactivity", {
    from: deployer,
    log: true,
  })

  const BeaconDkgValidator = await deployments.deploy("BeaconDkgValidator", {
    from: deployer,
    args: [BeaconSortitionPool.address],
    log: true,
  })

  const RandomBeacon = await deployments.deploy("RandomBeacon", {
    contract:
      deployments.getNetworkName() === "hardhat"
        ? "RandomBeaconStub"
        : undefined,
    from: deployer,
    args: [
      BeaconSortitionPool.address,
      T.address,
      TokenStaking.address,
      BeaconDkgValidator.address,
      ReimbursementPool.address,
    ],
    libraries: {
      BLS: BLS.address,
      BeaconAuthorization: BeaconAuthorization.address,
      BeaconDkg: BeaconDkg.address,
      BeaconInactivity: BeaconInactivity.address,
    },
    log: true,
  })

  await helpers.ownable.transferOwnership(
    "BeaconSortitionPool",
    RandomBeacon.address,
    deployer
  )

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "RandomBeacon",
      address: RandomBeacon.address,
    })
  }
}

export default func

func.tags = ["RandomBeacon"]
func.dependencies = [
  "T",
  "TokenStaking",
  "ReimbursementPool",
  "BeaconSortitionPool",
]
