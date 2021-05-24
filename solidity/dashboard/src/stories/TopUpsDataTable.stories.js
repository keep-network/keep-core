import React from "react"
import { TopUpsDataTable } from "../components/TopUpsDataTable"
import Tile from "../components/Tile"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import { LoadingOverlay } from "../components/Loadable"
import centered from "@storybook/addon-centered/react"

const mockedTopUps = [
  {
    availableTopUpAmount: "2000000000000000000000",
    createdAt: 1620803413,
    IsInUndelegation: false,
    operatorAddress: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
  },
  {
    availableTopUpAmount: "2000000000000000000000",
    createdAt: 1620803413,
    IsInUndelegation: false,
    operatorAddress: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
  },
]

export default {
  title: "TopUpsDataTable",
  component: TopUpsDataTable,
  decorators: [
    (Story) => (
      <section>
        <LoadingOverlay
          isFetching={false}
          skeletonComponent={<DataTableSkeleton columns={3} />}
        >
          <Tile>
            <Story />
          </Tile>
        </LoadingOverlay>
      </section>
    ),
    centered,
  ],
}

const Template = (args) => <TopUpsDataTable {...args} />

export const EmptyTable = Template.bind({})
EmptyTable.args = {
  topUps: [],
  commitTopUp: null,
}

export const WithMockedData = Template.bind({})
WithMockedData.args = {
  topUps: mockedTopUps,
  commitTopUp: null,
}
