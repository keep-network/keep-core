import React, { useEffect, useState, useRef } from "react"
import { CSSTransition } from 'react-transition-group'
import Loadable from "./Loadable"

const buttonContentTransitionTimeoutInMs = 500
const minimumLoaderDurationInMs = 400

const useMinimumLoaderDuration = (showLoader, setShowLoader, isFetching) => {
    useEffect(() => {
        if (isFetching)
            setShowLoader(true);

        if (!isFetching && showLoader) {
            const timeout = setTimeout(() => setShowLoader(false), minimumLoaderDurationInMs);

        return () => clearTimeout(timeout);
        }
    }, [isFetching, showLoader])
}

const useCurrentButtonDimensions = (buttonRef, children) => {
    const [width, setWidth] = useState(0)
    const [height, setHeight] = useState(0)
  
    useEffect(() => {
      if (buttonRef.current && buttonRef.current.getBoundingClientRect().width) {
        setWidth(buttonRef.current.getBoundingClientRect().width)
      }
      if (buttonRef.current && buttonRef.current.getBoundingClientRect().height) {
        setHeight(buttonRef.current.getBoundingClientRect().height)
      }
    }, [children])

    return [width, height]
}

export default function Button({ isFetching, children, ...props }) {
  const [showLoader, setShowLoader] = React.useState(false)
  const buttonRef = useRef(null)
  const [width, height] = useCurrentButtonDimensions(buttonRef, children)

  useMinimumLoaderDuration(showLoader, setShowLoader, isFetching)
  
  return (
    <button
      {...props}
      ref={buttonRef}
      style={showLoader ? { width: `${width}px`, height: `${height}px` } : {} }
      disabled={showLoader}
    >
      <CSSTransition
        in={showLoader}
        timeout={buttonContentTransitionTimeoutInMs}
        classNames="button-content"
      >
        <div className="button-content">
            { showLoader ? <Loadable text="In progress" /> : children }
        </div>
      </CSSTransition>
    </button>
  )
}

export const SubmitButton = ({ onSubmitAction, ...props }) => {
  const [isFetching, setIsFetching] = useState(false)

  const onButtonClick = async (event) => {
    event.preventDefault()
    setIsFetching(true)

    try {
      await onSubmitAction()
      setIsFetching(false)
    } catch(error) {
      setIsFetching(false)
    }
  }

  return <Button {...props} onClick={onButtonClick} isFetching={isFetching} />
}
