const eventList = ['Unauthorized', 'NetworkError', 'UnknownError', 'Message'] as const

export type EventList = (typeof eventList)[number]

export type EventCallback = (...args: unknown[]) => void

const _listenClient = new Map<string, Set<EventCallback>>()

export const on = (event: EventList, callback: EventCallback) => {
  if (!_listenClient.has(event)) _listenClient.set(event, new Set<EventCallback>())
  _listenClient.get(event)?.add(callback)
}

export const emit = (event: EventList, ...args: unknown[]) => {
  _listenClient.get(event)?.forEach((callback) => {
    callback(...args)
  })
}

const requestEvent = {
  on,
  emit,
  eventList
}

export default requestEvent