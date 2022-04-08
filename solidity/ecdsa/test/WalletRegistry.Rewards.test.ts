import { helpers } from "hardhat"
import { expect } from "chai"

import { walletRegistryFixture } from "./fixtures"
import ecdsaData from "./data/ecdsa"
import { createNewWallet } from "./utils/wallets"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { FakeContract } from "@defi-wonderland/smock"
import type { Operator, OperatorID } from "./utils/operators"
import type {
  SortitionPool,
  WalletRegistry,
  TokenStaking,
  IWalletOwner,
  T,
} from "../typechain"

const { to1e18 } = helpers.number

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("WalletRegistry - Rewards", () => {
  let tToken: T
  let walletRegistry: WalletRegistry
  let staking: TokenStaking
  let sortitionPool: SortitionPool
  let walletOwner: FakeContract<IWalletOwner>

  let deployer: SignerWithAddress
  let thirdParty: SignerWithAddress

  let members: Operator[]

  const walletPublicKey: string = ecdsaData.group1.publicKey

  const rewardAmount = to1e18(100000)

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({
      tToken,
      walletRegistry,
      staking,
      sortitionPool,
      walletOwner,
      deployer,
      thirdParty,
    } = await walletRegistryFixture())
    ;({ members } = await createNewWallet(
      walletRegistry,
      walletOwner.wallet,
      walletPublicKey
    ))
  })

  describe("withdrawRewards", () => {
    context("when called for an unknown operator", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.withdrawRewards(thirdParty.address)
        ).to.be.revertedWith("Unknown operator")
      })
    })

    context("when called for a known operator", () => {
      before(async () => {
        await createSnapshot()

        // Allocate sortition pool rewards
        await tToken.connect(deployer).mint(deployer.address, rewardAmount)
        await tToken
          .connect(deployer)
          .approveAndCall(sortitionPool.address, rewardAmount, [])
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should withdraw rewards", async () => {
        const operator = members[0].signer.address
        const stakingProvider = await walletRegistry.operatorToStakingProvider(
          operator
        )
        const { beneficiary } = await staking.rolesOf(stakingProvider)

        expect(await tToken.balanceOf(beneficiary)).to.equal(0)
        await walletRegistry.withdrawRewards(stakingProvider)
        expect(await tToken.balanceOf(beneficiary)).to.be.gt(0)
      })
    })
  })
})
