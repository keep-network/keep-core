import React, { useContext, useEffect } from "react"
import { CSSTransition, TransitionGroup } from "react-transition-group"
import Banner from "./Banner"
import { showMessage, closeMessage } from "../actions/messages"
import { connect } from "react-redux"
import { ViewInBlockExplorer } from "./ViewInBlockExplorer"
import * as Icons from "./Icons"

export const MessagesContext = React.createContext({})

export const messageType = {
  SUCCESS: { icon: Icons.OK },
  ERROR: { icon: Icons.OK },
  PENDING_ACTION: { icon: Icons.Time },
  INFO: { icon: Icons.Question },
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
    <Banner>
      <div className="flex row center">
        <Banner.Icon icon={message.type.icon} className="mr-1" />
        <div>
          <Banner.Title>{message.title}</Banner.Title>
          <Banner.Action>
            {message.withTransactionHash ? (
              <ViewInBlockExplorer
                type="tx"
                className="arrow-link"
                id={message.id}
              />
            ) : (
              message.subtitle
            )}
          </Banner.Action>
        </div>
        <Banner.CloseIcon onClick={() => onMessageClose(message)} />
      </div>
    </Banner>
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
