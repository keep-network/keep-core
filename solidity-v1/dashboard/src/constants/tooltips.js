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
  claimTokensFromUndelegation: {
    title: "Claim Tokens",
    content:
      "Click claim to return the undelegated tokens to your token balance.",
  },
  beaconEarnings: {
    title: "Beacon Earnings",
    content:
      "The total balance reflects the total Available and Active earnings. Available earningss are ready to be withdrawn. Active earnings become available after a signing group expires.",
    withRedirectLink: false,
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
  pendingWithdrawal: {
    title: "Pending Withdrawal",
    content:
      "After the 21 day cooldown ends, you have a 2 day claim window to claim your tokens and rewards. Your deposit and rewards will be sent in one transaction. If you do not claim your tokens within 2 days, your tokens will return to the pool and you will have to withdraw them again.",
    linkText: "How it works",
    redirectLink: "/coverage-pools/how-it-works",
  },
  totalValueLocked: {
    title: "Total Value Locked",
    content: "The total amount of KEEP deposited into the coverage pool.",
    withRedirectLink: true,
    redirectLink: "/coverage-pools/how-it-works",
    linkText: "How it works",
  },
  covPoolsDeposit: {
    title: "Coverage pool deposit",
    content:
      'Deposit into the coverage pool to secure the network and earn rewards. A coverage pool functions as a form of insurance that can be used as a back-stop or "buyer of last resort" in on-chain financial systems.',
    withRedirectLink: true,
    redirectLink: "/coverage-pools/how-it-works",
    linkText: "How it works",
  },
  covPoolsAvailableToWithdraw: {
    title: "Available to withdraw",
    content:
      "The amount of KEEP you have available to withdraw from the coverage pool. Note that there is a 21 day cooldown period before you can claim your tokens after you withdraw.",
    withRedirectLink: true,
    redirectLink: "/coverage-pools/how-it-works",
    linkText: "How it works",
  },
  thresholdPageGrantAllocation: {
    title: "Grant Allocation",
    content: "A grant is something that vests KEEP over a set period of time",
    withRedirectLink: false,
  },
}

export default resourceTooltipProps
