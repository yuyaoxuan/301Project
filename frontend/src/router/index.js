import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import LoginPage from '../views/LoginPage.vue'
import ClientPage from '../views/ClientPage.vue'
import AccountPage from '../views/AccountPage.vue'
import TransactionPage from '../views/TransactionPage.vue'
import LogsPage from '../views/LogsPage.vue'
import ClientForm from '../components/Client/ClientForm.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    component: LoginPage
  },
  {
    path: '/admin-dashboard',
    component: () => import('../views/AdminDashboard.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/agent-dashboard',
    component: () => import('../views/AgentDashboard.vue'),
    meta: { requiresAuth: true, requiresAgent: true }
  },
  {
    path: '/clients',
    component: ClientPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/clients/new',
    component: ClientForm,
    meta: { requiresAuth: true }
  },
  {
    path: '/clients/:id/edit',
    component: () => import('../components/Client/ClientUpdateForm.vue'),
    meta: { requiresAuth: true, requiresAgent: true }
  },
  {
    path: '/accounts',
    component: AccountPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/accounts/unassigned',
    component: () => import('../views/UnassignedAccountPage.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/transactions',
    component: TransactionPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/logs',
    component: LogsPage,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'client',
        component: () => import('../components/Logs/ClientLogs.vue')
      },
      {
        path: 'account',
        component: () => import('../components/Logs/AccountLogs.vue')
      },
      {
        path: 'email',
        component: () => import('../components/Logs/EmailLogs.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('userRole')

  // Allow login page access
  if (to.path === '/login') {
    if (token) {
      // If already logged in, redirect to appropriate dashboard
      return next(userRole === 'Admin' ? '/admin-dashboard' : '/agent-dashboard')
    }
    return next()
  }

  // Check authentication for protected routes
  if (to.meta.requiresAuth && !token) {
    return next('/login')
  }

  // Verify token if exists
  if (token) {
    try {
      await authService.authenticate()
    } catch (error) {
      localStorage.removeItem('token')
      localStorage.removeItem('userRole')
      return next('/login')
    }
  }

  // Role-based access control
  if (to.meta.requiresAdmin && userRole !== 'Admin') {
    next('/agent-dashboard')
    return
  }

  if (to.meta.requiresAgent && userRole !== 'Agent') {
    next('/admin-dashboard')
    return
  }

  next()
})

export default router