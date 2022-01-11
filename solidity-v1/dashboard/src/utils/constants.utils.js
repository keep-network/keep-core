export const renderDynamicConstant = (constant, ...args) => {
  let finalConst = constant
  const vars = constant.split("${")
  vars.shift()
  let index = 0
  for (const v of vars) {
    if (!v.includes("}")) break
    if (index > args.length - 1) {
      throw new Error("Too few arguments were given!")
    }
    const temp = v.split("}")[0]
    finalConst = finalConst.replace("${" + temp + "}", args[index])
    index++
  }
  return finalConst
}
