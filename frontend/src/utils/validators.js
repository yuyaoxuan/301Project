
export const required = (value) => {
  return !!value || 'This field is required'
}

export const email = (value) => {
  const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return pattern.test(value) || 'Invalid email address'
}

export const phone = (value) => {
  const pattern = /^\+?[\d\s-]{10,}$/
  return pattern.test(value) || 'Invalid phone number'
}

export const minLength = (length) => (value) => {
  return !value || value.length >= length || `Must be at least ${length} characters`
}
