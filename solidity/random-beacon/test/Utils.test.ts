import { ethers } from "hardhat"
import { expect } from "chai"
import { hashDKGMembers } from "./utils/dkg"

describe("utils/dkg", () => {
  context("hashDKGMembers", () => {
    const members = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

    context("when there are no misbehaved members", () => {
      const misbehavedMembers = []
      const expectedMembers = members

      const actualHash = hashDKGMembers(members, misbehavedMembers)
      const expectedHash = ethers.utils.keccak256(
        ethers.utils.defaultAbiCoder.encode(["uint32[]"], [expectedMembers])
      )

      expect(expectedHash).to.be.equal(actualHash)
    })

    context("when the first member id is a misbehaved one", () => {
      // misbehaved member count starts from 1, not 0
      const misbehavedMembers = [1]
      // expectedMembers[misbehavedMembers - 1], hence removing "0"
      const expectedMembers = [1, 2, 3, 4, 5, 6, 7, 8, 9]

      const actualHash = hashDKGMembers(members, misbehavedMembers)
      const expectedHash = ethers.utils.keccak256(
        ethers.utils.defaultAbiCoder.encode(["uint32[]"], [expectedMembers])
      )

      expect(expectedHash).to.be.equal(actualHash)
    })

    context("when the last member id is a misbehaved", () => {
      // misbehaved member count starts from 1, not 0
      const misbehavedMembers = [10]
      // expectedMembers[misbehavedMembers - 1], hence removing "9"
      const expectedMembers = [0, 1, 2, 3, 4, 5, 6, 7, 8]

      const actualHash = hashDKGMembers(members, misbehavedMembers)
      const expectedHash = ethers.utils.keccak256(
        ethers.utils.defaultAbiCoder.encode(["uint32[]"], [expectedMembers])
      )

      expect(expectedHash).to.be.equal(actualHash)
    })

    context("when there are multiple misbehaved members", () => {
      // misbehaved member count starts from 1, not 0
      const misbehavedMembers = [2, 6, 8]
      // expectedMembers[misbehavedMembers - 1], hence removing "1,5,7"
      const expectedMembers = [0, 2, 3, 4, 6, 8, 9]

      const actualHash = hashDKGMembers(members, misbehavedMembers)
      const expectedHash = ethers.utils.keccak256(
        ethers.utils.defaultAbiCoder.encode(["uint32[]"], [expectedMembers])
      )

      expect(expectedHash).to.be.equal(actualHash)
    })

    context("when a misbehaved member is not present in members array", () => {
      const misbehavedMembers = [42]
      const expectedMembers = members

      const actualHash = hashDKGMembers(members, misbehavedMembers)
      const expectedHash = ethers.utils.keccak256(
        ethers.utils.defaultAbiCoder.encode(["uint32[]"], [expectedMembers])
      )

      expect(expectedHash).to.be.equal(actualHash)
    })
  })
})
