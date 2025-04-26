<script setup lang="ts">
import type { TableColumnsType } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { h } from 'vue'

interface Protocol {
    name: string
    proto: string
    table: string
    state: string
    since: string
    info: string
}

const props = defineProps<{
    protocols: Protocol[]
    serverId: string
}>()

const router = useRouter()

const columns: TableColumnsType = [
    {
        title: 'Name',
        dataIndex: 'name',
        width: 200,
        customRender: ({ text }) => {
            return h(
                'a',
                {
                    onClick: () => {
                        router.push(`/detail/${props.serverId}/${text}`)
                    },
                    style: {
                        cursor: 'pointer',
                        color: '#1890ff',
                        textDecoration: 'none'
                    }
                },
                text
            )
        }
    },
    {
        title: 'Proto',
        dataIndex: 'proto',
        width: 120
    },
    {
        title: 'Table',
        dataIndex: 'table',
        width: 120
    },
    {
        title: 'State',
        dataIndex: 'state',
        width: 120,
        customRender: ({ text }) => {
            return text
        }
    },
    {
        title: 'Since',
        dataIndex: 'since',
        width: 180
    },
    {
        title: 'Info',
        dataIndex: 'info'
    }
]

// 添加表格行的自定义样式
const customRowStyle = (record: Protocol) => {
    const state = record.state.toLowerCase()
    
    const styleMap: Record<string, Record<string, string>> = {
        'up': { background: '#eafbf4' },      // 浅绿色背景
        'down': { background: '#FAFAFA' },    // 浅灰色背景
        'start': { background: '#fff1f0' },   // 浅红色背景
        'passive': { background: '#e6f7ff' }  // 浅蓝色背景
    }
    
    return { style: styleMap[state] || {} }
}
</script>

<template>
    <div class="table-responsive">
        <a-table
            :columns="columns"
            :data-source="protocols"
            :pagination="false"
            :bordered="true"
            :row-key="(record: Protocol) => record.name"
            size="middle"
            :customRow="customRowStyle"
            :scroll="{ x: '800px' }"
        />
    </div>
</template>

<style scoped>
:deep(table tbody tr:hover > td) {
    background-color: rgba(0, 0, 0, 0.08) !important;
}

.table-responsive {
    width: 100%;
    overflow-x: auto;
}

@media screen and (max-width: 768px) {
    :deep(.ant-table-thead > tr > th),
    :deep(.ant-table-tbody > tr > td) {
        padding: 8px 6px;
        font-size: 13px;
    }

    :deep(.ant-table-thead > tr > th) {
        white-space: nowrap;
    }
}
</style>
