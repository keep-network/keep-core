import "./../src/css/app.css"
import "./storybookStyleFix.css"
import React from "react";
import { BrowserRouter } from "react-router-dom";
export const parameters = {
  layout: 'centered',
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
