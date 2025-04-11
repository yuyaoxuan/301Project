
import api from './api'

export const logService = {
  // Agent Client Logs
  async getClientLogs(clientId) {
    return await api.get(`/agentclient_logs/client/${clientId}`)
  },

  async getAgentLogs(agentId) {
    return await api.get(`/agentclient_logs/agent/${agentId}`)
  },

  async getAllLogs() {
    return await api.get('/agentclient_logs')
  },

  // Account Logs
  async getAccountLogsByClient(clientId) {
    return await api.get(`/agentclient_logs/account/client/${clientId}`)
  },

  async getAccountLogsByAgent(agentId) {
    return await api.get(`/agentclient_logs/account/agent/${agentId}`)
  },

  async getAllAccountLogs() {
    return await api.get('/agentclient_logs/account')
  },

  // Combined Logs
  async getAllClientAndAccountLogs(clientId) {
    return await api.get(`/agentclient_logs/all/client/${clientId}`)
  },

  async getAllAgentAndAccountLogs(agentId) {
    return await api.get(`/agentclient_logs/all/agent/${agentId}`)
  },

  async getCommunicationLog(logId) {
    return await api.get(`/communication_logs/${logId}`)
  },

  async deleteLog(logId) {
    return await api.delete(`/agentclient_logs/${logId}`)
  }
}

export default logService
