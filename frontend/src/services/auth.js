
import api from './api'

export const authService = {
  async login(credentials) {
    const response = await api.post('/api/users/login', credentials)
    if (!response.data?.id_token) {
      throw new Error('Invalid response from server')
    }
    // Set token in api instance for subsequent requests
    api.defaults.headers.common['Authorization'] = `Bearer ${response.data.id_token}`
    return response
  },

  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userRole')
  },

  getCurrentUserRole() {
    return localStorage.getItem('userRole')
  },

  isAuthenticated() {
    return !!localStorage.getItem('token')
  }
}

export default authService
