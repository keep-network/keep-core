const {contract, web3} = require("@openzeppelin/test-environment")
var assert = require('chai').assert
const OperatorParamsStub = contract.fromArtifact('OperatorParamsStub');

describe('OperatorParams', () => {
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

  describe("pack", async () => {
    it("should roundtrip values", async () => {
      const params = await opUtils.publicPack(
        allKeepEver,
        recently,
        billionYearsFromNow);
      const amount = await opUtils.publicGetAmount(params);
      const createdAt = await opUtils.publicGetCreationTimestamp(params);
      const undelegatedAt = await opUtils.publicGetUndelegationTimestamp(params);

      assert.equal(
        amount.toJSON(),
        allKeepEver.toJSON(),
        "The amount should be the same");
      assert.equal(
        createdAt.toJSON(),
        recently.toJSON(),
        "The creation timestamp should be the same");
      assert.equal(
        undelegatedAt.toJSON(),
        billionYearsFromNow.toJSON(),
        "The undelegation timestamp should be the same");
    })
  })

  describe("setAmount", async () => {
    it("should set the amount", async () => {
      const params = await opUtils.publicPack(allKeepEver, recently, 0);
      const newParams = await opUtils.publicSetAmount(params, billion);
      const amount = await opUtils.publicGetAmount(newParams);
      assert.equal(
        amount.toJSON(),
        billion.toJSON(),
        "The amount should be the same");
    })
  })

  describe("setCreationTimestamp", async () => {
    it("should set the creation timestamp", async () => {
      const params = await opUtils.publicPack(allKeepEver, recently, 0);
      const newParams = await opUtils.publicSetCreationTimestamp(
        params,
        billionYearsFromNow);
      const creationBlock = await opUtils.publicGetCreationTimestamp(newParams);
      assert.equal(
        creationBlock.toJSON(),
        billionYearsFromNow.toJSON(),
        "The creation timestamp should be the same");
    })
  })

  describe("setUndelegationTimestamp", async () => {
    it("should set the undelegation timestamp", async () => {
      const params = await opUtils.publicPack(allKeepEver, recently, 0);
      const newParams = await opUtils.publicSetUndelegationTimestamp(
        params,
        recently);
      const undelegationTimestamp = await opUtils.publicGetUndelegationTimestamp(newParams);
      assert.equal(
        undelegationTimestamp.toJSON(),
        recently.toJSON(),
        "The undelegationTimestamp should be the same");
    })
  })

  describe("setAmountAndCreationTimestamp", async () => {
    it("should set the creation timestamp", async () => {
      const params = await opUtils.publicPack(allKeepEver, recently, 0)
      const newParams = await opUtils.publicSetAmountAndCreationTimestamp(
        params,
        billion,
        billionYearsFromNow        
      )
      const creationBlock = await opUtils.publicGetCreationTimestamp(newParams)
      assert.equal(
        creationBlock.toJSON(),
        billionYearsFromNow.toJSON(),
        "The creation timestamp should be the same"
      )
    })

    it("should set the amount", async () => {
      const params = await opUtils.publicPack(allKeepEver, recently, 0);
      const newParams = await opUtils.publicSetAmountAndCreationTimestamp(
        params,
        billion,
        billionYearsFromNow
      )
      const amount = await opUtils.publicGetAmount(newParams);
      assert.equal(
        amount.toJSON(),
        billion.toJSON(),
        "The amount should be the same"
      )
    })
  })
})
