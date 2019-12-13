import React, { useContext, useEffect } from 'react'
import { CSSTransition, TransitionGroup } from 'react-transition-group';

export const MessagesContext = React.createContext({})

let messageId = 0
const messageTransitionTimeoutInMs = 500

export class Messages extends React.Component { 
    constructor(props) {
        super(props)
        this.state = { messages: [] }
    }

    showMessage = (value) => {
        value.id = messageId++
        this.setState({ messages: this.state.messages ? [...this.state.messages, value] : [value]})
    }

    onMessageClose = (message) => {
        const updatedMessages = this.state.messages.filter(m => m.id !== message.id)
        this.setState({ messages: updatedMessages })
    }

    render() {
        return (
            <MessagesContext.Provider value={{ showMessage: this.showMessage }} >
                <div className="messages-container">
                    <TransitionGroup >
                        {this.state.messages.map(message => (
                            <CSSTransition
                                timeout={messageTransitionTimeoutInMs}
                                key={message.id}
                                classNames="message"
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

const messageIconMap = {
    error: 'glyphicon-remove',
    success: 'glyphicon-ok'
}

const closeMessageTimeoutInMs = 3250

const Message = ({ message, ...props }) => {
    useEffect(() => {
        const timeout = setTimeout(onMessageClose, closeMessageTimeoutInMs);
        return () => clearTimeout(timeout)
    }, [message.id])

    const onMessageClose = () => {
       props.onMessageClose(message)
    }

    return (
        <div className={`message message-${message.type || 'success'}`}>
            <div className='message-content-wrapper'>
                <div className="message-icon">
                    <span className={`glyphicon ${messageIconMap[message.type]}`} aria-hidden='true' />
                </div>
                <div className='message-content'>
                    <span className="message-title">{message.title}</span>
                    <div>{message.content}</div>
                </div>
                <div className='message-icon-close' onClick={onMessageClose}>
                    <span className="glyphicon glyphicon-remove" aria-hidden='true' />
                </div>
            </div>
        </div>
    )
}

export const useShowMessage = () => {
    const { showMessage } = useContext(MessagesContext)

    return showMessage
}