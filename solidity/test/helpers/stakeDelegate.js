async function stakeDelegate(stakingContract, token, owner, operator, beneficiary, authorizer, stake) {
  let delegation = Buffer.concat([
    Buffer.from(beneficiary.substr(2), 'hex'),
    Buffer.from(operator.substr(2), 'hex'),
    Buffer.from(authorizer.substr(2), 'hex')
  ]);
  return token.approveAndCall(stakingContract.address, stake, delegation, {from: owner});
}
module.exports = stakeDelegate
