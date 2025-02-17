<script setup lang="ts">
import { ref } from 'vue';
import { useUserStore } from '@/stores/user';
import router from '@/router';


const roomName = ref('');
const description = ref('');
const isPrivate = ref(false);
const userStore = useUserStore();
const errorMessage = ref('');

const handleRoomSubmit = async () => {
  try {
    const response = await fetch('http://localhost:8080/room/create', {
      method: "POST",
      headers: {
        'Authorization': `Bearer ${userStore.jwt}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: roomName.value,
        description: description.value,
        public: isPrivate.value
      })
    });

    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.error || 'Failed to create room');
    }

    router.push({ name: 'home', query: { successMessage: "Room Created" } });

  }catch(error){
    errorMessage.value = error instanceof Error ? error.message : 'An unexpected error occurred';
  }
}

</script>


<template>
  <div class="lex items-start pt-10 bg-slate-950">
      <div class="bg-slate-900 p-8 rounded-2xl shadow-lg w-full mx-auto text-center">
        <form class="space-y-6 w-100" @submit.prevent="handleRoomSubmit">
          <h1 class="text-2xl font-bold text-center text-sky-500 mb-6">Create Room</h1>

          <div v-if="errorMessage" class="mb-4 p-3 bg-red-500/20 text-red-300 rounded-lg">
            {{ errorMessage }}
          </div>

          <div>
            <input
                type="text"
                id="roomName"
                v-model="roomName"
                placeholder="Enter room name"
                class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500 text-white bg-transparent"
            />
          </div>

          <div>
            <textarea
              id="description"
              v-model="description"
              placeholder="Enter room description"
              class="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500 text-white bg-transparent resize-none"
              rows="3"
            ></textarea>
          </div>

          <div class="flex items-center justify-between p-3 border rounded-lg">
            <span class="text-white">Public Room</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="isPrivate" class="sr-only peer">
              <div class="w-11 h-6 bg-slate-700 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[2px] after:bg-white after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
            </label>
          </div>

          <button type="submit" class="w-full mt-4 hover:bg-sky-500/100 border text-white font-semibold py-3 rounded-lg hover:bg-sky-600 transition duration-300">
            Create Room
          </button>
        </form>
      </div>
    </div>
</template>
