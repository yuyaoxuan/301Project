
<template>
  <div class="client-form">
    <h2>Create Client Profile</h2>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label>First name</label>
        <input v-model="form.firstName" type="text" required>
      </div>
      <div class="form-group">
        <label>Last name</label>
        <input v-model="form.lastName" type="text" required>
      </div>
      <div class="form-group">
        <label>Date Of Birth</label>
        <input v-model="form.dob" type="date" required>
      </div>
      <div class="form-group">
        <label>Gender</label>
        <input v-model="form.gender" type="text" required>
      </div>
      <div class="form-group">
        <label>Email</label>
        <input v-model="form.email" type="email" required>
      </div>
      <div class="form-group">
        <label>Phone</label>
        <input v-model="form.phone" type="tel" required>
      </div>
      <div class="form-group">
        <label>Address</label>
        <input v-model="form.address" type="text" required>
      </div>
      <div class="form-group">
        <label>City</label>
        <input v-model="form.city" type="text" required>
      </div>
      <div class="form-group">
        <label>State</label>
        <input v-model="form.state" type="text" required>
      </div>
      <div class="form-group">
        <label>Country</label>
        <input v-model="form.country" type="text" required>
      </div>
      <div class="form-group">
        <label>Postal Code</label>
        <input v-model="form.postalCode" type="text" required>
      </div>
      <div class="button-group">
        <button type="submit" class="save-btn">Save</button>
        <button type="button" class="cancel-btn" @click="$router.go(-1)">Cancel</button>
      </div>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'ClientForm',
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
      firstName: '',
      lastName: '',
      email: '',
      phone: '',
      address: '',
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
.client-form {
  max-width: 600px;
  margin: 20px auto;
  padding: 20px;
}

h2 {
  text-align: center;
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 15px;
  display: flex;
  flex-direction: column;
}

label {
  color: #654321;
  margin-bottom: 5px;
}

input {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.button-group {
  display: flex;
  justify-content: space-between;
  margin-top: 30px;
}

button {
  padding: 10px 30px;
  border: 1px solid #000;
  background: white;
  cursor: pointer;
  font-size: 16px;
}

.save-btn {
  background: white;
}

.cancel-btn {
  background: white;
}
</style>
