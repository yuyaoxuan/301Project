
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
        <div class="error-message" v-if="error">{{ error }}</div>
        <button type="submit" class="login-button">Login</button>
      </form>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { authService } from '../services/auth'

export default {
  name: 'LoginPage',
  setup() {
    const router = useRouter()
    const store = useStore()
    const credentials = ref({
      username: '',
      password: ''
    })
    const error = ref('')

    const handleLogin = async () => {
      try {
        console.log('Starting login process...')
        error.value = ''
        
        console.log('Credentials:', {
          email: credentials.value.username,
          password: '***' // Masked for security
        })
        
        const response = await authService.login({
          email: credentials.value.username,
          password: credentials.value.password
        })
        
        console.log('Login successful:', {
          role: response.role,
          email: response.email
        })
        
        const dashboard = response.role === 'admin' ? '/admin-dashboard' : '/agent-dashboard'
        console.log('Login: User role:', response.role)
        console.log('Login: Target dashboard:', dashboard)
        console.log('Login: Starting navigation...')
        
        try {
          await router.push(dashboard)
          console.log('Login: Navigation successful')
        } catch (navError) {
          console.error('Login: Navigation failed:', navError)
        }
      } catch (err) {
        console.error('Login error:', err)
        error.value = err.response?.data?.message || 'Login failed'
        console.error('Login failed:', err)
      }
    }

    return {
      credentials,
      handleLogin,
      error
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

.login-button {
  width: 100%;
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  background-color: #D9D2C6;
  margin-top: 20px;
}

.error-message {
  color: red;
  margin-top: 10px;
}
</style>
