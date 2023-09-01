<template>
  <div class="p-8">
    <z-h1>Create a New Application</z-h1>
    <div class="flex gap-2 my-8">
      <z-input class="flex-1" placeholder="Git" v-model="url" />
      <z-button @click="clone" :is-loading="isLoading">Clone</z-button>
    </div>
    <pre class="mt-8 w-full whitespace-pre-wrap border border-input bg-secondary p-4 rounded min-h-[64px] text-secondary-foreground">{{ plan }}</pre>
  </div>
</template>

<script setup lang="ts">
const plan = ref('');
const url = ref('https://github.com/skyzh/prisma-edge-vercel');
const isLoading = ref(false);

const { $client } = useNuxtApp()
 
const clone = () => {
  isLoading.value = true;
  $client.createApp.mutate({ gitUrl: url.value }).then((res) => {
    plan.value = JSON.stringify(res);
    isLoading.value = false;
  });
}
</script>