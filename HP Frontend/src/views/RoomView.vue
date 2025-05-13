<script setup lang="ts">
// import HousePartyLogo from '@/components/HousePartyLogo.vue';
import { onMounted, ref, onBeforeUnmount, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import HousePartyLogo from '@/components/HousePartyLogo.vue'
import type { Song } from '@/types/spotify'
import { RefSymbol } from '@vue/reactivity'

const route = useRoute()
const socket = ref<WebSocket | null>(null)
const roomId = route.params.id as string
const user = useUserStore()
const messages = ref<string[]>([])
const usersCount = ref(1)
const currentSong = ref<Song | null>()
const queuedSongs = ref<Song[]>([])
const showSearchPanel = ref(false)
const searchQuery = ref('')
const isHost = ref(false)
const songPosition = ref<number>(0)
const searchResults = ref<Song[]>([])
const containerRef = ref<HTMLElement | null>(null)
const shouldScroll = ref(false)
let apiToken: string
const player = ref<any>()
const deviceId = ref<string | null>(null)
const remainingTime = ref<number>(0)
let timerInterval: ReturnType<typeof setInterval> | null = null
const skipCount = ref<number>(0)
const wsBaseUrl = import.meta.env.VITE_WS_BASE_URL

declare global {
  interface Window {
    onSpotifyWebPlaybackSDKReady: () => void
    Spotify: any
  }
}

watch(
  queuedSongs,
  () => {
    nextTick(checkHeight)
  },
  { deep: true },
)

const checkHeight = (): void => {
  if (containerRef.value) {
    const containerHeight = containerRef.value.scrollHeight
    shouldScroll.value = containerHeight > window.innerHeight * 0.42
  }
}

const toggleSearchPanel = () => {
  showSearchPanel.value = !showSearchPanel.value
}

onMounted(() => {
  joinRoom(roomId)
  checkHeight()
  window.addEventListener('resize', checkHeight)
})

const joinRoom = (roomId: string) => {
  console.log('Joining room', roomId)

  const wsUrl = `${wsBaseUrl}/join/room/${roomId}?token=${user.jwt}`

  socket.value = new WebSocket(wsUrl)

  socket.value.onopen = () => {
    console.log('WebSocket connection established')
    socket.value?.send(
      JSON.stringify({
        type: 'joined-room',
        payload: {
          from: user.credentials?.username,
        },
      }),
    )
  }

  socket.value.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)

      handleSocketMessage(message)
    } catch (error) {
      console.error('Error parsing message:', error)
    }
  }

  socket.value.onerror = (error) => {
    console.error('WebSocket error:', error)
  }

  socket.value.onclose = (event) => {
    console.log(`Connection closed: ${event.reason || 'Unknown reason'}`)
  }
}

const searchSongs = () => {
  const lower = <T extends string>(value: T) => value.toLowerCase()
  const query = lower(searchQuery.value)

  console.log('Searching for songs:', query)

  socket.value?.send(
    JSON.stringify({
      type: 'search-songs',
      payload: {
        search: query,
      },
    }),
  )
}

const skipSong = () => {
  socket.value?.send(
    JSON.stringify({
      type: 'skip-request',
      payload: {
        from: user.credentials?.username,
      },
    }),
  )
}

const leaveRoom = () => {
  socket.value?.send(
    JSON.stringify({
      type: 'user-left',
      payload: {
        from: user.credentials?.username,
      },
    }),
  )
}

const addSong = (songId: string) => {
  socket.value?.send(
    JSON.stringify({
      type: 'add-song',
      payload: {
        from: user.credentials?.username,
        song_id: songId,
      },
    }),
  )
  toggleSearchPanel()
}

const handleSocketMessage = (message: any) => {
  console.log(message)
  switch (message.type) {
    case 'user-left':
      usersCount.value -= 1
      break
    case 'search-songs':
      searchResults.value = message.payload.songs
      break

    case 'final-song-ended':
      currentSong.value = null
      player.value.pause().then(() => {
        console.log('Paused!')
      })
      if (timerInterval) {
        clearInterval(timerInterval)
        timerInterval = null
      }
      break

    case 'skip-request':
      skipCount.value += 1
      break

    case 'set-and-play-song':
      const incomingSong = message.payload.song as Song
      queuedSongs.value = queuedSongs.value.filter((song) => song.id !== incomingSong.id)
      currentSong.value = incomingSong
      songPosition.value = 0

      skipCount.value = 0

      if (currentSong.value?.duration_ms) {
        remainingTime.value = currentSong.value.duration_ms
        startCountdownTimer()
      }

      if (isHost.value) {
        playSong()
      }

      break

    case 'added-song-playlist':
      queuedSongs.value.push(message.payload.song)
      break

    case 'room-information':
      if (message.payload.host_id === user.credentials.id) {
        apiToken = message.payload.api_token
        isHost.value = true
        initSpotifyPlayer()
      }

      if (message.payload.current_song?.uri != '') {
        currentSong.value = message.payload.current_song

        if (currentSong.value?.duration_ms) {
          remainingTime.value = currentSong.value.duration_ms - message.payload.song_position
          startCountdownTimer()
        }
      }
      if (message.payload.playlist?.length > 0) {
        queuedSongs.value = message.payload.playlist
      }
      usersCount.value = message.payload.user_count
      songPosition.value = message.payload.song_position

      break

    case 'joined-room':
      console.log(message.payload)
      messages.value.push(`${message.payload.from} joined the room`)
      usersCount.value += 1

      break

    case 'error':
      break
    default:
      console.warn('Unknown message type:', message.type)
  }
}

