
<template>
  <div class="account-form">
    <h3>{{ editMode ? 'Edit Account' : 'Create New Account' }}</h3>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <select v-model="form.accountType" required>
          <option value="">Select Account Type</option>
          <option value="savings">Savings</option>
          <option value="checking">Checking</option>
        </select>
      </div>
      <div class="form-group">
        <input v-model="form.initialDeposit" type="number" step="0.01" placeholder="Initial Deposit" required>
      </div>
      <div class="form-group">
        <select v-model="form.currency" required>
          <option value="">Select Currency</option>
          <option value="USD">USD</option>
          <option value="EUR">EUR</option>
          <option value="GBP">GBP</option>
        </select>
      </div>
      <div class="form-group">
        <input v-model="form.branchId" placeholder="Branch ID" required>
      </div>
      <button type="submit">{{ editMode ? 'Update' : 'Create' }}</button>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'AccountForm',
  props: {
    editMode: {
      type: Boolean,
      default: false
    },
    initialData: {
      type: Object,
      default: () => ({})
    }
  },
  setup(props, { emit }) {
    const form = ref({
      accountType: '',
      initialDeposit: '',
      currency: '',
      branchId: '',
      ...props.initialData
    })

    const handleSubmit = () => {
      emit('submit', form.value)
    }

    return {
      form,
      handleSubmit
    }
  }
}
</script>

<style scoped>
.account-form {
  max-width: 500px;
  margin: 20px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 15px;
}

input, select {
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
