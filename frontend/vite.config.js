import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import svgr from 'vite-plugin-svgr'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), svgr()],
  server: {
    host: "https://dev.sumcrowds.com",
    port: 3000,
    allowedHosts: ["testing.sumcrowds.com", "sumcrowds.com", "dev.sumcrowds.com"],
  },
  optimizeDeps: {
    include: ['i17next', 'react-i18next', 'i18next-http-backend']
  }
})
