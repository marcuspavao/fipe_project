import { createRouter, createWebHistory } from 'vue-router'
// import HomeView from '../views/HomeView.vue' // No longer needed
import MainLayout from '../layouts/MainLayout.vue'
import FipeConsultaView from '../views/FipeConsultaView.vue'
import CompararValoresView from '../views/CompararValoresView.vue'
import DashboardMarcasView from '../views/DashboardMarcasView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: MainLayout,
      children: [
        { 
          path: '', // Matches the root path '/'
          name: 'fipe-consulta', 
          component: FipeConsultaView 
        }
      ]
    },
    {
      path: '/comparar',
      component: MainLayout,
      children: [
        { 
          path: '', // Matches '/comparar'
          name: 'comparar-valores', 
          component: CompararValoresView 
        }
      ]
    },
    {
      path: '/dashboard',
      component: MainLayout,
      children: [
        { 
          path: '', // Matches '/dashboard'
          name: 'dashboard-marcas', 
          component: DashboardMarcasView 
        }
      ]
    }
    // Example of a route not using MainLayout:
    // {
    //   path: '/login',
    //   name: 'login',
    //   component: () => import('../views/LoginView.vue') // Lazy load
    // }
  ]
})

export default router
