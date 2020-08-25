import React, { useEffect } from "react"
import TokenAmount from "../TokenAmount"
import Button from "../Button"
import { KeepLoadingIndicator } from "../Loadable"

const styles = {
  title: {
    marginBottom: "1rem",
  },
  tileTitle: {
    textAlign: "justify",
    marginBottom: "2rem",
  },
  addressWrapper: {
    marginBottom: "0.5rem",
  },
  delegationsList: { width: "100%" },
  delegationItem: {
    borderRadius: "10px",
  },
}

const CopyStakeStep1 = ({
  fetchDelegations,
  isFetching,
  delegations,
  decrementStep,
  incrementStep,
  onSelectDelegation,
  selectedDelegation,
}) => {
  useEffect(() => {
    fetchDelegations()
  }, [fetchDelegations])

  console.log("selected delegation", selectedDelegation)
  return (
    <>
      <h2 style={styles.title}>Stake balances to be upgraded.</h2>
      <h3 className="text-grey-70" style={styles.tileTitle}>
        Choose the stake delegations to be upgraded. You can select one or
        multiple stake delegations.
      </h3>
      <ul style={styles.delegationsList}>
        {isFetching ? (
          <KeepLoadingIndicator />
        ) : (
          delegations.map((delegation) => (
            <DelegationItem
              key={delegation.operatorAddress}
              delegation={delegation}
              onSelect={onSelectDelegation}
              isSelected={
                selectedDelegation &&
                selectedDelegation.operatorAddress ===
                  delegation.operatorAddress
              }
            />
          ))
        )}
      </ul>
      <div className="flex row space-between self-end">
        <Button
          className="btn btn-transparent btn-lg mr-2"
          onClick={decrementStep}
        >
          back
        </Button>
        <Button
          disabled={!selectedDelegation}
          className="btn btn-primary btn-lg"
          onClick={incrementStep}
        >
          continue
        </Button>
      </div>
    </>
  )
}

const DelegationItem = ({ delegation, onSelect }) => {
  const {
    amount,
    authorizerAddress,
    operatorAddress,
    beneficiary,
    isSelected,
  } = delegation

  return (
    <li className="tile" style={styles.delegationItem}>
      <input
        type="radio"
        name="option"
        value={delegation}
        id={`option-${operatorAddress}`}
        checked={isSelected}
        onChange={() => onSelect(delegation)}
      />
      <label htmlFor={`option-${operatorAddress}`} style={{ width: "100%" }}>
        <div className="flex row">
          <TokenAmount
            amount={amount}
            currencySymbol="KEEP"
            wrapperClassName="self-start"
            amountClassName="text-primary h3"
          />
          <div style={{ marginLeft: "auto" }}>
            <Address address={authorizerAddress} label="authorizer" />
            <Address address={operatorAddress} label="operator" />
            <Address address={beneficiary} label="beneficiary" />
          </div>
        </div>
      </label>
    </li>
  )
}

const Address = ({ label, address }) => (
  <div className="flex row center" style={styles.addressWrapper}>
    <h5 className="text-grey-70 flex-1">{label}</h5>
    <div className="text-big text-grey-50">{address}</div>
  </div>
)

export default CopyStakeStep1
