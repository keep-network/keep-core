import React from "react"
import { useFetchData } from "../hooks/useFetchData"
import { authorizationService } from "../services/authorization.service"
import { LoadingOverlay } from "../components/Loadable"
import AddressShortcut from "../components/AddressShortcut"
import Tile from "./Tile"
import ViewAddressInBlockExplorer from "./ViewAddressInBlockExplorer"

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
            <ViewAddressInBlockExplorer
              address={contract.contractAddress}
              text="operator contract"
            />
            .
          </div>
        ))}
      </Tile>
    </LoadingOverlay>
  )
}

export default AuthorizationInfo
