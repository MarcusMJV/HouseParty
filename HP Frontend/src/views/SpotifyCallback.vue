<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import HousePartyLogo from '@/components/HousePartyLogo.vue'
import router from '@/router'

const route = useRoute()
const errorMessage = ref('')
const successMessage = ref('')
const userStore = useUserStore()
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL

onMounted(() => {
  const code = route.query.code as string

  if (code) {
    sendCodeToServer(code)
  }
})

const sendCodeToServer = async (code: string) => {
  try {
    const response = await fetch(`${apiBaseUrl}/spotify/token/callback/${code}`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${userStore.jwt}`,
      },
    })

    const data = await response.json()
    if (!response.ok) {
      throw new Error(data.error || 'Signup failed')
    }

    console.log(data)

    userStore.activateSpotifyConnection()

    successMessage.value = 'Spotify Connected! Redirecting....'

    router.push({ name: 'home' })
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'An unexpected error occurred'
  }
}
</script>

<template>
  <div
    class="flex flex-col items-center pt-10 min-h-screen font-mono bg-slate-950 relative overflow-hidden"
  >
    <HousePartyLogo />

    <div v-if="errorMessage" class="mt-4 p-3 bg-red-500/20 text-red-300 rounded-lg">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="mt-4 p-3 bg-green-500/20 text-green-300 rounded-lg">
      {{ successMessage }}
    </div>
  </div>
</template>

<style scoped></style>
