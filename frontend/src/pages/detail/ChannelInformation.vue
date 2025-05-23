<script setup lang="ts">
import RouteStatsComponent from './RouteStats.vue'
import type { RouteStats } from '@/api/detail'

interface Channel {
    name: string
    state: string
    bgp_next_hop: string
    route_stats: {
        import_updates: RouteStats
        import_withdraws: RouteStats
        export_updates: RouteStats
        export_withdraws: RouteStats
    }
}

defineProps<{
    channels: Channel[]
}>()
</script>

<template>
    <div class="detail-section" v-if="channels && channels.length > 0">
        <h3>Channel Information</h3>
        <div v-for="(channel, index) in channels" :key="index" class="channel-item">
            <a-divider>{{ channel.name }}</a-divider>

            <h4>Route Statistics</h4>
            <template
                v-if="
                    channel.route_stats &&
                    channel.route_stats.import_updates &&
                    channel.route_stats.import_withdraws &&
                    channel.route_stats.export_updates &&
                    channel.route_stats.export_withdraws
                "
            >
                <div class="stats-container">
                    <div class="stats-item full-width">
                        <RouteStatsComponent
                            :importUpdates="channel.route_stats.import_updates"
                            :importWithdraws="channel.route_stats.import_withdraws"
                            :exportUpdates="channel.route_stats.export_updates"
                            :exportWithdraws="channel.route_stats.export_withdraws"
                        />
                    </div>
                </div>
            </template>

            <a-descriptions bordered :column="2">
                <a-descriptions-item label="Stats">{{
                    channel.state
                }}</a-descriptions-item>
                <a-descriptions-item label="BGP Next Hop">{{
                    channel.bgp_next_hop
                }}</a-descriptions-item>
            </a-descriptions>
        </div>
    </div>
</template>

<style scoped>
.detail-section {
    margin-bottom: 30px;
}

.channel-item {
    margin-bottom: 20px;
}

.stats-container {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
}

.stats-item {
    flex: 1 1 calc(50% - 20px);
    margin-bottom: 20px;
    max-width: 100%;
    overflow-x: hidden;
}

.stats-item.full-width {
    flex: 1 1 100%;
    max-width: 100%;
    overflow-x: hidden;
}
</style>
