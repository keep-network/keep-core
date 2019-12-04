import React, { useContext, useEffect } from 'react'
import { CSSTransition, TransitionGroup } from 'react-transition-group';

var messageId = 0

export class Messages extends React.Component { 
    constructor(props) {
        super(props)
        this.state = { messages: [] }
    }

    show = (value) => {
        value.id = messageId++
        this.setState({ messages: this.state.messages ? [...this.state.messages, value] : [value]})
    }

    onClose = (message) => {
        const updatedMessages = this.state.messages.filter(m => m.id !== message.id)
        this.setState({ messages: updatedMessages })
    }

    render() {
        return (
            <MessagesContext.Provider value={{ show: this.show }} >
                <div className="messages-container">
                    <TransitionGroup >
                        {this.state.messages.map(message => (
                            <CSSTransition
                                key={message.id}
                                timeout={500}
                                classNames="message"
                            >
                                <Message key={message.id} message={message} onClose={this.onClose} />
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

const Message = ({ message, ...props }) => {
    useEffect(() => {
        const timeout = setTimeout(onClose, 3250);
        return () => clearTimeout(timeout)
    }, [message.id])

    const onClose = () => {
       props.onClose(message)
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
                <div className='message-icon-close' onClick={onClose}>
                    <span className="glyphicon glyphicon-remove" aria-hidden='true' />
                </div>
            </div>
        </div>
    )
}

export const useShowMessage = () => {
    const { show } = useContext(MessagesContext)

    return show
}

const MessagesContext = React.createContext({})