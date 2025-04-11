import api from './api'

export const clientService = {
  async createClient(agentId, clientData) {
    return await api.post(`/api/clients/${agentId}`, clientData)
  },

  async getClient(clientId) {
    return await api.get(`/api/clients/${clientId}`)
  },

  async updateClient(agentId, clientId, clientData) {
    return await api.put(`/api/clients/${agentId}/${clientId}`, clientData)
  },

  async deleteClient(clientId) {
    return await api.delete(`/api/clients/${clientId}`)
  },

  async verifyClient(clientId) {
    return await api.post(`/api/clients/${clientId}/verify`)
  },

  async getAllClients() {
    const response = await api.get('/clients')
    return response.data
  },

  async getClientById(clientId) {
    const response = await api.get(`/clients/${clientId}`)
    return response.data
  },

  async getClientByAgent(agentId) {
    const response = await api.get(`users/${agentId}/clients`)
    return response.data
  },
  async getUnassignedClients() {
    const response = await api.get(`/agentclient/unassignedList`)
    return response.data
  },

  async updateClientByAgent(agentId, clientId) {
    const response = await api.put(`/clients/${agentId}/${clientId}`);
    return response.data;
  },

  async AssignAgentsUnasignedClients() {
    const response = await api.put(`/agentclient`)
  },
}

export default clientService