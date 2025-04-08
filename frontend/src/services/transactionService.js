
import api from './api'

export const transactionService = {
  async getAllTransactions() {
    const response = await api.get('/transactions')
    return response.data
  },

  async getTransactionsByAccountId(accountId) {
    const response = await api.get(`/transactions/account/${accountId}`)
    return response.data
  },

  async getTransactionsByClientId(clientId) {
    const response = await api.get(`/transactions/client/${clientId}`)
    return response.data
  }
}

export default transactionService
