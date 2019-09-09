import { sign } from './signature';

export default async function stakeDelegate(stakingContract, token, owner, operator, magpie, stake) {
  let signature = await sign(web3.utils.soliditySha3(owner), operator);

  let delegation = Buffer.concat([
    Buffer.from(magpie.substr(2), 'hex'),
    Buffer.from(signature.substr(2), 'hex')
  ]);
  token.approveAndCall(stakingContract.address, stake, delegation, {from: owner});
}