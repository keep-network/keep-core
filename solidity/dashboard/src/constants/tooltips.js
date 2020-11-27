const resourceTooltipProps = {
  delegation: {
    title: "Delegation",
    content:
      "Delegation sets aside an amount of KEEP to be staked by a trusted third party, referred to within the dApp as an operator.",
  },
  cliff: {
    title: "Cliff",
    content: "A cliff is a set period of time before vesting begins.",
  },
  recoverTokens: {
    title: "Recover Tokens",
    content:
      "Click recover to return undelegated tokens to your granted token balance.",
  },
  beaconEarnings: {
    title: "Beacon Earnings",
    content:
      "The total balance reflects the total Available and Active earnings. Available earningss are ready to be withdrawn. Active earnings become available after a signing group expires.",
    withRedirectButton: false,
  },
  slashing: {
    title: "Slashing",
    content:
      "A slash is a penalty for signing group misbehavior. It results in a removal of a portion of your delegated KEEP tokens.",
  },
  authorize: {
    title: "Authorize",
    content:
      "By authorizing a contract, you are approving a set of terms for the governance of an operator, e.g. the rules for slashing tokens.",
  },
  tokenGrant: {
    title: "Token Grant",
    content:
      "A grant that contains KEEP tokens that unlocks at a set schedule over a period of time.",
  },
}

export default resourceTooltipProps
