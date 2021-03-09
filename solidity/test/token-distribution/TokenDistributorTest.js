const { expectRevert, expectEvent } = require("@openzeppelin/test-helpers")
const {
  accounts,
  privateKeys,
  contract,
  web3,
} = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { parseBalanceMap } = require("@keep-network/merkle-distributor-helper")

const TokenDistributor = contract.fromArtifact("TokenDistributor")
const TestToken = contract.fromArtifact("TestToken")

const { BN, toBN } = web3.utils
const chai = require("chai")
chai.use(require("bn-chai")(BN))
const { expect } = chai

describe("TokenDistributor", () => {
  const ownerIndex = 1
  const recipientIndex = 2
  const destinationIndex = 3
  const thirdPartyIndex = 4

  const owner = accounts[ownerIndex]
  const recipient = accounts[recipientIndex]
  const destination = accounts[destinationIndex]
  const thirdParty = accounts[thirdPartyIndex]

  let testToken
  let tokenDistributor

  const allocationsMap = {}
  allocationsMap[recipient] = (100000).toString(16)
  allocationsMap[thirdParty] = (200000).toString(16)

  const testData = {
    signature: web3.eth.accounts.sign(
      web3.utils.sha3(destination),
      privateKeys[recipientIndex]
    ),
    merkle: parseBalanceMap(allocationsMap),
  }

  before(async () => {
    testToken = await TestToken.new({ from: owner })
    tokenDistributor = await TokenDistributor.new(testToken.address, {
      from: owner,
    })

    await testToken.mint(owner, testData.merkle.tokenTotal)
    await testToken.approve(
      tokenDistributor.address,
      testData.merkle.tokenTotal,
      {
        from: owner,
      }
    )
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("allocate", async () => {
    it("transfers tokens", async function () {
      await tokenDistributor.allocate(
        testData.merkle.merkleRoot,
        testData.merkle.tokenTotal,
        { from: owner }
      )

      expect(await testToken.balanceOf(tokenDistributor.address)).to.eq.BN(
        toBN(testData.merkle.tokenTotal)
      )
    })

    it("sets merkle root", async function () {
      await tokenDistributor.allocate(
        testData.merkle.merkleRoot,
        testData.merkle.tokenTotal,
        { from: owner }
      )

      expect(await tokenDistributor.merkleRoot()).to.equal(
        testData.merkle.merkleRoot
      )
    })

    it("emits event", async function () {
      expectEvent(
        await tokenDistributor.allocate(
          testData.merkle.merkleRoot,
          testData.merkle.tokenTotal,
          { from: owner }
        ),
        "TokensAllocated",
        {
          merkleRoot: testData.merkle.merkleRoot,
          amount: toBN(testData.merkle.tokenTotal),
        }
      )
    })

    it("reverts on merkle root overwrite", async function () {
      await tokenDistributor.allocate(
        testData.merkle.merkleRoot,
        testData.merkle.tokenTotal,
        { from: owner }
      )

      await expectRevert(
        tokenDistributor.allocate("0x1234567890", testData.merkle.tokenTotal, {
          from: owner,
        }),
        "tokens were already allocated"
      )
    })

    it("reverts on empty merkle root", async function () {
      await expectRevert(
        tokenDistributor.allocate([], testData.merkle.tokenTotal, {
          from: owner,
        }),
        "merkle root cannot be empty"
      )
    })

    it("reverts on token transfer failure", async function () {
      await expectRevert(
        tokenDistributor.allocate(
          testData.merkle.merkleRoot,
          testData.merkle.tokenTotal + 1,
          {
            from: owner,
          }
        ),
        "SafeERC20: low-level call failed"
      )
    })
  })

  describe("claim", async () => {
    beforeEach(async () => {
      await tokenDistributor.allocate(
        testData.merkle.merkleRoot,
        testData.merkle.tokenTotal,
        { from: owner }
      )
    })

    it("destination address signature verification", async function () {
      const testCases = new Map([
        [
          "completes when signed by recipient, submitted by third-party",
          {
            signer: recipientIndex,
            submitter: thirdPartyIndex,
            expectRevert: false,
          },
        ],
        [
          "completes when signed by recipient, submitted by recipient",
          {
            signer: recipientIndex,
            submitter: recipientIndex,
            expectRevert: false,
          },
        ],
        [
          "reverts when signed by third-party, submitted by recipient",
          {
            signer: thirdPartyIndex,
            submitter: recipientIndex,
            expectRevert: true,
          },
        ],
        [
          "completes when signed by recipient, submitted by destination",
          {
            signer: recipientIndex,
            submitter: destinationIndex,
            expectRevert: false,
          },
        ],
        [
          "reverts when signed by destination, submitted by recipient",
          {
            signer: destinationIndex,
            submitter: recipientIndex,
            expectRevert: true,
          },
        ],
        [
          "reverts when signed by destination, submitted by destination",
          {
            signer: destinationIndex,
            submitter: destinationIndex,
            expectRevert: true,
          },
        ],
      ])

      for (const [testCaseName, testCaseData] of testCases) {
        await createSnapshot()

        console.log(`      ${testCaseName}`)
        try {
          const signature = web3.eth.accounts.sign(
            web3.utils.sha3(destination),
            privateKeys[testCaseData.signer]
          )

          claimFuncCall = tokenDistributor.claim(
            recipient,
            destination,
            signature.v,
            signature.r,
            signature.s,
            testData.merkle.claims[recipient].index,
            testData.merkle.claims[recipient].amount,
            testData.merkle.claims[recipient].proof,
            { from: accounts[testCaseData.submitter] }
          )

          if (testCaseData.expectRevert) {
            await expectRevert(claimFuncCall, "invalid signature")
          } else {
            await claimFuncCall
          }
        } catch (err) {
          throw new Error(
            `Test case [${testCaseName}] failed with error: ${err}`
          )
        } finally {
          await restoreSnapshot()
        }
      }
    })

    it("detects malleable signatures", async function () {
      const secp256k1N = toBN(
        "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141"
      )

      const malleableS =
        "0x" + secp256k1N.sub(toBN(testData.signature.s)).toJSON()

      await expectRevert(
        tokenDistributor.claim(
          recipient,
          destination,
          testData.signature.v,
          testData.signature.r,
          malleableS,
          testData.merkle.claims[recipient].index,
          testData.merkle.claims[recipient].amount,
          testData.merkle.claims[recipient].proof
        ),
        "Invalid signature 's' value"
      )

      await expectRevert(
        tokenDistributor.claim(
          recipient,
          destination,
          testData.signature.v - 27,
          testData.signature.r,
          testData.signature.s,
          testData.merkle.claims[recipient].index,
          testData.merkle.claims[recipient].amount,
          testData.merkle.claims[recipient].proof
        ),
        "Invalid signature 'v' value"
      )
    })
    //  TODO: Add more tests
  })
})
