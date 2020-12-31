const fs = require("fs")

let doc = `:toc: true
:toclevels: 2

= Keep Random Beacon API Documentation

Welcome to the Keep Random Beacon API Documentation. The primary contracts involved
are listed below, along with their public methods.

toc::[]
`

const jsonFiles = [
  "./build/contracts/KeepRandomBeaconServiceImplV1.json",
  "./build/contracts/KeepRandomBeaconOperator.json",
  "./build/contracts/KeepToken.json",
  "./build/contracts/TokenStaking.json",
  "./build/contracts/TokenGrant.json",
]

jsonFiles.forEach((file) => {
  const json = JSON.parse(fs.readFileSync(file, { encoding: "utf8" }))
  let section = "== `" + json.contractName + "`\n\n"

  for (let i = 0; i < json.devdoc.methods.length; i++) {
    const signature = json.devdoc.methods[i]
    const props = json.devdoc.methods[signature]

    let subsection = "=== `" + signature + "`\n\n"
    if (props.details) {
      subsection += `${props.details}\n\n`
    }

    if (props.params) {
      for (let j = 0; j < props.params.length; j++) {
        const paramName = props.params[j]
        const paramDoc = props.params[paramName]
        subsection += `\`${paramName}\`:: ` + paramDoc + "\n"
      }
    }

    if (props.return) {
      subsection += `Returns:: ${props["return"]}`
    }

    subsection += "\n\n"
    section += subsection
  }

  doc += section
})

console.log(doc)
