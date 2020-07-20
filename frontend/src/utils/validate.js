export function isValidUsername(str) {
  if (!str) return false
  return str.length <= 100
}

export function isExternal(path) {
  return /^(https?:|mailto:|tel:)/.test(path)
}

export function isUUID(str) {
  return /^[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(str)
}

export function isNilUUID(str) {
  return /^00000000-0000-0000-0000-000000000000$/.test(str)
}
