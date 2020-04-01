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

let token, registry, stakingContract, authorityDelegator, badAuthorityDelegator;
let minimumStake, stakingAmount;

contract("TokenStaking/DelegatedAuthority", async (accounts) => {
  const owner = accounts[0];
  const operator = accounts[1];
  const magpie = accounts[2];
  const authorizer = accounts[3];
  const recognizedContract = accounts[4];
  const unrecognizedContract = accounts[5];
  const unapprovedContract = accounts[6];

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

  describe("claimDelegatedAuthority", async () => {
    it("lets contracts claim delegated authority", async () => {
      await stakingContract.claimDelegatedAuthority(
        authorityDelegator.address,
        {from: recognizedContract}
      );

      expect(await stakingContract.getDelegatedAuthority(recognizedContract))
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
  })
})
