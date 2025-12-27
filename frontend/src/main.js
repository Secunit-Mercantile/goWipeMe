import { mount } from 'svelte'
import App from './App.svelte'
import './styles/themes.css'

// Set initial theme
const storedTheme = localStorage.getItem('theme') || 'refined'
document.body.setAttribute('data-theme', storedTheme)

const app = mount(App, {
  target: document.getElementById('app'),
})

export default app
