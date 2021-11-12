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

export const WithdrawGrantedTokens = withBaseModal(({ grants, onClose }) => {
  const releaseTokens = useReleaseTokens()

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
    style={{
      display: "flex",
    }}
  >
    <FormCheckboxBase
      key={`Grant-${grant.id}`}
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
      withoutLabel
    >
      <Accordion allowZeroExpanded>
        <AccordionItem>
          <AccordionItemHeading>
            <AccordionItemButton>Grant #{grant.id} details</AccordionItemButton>
          </AccordionItemHeading>
          <AccordionItemPanel>SuperTest</AccordionItemPanel>
        </AccordionItem>
      </Accordion>
    </FormCheckboxBase>
  </div>
)
