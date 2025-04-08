
import api from './api'

export const clientService = {
  async getAllClients() {
    const response = await api.get('/clients')
    return response.data
  },

  async getClientById(clientId) {
    const response = await api.get(`/clients/${clientId}`)
    return response.data
  },

  async createClient(clientData) {
    const response = await api.post('/clients', clientData)
    return response.data
  },

  async updateClient(clientId, clientData) {
    const response = await api.put(`/clients/${clientId}`, clientData)
    return response.data
  },

  async deleteClient(clientId) {
    const response = await api.delete(`/clients/${clientId}`)
    return response.data
  }
}

export default clientService
