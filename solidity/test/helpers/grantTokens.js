export default async function grantTokens(
    grantContract,
    token, amount,
    from, grantee,
    vestingDuration, start, cliff,
    revocable) {
  let grantData = Buffer.concat([
    Buffer.from(grantee.substr(2), 'hex'),
    web3.utils.toBN(vestingDuration).toBuffer('be', 32),
    web3.utils.toBN(start).toBuffer('be', 32),
    web3.utils.toBN(cliff).toBuffer('be', 32),
    Buffer.from(revocable ? "01" : "00", 'hex'),
  ]);

  await token.approveAndCall(grantContract.address, amount, grantData, {from: from})
  return (await grantContract.getPastEvents())[0].args[0].toNumber()
}
