import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        vueDevTools(),
        Components({
            dts: true,
            resolvers: [
                AntDesignVueResolver({
                    importStyle: false // css in js
                })
            ]
        })
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    }
})
