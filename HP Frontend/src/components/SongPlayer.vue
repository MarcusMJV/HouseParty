<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { PropType } from 'vue'
import type { Song } from '@/types/song'

// Extend the global Window interface for TypeScript.
declare global {
  interface Window {
    onSpotifyWebPlaybackSDKReady: () => void
    Spotify: any
  }
}

// Define component props.
const props = defineProps({
  song: {
    type: Object as PropType<Song>,
    required: true,
  },
  apiToken: {
    type: String,
    required: true,
  },
})

// Local reactive references.
const currentSong = ref<Song>(props.song)
const token: string = "BQADq-OU4QvIBP0gLx1hWtiMXm1FQ04tP3zQKJ2cDVcvNQPezg9VIWMugzNtyVJR84Gkh0vYPhUvXfA4Ztevaz3fLqcDoyO9PA8l4Z_Y9MX89aKhiu5X8m0GSRylHNy7oMsgqo2qpY9xyQvpMXGYVcddzZi1R6co3ZUgzCL68szk9N2Xd632JemK3jTUzaObjYraa47sBbD3jgHzVfhsEODlLY5kRr64bRIqLSBzdkJ6zaAb253CBZKIgQON"
const player = ref<any>(null)
const deviceId = ref<string | null>(null)

/**
 * Loads the Spotify Web Playback SDK from the CDN.
 * Note: Define the onSpotifyWebPlaybackSDKReady callback before adding the script.
 */
function loadSpotifySDK(): Promise<void> {
  return new Promise((resolve, reject) => {
    // If the SDK is already loaded, resolve immediately.
    if (window.Spotify) {
      resolve()
      return
    }

    // Define the callback before the script is loaded.
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

/**
 * Initialize the Spotify player.
 */
async function initSpotifyPlayer() {
  await loadSpotifySDK()

  console.log("Token:" + token)

  const newPlayer = new window.Spotify.Player({
    name: 'House Party Player',
    getOAuthToken: (cb: (token: string) => void) => {
      cb(token)
    },
    volume: 0.5,
  })

  // Listen for when the player is ready.
  newPlayer.addListener('ready', ({ device_id }: any) => {
    console.log('Spotify Player Ready with Device ID:', device_id)
    deviceId.value = device_id
    playSong()
  })

  // Optionally listen for errors.
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
  
  await newPlayer.connect()
  player.value = newPlayer
}

// Initialize the Spotify player when the component mounts.
onMounted(() => {
  initSpotifyPlayer()
})

async function playSong() {
  if (deviceId.value && currentSong.value?.uri) {
      try {
        await fetch(`https://api.spotify.com/v1/me/player/play?device_id=${deviceId.value}`, {
          method: 'PUT',
          body: JSON.stringify({ uris: [currentSong.value.uri] }),
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
        })
        console.log('Playback started for song:', currentSong.value.uri)
      } catch (error) {
        console.error('Error starting playback:', error)
      }
    } else {
      console.warn('Device ID or song URI is not available yet.')
      console.log("uri: "+currentSong.value.uri)
      console.log("Device id: "+deviceId.value)
    }
}

/**
 * Watch for changes to the song prop and use the Spotify Web API to play the song.
 */
// watch(
//   () => props.song,
//   (song) => {

//     console.log('Song Set:', song)

//   },
// )
</script>

<template>
  <div>

  </div>
</template>
