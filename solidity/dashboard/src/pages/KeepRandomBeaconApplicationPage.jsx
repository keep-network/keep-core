import React from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"
import AuthorizationHistory from "../components/AuthorizationHistory"
import * as Icons from "../components/Icons"

const KeepRandomBeaconApplicationPage = () => {
  const data = [
    {
      operatorAddress: "address",
      stakeAmount: "1000",
      contracts: [
        {
          contractName: "Keep Random Beacon Operator Contract",
          operatorContractAddress: "address",
          isAuthorized: false,
        },
      ],
    },
  ]

  return (
    <PageWrapper
      className=""
      title="Random Beacon"
      nextPageLink="/rewards"
      nextPageTitle="Rewards"
      nextPageIcon={Icons.KeepBlackGreen}
    >
      <nav className="mb-2">
        <a
          href="https://keep.network/"
          className="arrow-link h4"
          rel="noopener noreferrer"
          target="_blank"
        >
          Keep Website
        </a>
      </nav>
      <AuthorizeContracts data={data} />
      <AuthorizationHistory contracts={data} />
    </PageWrapper>
  )
}

export default KeepRandomBeaconApplicationPage
