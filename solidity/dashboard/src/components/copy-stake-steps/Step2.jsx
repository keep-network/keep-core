import React from "react"
import Button from "../Button"
import Badge from "../Badge"

const options = [
  {
    id: 1,
    type: "COPY_STAKE_FLOW",
    title: "Copy stake balance to an upgraded delegation.",
    subtitle:
      "Avoid waiting the undelegation period and copy over your stake balance to a newly upgraded delegation. Your current stake will undelegate.",
    recomended: true,
  },
  {
    id: 2,
    type: "WAIT_FLOW",
    title:
      "Undelegate and wait the undelegation period to stake to a newly upgraded delegation.",
    subtitle:
      "Undelegate and wait the 60 day undelegation period in order to stake on the upgraded staking contract. You can start new delegations with any tokens not already staked.",
    recomended: false,
  },
]

const CopyStakeStep2 = ({
  incrementStep,
  decrementStep,
  selectedStrategy,
  setStrategy,
}) => {
  const onSetStrategy = (event) => {
    setStrategy(event.target.value)
  }

  return (
    <>
      <h2 className="text-center">
        Choose how to move your current stake to an upgraded delegation.
      </h2>
      <ul className="mt-2">
        {options.map((option) => (
          <Option
            key={option.id}
            {...option}
            isSelected={option.type === selectedStrategy}
            onChange={onSetStrategy}
          />
        ))}
      </ul>
      <div className="flex row space-between self-end">
        <Button
          className="btn btn-transparent btn-lg mr-2"
          onClick={decrementStep}
        >
          back
        </Button>
        <Button
          className="btn btn-primary btn-lg"
          onClick={incrementStep}
          disabled={!selectedStrategy}
        >
          review stake
        </Button>
      </div>
    </>
  )
}

const styles = {
  optionTile: {
    borderRadius: "10px",
  },
  recomendeLabel: {
    marginLeft: "auto",
    marginTop: "-1.5rem",
    marginRight: "-2rem",
  },
}

const Option = ({
  title,
  subtitle,
  isSelected,
  type,
  id,
  onChange,
  recomended,
}) => {
  return (
    <li className="tile" style={styles.optionTile}>
      {recomended && (
        <div className="flex flex-1">
          <Badge text="recomended" style={styles.recomendeLabel} />
        </div>
      )}
      <input
        type="radio"
        name="option"
        value={type}
        id={`option-${type}-${id}`}
        checked={isSelected}
        onChange={onChange}
      />
      <label htmlFor={`option-${type}-${id}`}>
        <h3 className="text-grey-70 mb-1">{title}</h3>
        <p className="text-big text-grey-60 mb-0">{subtitle}</p>
      </label>
    </li>
  )
}

export default CopyStakeStep2
