import axios from 'axios'

const sendEvCrawlab = (eventCategory, eventAction, eventLabel) => {
  axios.get(process.env.VUE_APP_CRAWLAB_BASE_URL + '/track', {
    params: {
      uid: localStorage.getItem('uid'),
      sid: sessionStorage.getItem('sid'),
      ec: eventCategory,
      ea: eventAction,
      el: eventLabel,
      v: 'lite-' + sessionStorage.getItem('v')
    }
  }).catch(() => {})
}

export default {
  sendPv(page) {
    if (localStorage.getItem('useStats') !== '0') {
      window._hmt.push(['_trackPageview', page])
      sendEvCrawlab('访问页面', page, '')
    }
  },
  sendEv(category, eventName, optLabel, optValue) {
    if (localStorage.getItem('useStats') !== '0') {
      window._hmt.push(['_trackEvent', category, eventName, optLabel, optValue])
      sendEvCrawlab(category, eventName, optLabel)
    }
  }
}
