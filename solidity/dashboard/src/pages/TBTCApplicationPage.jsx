import React from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"
import { tbtcAuthorizationService } from "../services/tbtc-authorization.service"
import { useFetchData } from "../hooks/useFetchData"

const TBTCApplicationPage = () => {
  // fetch data from service
  const initialData = {}
  const [state] = useFetchData(
    tbtcAuthorizationService.fetchTBTCAuthorizationData,
    initialData
  )

  console.log("state:", state.data)

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
