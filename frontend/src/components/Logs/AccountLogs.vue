<template>
  <div class="account-logs">
    <h2>Account Activity Logs</h2>
    <table>
      <thead>
        <tr>
          <th>Date</th>
          <th>Agent</th>
          <th>Client</th>
          <th>Action</th>
          <th>Details</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="log in logs" :key="log.id">
          <td>{{ new Date(log.timestamp).toLocaleString() }}</td>
          <td>{{ log.agentId }}</td>
          <td>{{ log.clientId }}</td>
          <td>{{ log.action }}</td>
          <td>{{ JSON.stringify(log.modifiedFields) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'AccountLogs',
  props: {
    logs: {
      type: Array,
      default: () => []
    }
  }
}
</script>

<style scoped>
.account-logs {
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

td {
  font-size: 0.9em;
}
</style>