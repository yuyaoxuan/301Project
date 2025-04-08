
<template>
  <div class="client-page">
    <h1>Client Management</h1>
    <div class="client-list">
      <ClientList :clients="clients" @client-updated="fetchClients"/>
    </div>
    <div class="bottom-buttons">
      <button @click="$router.push('/agent-dashboard')" class="return-btn">Return to Dashboard</button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import ClientList from '../components/Client/ClientList.vue'
import { clientService } from '../services/clientService'

export default {
  name: 'ClientPage',
  components: {
    ClientList
  },
  setup() {
    const clients = ref([])

    const fetchClients = async () => {
      try {
        const response = await clientService.getAllClients()
        clients.value = response.data
      } catch (error) {
        console.error('Error fetching clients:', error)
      }
    }

    onMounted(fetchClients)

    return {
      clients,
      fetchClients
    }
  }
}
</script>

<style scoped>
.client-page {
  padding: 20px;
  display: flex;
  flex-direction: column;
  min-height: 80vh;
}

.client-list {
  flex: 1;
  margin: 20px 0;
}

.bottom-buttons {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.return-btn {
  padding: 10px 20px;
  background-color: #D9D2C6;
  border: none;
  cursor: pointer;
}
</style>
