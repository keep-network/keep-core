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
  });

  it("can be staked by the grantee", async () => {
    await stakeFromManagedGrant(
      operator, beneficiary, authorizer, minimumStake, grantee
    );
  })

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
})
