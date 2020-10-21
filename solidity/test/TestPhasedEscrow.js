const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {expectRevert, expectEvent} = require("@openzeppelin/test-helpers")

const KeepToken = contract.fromArtifact("KeepToken")
const PhasedEscrow = contract.fromArtifact("PhasedEscrow")
const CurveRewardsEscrowBeneficiary = contract.fromArtifact(
  "CurveRewardsEscrowBeneficiary"
)

const TestSimpleBeneficiary = contract.fromArtifact("TestSimpleBeneficiary")
const TestCurveRewards = contract.fromArtifact("TestCurveRewards")

const chai = require("chai")
chai.use(require("bn-chai")(web3.utils.BN))
const expect = chai.expect

describe("PhasedEscrow", () => {
  const owner = accounts[0]
  const updatedOwner = accounts[1]

  let beneficiary
  let updatedBeneficiary

  let token
  let phasedEscrow

  before(async () => {
    token = await KeepToken.new({from: owner})
    phasedEscrow = await PhasedEscrow.new(token.address, {from: owner})
    beneficiary = await TestSimpleBeneficiary.new()
    updatedBeneficiary = await TestSimpleBeneficiary.new()
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("setBeneficiary", async () => {
    it("can be called by owner", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      // ok, no revert
    })

    it("can be called by updated owner", async () => {
      await phasedEscrow.transferOwnership(updatedOwner, {from: owner})

      await expectRevert(
        phasedEscrow.setBeneficiary(beneficiary.address, {from: owner}),
        "Ownable: caller is not the owner"
      )
      await phasedEscrow.setBeneficiary(beneficiary.address, {
        from: updatedOwner,
      })
      // ok, no revert
    })

    it("can not be called by non-owner", async () => {
      await expectRevert(
        phasedEscrow.setBeneficiary(beneficiary.address, {from: updatedOwner}),
        "Ownable: caller is not the owner"
      )
    })

    it("sets beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})

      expect(await phasedEscrow.beneficiary()).to.equal(
        beneficiary.address,
        "Unexpected beneficiary"
      )
    })

    it("emits an event", async () => {
      const receipt = await phasedEscrow.setBeneficiary(beneficiary.address, {
        from: owner,
      })

      expectEvent(receipt, "BeneficiaryUpdated", {
        beneficiary: beneficiary.address,
      })
    })
  })

  describe("withdraw", async () => {
    it("can not be called if beneficiary wasn't set", async () => {
      await token.transfer(phasedEscrow.address, 100, {from: owner})
      await expectRevert(
        phasedEscrow.withdraw(100, {from: owner}),
        "Beneficiary not assigned"
      )
    })

    it("can not be called by non-owner", async () => {
      await token.transfer(phasedEscrow.address, 100, {from: owner})
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      await expectRevert(
        phasedEscrow.withdraw(100, {from: beneficiary.address}),
        "Ownable: caller is not the owner"
      )
    })

    it("can be called by owner", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      await token.transfer(phasedEscrow.address, 100, {from: owner})
      await phasedEscrow.withdraw(100, {from: owner})
      // ok, no reverts
    })

    it("fails when escrow is empty", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      await expectRevert(
        phasedEscrow.withdraw(100, {from: owner}),
        "Not enough tokens for withdrawal"
      )
    })

    it("withdraws specified tokens to updated beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      const amount = web3.utils.toBN(987654321)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      await phasedEscrow.withdraw(100, {from: owner})

      await phasedEscrow.setBeneficiary(updatedBeneficiary.address, {
        from: owner,
      })
      await phasedEscrow.withdraw(100, {from: owner})

      expect(await token.balanceOf(beneficiary.address)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(updatedBeneficiary.address)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
        987654321 - 200,
        "Unexpected amount withdrawn"
      )
    })

    it("withdraws specified tokens to beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      const amount = web3.utils.toBN(123456789)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      await phasedEscrow.withdraw(100, {from: owner})

      expect(await token.balanceOf(beneficiary.address)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
        123456789 - 100,
        "Unexpected amount withdrawn"
      )
    })

    it("emits an event", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      const amount = web3.utils.toBN(100)
      await token.transfer(phasedEscrow.address, amount.muln(2), {from: owner})

      const receipt = await phasedEscrow.withdraw(amount, {from: owner})

      await expectEvent(receipt, "TokensWithdrawn", {
        beneficiary: beneficiary.address,
        amount: amount,
      })
    })

    describe("when withdrawing to a CurveRewardsEscrowBeneficiary", () => {
      let curveRewards
      let rewardsBeneficiary

      const baseBalance = 123456789
      const transferAmount = 100

      before(async () => {
        curveRewards = await TestCurveRewards.new(token.address)
        rewardsBeneficiary = await CurveRewardsEscrowBeneficiary.new(
          token.address,
          curveRewards.address
        )

        const amount = web3.utils.toBN(baseBalance)
        await token.transfer(phasedEscrow.address, amount, {from: owner})

        await phasedEscrow.setBeneficiary(rewardsBeneficiary.address, {
          from: owner,
        })
      })

      beforeEach(createSnapshot)
      afterEach(restoreSnapshot)

      it("withdraws specified tokens from escrow", async () => {
        await phasedEscrow.withdraw(transferAmount, {from: owner})

        expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
          baseBalance - transferAmount,
          "Unexpected amount withdrawn"
        )
      })

      it("transfers specified tokens to curve rewards contract", async () => {
        await phasedEscrow.withdraw(transferAmount, {from: owner})

        expect(await token.balanceOf(curveRewards.address)).to.eq.BN(
          transferAmount,
          "Unexpected amount deposited"
        )
      })

      it("leaves no tokens in the rewards beneficiary", async () => {
        await phasedEscrow.withdraw(100, {from: owner})

        expect(await token.balanceOf(rewardsBeneficiary.address)).to.eq.BN(
          0,
          "Unexpected amount left in rewards beneficiary"
        )
      })

      it("emits a TokensWithdrawn event to the rewards beneficiary", async () => {
        const receipt = await phasedEscrow.withdraw(transferAmount, {
          from: owner,
        })

        expectEvent(receipt, "TokensWithdrawn", {
          beneficiary: rewardsBeneficiary.address,
          amount: web3.utils.toBN(transferAmount),
        })
      })

      it("emits a RewardAdded event from the rewards beneficiary", async () => {
        const receipt = resolveAllLogs(
          (await phasedEscrow.withdraw(transferAmount, {from: owner})).receipt,
          {curveRewards}
        )

        expectEvent(receipt, "RewardAdded", {
          reward: web3.utils.toBN(transferAmount),
        })
      })
    })
  })
})

