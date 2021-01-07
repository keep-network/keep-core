import React from "react"

const DepositLPTokensMsgContent = ({ liquidityRewardPair }) => {
  return (
    <a
      href={liquidityRewardPair.viewPoolLink}
      rel="noopener noreferrer"
      target="_blank"
    >
      Deposit them and earn rewards
    </a>
  )
}

export default React.memo(DepositLPTokensMsgContent)