import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react-swc';
import { fileURLToPath, URL } from 'node:url';

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
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  test: {
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
  },
});
