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
    meta: { requiresAuth: false, requiresAdmin: false }
  },
  {
    path: '/agent-dashboard',
    component: () => import('../views/AgentDashboard.vue'),
    meta: { requiresAuth: false, requiresAgent: false }
  },
  {
    path: '/clients',
    component: ClientPage,
    meta: { requiresAuth: false }
  },
  {
    path: '/clients/new',
    component: ClientForm,
    meta: { requiresAuth: false }
  },
  {
    path: '/clients/:id/edit',
    component: () => import('../components/Client/ClientUpdateForm.vue'),
    meta: { requiresAuth: false, requiresAgent: false }
  },
  {
    path: '/accounts',
    component: AccountPage,
    meta: { requiresAuth: false }
  },
  {
    path: '/accounts/unassigned',
    component: () => import('../views/UnassignedAccountPage.vue'),
    meta: { requiresAuth: false, requiresAdmin: false }
  },
  {
    path: '/transactions',
    component: TransactionPage,
    meta: { requiresAuth: false }
  },
  {
    path: '/logs',
    component: LogsPage,
    meta: { requiresAuth: false },
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

  // Handle unknown roles
  if (userRole && userRole !== 'Admin' && userRole !== 'Agent') {
    next('/login')
    return
  }

  next()
})

export default router