
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

  async getClientByAgent(agentId) {
    const response = await api.get(`users/${agentId}/clients`)
    return response.data
  },

  async getUnassignedClients() {
    const response = await api.get(`/agentclient/unassigned`)
    return response.data
  },

  async createClient(clientData) {
    const response = await api.post('/clients', clientData)
    return response.data
  },

  async createClientForAgent(agentId) {
    const response = await api.post(`/clients/${agentId}`);
    return response.data;
  },
  
  // async updateClient(clientId, clientData) {
  //   const response = await api.put(`/clients/${clientId}`, clientData)
  //   return response.data
  // },

  async updateClientByAgent(agentId, clientId) {
    const response = await api.put(`/clients/${agentId}/${clientId}`);
    return response.data;
  },
  
  async deleteClient(clientId) {
    const response = await api.delete(`/clients/${clientId}`)
    return response.data
  },

  async verifyClient(clientId) {
    const response = await api.post(`/clients/${clientId}/verify`)
    return response.data
  },

  async AssignAgentsUnasignedClients() {
    const response = await api.put(`/agentclient`)
  },

  async getUnassignedClients() {
    const response = await api.get(`/agentclient/unassignedList`)
  }
}

export default clientService
