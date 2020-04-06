const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const grantTokens = require('../helpers/grantTokens');
const { createSnapshot, restoreSnapshot } = require('../helpers/snapshot');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const KeepToken = contract.fromArtifact('KeepToken');
const TokenStaking = contract.fromArtifact('TokenStaking');
const TokenGrant = contract.fromArtifact('TokenGrant');
const Registry = contract.fromArtifact("Registry");
const PermissiveStakingPolicy = contract.fromArtifact("PermissiveStakingPolicy");
const GuaranteedMinimumStakingPolicy = contract.fromArtifact("GuaranteedMinimumStakingPolicy");

const ManagedGrant = contract.fromArtifact('ManagedGrant');
const ManagedGrantFactory = contract.fromArtifact('ManagedGrantFactory');

const nullAddress = '0x0000000000000000000000000000000000000000';
const nullBytes = '0x';

describe('ManagedGrantFactory', () => {
  let token, registry, tokenGrant, staking;
  let permissivePolicy, minimumPolicy;
  let minimumStake, grantAmount;

  const grantCreator = accounts[0],
        grantee = accounts[2],
        operator = accounts[3],
        beneficiary = accounts[4],
        authorizer = accounts[5],
        newGrantee = accounts[6],
        unrelatedAddress = accounts[7];

  let grantId, grantStart, returnedId;

  const grantUnlockingDuration = time.duration.days(60);
  const grantCliff = time.duration.days(10);

  const initializationPeriod = time.duration.minutes(10);
  const undelegationPeriod = time.duration.minutes(30);

  let factory;

  let stakeFromManagedGrant;

  before(async () => {
    token = await KeepToken.new({from: accounts[0]});
    registry = await Registry.new({from: accounts[0]});
    staking = await TokenStaking.new(
      token.address,
      registry.address,
      initializationPeriod,
      undelegationPeriod,
      {from: accounts[0]}
    );

    tokenGrant = await TokenGrant.new(token.address, {from: accounts[0]});

    await tokenGrant.authorizeStakingContract(staking.address, {from: accounts[0]});

    minimumStake = await staking.minimumStake()

    permissivePolicy = await PermissiveStakingPolicy.new()
    minimumPolicy = await GuaranteedMinimumStakingPolicy.new(staking.address);
    grantAmount = minimumStake.muln(10);

    factory = await ManagedGrantFactory.new(
      token.address,
      tokenGrant.address,
      permissivePolicy.address,
      minimumPolicy.address,
      {from: accounts[0]}
    );
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  describe("funding", async () => {
    it("adds tokens to grant pool correctly", async () => {
      await token.approveAndCall(
        factory.address, grantAmount, nullBytes, {from: grantCreator}
      );
      expect(await factory.grantFundingPool(grantCreator)).to.eq.BN(grantAmount);
    });
  });

  describe("creating managed grants", async () => {
    it("works", async () => {
      await token.approveAndCall(
        factory.address, grantAmount, nullBytes, {from: grantCreator}
      );
      grantStart = await time.latest();
      let managedGrantAddress = await factory.createGrant.call(
        grantee,
        grantAmount,
        grantUnlockingDuration,
        grantStart,
        grantCliff,
        false,
        {from: grantCreator}
      );
      await factory.createGrant(
        grantee,
        grantAmount,
        grantUnlockingDuration,
        grantStart,
        grantCliff,
        false,
        {from: grantCreator}
      );
      let managedGrant = await ManagedGrant.at(managedGrantAddress);
      let grantId = await managedGrant.grantId();
      expect(await tokenGrant.availableToStake(grantId)).to.eq.BN(grantAmount);
      expect(await managedGrant.grantee()).to.equal(grantee);
      expect(await managedGrant.grantManager()).to.equal(grantCreator);
    });

    it("doesn't let one grant more than is in their pool", async () => {
      await token.approveAndCall(
        factory.address, grantAmount, nullBytes, {from: grantCreator}
      );
      await token.transfer(unrelatedAddress, grantAmount, {from: grantCreator});
      await token.approveAndCall(
        factory.address, grantAmount, nullBytes, {from: unrelatedAddress}
      );
      grantStart = await time.latest();
      await expectRevert(
        factory.createGrant(
          grantee,
          grantAmount.addn(1),
          grantUnlockingDuration,
          grantStart,
          grantCliff,
          false,
          {from: unrelatedAddress}
        ),
        "Insufficient funding"
      );
    });
  });
});
