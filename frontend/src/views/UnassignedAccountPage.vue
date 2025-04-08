
<template>
  <div class="unassigned-accounts">
    <h1>Unassigned Accounts</h1>
    <div class="account-list">
      <AccountList :accounts="unassignedAccounts" />
    </div>
    <button @click="$router.push('/admin-dashboard')" class="return-btn">Return to Dashboard</button>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import AccountList from '../components/Account/AccountList.vue'
import { accountService } from '../services/accountService'

export default {
  name: 'UnassignedAccountPage',
  components: { AccountList },
  setup() {
    const unassignedAccounts = ref([])

    onMounted(async () => {
      try {
        unassignedAccounts.value = await accountService.getUnassignedAccounts()
      } catch (error) {
        console.error('Error fetching unassigned accounts:', error)
      }
    })

    return {
      unassignedAccounts
    }
  }
}
</script>

<style scoped>
.unassigned-accounts {
  padding: 20px;
}

.return-btn {
  margin-top: 20px;
  padding: 10px 20px;
  background-color: #D9D2C6;
  border: none;
  cursor: pointer;
}
</style>
