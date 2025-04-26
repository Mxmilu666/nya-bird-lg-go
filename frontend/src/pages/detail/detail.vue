<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getProtocolDetail } from '@/api/detail'
import type { NodeProtocolDetail } from '@/api/detail'
import DetailHeader from './DetailHeader.vue'
import BasicInformation from './BasicInformation.vue'
import ChannelInformation from './ChannelInformation.vue'
import RawDataView from './RawDataView.vue'

const route = useRoute()
const serverName = ref(route.params.server as string)
const protocolName = ref(route.params.protocol as string)

const loading = ref(true)
const error = ref('')
const detailData = ref<NodeProtocolDetail | null>(null)
const currentView = ref('raw') // Default to show raw data view

onMounted(async () => {
    try {
        const response = await getProtocolDetail(serverName.value, protocolName.value)
        if (response.status === 200 && response.data) {
            const key = Object.keys(response.data)[0]
            detailData.value = response.data[key]
        } else {
            error.value = response.msg || 'Detail data not found'
        }
    } catch (err) {
        error.value =
            err instanceof Error ? err.message : 'Error occurred while fetching data'
        console.error(err)
    } finally {
        loading.value = false
    }
})
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
        </div>
        <div v-else class="detail-content">
            <a-card class="detail-card">
                <template #title>
                    <DetailHeader 
                        :protocol-name="protocolName"
                        :server-name="serverName"
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
}

.detail-content {
    max-width: 1000px;
    margin: 0 auto;
}

.detail-card {
    margin-bottom: 20px;
}
</style>
