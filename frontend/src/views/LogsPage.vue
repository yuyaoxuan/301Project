
<template>
  <div class="logs-page">
    <h1>System Logs</h1>
    <div class="log-navigation">
      <router-link to="/logs/client" class="nav-link">Client Logs</router-link>
      <router-link to="/logs/account" class="nav-link">Account Logs</router-link>
      <router-link to="/logs/email" class="nav-link">Email Logs</router-link>
    </div>
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: 'LogsPage'
}
</script>

<style scoped>
.logs-page {
  padding: 20px;
}

.log-navigation {
  display: flex;
  gap: 20px;
  margin: 20px 0;
}

.nav-link {
  padding: 10px 20px;
  background-color: #D9D2C6;
  text-decoration: none;
  color: black;
}

.router-link-active {
  background-color: #a39b8e;
}
</style>
