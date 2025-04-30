<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getServerList } from '@/api/serverlist'
import { getNodeTraceroute, type NodeTracerouteResult } from '@/api/traceroute'

const target = ref('')
// 数组类型，支持多选服务器
const selectedServers = ref<string[]>([])
const servers = ref<Array<{ id: string; display_name: string }>>([])
const loading = ref(false)
const error = ref('')
const tracerouteResult = ref<Record<string, NodeTracerouteResult> | null>(null)
const showRawOutput = ref(false)
// 当前查看的服务器ID
const currentViewServer = ref('')

// 当前选中服务器的结果
const currentResult = computed(() => {
    if (!tracerouteResult.value || !currentViewServer.value) return null
    return tracerouteResult.value[currentViewServer.value]
})

// 获取服务器列表
async function fetchServers() {
    try {
        const response = await getServerList()
        if (response.status === 200 && response.data) {
            console.log(response.data)
            servers.value = response.data
            if (servers.value.length > 0) {
                // 默认不选择任何服务器
                selectedServers.value = []
            }
        } else {
            error.value = response.msg || 'Failed to get server list'
        }
    } catch (err) {
        error.value = err instanceof Error ? err.message : 'Network Request Error'
        console.error(err)
    }
}

// 执行路由追踪
async function runTraceroute() {
    if (!target.value || selectedServers.value.length === 0) {
        error.value = 'Please select at least one server and enter the destination address'
        return
    }

    loading.value = true
    error.value = ''
    tracerouteResult.value = null

    try {
        // 使用一个请求发送所有选定的服务器ID
        const response = await getNodeTraceroute(selectedServers.value, target.value)

        if (response.status === 200 && response.data) {
            tracerouteResult.value = response.data

            // 设置当前查看的服务器为第一个有结果的服务器
            const firstServerId = Object.keys(response.data)[0]
            if (firstServerId) {
                currentViewServer.value = firstServerId
            }

            // 检查是否有错误
            if (Object.keys(response.data).length === 0) {
                error.value = 'All route trace requests fail'
            }
        } else {
            error.value = response.msg || 'Route Trace Request Failed'
        }
    } catch (err) {
        error.value = err instanceof Error ? err.message : 'Network Request Error'
        console.error(err)
    } finally {
        loading.value = false
    }
}

// 切换原始输出显示
function toggleRawOutput() {
    showRawOutput.value = !showRawOutput.value
}

// 根据ID获取服务器显示名称
function getServerName(id: string): string {
    const server = servers.value.find((server) => server.id === id)
    return server ? server.display_name : id
}

onMounted(fetchServers)
</script>

