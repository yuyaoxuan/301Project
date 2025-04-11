<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin">
      <h2>Login</h2>
      <div class="form-group">
        <input type="text" v-model="username" placeholder="Username" required>
      </div>
      <div class="form-group">
        <input type="password" v-model="password" placeholder="Password" required>
      </div>
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import authService from '@/services/authService' // Assuming authService is imported correctly

export default {
  name: 'Login',
  setup() {
    const username = ref('')
    const password = ref('')
    const router = useRouter()

    const handleLogin = async () => {
      try {
        const response = await authService.login({
          username: username.value,
          password: password.value
        })
        const userRole = authService.getCurrentUserRole()
        if (userRole === 'admin') {
          router.push('/admin-dashboard')
        } else if (userRole === 'agent') {
          router.push('/agent-dashboard')
        }
      } catch (error) {
        console.error('Login failed:', error)
      }
    }

    return {
      username,
      password,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 40px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 15px;
}

input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #3aa876;
}
</style>