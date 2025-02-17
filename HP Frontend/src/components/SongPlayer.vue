<script setup lang="ts">
import type { Song } from '@/types/song';
import { type PropType, watch, ref } from 'vue';

const props = defineProps({
  song: {
    type: Object as PropType<Song>,
    required: true
  }
})

const currenntSong = ref<Song>()

watch(
  () => props.song,
  (song) => {
    currenntSong.value = song;
    console.log("Song Set")
  },
  { immediate: true }
)

</script>


<template>
  <div></div>
</template>
