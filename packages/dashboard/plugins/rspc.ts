import { FetchTransport, createClient } from '@rspc/client'
import type { Procedures } from '@zyreva/rspc'

export default defineNuxtPlugin(() => {
  const client = createClient<Procedures>({
    transport: new FetchTransport("http://localhost:3000/rspc"),
  })

  return {
    provide: {
      client,
    },
  }
})