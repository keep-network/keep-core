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

const nullAddress = '0x0000000000000000000000000000000000000000';

describe('ManagedGrant', () => {
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

  let managedGrant;

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

    grantId = await tokenGrant.numGrants();

    managedGrant = await ManagedGrant.new(
      token.address,
      tokenGrant.address,
      grantCreator,
      grantId,
      grantee
    );

    grantStart = await time.latest();

    returnedId = await grantTokens(
      tokenGrant, token,
      grantAmount, grantCreator, managedGrant.address,
      grantUnlockingDuration, grantStart, grantCliff, false,
      permissivePolicy.address,
      {from: grantCreator}
    );

    stakeFromManagedGrant = async (
      operator,
      beneficiary,
      authorizer,
      amount,
      sender
    ) => {
      let delegation = Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operator.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ]);

      return managedGrant.stake(
        staking.address,
        amount,
        delegation,
        {from: sender}
      );
    }
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("is created correctly", async () => {
    expect(returnedId).to.eq.BN(grantId);
    expect(await managedGrant.grantManager()).to.equal(grantCreator);
    expect(await managedGrant.grantId()).to.eq.BN(grantId);
    expect(await managedGrant.grantee()).to.equal(grantee);
    expect(await managedGrant.requestedNewGrantee()).to.equal(nullAddress);
  });

  describe("staking", async () => {
    it("can be staked by the grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
    });

    it("can not be staked by the grant creator", async () => {
      await expectRevert(
        stakeFromManagedGrant(
          operator, beneficiary, authorizer, minimumStake, grantCreator
        ),
        "Only grantee may perform this action"
      );
    });

    it("can not be staked by a third party", async () => {
      await expectRevert(
        stakeFromManagedGrant(
          operator, beneficiary, authorizer, minimumStake, unrelatedAddress
        ),
        "Only grantee may perform this action"
      );
    });
  });

  describe("withdrawal", async () => {
    it("can be withdrawn", async () => {
      await time.increase(grantUnlockingDuration);
      await managedGrant.withdraw({from: grantee});
      expect(await token.balanceOf(grantee)).to.eq.BN(grantAmount);
    });

    it("can only be withdrawn by the grantee", async () => {
      await time.increase(grantUnlockingDuration);
      await expectRevert(
        managedGrant.withdraw({from: grantCreator}),
        "Only grantee may perform this action"
      );
      await expectRevert(
        managedGrant.withdraw({from: unrelatedAddress}),
        "Only grantee may perform this action"
      );
    });
  });

  describe("reassignment + withdrawal", async () => {
    it("can not be withdrawn when reassignment is pending", async () => {
      await time.increase(grantUnlockingDuration);
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.withdraw({from: grantee}),
        "Can not withdraw with pending reassignment"
      );
    });

    it("can be withdrawn to a reassigned grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await time.increase(grantUnlockingDuration);
      await managedGrant.withdraw({from: newGrantee});
      expect(await token.balanceOf(newGrantee)).to.eq.BN(grantAmount);
    });
  });

  describe("reassignment", async () => {
    it("can be reassigned", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      expect(await managedGrant.grantee()).to.equal(grantee);
      expect(await managedGrant.requestedNewGrantee()).to.equal(newGrantee);

      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      expect(await managedGrant.grantee()).to.equal(newGrantee);
      expect(await managedGrant.requestedNewGrantee()).to.equal(nullAddress);
    });

    it("only grantee can request reassignment", async () => {
      await expectRevert(
        managedGrant.requestGranteeReassignment(newGrantee, {from: newGrantee}),
        "Only grantee may perform this action"
      );
      await expectRevert(
        managedGrant.requestGranteeReassignment(newGrantee, {from: grantCreator}),
        "Only grantee may perform this action"
      );
      await expectRevert(
        managedGrant.requestGranteeReassignment(newGrantee, {from: unrelatedAddress}),
        "Only grantee may perform this action"
      );
    });

    it("rejects reassignment to the null address", async () => {
      await expectRevert(
        managedGrant.requestGranteeReassignment(nullAddress, {from: grantee}),
        "Invalid new grantee address"
      );
    });

    it("rejects reassignment to the old grantee's address", async () => {
      await expectRevert(
        managedGrant.requestGranteeReassignment(grantee, {from: grantee}),
        "Unchanged new grantee address"
      );
    });

    it("only grantManager can confirm reassignment", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});

      await expectRevert(
        managedGrant.confirmGranteeReassignment(newGrantee, {from: grantee}),
        "Only grantManager may perform this action"
      );
      await expectRevert(
        managedGrant.confirmGranteeReassignment(newGrantee, {from: newGrantee}),
        "Only grantManager may perform this action"
      );
      await expectRevert(
        managedGrant.confirmGranteeReassignment(newGrantee, {from: unrelatedAddress}),
        "Only grantManager may perform this action"
      );
    });

    it("requires the grant manager to confirm the new grantee address", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.confirmGranteeReassignment(grantee, {from: grantCreator}),
        "Reassignment address mismatch"
      );
      await expectRevert(
        managedGrant.confirmGranteeReassignment(grantCreator, {from: grantCreator}),
        "Reassignment address mismatch"
      );
      await expectRevert(
        managedGrant.confirmGranteeReassignment(unrelatedAddress, {from: grantCreator}),
        "Reassignment address mismatch"
      );
    });
  });
});
