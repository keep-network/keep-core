import React, { useMemo, useState } from "react"
import { withBaseModal } from "../withBaseModal"
import { ModalBody, ModalFooter, ModalHeader } from "../Modal"
import Button, { SubmitButton } from "../../Button"
import useReleaseTokens from "../../../hooks/useReleaseTokens"
import {
  Accordion,
  AccordionItem,
  AccordionItemButton,
  AccordionItemHeading,
  AccordionItemPanel,
  AccordionItemState,
} from "react-accessible-accordion"
import { KEEP } from "../../../utils/token.utils"
import TokenAmount from "../../TokenAmount"
import { shortenAddress } from "../../../utils/general.utils"
import moment from "moment"
import { add } from "../../../utils/arithmetics.utils"
import OnlyIf from "../../OnlyIf"
import List from "../../List"
import * as Icons from "../../Icons"

export const WithdrawGrantedTokens = withBaseModal(({ grants, onClose }) => {
  const DEFAULT_GRANTS_TO_DISPLAY = 5
  const [numberOfGrantsDisplayed, setNumberOfGrantsDisplayed] = useState(
    DEFAULT_GRANTS_TO_DISPLAY
  )
  const [selectedGrant, setSelectedGrant] = useState(
    grants.length === 1 ? grants[0].id : null
  )
  const releaseTokens = useReleaseTokens()

  const totalReadyToRelease = useMemo(() => {
    return grants
      .map((grant) => grant.readyToRelease)
      .reduce((previous, current) => add(previous, current), 0)
      .toString()
  }, [grants])

  const onWithdrawClick = (awaitingPromise) => {
    releaseTokens(selectedGrant, awaitingPromise)
  }

  return (
    <>
      <ModalHeader>Withdraw</ModalHeader>
      <ModalBody>
        <h3 className={"mb-1"}>Withdraw granted tokens</h3>
        <p className={"mb-2 text-grey-70"}>
          The following tokens are available to be released and withdrawn from
          your token grant.
        </p>
        <div>
          <OnlyIf condition={grants.length > 1}>
            <div className="flex row center space-between">
              <h4 className="text-grey-70">Withdraw:</h4>
              <span className="text-grey-70 small">
                Total:{" "}
                <TokenAmount
                  amount={totalReadyToRelease}
                  token={KEEP}
                  wrapperClassName="withdraw-granted-tokens__total-ready-to-release-amount"
                  amountClassName="text-grey-60"
                  symbolClassName="text-grey-60"
                />
              </span>
            </div>
          </OnlyIf>
          <form>
            <Accordion
              allowZeroExpanded={grants.length > 1}
              className={
                "withdraw-granted-tokens__grants-accordion withdraw-granted-tokens__grants-accordion-container"
              }
              preExpanded={grants.length === 0 ? [`grant-${grants[0].id}`] : []}
            >
              {grants
                .slice(0, numberOfGrantsDisplayed)
                .map((grant) =>
                  renderGrant(
                    grant,
                    grants.length,
                    selectedGrant,
                    setSelectedGrant
                  )
                )}
            </Accordion>
          </form>
          <OnlyIf condition={grants.length > numberOfGrantsDisplayed}>
            <Button
              className="withdraw-granted-tokens__view-more-button"
              onClick={() =>
                setNumberOfGrantsDisplayed(
                  numberOfGrantsDisplayed + DEFAULT_GRANTS_TO_DISPLAY
                )
              }
            >
              <Icons.Add />
              &nbsp;View More
            </Button>
          </OnlyIf>
        </div>
      </ModalBody>
      <ModalFooter>
        <SubmitButton
          className="btn btn-primary btn-lg mr-2"
          type="submit"
          onSubmitAction={onWithdrawClick}
          disabled={grants.length > 1 && !selectedGrant}
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

const renderGrant = (
  grant,
  totalNumberOfGrants,
  selectedGrant,
  setSelectedGrant
) => {
  return (
    <>
      <AccordionItem uuid={`grant-${grant.id}`}>
        <input
          className="radio-without-label"
          type="radio"
          name="selectedGrantId"
          value={grant.id}
          id={`grant-${grant.id}`}
          checked={selectedGrant?.id === grant?.id}
          onChange={() => {
            setSelectedGrant(grant)
          }}
        />
        <label htmlFor={`grant-${grant.id}`} />
        <AccordionItemState>
          {({ expanded }) => {
            return (
              <AccordionItemHeading
                className={`accordion__heading ${
                  expanded ? "accordion__heading--without-bottom-border" : ""
                }`}
              >
                <AccordionItemButton>
                  <TokenAmount
                    amount={grant.readyToRelease}
                    amountClassName={"h4 text-mint-100"}
                    symbolClassName={"h4 text-mint-100"}
                    token={KEEP}
                  />
                  <OnlyIf condition={totalNumberOfGrants > 1}>
                    <span className="withdraw-granted-tokens__details-text">
                      Details
                    </span>
                    {expanded ? <Icons.ChevronUp /> : <Icons.ChevronDown />}
                  </OnlyIf>
                </AccordionItemButton>
              </AccordionItemHeading>
            )
          }}
        </AccordionItemState>
        <AccordionItemPanel>
          <List>
            <List.Content>
              <List.Item className="flex row center space-between">
                <span className="text-label text-label--without-hover text-grey-70">
                  token grant id
                </span>
                <span className="text-grey-60 small">{grant.id}</span>
              </List.Item>
              <List.Item className="flex row center space-between">
                <span className="text-label text-label--without-hover text-grey-70">
                  date issued
                </span>
                <span className="text-grey-60 small">
                  {moment.unix(grant.start).format("MM/DD/YYYY")}
                </span>
              </List.Item>
              <List.Item className="flex row center space-between">
                <span className="text-label text-label--without-hover text-grey-70">
                  wallet
                </span>
                <span className="text-grey-60 small">
                  {shortenAddress(grant.grantee)}
                </span>
              </List.Item>
            </List.Content>
          </List>
        </AccordionItemPanel>
      </AccordionItem>
    </>
  )
}
