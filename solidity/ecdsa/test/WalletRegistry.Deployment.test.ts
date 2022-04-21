/* eslint-disable @typescript-eslint/no-unused-expressions */
import { deployments, ethers, upgrades, helpers } from "hardhat"
import chai, { expect } from "chai"
import chaiAsPromised from "chai-as-promised"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type {
  ProxyAdmin,
  WalletRegistry,
  WalletRegistryGovernance,
} from "../typechain"
import type { TransparentUpgradeableProxy } from "../typechain/TransparentUpgradeableProxy"

chai.use(chaiAsPromised)

describe("WalletRegistry - Deployment", async () => {
  let deployer: SignerWithAddress
  let governance: SignerWithAddress
  let esdm: SignerWithAddress

  let walletRegistry: WalletRegistry
  let walletRegistryGovernance: WalletRegistryGovernance
  let walletRegistryProxy: TransparentUpgradeableProxy
  let proxyAdmin: ProxyAdmin

  before(async () => {
    await deployments.fixture()
    ;({ deployer, governance, esdm } = await ethers.getNamedSigners())

    walletRegistry = await ethers.getContract<WalletRegistry>("WalletRegistry")

    walletRegistryGovernance =
      await ethers.getContract<WalletRegistryGovernance>(
        "WalletRegistryGovernance"
      )

    walletRegistryProxy =
      await ethers.getContractAt<TransparentUpgradeableProxy>(
        "TransparentUpgradeableProxy",
        walletRegistry.address
      )

    proxyAdmin = (await upgrades.admin.getInstance()) as ProxyAdmin

    expect(deployer.address, "deployer is the same as governance").not.equal(
      governance.address
    )
  })

  it("should set WalletRegistry proxy admin", async () => {
    expect(
      ethers.utils.defaultAbiCoder.decode(
        ["address"],
        await ethers.provider.getStorageAt(
          walletRegistry.address,
          "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103"
        )
      )[0],
      "invalid WalletRegistry proxy admin (read from storage slot)"
    ).to.be.equal(proxyAdmin.address)

    expect(
      await walletRegistryProxy.connect(proxyAdmin.address).callStatic.admin(),
      "invalid WalletRegistry proxy admin"
    ).to.be.equal(proxyAdmin.address)
  })

  it("should set ProxyAdmin owner", async () => {
    expect(await proxyAdmin.owner(), "invalid ProxyAdmin owner").to.be.equal(
      esdm.address
    )
  })

  it("should set WalletRegistry governance", async () => {
    expect(
      await walletRegistry.governance(),
      "invalid WalletRegistry governance"
    ).equal(walletRegistryGovernance.address)
  })

  it("should set WalletRegistryGovernance owner", async () => {
    expect(
      await walletRegistryGovernance.owner(),
      "invalid WalletRegistryGovernance owner"
    ).equal(governance.address)
  })
})
