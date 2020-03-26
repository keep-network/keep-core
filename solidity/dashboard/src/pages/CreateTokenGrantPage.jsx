import React, { useContext } from 'react'
import CreateTokenGrantForm from '../components/CreateTokenGrantForm'
import { ContractsDataContext } from '../components/ContractsDataContextProvider'

const CreateTokenGrantPage = () => {
  const {
    tokenBalance,
    refreshKeepTokenBalance,
  } = useContext(ContractsDataContext)
  return (
    <React.Fragment>
      <h2 className="mb-2">
        Create Token Grant
      </h2>
      <section className="rewards-history tile flex column">
        <h3 className="text-grey-70 mb-1">
            Create Grant
        </h3>
        <CreateTokenGrantForm
          keepBalance={tokenBalance}
          successCallback={refreshKeepTokenBalance}
        />
      </section>
    </React.Fragment>
  )
}

export default CreateTokenGrantPage
