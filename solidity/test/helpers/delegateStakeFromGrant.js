export default async function delegateStakeFromGrant(
    grantContract,
    stakingContractAddress,
    grantee,
    operator,
    magpie,
    authorizer,
    amount,
    grantId
) {
    let delegation = Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operator.substr(2), 'hex'),
      Buffer.from(authorizer.substr(2), 'hex')
    ]);

    return grantContract.stake(
      grantId, 
      stakingContractAddress, 
      amount, 
      delegation, 
      {from: grantee}
    );
  }