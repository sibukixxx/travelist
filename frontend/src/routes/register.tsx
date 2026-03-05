import { createRoute } from '@tanstack/react-router'
import { rootRoute } from '../routeTree'
import { RegisterForm } from '../components/RegisterForm'

export const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  component: RegisterForm,
})