<template>
    <div class="traceroute-container">
        <a-card class="main-card">
            <template #title>
                <div class="card-title">Traceroute</div>
            </template>

            <div class="form-container">
                <a-form :model="{ target, selectedServers }" layout="inline">
                    <a-form-item label="Servers">
                        <a-select
                            v-model:value="selectedServers"
                            placeholder="Select servers"
                            style="width: 300px"
                            :disabled="loading"
                            mode="multiple"
                            :maxTagCount="3"
                        >
                            <a-select-option
                                v-for="server in servers"
                                :key="server.id"
                                :value="server.id"
                            >
                                {{ server.display_name }}
                            </a-select-option>
                        </a-select>
                    </a-form-item>

                    <a-form-item label="Target Address">
                        <a-input
                            v-model:value="target"
                            placeholder="Enter IP address or domain name"
                            style="width: 250px"
                            :disabled="loading"
                        />
                    </a-form-item>

                    <a-form-item>
                        <a-button
                            type="primary"
                            @click="runTraceroute"
                            :loading="loading"
                        >
                            Start Trace
                        </a-button>
                    </a-form-item>
                </a-form>
            </div>

            <div v-if="!tracerouteResult && !loading && !error" class="welcome-container">
                <div class="welcome-content">
                    <h2>Ready to trace network routes</h2>
                    <p>
                        Select one or more servers and enter a target address to begin
                        tracing the network path.
                    </p>
                    <p class="tips">
                        The traceroute tool shows the path that network packets take from
                        the selected server to the target destination.
                    </p>
                </div>
            </div>

            <div v-if="error" class="error-container">
                <a-alert type="error" :message="error" show-icon />
            </div>

            <div v-if="loading" class="loading-container">
                <a-spin tip="Performing traceroute, please wait...">
                    <a-skeleton :rows="10" active />
                </a-spin>
            </div>

            <div v-if="tracerouteResult && !loading" class="result-container">
                <!-- 服务器结果选择器 -->
                <a-tabs v-model:activeKey="currentViewServer" class="server-tabs">
                    <a-tab-pane
                        v-for="(result, serverId) in tracerouteResult"
                        :key="serverId"
                        :tab="getServerName(serverId)"
                    >
                    </a-tab-pane>
                </a-tabs>

                <div v-if="currentResult" class="server-result">
                    <div class="result-header">
                        <h2>
                            Traceroute results from
                            {{
                                currentResult.displayName ||
                                getServerName(currentResult.id)
                            }}
                            to {{ target }}
                        </h2>
                        <a-button @click="toggleRawOutput" type="link">
                            {{
                                showRawOutput ? 'Show Parsed Results' : 'Show Raw Output'
                            }}
                        </a-button>
                    </div>

                    <!-- 显示服务器的描述信息 -->
                    <div v-if="currentResult.description">
                        <a-alert :message=currentResult.description type="info" show-icon/>
                    </div>

                    <!-- 显示错误信息 -->
                    <div v-if="currentResult.errorMsg">
                        <a-alert :message=currentResult.errorMsg type="error" show-icon/>
                    </div>

                    <div v-if="!showRawOutput" class="hop-list">
                        <a-table
                            :dataSource="currentResult.hops"
                            :pagination="false"
                            :bordered="true"
                            size="middle"
                            :scroll="{ x: '700px' }"
                            rowKey="hop_num"
                        >
                            <a-table-column title="Hop" dataIndex="hop_num" />
                            <a-table-column title="IP Address" dataIndex="address" />
                            <a-table-column title="Hostname" dataIndex="host" />
                            <a-table-column title="RTT" dataIndex="rtt" />
                        </a-table>
                    </div>

                    <div v-else class="raw-output">
                        <a-card title="Raw Output">
                            <pre>{{ currentResult.rawOutput }}</pre>
                        </a-card>
                    </div>
                </div>

                <div v-else class="no-result">
                    <a-empty description="No result available for selected server" />
                </div>
            </div>
        </a-card>
    </div>
</template>

<style scoped>
.traceroute-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

.main-card {
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    overflow: hidden;
}

.card-title {
    font-size: 18px;
    font-weight: 600;
}

.form-container {
    margin-bottom: 24px;
    padding: 18px;
    background-color: #f5f7fa;
    border-radius: 8px;
}

.welcome-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 300px;
    background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
    border-radius: 8px;
    text-align: center;
    padding: 40px;
}

.welcome-content {
    max-width: 600px;
}

.welcome-content h2 {
    font-size: 24px;
    margin-bottom: 16px;
    color: #222;
}

.welcome-content p {
    color: #666;
    font-size: 16px;
    line-height: 1.6;
    margin-bottom: 8px;
}

.welcome-content .tips {
    font-size: 14px;
    color: #888;
    margin-top: 16px;
    background-color: rgba(24, 144, 255, 0.1);
    padding: 12px;
    border-radius: 4px;
    border-left: 4px solid #1890ff;
}

.error-container {
    margin-bottom: 24px;
}

.loading-container {
    padding: 40px 0;
    text-align: center;
}

.result-container {
    margin-top: 24px;
}

.result-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
}

.hop-list {
    margin-top: 16px;
    box-sizing: border-box;
}

.raw-output {
    margin-top: 16px;
}

.raw-output pre {
    white-space: pre-wrap;
    word-wrap: break-word;
    background-color: #f5f5f5;
    padding: 12px;
    border-radius: 4px;
    font-family: monospace;
    font-size: 14px;
    max-height: 600px;
    overflow-y: auto;
}

.server-tabs {
    margin-bottom: 20px;
}

.server-result {
    background-color: #f9f9f9;
    padding: 16px;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.no-result {
    padding: 40px 0;
    text-align: center;
}

@media (max-width: 768px) {
    .form-container :deep(.ant-form) {
        display: flex;
        flex-direction: column;
    }

    .form-container :deep(.ant-form-item) {
        margin-bottom: 12px;
        width: 100%;
    }

    .form-container :deep(.ant-select),
    .form-container :deep(.ant-input) {
        width: 100% !important;
    }

    .welcome-container {
        padding: 20px;
    }
}
</style>
