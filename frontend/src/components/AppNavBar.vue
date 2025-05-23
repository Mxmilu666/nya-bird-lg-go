<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MenuOutlined } from '@ant-design/icons-vue'
import type { MenuProps } from 'ant-design-vue'

// 导入 ItemType 类型
type MenuItem = NonNullable<MenuProps['items']>[number]

interface MenuConfig {
    key: string
    label: string
    path: string
    icon?: string
}

const router = useRouter()
const route = useRoute()
const selectedKeys = ref<string[]>([])
const isMobile = ref(false)

// 菜单配置
const menuConfig: MenuConfig[] = [
    {
        key: '1',
        label: 'Summary',
        path: '/summary'
    },
    {
        key: '2',
        label: 'Traceroute',
        path: '/traceroute'
    }
]

// 修改计算属性，确保返回的是符合 Ant Design 要求的 MenuItem 类型
const menuItems = computed<MenuItem[]>(() =>
    menuConfig.map((item) => ({
        key: item.key,
        label: item.label,
        onClick: () => router.push(item.path)
    }))
)

// 使用泛型定义 throttle 函数
const throttle = <T extends (...args: unknown[]) => unknown>(fn: T, delay = 250) => {
    let lastCall = 0
    return function (this: unknown, ...args: Parameters<T>): ReturnType<T> | undefined {
        const now = Date.now()
        if (now - lastCall < delay) return undefined
        lastCall = now
        return fn.apply(this, args) as ReturnType<T>
    }
}

const checkMobile = throttle(() => {
    isMobile.value = window.innerWidth < 768
})

const updateSelectedMenu = () => {
    const path = route.path
    const matched = menuConfig.find((item) => path.includes(item.path))
    selectedKeys.value = matched ? [matched.key] : []
}

onMounted(() => {
    updateSelectedMenu()
    checkMobile()
    window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
    window.removeEventListener('resize', checkMobile)
})

watch(() => route.path, updateSelectedMenu)
</script>

<template>
    <a-layout-header class="header">
        <div class="logo">
            <router-link to="/">
                <h1>Nya Bird LG</h1>
            </router-link>
        </div>

        <!-- 桌面端显示水平菜单 -->
        <template v-if="!isMobile">
            <a-menu
                v-model:selectedKeys="selectedKeys"
                theme="light"
                mode="horizontal"
                :items="menuItems"
                class="desktop-menu"
            />
        </template>

        <!-- 移动端显示下拉菜单 -->
        <template v-else>
            <a-dropdown placement="bottomRight">
                <a-button class="menu-trigger" type="text">
                    <menu-outlined />
                </a-button>
                <template #overlay>
                    <a-menu
                        v-model:selectedKeys="selectedKeys"
                        theme="light"
                        :items="menuItems"
                    />
                </template>
            </a-dropdown>
        </template>
    </a-layout-header>
</template>

<style scoped>
.header {
    display: flex;
    align-items: center;
    padding: 0 24px;
    background: #ffffff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    height: 64px;
    justify-content: flex-start;
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
    transition: color 0.3s ease;
}

.logo h1:hover {
    color: #40a9ff;
}

.desktop-menu {
    flex: 1;
}

.menu-trigger {
    font-size: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    height: 32px;
    width: 32px;
    padding: 0;
    transition: background-color 0.3s;
}

.menu-trigger:hover {
    background-color: rgba(0, 0, 0, 0.04);
}

/* 移动端样式调整 */
@media (max-width: 767px) {
    .header {
        padding: 0 12px;
        justify-content: space-between;
    }

    .logo h1 {
        font-size: 16px;
    }
}
</style>
