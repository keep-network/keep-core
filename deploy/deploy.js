const hre = require('hardhat');
const { getChainId, ethers } = hre;

module.exports = async ({ deployments, getNamedAccounts }) => {
    console.log('running deploy script');
    console.log('network id ', await getChainId());

    const { deploy } = deployments;
    const { deployer } = await getNamedAccounts();

    const args = ['0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2'];

    const merkleDrop = await deploy('CumulativeMerkleDrop', {
        from: deployer,
        args: args,
        skipIfAlreadyDeployed: true,
    });

    const CumulativeMerkleDrop = await ethers.getContractFactory('CumulativeMerkleDrop');
    const cumulativeMerkleDrop = CumulativeMerkleDrop.attach(merkleDrop.address);
    await cumulativeMerkleDrop.setMerkleRoot('0xd76ea6876293c58ef1fa269fed8274d9784195fd3c2a12dd0bd35b5729d24f76');

    console.log('CumulativeMerkleDrop deployed to:', merkleDrop.address);

    if (await getChainId() !== '31337') {
        await hre.run('verify:verify', {
            address: merkleDrop.address,
            constructorArguments: args,
        });
    }
};

module.exports.skip = async () => true;
