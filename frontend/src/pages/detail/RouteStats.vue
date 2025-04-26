<template>
  <a-table :columns="columns" :data-source="dataSource" :pagination="false" size="small" bordered>
  </a-table>
</template>

<script setup lang="ts">
import { defineProps, computed } from 'vue'
import type { RouteStats } from '@/api/detail'

const props = defineProps<{
  importUpdates: RouteStats;
  importWithdraws: RouteStats;
  exportUpdates: RouteStats;
  exportWithdraws: RouteStats;
}>()

const columns = [
  {
    title: 'Stats',
    dataIndex: 'metric',
    key: 'metric',
  },
  {
    title: 'Import Updates',
    dataIndex: 'importUpdates',
    key: 'importUpdates',
  },
  {
    title: 'Import Withdraws',
    dataIndex: 'importWithdraws',
    key: 'importWithdraws',
  },
  {
    title: 'Export Updates',
    dataIndex: 'exportUpdates',
    key: 'exportUpdates',
  },
  {
    title: 'Export Withdraws',
    dataIndex: 'exportWithdraws',
    key: 'exportWithdraws',
  }
]

const dataSource = computed(() => [
  {
    key: '1',
    metric: 'accepted',
    importUpdates: props.importUpdates.accepted,
    importWithdraws: props.importWithdraws.accepted,
    exportUpdates: props.exportUpdates.accepted,
    exportWithdraws: props.exportWithdraws.accepted
  },
  {
    key: '2',
    metric: 'filtered',
    importUpdates: props.importUpdates.filtered,
    importWithdraws: props.importWithdraws.filtered,
    exportUpdates: props.exportUpdates.filtered,
    exportWithdraws: props.exportWithdraws.filtered
  },
  {
    key: '3',
    metric: 'ignored',
    importUpdates: props.importUpdates.ignored,
    importWithdraws: props.importWithdraws.ignored,
    exportUpdates: props.exportUpdates.ignored,
    exportWithdraws: props.exportWithdraws.ignored
  },
  {
    key: '4',
    metric: 'received',
    importUpdates: props.importUpdates.received,
    importWithdraws: props.importWithdraws.received,
    exportUpdates: props.exportUpdates.received,
    exportWithdraws: props.exportWithdraws.received
  },
  {
    key: '5',
    metric: 'rejected',
    importUpdates: props.importUpdates.rejected,
    importWithdraws: props.importWithdraws.rejected,
    exportUpdates: props.exportUpdates.rejected,
    exportWithdraws: props.exportWithdraws.rejected
  },
])
</script>
