// import React from "react"
// import Redirect from "../components/Redirect"
// import { MemoryRouter, Route } from "react-router-dom"
//
// export default {
//   title: "Redirect",
//   component: Redirect,
//   decorators: [
//     (Story) => {
//       ;<MemoryRouter initialEntries={["/liquidity"]}>
//         <Route path={"/liquidity"}>
//           <Story />
//         </Route>
//       </MemoryRouter>
//     },
//   ],
// }
//
// const Template = (args) => <Redirect {...args} />
//
// export const RedirectToLiquidity = Template.bind({})
// RedirectToLiquidity.args = {
//   to: "/liquidity",
// }
