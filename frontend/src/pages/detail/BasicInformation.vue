<script setup lang="ts">
import type { NodeProtocolDetail } from '@/api/detail'
import { ref, onMounted, onUnmounted } from 'vue'

defineProps<{
    detailData: NodeProtocolDetail
}>()

// 添加响应式判断
const isMobile = ref(false)

const checkMobile = () => {
    isMobile.value = window.innerWidth < 768
}

onMounted(() => {
    checkMobile()
    window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
    window.removeEventListener('resize', checkMobile)
})
</script>

<template>
    <div class="detail-section">
        <h3>Basic Information</h3>
        <a-descriptions bordered :column="isMobile ? 1 : 2">
            <a-descriptions-item label="Status">{{
                detailData.detail.state
            }}</a-descriptions-item>
            <a-descriptions-item label="Neighbor Address">{{
                detailData.detail.neighbor_address
            }}</a-descriptions-item>
            <a-descriptions-item label="Neighbor AS">{{
                detailData.detail.neighbor_as
            }}</a-descriptions-item>
            <a-descriptions-item label="Local AS">{{
                detailData.detail.local_as
            }}</a-descriptions-item>
            <a-descriptions-item label="Neighbor ID">{{
                detailData.detail.neighbor_id
            }}</a-descriptions-item>
        </a-descriptions>
    </div>
</template>

<style scoped>
.detail-section {
    margin-bottom: 30px;
}
</style>
