import React, { useEffect, useCallback, useState, useMemo } from "react"
import AuthorizeContracts from "../../components/AuthorizeContracts"
import { tbtcAuthorizationService } from "../../services/tbtc-authorization.service"
import { useFetchData } from "../../hooks/useFetchData"
import { BondingSection } from "../../components/BondingSection"
import { useSubscribeToContractEvent } from "../../hooks/useSubscribeToContractEvent"
import {
  findIndexAndObject,
  compareEthAddresses,
} from "../../utils/array.utils"
import { add, sub } from "../../utils/arithmetics.utils"
import web3Utils from "web3-utils"
import {
  KEEP_BONDING_CONTRACT_NAME,
  AUTH_CONTRACTS_LABEL,
} from "../../constants/constants"
import { LoadingOverlay } from "../../components/Loadable"
import { isSameEthAddress } from "../../utils/general.utils"
import DataTableSkeleton from "../../components/skeletons/DataTableSkeleton"
import {
  authorizeSortitionPoolContract,
  authorizeOperatorContract,
  deauthorizeSortitionPoolContract,
} from "../../actions/web3"
import { connect } from "react-redux"
import { getBondedECDSAKeepFactoryAddress } from "../../contracts"
import EmptyStatePage from "./EmptyStatePage"
import { useWeb3Address } from "../../components/WithWeb3Context"

