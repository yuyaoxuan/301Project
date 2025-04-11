
<template>
  <div class="client-update-form">
    <h2>Update Client</h2>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <input type="text" v-model="form.name" placeholder="Name" required>
      </div>
      <div class="form-group">
        <input type="email" v-model="form.email" placeholder="Email" required>
      </div>
      <div class="form-group">
        <input type="tel" v-model="form.phone" placeholder="Phone" required>
      </div>
      <button type="submit">Update Client</button>
    </form>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { clientService } from '@/services/clientService'

export default {
  name: 'ClientUpdateForm',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const form = ref({
      name: '',
      email: '',
      phone: ''
    })

    onMounted(async () => {
      try {
        const clientId = route.params.id
        const client = await clientService.getClientById(clientId)
        form.value = { ...client }
      } catch (error) {
        console.error('Failed to fetch client:', error)
      }
    })

    const handleSubmit = async () => {
      try {
        await clientService.updateClient(route.params.id, form.value)
        router.push('/clients')
      } catch (error) {
        console.error('Failed to update client:', error)
      }
    }

    return {
      form,
      handleSubmit
    }
  }
}
</script>
