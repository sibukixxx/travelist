import { createRootRoute, createRouter } from '@tanstack/react-router'
import { RootLayout } from './components/RootLayout'

export const rootRoute = createRootRoute({
  component: RootLayout,
})

// Lazy import to avoid circular dependency
import { indexRoute } from './routes/index'

const routeTree = rootRoute.addChildren([indexRoute])

export const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}