const initialData = []
const TBTCApplicationPage = ({
  authorizeSortitionPoolContract,
  authorizeOperatorContract,
  deauthorizeSortitionPoolContract,
}) => {
  const [selectedOperator, setOperator] = useState({})
  const address = useWeb3Address()

  // fetch data from service
  const [tbtcAuthState, updateTbtcAuthData, , setAuthDataArgs] = useFetchData(
    tbtcAuthorizationService.fetchTBTCAuthorizationData,
    initialData,
    address
  )
  // fetch bonding data
  const [bondingState, updateBondinData, , setBondingArgs] = useFetchData(
    tbtcAuthorizationService.fetchBondingData,
    initialData,
    address
  )

  useEffect(() => {
    setBondingArgs([address])
  }, [setBondingArgs, address])

  useEffect(() => {
    setAuthDataArgs([address])
  }, [setAuthDataArgs, address])

  const unbondedValueUpdated = useCallback(
    (event, arithmeticOpration = add) => {
      const {
        returnValues: { operator, amount },
      } = event
      const { indexInArray, obj: obsoleteData } = findIndexAndObject(
        "operatorAddress",
        operator,
        bondingState.data,
        compareEthAddresses
      )
      if (indexInArray === null) {
        return
      }

      const availableETHInWei = arithmeticOpration(
        obsoleteData.availableETHInWei,
        amount
      ).toString()
      const availableETH = web3Utils.fromWei(availableETHInWei, "ether")
      const updatedBondinData = [...bondingState.data]
      updatedBondinData[indexInArray] = {
        ...obsoleteData,
        availableETH,
        availableETHInWei,
      }
      updateBondinData(updatedBondinData)
    },
    [updateBondinData, bondingState.data]
  )

  const subscribeToUnbondedValueDepositedCallback = (event) => {
    unbondedValueUpdated(event)
  }

  useSubscribeToContractEvent(
    KEEP_BONDING_CONTRACT_NAME,
    "UnbondedValueDeposited",
    subscribeToUnbondedValueDepositedCallback
  )

  const unbondedValueWithdrawnCallback = (event) => {
    unbondedValueUpdated(event, sub)
  }

  useSubscribeToContractEvent(
    KEEP_BONDING_CONTRACT_NAME,
    "UnbondedValueWithdrawn",
    unbondedValueWithdrawnCallback
  )

  const onSuccessCallback = useCallback(
    (contractName, operatorAddress, isAuthorized = true) => {
      const { indexInArray: operatorIndexInArray, obj: obsoleteOperator } =
        findIndexAndObject(
          "operatorAddress",
          operatorAddress,
          tbtcAuthState.data,
          compareEthAddresses
        )
      if (operatorIndexInArray === null) {
        return
      }
      const { indexInArray: contractIndexInArray, obj: obsoleteContract } =
        findIndexAndObject(
          "contractName",
          contractName,
          obsoleteOperator.contracts
        )
      const updatedContracts = [...obsoleteOperator.contracts]
      updatedContracts[contractIndexInArray] = {
        ...obsoleteContract,
        isAuthorized,
      }
      const updatedOperators = [...tbtcAuthState.data]
      updatedOperators[operatorIndexInArray] = {
        ...obsoleteOperator,
        contracts: updatedContracts,
      }

      updateTbtcAuthData(updatedOperators)
    },
    [updateTbtcAuthData, tbtcAuthState.data]
  )

  const authorizeContract = useCallback(
    async (data, awaitingPromise) => {
      const { operatorAddress, contractName } = data
      if (contractName === AUTH_CONTRACTS_LABEL.TBTC_SYSTEM) {
        const sortitionPoolAddress =
          await tbtcAuthorizationService.fetchSortitionPoolForTbtc()

        authorizeSortitionPoolContract(
          {
            operatorAddress,
            sortitionPoolAddress,
          },
          awaitingPromise
        )
      } else {
        const operatorContractAddress = getBondedECDSAKeepFactoryAddress()
        authorizeOperatorContract(
          { operatorAddress, operatorContractAddress },
          awaitingPromise
        )
      }
    },
    [authorizeSortitionPoolContract, authorizeOperatorContract]
  )

  const deauthorizeTBTCSystem = useCallback(
    async (data, awaitingPromise) => {
      const { operatorAddress } = data
      const sortitionPoolAddress =
        await tbtcAuthorizationService.fetchSortitionPoolForTbtc()
      deauthorizeSortitionPoolContract(
        { operatorAddress, sortitionPoolAddress },
        awaitingPromise
      )
    },
    [deauthorizeSortitionPoolContract]
  )

  const tbtcAuthData = useMemo(() => {
    if (!selectedOperator.operatorAddress) {
      return tbtcAuthState.data
    }
    return tbtcAuthState.data.filter((data) =>
      isSameEthAddress(data.operatorAddress, selectedOperator.operatorAddress)
    )
  }, [selectedOperator.operatorAddress, tbtcAuthState.data])

  return (
    <>
      <nav className="mb-2">
        <a
          href="https://tbtc.network/"
          className="h4"
          rel="noopener noreferrer"
          target="_blank"
        >
          tBTC Website
        </a>
      </nav>
      <LoadingOverlay
        isFetching={tbtcAuthState.isFetching}
        skeletonComponent={
          <DataTableSkeleton columns={4} subtitleWidth="40%" />
        }
      >
        <AuthorizeContracts
          filterDropdownOptions={tbtcAuthState.data}
          onSelectOperator={setOperator}
          selectedOperator={selectedOperator}
          data={tbtcAuthData}
          onAuthorizeBtn={authorizeContract}
          onDeauthorizeBtn={deauthorizeTBTCSystem}
          onSuccessCallback={onSuccessCallback}
        />
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={bondingState.isFetching}
        skeletonComponent={<DataTableSkeleton subtitleWidth="70%" />}
      >
        <BondingSection data={bondingState.data} />
      </LoadingOverlay>
    </>
  )
}

const mapDispatchToProps = {
  authorizeSortitionPoolContract,
  authorizeOperatorContract,
  deauthorizeSortitionPoolContract,
}

const ConnectedTBTCApplicationPage = connect(
  null,
  mapDispatchToProps
)(TBTCApplicationPage)

ConnectedTBTCApplicationPage.route = {
  title: "tBTC",
  path: "/applications/tbtc",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default ConnectedTBTCApplicationPage
