<template>
    <div class="detail-section">
        <h3>Basic Information</h3>
        <a-descriptions bordered :column="2">
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

<script setup lang="ts">
import type { NodeProtocolDetail } from '@/api/detail';

defineProps<{
    detailData: NodeProtocolDetail;
}>();
</script>

<style scoped>
.detail-section {
    margin-bottom: 30px;
}
</style>