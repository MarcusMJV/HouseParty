import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AuthView from '@/views/AuthView.vue'
import { useUserStore } from '@/stores/user';
import CreateRoomView from '@/views/CreateRoomView.vue'
import RoomView from '@/views/RoomView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/create/room',
      name: 'create-room',
      component: CreateRoomView,
      meta: { requiresAuth: true },
    },
    {
      path: '/room/:id',
      name: 'join-room',
      component: RoomView,
      meta: { requiresAuth: true },
    },
    {
      path: '/signup-or-login',
      name: 'signup-or-login',
      component: AuthView,
    },

  ],
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore();


  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!userStore.isAuthenticated) {
      next({ name: 'signup-or-login' });
    } else {
      next();
    }
  } else {
    next();
  }
});


export default router
