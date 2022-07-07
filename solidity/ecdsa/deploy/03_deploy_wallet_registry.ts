import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, ethers, helpers } = hre
  const { deployer } = await getNamedAccounts()

  const SortitionPool = await deployments.get("EcdsaSortitionPool")
  const TokenStaking = await deployments.get("TokenStaking")
  const ReimbursementPool = await deployments.get("ReimbursementPool")
  const RandomBeacon = await deployments.get("RandomBeacon")
  const EcdsaDkgValidator = await deployments.get("EcdsaDkgValidator")

  const EcdsaInactivity = await deployments.deploy("EcdsaInactivity", {
    from: deployer,
    log: true,
  })

  const walletRegistry = await helpers.upgrades.deployProxy("WalletRegistry", {
    contractName:
      process.env.TEST_USE_STUBS_ECDSA === "true"
        ? "WalletRegistryStub"
        : undefined,
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
  "RandomBeacon",
  "EcdsaSortitionPool",
  "TokenStaking",
  "EcdsaDkgValidator",
]
