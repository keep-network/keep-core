/* eslint-disable @typescript-eslint/no-unused-expressions, no-await-in-loop */

import { BigNumber } from "ethers"
import { ethers, helpers, getUnnamedAccounts, waffle } from "hardhat"
import { expect } from "chai"

import blsData from "./data/bls"

import { constants } from "./fixtures"
import { selectGroup } from "./utils/groups"
import { signDkgResult, DkgResult, noMisbehaved } from "./utils/dkg"
import { Operator } from "./utils/operators"

import type { SortitionPool, DKGValidator } from "../typechain"

const { to1e18 } = helpers.number

const fixture = async () => {
  const SortitionPool = await ethers.getContractFactory("SortitionPool")
  const sortitionPool = (await SortitionPool.deploy(
    constants.poolWeightDivisor
  )) as SortitionPool

  const DKGValidator = await ethers.getContractFactory("DKGValidator")
  const dkgValidator = (await DKGValidator.deploy(
    sortitionPool.address
  )) as DKGValidator
  await dkgValidator.deployed()

  return {
    sortitionPool,
    dkgValidator,
  }
}

describe("DKGValidator", () => {
  const dkgSeed: BigNumber = BigNumber.from(
    "31415926535897932384626433832795028841971693993751058209749445923078164062862"
  )
  const dkgStartBlock = 1337
  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let selectedOperators

  let prepareDkgResult
  let validator: DKGValidator

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(fixture)
    const { sortitionPool } = contracts
    validator = contracts.dkgValidator

    const operators = (await getUnnamedAccounts()).slice(0, constants.groupSize)
    for (let i = 0; i < operators.length; i++) {
      await sortitionPool.insertOperator(operators[i], to1e18(100))
    }

    selectedOperators = await selectGroup(sortitionPool, dkgSeed)

    prepareDkgResult = async (
      _groupMembers: Operator[],
      _signers: Operator[],
      _groupPublicKey: string,
      _misbehaved: number[],
      _startBlock: number,
      _numberOfSignatures = 33,
      _submitterIndex = 1
    ) => {
      const { signingMembersIndices, signaturesBytes } = await signDkgResult(
        _signers,
        _groupPublicKey,
        _misbehaved,
        _startBlock,
        _numberOfSignatures
      )

      const dkgResult: DkgResult = {
        submitterMemberIndex: _submitterIndex,
        groupPubKey: _groupPublicKey,
        misbehavedMembersIndices: _misbehaved,
        signatures: signaturesBytes,
        signingMembersIndices,
        members: _groupMembers.map((m) => m.id),
      }

      return dkgResult
    }
  })

  describe("validate", () => {
    const testValidate = async (_groupMembers, _signers, _groupPublicKey) => {
      const dkgResult = await prepareDkgResult(
        _groupMembers,
        _signers,
        _groupPublicKey,
        noMisbehaved,
        dkgStartBlock
      )

      const result = await validator.validate(dkgResult, dkgSeed, dkgStartBlock)

      return {
        isValid: result[0],
        errorMsg: result[1],
      }
    }

    context("for a valid DKG result", () => {
      it("should pass", async () => {
        const result = await testValidate(
          selectedOperators,
          selectedOperators,
          groupPublicKey
        )

        expect(result.isValid).to.be.true
        expect(result.errorMsg).to.equal("")
      })
    })

    // just a basic test ensuring group members validation
    // is called; detailed test cases are implemented for
    // validateGroupMembers function
    context("for a DKG result with invalid group members", () => {
      it("should return validation error", async () => {
        const shuffledOperators = shuffle(selectedOperators)

        const result = await testValidate(
          shuffledOperators,
          shuffledOperators,
          groupPublicKey
        )

        expect(result.isValid).to.be.false
        expect(result.errorMsg).to.equal("Invalid group members")
      })
    })

    // just a basic test ensuring signatures validation
    // is called; detailed test cases are implemented for
    // validateSignatures function
    context("for a DKG result with invalid signatures", () => {
      it("should return validation error", async () => {
        const shuffledOperators = shuffle(selectedOperators)

        const result = await testValidate(
          selectedOperators,
          shuffledOperators,
          groupPublicKey
        )

        expect(result.isValid).to.be.false
        expect(result.errorMsg).to.equal("Invalid signatures")
      })
    })
  })

  describe("validateFields", () => {
    const testValidateFields = async (
      _groupMembers,
      _groupPublicKey,
      _misbehaved,
      _numberOfSignatures = 33
    ) => {
      const dkgResult = await prepareDkgResult(
        _groupMembers,
        _groupMembers,
        _groupPublicKey,
        _misbehaved,
        dkgStartBlock,
        _numberOfSignatures
      )

      const result = await validator.validateFields(dkgResult)

      return {
        isValid: result[0],
        errorMsg: result[1],
      }
    }

    context("for a valid DKG result", () => {
      it("should pass", async () => {
        const result = await testValidateFields(
          selectedOperators,
          groupPublicKey,
          noMisbehaved
        )

        expect(result.isValid).to.be.true
        expect(result.errorMsg).to.equal("")
      })
    })

    context("for malformed group public key", () => {
      it("should return validation error", async () => {
        const empty = "0x"
        const tooShort = groupPublicKey.substring(0, groupPublicKey.length - 2)
        const tooLong = `${groupPublicKey}ff`

        let result = await testValidateFields(
          selectedOperators,
          empty,
          noMisbehaved
        )

        expect(result.isValid).to.be.false
        expect(result.errorMsg).to.equal("Malformed group public key")

        result = await testValidateFields(
          selectedOperators,
          tooShort,
          noMisbehaved
        )

        expect(result.isValid).to.be.false
        expect(result.errorMsg).to.equal("Malformed group public key")

        result = await testValidateFields(
          selectedOperators,
          tooLong,
          noMisbehaved
        )

        expect(result.isValid).to.be.false
        expect(result.errorMsg).to.equal("Malformed group public key")
      })
    })

    context("for malformed misbehaved array", () => {
      it("should return validation error", async () => {
        const lessThanOne = [0, 1, 2]
        const higherThanGroupSize = [1, 2, 65]
        const unsorted = [1, 2, 3, 60, 4]

        let result = await testValidateFields(
          selectedOperators,
          groupPublicKey,
          lessThanOne
        )

        expect(result.isValid).to.be.false
        expect(result.errorMsg).to.equal("Corrupted misbehaved members indices")

        result = await testValidateFields(
          selectedOperators,
          groupPublicKey,
          higherThanGroupSize
        )

        expect(result.isValid).to.be.false
        expect(result.errorMsg).to.equal("Corrupted misbehaved members indices")

        result = await testValidateFields(
          selectedOperators,
          groupPublicKey,
          unsorted
        )

        expect(result.isValid).to.be.false
        expect(result.errorMsg).to.equal("Corrupted misbehaved members indices")
      })
    })

    // TODO: expand tests to ensure all possible cases of corrupted input
    //       data are covered;
  })

  describe("validateGroupMembers", () => {
    const testValidateGroupMembers = async (_groupMembers) => {
      const dkgResult = await prepareDkgResult(
        _groupMembers,
        _groupMembers,
        groupPublicKey,
        noMisbehaved,
        dkgStartBlock
      )

      return validator.validateGroupMembers(dkgResult, dkgSeed)
    }

    context("for a valid DKG result", () => {
      it("should pass", async () => {
        const isValid = await testValidateGroupMembers(selectedOperators)
        expect(isValid).to.be.true
      })
    })

    context("for operators other than selected", () => {
      it("should fail the validation", async () => {
        const isValid = await testValidateGroupMembers(
          shuffle(selectedOperators)
        )
        expect(isValid).to.be.false
      })
    })
  })

  describe("validateSignatures", () => {
    const testValidateSignatures = async (_groupMembers, _signers) => {
      const dkgResult = await prepareDkgResult(
        _groupMembers,
        _signers,
        groupPublicKey,
        noMisbehaved,
        dkgStartBlock
      )

      return validator.validateSignatures(dkgResult, dkgStartBlock)
    }

    context("for a valid DKG result", () => {
      it("should pass", async () => {
        const isValid = await testValidateSignatures(
          selectedOperators,
          selectedOperators
        )

        expect(isValid).to.be.true
      })
    })

    context(
      "for signature provided by one malicious, selected operator",
      () => {
        it("should fail the validation", async () => {
          const maliciousSigners = Array(constants.groupSize).fill(
            selectedOperators[0]
          )
          const isValid = await testValidateSignatures(
            selectedOperators,
            maliciousSigners
          )

          expect(isValid).to.be.false
        })
      }
    )

    // TODO: expand tests to ensure all possible cases of corrupted input
    //       data are covered;
  })
})

function shuffle(operators: Operator[]): Operator[] {
  return operators
    .map((v) => ({ v, sort: Math.random() }))
    .sort((a, b) => a.sort - b.sort)
    .map(({ v }) => v)
}
