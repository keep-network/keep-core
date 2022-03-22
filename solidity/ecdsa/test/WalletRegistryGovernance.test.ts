import { deployments, ethers, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"

import { constants, params, updateWalletRegistryParams } from "./fixtures"

import type { ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type {
  WalletRegistry,
  WalletRegistryStub,
  WalletRegistryGovernance,
} from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot
const { to1e18 } = helpers.number

const fixture = deployments.createFixture(async () => {
  await deployments.fixture(["WalletRegistry"])

  const walletRegistry: WalletRegistryStub & WalletRegistry =
    await ethers.getContract("WalletRegistry")
  const walletRegistryGovernance: WalletRegistryGovernance =
    await ethers.getContract("WalletRegistryGovernance")

  const governance: SignerWithAddress = await ethers.getNamedSigner(
    "governance"
  )

  const thirdParty: SignerWithAddress = await ethers.getSigner(
    (
      await getUnnamedAccounts()
    )[0]
  )

  await updateWalletRegistryParams(walletRegistryGovernance, governance)

  return {
    walletRegistry,
    walletRegistryGovernance,
    governance,
    thirdParty,
  }
})

describe("WalletRegistryGovernance", async () => {
  let governance: SignerWithAddress
  let walletRegistry: WalletRegistry
  let walletRegistryGovernance: WalletRegistryGovernance
  let thirdParty: SignerWithAddress

  const initialMinimumAuthorization = to1e18(400000)
  const initialAuthorizationDecreaseDelay = 5184000 // 60 days
  const initialMaliciousDkgResultSlashingAmount = to1e18(50000)
  const initialMaliciousDkgResultNotificationRewardMultiplier = 100
  const initialDkgResultSubmissionGas = 300000
  const initialDkgApprovalGasOffset = 65000
  const initialSortitionPoolRewardsBanDuration = 1209600 // 14 days

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletRegistryGovernance, governance, thirdParty } =
      await fixture())
  })

  describe("upgradeRandomBeacon", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .upgradeRandomBeacon(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      context("when new address is zero", () => {
        it("should revert when a new random beacon address is zero", async () => {
          await expect(
            walletRegistryGovernance
              .connect(governance)
              .upgradeRandomBeacon(ethers.constants.AddressZero)
          ).to.be.revertedWith("New random beacon address cannot be zero")
        })
      })

      context("when new address is not zero", () => {
        before(async () => {
          await createSnapshot()

          tx = await walletRegistryGovernance
            .connect(governance)
            .upgradeRandomBeacon(thirdParty.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the random beacon", async () => {
          expect(await walletRegistry.randomBeacon()).to.be.equal(
            thirdParty.address
          )
        })

        it("should emit RandomBeaconUpgraded event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "RandomBeaconUpgraded")
            .withArgs(thirdParty.address)
        })
      })
    })
  })

  describe("initializeWalletOwner", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .initializeWalletOwner(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      context("when new address is zero", () => {
        it("should revert when a new address is zero", async () => {
          await expect(
            walletRegistryGovernance
              .connect(governance)
              .initializeWalletOwner(ethers.constants.AddressZero)
          ).to.be.revertedWith("Wallet Owner address cannot be zero")
        })
      })

      context("when new address is not zero", () => {
        before(async () => {
          await createSnapshot()

          tx = await walletRegistryGovernance
            .connect(governance)
            .initializeWalletOwner(thirdParty.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the wallet owner", async () => {
          expect(await walletRegistry.walletOwner()).to.be.equal(
            thirdParty.address
          )
        })

        it("should emit WalletOwnerUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "WalletOwnerUpdated")
            .withArgs(thirdParty.address)
        })

        it("should revert when called for the second time", async () => {
          await expect(
            walletRegistryGovernance
              .connect(governance)
              .initializeWalletOwner(thirdParty.address)
          ).to.be.revertedWith("Wallet Owner already initialized")
        })
      })
    })
  })

  describe("beginGovernanceDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginGovernanceDelayUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginGovernanceDelayUpdate(1337)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the governance delay", async () => {
        expect(await walletRegistryGovernance.governanceDelay()).to.be.equal(
          constants.governanceDelay
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingGovernanceDelayUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit GovernanceDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(walletRegistryGovernance, "GovernanceDelayUpdateStarted")
          .withArgs(1337, blockTimestamp)
      })
    })
  })

  describe("finalizeGovernanceDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeGovernanceDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeGovernanceDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      before(async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginGovernanceDelayUpdate(7331)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeGovernanceDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginGovernanceDelayUpdate(7331)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeGovernanceDelayUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the governance delay", async () => {
          expect(await walletRegistryGovernance.governanceDelay()).to.be.equal(
            7331
          )
        })

        it("should emit GovernanceDelayUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "GovernanceDelayUpdated")
            .withArgs(7331)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingGovernanceDelayUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginWalletRegistryOwnershipTransfer", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginWalletRegistryOwnershipTransfer(
              "0x00Ea7D21bcCEeD400aCe08B583554aA619D3e537"
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginWalletRegistryOwnershipTransfer(
            "0x00Ea7D21bcCEeD400aCe08B583554aA619D3e537"
          )
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when new owner is the zero address", () => {
        it("should revert", async () => {
          await expect(
            walletRegistryGovernance
              .connect(governance)
              .beginWalletRegistryOwnershipTransfer(
                ethers.constants.AddressZero
              )
          ).to.be.revertedWith(
            "New wallet registry owner address cannot be zero"
          )
        })
      })

      it("should not transfer the ownership", async () => {
        expect(await walletRegistry.owner()).to.be.equal(
          walletRegistryGovernance.address
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingWalletRegistryOwnershipTransferDelayTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit WalletRegistryOwnershipTransferStarted", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "WalletRegistryOwnershipTransferStarted"
          )
          .withArgs(
            "0x00Ea7D21bcCEeD400aCe08B583554aA619D3e537",
            blockTimestamp
          )
      })
    })
  })

  describe("finalizeWalletRegistryOwnershipTransfer", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeWalletRegistryOwnershipTransfer()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeWalletRegistryOwnershipTransfer()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      before(async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginWalletRegistryOwnershipTransfer(
            "0x00Ea7D21bcCEeD400aCe08B583554aA619D3e537"
          )

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeWalletRegistryOwnershipTransfer()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginWalletRegistryOwnershipTransfer(
              "0x00Ea7D21bcCEeD400aCe08B583554aA619D3e537"
            )

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeWalletRegistryOwnershipTransfer()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should transfer wallet registry ownership", async () => {
          expect(await walletRegistry.owner()).to.be.equal(
            "0x00Ea7D21bcCEeD400aCe08B583554aA619D3e537"
          )
        })

        it("should emit WalletRegistryOwnershipTransferred event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "WalletRegistryOwnershipTransferred"
            )
            .withArgs("0x00Ea7D21bcCEeD400aCe08B583554aA619D3e537")
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingWalletRegistryOwnershipTransferDelayTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginWalletOwnerUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginWalletOwnerUpdate(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginWalletOwnerUpdate(thirdParty.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      describe("when new wallet owner is zero address", async () => {
        it("should revert", async () => {
          await expect(
            walletRegistryGovernance
              .connect(governance)
              .beginWalletOwnerUpdate(ethers.constants.AddressZero)
          ).to.be.revertedWith("New wallet owner address cannot be zero")
        })
      })

      it("should not update the wallet owner", async () => {
        expect(await walletRegistry.walletOwner()).to.be.equal(
          ethers.constants.AddressZero
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingWalletOwnerUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the WalletOwnerUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(walletRegistryGovernance, "WalletOwnerUpdateStarted")
          .withArgs(thirdParty.address, blockTimestamp)
      })
    })
  })

  describe("finalizeWalletOwnerUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeWalletOwnerUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeWalletOwnerUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginWalletOwnerUpdate(thirdParty.address)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeWalletOwnerUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginWalletOwnerUpdate(thirdParty.address)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeWalletOwnerUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the wallet owner", async () => {
          expect(await walletRegistry.walletOwner()).to.be.equal(
            thirdParty.address
          )
        })

        it("should emit WalletOwnerUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "WalletOwnerUpdated")
            .withArgs(thirdParty.address)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingWalletOwnerUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginMinimumAuthorizationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginMinimumAuthorizationUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginMinimumAuthorizationUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the minimum authorization amount", async () => {
        expect(await walletRegistry.minimumAuthorization()).to.be.equal(
          initialMinimumAuthorization
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingMimimumAuthorizationUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the MinimumAuthorizationUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "MinimumAuthorizationUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeMinimumAuthorizationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeMinimumAuthorizationUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMinimumAuthorizationUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginMinimumAuthorizationUpdate(123)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMinimumAuthorizationUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginMinimumAuthorizationUpdate(123)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeMinimumAuthorizationUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the minimum authorization amount", async () => {
          expect(await walletRegistry.minimumAuthorization()).to.be.equal(123)
        })

        it("should emit MinimumAuthorizationUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "MinimumAuthorizationUpdated")
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingMimimumAuthorizationUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginAuthorizationDecreaseDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginAuthorizationDecreaseDelayUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginAuthorizationDecreaseDelayUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the authorization decrease delay", async () => {
        expect(await walletRegistry.authorizationDecreaseDelay()).to.be.equal(
          initialAuthorizationDecreaseDelay
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingAuthorizationDecreaseDelayUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the AuthorizationDecreaseDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "AuthorizationDecreaseDelayUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeAuthorizationDecreaseDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeAuthorizationDecreaseDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeAuthorizationDecreaseDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginAuthorizationDecreaseDelayUpdate(123)
        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeAuthorizationDecreaseDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginAuthorizationDecreaseDelayUpdate(123)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeAuthorizationDecreaseDelayUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the authorization decrease delay", async () => {
          expect(await walletRegistry.authorizationDecreaseDelay()).to.be.equal(
            123
          )
        })

        it("should emit AuthorizationDecreaseDelayUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "AuthorizationDecreaseDelayUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingAuthorizationDecreaseDelayUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginMaliciousDkgResultSlashingAmountUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the malicious DKG result slashing amount", async () => {
        expect(
          await walletRegistry.maliciousDkgResultSlashingAmount()
        ).to.be.equal(initialMaliciousDkgResultSlashingAmount)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the MaliciousDkgResultSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "MaliciousDkgResultSlashingAmountUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginMaliciousDkgResultSlashingAmountUpdate(123)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the malicious DKG result slashing amount", async () => {
          expect(
            await walletRegistry.maliciousDkgResultSlashingAmount()
          ).to.be.equal(123)
        })

        it("should emit MaliciousDkgResultSlashingAmountUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "MaliciousDkgResultSlashingAmountUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgResultSubmissionGasUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgResultSubmissionGasUpdate(100)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner and value is correct", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionGasUpdate(100)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG submission result gas", async () => {
        expect(await walletRegistry.dkgResultSubmissionGas()).to.be.equal(
          initialDkgResultSubmissionGas
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgResultSubmissionGasUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the DkgResultSubmissionGasUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgResultSubmissionGasUpdateStarted"
          )
          .withArgs(100, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultSubmissionGasUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgResultSubmissionGasUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionGasUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionGasUpdate(100)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionGasUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionGasUpdate(100)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionGasUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result submission gas", async () => {
          expect(await walletRegistry.dkgResultSubmissionGas()).to.be.equal(100)
        })

        it("should emit DkgResultSubmissionGasUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "DkgResultSubmissionGasUpdated")
            .withArgs(100)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgResultSubmissionGasUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgApprovalGasOffsetUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgApprovalGasOffsetUpdate(100)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner and value is correct", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgApprovalGasOffsetUpdate(100)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG approval gas offset", async () => {
        expect(await walletRegistry.dkgApprovalGasOffset()).to.be.equal(
          initialDkgApprovalGasOffset
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgApprovalGasOffsetUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the DkgApprovalGasOffsetUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgApprovalGasOffsetUpdateStarted"
          )
          .withArgs(100, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgApprovalGasOffsetUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgApprovalGasOffsetUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgApprovalGasOffsetUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgApprovalGasOffsetUpdate(100)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgApprovalGasOffsetUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgApprovalGasOffsetUpdate(100)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgApprovalGasOffsetUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result approval gas offset", async () => {
          expect(await walletRegistry.dkgApprovalGasOffset()).to.be.equal(100)
        })

        it("should emit DkgApprovalGasOffsetUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "DkgApprovalGasOffsetUpdated")
            .withArgs(100)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgApprovalGasOffsetUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginMaliciousDkgResultNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(100)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called with value >100", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(101)
        ).to.be.revertedWith("Maximum value is 100")
      })
    })

    context("when the caller is the owner and value is correct", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(100)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG malicious result notification reward multiplier", async () => {
        expect(
          await walletRegistry.maliciousDkgResultNotificationRewardMultiplier()
        ).to.be.equal(initialMaliciousDkgResultNotificationRewardMultiplier)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingMaliciousDkgResultNotificationRewardMultiplierUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the MaliciousDkgResultNotificationRewardMultiplierUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "MaliciousDkgResultNotificationRewardMultiplierUpdateStarted"
          )
          .withArgs(100, blockTimestamp)
      })
    })
  })

  describe("finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(100)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(100)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG malicious result notification reward multiplier", async () => {
          expect(
            await walletRegistry.maliciousDkgResultNotificationRewardMultiplier()
          ).to.be.equal(100)
        })

        it("should emit MaliciousDkgResultNotificationRewardMultiplierUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "MaliciousDkgResultNotificationRewardMultiplierUpdated"
            )
            .withArgs(100)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingMaliciousDkgResultNotificationRewardMultiplierUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginSortitionPoolRewardsBanDurationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginSortitionPoolRewardsBanDurationUpdate(86400)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginSortitionPoolRewardsBanDurationUpdate(86400)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the sortition pool rewards ban duration", async () => {
        expect(
          await walletRegistry.sortitionPoolRewardsBanDuration()
        ).to.be.equal(initialSortitionPoolRewardsBanDuration)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingSortitionPoolRewardsBanDurationUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the SortitionPoolRewardsBanDurationUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "SortitionPoolRewardsBanDurationUpdateStarted"
          )
          .withArgs(86400, blockTimestamp)
      })
    })
  })

  describe("finalizeSortitionPoolRewardsBanDurationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeSortitionPoolRewardsBanDurationUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeSortitionPoolRewardsBanDurationUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginSortitionPoolRewardsBanDurationUpdate(86400)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeSortitionPoolRewardsBanDurationUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginSortitionPoolRewardsBanDurationUpdate(86400)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeSortitionPoolRewardsBanDurationUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the sortition pool rewards ban duration", async () => {
          expect(
            await walletRegistry.sortitionPoolRewardsBanDuration()
          ).to.be.equal(86400)
        })

        it("should emit SortitionPoolRewardsBanDurationUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "SortitionPoolRewardsBanDurationUpdated"
            )
            .withArgs(86400)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingSortitionPoolRewardsBanDurationUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgSeedTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgSeedTimeoutUpdate(11)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is equal 0", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSeedTimeoutUpdate(0)
        ).to.be.revertedWith("DKG seed timeout must be > 0")
      })
    })

    context("when the update value is at least 1", () => {
      before(async () => {
        await createSnapshot()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should accept the value", async () => {
        await createSnapshot()

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSeedTimeoutUpdate(1)
        ).not.to.be.reverted
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSeedTimeoutUpdate(11)
        ).not.to.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgSeedTimeoutUpdate(11)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG seed timeout", async () => {
        expect((await walletRegistry.dkgParameters()).seedTimeout).to.be.equal(
          params.dkgSeedTimeout
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgSeedTimeoutUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the DkgSeedTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(walletRegistryGovernance, "DkgSeedTimeoutUpdateStarted")
          .withArgs(11, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgSeedTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgSeedTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSeedTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgSeedTimeoutUpdate(11)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSeedTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgSeedTimeoutUpdate(11)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSeedTimeoutUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG seed timeout", async () => {
          expect(
            (await walletRegistry.dkgParameters()).seedTimeout
          ).to.be.equal(11)
        })

        it("should emit DkgSeedTimeoutUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "DkgSeedTimeoutUpdated")
            .withArgs(11)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgSeedTimeoutUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgResultChallengePeriodLengthUpdate(11)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is less than 10", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(9)
        ).to.be.revertedWith("DKG result challenge period length must be >= 10")
      })
    })

    context("when the update value is at least 10", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(10)
        ).to.not.be.reverted

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(11)
        ).to.not.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG result challenge period length", async () => {
        expect(
          (await walletRegistry.dkgParameters()).resultChallengePeriodLength
        ).to.be.equal(params.dkgResultChallengePeriodLength)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgResultChallengePeriodLengthUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the DkgResultChallengePeriodLengthUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgResultChallengePeriodLengthUpdateStarted"
          )
          .withArgs(11, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(11)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result challenge period length", async () => {
          expect(
            (await walletRegistry.dkgParameters()).resultChallengePeriodLength
          ).to.be.equal(11)
        })

        it("should emit DkgResultChallengePeriodLengthUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "DkgResultChallengePeriodLengthUpdated"
            )
            .withArgs(11)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgResultChallengePeriodLengthUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgResultSubmissionTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgResultSubmissionTimeoutUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(0)
        ).to.be.revertedWith(
          "DKG result submission eligibility delay must be > 0"
        )
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(1)
        ).to.not.be.reverted
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(2)
        ).to.not.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG result submission eligibility delay", async () => {
        expect(
          (await walletRegistry.dkgParameters()).resultSubmissionTimeout
        ).to.be.equal(params.dkgResultSubmissionTimeout)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgResultSubmissionTimeoutUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the DkgResultSubmissionTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgResultSubmissionTimeoutUpdateStarted"
          )
          .withArgs(1, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultSubmissionTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(1)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(10)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result submission eligibility delay", async () => {
          expect(
            (await walletRegistry.dkgParameters()).resultSubmissionTimeout
          ).to.be.equal(10)
        })

        it("should emit DkgResultSubmissionTimeoutUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "DkgResultSubmissionTimeoutUpdated"
            )
            .withArgs(10)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgResultSubmissionTimeoutUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgSubmitterPrecedencePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(0)
        ).to.be.revertedWith(
          "DKG submitter precedence period length must be > 0"
        )
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
        ).to.not.be.reverted
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(2)
        ).to.not.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG submitter precedence period length", async () => {
        expect(
          (await walletRegistry.dkgParameters()).submitterPrecedencePeriodLength
        ).to.be.equal(params.dkgSubmitterPrecedencePeriodLength)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingSubmitterPrecedencePeriodLengthUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the DkgSubmitterPrecedencePeriodLengthUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgSubmitterPrecedencePeriodLengthUpdateStarted"
          )
          .withArgs(1, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgSubmitterPrecedencePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(10)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG submitter precedence period length", async () => {
          expect(
            (await walletRegistry.dkgParameters())
              .submitterPrecedencePeriodLength
          ).to.be.equal(10)
        })

        it("should emit DkgSubmitterPrecedencePeriodLengthUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "DkgSubmitterPrecedencePeriodLengthUpdated"
            )
            .withArgs(10)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingSubmitterPrecedencePeriodLengthUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginReimbursementPoolUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginReimbursementPoolUpdate(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginReimbursementPoolUpdate(thirdParty.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the wallet owner", async () => {
        expect(await walletRegistry.walletOwner()).to.be.equal(
          ethers.constants.AddressZero
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingReimbursementPoolUpdateTime()
        ).to.be.equal(constants.governanceDelay)
      })

      it("should emit the ReimbursementPoolUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(walletRegistryGovernance, "ReimbursementPoolUpdateStarted")
          .withArgs(thirdParty.address, blockTimestamp)
      })
    })
  })

  describe("finalizeReimbursementPoolUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeReimbursementPoolUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeReimbursementPoolUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginReimbursementPoolUpdate(thirdParty.address)

        await helpers.time.increaseTime(constants.governanceDelay - 60) // -1min

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeReimbursementPoolUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginReimbursementPoolUpdate(thirdParty.address)

          await helpers.time.increaseTime(constants.governanceDelay)

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeReimbursementPoolUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the reimbursement pool", async () => {
          expect(await walletRegistry.reimbursementPool()).to.be.equal(
            thirdParty.address
          )
        })

        it("should emit ReimbursementPoolUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "ReimbursementPoolUpdated")
            .withArgs(thirdParty.address)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingReimbursementPoolUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })
})
