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

const Message = (props) => {
    useEffect(() => {
        const timeout = setTimeout(onClose, 3250);
        return () => clearTimeout(timeout)
    }, [props.message.id])


    const onClose = () => {
       props.onClose(props.message)
    }

    return (
        <div className={`message message-${props.message.type || 'success'}`}>
            <button onClick={onClose}>Close</button>
            {props.message.content}
        </div>
    )
}

export const useShowMessage = () => {
    const { show } = useContext(MessagesContext)

    return show
}

const MessagesContext = React.createContext({})