import React, { useCallback } from "react"
import { useDispatch } from "react-redux"
import web3Utils from "web3-utils"

import PageWrapper from "../../components/PageWrapper"
import Tile from "../../components/Tile"

import { MintBondTokensFormik } from "./MintBondTokens"
import { TransferKeepTokensFormik } from "./TransferKeepTokens"

const DebugMintingPage = (props) => {
  const dispatch = useDispatch()

  const handleMint = useCallback((amount, address, meta) => {
    dispatch({
      type: 'debug-minting/mint-bondTokens',
      payload: { amount, address },
      meta
    })
  }, [dispatch])

  const onBondSubmit = useCallback(
    async (formValues, awaitingPromise) => {
      const { amount, address } = formValues
      const weiAmount = web3Utils.toWei(amount.toString(), "ether")

      handleMint(weiAmount, address, awaitingPromise)
    },
    [handleMint]
  )

  const handleTransferKeeps = useCallback((amount, address, meta) => {
    dispatch({
      type: 'debug-minting/transfer-keep-tokens',
      payload: { amount, address },
      meta
    })
  }, [dispatch])

  const onKeepSubmit = useCallback(
    async (formValues, awaitingPromise) => {
      const { amount, address } = formValues
      const weiAmount = web3Utils.toWei(amount.toString(), "ether")

      handleTransferKeeps(weiAmount, address, awaitingPromise)
    },
    [handleTransferKeeps]
  )

  return <>
    <PageWrapper {...props} >
      <Tile title="Mint BondERC20 tokens" titleClassName="h2 mb-2">
        <MintBondTokensFormik onSubmit={onBondSubmit} />
      </Tile>

      <Tile title="Transfer KEEP tokens" titleClassName="h2 mb-2">
        <TransferKeepTokensFormik onSubmit={onKeepSubmit} />
      </Tile>
    </PageWrapper>
  </>
}

DebugMintingPage.route = {
  title: "Debug Minting",
  path: "/minting",
  exact: true,
}

export default DebugMintingPage
