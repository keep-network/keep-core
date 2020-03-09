import React, { useContext, useEffect } from 'react'
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import Banner, { BANNER_TYPE } from './Banner'

export const MessagesContext = React.createContext({})

export const messageType = {
  'SUCCESS': BANNER_TYPE.SUCCESS,
  'ERROR': BANNER_TYPE.ERROR,
  'PENDING_ACTION': BANNER_TYPE.PENDING,
  'INFO': BANNER_TYPE.DISABLED,
}

let messageId = 1
const messageTransitionTimeoutInMs = 500

export class Messages extends React.Component {
  constructor(props) {
    super(props)
    this.state = { messages: [] }
  }

    showMessage = (value) => {
      value.id = messageId++
      this.setState({ messages: this.state.messages ? [...this.state.messages, value] : [value] })
      return value
    }

    onMessageClose = (message) => {
      const updatedMessages = this.state.messages.filter((m) => m.id !== message.id)
      this.setState({ messages: updatedMessages })
    }

    render() {
      return (
        <MessagesContext.Provider value={{ showMessage: this.showMessage, closeMessage: this.onMessageClose }} >
          <div className="messages-container">
            <TransitionGroup >
              {this.state.messages.map((message) => (
                <CSSTransition
                  timeout={messageTransitionTimeoutInMs}
                  key={message.id}
                  classNames="banner"
                >
                  <Message key={message.id} message={message} onMessageClose={this.onMessageClose} />
                </CSSTransition>
              ))}
            </TransitionGroup>
          </div>
          {this.props.children}
        </MessagesContext.Provider>
      )
    }
}

const closeMessageTimeoutInMs = 3250

const Message = ({ message, ...props }) => {
  useEffect(() => {
    if (!message.sticky) {
      const timeout = setTimeout(onMessageClose, closeMessageTimeoutInMs)
      return () => clearTimeout(timeout)
    }
  }, [message.id])

  const onMessageClose = () => {
    props.onMessageClose(message)
  }

  return (
    <Banner
      type={message.type}
      title={message.title}
      subtitle={message.content}
      withIcon
      withCloseIcon
      onCloseIcon={onMessageClose}
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
