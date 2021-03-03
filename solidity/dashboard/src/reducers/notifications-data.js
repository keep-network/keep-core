const initialState = {
  liquidityRewardNotification: {
    pairsDisplayed: [],
  },
}

const notificationsDataReducer = (state = initialState, action) => {
  if (!action.payload) {
    return state
  }

  switch (action.type) {
    case "notifications_data/liquidityRewardNotification/pairs_displayed_updated":
      return {
        ...state,
        liquidityRewardNotification: {
          ...state["liquidityRewardNotification"],
          pairsDisplayed: action.payload,
        },
      }
    default:
      return state
  }
}

export default notificationsDataReducer
