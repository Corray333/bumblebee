<script setup lang="ts">
import { Toast, useToast } from 'primevue';
import NavBar from './components/NavBar.vue'
import { useSystemStore } from './stores/counter';
import { onMounted } from 'vue';

const systemStore = useSystemStore()

onMounted(()=>{
  systemStore.setToast(useToast())
})

</script>

<template>
  <Toast position="top-right" />
  
  <NavBar/>
  
  <RouterView />
</template>

<style scoped>



</style>
