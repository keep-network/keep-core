import { ethers } from "hardhat"
import { expect } from "chai"

// This is an example of Typescript test.
// Should be removed once we have actual tests.
describe("TestToken", () => {
  let testToken

  beforeEach(async () => {
    const TestToken = await ethers.getContractFactory("TestToken")
    testToken = await TestToken.deploy()
    await testToken.deployed()
  })

  it("should test token name", async () => {
    expect(await testToken.name()).to.equal("Test Token")
  })
})
