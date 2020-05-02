import React from "react"
import { useFetchData } from "../hooks/useFetchData"
import { authorizationService } from "../services/authorization.service"
import { LoadingOverlay } from "../components/Loadable"
import AddressShortcut from "../components/AddressShortcut"
import { ETHERSCAN_DEFAULT_URL } from "../constants/constants"
import Tile from "./Tile"

const initialData = { isOperator: false, contracts: [] }

const AuthorizationInfo = (props) => {
  const [state] = useFetchData(
    authorizationService.fetchOperatorAuthorizedContracts,
    initialData
  )
  const { isFetching, data } = state
  const { isOperator, contracts } = data

  if (!isOperator) {
    return null
  }

  return (
    <LoadingOverlay isFetching={isFetching}>
      <Tile title="Authorization">
        {contracts.map((contract) => (
          <div key={contract.contractAddress} className="mb-1">
            You have been successfully authorized by authorizer&nbsp;
            <AddressShortcut address={contract.authorizer} />
            &nbsp; to &nbsp;
            <a
              href={ETHERSCAN_DEFAULT_URL + contract.contractAddress}
              rel="noopener noreferrer"
              target="_blank"
            >
              operator contract
            </a>
            .
          </div>
        ))}
      </Tile>
    </LoadingOverlay>
  )
}

export default AuthorizationInfo
