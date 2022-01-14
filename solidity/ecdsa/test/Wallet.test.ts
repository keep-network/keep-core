import {
  deployments,
  ethers,
  waffle,
  helpers,
  getUnnamedAccounts,
} from "hardhat"
import { expect } from "chai"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { ContractFactory, ContractTransaction } from "ethers"
import { keccak256 } from "ethers/lib/utils"
import type { Wallet, CloneFactoryStub } from "../typechain"
import ecdsaData from "./data/ecdsa"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

const fixture = async () => {
  await deployments.fixture(["MasterWallet"])

  const masterWallet: Wallet = await ethers.getContract("MasterWallet")

  const CloneFactory: ContractFactory = await ethers.getContractFactory(
    "CloneFactoryStub"
  )
  const cloneFactory: CloneFactoryStub = (await CloneFactory.deploy(
    masterWallet.address
  )) as CloneFactoryStub

  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")

  const owner = await ethers.getSigner((await getUnnamedAccounts())[0])
  const thirdParty = await ethers.getSigner((await getUnnamedAccounts())[1])

  return {
    masterWallet,
    cloneFactory,
    deployer,
    owner,
    thirdParty,
  }
}

describe("Wallet", () => {
  const membersIdsHash = hashUint32Array([101, 102, 103])
  const publicKeyHash = keccak256(ecdsaData.publicKey)
  const { digest1: digest } = ecdsaData

  let masterWallet: Wallet
  let cloneFactory: CloneFactoryStub
  let deployer: SignerWithAddress
  let owner: SignerWithAddress
  let thirdParty: SignerWithAddress

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ masterWallet, cloneFactory, deployer, owner, thirdParty } =
      await waffle.loadFixture(fixture))
  })

  context("for master contract", async () => {
    let wallet: Wallet

    before(async () => {
      await createSnapshot()
      wallet = masterWallet
    })

    after(async () => {
      await restoreSnapshot()
    })

    describe("init", async () => {
      context("called by a deployer", async () => {
        it("should revert", async () => {
          await expect(
            wallet
              .connect(deployer)
              .init(owner.address, membersIdsHash, publicKeyHash)
          ).to.be.revertedWith("Initialization of master wallet is not allowed")
        })
      })

      context("called by an owner", async () => {
        it("should revert", async () => {
          await expect(
            wallet
              .connect(owner)
              .init(owner.address, membersIdsHash, publicKeyHash)
          ).to.be.revertedWith("Initialization of master wallet is not allowed")
        })
      })

      context("called by a third party", async () => {
        it("should revert", async () => {
          await expect(
            wallet
              .connect(thirdParty)
              .init(owner.address, membersIdsHash, publicKeyHash)
          ).to.be.revertedWith("Initialization of master wallet is not allowed")
        })
      })
    })

    describe("activate", async () => {
      context("called by a third party", async () => {
        it("should revert", async () => {
          await expect(
            wallet.connect(thirdParty).activate()
          ).to.be.revertedWith("Ownable: caller is not the owner")
        })
      })
    })

    describe("sign", async () => {
      context("called by a third party", async () => {
        it("should revert", async () => {
          await expect(
            wallet.connect(thirdParty).sign(keccak256("0x0102"))
          ).to.be.revertedWith("Ownable: caller is not the owner")
        })
      })
    })
  })

  context("for cloned contract", async () => {
    let wallet: Wallet

    before("create clone instance", async () => {
      await createSnapshot()

      const walletAddress: string = await cloneFactory
        .connect(deployer)
        .callStatic.createClone()
      await cloneFactory.connect(deployer).createClone()

      wallet = await ethers.getContractAt("Wallet", walletAddress)
    })

    after(async () => {
      await restoreSnapshot()
    })

    describe("init", async () => {
      beforeEach(async () => {
        await createSnapshot()
      })

      afterEach(async () => {
        await restoreSnapshot()
      })

      context("called by a deployer", async () => {
        it("should succeed", async () => {
          await expect(
            wallet
              .connect(deployer)
              .init(owner.address, membersIdsHash, publicKeyHash)
          ).to.not.be.reverted
        })
      })

      context("called by an owner", async () => {
        it("should succeed", async () => {
          await expect(
            wallet
              .connect(owner)
              .init(owner.address, membersIdsHash, publicKeyHash)
          ).to.not.be.reverted
        })
      })

      context("called by a third party", async () => {
        it("should succeed", async () => {
          await expect(
            wallet
              .connect(thirdParty)
              .init(owner.address, membersIdsHash, publicKeyHash)
          ).to.not.be.reverted
        })

        it("should store data", async () => {
          await wallet
            .connect(thirdParty)
            .init(owner.address, membersIdsHash, publicKeyHash)

          await expect(await wallet.membersIdsHash()).to.be.equal(
            membersIdsHash
          )
          await expect(await wallet.publicKeyHash()).to.be.equal(publicKeyHash)
        })

        it("should transfer ownership", async () => {
          await expect(
            await wallet.owner(),
            "invalid initial owner"
          ).to.be.not.equal(owner.address)

          await wallet
            .connect(thirdParty)
            .init(owner.address, membersIdsHash, publicKeyHash)

          await expect(await wallet.owner(), "invalid new owner").to.be.equal(
            owner.address
          )
        })
      })
    })

    describe("activate", async () => {
      context("for initialized contract", async () => {
        before("initialize", async () => {
          await createSnapshot()

          await wallet
            .connect(deployer)
            .init(owner.address, membersIdsHash, publicKeyHash)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("called by a third party", async () => {
          it("should revert", async () => {
            await expect(
              wallet.connect(thirdParty).activate()
            ).to.be.revertedWith("Ownable: caller is not the owner")
          })
        })

        context("called by the owner", async () => {
          it("should set activation block number", async () => {
            await createSnapshot()

            const tx: ContractTransaction = await wallet
              .connect(owner)
              .activate()

            await expect(await wallet.activationBlockNumber()).to.be.equal(
              tx.blockNumber
            )

            await restoreSnapshot()
          })
        })

        context("for activated contract", async () => {
          before("activate", async () => {
            await createSnapshot()

            await wallet.connect(owner).activate()
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("called by the owner", async () => {
            it("should revert", async () => {
              await expect(wallet.connect(owner).activate()).to.be.revertedWith(
                "Wallet was already activated"
              )
            })
          })
        })
      })
    })

    describe("sign", async () => {
      context("for initialized contract", async () => {
        before("initialize", async () => {
          await createSnapshot()

          await wallet
            .connect(deployer)
            .init(owner.address, membersIdsHash, publicKeyHash)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("called by a third party", async () => {
          it("should revert", async () => {
            await expect(
              wallet.connect(thirdParty).sign(digest)
            ).to.be.revertedWith("Ownable: caller is not the owner")
          })
        })

        context("called by the owner", async () => {
          it("should emit SignatureRequested event", async () => {
            const tx: ContractTransaction = await wallet
              .connect(owner)
              .sign(digest)

            await expect(tx)
              .to.emit(wallet, "SignatureRequested")
              .withArgs(digest)
          })
        })
      })
    })
  })
})
