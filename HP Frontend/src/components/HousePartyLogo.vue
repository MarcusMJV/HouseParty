<template>
   <div class="w-full flex justify-center">
      <h1 class="text-6xl font-bold text-sky-500/100 mt-6">
        House<span class="text-white">Party</span>
      </h1>
      <div class="flex items-center bg-slate-800 text-white p-4 rounded-2xl shadow-md space-x-4">
        <div class="flex space-x-1">
          <span class="w-1.5 h-6 bg-pink-500 animate-sound-bar"></span>
          <span class="w-1.5 h-8 bg-blue-500 animate-sound-bar delay-100"></span>
          <span class="w-1.5 h-5 bg-yellow-500 animate-sound-bar delay-200"></span>
          <span class="w-1.5 h-7 bg-green-500 animate-sound-bar delay-300"></span>
        </div>
      </div>
    </div>
</template>

<style scoped>
/* Sound bar animation */
@keyframes sound-bar {
  0%, 100% {
    transform: scaleY(1);
  }
  50% {
    transform: scaleY(1.5);
  }
}

.animate-sound-bar {
  animation: sound-bar 1s infinite ease-in-out;
}

/* Add a delay for each bar */
.delay-100 {
  animation-delay: 0.2s;
}
.delay-200 {
  animation-delay: 0.4s;
}
.delay-300 {
  animation-delay: 0.6s;
}




</style>
