const { ethers } = require("hardhat")

module.exports = async ({
  getNamedAccounts,
  deployments,
  getChainId,
  getUnnamedAccounts,
}) => {
  const { deploy, execute } = deployments
  const { deployer } = await getNamedAccounts()

  const totalAmount = '67776031738882'

  const tokenDeployment = await deployments.get('TokenMock')
  const merkleDist = await deployments.get('CumulativeMerkleDrop')

  const Token = await ethers.getContractFactory('TokenMock')
  const token = await Token.attach(tokenDeployment.address)

  const mintTx = await token.mint(merkleDist.address, totalAmount)
  await mintTx.wait()

  const minted = (await token.totalSupply()).toString()

  console.log('Token minted:', minted)
}

module.exports.tags = ["MintToken"]
