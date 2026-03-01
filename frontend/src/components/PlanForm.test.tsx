import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { PlanForm } from './PlanForm'

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>,
  )
}

describe('PlanForm', () => {
  describe('start_date', () => {
    it('has min attribute set to today', () => {
      renderWithProviders(<PlanForm onResult={() => {}} />)

      const input = screen.getByLabelText('開始日')
      const today = new Date().toISOString().split('T')[0]

      expect(input).toHaveAttribute('min', today)
    })

    it('shows validation error for past date', async () => {
      renderWithProviders(<PlanForm onResult={() => {}} />)

      const dateInput = screen.getByLabelText('開始日')
      fireEvent.change(dateInput, { target: { value: '2020-01-01' } })

      // Fill required destination to avoid that validation error
      fireEvent.change(screen.getByLabelText('目的地'), { target: { value: '京都' } })

      fireEvent.submit(screen.getByRole('button', { name: 'プランを生成' }))

      await waitFor(() => {
        expect(screen.getByText('過去の日付は選択できません')).toBeInTheDocument()
      })
    })
  })

  describe('interests', () => {
    it('renders interests label without comma instruction', () => {
      renderWithProviders(<PlanForm onResult={() => {}} />)

      expect(screen.getByText('興味・関心')).toBeInTheDocument()
      expect(screen.queryByText(/カンマ区切り/)).not.toBeInTheDocument()
    })

    it('renders tag suggestions', () => {
      renderWithProviders(<PlanForm onResult={() => {}} />)

      expect(screen.getByRole('button', { name: '文化 を追加' })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: '歴史 を追加' })).toBeInTheDocument()
    })

    it('adds interest tag via suggestion click', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanForm onResult={() => {}} />)

      await user.click(screen.getByRole('button', { name: '食事 を追加' }))

      expect(screen.getByText('食事')).toBeInTheDocument()
    })
  })

  describe('form submission', () => {
    it('submits with required fields filled', async () => {
      const onResult = vi.fn()
      const user = userEvent.setup()
      renderWithProviders(<PlanForm onResult={onResult} />)

      await user.type(screen.getByLabelText('目的地'), '京都')

      const today = new Date().toISOString().split('T')[0]
      await user.type(screen.getByLabelText('開始日'), today)

      await user.click(screen.getByRole('button', { name: '文化 を追加' }))

      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      // Form should attempt submission (mutation will be triggered)
      // We verify no validation errors are shown
      await waitFor(() => {
        expect(screen.queryByText('目的地を入力してください')).not.toBeInTheDocument()
        expect(screen.queryByText('開始日を入力してください')).not.toBeInTheDocument()
      })
    })
  })
})
