import { useEffect, useReducer, useContext } from 'react'
import { Web3Context } from '../components/WithWeb3Context'
import { wait } from '../utils'

const FETCH_REQUEST_START = 'FETCH_REQUEST_START'
const FETCH_REQUEST_SUCCESS = 'FETCH_REQUEST_SUCCESS'
const FETCH_REQUEST_FAILURE = 'FETCH_REQUEST_FAILURE'

const requestTimeDelay = 500 // 0.5s

export const useFetchData = (serviceMethod, initialData) => {
  const web3Context = useContext(Web3Context)
  const [state, dispatch] = useReducer(dataFetchReducer, {
    isFetching: true,
    isError: false,
    data: initialData,
  })

  useEffect(() => {
    let shouldSetState = true

    dispatch({ type: FETCH_REQUEST_START })
    Promise.all([serviceMethod(web3Context), wait(requestTimeDelay)])
      .then(([data]) => {
        shouldSetState && dispatch({ type: FETCH_REQUEST_SUCCESS, payload: data })
      })
      .catch((error) => {
        shouldSetState && dispatch({ type: FETCH_REQUEST_FAILURE })
      })

    return () => {
      shouldSetState = false
    }
  }, [])

  return state
}

const dataFetchReducer = (state, action) => {
  switch (action.type) {
  case FETCH_REQUEST_START:
    return {
      ...state,
      isFetching: true,
      isError: false,
    }
  case FETCH_REQUEST_SUCCESS:
    return {
      ...state,
      isFetching: false,
      isError: false,
      data: action.payload,
    }
  case FETCH_REQUEST_FAILURE:
    return {
      ...state,
      isFetching: false,
      isError: true,
    }
  default:
    return { ...state }
  }
}
