import React from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"
import * as Icons from "../components/Icons"

const TBTCApplicationPage = () => {
  // fetch data from service
  const data = [
    {
      operatorAddress: "address",
      stakeAmount: "1000",
      contracts: [
        {
          contractName: "BondedECDSAKeepFactory",
          operatorContractAddress: "address",
          isAuthorized: false,
        },
        {
          contractName: "TBTCSystem",
          operatorContractAddress: "address",
          isAuthorized: true,
        },
      ],
    },
  ]

  return (
    <PageWrapper
      className=""
      title="tBTC"
      nextPageLink="/rewards"
      nextPageTitle="Rewards"
      nextPageIcon={Icons.TBTC}
    >
      <nav className="mb-2">
        <a
          href="https://tbtc.network/"
          className="arrow-link h4"
          rel="noopener noreferrer"
          target="_blank"
        >
          tBTC Website
        </a>
      </nav>
      <AuthorizeContracts data={data} />
    </PageWrapper>
  )
}

export default TBTCApplicationPage
