import { describe, it, expect, vi, beforeEach } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { HomePage } from './HomePage'
import { renderWithProviders } from '../test/helpers'

vi.mock('../api/client', () => ({
  generatePlan: vi.fn(),
}))

import { generatePlan } from '../api/client'

const mockedGeneratePlan = vi.mocked(generatePlan)

async function submitFormWithDefaults(user: ReturnType<typeof userEvent.setup>) {
  await user.type(screen.getByLabelText('目的地'), '京都')
  await user.type(screen.getByLabelText('開始日'), '2026-04-01')
  await user.click(screen.getByRole('button', { name: 'プランを生成' }))
}

describe('HomePage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders PlanForm without crashing', () => {
    renderWithProviders(<HomePage />)
    expect(screen.getByRole('button', { name: 'プランを生成' })).toBeInTheDocument()
  })

  it('does not show result action buttons when there is no result', () => {
    renderWithProviders(<HomePage />)
    expect(screen.queryByRole('button', { name: 'JSONダウンロード' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'プランを再生成' })).not.toBeInTheDocument()
  })

  it('does not show result section before form submission', () => {
    renderWithProviders(<HomePage />)
    expect(screen.queryByRole('heading', { level: 2 })).not.toBeInTheDocument()
  })

  describe('after successful plan generation', () => {
    const mockResult = {
      itinerary: {
        id: '1',
        title: '京都3日間の旅',
        destination: '京都',
        start_date: '2026-04-01',
        end_date: '2026-04-03',
        days: [
          {
            day_number: 1,
            date: '2026-04-01',
            activities: [
              {
                order: 1,
                place_id: 'place_1',
                place: {
                  id: 'place_1',
                  google_place_id: 'gp_1',
                  name: '金閣寺',
                  lat: 35.0394,
                  lng: 135.7292,
                  types: ['temple'],
                  price_level: 1,
                  rating: 4.5,
                  address: '京都市北区',
                },
                start_time: '09:00',
                end_time: '10:30',
                duration_min: 90,
                estimated_cost_yen: 500,
                note: '朝一番がおすすめ',
              },
            ],
          },
        ],
        created_at: '2026-04-01T00:00:00Z',
        updated_at: '2026-04-01T00:00:00Z',
      },
      violations: [],
      budget_summary: {
        total_cost_yen: 500,
        daily_costs: [{ day_number: 1, cost_yen: 500 }],
      },
    }

    it('displays itinerary title after successful submission', async () => {
      const user = userEvent.setup()
      mockedGeneratePlan.mockResolvedValueOnce(mockResult)
      renderWithProviders(<HomePage />)

      await submitFormWithDefaults(user)

      expect(await screen.findByText('京都3日間の旅')).toBeInTheDocument()
    })

    it('displays day plan with activities', async () => {
      const user = userEvent.setup()
      mockedGeneratePlan.mockResolvedValueOnce(mockResult)
      renderWithProviders(<HomePage />)

      await submitFormWithDefaults(user)

      await waitFor(() => {
        expect(screen.getByText(/Day 1/)).toBeInTheDocument()
      })
      expect(screen.getByText('金閣寺')).toBeInTheDocument()
      expect(screen.getByText(/09:00 - 10:30/)).toBeInTheDocument()
      expect(screen.getByText('500円', { selector: '.timeline-cost' })).toBeInTheDocument()
      expect(screen.getByText(/朝一番がおすすめ/)).toBeInTheDocument()
    })

    it('displays budget summary', async () => {
      const user = userEvent.setup()
      mockedGeneratePlan.mockResolvedValueOnce(mockResult)
      renderWithProviders(<HomePage />)

      await submitFormWithDefaults(user)

      await waitFor(() => {
        expect(screen.getByText('500円', { selector: '.total-cost' })).toBeInTheDocument()
      })
    })
  })

  describe('violations display', () => {
    it('shows violations when returned from API', async () => {
      const user = userEvent.setup()
      mockedGeneratePlan.mockResolvedValueOnce({
        itinerary: {
          id: '1',
          title: '京都旅行',
          destination: '京都',
          start_date: '2026-04-01',
          end_date: '2026-04-03',
          days: [],
          created_at: '',
          updated_at: '',
        },
        violations: [
          {
            type: 'budget_exceeded',
            day_number: 0,
            activity_idx: 0,
            message: '予算を超過しています',
          },
          {
            type: 'too_many_activities',
            day_number: 1,
            activity_idx: 0,
            message: 'アクティビティが多すぎます',
          },
        ],
        budget_summary: { total_cost_yen: 80000, daily_costs: [] },
      })
      renderWithProviders(<HomePage />)

      await submitFormWithDefaults(user)

      expect(await screen.findByText('注意事項')).toBeInTheDocument()
      expect(screen.getByText('予算を超過しています')).toBeInTheDocument()
      expect(screen.getByText('アクティビティが多すぎます')).toBeInTheDocument()
    })

    it('does not show violations section when there are none', async () => {
      const user = userEvent.setup()
      mockedGeneratePlan.mockResolvedValueOnce({
        itinerary: {
          id: '1',
          title: '京都旅行',
          destination: '京都',
          start_date: '2026-04-01',
          end_date: '2026-04-03',
          days: [],
          created_at: '',
          updated_at: '',
        },
        violations: [],
        budget_summary: { total_cost_yen: 5000, daily_costs: [] },
      })
      renderWithProviders(<HomePage />)

      await submitFormWithDefaults(user)

      await waitFor(() => {
        expect(screen.getByText('京都旅行')).toBeInTheDocument()
      })
      expect(screen.queryByText('注意事項')).not.toBeInTheDocument()
    })
  })

  describe('error flow', () => {
    it('shows API error and allows retry', async () => {
      const user = userEvent.setup()
      mockedGeneratePlan.mockRejectedValueOnce(
        new Error('API error 503: Service Unavailable'),
      )
      renderWithProviders(<HomePage />)

      await submitFormWithDefaults(user)

      // Assert error is shown
      expect(
        await screen.findByText('エラー: API error 503: Service Unavailable'),
      ).toBeInTheDocument()
      // Assert no result section
      expect(screen.queryByRole('heading', { level: 2 })).not.toBeInTheDocument()

      // Retry with success
      const retryResult = {
        itinerary: {
          id: '2',
          title: 'リトライ成功',
          destination: '京都',
          start_date: '2026-04-01',
          end_date: '2026-04-03',
          days: [],
          created_at: '',
          updated_at: '',
        },
        violations: [],
        budget_summary: { total_cost_yen: 0, daily_costs: [] },
      }
      mockedGeneratePlan.mockResolvedValueOnce(retryResult)

      await user.click(screen.getByRole('button', { name: 'プランを生成' }))

      expect(await screen.findByText('リトライ成功')).toBeInTheDocument()
    })
  })
})
