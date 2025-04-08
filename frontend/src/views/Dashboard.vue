
<template>
  <div class="dashboard">
    <div v-if="userRole === 'admin'">
      <h1>Admin Dashboard</h1>
      <div class="nav-buttons">
        <button @click="$router.push('/accounts/new')">CREATE NEW ACCOUNT</button>
        <button @click="$router.push('/accounts')">MANAGE ACCOUNTS</button>
        <button @click="$router.push('/transactions')">VIEW TRANSACTIONS</button>
        <button @click="$router.push('/clients/unassigned')">MANAGE UNASSIGNED CLIENT</button>
        <button @click="$router.push('/settings')">SETTINGS</button>
      </div>
      <div class="recent-activities">
        <h3>RECENT ACTIVITIES</h3>
        <ul>
          <li v-for="(log, index) in recentLogs" :key="index">
            {{ log }}
          </li>
        </ul>
      </div>
    </div>
    <div v-else>
      <h1>Agent Dashboard</h1>
      <div class="nav-buttons">
        <button @click="$router.push('/clients/new')">CREATE CLIENT PROFILE</button>
        <button @click="$router.push('/clients')">MANAGE PROFILE</button>
        <button @click="$router.push('/transactions')">VIEW TRANSACTIONS</button>
      </div>
      <div class="recent-activities">
        <h3>RECENT ACTIVITIES</h3>
        <ul>
          <li v-for="(activity, index) in recentActivities" :key="index">
            {{ activity }}
          </li>
        </ul>
      </div>
    </div>
    <button class="logout-btn" @click="handleLogout">Logout</button>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

export default {
  name: 'Dashboard',
  setup() {
    const router = useRouter()
    const store = useStore()
    const userRole = ref('')
    const recentLogs = ref([])
    const recentActivities = ref([])

    onMounted(async () => {
      userRole.value = store.state.auth.userRole
      // Fetch recent activities based on role
      if (userRole.value === 'admin') {
        // Get admin logs
      } else {
        // Get agent logs
      }
    })

    const handleLogout = () => {
      store.dispatch('auth/logout')
      router.push('/login')
    }

    return {
      userRole,
      recentLogs,
      recentActivities,
      handleLogout
    }
  }
}
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.nav-buttons {
  display: flex;
  gap: 10px;
  margin: 20px 0;
  background-color: #D9D2C6;
  padding: 10px;
}

button {
  padding: 10px 20px;
  background: none;
  border: none;
  cursor: pointer;
  font-weight: bold;
}

.recent-activities {
  background-color: #F5F5F5;
  padding: 20px;
  margin-top: 20px;
}

.logout-btn {
  position: absolute;
  top: 20px;
  right: 20px;
}
</style>
