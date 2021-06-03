import winston from "winston"

const transports = () => {
  const transports = []

  transports.push(
    new winston.transports.Console({
      level: "debug",
      handleExceptions: true,

      format: winston.format.combine(
        winston.format.colorize({
          colors: { debug: "grey" },
        }),
        winston.format.simple(),
        winston.format.align()
      ),
    })
  )

  return transports
}

export const logger = winston.createLogger({
  format: winston.format.errors({ stack: true }),
  transports: transports(),
  exitOnError: false, // do not exit on handled exceptions
})
