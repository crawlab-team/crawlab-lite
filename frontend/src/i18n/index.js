import en from './en'
import zh from './zh'
import { createI18n } from 'vue-i18n'

const i18n = createI18n({
  locale: localStorage.getItem('lang') || 'zh',
  messages: {
    en,
    zh,
  },
  silentTranslationWarn: true,
})

export default i18n
