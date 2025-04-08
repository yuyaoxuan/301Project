
<template>
  <div class="account-list">
    <h2>Accounts</h2>
    <table>
      <thead>
        <tr>
          <th>Account ID</th>
          <th>Type</th>
          <th>Status</th>
          <th>Balance</th>
          <th>Currency</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="account in accounts" :key="account.accountId">
          <td>{{ account.accountId }}</td>
          <td>{{ account.accountType }}</td>
          <td>{{ account.accountStatus }}</td>
          <td>{{ account.balance }}</td>
          <td>{{ account.currency }}</td>
          <td>
            <button @click="$emit('edit', account)" class="edit-btn">Edit</button>
            <button @click="$emit('delete', account.accountId)" class="delete-btn">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'AccountList',
  props: {
    accounts: {
      type: Array,
      default: () => []
    }
  }
}
</script>

<style scoped>
.account-list {
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
