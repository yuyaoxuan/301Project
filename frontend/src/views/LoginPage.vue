
<template>
  <div class="login-page">
    <div class="welcome-section">
      <h1>Welcome</h1>
    </div>
    <div class="login-form">
      <h2>CRM Login</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>Username</label>
          <input 
            type="text" 
            v-model="credentials.username" 
            placeholder="Enter username"
            required
          />
        </div>
        <div class="form-group">
          <label>Password</label>
          <input 
            type="password" 
            v-model="credentials.password" 
            placeholder="Enter your Password"
            required
          />
        </div>
        <div class="login-buttons">
          <button type="submit" class="admin-login">Admin Login</button>
          <button type="submit" class="agent-login">Agent Login</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

export default {
  name: 'LoginPage',
  setup() {
    const router = useRouter()
    const store = useStore()
    const credentials = ref({
      username: '',
      password: ''
    })

    const handleLogin = async () => {
      try {
        await store.dispatch('auth/login', credentials.value)
        router.push('/dashboard')
      } catch (error) {
        console.error('Login failed:', error)
      }
    }

    return {
      credentials,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-page {
  display: grid;
  grid-template-columns: 1fr 1fr;
  height: 100vh;
}

.welcome-section {
  background-color: #D9D2C6;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-form {
  padding: 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-group {
  margin-bottom: 20px;
}

input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.login-buttons {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

button {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  background-color: #D9D2C6;
}
</style>
