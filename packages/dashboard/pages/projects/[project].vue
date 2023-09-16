<template>
  {{ project }}
</template>

<script setup lang="ts">
const project = useRoute().params.project;
</script>