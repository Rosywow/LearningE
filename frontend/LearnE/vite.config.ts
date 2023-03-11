import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: 'localhost',
    port: 8080,
    open: true,
    https: false,
    proxy: {
      '/api': 'http://localhost:1235/',
      '/trans':{
        target:'http://api.fanyi.baidu.com/api/trans/vip/translate',
        ws:true,
        changeOrigin:true,
        rewrite:(path)=>path.replace(/^\/trans/,'')
      }
    },
  },
})
