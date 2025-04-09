
<template>
  <div class="email-logs">
    <h2>Email Communication Logs</h2>
    <table>
      <thead>
        <tr>
          <th>Date</th>
          <th>Agent</th>
          <th>Client</th>
          <th>Subject</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="log in logs" :key="log.id">
          <td>{{ new Date(log.timestamp).toLocaleString() }}</td>
          <td>{{ log.agent_id }}</td>
          <td>{{ log.client_id }}</td>
          <td>{{ log.email_subject }}</td>
          <td>{{ log.email_status }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'EmailLogs',
  props: {
    logs: {
      type: Array,
      default: () => []
    }
  }
}
</script>

<style scoped>
.email-logs {
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
