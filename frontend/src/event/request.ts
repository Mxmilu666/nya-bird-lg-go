const eventList = ['Unauthorized', 'NetworkError', 'UnknownError', 'Message'] as const

export type EventList = (typeof eventList)[number]

// 定义每个事件的类型
export interface EventPayloads {
  NetworkError: [statusResult: string | boolean]
  UnknownError: []
  Message: [type: 'success' | 'error' | 'warn', msg: string]
  Unauthorized: unknown[]
}

export type EventCallback<E extends EventList = EventList> = 
  (...args: EventPayloads[E]) => void

const _listenClient = new Map<string, Set<EventCallback<EventList>>>()

export const on = <E extends EventList>(event: E, callback: EventCallback<E>) => {
  if (!_listenClient.has(event)) {
    _listenClient.set(event, new Set<EventCallback<EventList>>())
  }
  _listenClient.get(event)?.add(callback as EventCallback<EventList>)
}

export const emit = <E extends EventList>(event: E, ...args: EventPayloads[E]) => {
  _listenClient.get(event)?.forEach((callback) => {
    callback(...args as unknown as EventPayloads[EventList])
  })
}

const requestEvent = {
  on,
  emit,
  eventList
}

export default requestEvent