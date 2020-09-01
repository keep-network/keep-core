import React, { useEffect, useState, useRef } from "react"
import { CSSTransition } from "react-transition-group"
import { ClockIndicator } from "./Loadable"
import * as Icons from "./Icons"
import { Deferred } from "../contracts"

const buttonContentTransitionTimeoutInMs = 500
const minimumLoaderDurationInMs = 400
const minWidthPendingButton = 130
const minHeightPendingButton = 38

const useMinimumLoaderDuration = (showLoader, setShowLoader, isFetching) => {
  useEffect(() => {
    if (isFetching) {
      setShowLoader(true)
    }

    if (!isFetching && showLoader) {
      const timeout = setTimeout(
        () => setShowLoader(false),
        minimumLoaderDurationInMs
      )

      return () => clearTimeout(timeout)
    }
  }, [isFetching, showLoader, setShowLoader])
}

const useCurrentButtonDimensions = (buttonRef, children) => {
  const [width, setWidth] = useState(0)
  const [height, setHeight] = useState(0)

  useEffect(() => {
    if (buttonRef.current && buttonRef.current.getBoundingClientRect().width) {
      const width = buttonRef.current.getBoundingClientRect().width
      setWidth(width < minWidthPendingButton ? minWidthPendingButton : width)
    } else {
      setWidth(minWidthPendingButton)
    }
    if (buttonRef.current && buttonRef.current.getBoundingClientRect().height) {
      setHeight(buttonRef.current.getBoundingClientRect().height)
    } else {
      setHeight(minHeightPendingButton)
    }
  }, [buttonRef, children])

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
      style={showLoader ? { width: `${width}px`, height: `${height}px` } : {}}
      disabled={props.disabled || showLoader}
      className={`${props.className}${showLoader ? " pending" : ""}`}
    >
      <CSSTransition
        in={showLoader}
        timeout={buttonContentTransitionTimeoutInMs}
        classNames="button-content"
      >
        <div className="button-content">
          {showLoader ? (
            <div className="flex full-center">
              <span style={{ display: "inline-block" }}>
                {" "}
                <ClockIndicator color="primary" />
              </span>
              <span className="ml-1 text-primary">pending</span>
            </div>
          ) : (
            children
          )}
        </div>
      </CSSTransition>
    </button>
  )
}

const successBtnVisibilityDuration = 5000 // 5s

export const SubmitButton = ({
  onSubmitAction,
  withMessageActionIsPending,
  pendingMessageTitle,
  pendingMessageContent,
  triggerManuallyFetch,
  successCallback,
  confirmationModalTitle,
  ...props
}) => {
  const [isSubmitting, setSubmitting] = useState(false)
  const [showSuccessBtn, setShowSuccessBtn] = useState(false)

  useEffect(() => {
    if (showSuccessBtn) {
      const timeout = setTimeout(() => {
        setShowSuccessBtn(false)
        successCallback()
      }, successBtnVisibilityDuration)
      return () => clearTimeout(timeout)
    }
  }, [showSuccessBtn, successCallback])

  const onButtonClick = async (event) => {
    event.preventDefault()
    const awaitingPromise = new Deferred()
    try {
      setSubmitting(true)

      await onSubmitAction(awaitingPromise)
      await awaitingPromise.promise

      setSubmitting(false)
      setShowSuccessBtn(true)
    } catch (error) {
      setSubmitting(false)
    }
  }

  return (
    <>
      <Button
        {...props}
        className={`${props.className} ${showSuccessBtn && `btn btn-success`}`}
        onClick={onButtonClick}
        isFetching={isSubmitting}
        disabled={showSuccessBtn || props.disabled}
      >
        {showSuccessBtn ? (
          <div className="flex row full-center flex-1">
            <Icons.OK />
            <span className="ml-1 text-black">success</span>
          </div>
        ) : (
          props.children
        )}
      </Button>
    </>
  )
}

SubmitButton.defaultProps = {
  withMessageActionIsPending: true,
  triggerManuallyFetch: false,
  pendingMessageTitle: "Action is pending",
  pendingMessageContent: "",
  successCallback: () => {},
  confirmationModalTitle: "Are you sure?",
}
