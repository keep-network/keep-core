import { helpers } from "hardhat"

const { to1e18 } = helpers.number

export const constants = {
  groupSize: 100,
  groupThreshold: 51,
  offchainDkgTime: 72, // 5 * (1 + 5) + 2 * (1 + 10) + 20
  minimumStake: to1e18(100000),
  poolWeightDivisor: to1e18(1),
}

export const dkgState = {
  IDLE: 0,
  AWAITING_SEED: 1,
  KEY_GENERATION: 2,
  AWAITING_RESULT: 3,
  CHALLENGE: 4,
}

export const params = {
  dkgResultChallengePeriodLength: 11520,
  dkgResultSubmissionEligibilityDelay: 20,
}
