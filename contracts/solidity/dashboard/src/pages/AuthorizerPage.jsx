import React, { useState, useEffect, useContext, useMemo } from 'react'
import AuthorizeContracts from '../components/AuthorizeContracts'
import AuthorizationHistory from '../components/AuthorizationHistory'
import { authorizationService } from '../services/authorization.service'
import { LoadingOverlay } from '../components/Loadable'
import { useFetchData } from '../hooks/useFetchData'
import Dropdown from '../components/Dropdown'
import { Web3Context } from '../components/WithWeb3Context'
import { useShowMessage, messageType } from '../components/Message'

const initialData = {}

const AuthorizerPage = () => {
  const [
    state,
    updateData,
    refreshData,
  ] = useFetchData(authorizationService.fetchAuthorizationPageData, initialData)
  const { yourAddress, stakingContract } = useContext(Web3Context)
  const showMessage = useShowMessage()

  const { isFetching, data } = state
  const [operator, setOperator] = useState({})

  useEffect(() => {
    if (data && Object.keys(data).length >= 0) {
      setOperator({ address: Object.keys(data)[0] })
    }
  }, [data])

  const authorizeContaract = async (contract, onTransactionHashCallback) => {
    try {
      await stakingContract
        .methods
        .authorizeOperatorContract(operator.address, contract.contractAddress)
        .send({ from: yourAddress })
        .on('transactionHash', onTransactionHashCallback)
      showMessage({ type: messageType.SUCCESS, title: 'Success', content: 'You have successfully authorized operator' })
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: 'Error', content: error.message })
      throw error
    }
  }

  const dropdownOptions = useMemo(() => {
    return Object.keys(data).map((key) => ({ address: key }))
  }, [data])

  return (
    <>
      <h2 className="mb-2">Authorizations</h2>
      <section className="tile">
        <h3 className="text-grey-60 mb-1">Available Operators</h3>
        <Dropdown
          options={dropdownOptions}
          onSelect={(operator) => setOperator(operator)}
          valuePropertyName='address'
          labelPropertyName='address'
          selectedItem={operator}
          labelPrefix='Operator:'
          noItemSelectedText='Select Operator'
          label=''
        />
      </section>
      <LoadingOverlay isFetching={isFetching}>
        <AuthorizeContracts
          contracts={data[operator.address] && data[operator.address].contractsToAuthorize}
          onAuthorizeBtn={authorizeContaract}
          onAuthorizeSuccessCallback={refreshData}
        />
      </LoadingOverlay>
      <LoadingOverlay isFetching={isFetching}>
        <AuthorizationHistory
          contracts={data[operator.address] && data[operator.address].authorizedContracts}
        />
      </LoadingOverlay>
    </>
  )
}

export default AuthorizerPage
