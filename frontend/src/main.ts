import './assets/base.css'
import 'ant-design-vue/dist/reset.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Event
import('@/event/request')

// Global Function
import('@/utils/globalMessage')

app.mount('#app')
