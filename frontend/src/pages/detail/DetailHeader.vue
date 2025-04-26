<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

defineProps<{
    protocolName: string;
    serverName: string;
    currentView: string;
}>();

const emit = defineEmits<{
    (e: 'update:currentView', view: string): void;
}>();

const updateView = (view: string) => {
    emit('update:currentView', view);
};

// 添加响应式判断
const isMobile = ref(false);

const checkMobile = () => {
    isMobile.value = window.innerWidth < 768;
};

onMounted(() => {
    checkMobile();
    window.addEventListener('resize', checkMobile);
});

onUnmounted(() => {
    window.removeEventListener('resize', checkMobile);
});
</script>

<template>
    <a-page-header
        :title="`Protocol Details: ${protocolName}`"
        :subTitle="`Server: ${serverName}`"
        @back="() => $router.go(-1)"
        class="detail-header"
    >
        <template #extra>
            <a-button-group :class="{ 'mobile-button-group': isMobile }">
                <a-button
                    type="primary"
                    :ghost="currentView !== 'raw'"
                    @click="updateView('raw')"
                >
                    Raw Data
                </a-button>
                <a-button
                    type="primary"
                    :ghost="currentView !== 'analysis'"
                    @click="updateView('analysis')"
                >
                    Analysis View
                </a-button>
            </a-button-group>
        </template>
    </a-page-header>
</template>

<style scoped>
.detail-header :deep(.ant-page-header-heading) {
    flex-wrap: wrap;
}

.detail-header :deep(.ant-page-header-heading-title),
.detail-header :deep(.ant-page-header-heading-sub-title) {
    word-break: break-word;
    white-space: normal;
    overflow-wrap: break-word;
    max-width: 100%;
}

@media (max-width: 767px) {
    .detail-header :deep(.ant-page-header-heading-title) {
        font-size: 16px;
        margin-right: 0;
    }
    
    .detail-header :deep(.ant-page-header-heading-sub-title) {
        font-size: 14px;
        display: block;
        margin: 4px 0;
    }
    
    .detail-header :deep(.ant-page-header-heading-extra) {
        margin: 8px 0 0;
        width: 100%;
    }
    
    .mobile-button-group {
        display: flex;
        width: 100%;
    }
    
    .mobile-button-group button {
        flex: 1;
    }
}
</style>