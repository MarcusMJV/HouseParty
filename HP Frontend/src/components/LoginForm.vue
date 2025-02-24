<script setup lang="ts">
import { useUserStore } from '@/stores/user';
import { ref } from 'vue';
import router from '@/router';

const password = ref('');
const errorMessage = ref('');
const nameField = ref('');
const successMessage = ref('');
let body = '';

const emit = defineEmits(['switch-to-signup']);

const emitSwitchEvent = () => {
  emit('switch-to-signup');
};

const isValidEmail = (email: string) => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;
  return re.test(email);
};

const handleSubmit = async (e: Event) => {
  e.preventDefault();
  errorMessage.value = '';
  successMessage.value = '';

  if (!nameField.value || !password.value) {
    errorMessage.value = 'All fields are required';
    return;
  }

  try {
    if(isValidEmail(nameField.value)){
        body = JSON.stringify({
        username: "",
        email: nameField.value,
        password: password.value
      })
    }else{
        body = JSON.stringify({
        username: nameField.value,
        email: "",
        password: password.value
      })
    }

    const response = await fetch('http://localhost:8080/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: body,
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Signup failed');
    }

    successMessage.value = 'Login successful! Redirecting...';

    nameField.value = '';
    password.value = '';

    const userStore = useUserStore();
    userStore.setJwt(data.token);
    console.log(data.user)
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
      <h1 class="text-2xl font-bold text-center text-sky-500 mb-6">LOGIN</h1>

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
            id="nameField"
            v-model="nameField"
            placeholder="Enter your username or password"
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
          Login
        </button>
      </form>
      <p>Don't have a account? <span @click="emitSwitchEvent" class="text-sky-500/100 hover:underline hover:cursor-pointer">Sign Up</span></p>
    </div>
  </div>
</template>
