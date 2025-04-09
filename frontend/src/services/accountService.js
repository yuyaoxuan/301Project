
import api from './api'

export const accountService = {
  async getAllAccounts() {
    const response = await api.get('/accounts')
    return response.data
  },

  async getAccountById(clientId) {
    const response = await api.get(`/clients/${clientId}/accounts`)
    return response.data
  },

  async createAccount(accountData) {
    const response = await api.post('/accounts', accountData)
    return response.data
  },

  async updateAccount(accountId, accountData) {
    const response = await api.put(`/accounts/${accountId}`, accountData)
    return response.data
  },

  async deleteAccount(accountId) {
    const response = await api.delete(`/accounts/${accountId}`)
    return response.data
  }
}

export default accountService
