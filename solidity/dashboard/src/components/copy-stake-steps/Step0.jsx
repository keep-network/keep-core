import React from "react"
import Button from "../Button"

const styles = {
  title: { marginTop: "5rem", marginBottom: "1.5rem" },
  subtitle: { textAlign: "justify" },
}

const CopyStakeStepO = () => {
  return (
    <>
      <h2 style={styles.title}>The Keep staking contract is brand new!</h2>
      <h3 className="mb-1" style={styles.subtitle}>
        To continue running smoothly on the Keep network, any stake that was
        delegated to the previous contract version will need to move to the
        upgraded staking contract.
        <br />
        <br />
        This process will guide you through the quick upgrade.
      </h3>
      <Button className="btn btn-primary btn-lg mt-2">continue</Button>
    </>
  )
}

export default CopyStakeStepO
