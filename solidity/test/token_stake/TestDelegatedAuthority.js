import increaseTime, {duration, increaseTimeTo} from '../helpers/increaseTime';
import latestTime from '../helpers/latestTime';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage'
import {createSnapshot, restoreSnapshot} from "../helpers/snapshot"

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

// Depending on test network increaseTimeTo can be inconsistent and add
// extra time. As a workaround we subtract timeRoundMargin in all cases
// that test times before initialization/undelegation periods end.
const timeRoundMargin = duration.minutes(1)

const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const Registry = artifacts.require("./Registry.sol");
const DelegatedAuthorityStub = artifacts.require("./stubs/DelegatedAuthorityStub.sol");

const initializationPeriod = 10;
const undelegationPeriod = 30;

let token, registry, stakingContract;
let authorityDelegator, badAuthorityDelegator;
let innerRecursiveDelegator, outerRecursiveDelegator;
let minimumStake, stakingAmount;

contract("TokenStaking/DelegatedAuthority", async (accounts) => {
  const owner = accounts[0];
  const operator = accounts[1];
  const magpie = accounts[2];
  const authorizer = accounts[3];
  const recognizedContract = accounts[4];
  const unrecognizedContract = accounts[5];
  const unapprovedContract = accounts[6];
  const recursivelyAuthorizedContract = accounts[7];

  before(async () => {
    token = await KeepToken.new();
    registry = await Registry.new();

    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod
    );
    minimumStake = await stakingContract.minimumStake();
    stakingAmount = minimumStake.muln(20);
    let tx = await delegate(operator, stakingAmount);
    let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
    await increaseTimeTo(createdAt + initializationPeriod + 1)

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
      await expectThrowWithMessage(
        stakingContract.claimDelegatedAuthority(
          authorityDelegator.address,
          {from: unrecognizedContract}
        ),
        "Unrecognized claimant"
      );
    })

    it("doesn't give delegated authority through unapproved contracts", async () => {
      await expectThrowWithMessage(
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
      await expectThrowWithMessage(
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

  describe("slash", async () => {
    it("uses delegated authorization correctly", async () => {
      await expectThrowWithMessage(
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
      await expectThrowWithMessage(
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

  describe("lockStake", async () => {
    it("uses delegated authorization correctly", async () => {
      let lockPeriod = duration.weeks(12);
      await expectThrowWithMessage(
        stakingContract.lockStake(
          operator,
          lockPeriod,
          {from: recognizedContract}
        ),
        "Not authorized"
      );
      await authorize(authorityDelegator);
      await stakingContract.lockStake(
        operator,
        lockPeriod,
        {from: recognizedContract}
      );
      // no error
    })
  })
})
