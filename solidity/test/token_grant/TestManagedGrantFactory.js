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

const nullBytes = '0x';

describe('TokenGrant/ManagedGrantFactory', () => {
  let token, registry, tokenGrant, staking;
  let permissivePolicy, minimumPolicy;
  let minimumStake, grantAmount;

  const grantCreator = accounts[0],
        grantee = accounts[2],
        unrelatedAddress = accounts[3];

  let grantStart;

  const grantUnlockingDuration = time.duration.days(60);
  const grantCliff = time.duration.days(10);

  const initializationPeriod = time.duration.minutes(10);
  const undelegationPeriod = time.duration.minutes(30);

  let factory;

  before(async () => {
    token = await KeepToken.new({from: grantCreator});
    registry = await Registry.new({from: grantCreator});
    staking = await TokenStaking.new(
      token.address,
      registry.address,
      initializationPeriod,
      undelegationPeriod,
      {from: grantCreator}
    );

    tokenGrant = await TokenGrant.new(token.address, {from: grantCreator});

    await tokenGrant.authorizeStakingContract(staking.address, {from: grantCreator});

    minimumStake = await staking.minimumStake()

    permissivePolicy = await PermissiveStakingPolicy.new()
    minimumPolicy = await GuaranteedMinimumStakingPolicy.new(staking.address);
    grantAmount = minimumStake.muln(10);

    factory = await ManagedGrantFactory.new(
      token.address,
      tokenGrant.address,
      permissivePolicy.address,
      minimumPolicy.address,
      {from: grantCreator}
    );
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  describe("creating managed grants", async () => {
    it("works", async () => {
      await token.approve(
        factory.address, grantAmount, {from: grantCreator}
      );
      grantStart = await time.latest();
      let managedGrantAddress = await factory.createManagedGrant.call(
        grantee,
        grantAmount,
        grantUnlockingDuration,
        grantStart,
        grantCliff,
        false,
        {from: grantCreator}
      );
      await factory.createManagedGrant(
        grantee,
        grantAmount,
        grantUnlockingDuration,
        grantStart,
        grantCliff,
        false,
        {from: grantCreator}
      );
      let event = (await factory.getPastEvents())[0];
      let managedGrant = await ManagedGrant.at(managedGrantAddress);
      let grantId = await managedGrant.grantId();
      expect(await tokenGrant.availableToStake(grantId)).to.eq.BN(grantAmount);
      expect(await managedGrant.grantee()).to.equal(grantee);
      expect(await managedGrant.grantManager()).to.equal(grantCreator);
      expect(event.args['grantAddress']).to.equal(managedGrantAddress);
      expect(event.args['grantee']).to.equal(grantee);
    });

    it("doesn't let one grant more than they've approved on the token", async () => {
      await token.transfer(unrelatedAddress, grantAmount, {from: grantCreator});
      await token.approve(
        factory.address, grantAmount.subn(1), {from: unrelatedAddress}
      );
      grantStart = await time.latest();
      await expectRevert(
        factory.createManagedGrant(
          grantee,
          grantAmount,
          grantUnlockingDuration,
          grantStart,
          grantCliff,
          false,
          {from: unrelatedAddress}
        ),
        "SafeERC20: low-level call failed -- Reason given: SafeERC20: low-level call failed."
      );
    });

    it("doesn't let one grant more than they have on the token", async () => {
      await token.transfer(unrelatedAddress, grantAmount, {from: grantCreator});
      await token.approve(
        factory.address, grantAmount.addn(1), {from: unrelatedAddress}
      );
      grantStart = await time.latest();
      await expectRevert(
        factory.createManagedGrant(
          grantee,
          grantAmount.addn(1),
          grantUnlockingDuration,
          grantStart,
          grantCliff,
          false,
          {from: unrelatedAddress}
        ),
        "SafeERC20: low-level call failed -- Reason given: SafeERC20: low-level call failed."
      );
    });
  });
});
