import React, { useEffect } from "react"
import List from "../../components/List"
import TokenAmount from "../../components/TokenAmount"
import {
  MigrationPortalForm,
  ConfirmMigrationModal,
} from "../../components/tbtc-migration"
import { TBTC } from "../../utils/token.utils"
import { useModal } from "../../hooks/useModal"
import { useWeb3Address } from "../../components/WithWeb3Context"
import { useDispatch, useSelector } from "react-redux"
import { tbtcV2Migration } from "../../actions"

const TokenUpgradePortalPage = () => {
  const { openConfirmationModal } = useModal()
  const address = useWeb3Address()
  const dispatch = useDispatch()

  const { tbtcV1Balance, tbtcV2Balance, unmintFee } = useSelector(
    (state) => state.tbtcV2Migration
  )

  useEffect(() => {
    dispatch(tbtcV2Migration.fetchDataRequest(address))
  }, [address, dispatch])

  const onSubmitMigrationForm = async (values, awaitingPromise) => {
    const { amount, from, to } = values
    const _amount = TBTC.fromTokenUnit(amount).toString()

    await openConfirmationModal(
      {
        modalOptions: { title: to === "v2" ? "Upgrade" : "Downgrade" },
        from,
        to,
        amount: _amount,
        fee: unmintFee,
      },
      ConfirmMigrationModal
    )

    if (to === "v2") {
      dispatch(tbtcV2Migration.mint(_amount, awaitingPromise))
    } else {
      dispatch(tbtcV2Migration.unmint(_amount, awaitingPromise))
    }
  }

  return (
    <section className="tbtc-migration-portal">
      <List className="tbtc-migration-portal__tbtc-balances">
        <List.Title className="h3 text-grey-70">Balance</List.Title>
        <List.Content className="tbtc-balances">
          <List.Item className="tbtc-balance tbtc-balance--v1">
            <TokenAmount
              token={TBTC}
              amount={tbtcV1Balance}
              symbol="tBTC v1"
              amountClassName="h2 text-white"
              symbolClassName="h3 text-white"
              withIcon
            />
          </List.Item>
          <List.Item className="tbtc-balance tbtc-balance--v2">
            <TokenAmount
              token={TBTC}
              amount={tbtcV2Balance}
              symbol="tBTC v2"
              amountClassName="h2 text-black"
              symbolClassName="h3 text-black"
              withIcon
            />
          </List.Item>
        </List.Content>
      </List>
      <section className="tbtc-migration-portal__form-wrapper">
        <h3 className="text-grey-70 mb-1">Migration Portal</h3>
        <MigrationPortalForm
          mintingFee={unmintFee}
          tbtcV1Balance={tbtcV1Balance}
          tbtcV2Balance={tbtcV2Balance}
          onSubmit={onSubmitMigrationForm}
        />
      </section>
    </section>
  )
}

TokenUpgradePortalPage.route = {
  title: "Token Upgrade Portal",
  path: "/tbtc-migration/portal",
  exact: true,
}

export default TokenUpgradePortalPage
