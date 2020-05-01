import React, { useContext } from "react"
import CreateTokenGrantForm from "../components/CreateTokenGrantForm"
import { ContractsDataContext } from "../components/ContractsDataContextProvider"
import PageWrapper from "../components/PageWrapper"
import Tile from "../components/Tile"

const CreateTokenGrantPage = () => {
  const { tokenBalance, refreshKeepTokenBalance } = useContext(
    ContractsDataContext
  )
  return (
    <PageWrapper title="Create Token Grant">
      <Tile title="Create Grant" className="rewards-history tile flex column">
        <CreateTokenGrantForm
          keepBalance={tokenBalance}
          successCallback={refreshKeepTokenBalance}
        />
      </Tile>
    </PageWrapper>
  )
}

export default CreateTokenGrantPage
