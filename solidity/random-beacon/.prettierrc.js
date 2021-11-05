module.exports = {
  ...require("@keep-network/prettier-config-keep"),
  overrides: [
    {
      files: "*.sol",
      options: {
        tabWidth: 4,
      },
    },
  ],
};
