
<template>
  <div class="transaction-page">
    <h1>View Transactions</h1>
    
    <div class="search-container">
      <input 
        v-model="searchQuery" 
        placeholder="Enter Transaction ID"
        class="search-input"
      />
      <button @click="searchTransaction" class="btn search-btn">Search</button>
      <button @click="clearSearch" class="btn clear-btn">Clear</button>
    </div>

    <div class="transaction-list">
      <div v-for="transaction in filteredTransactions" :key="transaction.id" class="transaction-item">
        <div class="transaction-info">
          <span>Transaction ID: {{ transaction.id }}</span>
          <span>Amount: ${{ transaction.amount }}</span>
          <span>Status: {{ transaction.status }}</span>
        </div>
        <div class="action-buttons">
          <button @click="viewDetails(transaction)" class="btn view-btn">View Details</button>
          <button v-if="transaction.status === 'Failed'" @click="retryTransaction(transaction)" class="btn retry-btn">Retry</button>
        </div>
      </div>
    </div>

    <button @click="returnToDashboard" class="return-btn">Return to Dashboard</button>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useStore } from 'vuex'
import { transactionService } from '../services/transactionService'

export default {
  name: 'TransactionPage',
  setup() {
    const store = useStore()
    const transactions = ref([])
    const searchQuery = ref('')
    const userRole = computed(() => store.state.auth.userRole)
    const userId = computed(() => store.state.auth.userId)

    const filteredTransactions = computed(() => {
      return transactions.value.filter(t => 
        t.id.toString().includes(searchQuery.value)
      )
    })

    const fetchTransactions = async () => {
      try {
        if (userRole.value === 'admin') {
          transactions.value = await transactionService.getAllTransactions()
        } else {
          transactions.value = await transactionService.getAgentTransactions(userId.value)
        }
      } catch (error) {
        console.error('Error fetching transactions:', error)
      }
    }

    const searchTransaction = () => {
      // Search functionality implemented through computed property
    }

    const clearSearch = () => {
      searchQuery.value = ''
    }

    const viewDetails = (transaction) => {
      // Implement view details logic
    }

    const retryTransaction = async (transaction) => {
      try {
        await transactionService.retryTransaction(transaction.id)
        await fetchTransactions()
      } catch (error) {
        console.error('Error retrying transaction:', error)
      }
    }

    const returnToDashboard = () => {
      const role = store.state.auth.userRole;
      if (role === 'admin') {
        router.push('/admin-dashboard');
      } else {
        router.push('/agent-dashboard');
      }
    }

    onMounted(fetchTransactions)

    return {
      transactions,
      searchQuery,
      filteredTransactions,
      searchTransaction,
      clearSearch,
      viewDetails,
      retryTransaction,
      returnToDashboard
    }
  }
}
</script>

<style scoped>
.transaction-page {
  padding: 20px;
}

.search-container {
  margin: 20px 0;
  display: flex;
  gap: 10px;
}

.search-input {
  padding: 8px;
  flex: 1;
}

.btn {
  padding: 8px 16px;
  cursor: pointer;
  border: none;
}

.transaction-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  margin: 10px 0;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.action-buttons {
  display: flex;
  gap: 10px;
}

.return-btn {
  margin-top: 20px;
  padding: 10px 20px;
  background-color: #D9D2C6;
}
</style>
