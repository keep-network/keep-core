import React from "react"
import TokenAmount from "../TokenAmount"
import Button from "../Button"

const styles = {
  title: {
    marginBottom: "2.5rem",
  },
  tileTitle: {
    textAlign: "justify",
    marginBottom: "2.5rem",
  },
  addressWrapper: {
    marginBottom: "0.5rem",
  },
}

const CopyStakeStep2 = ({
  amount,
  beneficiary,
  operatorAddress,
  authorizerAddress,
}) => {
  return (
    <>
      <h2 style={styles.title}>Review your stake details below.</h2>
      <section className="tile">
        <h3 style={styles.tileTitle}>
          This stake balance will be copied to a newly upgraded delegation. The
          undelegation process will be started on the old staking contract.
        </h3>
        <TokenAmount
          amount={amount}
          currencySymbol={"KEEP"}
          wrapperClassName="mb-1"
        />
        <Address address={authorizerAddress} label="authorizer" />
        <Address address={operatorAddress} label="operator" />
        <Address address={beneficiary} label="beneficiary" />
      </section>
      <div className="flex row space-between self-end">
        <Button className="btn btn-transparent btn-lg mr-2">back</Button>
        <Button className="btn btn-primary btn-lg">copy stake</Button>
      </div>
    </>
  )
}

const Address = ({ label, address }) => (
  <div className="flex row center" style={styles.addressWrapper}>
    <h5 className="text-grey-70 flex-1">{label}</h5>
    <div className="text-big text-grey-50">{address}</div>
  </div>
)

export default CopyStakeStep2
