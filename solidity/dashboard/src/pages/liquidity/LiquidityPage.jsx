import React from "react"
import PageWrapper from "../../components/PageWrapper"
import CardContainer from "../../components/CardContainer";
import Card from "../../components/Card";

const LiquidityPage = ({ title }) => {
  return (
    <PageWrapper title={title}>
      <CardContainer>
        <Card>KEEP + ETH</Card>
        <Card>KEEP + TBTC</Card>
        <Card>TBTC + ETH</Card>
      </CardContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
