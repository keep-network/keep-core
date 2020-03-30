import {createSnapshot, restoreSnapshot} from '../helpers/snapshot';

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const LockStub = artifacts.require('./stubs/LockStub.sol');

contract.only('LockStub', (accounts) => {
  let locks;

  const alice = accounts[0];
  const bob = accounts[1];
  const carol = accounts[2];

  before(async () => {
      locks = await LockStub.new();
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("sets locks", async () => {
    expect(await locks.publicContains(alice)).to.be.false;
    await locks.publicSetLock(alice, 1);
    expect(await locks.publicContains(alice)).to.be.true;
  })

  it("releases locks", async () => {
    await locks.publicSetLock(alice, 1);
    expect(await locks.publicContains(alice)).to.be.true;
    await locks.publicReleaseLock(alice);
    expect(await locks.publicContains(alice)).to.be.false;
  })

  it("overwrites existing locks", async () => {
    await locks.publicSetLock(alice, 2);
    await locks.publicSetLock(bob, 2);
    expect(await locks.publicGetLockTime(alice)).to.eq.BN(2);
    expect(await locks.publicGetLockTime(bob)).to.eq.BN(2);
    await locks.publicSetLock(alice, 1);
    await locks.publicSetLock(bob, 3);
    expect(await locks.publicGetLockTime(alice)).to.eq.BN(1);
    expect(await locks.publicGetLockTime(bob)).to.eq.BN(3);
  })

  it("reorders locks correctly when deleting", async () => {
    await locks.publicSetLock(alice, 1);
    await locks.publicSetLock(bob, 2);
    await locks.publicSetLock(carol, 3);
    let creatorsPre = await locks.publicEnumerateCreators();
    expect(creatorsPre[1]).to.equal(bob);

    await locks.publicReleaseLock(bob);
    let creatorsPost = await locks.publicEnumerateCreators();
    expect(creatorsPost[1]).to.equal(carol);
  })
})
