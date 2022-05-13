/* eslint-disable @typescript-eslint/no-unused-expressions */
import env, { ethers } from "hardhat"
import { expect } from "chai"

import { getNamedSigners, getUnnamedSigners } from "../../utils/signers"

import type { HardhatNetworkHDAccountsUserConfig } from "hardhat/types"

describe("signers", () => {
  describe("getNamedSigners", () => {
    it("should return named signers", async () => {
      const { deployer, governance } = await getNamedSigners()

      expect(ethers.utils.isAddress(deployer.address)).to.be.true

      expect(deployer.address).to.be.not.equal(ethers.constants.AddressZero)

      expect(ethers.utils.isAddress(governance.address)).to.be.true

      expect(governance.address).to.be.not.equal(ethers.constants.AddressZero)

      expect(governance.address).to.be.not.equal(deployer.address)
    })
  })

  describe("getUnnamedSigners", () => {
    it("should return unnamed signers", async () => {
      const namedSigners = await getNamedSigners()
      const unnamedSigners = await getUnnamedSigners()

      expect(unnamedSigners).to.have.length(
        (env.network.config.accounts as HardhatNetworkHDAccountsUserConfig)
          .count - Object.keys(namedSigners).length
      )

      // eslint-disable-next-line no-restricted-syntax
      for (const namedSigner of Object.values(namedSigners)) {
        expect(unnamedSigners).to.not.contain(namedSigner.address)
      }
    })
  })
})
