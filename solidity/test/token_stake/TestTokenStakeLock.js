const {contract, accounts, web3} = require("@openzeppelin/test-environment");
const {expectRevert, time} = require("@openzeppelin/test-helpers");
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot.js');

const KeepToken = contract.fromArtifact('KeepToken');
const TokenStaking = contract.fromArtifact('TokenStaking');
const Registry = contract.fromArtifact("Registry");

describe('TokenStaking/Lock', () => {
  let token, registry, stakingContract, stakingAmount, minimumStake;
  const owner = accounts[0],
    operator1 = accounts[1],
    operator2 = accounts[2],
    operator3 = accounts[3],
    magpie = accounts[4],
    authorizer = accounts[5],
    operatorContract = accounts[6],
    operatorContract2 = accounts[7];

  const initializationPeriod = time.duration.minutes(10);
  const undelegationPeriod = time.duration.minutes(10);
  const lockPeriod = time.duration.weeks(12);

  let createdAt;
  let operator;

  before(async () => {
    token = await KeepToken.new({from: owner});
    registry = await Registry.new({from: owner});
    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod,
      {from: owner}
    );

    await registry.approveOperatorContract(operatorContract, {from: owner});
    await registry.approveOperatorContract(operatorContract2, {from: owner});

    minimumStake = await stakingContract.minimumStake();
    stakingAmount = minimumStake.muln(20);
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

  async function timestampOf(tx) {
    return web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp);
  }

  async function undelegate(operator) {
    let tx = await stakingContract.undelegate(operator, {from: operator})
    let undelegatedAt = await timestampOf(tx);
    await time.increaseTo(undelegationPeriod.add(undelegatedAt).addn(1))
  }

  describe("setting locks", async () => {
    before(async () => {
      operator = operator1;
      let tx = await delegate(operator, stakingAmount)
      await stakingContract.authorizeOperatorContract(
        operator,
        operatorContract,
        { from: authorizer },
      );

      createdAt = await timestampOf(tx);
    })

    it("should not permit locks on non-initialized operators", async () => {
      await expectRevert(
        stakingContract.lockStake(operator, lockPeriod, {from: operatorContract}),
        "Operator stake must be active"
      )
    })

    it("should not permit locks on undelegating operators", async () => {
      await time.increaseTo(initializationPeriod.add(createdAt).addn(1))
      let tx = await stakingContract.undelegate(operator, {from: operator})
      let undelegatedAt = await timestampOf(tx)
      await time.increaseTo(undelegatedAt.addn(1))
      await expectRevert(
        stakingContract.lockStake(operator, lockPeriod, {from: operatorContract}),
        "Operator undelegating"
      )
    })

    it("should not permit locks from unauthorized operator contracts", async () => {
      await time.increaseTo(initializationPeriod.add(createdAt).addn(1))
      await expectRevert(
        stakingContract.lockStake(operator, lockPeriod, {from: operatorContract2}),
        "Not authorized"
      )
    })

    it("should not permit locks from disabled operator contracts", async () => {
      await time.increaseTo(initializationPeriod.add(createdAt).addn(1))
      await registry.disableOperatorContract(operatorContract, {from: owner})
      await expectRevert(
        stakingContract.lockStake(operator, lockPeriod, {from: operatorContract}),
        "Operator contract is not approved"
      )
    })

    it("should not permit locks from unapproved operator contracts", async () => {
      await time.increaseTo(initializationPeriod.add(createdAt).addn(1))
      await expectRevert(
        stakingContract.lockStake(operator, lockPeriod, {from: operator}),
        "Operator contract is not approved"
      )
    })

    it("should not permit locks that exceed the maximum lock duration", async () => {
      await time.increaseTo(initializationPeriod.add(createdAt).addn(1))
      let maximumDuration = await stakingContract.maximumLockDuration();
      let longPeriod = maximumDuration.addn(1);
      await expectRevert(
        stakingContract.lockStake(operator, longPeriod, {from: operatorContract}),
        "Lock duration too long"
      )

    })
  })

  describe("single lock", async () => {
    before(async () => {
      operator = operator2;
      let tx = await delegate(operator, stakingAmount)
      await stakingContract.authorizeOperatorContract(
        operator,
        operatorContract,
        { from: authorizer },
      );

      createdAt = await timestampOf(tx);
      await time.increaseTo(initializationPeriod.add(createdAt).addn(1))
      await stakingContract.lockStake(operator, lockPeriod, {from: operatorContract})
    })

    it("should only permit recover unlocked stake", async () => {
      await undelegate(operator)
      await expectRevert(
        stakingContract.recoverStake(operator),
        "Can not recover locked stake"
      )

      await stakingContract.unlockStake(operator, {from: operatorContract})
      await stakingContract.recoverStake(operator)
      // ok, no revert
    })

    it("should allow recover locked stake after lock duration has expired", async () => {
      await undelegate(operator)
      await expectRevert(
        stakingContract.recoverStake(operator),
        "Can not recover locked stake"
      )

      await time.increase(lockPeriod)
      await stakingContract.recoverStake(operator, {from: operator})
      // ok, no revert
    })

    it("should allow recover locked stake after operator contract has been disabled", async () => {
      await undelegate(operator)
      await expectRevert(
        stakingContract.recoverStake(operator),
        "Can not recover locked stake"
      )

      // disable operator contract with panic button
      await registry.disableOperatorContract(operatorContract, {from: owner});

      await stakingContract.recoverStake(operator, {from: operator})
      // ok, no revert
    })

    it("should be able to reduce the duration of existing locks", async () => {
      await stakingContract.lockStake(
        operator,
        undelegationPeriod.add(time.duration.minutes(5)),
        {from: operatorContract}
      )

      await undelegate(operator)
      // 5 minutes left in lock
      await expectRevert(
        stakingContract.recoverStake(operator),
        "Can not recover locked stake"
      )

      await time.increase(time.duration.minutes(5))
      await stakingContract.recoverStake(operator, {from: operator})
      // ok, no revert
    })

    it("should allow slashing/seizing locked stake after undelegation", async () => {
      await undelegate(operator)

      await stakingContract.slash(
        minimumStake, [operator],
        {from: operatorContract}
      )
      await stakingContract.seize(
        minimumStake, 100, magpie, [operator],
        {from: operatorContract}
      )
      // ok, no revert
    })

    it("should not allow slashing/seizing non-locked stake after undelegation", async () => {
      await undelegate(operator)

      await time.increase(lockPeriod)

      await expectRevert(
        stakingContract.slash(
          minimumStake, [operator],
          {from: operatorContract}
        ),
        "Stake is released"
      )
      await expectRevert(
        stakingContract.seize(
          minimumStake, 100, magpie, [operator],
          {from: operatorContract}
        ),
        "Stake is released"
      )
    })

    it("should not allow slashing/seizing unlocked stake after undelegation", async () => {
      await undelegate(operator)
      await stakingContract.unlockStake(operator, {from: operatorContract})

      await expectRevert(
        stakingContract.slash(
          minimumStake, [operator],
          {from: operatorContract}
        ),
        "Stake is released"
      )

      await expectRevert(
        stakingContract.seize(
          minimumStake, 100, magpie, [operator],
          {from: operatorContract}
        ),
        "Stake is released"
      )
    })

    it("should only allow the lock creator to slash/seize after undelegation", async () => {
      await stakingContract.authorizeOperatorContract(
        operator,
        operatorContract2,
        { from: authorizer },
      );
      await undelegate(operator)

      await expectRevert(
        stakingContract.slash(
          minimumStake, [operator],
          {from: operatorContract2}
        ),
        "Stake is released"
      )
      await expectRevert(
        stakingContract.seize(
          minimumStake, 100, magpie, [operator],
          {from: operatorContract2}
        ),
        "Stake is released"
      )
    })
  })


  describe("multiple locks", async () => {
    before(async () => {
      operator = operator3;
      let tx = await delegate(operator, stakingAmount)
      await stakingContract.authorizeOperatorContract(
        operator,
        operatorContract,
        { from: authorizer },
      );

      await stakingContract.authorizeOperatorContract(
        operator,
        operatorContract2,
        { from: authorizer },
      );

      createdAt = await timestampOf(tx);
      await time.increaseTo(initializationPeriod.add(createdAt).addn(1));
      await stakingContract.lockStake(operator, lockPeriod, {from: operatorContract})
      await stakingContract.lockStake(operator, lockPeriod, {from: operatorContract2})
    })

    it("should require all locks to be released before recovering tokens", async () => {
      await undelegate(operator)
      await stakingContract.unlockStake(operator, {from: operatorContract})

      await expectRevert(
        stakingContract.recoverStake(operator),
        "Can not recover locked stake"
      )

      await stakingContract.unlockStake(operator, {from: operatorContract2})

      await stakingContract.recoverStake(operator, {from: operator})
      // ok, no revert
    })

    it("should count disabled contracts' locks as invalid", async () => {
      await undelegate(operator)
      await stakingContract.unlockStake(operator, {from: operatorContract})

      await expectRevert(
        stakingContract.recoverStake(operator),
        "Can not recover locked stake"
      )

      await registry.disableOperatorContract(operatorContract2, {from: owner});

      await stakingContract.recoverStake(operator, {from: operator})
      // ok, no revert
    })
  })
});

