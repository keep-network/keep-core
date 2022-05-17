import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, ethers, helpers } = hre
  const { deployer } = await getNamedAccounts()
  const { log } = deployments

  const SortitionPool = await deployments.get("EcdsaSortitionPool")
  const TokenStaking = await deployments.get("TokenStaking")
  const ReimbursementPool = await deployments.get("ReimbursementPool")
  const EcdsaDkgValidator = await deployments.get("EcdsaDkgValidator")

  // TODO: RandomBeaconStub contract should be replaced by actual implementation of
  // RandomBeacon contract, once @keep-network/random-beacon hardhat deployments
  // scripts are implemented.
  log("deploying RandomBeaconStub contract instead of RandomBeacon")
  const RandomBeacon = await deployments.deploy("RandomBeaconStub", {
    from: deployer,
    log: true,
  })

  const EcdsaInactivity = await deployments.deploy("EcdsaInactivity", {
    from: deployer,
    log: true,
  })

  const walletRegistry = await helpers.upgrades.deployProxy("WalletRegistry", {
    contractName:
      process.env.TEST_USE_STUBS === "true" ? "WalletRegistryStub" : undefined,
    initializerArgs: [
      EcdsaDkgValidator.address,
      RandomBeacon.address,
      ReimbursementPool.address,
    ],
    factoryOpts: {
      signer: await ethers.getSigner(deployer),
      libraries: {
        EcdsaInactivity: EcdsaInactivity.address,
      },
    },
    proxyOpts: {
      constructorArgs: [SortitionPool.address, TokenStaking.address],
      unsafeAllow: ["external-library-linking"],
      kind: "transparent",
    },
  })

  await helpers.ownable.transferOwnership(
    "EcdsaSortitionPool",
    walletRegistry.address,
    deployer
  )

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "WalletRegistry",
      address: walletRegistry.address,
    })
  }
}

export default func

func.tags = ["WalletRegistry"]
func.dependencies = [
  "ReimbursementPool",
  "EcdsaSortitionPool",
  "TokenStaking",
  "EcdsaDkgValidator",
]
