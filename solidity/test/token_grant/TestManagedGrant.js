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

  describe("requestGranteeReassignment", async () => {
    it("can be done by the grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      expect(await managedGrant.grantee()).to.equal(grantee);
      expect(await managedGrant.requestedNewGrantee()).to.equal(newGrantee);
    });

    it("cannot be done by the new grantee", async () => {
      await expectRevert(
        managedGrant.requestGranteeReassignment(newGrantee, {from: newGrantee}),
        "Only grantee may perform this action"
      );
    });

    it("cannot be done by the creator", async () => {
      await expectRevert(
        managedGrant.requestGranteeReassignment(newGrantee, {from: grantCreator}),
        "Only grantee may perform this action"
      );
    });

    it("cannot be done by a third party", async () => {
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
  });

  describe("confirmGranteeReassignment", async () => {
    it("can be done by the grant manager", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      expect(await managedGrant.grantee()).to.equal(newGrantee);
      expect(await managedGrant.requestedNewGrantee()).to.equal(nullAddress);
    });

    it("can't be done by the old grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.confirmGranteeReassignment(newGrantee, {from: grantee}),
        "Only grantManager may perform this action"
      );
    });

    it("can't be done by the new grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.confirmGranteeReassignment(newGrantee, {from: newGrantee}),
        "Only grantManager may perform this action"
      );
    });

    it("can't be done by a third party", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.confirmGranteeReassignment(newGrantee, {from: unrelatedAddress}),
        "Only grantManager may perform this action"
      );
    });

    it("can't confirm the old grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.confirmGranteeReassignment(grantee, {from: grantCreator}),
        "Reassignment address mismatch"
      );
    });

    it("can't confirm the grant manager", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.confirmGranteeReassignment(grantCreator, {from: grantCreator}),
        "Reassignment address mismatch"
      );
    });

    it("can't confirm a third party", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.confirmGranteeReassignment(unrelatedAddress, {from: grantCreator}),
        "Reassignment address mismatch"
      );
    });
  });

  describe("withdrawal", async () => {
    it("can be done by the grantee", async () => {
      await time.increase(grantUnlockingDuration);
      await managedGrant.withdraw({from: grantee});
      expect(await token.balanceOf(grantee)).to.eq.BN(grantAmount);
    });

    it("can't be done by the grant manager", async () => {
      await time.increase(grantUnlockingDuration);
      await expectRevert(
        managedGrant.withdraw({from: grantCreator}),
        "Only grantee may perform this action"
      );
    });

    it("can't be done by a third party", async () => {
      await time.increase(grantUnlockingDuration);
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

  describe("stake", async () => {
    it("can be done by the grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      expect(await staking.balanceOf(operator)).to.eq.BN(minimumStake);
      expect(await staking.magpieOf(operator)).to.equal(beneficiary);
      expect(await staking.authorizerOf(operator)).to.equal(authorizer);
    });

    it("can not be done by the grant creator", async () => {
      await expectRevert(
        stakeFromManagedGrant(
          operator, beneficiary, authorizer, minimumStake, grantCreator
        ),
        "Only grantee may perform this action"
      );
    });

    it("can not be done by a third party", async () => {
      await expectRevert(
        stakeFromManagedGrant(
          operator, beneficiary, authorizer, minimumStake, unrelatedAddress
        ),
        "Only grantee may perform this action"
      );
    });
  });

  describe("cancelStake", async () => {
    it("can be done by the grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await managedGrant.cancelStake(operator, {from: grantee});
      expect(await staking.balanceOf(operator)).to.eq.BN(0);
    });

    it("can be done by the operator", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await managedGrant.cancelStake(operator, {from: operator});
      expect(await staking.balanceOf(operator)).to.eq.BN(0);
    });

    it("can not be done by the grant manager", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await expectRevert(
        managedGrant.cancelStake(operator, {from: grantCreator}),
        "Only grantee or operator may perform this action"
      );
    });

    it("can not be done by a third party", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await expectRevert(
        managedGrant.cancelStake(operator, {from: unrelatedAddress}),
        "Only grantee or operator may perform this action"
      );
    });
  });

  describe("undelegate", async () => {
    it("can be done by the grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      expect(await token.balanceOf(managedGrant.address)).to.eq.BN(0);
      await time.increase(initializationPeriod * 2);
      await managedGrant.undelegate(operator, {from: grantee});
    });

    it("can be done by the operator", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      expect(await token.balanceOf(managedGrant.address)).to.eq.BN(0);
      await time.increase(initializationPeriod * 2);
      await managedGrant.undelegate(operator, {from: operator});
    });

    it("can't be done by the grant manager", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await time.increase(initializationPeriod * 2);
      await expectRevert(
        managedGrant.undelegate(operator, {from: grantCreator}),
        "Only grantee or operator may perform this action"
      );
    });

    it("can't be done by a third party", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await time.increase(initializationPeriod * 2);
      await expectRevert(
        managedGrant.undelegate(operator, {from: unrelatedAddress}),
        "Only grantee or operator may perform this action"
      );
    });
  });

  describe("recoverStake", async () => {
    it("can be done by the grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod * 2);
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod + grantUnlockingDuration);
      await managedGrant.recoverStake(operator, {from: grantee});
      await managedGrant.withdraw({from: grantee});
      expect(await token.balanceOf(grantee)).to.eq.BN(grantAmount);
    });

    it("can't be done by the operator", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod * 2);
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod + grantUnlockingDuration);
      await expectRevert(
        managedGrant.recoverStake(operator, {from: operator}),
        "Only grantee may perform this action"
      );
    });

    it("can't be done by the grant manager", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod * 2);
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod + grantUnlockingDuration);
      await expectRevert(
        managedGrant.recoverStake(operator, {from: grantCreator}),
        "Only grantee may perform this action"
      );
    });

    it("can't be done by a third party", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod * 2);
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod + grantUnlockingDuration);
      await expectRevert(
        managedGrant.recoverStake(operator, {from: unrelatedAddress}),
        "Only grantee may perform this action"
      );
    });
  });
});
