import React, { useEffect, useState, useRef, useContext } from 'react'
import { CSSTransition } from 'react-transition-group'
import { ClockIndicator } from './Loadable'
import { messageType, MessagesContext } from './Message'
import * as Icons from './Icons'

const buttonContentTransitionTimeoutInMs = 500
const minimumLoaderDurationInMs = 400
const minWidthPendingButton = 130

const useMinimumLoaderDuration = (showLoader, setShowLoader, isFetching) => {
  useEffect(() => {
    if (isFetching) {
      setShowLoader(true)
    }

    if (!isFetching && showLoader) {
      const timeout = setTimeout(() => setShowLoader(false), minimumLoaderDurationInMs)

      return () => clearTimeout(timeout)
    }
  }, [isFetching, showLoader])
}

const useCurrentButtonDimensions = (buttonRef, children) => {
  const [width, setWidth] = useState(0)
  const [height, setHeight] = useState(0)

  useEffect(() => {
    if (buttonRef.current && buttonRef.current.getBoundingClientRect().width) {
      const width = buttonRef.current.getBoundingClientRect().width
      setWidth(width < minWidthPendingButton ? minWidthPendingButton : width)
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
      disabled={props.disabled || showLoader}
      className={`${props.className}${showLoader ? ' pending' : ''}`}
    >
      <CSSTransition
        in={showLoader}
        timeout={buttonContentTransitionTimeoutInMs}
        classNames="button-content"
      >
        <div className="button-content">
          { showLoader ?
            <div className="flex flex-1 flex-center">
              <ClockIndicator color='primary' />
              <span className="ml-1 text-primary">pending</span>
            </div> :
            children
          }
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
  ...props
}) => {
  const [isFetching, setIsFetching] = useState(false)
  const { showMessage, closeMessage } = useContext(MessagesContext)
  const [showSuccessBtn, setShowSuccessBtn] = useState(false)
  useEffect(() => {
    if (showSuccessBtn) {
      const timeout = setTimeout(() => {
        setShowSuccessBtn(false)
        successCallback()
      }, successBtnVisibilityDuration)
      return () => clearTimeout(timeout)
    }
  }, [showSuccessBtn])

  let pendingMessage = { type: messageType.PENDING_ACTION, sticky: true, title: pendingMessageTitle, content: pendingMessageContent }
  let infoMessage = { type: messageType.INFO, sticky: true, title: 'Waiting for the transaction confirmation...' }

  const onTransactionHashCallback = (hash) => {
    pendingMessage = showMessage({ ...pendingMessage, content: `Transaction hash: ${hash}` })
    closeMessage(infoMessage)
  }

  const openMessageInfo = () => {
    infoMessage = showMessage(infoMessage)
  }

  const setFetching = () => setIsFetching(true)

  const onButtonClick = async (event) => {
    event.preventDefault()
    if (!triggerManuallyFetch) {
      setIsFetching(true)
    }
    if (withMessageActionIsPending) {
      infoMessage = showMessage(infoMessage)
    }

    try {
      await onSubmitAction(onTransactionHashCallback, openMessageInfo, setFetching)
      setIsFetching(false)
      setShowSuccessBtn(true)
    } catch (error) {
      setIsFetching(false)
    }

    closeMessage(pendingMessage)
    closeMessage(infoMessage)
  }

  return (
    <Button
      {...props}
      className={`${props.className} ${showSuccessBtn && `btn btn-success`}`}
      onClick={onButtonClick}
      isFetching={isFetching}
      disabled={showSuccessBtn}
    >
      {showSuccessBtn ?
        <div className="flex flex-1 flex-center">
          <Icons.OK />
          <span className="ml-1 text-black">success</span>
        </div> :
        props.children
      }
    </Button>
  )
}

SubmitButton.defaultProps = {
  withMessageActionIsPending: true,
  triggerManuallyFetch: false,
  pendingMessageTitle: 'Action is pending',
  pendingMessageContent: '',
  successCallback: () => {},
}
