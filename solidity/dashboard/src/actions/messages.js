export const ADD_MESSAGE = "ADD_MESSAGE"
export const REMOVE_MESSAGE = "REMOVE_MESSAGE"
export const CLOSE_MEESSAGE = "CLOSE_MEESSAGE"
export const SHOW_MESSAGE = "SHOW_MESSAGE"

let messageId = 1

export class Message {
  static create(options) {
    return new Message(options)
  }

  constructor(options) {
    Object.assign(this, options)
    this.id = messageId++
  }

  id
  title
  content
  type
  sticky
}

export const showMessage = (options) => {
  return {
    type: SHOW_MESSAGE,
    payload: Message.create(options),
  }
}

export const showCreatedMessage = (message) => {
  return {
    type: SHOW_MESSAGE,
    payload: message,
  }
}

export const closeMessage = (id) => {
  return {
    type: CLOSE_MEESSAGE,
    payload: id,
  }
}
