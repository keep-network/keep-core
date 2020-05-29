import React from "react"
import { DataTable, Column } from "../DataTable"

const TerminologyDataTable = () => (
  <section className="tile" id="quick-terminology">
    <header className="flex row wrap mb-1">
      <h3 className="text-grey-70">Quick Terminology</h3>
      <a
        href="https://github.com/keep-network/keep-core/blob/master/docs/glossary.adoc"
        className="arrow-link"
        style={{ marginLeft: "auto", marginRight: "2rem" }}
      >
        Further Reading in GitHub
      </a>
    </header>
    <DataTable data={terminology}>
      <Column header="term" field="term" />
      <Column
        header="explanation"
        field="explanation"
        headerStyle={{ width: "65%" }}
      />
    </DataTable>
  </section>
)

export default TerminologyDataTable

const terminology = [
  {
    term: "Address",
    explanation: "A unique identifier that serves as a virtual location.",
  },
  {
    term: "Cliff",
    explanation:
      "The amount of time between when the grant is issued and when tokens start to vest and become unlocked.",
  },
  {
    term: "Delegation",
    explanation:
      "Setting aside an amount of KEEP to be staked by a trusted third party, referred to within the dApp as an operator.",
  },
  {
    term: "KEEP",
    explanation: "Keep Network’s native token and the token used to stake.",
  },
  {
    term: "Minimum Stake Amount",
    explanation:
      "The minimum stake amount as required by the staking smart contract.",
  },
  {
    term: "Signing Group",
    explanation:
      "The signing group that will produce the next relay entry candidate. If this group fails to respond to the request in time, the lead responsibility may shift to another group. Signing members are members of one complete signing group in the threshold relay.",
  },
  {
    term: "Stake",
    explanation:
      "An amount of KEEP that is bonded in order to participate in the threshold relay and, optionally, the Keep network. Part or all of this can be removed from escrow as penalties for misbehavior, while part or all of it can be refunded if and when a participant chooses to withdraw from the network and relay. There will be a period of delay called the undelegation period until the tokens are returned to the owner.",
  },
  {
    term: "Staking",
    explanation:
      "A mechanism by which token holders can support the decentralization of the KEEP network and earn rewards on their tokens. Reward earners are referred to within the dApp as a beneficiary.",
  },
  {
    term: "Token Grant",
    explanation:
      "A grant that contains KEEP tokens that unlocks at a set schedule over a period of time.",
  },
  {
    term: "Unlocking Schedule",
    explanation:
      "The amount of time it will take for KEEP tokens within a grant to be fully unlocked and available for release into the Owned token balance.",
  },
]
