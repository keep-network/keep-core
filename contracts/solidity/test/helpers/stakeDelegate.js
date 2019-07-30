import { sign } from './signature';

export default async function stakeDelegate(stakingContract, token, owner, operator, magpie, stake) {
  let signature = Buffer.from(await sign(owner, operator), 'hex');

  let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');
  token.approveAndCall(stakingContract.address, stake, delegation, {from: owner});
}