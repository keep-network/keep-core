module.exports = {
  ...require("@thesis-co/prettier-config"),
  overrides: [
    {
      files: "*.sol",
      options: {
        tabWidth: 4,
      },
    },
  ],
}