const startCountdownTimer = () => {
  if (timerInterval) {
    clearInterval(timerInterval)
  }

  timerInterval = setInterval(() => {
    if (remainingTime.value > 0) {
      remainingTime.value = Math.max(0, remainingTime.value - 1000)
    } else if (timerInterval) {
      clearInterval(timerInterval)
      timerInterval = null
    }
  }, 1000)
}

const formatDuration = (ms: number | undefined) => {
  if (!ms) return '0:00'
  const minutes = Math.floor(ms / 60000)
  const seconds = Math.floor((ms % 60000) / 1000)
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}

onBeforeUnmount(() => {
  if (socket.value) {
    leaveRoom()
    socket.value.close(1000, 'Component unmounted')
    socket.value = null
  }
  if (player.value) {
    player.value.disconnect()
    player.value = null
  }
  if (timerInterval) {
    clearInterval(timerInterval)
    timerInterval = null
  }

  window.removeEventListener('resize', checkHeight)
})

function loadSpotifySDK(): Promise<void> {
  return new Promise((resolve, reject) => {
    if (window.Spotify) {
      resolve()
      return
    }

    window.onSpotifyWebPlaybackSDKReady = () => {
      resolve()
    }
    const script = document.createElement('script')
    script.src = 'https://sdk.scdn.co/spotify-player.js'
    script.async = true
    script.onerror = reject
    document.head.appendChild(script)
  })
}

async function initSpotifyPlayer() {
  await loadSpotifySDK()
  const token: string = apiToken

  const newPlayer = new window.Spotify.Player({
    name: 'House Party Player',
    getOAuthToken: (cb: (token: string) => void) => {
      cb(token)
    },
    volume: 0.5,
  })

  newPlayer.addListener('ready', async ({ device_id }: any) => {
    deviceId.value = device_id

    await fetch('https://api.spotify.com/v1/me/player', {
      method: 'PUT',
      headers: {
        Authorization: 'Bearer ' + token,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        device_ids: [device_id],
        play: false, // Don’t auto-play yet
      }),
    })

    if (currentSong.value?.uri) {
      playSong()
    }
  })

  newPlayer.addListener('initialization_error', ({ message }: any) =>
    console.error('Initialization Error:', message),
  )
  newPlayer.addListener('authentication_error', ({ message }: any) =>
    console.error('Authentication Error:', message),
  )
  newPlayer.addListener('account_error', ({ message }: any) =>
    console.error('Account Error:', message),
  )
  newPlayer.addListener('playback_error', ({ message }: any) =>
    console.error('Playback Error:', message),
  )

  const connected = await newPlayer.connect()
  if (!connected) {
    console.error('Failed to connect the Spotify player')
  }

  player.value = newPlayer
  console.log('Player has been initialized')
}

