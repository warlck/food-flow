import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import path from 'path';

export default defineConfig({
  server: {
    host: '::',
    port: 8081,
    proxy: {
      '/v1/auth': { target: 'http://localhost:6000', changeOrigin: true },
      '/v1': { target: 'http://localhost:3000', changeOrigin: true },
    },
  },
  plugins: [react()],
  resolve: {
    alias: { '@': path.resolve(__dirname, './src') },
  },
});
