import axios, { type Response } from '@/utils/axios'

// 跳跃点信息接口
export interface Hop {
    hop_num: number
    address: string
    rtt: string
    host: string
}

// 节点路由追踪结果接口
export interface NodeTracerouteResult {
    displayName: string
    id: string
    hops: Hop[]
    rawOutput: string
    error?: string
    description?: string
    errorMsg?: string
}

// 获取特定节点到指定目标的路由追踪结果
export function getNodeTraceroute(
    server: string | string[],
    target: string
): Promise<Response<Record<string, NodeTracerouteResult>>> {
    // 将服务器数组转换为逗号分隔的字符串
    const serverParam = Array.isArray(server) ? server.join(',') : server
    
    return axios
        .get(`/traceroute?target=${target}&server=${serverParam}`)
        .then((res) => {
            if (res.data.data && !res.data.data[typeof server === 'string' ? server : server[0]]) {
                return {
                    ...res.data,
                    data: { [typeof server === 'string' ? server : server[0]]: res.data.data }
                };
            }
            return res.data;
        });
}