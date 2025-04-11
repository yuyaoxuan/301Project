
import api from './api'

export const authService = {
  async login(credentials) {
    try {
      console.log('Sending login request:', credentials);
      const response = await api.post('/users/login', {
        email: credentials.email,
        password: credentials.password
      })
      if (!response.data?.id_token) {
        throw new Error('Invalid response from server')
      }
      localStorage.setItem('token', response.data.id_token)
      api.defaults.headers.common['Authorization'] = `Bearer ${response.data.id_token}`
      return response
    } catch (error) {
      console.error('Login error:', error.response?.data || error.message)
      throw error
    }
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
