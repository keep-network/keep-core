const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const {grantTokens} = require('../helpers/grantTokens');
const { createSnapshot, restoreSnapshot } = require('../helpers/snapshot');
const {initTokenStaking} = require('../helpers/initContracts')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const KeepToken = contract.fromArtifact('KeepToken');
const TokenGrant = contract.fromArtifact('TokenGrant');
const KeepRegistry = contract.fromArtifact("KeepRegistry");
const PermissiveStakingPolicy = contract.fromArtifact("PermissiveStakingPolicy");
const GuaranteedMinimumStakingPolicy = contract.fromArtifact("GuaranteedMinimumStakingPolicy");

const ManagedGrant = contract.fromArtifact('ManagedGrant');

const nullAddress = '0x0000000000000000000000000000000000000000';

describe('TokenGrant/ManagedGrant', () => {
  let token, registry, tokenGrant, staking, stakingEscrow;
  let permissivePolicy, minimumPolicy;
  let minimumStake, grantAmount;

  const grantCreator = accounts[0],
        grantee = accounts[2],
        operator = accounts[3],
        beneficiary = accounts[4],
        authorizer = accounts[5],
        newGrantee = accounts[6],
        anotherGrantee = accounts[7],
        unrelatedAddress = accounts[8];

  let grantId, grantStart, returnedId;

  const grantUnlockingDuration = time.duration.days(60);
  const grantCliff = time.duration.days(10);

  const initializationPeriod = time.duration.hours(12);
  let undelegationPeriod

  let managedGrant;

  let stakeFromManagedGrant;

  before(async () => {
    token = await KeepToken.new({from: grantCreator});
    registry = await KeepRegistry.new({from: grantCreator});
    tokenGrant = await TokenGrant.new(token.address, {from: grantCreator});
    const contracts = await initTokenStaking(
      token.address,
      tokenGrant.address,
      registry.address,
      initializationPeriod,
      contract.fromArtifact('TokenStakingEscrow'),
      contract.fromArtifact('TokenStaking')
    );
    staking = contracts.tokenStaking;
    stakingEscrow = contracts.tokenStakingEscrow;

    await tokenGrant.authorizeStakingContract(staking.address, {from: grantCreator});

    undelegationPeriod = await staking.undelegationPeriod()
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
        "New grantee same as current grantee"
      );
    });

    it("emits an event", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      expect((await managedGrant.getPastEvents())[0].args['newGrantee'])
        .to.equal(newGrantee);
    });
  });

  describe("cancelReassignmentRequest", async () => {
    it("can be done by the grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.cancelReassignmentRequest({from: grantee});
      expect(await managedGrant.grantee()).to.equal(grantee);
      expect(await managedGrant.requestedNewGrantee()).to.equal(nullAddress);
    });

    it("cannot be done by the new grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.cancelReassignmentRequest({from: newGrantee}),
        "Only grantee may perform this action"
      );
    });

    it("cannot be done by the creator", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.cancelReassignmentRequest({from: grantCreator}),
        "Only grantee may perform this action"
      );
    });

    it("cannot be done by a third party", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.cancelReassignmentRequest({from: unrelatedAddress}),
        "Only grantee may perform this action"
      );
    });

    it("emits an event", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.cancelReassignmentRequest({from: grantee});
      expect((await managedGrant.getPastEvents())[0].args['cancelledRequestedGrantee'])
        .to.equal(newGrantee);
    });
  });

  describe("changeReassignmentRequest", async () => {
    it("can be done by the grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.changeReassignmentRequest(anotherGrantee, {from: grantee});
      expect(await managedGrant.grantee()).to.equal(grantee);
      expect(await managedGrant.requestedNewGrantee()).to.equal(anotherGrantee);
    });

    it("cannot be done by the previous new grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.changeReassignmentRequest(anotherGrantee, {from: newGrantee}),
        "Only grantee may perform this action"
      );
    });

    it("cannot be done by the changed new grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.changeReassignmentRequest(anotherGrantee, {from: anotherGrantee}),
        "Only grantee may perform this action"
      );
    });

    it("cannot be done by the creator", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.requestGranteeReassignment(anotherGrantee, {from: grantCreator}),
        "Only grantee may perform this action"
      );
    });

    it("cannot be done by a third party", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.requestGranteeReassignment(anotherGrantee, {from: unrelatedAddress}),
        "Only grantee may perform this action"
      );
    });

    it("rejects reassignment to the null address", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.changeReassignmentRequest(nullAddress, {from: grantee}),
        "Invalid new grantee address"
      );
    });

    it("rejects reassignment to the old grantee's address", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.changeReassignmentRequest(grantee, {from: grantee}),
        "New grantee same as current grantee"
      );
    });

    it("rejects reassignment to the previously requested address", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await expectRevert(
        managedGrant.changeReassignmentRequest(newGrantee, {from: grantee}),
        "Unchanged reassignment request"
      );
    });

    it("emits an event", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.changeReassignmentRequest(anotherGrantee, {from: grantee});
      let event = (await managedGrant.getPastEvents())[0];
      expect(event.args['previouslyRequestedGrantee'])
        .to.equal(newGrantee);
      expect(event.args['newRequestedGrantee'])
        .to.equal(anotherGrantee);
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

    it("emits an event", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      let event = (await managedGrant.getPastEvents())[0];
      expect(event.args['oldGrantee'])
        .to.equal(grantee);
      expect(event.args['newGrantee'])
        .to.equal(newGrantee);
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

    it("can't withdraw when reassignment is pending", async () => {
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

    it("emits an event", async () => {
      await time.increase(grantUnlockingDuration);
      await managedGrant.withdraw({from: grantee});
      let event = (await managedGrant.getPastEvents())[0];
      expect(event.args['destination'])
        .to.equal(grantee);
      expect(event.args['amount'])
        .to.eq.BN(grantAmount);
    });
  });

  describe("stake", async () => {
    it("can be done by the grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      expect(await staking.balanceOf(operator)).to.eq.BN(minimumStake);
      expect(await staking.beneficiaryOf(operator)).to.equal(beneficiary);
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

    it("can be done by a reassigned grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, newGrantee
      );
      expect(await staking.balanceOf(operator)).to.eq.BN(minimumStake);
      expect(await staking.beneficiaryOf(operator)).to.equal(beneficiary);
      expect(await staking.authorizerOf(operator)).to.equal(authorizer);
    });

    it("can't be done by the old grantee", async () => {
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await expectRevert(
        stakeFromManagedGrant(
          operator, beneficiary, authorizer, minimumStake, grantee
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

    it("can be done by a reassigned grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await managedGrant.cancelStake(operator, {from: newGrantee});
      expect(await staking.balanceOf(operator)).to.eq.BN(0);
    });

    it("can't be done by the old grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await expectRevert(
        managedGrant.cancelStake(operator, {from: grantee}),
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
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.undelegate(operator, {from: grantee});
    });

    it("can be done by the operator", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      expect(await token.balanceOf(managedGrant.address)).to.eq.BN(0);
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.undelegate(operator, {from: operator});
    });

    it("can't be done by the grant manager", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await time.increase(initializationPeriod.muln(2));
      await expectRevert(
        managedGrant.undelegate(operator, {from: grantCreator}),
        "Only grantee or operator may perform this action"
      );
    });

    it("can't be done by a third party", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, minimumStake, grantee
      );
      await time.increase(initializationPeriod.muln(2));
      await expectRevert(
        managedGrant.undelegate(operator, {from: unrelatedAddress}),
        "Only grantee or operator may perform this action"
      );
    });

    it("can be done by a reassigned grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      expect(await token.balanceOf(managedGrant.address)).to.eq.BN(0);
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await managedGrant.undelegate(operator, {from: newGrantee});
    });

    it("can't be done by the old grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      expect(await token.balanceOf(managedGrant.address)).to.eq.BN(0);
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await expectRevert(
        managedGrant.undelegate(operator, {from: grantee}),
        "Only grantee or operator may perform this action"
      );
    });
  });

  describe("recoverStake", async () => {
    it("can be done by the grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod.add(grantUnlockingDuration));
      await managedGrant.recoverStake(operator, {from: grantee});
      await stakingEscrow.withdrawToManagedGrantee(operator, {from: grantee});
      expect(await token.balanceOf(grantee)).to.eq.BN(grantAmount);
    });

    it("can be done by the operator", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod.add(grantUnlockingDuration));
      await managedGrant.recoverStake(operator, {from: operator});
      await stakingEscrow.withdrawToManagedGrantee(operator, {from: operator});
      expect(await token.balanceOf(grantee)).to.eq.BN(grantAmount);
    });

    it("can be done by the grant manager", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod.add(grantUnlockingDuration));
      await managedGrant.recoverStake(operator, {from: grantCreator});
      await stakingEscrow.withdrawToManagedGrantee(operator, {from: grantee});
      expect(await token.balanceOf(grantee)).to.eq.BN(grantAmount);
    });

    it("can be done by a third party", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod.add(grantUnlockingDuration));
      await managedGrant.recoverStake(operator, {from: unrelatedAddress});
      await stakingEscrow.withdrawToManagedGrantee(operator, {from: grantee});
      expect(await token.balanceOf(grantee)).to.eq.BN(grantAmount);
    });

    it("can be done by a reassigned grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod.add(grantUnlockingDuration));
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await managedGrant.recoverStake(operator, {from: newGrantee});
      await stakingEscrow.withdrawToManagedGrantee(operator, {from: newGrantee});
      expect(await token.balanceOf(newGrantee)).to.eq.BN(grantAmount);
    });

    it("can be done by the old grantee", async () => {
      await stakeFromManagedGrant(
        operator, beneficiary, authorizer, grantAmount, grantee
      );
      await time.increase(initializationPeriod.muln(2));
      await managedGrant.undelegate(operator, {from: grantee});
      await time.increase(undelegationPeriod.add(grantUnlockingDuration));
      await managedGrant.requestGranteeReassignment(newGrantee, {from: grantee});
      await managedGrant.confirmGranteeReassignment(newGrantee, {from: grantCreator});
      await managedGrant.recoverStake(operator, {from: grantee});
      await stakingEscrow.withdrawToManagedGrantee(operator, {from: newGrantee});
      expect(await token.balanceOf(newGrantee)).to.eq.BN(grantAmount);
    });
  });
});
