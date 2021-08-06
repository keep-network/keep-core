const hre = require('hardhat');
const { getChainId } = hre;

module.exports = async ({ deployments, getNamedAccounts }) => {
    console.log('running deploy script');
    console.log('network id ', await getChainId());

    const { deploy } = deployments;
    const { deployer } = await getNamedAccounts();

    const example = await deploy('Example', {
        from: deployer,
    });

    console.log('Example deployed to:', example.address);

    if (await getChainId() !== '31337') {
        await hre.run('verify:verify', {
            address: example.address,
        });
    }
};

module.exports.skip = async () => true;
