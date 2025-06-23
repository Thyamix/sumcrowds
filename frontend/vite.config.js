import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: "192.168.1.35",
    port: 3000,
    allowedHosts: ["testing.sumcrowds.com", "sumcrowds.com", "dev.sumcrowds.com"],
  }
})
