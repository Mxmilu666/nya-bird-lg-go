<script setup lang="ts">
defineOptions({
    name: 'HomeComponent'
})

import { message } from 'ant-design-vue'
import { getSummary } from '@/api/summary'

const [messageApi, contextHolder] = message.useMessage()
const info = () => {
    messageApi.info('Hello, Ant Design Vue!')
}

getSummary()
    .then((res) => {
        console.log(res)
    })
    .catch((err) => {
        console.error(err)
    })
</script>

<template>
    <context-holder />
    <a-button type="primary" @click="info">Display normal message</a-button>
</template>
