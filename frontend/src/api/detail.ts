import axios, { type Response } from '@/utils/axios'

// Bird协议通道的路由统计信息
export interface RouteStats {
    accepted: string
    filtered: string
    ignored: string
    received: string
    rejected: string
}

// Bird协议通道信息
export interface Channel {
    name: string
    state: string
    route_stats: {
        import_updates: RouteStats
        import_withdraws: RouteStats
        export_updates: RouteStats
        export_withdraws: RouteStats
    }
    bgp_next_hop: string
}

// Bird协议详情信息
export interface BirdDetail {
    state: string
    neighbor_address: string
    neighbor_as: string
    local_as: string
    neighbor_id: string
    channels: Channel[]
}

// 节点协议详情响应
export interface NodeProtocolDetail {
    detail: BirdDetail
    displayName: string
    id: string
    rawOutput: string
    error?: string
}

// 获取指定节点和协议的详情信息
export async function getProtocolDetail(
    server: string,
    protocol: string
): Promise<
    Response<{
        [key: string]: NodeProtocolDetail
    }>
> {
    return axios
        .get(`/bird/detail?server=${server}&protocol=${protocol}`)
        .then((res) => res.data)
}

// 指定多个节点或协议
export async function queryProtocolDetail(params: {
    server?: string | string[]
    protocol?: string | string[]
}): Promise<
    Response<{
        [key: string]: NodeProtocolDetail
    }>
> {
    return axios.get('/bird/detail', { params }).then((res) => res.data)
}
