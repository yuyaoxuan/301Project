<template>
  <div class="account-page">
    <h1>Account Management</h1>
    <AccountList :accounts="accounts" @account-updated="fetchAccounts"/>
    <div class="bottom-buttons">
      <button @click="$router.push('/admin-dashboard')" class="return-btn">Return to Dashboard</button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import AccountForm from '../components/Account/AccountForm.vue'
import AccountList from '../components/Account/AccountList.vue'
import { accountService } from '../services/accountService'

export default {
  name: 'AccountPage',
  components: {
    AccountForm,
    AccountList
  },
  setup() {
    const accounts = ref([])
    const showForm = ref(false)

    const fetchAccounts = async () => {
      try {
        accounts.value = await accountService.getAllAccounts()
      } catch (error) {
        console.error('Error fetching accounts:', error)
      }
    }

    const handleAccountCreated = () => {
      showForm.value = false
      fetchAccounts()
    }

    onMounted(fetchAccounts)

    return {
      accounts,
      showForm,
      handleAccountCreated,
      fetchAccounts
    }
  }
}
</script>

<style scoped>
.account-page {
  padding: 20px;
}

.account-actions {
  margin: 20px 0;
}

button {
  padding: 10px 20px;
  background-color: #D9D2C6;
  border: none;
  cursor: pointer;
}
</style>