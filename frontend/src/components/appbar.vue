<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const selectedKeys = ref<string[]>(['1'])

const menuItems = [
    {
        key: '1',
        label: 'Summary',
        onClick: () => router.push('/summary')
    }
]
</script>

<template>
    <a-layout-header class="header">
        <div class="logo">
            <router-link to="/">
                <h1>Nya Bird Looking Glass</h1>
            </router-link>
        </div>
        <a-menu
            v-model:selectedKeys="selectedKeys"
            theme="light"
            mode="horizontal"
            :items="menuItems"
        />
    </a-layout-header>
</template>

<style scoped>
.header {
    display: flex;
    align-items: center;
    padding: 0 24px;
    background: #ffffff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.logo {
    margin-right: 24px;
}

.logo a {
    text-decoration: none;
}

.logo h1 {
    color: #1890ff;
    margin: 0;
    font-size: 18px;
    font-weight: 600;
}
</style>
