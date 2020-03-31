const delegateStakeFromGrant = require('../helpers/delegateStakeFromGrant')
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const { duration, increaseTimeTo } = require('../helpers/increaseTime');
const latestTime = require('../helpers/latestTime');
const expectThrowWithMessage = require('../helpers/expectThrowWithMessage');
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

describe('TokenGrant/Withdraw', function() {

  let tokenContract, registryContract, grantContract, stakingContract;

  const tokenOwner = accounts[0],
    grantee = accounts[1],
    operatorOne = accounts[2],
    magpie = accounts[4],
    authorizer = accounts[5];

  let grantId, grantStart, grantAmount;

  const grantRevocable = false;
  const grantDuration = duration.seconds(60);
  const grantCliff = duration.seconds(1);
    
  const initializationPeriod = 10;
  const undelegationPeriod = 30;

  before(async () => {
    tokenContract = await KeepToken.new({from: accounts[0]});
    registryContract = await Registry.new({from: accounts[0]});
    stakingContract = await TokenStaking.new(
      tokenContract.address, 
      registryContract.address, 
      initializationPeriod, 
      undelegationPeriod,
      {from: accounts[0]}
    );
    grantAmount = (await stakingContract.minimumStake()).muln(10);

    grantContract = await TokenGrant.new(tokenContract.address, {from: accounts[0]});
    
    await grantContract.authorizeStakingContract(stakingContract.address, {from: accounts[0]});

    grantStart = await latestTime();

    grantId = await grantTokens(
      grantContract, 
      tokenContract,
      grantAmount,
      tokenOwner, 
      grantee, 
      grantDuration, 
      grantStart, 
      grantCliff, 
      grantRevocable,
      {from: accounts[0]}
    );
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should allow to wihtdraw some tokens", async () => {
    await increaseTimeTo(grantStart + grantDuration - 30)

    const withdrawable = await grantContract.withdrawable(grantId)
    const granteeTokenGrantBalancePre = await grantContract.balanceOf(grantee)
    await grantContract.withdraw(grantId)
    const granteeTokenGrantBalancePost = await grantContract.balanceOf(grantee)

    const granteeTokenBalance = await tokenContract.balanceOf(grantee)
    const grantDetails = await grantContract.getGrant(grantId)
    
    expect(withdrawable).to.be.gt.BN(
      0,
      "Should allow to withdraw more than 0"
    )
    expect(granteeTokenBalance).to.eq.BN(
      grantDetails.withdrawn,
      "Grantee KEEP token balance should be equal to the grant withdrawn amount"
    )
    expect(granteeTokenGrantBalancePre.sub(granteeTokenGrantBalancePost)).to.eq.BN(
      grantDetails.withdrawn,
      "Grantee token grant balance should be updated"
    )
  })

  it("should allow to wihtdraw the whole grant amount ", async () => {
    await increaseTimeTo(grantStart + grantDuration)

    const withdrawablePre = await grantContract.withdrawable(grantId)
    const granteeTokenGrantBalancePre = await grantContract.balanceOf(grantee)
    await grantContract.withdraw(grantId)
    const withdrawablePost = await grantContract.withdrawable(grantId)
    const granteeTokenGrantBalancePost = await grantContract.balanceOf(grantee)

    const granteeTokenBalance = await tokenContract.balanceOf(grantee)
    const grantDetails = await grantContract.getGrant(grantId)

    expect(withdrawablePre).to.eq.BN(
      grantAmount,
      "The withdrawable amount should be equal to the whole grant amount"
    )
    expect(granteeTokenBalance).to.eq.BN(
      grantAmount,
      "Grantee KEEP token balance should be equal to the grant amount"
    )
    expect(withdrawablePost).to.eq.BN(
      0,
      "The withdrawable amount should be equal to 0, when the whole grant amount has been withdrawn"
    )
    expect(granteeTokenGrantBalancePre.sub(grantAmount)).to.eq.BN(
      granteeTokenGrantBalancePost,
      "Grantee token grant balance should be updated"
    )
    expect(grantDetails.withdrawn).to.eq.BN(
      grantAmount,
      "The grant withdrawan amount should be updated"
    )
  })

  it("should not allow to withdraw delegated tokens", async () => {
    await increaseTimeTo(grantStart + grantDuration)
    const withdrawable = await grantContract.withdrawable(grantId)
    await delegateStakeFromGrant(
        grantContract,
        stakingContract.address,
        grantee,
        operatorOne,
        magpie,
        authorizer,
        grantAmount,
        grantId
    )
    const withdrawableAfterStake = await grantContract.withdrawable(grantId)

    await expectThrowWithMessage(
      grantContract.withdraw(grantId),
      "Grant available to withdraw amount should be greater than zero."
    )
    expect(withdrawableAfterStake).to.eq.BN(
      0,
      "The withdrawable amount should be equal to 0"
    )
  })
});
