
import api from './api'

export const authService = {
  async login(credentials) {
    const response = await api.post('/api/users/login', {
      username: credentials.username,
      password: credentials.password
    })
    if (!response.data?.id_token) {
      throw new Error('Invalid response from server')
    }
    api.defaults.headers.common['Authorization'] = `Bearer ${response.data.id_token}`
    return response
  },

  async authenticate() {
    return await api.get('/api/users/authenticate')
  },

  async register(userData) {
    return await api.post('/api/users', userData)
  },

  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userRole')
    api.defaults.headers.common['Authorization'] = ''
  },

  getCurrentUserRole() {
    return localStorage.getItem('userRole')
  },

  isAuthenticated() {
    return !!localStorage.getItem('token')
  }
}

export default authService
