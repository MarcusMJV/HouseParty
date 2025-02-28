<script setup lang="ts">
import { ref } from 'vue';
import { useUserStore } from '@/stores/user';
import router from '@/router';

const username = ref('');
const email = ref('');
const password = ref('');
const errorMessage = ref('');
const successMessage = ref('');

const emit = defineEmits(['switch-to-login']);

const emitSwitchEvent = () => {
  emit('switch-to-login');
};

const handleSubmit = async (e: Event) => {
  e.preventDefault();
  errorMessage.value = '';
  successMessage.value = '';

  if (!username.value || !email.value || !password.value) {
    errorMessage.value = 'All fields are required';
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        email: email.value,
        password: password.value
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Signup failed');
    }

    successMessage.value = 'Signup successful! Redirecting...';

    username.value = '';
    email.value = '';
    password.value = '';

    const userStore = useUserStore();
    userStore.setJwt(data.token);
    userStore.setCredentials(data.user);

    router.push({ name: 'home' });

  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'An unexpected error occurred';
  }
};
</script>

<template>
  <div class="lex items-start pt-10 bg-slate-950">
    <div class="bg-slate-900 p-8 rounded-2xl shadow-lg max-w-md w-full mx-auto text-center">
      <h1 class="text-2xl font-bold text-center text-sky-500 mb-6">SIGN UP</h1>

      <!-- Error/Success Messages -->
      <div v-if="errorMessage" class="mb-4 p-3 bg-red-500/20 text-red-300 rounded-lg">
        {{ errorMessage }}
      </div>
      <div v-if="successMessage" class="mb-4 p-3 bg-green-500/20 text-green-300 rounded-lg">
        {{ successMessage }}
      </div>

      <form class="space-y-6" @submit.prevent="handleSubmit">
        <div>
          <input
            type="text"
            id="username"
            v-model="username"
            placeholder="Enter your username"
            class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500 text-white"
          />
        </div>

        <div>
          <input
            type="email"
            id="email"
            v-model="email"
            placeholder="Enter your email"
            class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500 text-white"
          />
        </div>

        <div>
          <input
            type="password"
            id="password"
            v-model="password"
            placeholder="Enter your password"
            class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500 text-white"
          />
        </div>

        <button type="submit" class="w-full mt-4 hover:bg-sky-500/100 border text-white font-semibold py-3 rounded-lg hover:bg-sky-600 transition duration-300">
          Sign Up
        </button>
      </form>
      <p>Already have a account? <span @click="emitSwitchEvent" class="text-sky-500/100 hover:underline hover:cursor-pointer">Login</span></p>
    </div>
  </div>
</template>
