import React from "react";

export const centeredWithFullWidth = (Story) => (
  <div style={{display: "flex", overflow: "auto", alignItems: "center", inset: "0px", position: "fixed"}}>
    <div style={{margin: "auto", maxHeight: "100%", width: "100%"}}>
      <Story />
    </div>
  </div>
)
