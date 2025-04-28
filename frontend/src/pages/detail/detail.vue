<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getProtocolDetail } from '@/api/detail'
import type { NodeProtocolDetail } from '@/api/detail'
import DetailHeader from './DetailHeader.vue'
import BasicInformation from './BasicInformation.vue'
import ChannelInformation from './ChannelInformation.vue'
import RawDataView from './RawDataView.vue'

const route = useRoute()
const router = useRouter()
const serverName = ref(route.params.server as string)
const protocolName = ref(route.params.protocol as string)

const loading = ref(true)
const error = ref('')
const detailData = ref<NodeProtocolDetail | null>(null)
const currentView = ref('raw') // Default to show raw data view

async function fetchProtocolDetail() {
    loading.value = true
    error.value = ''

    try {
        const response = await getProtocolDetail(serverName.value, protocolName.value)
        if (response.status === 200) {
            if (response.data) {
                const keys = Object.keys(response.data)
                if (keys.length === 0) {
                    error.value = 'No protocol detail data found'
                    return
                }
                const key = keys[0]
                detailData.value = response.data[key]

                if (response.data[key]?.error) {
                    error.value = response.data[key].error
                }
            } else {
                error.value = 'Empty response data'
            }
        } else {
            error.value = response.msg || 'Failed to fetch data'
        }
    } catch (err) {
        error.value = err instanceof Error ? err.message : 'Network request error'
        console.error(err)
    } finally {
        loading.value = false
    }
}

function goBack() {
    router.back()
}

onMounted(fetchProtocolDetail)
</script>

<template>
    <div class="detail-container">
        <div v-if="loading" class="loading">
            <a-skeleton :rows="10" active />
        </div>
        <div v-else-if="error" class="error">
            <a-alert
                message="Failed to retrieve data"
                type="error"
                :description="error"
                show-icon
            />
            <div class="button-group">
                <a-button type="primary" @click="fetchProtocolDetail" class="action-button">
                    Reload
                </a-button>
                <a-button @click="goBack" class="action-button">
                    Go Back
                </a-button>
            </div>
        </div>
        <div v-else class="detail-content">
            <a-card class="detail-card">
                <template #title>
                    <DetailHeader
                        :protocol-name="protocolName"
                        :server-name="detailData?.displayName || ''"
                        :current-view="currentView"
                        @update:current-view="currentView = $event"
                    />
                </template>

                <div v-if="detailData">
                    <div v-if="currentView === 'analysis'">
                        <BasicInformation :detail-data="detailData" />
                        <ChannelInformation
                            v-if="detailData.detail.channels"
                            :channels="detailData.detail.channels"
                        />
                    </div>
                    <div v-else>
                        <RawDataView :raw-output="detailData.rawOutput" />
                    </div>
                </div>
            </a-card>
        </div>
    </div>
</template>

<style scoped>
.detail-container {
    padding: 20px;
}

.loading,
.error {
    margin: 50px auto;
    max-width: 800px;
    text-align: center;
}

.button-group {
    margin-top: 16px;
    display: flex;
    justify-content: center;
    gap: 12px;
}

.action-button {
    min-width: 90px;
}

.detail-content {
    max-width: 1000px;
    margin: 0 auto;
}

.detail-card {
    margin-bottom: 20px;
}
</style>