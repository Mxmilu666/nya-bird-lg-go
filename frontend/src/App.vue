<script setup lang="ts">
import Appbar from '@/components/appbar.vue'

defineOptions({
    name: 'AppComponent'
})
</script>

<template>
    <a-layout>
        <Appbar />
        <a-layout-content>
            <router-view />
        </a-layout-content>
    </a-layout>
</template>

<style>
.ant-layout {
    min-height: 100vh;
}

.ant-layout-content {
    padding: 24px;
}
</style>
