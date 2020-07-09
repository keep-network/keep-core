const {contract, accounts, web3} = require("@openzeppelin/test-environment");
const {createSnapshot, restoreSnapshot} = require('../helpers/snapshot.js');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const LockStub = contract.fromArtifact('LockStub');

describe('LockUtils', () => {
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
    let creatorsPre = await locks.publicEnumerateCreators();
    expect(creatorsPre[0]).to.equal(alice);
    expect(creatorsPre[1]).to.equal(bob);
    expect(await locks.publicGetLockTime(alice)).to.eq.BN(2);
    expect(await locks.publicGetLockTime(bob)).to.eq.BN(2);

    await locks.publicSetLock(alice, 1);
    await locks.publicSetLock(bob, 3);
    let creatorsPost = await locks.publicEnumerateCreators();
    expect(creatorsPost[0]).to.equal(alice);
    expect(creatorsPost[1]).to.equal(bob);
    expect(await locks.publicGetLockTime(alice)).to.eq.BN(1);
    expect(await locks.publicGetLockTime(bob)).to.eq.BN(3);
  })

  it("reorders locks correctly when deleting", async () => {
    await locks.publicSetLock(alice, 1);
    await locks.publicSetLock(bob, 2);
    await locks.publicSetLock(carol, 3);
    let creatorsPre = await locks.publicEnumerateCreators();
    expect(creatorsPre[0]).to.equal(alice);
    expect(creatorsPre[1]).to.equal(bob);
    expect(creatorsPre[2]).to.equal(carol);

    await locks.publicReleaseLock(bob);
    let creatorsPost = await locks.publicEnumerateCreators();
    expect(creatorsPost[0]).to.equal(alice);
    expect(creatorsPost[1]).to.equal(carol);

    await locks.publicSetLock(bob, 2);
    let creatorsDoublePost = await locks.publicEnumerateCreators();
    expect(creatorsDoublePost[0]).to.equal(alice);
    expect(creatorsDoublePost[1]).to.equal(carol);
    expect(creatorsDoublePost[2]).to.equal(bob);

    await locks.publicReleaseLock(alice);
    let creatorsTriplePost = await locks.publicEnumerateCreators();
    expect(creatorsTriplePost[0]).to.equal(bob);
    expect(creatorsTriplePost[1]).to.equal(carol);

    await locks.publicReleaseLock(bob);
    let creatorsQuadPost = await locks.publicEnumerateCreators();
    expect(creatorsQuadPost[0]).to.equal(carol);

    await locks.publicReleaseLock(carol);
    let creatorsPentaPost = await locks.publicEnumerateCreators();
    expect(creatorsPentaPost.length).to.equal(0);
  })
})
