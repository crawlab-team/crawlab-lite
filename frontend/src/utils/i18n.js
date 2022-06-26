import { getCurrentInstance } from 'vue'

// translate router.meta.title, be used in breadcrumb sidebar tagsview
export function generateTitle(title) {
  const currIns = getCurrentInstance()
  const hasKey = currIns.t('route.' + title)

  if (hasKey) {
    // $t :this method from vue-i18n, inject in @/lang/index.js
    const translatedTitle = currIns.t('route.' + title)

    return translatedTitle
  }
  return title
}
