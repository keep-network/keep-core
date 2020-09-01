import { useEffect, useReducer, useContext, useCallback, useState } from "react"
import { Web3Context } from "../components/WithWeb3Context"
import { wait } from "../utils/general.utils"
import { usePrevious } from "./usePrevious"

const FETCH_REQUEST_START = "FETCH_REQUEST_START"
const FETCH_REQUEST_SUCCESS = "FETCH_REQUEST_SUCCESS"
const FETCH_REQUEST_FAILURE = "FETCH_REQUEST_FAILURE"
const UPDATE_DATA = "UPDATE_DATA"
const REFRESH_DATA = "REFRESH_DATA"
const syncState = { UP_TO_DATE: "UP_TO_DATE", OBSOLETE: "OBSOLETE" }

const requestTimeDelay = 500 // 0.5s

export const useFetchData = (
  serviceMethod,
  initialData,
  ...initialServiceMethodArgs
) => {
  const [serviceMethodArgs, setServiceMethodArgs] = useState(
    initialServiceMethodArgs
  )
  const web3Context = useContext(Web3Context)
  const [state, dispatch] = useReducer(dataFetchReducer, {
    isFetching: true,
    isError: false,
    error: null,
    data: initialData,
    syncState: syncState.UP_TO_DATE,
  })
  const prevSyncState = usePrevious(state.syncState)

  const fetchData = () => {
    let shouldSetState = true

    dispatch({ type: FETCH_REQUEST_START })
    Promise.all([
      serviceMethod(web3Context, ...serviceMethodArgs),
      wait(requestTimeDelay),
    ])
      .then(([data]) => {
        shouldSetState &&
          dispatch({ type: FETCH_REQUEST_SUCCESS, payload: data })
      })
      .catch((error) => {
        shouldSetState &&
          dispatch({ type: FETCH_REQUEST_FAILURE, payload: error })
      })

    return () => {
      shouldSetState = false
    }
  }

  useEffect(fetchData, [serviceMethodArgs])
  useEffect(() => {
    if (
      prevSyncState === syncState.UP_TO_DATE &&
      state.syncState === syncState.OBSOLETE
    ) {
      fetchData()
    }
  })

  const updateData = useCallback((updatedData) => {
    dispatch({ type: UPDATE_DATA, payload: updatedData })
  }, [])

  const refreshData = useCallback(() => {
    dispatch({ type: REFRESH_DATA })
  }, [])

  return [state, updateData, refreshData, setServiceMethodArgs]
}

const dataFetchReducer = (state, action) => {
  switch (action.type) {
    case FETCH_REQUEST_START:
      return {
        ...state,
        isFetching: true,
        isError: false,
        error: null,
      }
    case FETCH_REQUEST_SUCCESS:
      return {
        ...state,
        syncState: syncState.UP_TO_DATE,
        isFetching: false,
        isError: false,
        data: action.payload,
        error: null,
      }
    case FETCH_REQUEST_FAILURE:
      return {
        ...state,
        isFetching: false,
        isError: true,
        error: action.payload,
      }
    case UPDATE_DATA:
      return {
        ...state,
        data: action.payload,
      }
    case REFRESH_DATA:
      return {
        ...state,
        syncState: syncState.OBSOLETE,
      }
    default:
      return { ...state }
  }
}