const playSong = async () => {
  if (deviceId.value && currentSong.value?.uri) {
    await fetch(`https://api.spotify.com/v1/me/player/play?device_id=${deviceId.value}`, {
      method: 'PUT',
      body: JSON.stringify({ uris: [currentSong.value?.uri], position_ms: songPosition.value }),
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${apiToken}`,
      },
    })
  } else {
    console.warn('Device ID or song URI is not available yet.')
    console.log('uri: ' + currentSong.value?.uri)
    console.log('Device id: ' + deviceId.value)
  }
}
</script>

<template>
  <div
    class="flex flex-col items-center pt-10 min-h-screen font-mono bg-slate-950 relative overflow-hidden"
  >
    <HousePartyLogo />
    <p class="text-white">{{ usersCount }} Users Connnected To Room</p>

    <div class="relative w-120 mt-8">
      <div class="absolute inset-0 flex items-center">
        <div class="w-full border-t border-sky-500/100"></div>
      </div>
      <div class="relative flex justify-center">
        <span class="px-2 text-white bg-slate-950 text-sm"> Currently Playing </span>
      </div>
    </div>

    <div class="flex flex-col items-center mt-4 space-y-4">
      <div
        v-if="currentSong"
        class="w-150 bg-slate-900 rounded-2xl items-center justify-between relative"
      >
        <div class="p-3 hover:bg-slate-800 rounded-lg cursor-pointer transition-colors">
          <div class="flex justify-between items-center">
            <div class="flex items-center gap-3">
              <img
                :src="currentSong?.image.url"
                :alt="currentSong?.name"
                class="w-12 h-12 rounded-md object-cover"
              />
              <div>
                <p class="text-white font-medium">{{ currentSong?.name }}</p>
                <p class="text-slate-400 text-sm">{{ currentSong?.artists.join(', ') }}</p>
              </div>
            </div>
            <p></p>
            <span class="text-white text-sm">
              {{ formatDuration(remainingTime) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <p v-if="currentSong?.name" class="mt-2">Skip Song Votes: {{ skipCount }}</p>

    <div class="relative w-120 mt-4">
      <div class="absolute inset-0 flex items-center">
        <div class="w-full border-t border-sky-500/100"></div>
      </div>
      <div class="relative flex justify-center">
        <span class="px-2 text-white bg-slate-950 text-sm"> Queue </span>
      </div>
    </div>

    <div
      ref="containerRef"
      class="flex flex-col items-center mt-4 space-y-4 custom-scroll"
      :class="{ 'max-h-[42vh] overflow-y-auto': shouldScroll }"
    >
      <div
        v-for="song in queuedSongs"
        :key="song.id"
        class="w-150 bg-slate-900 rounded-2xl items-center justify-between relative"
      >
        <div class="p-3 hover:bg-slate-800 rounded-lg transition-colors">
          <div class="flex justify-between items-center">
            <div class="flex items-center gap-3">
              <img
                :src="song?.image.url"
                :alt="song?.name"
                class="w-12 h-12 rounded-md object-cover"
              />
              <div>
                <p class="text-white font-medium">{{ song?.name }}</p>
                <p class="text-slate-400 text-sm">{{ song?.artists.join(', ') }}</p>
              </div>
            </div>
            <span class="text-slate-400 text-sm">
              {{ formatDuration(song?.duration_ms) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <div class="flex flex-row items-center justify-center space-x-4 mt-4">
      <button
        @click="toggleSearchPanel"
        class="p-2 hover:cursor-pointer rounded-lg shadow-lg flex items-center justify-center bg-sky-500 hover:bg-sky-600 transition-colors group"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-5 h-5 text-white group-hover:rotate-90 transition-transform"
        >
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
      </button>

      <button
        @click="skipSong"
        class="p-3 bg-sky-500/20 rounded-full hover:cursor-pointer hover:bg-sky-500/30 transition-colors"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 h-6"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M3 8.689c0-.864.933-1.406 1.683-.977l7.108 4.061a1.125 1.125 0 0 1 0 1.954l-7.108 4.061A1.125 1.125 0 0 1 3 16.811V8.69ZM12.75 8.689c0-.864.933-1.406 1.683-.977l7.108 4.061a1.125 1.125 0 0 1 0 1.954l-7.108 4.061a1.125 1.125 0 0 1-1.683-.977V8.69Z"
          />
        </svg>
      </button>
    </div>

    <div
      v-if="showSearchPanel"
      @click="toggleSearchPanel"
      class="z-20 fixed inset-0 bg-black/50 backdrop-blur-sm transition-opacity"
    ></div>

    <div
      class="z-30 fixed bottom-0 w-150 bg-slate-900 rounded-t-2xl shadow-xl transition-transform duration-300 transform"
      :class="showSearchPanel ? 'translate-y-0' : 'translate-y-full'"
    >
      <div class="p-4">
        <div class="flex items-center space-x-2">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search for songs..."
            class="w-full p-2 bg-slate-800 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
          />
          <button
            @click="searchSongs"
            class="p-2 hover:cursor-pointer hover:bg-sky-500/100 border text-white font-semibold rounded-lg hover:bg-sky-600 transition duration-300"
          >
            Search
          </button>
          <button
            @click="toggleSearchPanel"
            class="p-2 text-slate-400 hover:text-white hover:cursor-pointer"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div class="mt-4 space-y-2 max-h-60 overflow-y-auto custom-scroll">
          <div
            @click="addSong(song.id)"
            v-for="song in searchResults"
            :key="song.id"
            class="p-3 hover:bg-slate-800 rounded-lg cursor-pointer transition-colors"
          >
            <div class="flex justify-between items-center">
              <div class="flex items-center gap-3">
                <img
                  :src="song.image.url"
                  :alt="song.name"
                  class="w-12 h-12 rounded-md object-cover"
                />
                <div>
                  <p class="text-white font-medium">{{ song.name }}</p>
                  <p class="text-slate-400 text-sm">{{ song.artists.join(', ') }}</p>
                </div>
              </div>
              <span class="text-slate-400 text-sm">
                {{ formatDuration(song.duration_ms) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* WebKit browsers (Chrome, Safari) */
.custom-scroll::-webkit-scrollbar {
  width: 8px;
}

.custom-scroll::-webkit-scrollbar-track {
  background: rgba(15, 23, 42, 0.5); /* slate-900 with transparency */
  border-radius: 4px;
}

.custom-scroll::-webkit-scrollbar-thumb {
  background: #334155; /* slate-700 */
  border-radius: 4px;
  border: 1px solid #475569; /* slate-600 */
}

.custom-scroll::-webkit-scrollbar-thumb:hover {
  background: #475569; /* slate-600 */
}

/* Firefox */
.custom-scroll {
  scrollbar-width: thin;
  scrollbar-color: #334155 rgba(15, 23, 42, 0.5);
}

/* Optional: Smooth transition for hover effects */
.custom-scroll::-webkit-scrollbar-thumb {
  transition: background 0.2s ease-in-out;
}
</style>
