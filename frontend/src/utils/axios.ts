import axios from 'axios'
import requestEvent from '@/event/request'

axios.defaults.baseURL = import.meta.env.VITE_HTTP_BASE_URL || '/api'

// http request 拦截器
axios.interceptors.request.use(
  (config) => {
    return config
  },
  (err) => {
    return Promise.reject(err)
  }
)

// http response 拦截器
axios.interceptors.response.use(
  (response) => {
    if (response.status >= 500 && response.status <= 599) {
      requestEvent.emit('UnknownError')
    }

    if (response.status === 418) {
      requestEvent.emit('Message', response.data?.type, response.data?.msg)
    }

    return response
  },
  (err) => {
    console.log('request error', err)
    requestEvent.emit('UnknownError')
    return Promise.reject(err)
  }
)

export type Response<T = any> = { status: number; msg: string; data: T }

export default axios