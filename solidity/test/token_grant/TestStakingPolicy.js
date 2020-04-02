const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const PermissiveStakingPolicy = artifacts.require('./PermissiveStakingPolicy.sol');
const GuaranteedMinimumStakingPolicy = artifacts.require('./GuaranteedMinimumStakingPolicy.sol');

contract('PermissiveStakingPolicy', function(accounts) {
  let policy;
  let amount = 10000;
  let start = 1000;
  let duration = 2000;
  let cliff = 1500;
  let withdrawn = 4000;

  before(async () => {
    policy = await PermissiveStakingPolicy.new();
  });

  async function calculate(atTimestamp, withdrawnAmount) {
    return await policy.getStakeableAmount(
      atTimestamp,
      amount,
      duration,
      start,
      cliff,
      withdrawnAmount
    );
  }

  it("should permit staking all tokens before cliff", async function() {
    expect(await calculate(1400, 0)).to.eq.BN(amount);
  });

  it("should permit staking all tokens after cliff", async function() {
    expect(await calculate(1600, 0)).to.eq.BN(amount);
  });

  it("should permit staking all tokens after unlocking", async function() {
    expect(await calculate(3100, 0)).to.eq.BN(amount);
  });

  it("should permit staking remaining tokens before cliff", async function() {
    expect(await calculate(1400, withdrawn)).to.eq.BN(amount - withdrawn);
  });

  it("should permit staking remaining tokens after cliff", async function() {
    expect(await calculate(1600, withdrawn)).to.eq.BN(amount - withdrawn);
  });

  it("should permit staking remaining tokens after unlocking", async function() {
    expect(await calculate(3100, withdrawn)).to.eq.BN(amount - withdrawn);
  });
});

contract('GuaranteedMinimumStakingPolicy', function(accounts) {
  let policy;
  let minimumStake = 2000;
  let largeGrant = 10000;
  let mediumGrant = 5000;
  let smallGrant = 1000;
  let start = 1000;
  let duration = 2000;
  let cliff = 1500;

  before(async () => {
    policy = await GuaranteedMinimumStakingPolicy.new(minimumStake);
  });

  async function calculate(atTimestamp, givenAmount, withdrawnAmount) {
    return await policy.getStakeableAmount(
      atTimestamp,
      givenAmount,
      duration,
      start,
      cliff,
      withdrawnAmount
    );
  }

  describe("with nothing withdrawn", async () => {
    it("should calculate stakeable amount correctly before cliff", async () => {
      expect(await calculate(1499, largeGrant, 0)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with large grant before cliff");
      expect(await calculate(1499, mediumGrant, 0)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with medium grant before cliff");
      expect(await calculate(1499, smallGrant, 0)).to.eq.BN(
        smallGrant,
        "Should permit entire grant with small grant before cliff");
    });

    it("should calculate stakeable amount correctly just after cliff", async () => {
      expect(await calculate(1500, largeGrant, 0)).to.eq.BN(
        2500,
        "Should permit unlocked amount with large grant just after cliff");
      expect(await calculate(1500, mediumGrant, 0)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with medium grant just after cliff");
      expect(await calculate(1500, smallGrant, 0)).to.eq.BN(
        smallGrant,
        "Should permit entire grant with small grant just after cliff");
    });

    it("should calculate stakeable amount correctly halfway through", async () => {
      expect(await calculate(2000, largeGrant, 0)).to.eq.BN(
        5000,
        "Should permit unlocked amount with large grant halfway through");
      expect(await calculate(2000, mediumGrant, 0)).to.eq.BN(
        2500,
        "Should permit unlocked amount with medium grant halfway through");
      expect(await calculate(2000, smallGrant, 0)).to.eq.BN(
        smallGrant,
        "Should permit entire grant with small grant halfway through");
    });

    it("should calculate stakeable amount correctly after unlocking period", async () => {
      expect(await calculate(3000, largeGrant, 0)).to.eq.BN(
        largeGrant,
        "Should permit unlocked amount with large grant after unlocking period");
      expect(await calculate(3000, mediumGrant, 0)).to.eq.BN(
        mediumGrant,
        "Should permit unlocked amount with medium grant after unlocking period");
      expect(await calculate(3000, smallGrant, 0)).to.eq.BN(
        smallGrant,
        "Should permit entire grant with small grant after unlocking period");
    });
  })

  describe("with all unlocked tokens withdrawn", async () => {
    it("should calculate stakeable amount correctly just after cliff", async () => {
      expect(await calculate(1500, largeGrant, 2500)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with large grant just after cliff");
      expect(await calculate(1500, mediumGrant, 1250)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with medium grant just after cliff");
      expect(await calculate(1500, smallGrant, 250)).to.eq.BN(
        750,
        "Should permit remaining amount with small grant just after cliff");
    });

    it("should calculate stakeable amount correctly halfway through", async () => {
      expect(await calculate(2000, largeGrant, 5000)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with large grant halfway through");
      expect(await calculate(2000, mediumGrant, 2500)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with medium grant halfway through");
      expect(await calculate(2000, smallGrant, 500)).to.eq.BN(
        500,
        "Should permit remaining amount with small grant halfway through");
    });

    it("should calculate stakeable amount correctly at three quarters", async () => {
      expect(await calculate(2500, largeGrant, 7500)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with large grant at three quarters");
      expect(await calculate(2500, mediumGrant, 3750)).to.eq.BN(
        1250,
        "Should permit remaining amount with medium grant at three quarters");
      expect(await calculate(2500, smallGrant, 750)).to.eq.BN(
        250,
        "Should permit remaining amount with small grant at three quarters");
    });
  })

  describe("with half of unlocked tokens withdrawn", async () => {
    it("should calculate stakeable amount correctly just after cliff", async () => {
      expect(await calculate(1500, largeGrant, 1250)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with large grant just after cliff");
      expect(await calculate(1500, mediumGrant, 625)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with medium grant just after cliff");
      expect(await calculate(1500, smallGrant, 125)).to.eq.BN(
        875,
        "Should permit remaining amount with small grant just after cliff");
    });

    it("should calculate stakeable amount correctly halfway through", async () => {
      expect(await calculate(2000, largeGrant, 2500)).to.eq.BN(
        2500,
        "Should permit remaining unlocked amount with large grant halfway through");
      expect(await calculate(2000, mediumGrant, 1250)).to.eq.BN(
        minimumStake,
        "Should permit minimum stake with medium grant halfway through");
      expect(await calculate(2000, smallGrant, 250)).to.eq.BN(
        750,
        "Should permit remaining amount with small grant halfway through");
    });
  })
});