// FIXME Move to a shared test utils library for all Keep projects.
/**
 * Uses the ABIs of all contracts in the `contractContainer` to resolve any
 * events they may have emitted into the given `receipt`'s logs. Typically
 * Truffle only resolves the events on the calling contract; this function
 * resolves all of the ones that can be resolved.
 *
 * @param {TruffleReceipt} receipt The receipt of a contract function call
 *        submission.
 * @param {ContractContainer} contractContainer An object that contains
 *        properties that are TruffleContracts. Not all properties in the
 *        container need be contracts, nor do all contracts need to have events
 *        in the receipt.
 *
 * @return {TruffleReceipt} The receipt, with its `logs` property updated to
 *         include all resolved logs.
 */
function resolveAllLogs(receipt, contractContainer) {
  const contracts = Object.entries(contractContainer)
    .map(([, value]) => value)
    .filter((_) => _.contract && _.address)

  const {resolved: resolvedLogs} = contracts.reduce(
    ({raw, resolved}, contract) => {
      const events = contract.contract._jsonInterface.filter(
        (_) => _.type === "event"
      )
      const contractLogs = raw.filter((_) => _.address == contract.address)

      const decoded = contractLogs.map((log) => {
        const event = events.find((_) => log.topics.includes(_.signature))
        const decoded = web3.eth.abi.decodeLog(
          event.inputs,
          log.data,
          log.topics.slice(1)
        )

        return Object.assign({}, log, {
          event: event.name,
          args: decoded,
        })
      })

      return {
        raw: raw.filter((_) => _.address != contract.address),
        resolved: resolved.concat(decoded),
      }
    },
    {raw: receipt.rawLogs, resolved: []}
  )

  return Object.assign({}, receipt, {
    logs: resolvedLogs,
  })
}
