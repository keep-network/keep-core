import { sign } from './signature';

export default async function stakeDelegate(stakingContract, token, owner, operator, magpie, stake) {
  let delegation = Buffer.concat([
    Buffer.from(magpie.substr(2), 'hex'),
    Buffer.from(operator.substr(2), 'hex')
  ]);
  token.approveAndCall(stakingContract.address, stake, delegation, {from: owner});
}