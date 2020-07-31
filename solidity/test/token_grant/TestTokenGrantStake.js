const { delegateStake, delegateStakeFromGrant } = require('../helpers/delegateStake')
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const {initTokenStaking} = require('../helpers/initContracts')
const {grantTokens} = require('../helpers/grantTokens');
const { createSnapshot, restoreSnapshot } = require('../helpers/snapshot');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

// Depending on test network increaseTimeTo can be inconsistent and add
// extra time. As a workaround we subtract timeRoundMargin in all cases
// that test times before initialization/undelegation periods end.
const timeRoundMargin = time.duration.minutes(1)

const KeepToken = contract.fromArtifact('KeepToken');
const TokenGrant = contract.fromArtifact('TokenGrant');
const KeepRegistry = contract.fromArtifact("KeepRegistry");
const PermissiveStakingPolicy = contract.fromArtifact("PermissiveStakingPolicy");
const GuaranteedMinimumStakingPolicy = contract.fromArtifact("GuaranteedMinimumStakingPolicy");
const EvilStakingPolicy = contract.fromArtifact("EvilStakingPolicy");

describe('TokenGrant/Stake', function() {

  let tokenContract, registryContract, grantContract, stakingContract,
    permissivePolicy, minimumPolicy, evilPolicy,
    minimumStake, grantAmount;

  const tokenOwner = accounts[0],
    grantee = accounts[1],
    operatorOne = accounts[2],
    operatorTwo = accounts[3],
    beneficiary = accounts[4],
    authorizer = accounts[5],
    revocableGrantee = accounts[6],
    evilGrantee = accounts[7];

  let grantId;
  let revocableGrantId;
  let evilGrantId;
  let grantStart;

  const grantUnlockingDuration = time.duration.days(180);
  const grantCliff = time.duration.days(10);

  const initializationPeriod = time.duration.days(10);
  let undelegationPeriod;

  before(async () => {
    tokenContract = await KeepToken.new({from: accounts[0]});
    grantContract = await TokenGrant.new(tokenContract.address, {from: accounts[0]});
    registryContract = await KeepRegistry.new({from: accounts[0]});
    const stakingContracts = await initTokenStaking(
      tokenContract.address,
      grantContract.address,
      registryContract.address,
      initializationPeriod,
      contract.fromArtifact('TokenStakingEscrow'),
      contract.fromArtifact('TokenStaking')
    );
    stakingContract = stakingContracts.tokenStaking;
    stakingEscrowContract = stakingContracts.tokenStakingEscrow;

    await grantContract.authorizeStakingContract(stakingContract.address, {from: accounts[0]});

    undelegationPeriod = await stakingContract.undelegationPeriod()

    grantStart = await time.latest();
    minimumStake = await stakingContract.minimumStake()

    permissivePolicy = await PermissiveStakingPolicy.new()
    minimumPolicy = await GuaranteedMinimumStakingPolicy.new(stakingContract.address);
    evilPolicy = await EvilStakingPolicy.new()
    grantAmount = minimumStake.muln(10),

    // Grant tokens
    grantId = await grantTokens(
      grantContract,
      tokenContract,
      grantAmount,
      tokenOwner,
      grantee,
      grantUnlockingDuration,
      grantStart,
      grantCliff,
      false,
      permissivePolicy.address,
      {from: accounts[0]}
    );

    revocableGrantId = await grantTokens(
      grantContract,
      tokenContract,
      grantAmount,
      tokenOwner,
      revocableGrantee,
      grantUnlockingDuration,
      grantStart,
      grantCliff,
      true,
      minimumPolicy.address,
      {from: accounts[0]}
    );

    evilGrantId = await grantTokens(
      grantContract,
      tokenContract,
      grantAmount,
      tokenOwner,
      evilGrantee,
      grantUnlockingDuration,
      grantStart,
      grantCliff,
      false,
      evilPolicy.address,
      {from: accounts[0]}
    );
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  async function delegate(grantee, operator, amount) {
    return await delegateStakeFromGrant(
      grantContract,
      stakingContract.address,
      grantee,
      operator,
      beneficiary,
      authorizer,
      amount,
      grantId
    )
  }

  async function delegateRevocable(grantee, operator, amount) {
    return await delegateStakeFromGrant(
      grantContract,
      stakingContract.address,
      grantee,
      operator,
      beneficiary,
      authorizer,
      amount,
      revocableGrantId
    )
  }

  async function delegateEvil(grantee, operator, amount) {
    return await delegateStakeFromGrant(
      grantContract,
      stakingContract.address,
      grantee,
      operator,
      beneficiary,
      authorizer,
      amount,
      evilGrantId
    )
  }

  it("should update balances when delegating", async () => {
    let amountToDelegate = minimumStake.muln(5);
    let remaining = grantAmount.sub(amountToDelegate)

    await delegate(grantee, operatorOne, amountToDelegate);

    let availableForStaking = await grantContract.availableToStake.call(grantId)
    let operatorBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(availableForStaking).to.eq.BN(
      remaining,
      "All granted tokens delegated, should be nothing more available"
    )
    expect(operatorBalance).to.eq.BN(
      amountToDelegate,
      "Staking amount should be added to the operator balance"
    );
  })

  it("should allow to delegate, undelegate, and recover grant", async () => {
    let tx = await delegate(grantee, operatorOne, grantAmount)
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

    await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
    tx = await grantContract.undelegate(operatorOne, {from: grantee})
    let undelegatedAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
    await time.increaseTo(undelegatedAt.add(undelegationPeriod).addn(1))
    await grantContract.recoverStake(operatorOne);

    let availableForStaking = await stakingEscrowContract.depositedAmount(operatorOne);
    let operatorBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(availableForStaking).to.eq.BN(
      grantAmount,
      "All granted tokens should be again available for staking"
    )
    expect(operatorBalance).to.eq.BN(
      0,
      "Staking amount should be removed from operator balance"
    );
  })

  it("should allow to cancel delegation right away", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await grantContract.cancelStake(operatorOne, {from: grantee});

    let availableForStaking = await stakingEscrowContract.depositedAmount.call(operatorOne)
    let operatorBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(availableForStaking).to.eq.BN(
      grantAmount,
      "All granted tokens should be again available for staking"
    )
    expect(operatorBalance).to.eq.BN(
      0,
      "Staking amount should be removed from operator balance"
    );
  })

  it("should allow to cancel delegation just before initialization period is over", async () => {
    let tx = await delegate(grantee, operatorOne, grantAmount)
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

    await time.increaseTo(createdAt.add(initializationPeriod).sub(timeRoundMargin))

    await grantContract.cancelStake(operatorOne, {from: grantee});

    let availableForStaking = await stakingEscrowContract.depositedAmount.call(operatorOne)
    let operatorBalance = await stakingContract.balanceOf.call(operatorOne);

    expect(availableForStaking).to.eq.BN(
      grantAmount,
      "All granted tokens should be again available for staking"
    )
    expect(operatorBalance).to.eq.BN(
      0,
      "Staking amount should be removed from operator balance"
    );
  })

  it("should not allow to cancel delegation after initialization period is over", async () => {
    let tx = await delegate(grantee, operatorOne, grantAmount)
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

    await time.increaseTo(createdAt.add(initializationPeriod).addn(1))

    await expectRevert(
      grantContract.cancelStake(operatorOne, {from: grantee}),
      "Initialized stake"
    );
  })

  it("should not allow to recover stake before undelegation period is over", async () => {
    let tx = await delegate(grantee, operatorOne, grantAmount)
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

    await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
    tx = await grantContract.undelegate(operatorOne, {from: grantee})
    let undelegatedAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
    await time.increaseTo(undelegatedAt.add(undelegationPeriod).sub(timeRoundMargin));

    await expectRevert(
      stakingContract.recoverStake(operatorOne),
      "Still undelegating"
    )
  })

  it("should not allow to delegate to the same operator after recovering stake", async () => {
    let tx = await delegate(grantee, operatorOne, minimumStake)
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)
    await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
    tx = await grantContract.undelegate(operatorOne, {from: grantee})
    await time.increaseTo(createdAt.add(grantUnlockingDuration))
    await grantContract.recoverStake(operatorOne, {from: grantee});

    await expectRevert(
      delegate(grantee, operatorOne, minimumStake),
      "Stake undelegated"
    )
  })

  it("should not allow to delegate to the same operator after cancelling stake", async () => {
    await delegate(grantee, operatorOne, minimumStake)
    await grantContract.cancelStake(operatorOne, {from: grantee});

    await expectRevert(
      delegate(grantee, operatorOne, minimumStake),
      "Stake for the operator already deposited in the escrow"
    )

    await time.increase(initializationPeriod.addn(1))

    await expectRevert(
      delegate(grantee, operatorOne, minimumStake),
      "Stake for the operator already deposited in the escrow"
    )
  })

  it("should allow to delegate to two different operators", async () => {
    let amountToDelegate = minimumStake.muln(5);

    await delegate(grantee, operatorOne, amountToDelegate);
    await delegate(grantee, operatorTwo, amountToDelegate);

    let availableForStaking = await grantContract.availableToStake.call(grantId)
    let operatorOneBalance = await stakingContract.balanceOf.call(operatorOne);
    let operatorTwoBalance = await stakingContract.balanceOf.call(operatorTwo);

    expect(availableForStaking).to.eq.BN(
      grantAmount.sub(amountToDelegate).sub(amountToDelegate),
      "All granted tokens delegated, should be nothing more available"
    )
    expect(operatorOneBalance).to.eq.BN(
      amountToDelegate,
      "Staking amount should be added to the operator balance"
    );
    expect(operatorTwoBalance).to.eq.BN(
      amountToDelegate,
      "Staking amount should be added to the operator balance"
    );
  })

  it("should not allow to delegate to not authorized staking contract", async () => {
    const delegation = Buffer.concat([
      Buffer.from(beneficiary.substr(2), 'hex'),
      Buffer.from(operatorOne.substr(2), 'hex'),
      Buffer.from(authorizer.substr(2), 'hex')
    ]);

    const notAuthorizedContract = "0x9E8E3487dCCd6a50045792fAfe8Ac71600B649a9"

    await expectRevert(
      grantContract.stake(
        grantId,
        notAuthorizedContract,
        grantAmount,
        delegation,
        {from: grantee}
      ),
      "Provided staking contract is not authorized"
    )
  })

  it("should not allow anyone but grantee to delegate", async () => {
    await expectRevert(
      delegate(operatorOne, operatorOne, grantAmount),
      "Only grantee of the grant can stake it."
    );
  })

  it("should let operator cancel delegation", async () => {
    await delegate(grantee, operatorOne, grantAmount, grantId);

    await grantContract.cancelStake(operatorOne, {from: operatorOne})
    // ok, no exception
  })

  it("should not allow third party to cancel delegation", async () => {
    await delegate(grantee, operatorOne, grantAmount);

    await expectRevert(
      grantContract.cancelStake(operatorOne, {from: operatorTwo}),
      "Only operator or grantee can cancel the delegation."
    );
  })

  it("should let operator undelegate", async () => {
    let tx = await delegate(grantee, operatorOne, grantAmount)
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

    await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
    await grantContract.undelegate(operatorOne, {from: operatorOne})
    // ok, no exceptions
  })

  it("should not allow third party to undelegate", async () => {
    let tx = await delegate(grantee, operatorOne, grantAmount)
    let createdAt = web3.utils.toBN((await web3.eth.getBlock(tx.receipt.blockNumber)).timestamp)

    await time.increaseTo(createdAt.add(initializationPeriod).addn(1))
    await expectRevert(
      grantContract.undelegate(operatorOne, {from: operatorTwo}),
      "Only operator or grantee can undelegate"
    )
  })

  it("should allow delegation of revocable grants", async () => {
    await delegateRevocable(revocableGrantee, operatorTwo, minimumStake);
    // ok, no exceptions
  })

  it("should not allow delegation of more than permitted", async () => {
    await expectRevert(
      delegateRevocable(revocableGrantee, operatorTwo, minimumStake.addn(1)),
      "Must have available granted amount to stake."
    );
  })

  it("should allow delegation of evil grants", async () => {
    await delegateEvil(evilGrantee, operatorTwo, grantAmount);
    // ok, no exceptions
  })

  it("should not allow delegation of more than in the grant", async () => {
    await expectRevert(
      delegateEvil(evilGrantee, operatorTwo, grantAmount.addn(1)),
      "Must have available granted amount to stake."
    );
  })
});
