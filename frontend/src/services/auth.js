
import api from './api'

export const authService = {
  async login(credentials) {
    const response = await api.post('/api/users/login', credentials)
    if (response.data.id_token) {
      localStorage.setItem('token', response.data.id_token)
      localStorage.setItem('userRole', response.data.userRole)
      return response.data
    }
    throw new Error('Authentication failed')
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
