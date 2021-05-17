import React from "react"
import { TopUpsDataTable } from "../components/TopUpsDataTable"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import Tile from "../components/Tile"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import { LoadingOverlay } from "../components/Loadable"

storiesOf("TopUpsDataTable", module).addDecorator(centered)

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
  ],
}

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

const Template = (args) => <TopUpsDataTable {...args} />

export const Default = Template.bind({})
Default.args = {
  topUps: mockedTopUps,
  commitTopUp: null,
}
