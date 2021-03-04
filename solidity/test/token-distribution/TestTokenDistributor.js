const { expectRevert } = require("@openzeppelin/test-helpers")
const {
  accounts,
  privateKeys,
  contract,
  web3,
} = require("@openzeppelin/test-environment")
const TokenDistributor = contract.fromArtifact("TokenDistributor")

describe("TokenDistributor", () => {
  let tokenDistributor

  before(async () => {
    tokenDistributor = await TokenDistributor.new()
  })

  describe("claim", async () => {
    const thirdPartyIndex = 1
    const recipientIndex = 2
    const destinationIndex = 3

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

      const recipientAddress = accounts[recipientIndex]
      const destinationAddress = accounts[destinationIndex]

      for (const [testName, testData] of testCases) {
        console.log(`      ${testName}`)
        try {
          const signature = web3.eth.accounts.sign(
            web3.utils.sha3(destinationAddress),
            privateKeys[testData.signer]
          )

          claimFuncCall = tokenDistributor.claim(
            recipientAddress,
            destinationAddress,
            signature.v,
            signature.r,
            signature.s,
            { from: accounts[testData.submitter] }
          )

          if (testData.expectRevert) {
            await expectRevert(claimFuncCall, "invalid signature")
          } else {
            await claimFuncCall
          }
        } catch (err) {
          throw new Error(`Test case [${testName}] failed with error: ${err}`)
        }
      }
    })

    it("detects malleable signatures", async function () {
      const secp256k1N = web3.utils.toBN(
        "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141"
      )

      const recipientAddress = accounts[recipientIndex]
      const destinationAddress = accounts[destinationIndex]

      const signature = web3.eth.accounts.sign(
        web3.utils.sha3(destinationAddress),
        privateKeys[recipientIndex]
      )

      const malleableS =
        "0x" + secp256k1N.sub(web3.utils.toBN(signature.s)).toJSON()

      await expectRevert(
        tokenDistributor.claim(
          recipientAddress,
          destinationAddress,
          signature.v,
          signature.r,
          malleableS
        ),
        "Malleable signature - s should be in the low half of secp256k1 curve's order"
      )
    })
  })
})
