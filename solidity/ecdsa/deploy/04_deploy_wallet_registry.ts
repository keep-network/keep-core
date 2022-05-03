import type { Contract } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"
import type { DeployResult } from "hardhat-deploy/dist/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, ethers, helpers } = hre
  const { deployer, esdm } = await getNamedAccounts()
  const { log } = deployments

  const SortitionPool = await deployments.get("SortitionPool")
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

  // FIXME: As a workaround for a bug in hardhat-gas-reporter #86 we need to provide
  // alternative deployment script to obtain a gas report.
  // #86: https://github.com/cgewecke/hardhat-gas-reporter/issues/86
  let walletRegistry: DeployResult | Contract
  if (process.env.GAS_REPORTER_BUG_WORKAROUND === "true") {
    walletRegistry = await deployments.deploy("WalletRegistry", {
      contract: "WalletRegistryStub",
      from: deployer,
      args: [SortitionPool.address, TokenStaking.address],
      libraries: {
        EcdsaInactivity: EcdsaInactivity.address,
      },
      proxy: {
        proxyContract: "TransparentUpgradeableProxy",
        viaAdminContract: "DefaultProxyAdmin",
        owner: esdm,
        execute: {
          init: {
            methodName: "initialize",
            args: [
              EcdsaDkgValidator.address,
              RandomBeacon.address,
              ReimbursementPool.address,
            ],
          },
        },
      },
      log: true,
    })
  } else {
    walletRegistry = await helpers.upgrades.deployProxy("WalletRegistry", {
      contractName:
        deployments.getNetworkName() === "hardhat"
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
  }

  await helpers.ownable.transferOwnership(
    "SortitionPool",
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
func.dependencies = ["SortitionPool", "TokenStaking", "EcdsaDkgValidator"]
