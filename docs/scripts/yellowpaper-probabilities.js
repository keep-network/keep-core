const asciidoctor = require('asciidoctor.js')()
const math = require('mathjs')

const adoc = asciidoctor.loadFile('../random-beacon.adoc')

const varDefs =
    adoc.findBy({ role: 'variables' })[0]
        .rows
        .body
        .map(row => ({ variable: row[0].text.trim(), meaning: row[1].text.trim() }));
const derivedVarDefs =
    adoc.findBy({ role: 'derived-variables' })[0]
        .rows
        .body
        .map(row => ({ variable: row[0].text.trim(), meaning: row[2].text.trim() }));

const variables =
    varDefs.concat(derivedVarDefs)
        .reduce((map, { variable, meaning }) => map.set(variable, meaning), new Map())

const probabilities =
    adoc.findBy({ role: 'probabilities' })[0]
        .rows
        .body
        .map(row => ({ label: row[0].text.trim(), formula: row[1].text.trim() }))

function hgeo(N, K, n, k) {
    return math.combinations(math.bignumber(K), math.bignumber(k))
        .mul(math.combinations(math.bignumber(N).minus(math.bignumber(K)), math.bignumber(n).minus(math.bignumber(k))))
        .div(math.combinations(math.bignumber(N), math.bignumber(n)));
}

const C = 150000,
      G = 1000,
      g = 1000,
      t = math.floor(g / 2) - 1,
      M = 10000,
      B = 14,
      N = 8000000,
      D = 75000,
      R = 1,
      B_s = 8000000;

console.log("Given:");

for (variable of ['C', 'G', 'g', 't', 'M', 'B', 'N', 'D', 'R', 'B_s']) {
    let extra = "" 
    if (variables.has(variable)) {
        extra = "\t(" + variables.get(variable) + ")"
    }

    console.log("\t" + variable + " = " + eval(variable) + extra);
}

console.log("\nThe probability of:")

probabilities.forEach(function({ label, formula }) {
    const probabilityResult = math.bignumber(eval(formula))

    const probabilityString = probabilityResult.toExponential(4, math.ROUND_DOWN ).toString()

    //console.log(math.bignumber(1).minus(probabilityResult).pow(N).toString())

    console.log(
        "\t" + label.split("\n").join("\n\t") + "\n\t\t",
        probabilityString
    )
})