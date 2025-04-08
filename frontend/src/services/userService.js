
import api from './api'

export const userService = {
  async getAllUsers() {
    const response = await api.get('/users')
    return response.data
  },

  async getUserById(userId) {
    const response = await api.get(`/users/${userId}`)
    return response.data
  },

  async updateUser(userId, userData) {
    const response = await api.put(`/users/${userId}`, userData)
    return response.data
  },

  async deleteUser(userId) {
    const response = await api.delete(`/users/${userId}`)
    return response.data
  }
}

export default userService
