import React from "react"

const options = [
  {
    id: 1,
    type: "COPY_TO_NEW_CONTRACT",
    title: "Copy stake balance to an upgraded delegation.",
    subtitle:
      "Avoid waiting the undelegation period and copy over your stake balance to a newly upgraded delegation. Your current stake will undelegate.",
  },
  {
    id: 2,
    type: "WAIT",
    title:
      "Undelegate and wait the undelegation period to stake to a newly upgraded delegation.",
    subtitle:
      "Undelegate and wait the 60 day undelegation period in order to stake on the upgraded staking contract. You can start new delegations with any tokens not already staked.",
  },
]

const CopyStakeStep1 = () => {
  return (
    <>
      <h2 className="text-center">
        Choose how to move your current stake to an upgraded delegation.
      </h2>
      <ul className="mt-2">
        {options.map((option) => (
          <Option key={option.id} {...option} />
        ))}
      </ul>
    </>
  )
}

const styles = {
  optionTile: {
    borderRadius: "10px",
  },
}

const Option = ({ title, subtitle, isSelected, type, id }) => {
  return (
    <li className="tile" style={styles.optionTile}>
      <input
        type="radio"
        name="option"
        value={type}
        id={`option-${type}-${id}`}
        checked={isSelected}
        onChange={() => console.log("onChnage")}
      />
      <label htmlFor={`option-${type}-${id}`}>
        <h3 className="text-grey-70 mb-1">{title}</h3>
        <p className="text-big text-grey-60 mb-0">{subtitle}</p>
      </label>
    </li>
  )
}

export default CopyStakeStep1
