import React from "react"
import { CSSTransition, SwitchTransition } from "react-transition-group"
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
  ({ isFetching, children, skeletonComponent }) => {
    return (
      <div className="loading-overlay-container">
        <SwitchTransition mode={"out-in"}>
          <CSSTransition
            key={isFetching}
            addEndListener={(node, done) => {
              node.addEventListener("transitionend", done, false)
            }}
            classNames="loading-overlay"
          >
            <div className="loading-overlay">
              {isFetching ? skeletonComponent : children}
            </div>
          </CSSTransition>
        </SwitchTransition>
      </div>
    )
  }
)

export const KeepLoadingIndicator = () => <Icons.KeepLoadingIndicator />

export default Loadable
