const hre = require('hardhat');
const { getChainId, ethers } = hre;

module.exports = async ({ deployments, getNamedAccounts }) => {
    console.log('running mainnet test deployment script');
    console.log('network id ', await getChainId());

    const {deploy} = deployments;
    const {deployer, _, rewardsHolder, owner} = await getNamedAccounts();

    const tokenContract = await ethers.getContractFactory('TokenMock')
    const token = await tokenContract.deploy()
    console.log('Token mock deployed to:', token.address)
    const tokenTotal = '39223394181720035083539958'
    const mintTx = await token.mint(rewardsHolder, tokenTotal);
    await mintTx.wait()
    const minted = (await token.totalSupply()).toString()
    // console.log('Token minted:', minted)
    
    // const args = [token.address, rewardsHolder, owner];
    // const merkleDrop = await deploy('CumulativeMerkleDrop', {
    //     from: deployer,
    //     args,
    // });

    // console.log('CumulativeMerkleDrop deployed to:', merkleDrop.address);

    if (await getChainId() !== '31337') {
        await hre.run('verify:verify', {
            address: token.address,
            constructorArguments: args,
        });
    }

};