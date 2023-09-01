import { createTRPCNuxtClient, httpBatchLink } from 'trpc-nuxt/client'
import type { Router } from '@zyreva/server/src'

export default defineNuxtPlugin(() => {
  /**
   * createTRPCNuxtClient adds a `useQuery` composable
   * built on top of `useAsyncData`.
   */
  const client = createTRPCNuxtClient<Router>({
    links: [
      httpBatchLink({
        url: 'http://localhost:8080/trpc',
      }),
    ],
  })

  return {
    provide: {
      client,
    },
  }
})