const { expectRevert } = require("@openzeppelin/test-helpers")
const { contract } = require("@openzeppelin/test-environment")
const TestAltBn128 = contract.fromArtifact("TestAltBn128")

describe("AltBn128", () => {

    let g1 = "0x15c30f4b6cf6dbbcbdcc10fe22f54c8170aea44e198139b776d512d8f027319a1b9e8bfaf1383978231ce98e42bafc8129f473fc993cf60ce327f7d223460663"
    let g2 = "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d01659dc18b57722ecf6a4beb4d04dfe780a660c4c3bb2b165ab8486114c464c621bf37ecdba226629c20908c7f475c5b3a7628ce26d696436eab0b0148034dfcd"
    let g2Compressed = "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0"

    let altBn128

    before(async () => {
        altBn128 = await TestAltBn128.new()
    })

    describe("g1Unmarshal", async () => {
        it("does not accept less than 64 bytes", async () => {
            await expectRevert(
                altBn128.publicG1Unmarshal(g1.slice(0, -2)),
                "Invalid G1 bytes length"
            )
        })

        it("does accept 64 bytes", async () => {
            await altBn128.publicG1Unmarshal(g1)
            // ok, no revert
        })

        it("does not accept more than 64 bytes", async () => {
            await expectRevert(
                altBn128.publicG1Unmarshal(g1 + 'ff'),
                "Invalid G1 bytes length"
            )
        })
    })

    describe("g2Unmarshal", async () => {
        it("does not accept less than 128 bytes", async () => {
            await expectRevert(
                altBn128.publicG2Unmarshal(g2.slice(0, -2)),
                "Invalid G2 bytes length"
            )
        })

        it("does accept 128 bytes", async () => {
            await altBn128.publicG2Unmarshal(g2)
            // ok, no revert
        })

        it("does not accept more than 128 bytes", async () => {
            await expectRevert(
                altBn128.publicG2Unmarshal(g2 + 'ff'),
                "Invalid G2 bytes length"
            )
        })
    })

    describe("g2Decompress", async () => {
        it("does not accept less than 64 bytes", async () => {
            await expectRevert(
                altBn128.publicG2Decompress(g2Compressed.slice(0, -2)),
                "Invalid G2 compressed bytes length"
            )
        })

        it("does accept 64 bytes", async () => {
            await altBn128.publicG2Decompress(g2Compressed)
            // ok, no revert
        })

        it("does not accept more than 64 bytes", async () => {
            await expectRevert(
                altBn128.publicG2Decompress(g2Compressed + 'ff'),
                "Invalid G2 compressed bytes length"
            )
        })
    })

  it("runHashingTest()", async () => {
    await altBn128.runHashingTest()
    // ok, no revert
  })

  it("runHashAndAddTest()", async () => {
    await altBn128.runHashAndAddTest()
    // ok, no revert
  })

  it("runHashAndScalarMultiplyTest()", async () => {
    await altBn128.runHashAndScalarMultiplyTest()
    // ok, no revert
  })

  it("runGfP2AddTest()", async () => {
    await altBn128.runGfP2AddTest()
    // ok, no revert
  })

  it("runAddTest()", async () => {
    await altBn128.runAddTest()
    // ok, no revert
  })

  it("runScalarMultiplyTest()", async () => {
    await altBn128.runScalarMultiplyTest()
    // ok, no revert
  })

  it("runBasicPairingTest()", async () => {
    await altBn128.runBasicPairingTest()
    // ok, no revert
  })

  it("runVerifySignatureTest()", async () => {
    await altBn128.runVerifySignatureTest()
    // ok, no revert
  })

  it("runCompressG1InvertibilityTest()", async () => {
    await altBn128.runCompressG1InvertibilityTest()
    // ok, no revert
  })

  it("runCompressG2InvertibilityTest()", async () => {
    await altBn128.runCompressG2InvertibilityTest()
    // ok, no revert
  })

  it("runG2PointOnCurveTest()", async () => {
    await altBn128.runG2PointOnCurveTest()
    // ok, no revert
  })
})
