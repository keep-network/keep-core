import React, { useContext, useEffect } from "react"
import { CSSTransition, TransitionGroup } from "react-transition-group"
import Banner, { BANNER_TYPE } from "./Banner"
import { showMessage, closeMessage } from "../actions/messages"
import { connect } from "react-redux"

export const MessagesContext = React.createContext({})

export const messageType = {
  SUCCESS: BANNER_TYPE.SUCCESS,
  ERROR: BANNER_TYPE.ERROR,
  PENDING_ACTION: BANNER_TYPE.PENDING,
  INFO: BANNER_TYPE.DISABLED,
}

const messageTransitionTimeoutInMs = 500

class MessagesComponent extends React.Component {
  showMessage = (options) => {
    this.props.showMessage(options)
  }

  onMessageClose = (message) => {
    this.props.closeMessage(message.id)
  }

  render() {
    return (
      <MessagesContext.Provider
        value={{
          showMessage: this.showMessage,
          closeMessage: this.onMessageClose,
        }}
      >
        <div className="messages-container">
          <TransitionGroup>
            {this.props.messages.map((message) => (
              <CSSTransition
                timeout={messageTransitionTimeoutInMs}
                key={message.id}
                classNames="banner"
              >
                <Message
                  key={message.id}
                  message={message}
                  onMessageClose={this.onMessageClose}
                />
              </CSSTransition>
            ))}
          </TransitionGroup>
        </div>
        {this.props.children}
      </MessagesContext.Provider>
    )
  }
}

const mapStateToProps = (state) => {
  return { messages: state.messages }
}

const mapDispatchToProps = {
  showMessage,
  closeMessage,
}

export const Messages = connect(
  mapStateToProps,
  mapDispatchToProps
)(MessagesComponent)

const closeMessageTimeoutInMs = 3250

const Message = ({ message, onMessageClose }) => {
  useEffect(() => {
    if (!message.sticky) {
      const timeout = setTimeout(
        () => onMessageClose(message),
        closeMessageTimeoutInMs
      )
      return () => clearTimeout(timeout)
    }
  }, [message, onMessageClose])

  return (
    <Banner
      type={message.type}
      title={message.title}
      subtitle={message.content}
      withIcon
      withCloseIcon
      onCloseIcon={() => onMessageClose(message)}
    />
  )
}

export const useShowMessage = () => {
  const { showMessage } = useContext(MessagesContext)

  return showMessage
}

export const useCloseMessage = () => {
  const { closeMessage } = useContext(MessagesContext)

  return closeMessage
}
