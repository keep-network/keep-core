import React from "react"
import { withBaseModal } from "../withBaseModal"
import { ModalBody, ModalFooter, ModalHeader } from "../Modal"
import Button, { SubmitButton } from "../../Button"
import { FormCheckboxBase } from "../../FormCheckbox"
import { useFormik } from "formik"
import useReleaseTokens from "../../../hooks/useReleaseTokens"
import {
  Accordion,
  AccordionItem,
  AccordionItemButton,
  AccordionItemHeading,
  AccordionItemPanel,
} from "react-accessible-accordion"
import { KEEP } from "../../../utils/token.utils"
import TokenAmount from "../../TokenAmount"
import { shortenAddress } from "../../../utils/general.utils"
import moment from "moment"

export const WithdrawGrantedTokens = withBaseModal(({ grants, onClose }) => {
  const releaseTokens = useReleaseTokens()

  console.log("grants", grants)

  const onWithdrawClick = (awaitingPromise) => {
    for (const grantId of formik.values.grantIds) {
      const selectedTokenGrant = grants.find((grant) => grant.id === grantId)
      releaseTokens(selectedTokenGrant, awaitingPromise)
    }
    onClose()
  }

  const formik = useFormik({
    enableReinitialize: true,
    initialValues: {
      grantIds: [],
    },
    onSubmit: (values) => {
      alert(JSON.stringify(values, null, 2))
    },
  })

  return (
    <>
      <ModalHeader>Withdraw</ModalHeader>
      <ModalBody>
        <h3 className={"mb-1"}>Withdraw granted tokens</h3>
        <p className={"mb-2"}>
          The following tokens are available to be released and withdrawn from
          your token grant.
        </p>
        <div>
          <form>{grants.map((grant) => renderGrant(grant, formik))}</form>
        </div>
      </ModalBody>
      <ModalFooter>
        <SubmitButton
          className="btn btn-primary btn-lg mr-2"
          type="submit"
          onSubmitAction={(awaitingPromise) => {
            onWithdrawClick(awaitingPromise)
          }}
          disabled={formik.values.grantIds.length === 0}
        >
          withdraw
        </SubmitButton>
        <Button className={`btn btn-unstyled text-link`} onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </>
  )
})

const renderGrant = (grant, formik) => (
  <div
    key={`Grant-${grant.id}`}
    style={{
      display: "flex",
      alignItems: "baseline",
    }}
  >
    <FormCheckboxBase
      className="withdraw-granted-tokens__form_checkbox"
      name="grantIds"
      type="checkbox"
      value={grant.id}
      checked={formik.values.grantIds.includes(grant.id)}
      onChange={(e) => {
        const set = new Set(formik.values.grantIds)
        if (e.target.checked) {
          set.add(grant.id)
        } else {
          set.delete(grant.id)
        }
        formik.setFieldValue("grantIds", [...set])
      }}
    />
    <Accordion
      allowZeroExpanded
      className={"withdraw-granted-tokens__grants-accordion"}
    >
      <AccordionItem>
        <AccordionItemHeading>
          <AccordionItemButton>
            <div className="withdraw-granted-tokens__token-amount">
              <TokenAmount
                amount={grant.amount}
                amountClassName={"h4 text-mint-100"}
                symbolClassName={"h4 text-mint-100"}
                token={KEEP}
              />
            </div>
            <span className="withdraw-granted-tokens__details-text">
              Details
            </span>
          </AccordionItemButton>
        </AccordionItemHeading>
        <AccordionItemPanel>
          <div className="withdraw-granted-tokens__accordion-info-row">
            <span className="withdraw-granted-tokens__info-row-title">
              token grant id
            </span>
            <span className="withdraw-granted-tokens__info-row-value">
              {grant.id}
            </span>
          </div>
          <div className="withdraw-granted-tokens__accordion-info-row">
            <span className="withdraw-granted-tokens__info-row-title">
              date issued
            </span>
            <span className="withdraw-granted-tokens__info-row-value">
              {moment.unix(grant.start).format("MM/DD/YYYY")}
            </span>
          </div>
          <div className="withdraw-granted-tokens__accordion-info-row">
            <span className="withdraw-granted-tokens__info-row-title">
              wallet
            </span>
            <span className="withdraw-granted-tokens__info-row-value">
              {shortenAddress(grant.grantee)}
            </span>
          </div>
        </AccordionItemPanel>
      </AccordionItem>
    </Accordion>
  </div>
)
