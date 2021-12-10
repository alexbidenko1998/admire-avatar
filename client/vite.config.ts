import {defineConfig} from 'vite';
import vue from '@vitejs/plugin-vue';
import path = require('path')

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': 'http://192.168.43.182:7015',
    },
  },
  alias: {
    '~': path.resolve(__dirname, './src'),
  },
});
