<template>
  <div class="dashboard">
    <h1>Agent Dashboard</h1>
    <div class="nav-buttons">
      <button @click="$router.push('/clients/new')">CREATE CLIENT PROFILE</button>
      <button @click="$router.push('/clients')">MANAGE CLIENT PROFILES</button>
      <button @click="$router.push('/transactions')">VIEW TRANSACTIONS</button>
      <button @click="logout" class="logout-btn">LOGOUT</button>
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
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

export default {
  name: 'AgentDashboard',
  setup() {
    const recentActivities = ref([])
    const router = useRouter()
    const store = useStore()

    const logout = () => {
      store.dispatch('logout');
      router.push('/login');
    };

    onMounted(async () => {
      // Fetch agent activities
    })

    return {
      recentActivities,
      logout
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
</style>