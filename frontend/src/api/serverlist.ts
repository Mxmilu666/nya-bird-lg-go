import axios, { type Response } from '@/utils/axios'

// 获取所有可用服务器列表
export function getServerList(): Promise<
  Response<
    Array<{
      id: string
      display_name: string
    }>
  >
> {
  return axios.get('/servers').then((res) => res.data)
}