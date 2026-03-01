import { createRoute } from '@tanstack/react-router'
import { rootRoute } from '../routeTree'
import { HomePage } from '../components/HomePage'

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomePage,
})
