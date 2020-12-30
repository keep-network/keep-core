const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const KeepToken = contract.fromArtifact("KeepToken")
const assert = require("chai").assert

describe("TestToken", function () {
  let token
  const accountOne = accounts[0]
  const accountTwo = accounts[1]

  before(async () => {
    token = await KeepToken.new({from: accountOne})
  })

  it("sets token details", async function () {
    await token.name.call()

    assert.equal(await token.name.call(), "KEEP Token", "unexpected token name")
    assert.equal(await token.symbol.call(), "KEEP", "unexpected token symbol")
    assert.equal(await token.decimals.call(), 18, "unexpected decimals")
  })

  it("should send tokens correctly", async function () {
    const amount = web3.utils.toBN(1000000000)

    // Starting balances
    const accountOneStartingBalance = await token.balanceOf.call(accountOne)
    const accountTwoStartingBalance = await token.balanceOf.call(accountTwo)

    // Send tokens
    await token.transfer(accountTwo, amount, {from: accountOne})

    // Ending balances
    const accountOneEndingBalance = await token.balanceOf.call(accountOne)
    const accountTwoEndingBalance = await token.balanceOf.call(accountTwo)

    assert.equal(
      accountOneEndingBalance.eq(accountOneStartingBalance.sub(amount)),
      true,
      "Amount wasn't correctly taken from the sender"
    )
    assert.equal(
      accountTwoEndingBalance.eq(accountTwoStartingBalance.add(amount)),
      true,
      "Amount wasn't correctly sent to the receiver"
    )
  })
})
