const OperatorParamsStub = artifacts.require('./stubs/OperatorParamsStub.sol');

contract('OperatorParamsStub', (accounts) => {
  let opUtils;
  const eighteen = web3.utils.toBN(18)
  const ten = web3.utils.toBN(10)
  const keepDecimals = ten.pow(eighteen);
  const billion = web3.utils.toBN(1000000000)
  const allKeepEver = billion.mul(keepDecimals);

  const blocksPerYear = web3.utils.toBN(3153600);
  const recently = blocksPerYear.muln(5);
  const billionYearsFromNow = blocksPerYear.mul(billion);

  before(async () => {
      opUtils = await OperatorParamsStub.new();
  });

  it("should roundtrip values", async () => {
    const params = await opUtils.publicPack(allKeepEver, recently, billionYearsFromNow);
    const amount = await opUtils.publicGetAmount(params);
    const createdAt = await opUtils.publicGetCreationBlock(params);
    const undelegatedAt = await opUtils.publicGetUndelegationBlock(params);

    assert.equal(amount.toJSON(), allKeepEver.toJSON(), "The amount should be the same");
    assert.equal(createdAt.toJSON(), recently.toJSON(), "The creation block should be the same");
    assert.equal(undelegatedAt.toJSON(), billionYearsFromNow.toJSON(), "The undelegation block should be the same");
  })
})
