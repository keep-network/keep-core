const hre = require('hardhat');
const { getChainId } = hre;

module.exports = async ({ deployments, getNamedAccounts }) => {
    console.log('running deploy script');
    console.log('network id ', await getChainId());

    const { deploy } = deployments;
    const { deployer } = await getNamedAccounts();

    const args = ['0x111111111117dC0aa78b770fA6A738034120C302', '0x8d75e829b65729ab44ed5117daec9c21', 10];

    const merkleDrop128 = await deploy('MerkleDrop128', {
        from: deployer,
        args: args,
        skipIfAlreadyDeployed: false,
        maxFeePerGas: 100e9,
        maxPriorityFeePerGas: 2e9,
    });

    console.log('MerkleDrop128 deployed to:', merkleDrop128.address);

    if (await getChainId() !== '31337') {
        await hre.run('verify:verify', {
            address: merkleDrop128.address,
            constructorArguments: args,
        });
    }
};

module.exports.skip = async () => true;
