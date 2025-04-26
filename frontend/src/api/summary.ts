import axios, { type Response } from '@/utils/axios'

// 获取所有节点的Bird协议摘要信息
export async function getSummary(): Promise<
    Response<{
        [key: string]: {
            displayName: string
            id: string
            protocols: Array<{
                name: string
                proto: string
                table: string
                state: string
                since: string
                info: string
            }>
        }
    }>
> {
    return axios.get('/bird/summary').then((res) => res.data)
}

// 获取单个节点的Bird协议摘要信息
export async function getNodeSummary(node: string): Promise<
    Response<{
        displayName: string
        id: string
        protocols: Array<{
            name: string
            proto: string
            table: string
            state: string
            since: string
            info: string
        }>
    }>
> {
    return axios.get(`/bird/summary/${node}`).then((res) => res.data)
}