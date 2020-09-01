import React, { useEffect, useRef, useState } from "react"
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
    const ref = useRef(null)
    const [fetchingDelay, setFetchingDelay] = useState(true)

    // force switch between skeleton -> children when `isFetching` prop has a constant value.
    useEffect(() => {
      setFetchingDelay(isFetching)
    }, [isFetching])

    return (
      <div ref={ref} className="loading-overlay-container">
        <SwitchTransition mode={"out-in"}>
          <CSSTransition
            key={fetchingDelay}
            addEndListener={(node, done) => {
              node.addEventListener("transitionend", done, false)
            }}
            classNames="loading-overlay"
            onEntering={() => {
              if (ref.current) {
                ref.current.style.backgroundColor = "transparent"
              }
            }}
          >
            <div className="loading-overlay">
              {fetchingDelay ? skeletonComponent : children}
            </div>
          </CSSTransition>
        </SwitchTransition>
      </div>
    )
  }
)

export const KeepLoadingIndicator = () => <Icons.KeepLoadingIndicator />

export default Loadable
