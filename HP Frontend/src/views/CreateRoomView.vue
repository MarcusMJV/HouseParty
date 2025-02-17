<script setup lang="ts">

import HousePartyLogo from '@/components/HousePartyLogo.vue';
import CreateRoomForm from '@/components/CreateRoomForm.vue';
import { ref, onMounted } from 'vue';

const isMounted = ref(false);

onMounted(() => {
  setTimeout(() => {
    isMounted.value = true;
  }, 50);
});

</script>

<template>
  <div class="flex flex-col items-center pt-10 min-h-screen font-mono bg-slate-950 relative overflow-hidden">
    <HousePartyLogo/>

    <transition name="slide" mode="out-in" appear>
      <CreateRoomForm v-if="isMounted"/>
    </transition>


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
