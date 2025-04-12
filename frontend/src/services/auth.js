
import api from './api'

export const authService = {
  async login(credentials) {
    try {
      console.log('Sending login request:', credentials);
      const response = await api.post('/users/login', {
        email: credentials.email,
        password: credentials.password
      })
      console.log('Login response:', response);
      
      if (!response.data?.id_token) {
        throw new Error('Invalid response from server')
      }

      // Store token and user info
      const token = response.data.id_token
      const userInfo = response.data.user_info || {}
      const userRole = userInfo.groups?.[0] || 'Agent'

      localStorage.setItem('token', token)
      localStorage.setItem('userRole', userRole)
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`

      return { 
        id_token: token,
        role: userRole,
        email: userInfo.email
      }
    } catch (error) {
      console.error('Login error:', error.response?.data || error.message)
      throw error
    }
  },

  async authenticate() {
    return await api.get('/users/authenticate')
  },

  async register(userData) {
    return await api.post('/users', userData)
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
