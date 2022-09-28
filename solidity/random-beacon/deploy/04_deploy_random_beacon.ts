import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction, DeployOptions } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()
  const { to1e18 } = helpers.number

  const minimumAuthorization = to1e18(40_000) // 40k T
  const authorizationDecreaseDelay = 1_209_600 // 1 209 600 seconds = 14 days
  const authorizationDecreaseChangePeriod = 1_209_600 // 1 209 600 seconds = 14 days

  const T = await deployments.get("T")
  const TokenStaking = await deployments.get("TokenStaking")
  const ReimbursementPool = await deployments.get("ReimbursementPool")
  const BeaconSortitionPool = await deployments.get("BeaconSortitionPool")
  const BeaconDkgValidator = await deployments.get("BeaconDkgValidator")

  const deployOptions: DeployOptions = {
    from: deployer,
    log: true,
    waitConfirmations: 1,
  }

  const BLS = await deployments.deploy("BLS", deployOptions)

  const BeaconAuthorization = await deployments.deploy(
    "BeaconAuthorization",
    deployOptions
  )

  const BeaconDkg = await deployments.deploy("BeaconDkg", deployOptions)

  const BeaconInactivity = await deployments.deploy(
    "BeaconInactivity",
    deployOptions
  )

  const RandomBeacon = await deployments.deploy("RandomBeacon", {
    contract:
      process.env.TEST_USE_STUBS_BEACON === "true"
        ? "RandomBeaconStub"
        : undefined,
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
    ...deployOptions,
  })

  await deployments.execute(
    "RandomBeacon",
    { from: deployer, log: true, waitConfirmations: 1 },
    "updateAuthorizationParameters",
    minimumAuthorization,
    authorizationDecreaseDelay,
    authorizationDecreaseChangePeriod
  )

  await helpers.ownable.transferOwnership(
    "BeaconSortitionPool",
    RandomBeacon.address,
    deployer
  )

  if (hre.network.tags.etherscan) {
    await helpers.etherscan.verify(BLS)
    await helpers.etherscan.verify(BeaconAuthorization)
    await helpers.etherscan.verify(BeaconDkg)
    await helpers.etherscan.verify(BeaconInactivity)
    await helpers.etherscan.verify(RandomBeacon)
  }

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
  "BeaconDkgValidator",
]
