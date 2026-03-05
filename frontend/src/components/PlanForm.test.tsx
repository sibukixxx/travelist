import { describe, it, expect, vi, beforeEach } from 'vitest'
import { screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { PlanForm } from './PlanForm'
import { renderWithProviders } from '../test/helpers'

vi.mock('../api/client', () => ({
  generatePlan: vi.fn(),
}))

import { generatePlan } from '../api/client'

const mockedGeneratePlan = vi.mocked(generatePlan)

describe('PlanForm', () => {
  const onResult = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('start_date', () => {
    it('has min attribute set to today', () => {
      renderWithProviders(<PlanForm onResult={onResult} />)

      const input = screen.getByLabelText('開始日')
      const today = new Date().toISOString().split('T')[0]

      expect(input).toHaveAttribute('min', today)
    })

    it('shows validation error for past date', async () => {
      renderWithProviders(<PlanForm onResult={onResult} />)

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
      renderWithProviders(<PlanForm onResult={onResult} />)

      expect(screen.getByText('興味・関心')).toBeInTheDocument()
      expect(screen.queryByText(/カンマ区切り/)).not.toBeInTheDocument()
    })

    it('renders tag suggestions', () => {
      renderWithProviders(<PlanForm onResult={onResult} />)

      expect(screen.getByRole('button', { name: '文化 を追加' })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: '歴史 を追加' })).toBeInTheDocument()
    })

    it('adds interest tag via suggestion click', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanForm onResult={onResult} />)

      await user.click(screen.getByRole('button', { name: '食事 を追加' }))

      expect(screen.getByText('食事')).toBeInTheDocument()
    })
  })

  describe('validation error display', () => {
    it('shows error when destination is empty on submit', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanForm onResult={onResult} />)

      // Act - submit without filling destination
      await user.type(screen.getByLabelText('開始日'), '2026-04-01')
      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      // Assert
      expect(await screen.findByText('目的地を入力してください')).toBeInTheDocument()
      expect(mockedGeneratePlan).not.toHaveBeenCalled()
    })

    it('shows error when start_date is empty on submit', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanForm onResult={onResult} />)

      // Act - submit without filling start_date
      await user.type(screen.getByLabelText('目的地'), '京都')
      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      // Assert
      expect(await screen.findByText('開始日を入力してください')).toBeInTheDocument()
      expect(mockedGeneratePlan).not.toHaveBeenCalled()
    })

    it('shows both errors when all required fields are empty on submit', async () => {
      const user = userEvent.setup()
      renderWithProviders(<PlanForm onResult={onResult} />)

      // Act - submit without filling any fields
      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      // Assert
      expect(await screen.findByText('目的地を入力してください')).toBeInTheDocument()
      expect(screen.getByText('開始日を入力してください')).toBeInTheDocument()
      expect(mockedGeneratePlan).not.toHaveBeenCalled()
    })
  })

  describe('API failure', () => {
    it('shows error message when API returns an error', async () => {
      const user = userEvent.setup()
      mockedGeneratePlan.mockRejectedValueOnce(
        new Error('API error 500: Internal Server Error'),
      )
      renderWithProviders(<PlanForm onResult={onResult} />)

      // Arrange - fill required fields
      await user.type(screen.getByLabelText('目的地'), '京都')
      await user.type(screen.getByLabelText('開始日'), '2026-04-01')

      // Act
      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      // Assert
      expect(
        await screen.findByText('エラー: API error 500: Internal Server Error'),
      ).toBeInTheDocument()
      expect(onResult).not.toHaveBeenCalled()
    })
  })

  describe('loading state', () => {
    it('disables button and shows loading text during submission', async () => {
      const user = userEvent.setup()
      let resolvePromise: (value: unknown) => void
      mockedGeneratePlan.mockImplementationOnce(
        () => new Promise((resolve) => { resolvePromise = resolve }),
      )
      renderWithProviders(<PlanForm onResult={onResult} />)

      // Arrange
      await user.type(screen.getByLabelText('目的地'), '京都')
      await user.type(screen.getByLabelText('開始日'), '2026-04-01')

      // Act
      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      // Assert - loading state
      const button = screen.getByRole('button', { name: 'プラン生成中...' })
      expect(button).toBeDisabled()

      // Cleanup - resolve the pending promise
      resolvePromise!({
        itinerary: { id: '1', title: 'テスト', destination: '京都', start_date: '2026-04-01', end_date: '2026-04-03', days: [], created_at: '', updated_at: '' },
        violations: [],
        budget_summary: { total_cost_yen: 0, daily_costs: [] },
      })

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'プランを生成' })).toBeEnabled()
      })
    })
  })

  describe('success state', () => {
    it('calls onResult with the API response on success', async () => {
      const user = userEvent.setup()
      const mockResult = {
        itinerary: {
          id: '1',
          title: '京都3日間の旅',
          destination: '京都',
          start_date: '2026-04-01',
          end_date: '2026-04-03',
          days: [],
          created_at: '2026-04-01T00:00:00Z',
          updated_at: '2026-04-01T00:00:00Z',
        },
        violations: [],
        budget_summary: { total_cost_yen: 15000, daily_costs: [] },
      }
      mockedGeneratePlan.mockResolvedValueOnce(mockResult)
      renderWithProviders(<PlanForm onResult={onResult} />)

      // Arrange
      await user.type(screen.getByLabelText('目的地'), '京都')
      await user.type(screen.getByLabelText('開始日'), '2026-04-01')

      // Act
      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      // Assert - onSuccess receives (data, variables, context) from react-query
      await waitFor(() => {
        expect(onResult).toHaveBeenCalled()
      })
      expect(onResult.mock.calls[0][0]).toEqual(mockResult)
    })

    it('sends correct request payload including optional fields', async () => {
      const user = userEvent.setup()
      mockedGeneratePlan.mockResolvedValueOnce({
        itinerary: { id: '1', title: 'テスト', destination: '京都', start_date: '2026-04-01', end_date: '2026-04-03', days: [], created_at: '', updated_at: '' },
        violations: [],
        budget_summary: { total_cost_yen: 0, daily_costs: [] },
      })
      renderWithProviders(<PlanForm onResult={onResult} />)

      // Arrange
      await user.type(screen.getByLabelText('目的地'), '京都')
      await user.type(screen.getByLabelText('開始日'), '2026-04-01')
      await user.clear(screen.getByLabelText('日数'))
      await user.type(screen.getByLabelText('日数'), '5')
      // Add interests via tag suggestion buttons
      await user.click(screen.getByRole('button', { name: '文化 を追加' }))
      await user.click(screen.getByRole('button', { name: '食事 を追加' }))
      await user.type(screen.getByLabelText(/予算上限/), '50000')
      await user.selectOptions(screen.getByLabelText('予算'), 'luxury')
      await user.selectOptions(screen.getByLabelText('旅行スタイル'), 'active')

      // Act
      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      // Assert
      await waitFor(() => {
        expect(mockedGeneratePlan).toHaveBeenCalled()
      })
      const payload = mockedGeneratePlan.mock.calls[0][0]
      expect(payload.destination).toBe('京都')
      expect(payload.num_days).toBe(5)
      expect(payload.start_date).toBe('2026-04-01')
      expect(payload.preferences.interests).toEqual(['文化', '食事'])
      expect(payload.preferences.budget).toBe('luxury')
      expect(payload.preferences.travel_style).toBe('active')
      expect(payload.preferences.total_budget_yen).toBe(50000)
      expect(payload.constraint).toEqual({
        max_walk_distance_m: 2000,
        max_activities_day: 6,
        earliest_start_time: '08:00',
        latest_end_time: '21:00',
      })
    })
  })
})
