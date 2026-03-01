import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { HomePage } from './HomePage'

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>,
  )
}

describe('HomePage', () => {
  it('renders PlanForm without crashing', () => {
    renderWithProviders(<HomePage />)
    expect(screen.getByRole('button', { name: 'プランを生成' })).toBeInTheDocument()
  })
})
