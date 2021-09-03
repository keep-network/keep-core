const hre = require('hardhat');
const { getChainId, ethers } = hre;

module.exports = async ({ deployments, getNamedAccounts }) => {
    console.log('running deploy script');
    console.log('network id ', await getChainId());

    const { deploy } = deployments;
    const { deployer } = await getNamedAccounts();

    const args = ['0x111111111117dC0aa78b770fA6A738034120C302'];

    const merkleDrop = await deploy('CumulativeMerkleDrop', {
        from: deployer,
        args: args,
        skipIfAlreadyDeployed: true,
        maxFeePerGas: 100000000000,
        maxPriorityFeePerGas: 2000000000,
    });

    const CumulativeMerkleDrop = await ethers.getContractFactory('CumulativeMerkleDrop');
    const cumulativeMerkleDrop = CumulativeMerkleDrop.attach(merkleDrop.address);

    const txn = await cumulativeMerkleDrop.setMerkleRoot(
        '0x323e1a13446c2a6ed35c700e5f336cdd367b554f76fd7e8268eb3a302e963924',
        {
            maxFeePerGas: 100000000000,
            maxPriorityFeePerGas: 2000000000,
        }
    );
    await txn;

    console.log('CumulativeMerkleDrop deployed to:', merkleDrop.address);

    if (await getChainId() !== '31337') {
        await hre.run('verify:verify', {
            address: merkleDrop.address,
            constructorArguments: args,
        });
    }
};

module.exports.skip = async () => true;
