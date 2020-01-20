import { useEffect, useReducer } from 'react'

const FETCH_REQUEST_START = 'FETCH_REQUEST_START'
const FETCH_REQUEST_SUCCESS = 'FETCH_REQUEST_SUCCESS'
const FETCH_REQUEST_FAILURE = 'FETCH_REQUEST_FAILURE'

export const useFetchData = (serviceMethod, initialData) => {
  const [state, dispatch] = useReducer(dataFetchrReducer, {
    isFetching: false,
    isError: false,
    data: initialData,
  })

  useEffect(() => {
    const shouldSetState = true

    dispatch({ type: FETCH_REQUEST_START })
    serviceMethod()
      .then((data) => {
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

const dataFetchrReducer = (state, action) => {
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
