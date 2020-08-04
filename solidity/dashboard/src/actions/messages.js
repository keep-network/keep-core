export const ADD_MESSAGE = "ADD_MESSAGE"
export const REMOVE_MESSAGE = "REMOVE_MESSAGE"
export const CLOSE_MEESSAGE = "CLOSE_MEESSAGE"
export const SHOW_MESSAGE = "SHOW_MESSAGE"

let messageId = 0

export class Message {
  static create(options) {
    const { title, content, type, sticky } = options

    return new Message(title, content, type, sticky)
  }

  constructor(title, content, type, sticky = false) {
    this.id = messageId++
    this.title = title
    this.content = content
    this.type = type
    this.sticky = sticky
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
    payload: { ...options },
  }
}

export const closeMessage = (id) => {
  return {
    type: CLOSE_MEESSAGE,
    payload: id,
  }
}
