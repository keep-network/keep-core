/* eslint-disable @typescript-eslint/no-unused-expressions */
import { deployments, ethers, helpers, upgrades } from "hardhat"
import chai, { expect } from "chai"
import chaiAsPromised from "chai-as-promised"

import type { Contract } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryGovernance } from "../typechain"
import type { TransparentUpgradeableProxy } from "../typechain/TransparentUpgradeableProxy"

chai.use(chaiAsPromised)

const { AddressZero } = ethers.constants

describe("WalletRegistry - Deployment", async () => {
  let deployer: SignerWithAddress
  let governance: SignerWithAddress
  let esdm: SignerWithAddress

  let walletRegistry: WalletRegistry
  let walletRegistryGovernance: WalletRegistryGovernance
  let walletRegistryProxy: TransparentUpgradeableProxy
  let proxyAdmin: Contract
  let walletRegistryImplementationAddress: string

  before(async () => {
    await deployments.fixture()
    ;({ deployer, governance, esdm } = await helpers.signers.getNamedSigners())

    walletRegistry = await helpers.contracts.getContract<WalletRegistry>(
      "WalletRegistry"
    )

    walletRegistryImplementationAddress = (
      await deployments.get("WalletRegistry")
    ).implementation

    walletRegistryGovernance =
      await helpers.contracts.getContract<WalletRegistryGovernance>(
        "WalletRegistryGovernance"
      )

    walletRegistryProxy = await ethers.getContractAt(
      "TransparentUpgradeableProxy",
      walletRegistry.address
    )

    proxyAdmin = await upgrades.admin.getInstance()

    expect(deployer.address, "deployer is the same as governance").not.equal(
      governance.address
    )
  })

  it("should set WalletRegistry proxy admin", async () => {
    // To let a non-proxy-admin read the admin we have to read it directly from
    // the storage slot, see: https://docs.openzeppelin.com/contracts/4.x/api/proxy#TransparentUpgradeableProxy-admin--
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

  it("should set WalletRegistry implementation", async () => {
    // To let a non-proxy-admin read the implementation we have to read it directly from
    // the storage slot, see: https://docs.openzeppelin.com/contracts/4.x/api/proxy#TransparentUpgradeableProxy-implementation--
    expect(
      ethers.utils.defaultAbiCoder.decode(
        ["address"],
        await ethers.provider.getStorageAt(
          walletRegistry.address,
          "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"
        )
      )[0],
      "invalid WalletRegistry implementation (read from storage slot)"
    ).to.be.equal(walletRegistryImplementationAddress)

    expect(
      await walletRegistryProxy
        .connect(proxyAdmin.address)
        .callStatic.implementation(),
      "invalid WalletRegistry implementation"
    ).to.be.equal(walletRegistryImplementationAddress)
  })

  it("should set WalletRegistry implementation in ProxyAdmin", async () => {
    expect(
      await proxyAdmin.getProxyImplementation(walletRegistryProxy.address),
      "invalid proxy implementation"
    ).to.be.equal(walletRegistryImplementationAddress)
  })

  it("should set implementation address different than proxy address", async () => {
    expect(
      await walletRegistryProxy
        .connect(proxyAdmin.address)
        .callStatic.implementation(),
      "invalid ProxyAdmin owner"
    ).to.be.not.equal(walletRegistry.address)
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

  it("should set WalletRegistry address in artifact to the proxy address", async () => {
    expect(walletRegistry.address, "invalid WalletRegistry address").equal(
      walletRegistryProxy.address
    )
  })

  it("should revert when initialize called again", async () => {
    await expect(
      walletRegistry.initialize(AddressZero, AddressZero, AddressZero)
    ).to.be.revertedWith("Initializable: contract is already initialized")
  })
})
