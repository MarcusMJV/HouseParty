<script setup lang="ts">
import { ref, onMounted } from 'vue';
import HousePartyLogo from '@/components/HousePartyLogo.vue';
import SignUpFrom from '@/components/SignUpForm.vue';
import LoginForm from '@/components/LoginForm.vue';


const currentComponent = ref<'signup' | 'login'>('signup');
const isMounted = ref(false);

onMounted(() => {
  setTimeout(() => {
    currentComponent.value = 'signup';
    isMounted.value = true;
  }, 50);
});

</script>

<template>
  <div class="flex flex-col items-center pt-10 min-h-screen font-mono bg-slate-950 relative overflow-hidden">
    <HousePartyLogo/>

    <transition name="slide" mode="out-in" appear>
      <component
        v-if="isMounted"
        :is="currentComponent === 'signup' ? SignUpFrom : LoginForm"
        @switch-to-login="currentComponent = 'login'"
        @switch-to-signup="currentComponent = 'signup'"
        class="w-full max-w-md"
      />
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
