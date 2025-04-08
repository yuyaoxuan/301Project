<template>
  <div class="transaction-list">
    <h2>Transactions</h2>
    <table>
      <thead>
        <tr>
          <th>Date</th>
          <th>Type</th>
          <th>Amount</th>
          <th>Status</th>
          <th>Client ID</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="transaction in transactions" :key="transaction.id">
          <td>{{ new Date(transaction.transaction_date).toLocaleString() }}</td>
          <td>{{ transaction.transaction_type }}</td>
          <td>{{ transaction.amount }}</td>
          <td>{{ transaction.status }}</td>
          <td>{{ transaction.clientid }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'TransactionList',
  props: {
    transactions: {
      type: Array,
      default: () => []
    }
  }
}
</script>

<style scoped>
.transaction-list {
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

td:nth-child(3) {
  text-align: right;
}
</style>