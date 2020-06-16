import React, { useEffect } from "react"
import { useParams } from "react-router-dom"
import { useFetchData } from "../hooks/useFetchData"
import PageWrapper from "../components/PageWrapper"
import { tokenGrantsService } from "../services/token-grants.service"
import { LoadingOverlay } from "../components/Loadable"
import { Web3Context } from "../components/WithWeb3Context"
import TokenGrant from "@keep-network/keep-core/artifacts/TokenGrant.json"
import { getContractAddress } from "../contracts"
import Web3 from "web3"
import { getWsUrl } from "../connectors/utils"
import {
  TokenGrantDetails,
  TokenGrantStakedDetails,
  TokenGrantUnlockingdDetails,
} from "../components/TokenGrantOverview"

const TokenGrantPreview = () => {
  const { grantId } = useParams()
  const [state] = useFetchData(tokenGrantsService.fetchGrantById, {}, grantId)

  return (
    <LoadingOverlay isFetching={state.isFetching}>
      <PageWrapper title={`Grant ID ${grantId}`}>
        <section className="tile token-grant-overview">
          <div className="grant-amount">
            <TokenGrantDetails
              title="Grant Amount"
              selectedGrant={state.data}
            />
          </div>
          <div className="unlocking-details">
            <TokenGrantUnlockingdDetails
              selectedGrant={state.data}
              hideReleaseTokensBtn
            />
          </div>
          <div className="staked-details">
            <TokenGrantStakedDetails
              selectedGrant={state.data}
              stakedAmount={state.data.staked}
            />
          </div>
        </section>
      </PageWrapper>
    </LoadingOverlay>
  )
}

const getCustomWeb3Context = () => {
  const web3 = new Web3(getWsUrl())
  return {
    web3,
    grantContract: new web3.eth.Contract(
      TokenGrant.abi,
      getContractAddress(TokenGrant)
    ),
  }
}

const TokenGrantPreviewPage = () => {
  useEffect(() => {})
  return (
    <Web3Context.Provider value={getCustomWeb3Context()}>
      <TokenGrantPreview />
    </Web3Context.Provider>
  )
}

export default TokenGrantPreviewPage
