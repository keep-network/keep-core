import React from "react"
import * as Icons from "../Icons"
import Button from "../Button"
import moment from "moment"

const content = {
  WAIT_FLOW: {
    title:
      "Success! The undelegation process for your stake balance has started.",
    subtitles: [
      "Your tokens will be available in undelegationCompletedAt",
      "Once available, hit recover to transfer to your balance.",
      "In the meantime, continue staking unstaked tokens. Any new stakes will now use the upgraded staking contract.",
    ],
  },
  WAIT_FLOW_RECOVER: {
    title: "Success! The recovery process for your stake balance has finished.",
    subtitles: [
      "Your tokens are available.",
      "You can stake tokens to the upgraded staking contract.",
    ],
  },
  COPY_STAKE_FLOW: {
    title: "Success! Your stake balance copied and redelegated.",
    subtitles: [
      "Your former stake will be available in undelegationCompletedAt.",
      "You’ll need to initiate the recovery process in the dashboard.",
      "You’ll see a notification in the dashboard when it’s time to do this.",
    ],
  },
}

const styles = {
  title: { textAlign: "left" },
  listItem: { marginBottom: "1.25rem" },
  iconWrapper: {
    backgroundColor: "black",
    padding: "0.5rem",
    borderRadius: "50%",
  },
}

const CopyStakeStep4 = ({
  onClose,
  strategy,
  undelegationPeriod,
  selectedDelegation,
}) => {
  const getUndelegationCompleteAt = () => {
    if (selectedDelegation && selectedDelegation.undelegationCompleteAt) {
      return selectedDelegation.undelegationCompleteAt.fromNow(true)
    }

    return moment().add(undelegationPeriod, "seconds").fromNow(true)
  }

  const getContent = () => {
    if (strategy === "WAIT_FLOW" && selectedDelegation.canRecoverStake) {
      return content.WAIT_FLOW_RECOVER
    } else {
      return content[strategy]
    }
  }

  return (
    <>
      <div className="self-start flex center" style={styles.iconWrapper}>
        <Icons.Success />
      </div>
      <h2 className="mb-2 mt-2 self-start" style={styles.title}>
        {getContent().title}
      </h2>
      <ul className="list__colored-bullets--grey-60">
        {getContent().subtitles.map((subtitle, index) => (
          <li key={index} className="h3 text-grey-70" style={styles.listItem}>
            {subtitle.replace(
              "undelegationCompletedAt",
              getUndelegationCompleteAt()
            )}
          </li>
        ))}
      </ul>
      <Button className="btn btn-primary btn-lg mr-a mt-1" onClick={onClose}>
        close
      </Button>
    </>
  )
}

export default CopyStakeStep4
