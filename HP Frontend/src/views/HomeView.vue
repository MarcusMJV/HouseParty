<script setup lang="ts">
import { onMounted, ref } from 'vue';
import HousePartyLogo from '@/components/HousePartyLogo.vue';
import { useUserStore } from '@/stores/user';
import { useRoute } from 'vue-router';
import router from '@/router';

interface Room {
  id: string;
  name: string;
  description: string;
  host_id: string;
  host_name: string;
  is_public: boolean;
  created_at: string;
}

interface ApiResponse {
  message: string;
  user_rooms: Room[];
  public_rooms: Room[];
}

const userRooms = ref<Room[]>([]);
const publicRooms = ref<Room[]>([]);
const username = ref<string>('');
const errorMessage = ref('');
const successMessage = ref('');
const userStore = useUserStore();
const route = useRoute();

username.value = userStore.credentials?.username || '';


if(route.query.successMessage){
  successMessage.value = route.query.successMessage as string;

  setTimeout(() => {
    successMessage.value = '';
    router.push({ name: 'home' });
  }, 3000);
}

onMounted(() => {
  getRooms();


});


const deleteRoom = async (roomId: string) => {
  try{
    const response = await fetch(`http://localhost:8080/room/delete/${roomId}`, {
      method: "DELETE",
      headers: {
        'Authorization': `Bearer ${userStore.jwt}`
      }
    })

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Failed to delete room');
    }

    successMessage.value = data.message;
    userRooms.value = userRooms.value.filter(room => room.id !== roomId);

    setTimeout(() => {
      successMessage.value = '';
    }, 3000);


  }catch (error){
    errorMessage.value = error instanceof Error ? error.message : 'An unexpected error occurred';
  }
}

const getRooms = async () => {
  try {
    const response = await fetch('http://localhost:8080/rooms', {
      method: "GET",
      headers: {
        'Authorization': `Bearer ${userStore.jwt}`
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: ApiResponse = await response.json();

    console.log(data);

    userRooms.value = data.user_rooms;
    publicRooms.value = data.public_rooms;

  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'An unexpected error occurred';
  }
}
</script>

<template>
  <div class="flex flex-col items-center pt-10 min-h-screen font-mono bg-slate-950 relative">
    <HousePartyLogo/>

      <h1 class="text-white text-xl">Select a room to join or create a new one</h1>

      <div v-if="errorMessage" class="mt-4 p-3 bg-red-500/20 text-red-300 rounded-lg">
        {{ errorMessage }}
      </div>
      <div v-if="successMessage" class="mt-4 p-3 bg-green-500/20 text-green-300 rounded-lg">
        {{ successMessage }}
      </div>


      <div class="relative w-120 mt-8">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-sky-500/100"></div>
        </div>
        <div class="relative flex justify-center">
          <span class="px-2 text-white bg-slate-950 text-sm">
            {{ username }}'s Rooms
          </span>
        </div>
      </div>

      <div class="flex flex-col items-center mt-4 space-y-4">
        <div v-for="room in userRooms" :key="room.id" class="w-150 bg-slate-900 p-4 rounded-2xl shadow-lg flex items-center justify-between relative">
          <button @click="deleteRoom(room.id)" class="absolute -top-1 -right-1 text-white hover:text-red-500 hover:cursor-pointer text-xs font-bold p-1 rounded-full transition-colors">x</button>

          <span class="text-white">{{ room.name }}:</span>
          <span class="text-white">{{ room.description }}</span>
          <RouterLink :to="{name: 'join-room', params: {id: room.id}}" class="text-white bg-sky-500 p-2 rounded-lg hover:bg-sky-600 transition-colors hover:cursor-pointer">
            Join
          </RouterLink>
        </div>
      </div>

      <router-link to="/create/room" class="p-2 hover:cursor-pointer rounded-lg shadow-lg flex items-center justify-center bg-sky-500 hover:bg-sky-600 transition-colors group mt-4">
       <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 text-white group-hover:rotate-90 transition-transform">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
       </svg>
      </router-link>

      <div class="relative w-120 mt-4">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-sky-500/100"></div>
        </div>
        <div class="relative flex justify-center">
          <span class="px-2 text-white bg-slate-950 text-sm">
            Public Rooms
          </span>
        </div>
      </div>

      <div class="flex flex-col items-center mt-4 space-y-4">
        <div v-for="room in publicRooms" :key="room.id" class="w-150 bg-slate-900 p-4 rounded-2xl shadow-lg mt-4 flex items-center justify-between relative">
          <span class="absolute -top-5 text-white text-xs font-bold p-1 rounded-full transition-colors">{{ room.host_name }}</span>

          <span class="text-white">{{ room.name }}:</span>
          <span class="text-white">{{ room.description }}</span>
          <RouterLink :to="{name: 'join-room', params: {id: room.id}}" class="text-white bg-sky-500 p-2 rounded-lg hover:bg-sky-600 transition-colors hover:cursor-pointer">
            Join
          </RouterLink>
        </div>
      </div>

  </div>
</template>

<style scoped>


.slide-enter-active,
.slide-leave-active {
  transition: all 0.7s ease-in-out;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateY(100%);
  opacity: 0;
}


</style>
