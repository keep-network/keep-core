module.exports = async ({
  getNamedAccounts,
  deployments,
  getChainId,
  getUnnamedAccounts,
}) => {
  const {deploy} = deployments
  const {deployer} = await getNamedAccounts()

  const merkleRoot = '0x591189966065882f5d6222e97d91793524b7464c039fae114ae03d1c2fad5ef1'

  const token = await deployments.get("TokenMock")

  const merkleDistDeployment = await deploy('CumulativeMerkleDrop', {
    from: deployer,
    args: [
      token.address,
    ],
  })

  const MerkleDist = await ethers.getContractFactory('CumulativeMerkleDrop')
  const merkleDist = MerkleDist.attach(merkleDistDeployment.address)

  const setMerkleRootTx = await merkleDist.setMerkleRoot(merkleRoot)
  setMerkleRootTx.wait()

  console.log('Merkle Distribution deployed to:', merkleDist.address)
}

module.exports.tags = ["TMerkleDistributor"]
