<script setup lang="ts">
import type { ComponentPublicInstance } from 'vue'
import { Pressable } from '@ark-ui/vue'

const props = defineProps<{
  isLoading?: boolean;
  disabled?: boolean;
}>();

const isDisabled = computed(() => props.disabled || props.isLoading);

// Hack with @ark-ui/vue `disabled` bug

const buttonElement = ref<InstanceType<typeof Pressable> | null>(null);
const updateDisabledProp = (currentIsDisabled?: boolean) => {
  if (buttonElement.value) {
    ((buttonElement.value as ComponentPublicInstance).$el as HTMLButtonElement).disabled = currentIsDisabled ?? isDisabled.value;
  }
};

onMounted(updateDisabledProp);
watch(isDisabled, (currentIsDisabled) => {
  updateDisabledProp(currentIsDisabled);
});
</script>

<template>
  <Pressable
    ref="buttonElement"
    :disabled="isDisabled"
    class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2">
    <Icon v-if="isLoading" name="ph:circle-notch" class="mr-2 animate-spin" size="16" />
    <slot />
  </Pressable>
</template>