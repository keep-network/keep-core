async function delegateStake(
  tokenContract, 
  stakingContract, 
  tokenOwner,
  operator,
  beneficiary,
  authorizer,
  amount
) {
  let data = Buffer.concat([
    Buffer.from(beneficiary.substr(2), 'hex'),
    Buffer.from(operator.substr(2), 'hex'),
    Buffer.from(authorizer.substr(2), 'hex')
  ]);
    
  return tokenContract.approveAndCall(
    stakingContract.address, amount, 
    '0x' + data.toString('hex'), 
    {from: tokenOwner}
  );
}

async function delegateStakeFromGrant(
  grantContract,
  stakingContractAddress,
  grantee,
  operator,
  beneficiary,
  authorizer,
  amount,
  grantId
) {
  let delegation = Buffer.concat([
    Buffer.from(beneficiary.substr(2), 'hex'),
    Buffer.from(operator.substr(2), 'hex'),
    Buffer.from(authorizer.substr(2), 'hex')
  ])

  return grantContract.stake(
    grantId, 
    stakingContractAddress, 
    amount, 
    delegation, 
    {from: grantee}
  )
}

async function delegateStakeFromManagedGrant(
  managedGrant,
  stakingContractAddress,
  grantee,
  operator,
  beneficiary,
  authorizer,
  amount
) {
  let delegation = Buffer.concat([
    Buffer.from(beneficiary.substr(2), 'hex'),
    Buffer.from(operator.substr(2), 'hex'),
    Buffer.from(authorizer.substr(2), 'hex')
  ])

  return managedGrant.stake(
    stakingContractAddress, 
    amount, 
    delegation, 
    {from: grantee}
  )
}

module.exports.delegateStake = delegateStake
module.exports.delegateStakeFromGrant = delegateStakeFromGrant
module.exports.delegateStakeFromManagedGrant = delegateStakeFromManagedGrant
