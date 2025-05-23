<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSummary } from '@/api/summary'
import SummaryTable from './SummaryTable.vue'

interface Protocol {
    name: string
    proto: string
    table: string
    state: string
    since: string
    info: string
}

interface NodeData {
    id: string
    displayName: string
    protocols: Protocol[]
    error?: string
}

const loading = ref(true)
const groupedData = ref<NodeData[]>([])
const errorMsg = ref('')

const fetchData = async () => {
    loading.value = true
    errorMsg.value = ''

    try {
        const response = await getSummary()
        const actualData = response.data

        if (!actualData || typeof actualData !== 'object') {
            errorMsg.value = 'Invalid data format'
            return
        }

        // 转换数据
        const nodes = Object.entries(actualData).map(([nodeId, nodeInfo]) => {
            const info = nodeInfo as {
                displayName?: string
                protocols?: Protocol[]
                error?: string
            }
            return {
                id: nodeId,
                displayName: info?.displayName || nodeId,
                protocols: Array.isArray(info?.protocols) ? info.protocols : [],
                error: info?.error // 保留error信息
            }
        })

        groupedData.value = nodes
    } catch (error) {
        errorMsg.value = `Failed to fetch data: ${error instanceof Error ? error.message : String(error)}`
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    fetchData()
})
</script>

<template>
    <div class="summary-container">
        <a-alert
            v-if="errorMsg"
            type="error"
            :message="errorMsg"
            show-icon
            style="margin-bottom: 20px"
        />
        <a-spin :spinning="loading">
            <div v-if="groupedData.length === 0 && !loading" class="empty-state">
                <a-empty description="No Data" />
                <a-button type="primary" @click="fetchData" style="margin-top: 16px">
                    Reload
                </a-button>
            </div>
            <div v-else>
                <div v-for="node in groupedData" :key="node.id" class="node-section">
                    <h1 class="node-title">{{ node.displayName }}: show protocols</h1>
                    <a-alert
                        v-if="node.error"
                        type="error"
                        :message="node.error"
                        show-icon
                        style="margin-bottom: 12px"
                    />
                    <summary-table
                        v-if="!node.error"
                        :protocols="node.protocols"
                        :server-id="node.id"
                    />
                </div>
            </div>
        </a-spin>
    </div>
</template>

<style scoped>
.summary-container {
    padding: 30px 70px;
    width: 100%;
    box-sizing: border-box;
}

.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 0;
}

.node-section {
    margin-bottom: 30px;
}

.node-title {
    margin-bottom: 12px;
    font-size: 25px;
    font-weight: 500;
}

@media screen and (max-width: 768px) {
    .summary-container {
        padding: 25px 20px;
    }

    .node-title {
        font-size: 16px;
        margin-bottom: 8px;
    }

    .node-section {
        margin-bottom: 20px;
    }
}
</style>
