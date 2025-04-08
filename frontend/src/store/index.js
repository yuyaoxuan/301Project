
import { createStore } from 'vuex'

export default createStore({
  state: {
    user: null,
    token: localStorage.getItem('token') || null,
    clients: [],
    accounts: [],
    transactions: [],
    logs: []
  },
  mutations: {
    setUser(state, user) {
      state.user = user
    },
    setToken(state, token) {
      state.token = token
      localStorage.setItem('token', token)
    },
    clearAuth(state) {
      state.user = null
      state.token = null
      localStorage.removeItem('token')
    },
    setClients(state, clients) {
      state.clients = clients
    },
    setAccounts(state, accounts) {
      state.accounts = accounts
    },
    setTransactions(state, transactions) {
      state.transactions = transactions
    },
    setLogs(state, logs) {
      state.logs = logs
    }
  },
  actions: {
    async login({ commit }, credentials) {
      // TODO: Implement login logic
    },
    async logout({ commit }) {
      commit('clearAuth')
    }
  }
})
