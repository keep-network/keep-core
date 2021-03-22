const { toBN } = require("web3-utils")
const BN = require("bn.js")

const chai = require("chai")
chai.use(require("bn-chai")(BN))
const expect = chai.expect

/**
 * Checks if a value is close to the expected within an allowed deviation.
 * In case of a value that is outside the delta range it will fail chai expect.
 *
 * @param {BN|number|string} actual Actual value.
 * @param {BN|number|string} expected Expected value.
 * @param {string} [message] Message in case of failure (optional).
 * @param {number} [deltaPercent=1] Percent deviation from the expected value that is
 * acceptable. E.g. value of 5 means that the actual value cannot be different
 * from the expected value more than 5%. Default value: 1.
 */
function expectCloseTo(actual, expected, message, deltaPercent = 1) {
  actualBN = toBN(actual)
  expectedBN = toBN(expected)

  const delta = expectedBN.muln(deltaPercent).divn(100) // approx. `deltaPercent` %

  if (
    actualBN.lt(expectedBN.sub(delta)) ||
    actualBN.gt(expectedBN.add(delta))
  ) {
    expect.fail(
      `${message}\nexpected : ${expectedBN.toString()}\nactual   : ${actualBN.toString()}`
    )
  }
}

module.exports.expectCloseTo = expectCloseTo
