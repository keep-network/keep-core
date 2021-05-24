import React from "react"
import centered from "@storybook/addon-centered/react"
import TBTCRewardsDataTable from "../components/TBTCRewardsDataTable"

const mockedData = [
  {
    index: 0,
    amount: "200000000000000000000",
    proof: [
      "0x0112d0bca508715d464d0512dc14c3e32e008c8589087841983bafc9c2aea754",
      "0x6de0b077c9731dde9da1bef897d151490ae4f1a3b8e569cd3576f1b1911e5186",
      "0xd8efb31471d13dc0a8132dcbcebe8cc204e6a8721f29c4c2fa2789a5d15f2100",
    ],
    operator: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
    merkleRoot:
      "0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
    interval: 18,
    rewardsPeriod: "03/19/2021 - 03/26/2021",
    status: "AVAILABLE",
    id:
      "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756-0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
  },
  {
    index: 4,
    amount: "200000000000000000000",
    proof: [
      "0xb87d842718b981c0fa0b064594e3907209c2fabcd6298ef57ef7588f30d950f2",
      "0x5a92ce9e4026109734ad608859e6a45c621957b1ccf3d86eef107f61a85e5b48",
      "0xd8efb31471d13dc0a8132dcbcebe8cc204e6a8721f29c4c2fa2789a5d15f2100",
    ],
    operator: "0xf119557AC33585405467135eC9A343DCDb047517",
    merkleRoot:
      "0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
    interval: 18,
    rewardsPeriod: "03/19/2021 - 03/26/2021",
    status: "AVAILABLE",
    id:
      "0xf119557AC33585405467135eC9A343DCDb047517-0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
  },
  {
    index: 3,
    amount: "200000000000000000000",
    proof: [
      "0x29bd62772bbd6aa0d421a5ec1899b920480dd492f7cf1428b0fe892cdc1e706d",
      "0x6de0b077c9731dde9da1bef897d151490ae4f1a3b8e569cd3576f1b1911e5186",
      "0xd8efb31471d13dc0a8132dcbcebe8cc204e6a8721f29c4c2fa2789a5d15f2100",
    ],
    operator: "0xd2C6168Fd106908Df71Ab639f8b7e2F971Ab8205",
    merkleRoot:
      "0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
    interval: 18,
    rewardsPeriod: "03/19/2021 - 03/26/2021",
    status: "AVAILABLE",
    id:
      "0xd2C6168Fd106908Df71Ab639f8b7e2F971Ab8205-0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
  },
  {
    index: 2,
    amount: "200000000000000000000",
    proof: [
      "0x7a9cb5644eb8d5ab511b380271a20837a512eb0724770b3bdeb9222ef3ab4acc",
      "0x5a92ce9e4026109734ad608859e6a45c621957b1ccf3d86eef107f61a85e5b48",
      "0xd8efb31471d13dc0a8132dcbcebe8cc204e6a8721f29c4c2fa2789a5d15f2100",
    ],
    operator: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
    merkleRoot:
      "0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
    interval: 18,
    rewardsPeriod: "03/19/2021 - 03/26/2021",
    status: "AVAILABLE",
    id:
      "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953-0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
  },
  {
    index: 1,
    amount: "200000000000000000000",
    proof: [
      "0x1381ae8c1bd3aa152d29700ca51a9a225a545a48d46179038ee8b869281fa8ff",
    ],
    operator: "0xCDAfb5A23A1F1c6f80706Cc101BCcf4b9A1A3e3B",
    merkleRoot:
      "0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
    interval: 18,
    rewardsPeriod: "03/19/2021 - 03/26/2021",
    status: "AVAILABLE",
    id:
      "0xCDAfb5A23A1F1c6f80706Cc101BCcf4b9A1A3e3B-0x17f58452effed557fdf4c8763899a31c9d80f4a05334033c188a280702567929",
  },
]

export default {
  title: "TBTCRewardsDataTable",
  component: TBTCRewardsDataTable,
  decorators: [centered],
}

const Template = (args) => <TBTCRewardsDataTable {...args} />

export const EmptyTable = Template.bind({})
EmptyTable.args = {}

export const WithMockedData = Template.bind({})
WithMockedData.args = { data: mockedData}
