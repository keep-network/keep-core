import increaseTime, {duration} from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const Registry = artifacts.require("./Registry.sol");

contract('TokenStaking', function() {
  let token, registry, stakingContract, minimumStakeBase, minimumStakeSteps, minimumStakeSchedule;
  const initializationPeriod = 10;
  const undelegationPeriod = 30;

  before(async () => {
    token = await KeepToken.new();
    registry = await Registry.new();
    stakingContract = await TokenStaking.new(
      token.address, registry.address, initializationPeriod, undelegationPeriod
    );

    minimumStakeBase = await stakingContract.minimumStakeBase();
    minimumStakeSteps = await stakingContract.minimumStakeSteps();
    minimumStakeSchedule = await stakingContract.minimumStakeSchedule();
  });

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("minimumStake", async () => {
    it("returns max value when the schedule starts", async () => {
      expect(await stakingContract.minimumStake()).to.eq.BN(
        minimumStakeBase.mul(minimumStakeSteps),
        "Unexpected minimum stake amount"
      );
    })

    it("returns max value right before the next schedule step", async () => {
      let minimumStakeScheduleStart = await stakingContract.minimumStakeScheduleStart();
      let timeForStepOne = minimumStakeScheduleStart.add(minimumStakeSchedule.div(minimumStakeSteps))
      // Rounding timestamp jump to 1 minute less (looks like increaseTime() can occasionally add extra seconds)
      await increaseTime(timeForStepOne.toNumber() - await latestTime() - duration.minutes(1))
      expect(await stakingContract.minimumStake()).to.eq.BN(
        minimumStakeBase.mul(minimumStakeSteps),
        "Unexpected minimum stake amount"
      );
    })

    it("returns half value in the middle of the schedule", async () => {
      await increaseTime(minimumStakeSchedule.divn(2).toNumber());
      expect(await stakingContract.minimumStake()).to.eq.BN(
        minimumStakeBase.mul(minimumStakeSteps.divn(2)),
        "Unexpected minimum stake amount"
      );
    })

    it("returns min value when the schedule ends", async () => {
      await increaseTime(minimumStakeSchedule.toNumber());
      expect(await stakingContract.minimumStake()).to.eq.BN(
        minimumStakeBase,
        "Unexpected minimum stake amount"
      );
    })
  })
});
