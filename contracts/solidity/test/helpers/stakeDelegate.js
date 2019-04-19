export default async function stakeDelegate(stakingContract, token, owner, operator, magpie,stake) {
    let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator)).substr(2), 'hex');
    let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');
    token.approveAndCall(stakingContract.address, stake, delegation, {from: owner});
}