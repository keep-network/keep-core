const {accounts, contract, web3} = require('@openzeppelin/test-environment')
const {expectRevert} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot')

const assert = require('chai').assert
const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const GrantStaking = contract.fromArtifact('GrantStaking')
const GrantStakingStub = contract.fromArtifact('GrantStakingStub')

describe('GrantStaking', () => {
  
  const deployer = accounts[0],
    operator1 = accounts[1]
    operator2 = accounts[2]

  let info

  before(async () => {
    const infoLib = await GrantStaking.new({from: deployer})
    await GrantStakingStub.detectNetwork()
    await GrantStakingStub.link('GrantStaking', infoLib.address)
    info = await GrantStakingStub.new({from: deployer})
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe('hasGrantDelegated', async () => {
    it('returns true if grant is set for operator', async () => {
      await info.setGrantForOperator(operator1, 0)
      await info.setGrantForOperator(operator2, 1)

      assert.isTrue(await info.hasGrantDelegated(operator1))
      assert.isTrue(await info.hasGrantDelegated(operator2))
    })

    it('returns false if grant is not set for operator', async () => {        
      assert.isFalse(await info.hasGrantDelegated(operator1))
      await info.setGrantForOperator(operator2, 0)
      assert.isFalse(await info.hasGrantDelegated(operator1))
    })
  })

  describe('getGrantForOperator', async () => {
    it('returns grant ID for operator having grant staked', async () => {
      await info.setGrantForOperator(operator1, 0)
      await info.setGrantForOperator(operator2, 10)

      const operator1GrantId = await info.getGrantForOperator(operator1)
      const operator2GrantId = await info.getGrantForOperator(operator2)

      expect(operator1GrantId).to.eq.BN(0)
      expect(operator2GrantId).to.eq.BN(10)
    })

    it('reverts where there is no grant for operator', async () => {
      await expectRevert(
        info.getGrantForOperator(operator1),
        "No grant for the operator"
      )
    })
  })
})