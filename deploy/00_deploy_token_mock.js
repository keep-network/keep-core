module.exports = async ({
  getNamedAccounts,
  deployments,
  getChainId,
  getUnnamedAccounts,
}) => {
  const {deploy} = deployments
  const {deployer} = await getNamedAccounts()

  const token = await deploy('TokenMock', {
    from: deployer,
    args: [],
  })

  console.log('Token mock deployed to:', token.address)
}

module.exports.tags = ["TokenMock"]
