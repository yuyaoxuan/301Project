
<template>
  <div class="register-container">
    <form @submit.prevent="handleRegister">
      <h2>Register</h2>
      <div class="form-group">
        <input type="text" v-model="form.username" placeholder="Username" required>
      </div>
      <div class="form-group">
        <input type="password" v-model="form.password" placeholder="Password" required>
      </div>
      <div class="form-group">
        <input type="password" v-model="form.confirmPassword" placeholder="Confirm Password" required>
      </div>
      <div class="form-group">
        <input type="email" v-model="form.email" placeholder="Email" required>
      </div>
      <button type="submit">Register</button>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'Register',
  setup() {
    const form = ref({
      username: '',
      password: '',
      confirmPassword: '',
      email: ''
    })
    const router = useRouter()

    const handleRegister = async () => {
      if (form.value.password !== form.value.confirmPassword) {
        alert('Passwords do not match')
        return
      }

      try {
        // TODO: Implement register API call
        router.push('/login')
      } catch (error) {
        console.error('Registration failed:', error)
      }
    }

    return {
      form,
      handleRegister
    }
  }
}
</script>

<style scoped>
.register-container {
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
