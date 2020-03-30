import increaseTime, {duration, increaseTimeTo} from '../helpers/increaseTime';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage'
import {createSnapshot, restoreSnapshot} from "../helpers/snapshot"

const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const Registry = artifacts.require("./Registry.sol");

contract.only('TokenStaking', function(accounts) {
  let token, registry, stakingContract, stakingAmount, minimumStake;
  const owner = accounts[0],
    operator = accounts[1],
    magpie = accounts[2],
    authorizer = accounts[3],
    operatorContract = accounts[4];

  const initializationPeriod = duration.minutes(10);
  const undelegationPeriod = duration.minutes(10);
  const lockPeriod = duration.weeks(12);

  before(async () => {
    token = await KeepToken.new();
    registry = await Registry.new();
    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod
    );

    await registry.approveOperatorContract(operatorContract);

    minimumStake = await stakingContract.minimumStake();
    stakingAmount = minimumStake.muln(20);

    let tx = await delegate(operator, stakingAmount)
    let createdAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
    await stakingContract.lockStake(operator, lockPeriod, {from: operatorContract})

    await increaseTimeTo(createdAt + initializationPeriod + 1)
    tx = await stakingContract.undelegate(operator, {from: operator})
    let undelegatedAt = (await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp
    await increaseTimeTo(undelegatedAt + undelegationPeriod + 1)
  });

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

  it("should only permit recover unlocked stake", async () => {
    await expectThrowWithMessage(
      stakingContract.recoverStake(operator),
      "Can not recover locked stake"
    )

    await expectThrowWithMessage(
      stakingContract.unlockStake(operator),
      "Not authorized"
    )

    stakingContract.unlockStake(operator, {from: operatorContract})
    await stakingContract.recoverStake(operator)
    // ok, no revert
  })

  it("should allow recover locked stake after lock duration has expired", async () => {
    await expectThrowWithMessage(
      stakingContract.recoverStake(operator),
      "Can not recover locked stake"
    )

    await increaseTime(lockPeriod)
    await stakingContract.recoverStake(operator, {from: operator})
    // ok, no revert
  })

  it("should allow recover locked stake after operator contract has been disabled", async () => {
    await expectThrowWithMessage(
      stakingContract.recoverStake(operator),
      "Can not recover locked stake"
    )

    // disable operator contract with panic button
    await registry.disableOperatorContract(operatorContract)

    await stakingContract.recoverStake(operator, {from: operator})
    // ok, no revert
  })
});
