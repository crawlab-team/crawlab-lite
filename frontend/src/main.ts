//@ts-nocheck

import 'normalize.css/normalize.css' // A modern alternative to CSS resets
import 'element-plus/theme-chalk/index.css'
import '@/styles/index.scss' // global css
import 'font-awesome/scss/font-awesome.scss' // FontAwesome
import 'codemirror/lib/codemirror.js'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/darcula.css'
import '@/icons' // icon
import '@/permission' // permission control

import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import * as Vue from 'vue'

import {
  FontAwesomeIcon,
  FontAwesomeLayers,
  FontAwesomeLayersText,
} from '@fortawesome/vue-fontawesome'

import App from './App'
import ElementUI from 'element-plus'
import { codemirror } from 'vue-codemirror-lite'
import { createApp } from 'vue'
import { fab } from '@fortawesome/free-brands-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { getI18n } from 'crawlab-ui'
import { library } from '@fortawesome/fontawesome-svg-core'
import locale from 'element-plus/lib/locale/lang/en' // lang i18n
import request from './api/request'
import router from './router'
import store from './store'
// import i18n from './i18n'
import utils from './utils'

const app = (window.$vueApp = createApp(App))
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(getI18n())

// code mirror
app.use(codemirror)

// element-ui
app.use(ElementUI, { locale })

// font-awesome
library.add(fab)
library.add(far)
library.add(fas)
app.component('FontAwesomeIcon', FontAwesomeIcon)
app.component('FontAwesomeLayers', FontAwesomeLayers)
app.component('FontAwesomeLayersText', FontAwesomeLayersText)

// 百度统计
if (localStorage.getItem('useStats') !== '0') {
  window._hmt = window._hmt || []
  ;(function () {
    const hm = document.createElement('script')
    hm.src = 'https://hm.baidu.com/hm.js?c35e3a563a06caee2524902c81975add'
    const s = document.getElementsByTagName('script')[0]
    s.parentNode.insertBefore(hm, s)
  })()
}

// inject request api
app.config.globalProperties.$request = request

// inject utils
app.config.globalProperties.$utils = utils

// inject stats
app.config.globalProperties.$st = utils.stats

app.config.globalProperties.routerAppend = (path, pathToAppend) => {
  return path + (path.endsWith('/') ? '' : '/') + pathToAppend
}
app.use(store)
app.use(router)
app.mount('#app')

export default app
