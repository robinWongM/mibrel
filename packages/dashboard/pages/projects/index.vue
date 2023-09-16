<template>
  <div class="p-8">
    <z-h1>Projects</z-h1>
    <div v-if="data">
      <router-link v-for="{ id, name } in data" :key="id" :to="'/projects/' + id">
        <z-card>
          <z-card-header>{{ name || id }}</z-card-header>
          <z-card-content>
            <z-card-description>{{ name || id }}</z-card-description>
          </z-card-content>
        </z-card>
      </router-link>
    </div>
    <div v-else-if="isError">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query';
const { $client } = useNuxtApp();

const { data, isLoading, isError, error } = useQuery(['projects'], () => $client.query(['projects.list']));
</script>