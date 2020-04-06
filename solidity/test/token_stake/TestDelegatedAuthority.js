const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const { createSnapshot, restoreSnapshot } = require('../helpers/snapshot');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const KeepToken = contract.fromArtifact('KeepToken');
const TokenStaking = contract.fromArtifact('TokenStaking');
const Registry = contract.fromArtifact("Registry");
const DelegatedAuthorityStub = contract.fromArtifact("DelegatedAuthorityStub");

const initializationPeriod = time.duration.seconds(10);
const undelegationPeriod = time.duration.seconds(30);

let token, registry, stakingContract;
let authorityDelegator, badAuthorityDelegator;
let innerRecursiveDelegator, outerRecursiveDelegator;
let minimumStake, stakingAmount;

describe("TokenStaking/DelegatedAuthority", async () => {
  const owner = accounts[0];
  const operator = accounts[1];
  const magpie = accounts[2];
  const authorizer = accounts[3];
  const recognizedContract = accounts[4];
  const unrecognizedContract = accounts[5];
  const unapprovedContract = accounts[6];
  const recursivelyAuthorizedContract = accounts[7];

  before(async () => {
    token = await KeepToken.new({from: accounts[0]});
    registry = await Registry.new();

    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod
    );
    minimumStake = await stakingContract.minimumStake();
    stakingAmount = minimumStake.muln(20);
    let tx = await delegate(operator, stakingAmount);
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
    await time.increaseTo(createdAt.add(initializationPeriod).addn(1))

    authorityDelegator = await DelegatedAuthorityStub.new(recognizedContract);
    badAuthorityDelegator = await DelegatedAuthorityStub.new(unapprovedContract);
    await registry.approveOperatorContract(authorityDelegator.address);

    innerRecursiveDelegator = await DelegatedAuthorityStub.new(
      recursivelyAuthorizedContract);
    outerRecursiveDelegator = await DelegatedAuthorityStub.new(
      innerRecursiveDelegator.address);
    await registry.approveOperatorContract(outerRecursiveDelegator.address);
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  async function delegate(operator, amount) {
    let data = Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operator.substr(2), 'hex'),
      Buffer.from(authorizer.substr(2), 'hex')
    ]);

    return token.approveAndCall(
      stakingContract.address, amount,
      '0x' + data.toString('hex'),
      {from: owner}
    );
  }

  async function hasDelegatedAuthorization(operatorContract) {
    return stakingContract.isAuthorizedForOperator(operator, operatorContract);
  }

  async function authorize(operatorContract) {
    stakingContract.authorizeOperatorContract(
      operator,
      operatorContract.address,
      {from: authorizer}
    );
  }

  async function disable(operatorContract) {
    registry.disableOperatorContract(operatorContract.address);
  }

  describe("claimDelegatedAuthority", async () => {
    it("lets contracts claim delegated authority", async () => {
      await stakingContract.claimDelegatedAuthority(
        authorityDelegator.address,
        {from: recognizedContract}
      );

      expect(await stakingContract.getAuthoritySource(recognizedContract))
        .to.equal(authorityDelegator.address);
    })

    it("doesn't give unrecognized contracts delegated authority", async () => {
      await expectRevert(
        stakingContract.claimDelegatedAuthority(
          authorityDelegator.address,
          {from: unrecognizedContract}
        ),
        "Unrecognized claimant"
      );
    })

    it("doesn't give delegated authority through unapproved contracts", async () => {
      await expectRevert(
        stakingContract.claimDelegatedAuthority(
          badAuthorityDelegator.address,
          {from: unapprovedContract}
        ),
        "Operator contract is not approved"
      );
    })

    it("delegates authority recursively", async () => {
      await innerRecursiveDelegator.claimAuthorityRecursively(
        stakingContract.address,
        outerRecursiveDelegator.address
      );
      await stakingContract.claimDelegatedAuthority(
        innerRecursiveDelegator.address,
        {from: recursivelyAuthorizedContract}
      );

      expect(await stakingContract.getAuthoritySource(recursivelyAuthorizedContract))
        .to.equal(outerRecursiveDelegator.address);
    })
  })

  describe("isAuthorizedForOperator", async () => {
    before(async () => {
      await stakingContract.claimDelegatedAuthority(
        authorityDelegator.address,
        {from: recognizedContract}
      );
    })

    it("delegates authorization correctly", async () => {
      expect(await hasDelegatedAuthorization(recognizedContract)).to.be.false;
      await authorize(authorityDelegator);
      expect(await hasDelegatedAuthorization(recognizedContract)).to.be.true;
    })

    it("disables delegated authorization with the panic button", async () => {
      await authorize(authorityDelegator);
      await disable(authorityDelegator);
      // Indirect test;
      // `claimDelegatedAuthority` checks `onlyApprovedOperatorContract`
      await expectRevert(
        stakingContract.claimDelegatedAuthority(
          recognizedContract,
          {from: unrecognizedContract}
        ),
        "Operator contract is not approved"
      );
    })

    it("works recursively", async () => {
      await innerRecursiveDelegator.claimAuthorityRecursively(
        stakingContract.address,
        outerRecursiveDelegator.address
      );
      await stakingContract.claimDelegatedAuthority(
        innerRecursiveDelegator.address,
        {from: recursivelyAuthorizedContract}
      );
      await authorize(outerRecursiveDelegator);
      expect(await hasDelegatedAuthorization(recursivelyAuthorizedContract)).to.be.true;
    })
  })

  describe("authorizeOperatorContract", async () => {
    it("doesn't authorize contracts using delegated authority", async () => {
      await expectRevert(
        stakingContract.authorizeOperatorContract(
          operator,
          recognizedContract,
          {from: authorizer}
        ),
        "Contract uses delegated authority"
      );
    })
  })

  describe("slash", async () => {
    it("uses delegated authorization correctly", async () => {
      await expectRevert(
        stakingContract.slash(
          minimumStake,
          [operator],
          {from: recognizedContract}
        ),
        "Not authorized"
      );
      await authorize(authorityDelegator);
      await stakingContract.slash(
        minimumStake,
        [operator],
        {from: recognizedContract}
      );
      // no error
    })
  })

  describe("seize", async () => {
    it("uses delegated authorization correctly", async () => {
      await expectRevert(
        stakingContract.seize(
          minimumStake,
          100,
          magpie,
          [operator],
          {from: recognizedContract}
        ),
        "Not authorized"
      );
      await authorize(authorityDelegator);
      await stakingContract.seize(
        minimumStake,
        100,
        magpie,
        [operator],
        {from: recognizedContract}
      );
      // no error
    })
  })

  describe("eligibleStake", async () => {
    it("uses delegated authorization correctly", async () => {
      expect(await stakingContract.eligibleStake(operator, recognizedContract))
        .to.eq.BN(0);
      await authorize(authorityDelegator);
      expect(await stakingContract.eligibleStake(operator, recognizedContract))
        .to.eq.BN(stakingAmount);
    })
  })

  describe("activeStake", async () => {
    it("uses delegated authorization correctly", async () => {
      expect(await stakingContract.activeStake(operator, recognizedContract))
        .to.eq.BN(0);
      await authorize(authorityDelegator);
      expect(await stakingContract.activeStake(operator, recognizedContract))
        .to.eq.BN(stakingAmount);
    })
  })
})
