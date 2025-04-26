<template>
    <a-page-header
        :title="`Protocol Details: ${protocolName}`"
        :subTitle="`Server: ${serverName}`"
        @back="() => $router.go(-1)"
    >
        <template #extra>
            <a-button-group>
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

<script setup lang="ts">
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
</script>