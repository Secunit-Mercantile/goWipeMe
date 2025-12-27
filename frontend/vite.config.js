import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    strictPort: true, // Fail if port is already in use
    host: true, // Expose to network (helps with detection)
  },
  logLevel: 'info', // Ensure Vite outputs info level logs
})
