export default async function grantTokens(grantContract, token, amount, from, beneficiary, vestingDuration, start, cliff, revocable) {
  await token.approve(grantContract.address, amount, {from: from});
  return await grantContract.grant(amount, beneficiary, vestingDuration,
    start, cliff, revocable, {from: from}).then((result)=>{
    // Look for CreatedTokenGrant event in transaction receipt and get grant id
    for (var i = 0; i < result.logs.length; i++) {
      var log = result.logs[i];
      if (log.event == "CreatedTokenGrant") {
        return log.args.id.toNumber();
      }
    }
  })
}
