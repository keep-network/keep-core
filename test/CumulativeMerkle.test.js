const { ethers } = require('hardhat')
const { expect } = require('chai')
const { MerkleTree } = require('merkletreejs')
const fc = require('fast-check')
const keccak256 = require('keccak256')

const { dist } = require('./constants')
const { cumDist } = require('./constants')

function genMerkleLeaf (account, amount) {
  const amountHex = ethers.BigNumber.from(amount).toHexString()
  const leaf = account + amountHex.slice(2).padStart(64, '0')
  return MerkleTree.bufferToHex(keccak256(leaf))
}

describe('Cumulative Merkle Distribution', function () {
  let token

  before(function () {
    // numRuns must be less or equal to the number of accounts in `dist`
    const numRuns = Object.keys(dist.claims).length
    fc.configureGlobal({ numRuns: numRuns, skipEqualValues: true })
  })

  beforeEach(async function () {
    const Token = await ethers.getContractFactory('TokenMock')
    token = await Token.deploy()
  })

  context('when set Merkle Root for first time', async function () {
    let merkleDist

    beforeEach(async function () {
      const MerkleDist = await ethers.getContractFactory('CumulativeMerkleDrop')
      merkleDist = await MerkleDist.deploy(token.address)
    })

    it('should be 0 before setting it up', async function () {
      const contractMerkleRoot = await merkleDist.merkleRoot()
      expect(parseInt(contractMerkleRoot, 16)).to.equal(0)
    })

    it('should be possible to be set a new one', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.hexaString({ minLength: 64, maxLength: 64 }),
          async function (merkleRoot) {
            merkleRoot = '0x' + merkleRoot
            await merkleDist.setMerkleRoot(merkleRoot)
            const contractMerkleRoot = await merkleDist.merkleRoot()
            expect(contractMerkleRoot).to.equal(merkleRoot)
          }
        )
      )
    })

    it('should be emitted an event', async function () {
      const prevMerkleRoot = await merkleDist.merkleRoot()
      const nextMerkleRoot = '0xb2c0cd477fff5f352df19233236e02bac0c4170a783f11cd39589413132914cc'
      const tx = merkleDist.setMerkleRoot(nextMerkleRoot)
      await expect(tx).to.emit(merkleDist, 'MerkelRootUpdated').withArgs(prevMerkleRoot, nextMerkleRoot)
    })

    it("only contract's owner should can change Merkle Root", async function () {
      const [owner, addr1] = await ethers.getSigners()
      const merkleRoot = '0xb2c0cd477fff5f352df19233236e02bac0c4170a783f11cd39589413132914cc'
      await expect(merkleDist.connect(addr1).setMerkleRoot(merkleRoot)).to.be.revertedWith('Ownable: caller is not the owner')
    })
  })

  context('when batch claim array tokens', async function () {
    let merkleDist
    let merkleRoot
    let tokenTotal
    let proofAccounts

    before(function () {
      // numRuns must be less or equal to the number of accounts in `cum_dist`
      merkleRoot = dist.merkleRoot
      tokenTotal = ethers.BigNumber.from(dist.tokenTotal)
      proofAccounts = Object.keys(dist.claims)
    })

    beforeEach(async function () {
      const MerkleDist = await ethers.getContractFactory('CumulativeMerkleDrop')
      merkleDist = await MerkleDist.deploy(token.address)
      await merkleDist.setMerkleRoot(merkleRoot)
      await token.mint(merkleDist.address, tokenTotal)
    })

    it('should accounts get tokens', async function () {
      const claimAmounts = Array.from(proofAccounts).map((claimAccount, _) => ethers.BigNumber.from(dist.claims[claimAccount].amount))
      const claimProofs = Array.from(proofAccounts).map((claimAccount, _) => dist.claims[claimAccount].proof)

      const prevBalances = await Promise.all(proofAccounts.map(async (claimAccount, _) => await token.balanceOf(claimAccount)))
      const expBalances = Array.from(prevBalances).map((prevAmmount, index) =>  parseInt(prevAmmount + claimAmounts[index], 10))
           
      await merkleDist.batchClaimArray(proofAccounts, claimAmounts, merkleRoot, claimProofs)

      const afterBalancesHex = await Promise.all(proofAccounts.map(async (claimAccount, _) => await token.balanceOf(claimAccount)))
      const afterBalances = Array.from(afterBalancesHex).map((afterAmmount, _) => parseInt(afterAmmount["_hex"], 16))
      expBalances.forEach((expAmount, index) => {
        expect(expAmount).to.equal(afterBalances[index])
      })
    })
  })

  context('when batch claim struct tokens', async function () {
    let merkleDist
    let merkleRoot
    let tokenTotal
    let proofAccounts

    before(function () {
      // numRuns must be less or equal to the number of accounts in `cum_dist`
      merkleRoot = dist.merkleRoot
      tokenTotal = ethers.BigNumber.from(dist.tokenTotal)
      proofAccounts = Object.keys(dist.claims)
    })

    beforeEach(async function () {
      const MerkleDist = await ethers.getContractFactory('CumulativeMerkleDrop')
      merkleDist = await MerkleDist.deploy(token.address)
      await merkleDist.setMerkleRoot(merkleRoot)
      await token.mint(merkleDist.address, tokenTotal)
    })

    it('should accounts get tokens', async function () {
      const claimAmounts = Array.from(proofAccounts).map((claimAccount, _) => ethers.BigNumber.from(dist.claims[claimAccount].amount))
      const claimProofs = Array.from(proofAccounts).map((claimAccount, _) => dist.claims[claimAccount].proof)
      
      const claimStructs = Array.from(proofAccounts).map((claimAccount, index) => [claimAccount, claimAmounts[index], claimProofs[index]])

      const prevBalances = await Promise.all(proofAccounts.map(async (claimAccount, _) => await token.balanceOf(claimAccount)))
      const expBalances = Array.from(prevBalances).map((prevAmmount, index) =>  parseInt(prevAmmount + claimAmounts[index], 10))
      
      await merkleDist.batchClaimStruct(merkleRoot, claimStructs)

      const afterBalancesHex = await Promise.all(proofAccounts.map(async (claimAccount, _) => await token.balanceOf(claimAccount)))
      const afterBalances = Array.from(afterBalancesHex).map((afterAmmount, _) => parseInt(afterAmmount["_hex"], 16))
      expBalances.forEach((expAmount, index) => {
        expect(expAmount).to.equal(afterBalances[index])
      })
    })
  })

  context('when claim tokens', async function () {
    let merkleDist
    let merkleRoot
    let tokenTotal
    let proofAccounts

    before(function () {
      // numRuns must be less or equal to the number of accounts in `cum_dist`
      const numRuns = Object.keys(dist.claims).length
      fc.configureGlobal({ numRuns: numRuns, skipEqualValues: true })
      merkleRoot = dist.merkleRoot
      tokenTotal = ethers.BigNumber.from(dist.tokenTotal)
      proofAccounts = Object.keys(dist.claims)
    })

    beforeEach(async function () {
      const MerkleDist = await ethers.getContractFactory('CumulativeMerkleDrop')
      merkleDist = await MerkleDist.deploy(token.address)
      await merkleDist.setMerkleRoot(merkleRoot)
      await token.mint(merkleDist.address, tokenTotal)
    })

    it('should be emitted an event', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          async function (index) {
            const claimAccount = proofAccounts[index]
            const claimAmount = ethers.BigNumber.from(dist.claims[claimAccount].amount)
            const claimProof = dist.claims[claimAccount].proof

            claimAccChecksum = ethers.utils.getAddress(claimAccount)
            const tx = merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)
            await expect(tx).to.emit(merkleDist, 'Claimed').withArgs(claimAccChecksum, claimAmount)
          }
        )
      )
    })
    it('should accounts get tokens', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          async function (index) {
            const claimAccount = proofAccounts[index]
            const claimAmount = ethers.BigNumber.from(dist.claims[claimAccount].amount)
            const claimProof = dist.claims[claimAccount].proof

            const prevBalance = await token.balanceOf(claimAccount)
            const expBalance = prevBalance.add(claimAmount)
            await merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)
            const afterBalance = await token.balanceOf(claimAccount)
            expect(expBalance).to.equal(afterBalance)
          }
        )
      )
    })

    it('should merkle dist contract to reduce its balance', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          async function (index) {
            const claimAccount = proofAccounts[index]
            const claimAmount = ethers.BigNumber.from(dist.claims[claimAccount].amount)
            const claimProof = dist.claims[claimAccount].proof

            preBalance = await token.balanceOf(merkleDist.address)
            expBalance = preBalance.sub(claimAmount)
            await merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)
            const afterBalance = await token.balanceOf(merkleDist.address)
            expect(expBalance).to.equal(afterBalance)
          }
        )
      )
    })

    it('should not be possible to claim for fake accounts', async function () {
      claimAccount = ethers.Wallet.createRandom().address
      claimAmount = 100000
      claimProof = [
        '0xf558bba7dd8aef6fdfb36ea106d965fd7ef483aa217cc02e2c33b78cdfb74cab',
        '0x7a8326f3dfbbddc4a0bc1e3e5005d4cecf6a7c89d386692a27dc5235b55e92cd'
      ]
      await expect(merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)).to.be.revertedWith('CMD: Invalid proof')
    })

    it('should not be possible to claim a different amount of tokens', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          fc.integer({ min: 0, max: 10000000 }),
          async function (index, claimAmount) {
            const claimAccount = proofAccounts[index]
            const claimProof = dist.claims[claimAccount].proof
            await expect(merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)).to.be.revertedWith('CMD: Invalid proof')
          }
        )
      )
    })

    it('should not be possible to claim twice', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          async function (index) {
            const claimAccount = proofAccounts[index]
            const claimAmount = ethers.BigNumber.from(dist.claims[claimAccount].amount)
            const claimProof = dist.claims[claimAccount].proof
            await merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)
            await expect(merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)).to.be.revertedWith('CMD: Nothing to claim')
          }
        )
      )
    })
  })

  context('when set a new Merkle Distribution (cumulative)', async function () {
    let merkleDist
    let merkleRoot
    let cumulativeMerkleRoot
    let tokenTotal
    let cumulativeTokenTotal
    let proofAccounts
    let cumulativeProofAccounts

    before(function () {
      // numRuns must be less or equal to the number of accounts in `cum_dist`
      const numRuns = Object.keys(cumDist.claims).length
      fc.configureGlobal({ numRuns: numRuns, skipEqualValues: true })

      merkleRoot = dist.merkleRoot
      cumulativeMerkleRoot = cumDist.merkleRoot
      tokenTotal = ethers.BigNumber.from(dist.tokenTotal)
      cumulativeTokenTotal = ethers.BigNumber.from(cumDist.tokenTotal)
      proofAccounts = Object.keys(dist.claims)
      cumulativeProofAccounts = Object.keys(cumDist.claims)
    })

    beforeEach(async function () {
      const MerkleDist = await ethers.getContractFactory('CumulativeMerkleDrop')
      merkleDist = await MerkleDist.deploy(token.address)
      await merkleDist.setMerkleRoot(merkleRoot)
      await token.mint(merkleDist.address, tokenTotal)
    })

    it('should be possible to set a new Merkle Root after claiming', async function () {
      const claimAccount = proofAccounts[0]
      const claimAmount = ethers.BigNumber.from(dist.claims[claimAccount].amount)
      const claimProof = dist.claims[claimAccount].proof

      await merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)
      await merkleDist.setMerkleRoot(cumulativeMerkleRoot)
      const contractMerkleRoot = await merkleDist.merkleRoot()
      expect(contractMerkleRoot).to.equal(cumulativeMerkleRoot)
    })

    it('should not be possible to claim using old Merkle Root', async function () {
      const claimAccount = proofAccounts[0]
      const claimAmount = ethers.BigNumber.from(dist.claims[claimAccount].amount)
      const claimProof = dist.claims[claimAccount].proof
      await merkleDist.setMerkleRoot(cumulativeMerkleRoot)

      await expect(merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)).to.be.revertedWith('CMD: Merkle root was updated')
    })

    context('after claiming all tokens of the previous distribution', async function () {
      beforeEach(async function () {
        for (claimAccount of proofAccounts) {
          const claimAmount = ethers.BigNumber.from(dist.claims[claimAccount].amount)
          const claimProof = dist.claims[claimAccount].proof
          await merkleDist.claim(claimAccount, claimAmount, merkleRoot, claimProof)
        }
      })

      it('should not be possible to claim without enough balance in contract', async function () {
        await merkleDist.setMerkleRoot(cumulativeMerkleRoot)

        const claimAccount = cumulativeProofAccounts[0]
        const claimAmount = ethers.BigNumber.from(cumDist.claims[claimAccount].amount)
        const claimProof = cumDist.claims[claimAccount].proof

        await expect(merkleDist.claim(claimAccount, claimAmount, cumulativeMerkleRoot, claimProof)).to.be.revertedWith('Transfer amount exceeds balance')
      })

      it('should be possible to claim new distribution tokens', async function () {
        await token.mint(merkleDist.address, cumulativeTokenTotal)
        await merkleDist.setMerkleRoot(cumulativeMerkleRoot)

        await fc.assert(
          fc.asyncProperty(
            fc.integer({ min: 0, max: cumulativeProofAccounts.length - 1 }),
            async function (index) {
              const claimAccount = cumulativeProofAccounts[index]
              const claimAmount = ethers.BigNumber.from(cumDist.claims[claimAccount].amount)
              const claimProof = cumDist.claims[claimAccount].proof

              await merkleDist.claim(claimAccount, claimAmount, cumulativeMerkleRoot, claimProof)
              const balance = await token.balanceOf(claimAccount)
              expect(claimAmount).to.equal(balance)
            }
          )
        )
      })
    })
  })

  context('when verify Merkle Proof', async function () {
    let merkleDist
    let merkleRoot
    let proofAccounts

    before(function () {
      // numRuns must be less or equal to the number of accounts in `cum_dist`
      const numRuns = Object.keys(dist.claims).length
      fc.configureGlobal({ numRuns: numRuns, skipEqualValues: true })
      merkleRoot = dist.merkleRoot
      proofAccounts = Object.keys(dist.claims)
    })

    beforeEach(async function () {
      const MerkleDist = await ethers.getContractFactory('CumulativeMerkleDrop')
      merkleDist = await MerkleDist.deploy(token.address)
      await merkleDist.setMerkleRoot(merkleRoot)
    })

    it('should not be verified if no Merkle Proof', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          async function (index) {
            const account = proofAccounts[index]
            const amount = dist.claims[account].amount
            const leaf = genMerkleLeaf(account, amount)
            const claimProof = []
            const verif = await merkleDist.verify(claimProof, merkleRoot, leaf)
            expect(verif).to.be.false
          }
        )
      )
    })

    it('should not be verified with incorrect Merkle Proof', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          async function (index) {
            const account = proofAccounts[index]
            const amount = dist.claims[account].amount
            const leaf = genMerkleLeaf(account, amount)
            const claimProof = [
              MerkleTree.bufferToHex(keccak256('proof1')),
              MerkleTree.bufferToHex(keccak256('proof2'))
            ]
            const verif = await merkleDist.verify(claimProof, merkleRoot, leaf)
            expect(verif).to.be.false
          }
        )
      )
    })

    it('should a correct MerkleProof be verified', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          async function (index) {
            const account = proofAccounts[index]
            const amount = dist.claims[account].amount
            const leaf = genMerkleLeaf(account, amount)
            const claimProof = dist.claims[account].proof
            const verif = await merkleDist.verify(claimProof, merkleRoot, leaf)
            expect(verif).to.be.true
          }
        )
      )
    })

    it('should not be verified a Merkle Proof with incorrect root', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          fc.hexaString({ minLength: 64, maxLength: 64 }),
          async function (index, root) {
            root = '0x' + root
            const account = proofAccounts[index]
            const amount = dist.claims[account].amount
            const leaf = genMerkleLeaf(account, amount)
            const claimProof = dist.claims[account].proof
            const verif = await merkleDist.verify(claimProof, root, leaf)
            expect(verif).to.be.false
          }
        )
      )
    })

    it('should not be verified a Merkle Proof with incorrect amount', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          fc.integer({ min: 1, max: 1000000000 }),
          async function (index, amount) {
            const account = proofAccounts[index]
            const leaf = genMerkleLeaf(account, amount)
            const claimProof = dist.claims[account].proof
            const verif = await merkleDist.verify(claimProof, merkleRoot, leaf)
            expect(verif).to.be.false
          }
        )
      )
    })

    it('should not be verified a Merkle Proof with incorrect account', async function () {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 0, max: proofAccounts.length - 1 }),
          async function (index) {
            const fakeAccount = ethers.Wallet.createRandom().address
            const account = proofAccounts[index]
            const amount = dist.claims[account].amount
            const leaf = genMerkleLeaf(fakeAccount, amount)
            const claimProof = dist.claims[account].proof
            const verif = await merkleDist.verify(claimProof, merkleRoot, leaf)
            expect(verif).to.be.false
          }
        )
      )
    })
  })
})
