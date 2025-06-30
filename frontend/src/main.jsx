import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.jsx'
import { initI18n } from './utils/i18n.js'

initI18n().then(() => {
  createRoot(document.getElementById('root')).render(
    <App />
  )
})
