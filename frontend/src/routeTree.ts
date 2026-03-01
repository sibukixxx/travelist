import { createRootRoute, createRouter } from '@tanstack/react-router'
import { RootLayout } from './components/RootLayout'

export const rootRoute = createRootRoute({
  component: RootLayout,
})

// Lazy import to avoid circular dependency
import { indexRoute } from './routes/index'
import { registerRoute } from './routes/register'

const routeTree = rootRoute.addChildren([indexRoute, registerRoute])

export const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}
