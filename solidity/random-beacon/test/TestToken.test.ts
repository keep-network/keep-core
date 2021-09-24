import { ethers } from "hardhat";
import { expect } from "chai";

describe("Token", function () {
  let testToken

  beforeEach(async function () {
    const TestToken = await ethers.getContractFactory("TestToken")
    testToken = await TestToken.deploy()
    await testToken.deployed()
  });

  it("should test token name", async function () {
    expect(await testToken.name()).to.equal("Test Token")
  });
});