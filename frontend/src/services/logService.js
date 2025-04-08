
import api from './api'

export const logService = {
  // Client logs
  async getClientLogs(clientId) {
    const response = await api.get(`/agentclient_logs/client/${clientId}`)
    return response.data
  },

  async getClientLogsByAgent(agentId) {
    const response = await api.get(`/agentclient_logs/agent/${agentId}`)
    return response.data
  },

  async getAllClientLogs() {
    const response = await api.get('/agentclient_logs')
    return response.data
  },

  // Account logs
  async getAccountLogs(clientId) {
    const response = await api.get(`/agentclient_logs/account/client/${clientId}`)
    return response.data
  },

  async getAccountLogsByAgent(agentId) {
    const response = await api.get(`/agentclient_logs/account/agent/${agentId}`)
    return response.data
  },

  async getAllAccountLogs() {
    const response = await api.get('/agentclient_logs/account')
    return response.data
  },

  // Combined logs
  async getAllLogsByClient(clientId) {
    const response = await api.get(`/agentclient_logs/all/client/${clientId}`)
    return response.data
  },

  async getAllLogsByAgent(agentId) {
    const response = await api.get(`/agentclient_logs/all/agent/${agentId}`)
    return response.data
  },

  async getAllLogs() {
    const response = await api.get('/agentclient_logs/all')
    return response.data
  },

  async deleteLog(logId) {
    const response = await api.delete(`/agentclient_logs/${logId}`)
    return response.data
  }
}

export default logService
