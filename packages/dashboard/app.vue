<template>
  <div class="p-8">
    <z-h1>Create a New Application</z-h1>
    <div class="flex gap-2 my-8">
      <z-input class="flex-1" placeholder="Git" v-model="url" />
      <z-button @click="clone" :is-loading="isLoading">Clone</z-button>
    </div>
    <pre class="w-full whitespace-pre-wrap border border-input bg-secondary p-4 rounded min-h-[64px] text-secondary-foreground">{{ log }}</pre>
    <pre class="mt-8 w-full whitespace-pre-wrap border border-input bg-secondary p-4 rounded min-h-[64px] text-secondary-foreground">{{ plan }}</pre>
  </div>
</template>

<script setup lang="ts">
import { createConnectTransport } from "@bufbuild/connect-web";
import { createPromiseClient } from "@bufbuild/connect";
import { AppService } from "proto/zyreva/v1/zyreva_connect";

const log = ref('');
const plan = ref('');
const url = ref('https://github.com/skyzh/prisma-edge-vercel');
const isLoading = ref(false);

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});
const client = createPromiseClient(AppService, transport);

const clone = async () => {
  isLoading.value = true;

  const stream = client.analyze({ url: url.value });
  for await (const chunk of stream) {
    log.value += chunk.log;
    if (chunk.finished) {
      plan.value = JSON.stringify(chunk, null, 2);
    }
  }

  isLoading.value = false;
}
</script>