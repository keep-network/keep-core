import "./../src/css/app.css"
import "./storybookStyleFix.css"
import React from "react";
import { BrowserRouter } from "react-router-dom";
import {addDecorator, storiesOf} from "@storybook/react";
import centered from "@storybook/addon-centered/react";
export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
}

export const decorators = [
  (Story) => (
    <>
      <div></div>
      <BrowserRouter>
        <Story />
      </BrowserRouter>
    </>
  ),
];

// addDecorator((...args) => {
//   const params = (new URL(document.location)).searchParams;
//   const isInDockView = params.get('viewMode') === 'docs';
//
//   if (isInDockView) {
//     return args[0]();
//   }
//
//   return centered(...args);
// });



// storiesOf("KeepOnlyPool", module).addDecorator(customCentered)