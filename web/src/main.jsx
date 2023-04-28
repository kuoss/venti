import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)

if (localStorage.theme === 'dark') {
  localStorage.theme = 'light'
  document.documentElement.classList.remove('dark')
} else {
  localStorage.theme = 'dark'
  document.documentElement.classList.add('dark')
}