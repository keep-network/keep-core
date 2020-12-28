const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const KeepToken = contract.fromArtifact("KeepToken")
const assert = require("chai").assert

describe("TestToken", function () {
  let token
  const account_one = accounts[0]
  const account_two = accounts[1]

  before(async () => {
    token = await KeepToken.new({from: account_one})
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
    const account_one_starting_balance = await token.balanceOf.call(account_one)
    const account_two_starting_balance = await token.balanceOf.call(account_two)

    // Send tokens
    await token.transfer(account_two, amount, {from: account_one})

    // Ending balances
    const account_one_ending_balance = await token.balanceOf.call(account_one)
    const account_two_ending_balance = await token.balanceOf.call(account_two)

    assert.equal(
      account_one_ending_balance.eq(account_one_starting_balance.sub(amount)),
      true,
      "Amount wasn't correctly taken from the sender"
    )
    assert.equal(
      account_two_ending_balance.eq(account_two_starting_balance.add(amount)),
      true,
      "Amount wasn't correctly sent to the receiver"
    )
  })
})
