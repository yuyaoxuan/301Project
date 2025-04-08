
<template>
  <div class="app">
    <nav v-if="$store.state.token" class="navbar">
      <router-link to="/dashboard">Dashboard</router-link>
      <router-link to="/clients">Clients</router-link>
      <router-link to="/accounts">Accounts</router-link>
      <router-link to="/transactions">Transactions</router-link>
      <router-link to="/logs">Logs</router-link>
      <button @click="logout" class="logout-btn">Logout</button>
    </nav>
    <main class="main-content">
      <router-view></router-view>
    </main>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

export default {
  name: 'App',
  setup() {
    const router = useRouter()
    const store = useStore()

    const logout = () => {
      store.dispatch('logout')
      router.push('/login')
    }

    return {
      logout
    }
  }
}
</script>

<style>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background-color: #2c3e50;
  padding: 1rem;
  display: flex;
  gap: 1rem;
}

.navbar a {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
}

.navbar a.router-link-active {
  background-color: #34495e;
  border-radius: 4px;
}

.logout-btn {
  margin-left: auto;
  background-color: #e74c3c;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  cursor: pointer;
  border-radius: 4px;
}

.main-content {
  flex: 1;
  padding: 2rem;
  background-color: #f5f5f5;
}
</style>
