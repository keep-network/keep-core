import React, { useContext, useEffect } from "react"
import { CSSTransition, TransitionGroup } from "react-transition-group"
import Banner from "./Banner"
import { showMessage, closeMessage } from "../actions/messages"
import { connect } from "react-redux"
import { ViewInBlockExplorer } from "./ViewInBlockExplorer"
import LiquidityRewardsEarnedMessage from "./messages/LiquidityRewardsEarnedMessage"
import LPTokensInWalletMessage from "./messages/LPTokensInWalletMessage"
import WalletMessage from "./messages/WalletMessage"
import PendingActionMessage from "./messages/PendingActionMessage"
import SuccessMessage from "./messages/SuccessMessage"
import ErrorMessage from "./messages/ErrorMessage"
import DelegationAlreadyCopiedMessage from "./messages/DelegationAlreadyCopiedMessage"

export const MessagesContext = React.createContext({})

export const messageType = {
  SUCCESS: SuccessMessage,
  ERROR: ErrorMessage,
  PENDING_ACTION: PendingActionMessage,
  INFO: LiquidityRewardsEarnedMessage,
  WALLET: WalletMessage,
  NEW_LP_TOKENS_IN_WALLET: LPTokensInWalletMessage,
  LIQUIDITY_REWARDS_EARNED: LiquidityRewardsEarnedMessage,
  DELEGATION_ALREADY_COPIED: DelegationAlreadyCopiedMessage,
}

const messageTransitionTimeoutInMs = 500

class MessagesComponent extends React.Component {
  showMessage = (options) => {
    this.props.showMessage(options)
  }

  onMessageClose = (messageId) => {
    this.props.closeMessage(messageId)
  }

  renderSpecificMessageType = (message) => {
    const SpecificComponent = message.messageType
    return (
      <SpecificComponent
        {...message.messageProps}
        messageId={message.id}
        messageType={message.messageType}
        onMessageClose={this.onMessageClose}
      />
    )
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
                {this.renderSpecificMessageType(message)}
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

export const Message = ({
  icon,
  sticky,
  title,
  content,
  classes,
  withTransactionHash,
  txHash,
  messageId,
  onMessageClose,
}) => {
  useEffect(() => {
    if (!sticky) {
      const timeout = setTimeout(
        () => onMessageClose(messageId),
        closeMessageTimeoutInMs
      )
      return () => clearTimeout(timeout)
    }
  }, [sticky, messageId, onMessageClose])

  return (
    <Banner>
      <div className="flex row">
        <Banner.Icon icon={icon} className={`${classes?.iconClassName} mr-1`} />
        <div style={styles.messageContentWrapper}>
          <Banner.Title>{title}</Banner.Title>
          {content && (
            <Banner.Description className={classes?.bannerDescription}>
              {content}
            </Banner.Description>
          )}
          {withTransactionHash && (
            <Banner.Action>
              <ViewInBlockExplorer
                type="tx"
                className="arrow-link"
                id={txHash}
              />
            </Banner.Action>
          )}
        </div>
        <Banner.CloseIcon onClick={() => onMessageClose(messageId)} />
      </div>
    </Banner>
  )
}

const styles = {
  messageContentWrapper: { minWidth: 0, flex: 1 },
}

export const useShowMessage = () => {
  const { showMessage } = useContext(MessagesContext)

  return showMessage
}

export const useCloseMessage = () => {
  const { closeMessage } = useContext(MessagesContext)

  return closeMessage
}
