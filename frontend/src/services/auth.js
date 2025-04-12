
import api from './api'

export const authService = {
  async login(credentials) {
    try {
      console.log('Auth service: Sending login request to API...');
      const response = await api.post('/users/login', {
        email: credentials.email,
        password: credentials.password
      })
      console.log('Auth service: Received API response:', {
        status: response.status,
        data: response.data
      });
      
      if (!response.data?.access_token) {
        throw new Error('Invalid response from server')
      }

      // Store token and decode JWT to get user info
      const token = response.data.access_token
      const payload = JSON.parse(atob(token.split('.')[1]))
      const userRole = payload['cognito:groups']?.[0] || 'Agent'
      const userEmail = payload.email

      // Store auth data and set API header
      localStorage.setItem('token', token)
      localStorage.setItem('userRole', userRole)
      localStorage.setItem('userEmail', userEmail)
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`
      
      // Set auth header
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`

      return {
        token: token,
        role: userRole.toLowerCase(), // Ensure consistent casing
        email: userEmail
      }
    } catch (error) {
      console.error('Login error:', error.response?.data || error.message)
      throw error
    }
  },

  async authenticate() {
    try {
      const response = await api.get('/users/authenticate')
      // Handle redirect to callback
      if (response.data?.id_token) {
        const token = response.data.id_token
        const userRole = response.data.user_info?.groups?.[0] || 'Agent'
        
        localStorage.setItem('token', token)
        localStorage.setItem('userRole', userRole)
        api.defaults.headers.common['Authorization'] = `Bearer ${token}`
        
        return {
          id_token: token,
          role: userRole.toLowerCase()
        }
      }
      return response
    } catch (error) {
      console.error('Authentication error:', error)
      throw error
    }
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
