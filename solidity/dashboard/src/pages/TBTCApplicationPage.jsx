import React from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"

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
    <PageWrapper title="tBTC">
      <AuthorizeContracts data={data} />
    </PageWrapper>
  )
}

export default TBTCApplicationPage
