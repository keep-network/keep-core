import React from "react"
import * as Icons from "./Icons"

const Loadable = ({ text, isFetching, children }) =>
  isFetching ? <div>{text}</div> : children

Loadable.defaultProps = {
  text: "Loading...",
  isFetching: true,
}

export const ClockIndicator = ({ color }) => (
  <div className={`indicator-clock ${color}`} />
)

ClockIndicator.defaultProps = {
  color: "",
}

export const LoadingOverlay = React.memo(
  ({ isFetching, classNames, children, skeletonComponent }) => {
    return (
      <div className={`loading-overlay-container ${classNames}`}>
        {children}
        <div className={`loading-overlay ${isFetching ? "show" : "hidden"}`}>
          {skeletonComponent}
        </div>
      </div>
    )
  }
)

export const KeepLoadingIndicator = () => <Icons.KeepLoadingIndicator />

export default Loadable
