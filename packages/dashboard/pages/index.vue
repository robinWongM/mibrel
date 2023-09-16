<template>
  <div class="p-8">
    <z-h1>Create a New Application</z-h1>
    <div class="flex gap-2 my-8">
      <z-input class="flex-1" placeholder="Git" v-model="url" autocomplete="off" />
      <z-button @click="create" :is-loading="isLoading">Create</z-button>
    </div>

    <div>
      <z-card v-for="app in apps" :key="app.id">
        <z-card-header>{{ app.id }}</z-card-header>
        <z-card-content>
          <z-card-description>{{ app.git_url }}</z-card-description>
          <z-button @click="() => build(app.id)" :is-loading="isBuildLoading">Build</z-button>
          <z-button @click="() => deploy(app.id)" :is-loading="isBuildLoading">Deploy</z-button>
        </z-card-content>
      </z-card>
    </div>

    <pre
      class="mt-8 w-full whitespace-pre-wrap border border-input bg-secondary p-4 rounded min-h-[64px] text-secondary-foreground">{{ plan }}</pre>
  </div>
</template>

<script setup lang="ts">
import { Model } from '@mibrel/rspc';
import { WebsocketTransport, createClient } from '@rspc/client'
import type { Procedures } from '@mibrel/rspc'

const client = createClient<Procedures>({
  transport: new WebsocketTransport("ws://localhost:3000/rspc/ws"),
});

const plan = ref('');
const url = ref('https://github.com/skyzh/prisma-edge-vercel');
const isLoading = ref(false);
const isBuildLoading = ref(false);
const apps = ref<Model[]>([]);

const { $client } = useNuxtApp()

const refresh = () => {
  $client.query(['apps.list']).then((res) => {
    apps.value = res;
  });
}

const create = () => {
  isLoading.value = true;
  $client.mutation(['apps.create', { git_url: url.value }]).finally(() => {
    isLoading.value = false;
    refresh();
  });
}

const deploy = (id: number) => {
  $client.mutation(['apps.deploy', { id }]).then((res) => {
    plan.value = res;
  });
}

const build = (id: number) => {
  isBuildLoading.value = true;
  // @ts-ignore
  const unsubscribe = client.addSubscription(['apps.build', { id }], {
    onData: (data) => {
      plan.value += data + '\n';

      if (data === 'Build complete') {
        unsubscribe();
        isBuildLoading.value = false;
      }
    },
    onError: (err) => {
      console.error(err);
      isBuildLoading.value = false;
    }
  });
}

refresh();
</script>