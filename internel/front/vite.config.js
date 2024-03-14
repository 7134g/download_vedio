import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'


// https://vitejs.dev/config/
export default defineConfig({

  server: {
    port: 9999,
  },

  build: {
    outDir: 'dist',
    assetsDir: '',
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        entryFileNames: '[name].js', // 指定输出的js文件名称
        assetFileNames(chunkInfo) {
          return chunkInfo.name === 'index' ? '[name].[ext]' : `[name].[ext]`;
        }, // 指定输出的asset（如图片、字体等）文件名称
      },
    },
  },

  plugins: [
    vue(),
    vueJsx(),
    AutoImport({ /* options */ }),
    Components({
      imports: ['vue', 'vue-router', 'pinia'],
      resolvers: [
        ElementPlusResolver(),
      ]
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
