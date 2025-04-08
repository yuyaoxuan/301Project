
<template>
  <div class="client-list">
    <h2>Clients</h2>
    <table>
      <thead>
        <tr>
          <th>Client ID</th>
          <th>Name</th>
          <th>Email</th>
          <th>Phone</th>
          <th>Address</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="client in clients" :key="client.clientId">
          <td>{{ client.clientId }}</td>
          <td>{{ `${client.firstName} ${client.lastName}` }}</td>
          <td>{{ client.email }}</td>
          <td>{{ client.phone }}</td>
          <td>{{ client.address }}</td>
          <td>
            <button @click="$emit('edit', client)" class="edit-btn">Edit</button>
            <button @click="$emit('delete', client.clientId)" class="delete-btn">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'ClientList',
  props: {
    clients: {
      type: Array,
      default: () => []
    }
  }
}
</script>

<style scoped>
.client-list {
  padding: 20px;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

th {
  background-color: #f5f5f5;
}

.edit-btn, .delete-btn {
  padding: 5px 10px;
  margin: 0 5px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.edit-btn {
  background-color: #42b983;
  color: white;
}

.delete-btn {
  background-color: #ff4444;
  color: white;
}
</style>
